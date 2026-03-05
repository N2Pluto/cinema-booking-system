package mongodb

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type bookingRepository struct {
	col *mongo.Collection
}

func NewBookingRepository(col *mongo.Collection) repository.BookingRepository {
	return &bookingRepository{col: col}
}

func (r *bookingRepository) Create(ctx context.Context, booking *entity.Booking) (*entity.Booking, error) {
	booking.ID = bson.NewObjectID()
	booking.CreatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (r *bookingRepository) FindByID(ctx context.Context, id string) (*entity.Booking, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid booking id: %w", err)
	}
	var b entity.Booking
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&b)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &b, err
}

func (r *bookingRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Booking, error) {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	cursor, err := r.col.Find(ctx, bson.M{"user_id": oid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	bookings := make([]*entity.Booking, 0)
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id bson.ObjectID, status entity.BookingStatus) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": string(status)}},
	)
	return err
}

func (r *bookingRepository) FindExpiredPending(ctx context.Context, before time.Time) ([]*entity.Booking, error) {
	filter := bson.M{
		"status":     string(entity.BookingPending),
		"created_at": bson.M{"$lt": before},
	}
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	bookings := make([]*entity.Booking, 0)
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) FindAll(ctx context.Context, f repository.AdminBookingFilter) (*repository.AdminBookingResult, error) {
	page := f.Page
	if page < 1 {
		page = 1
	}
	limit := f.Limit
	if limit < 1 {
		limit = 20
	}

	// ── Build aggregation pipeline ─────────────────────────────────────────
	pipeline := mongo.Pipeline{
		// Join cinema to get movie info
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "cinemas"},
			{Key: "localField", Value: "cinema_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "cinema"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$cinema"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		// Promote cinema fields to top level
		{{Key: "$addFields", Value: bson.D{
			{Key: "movie_name", Value: "$cinema.movie_name"},
			{Key: "theater_no", Value: "$cinema.theater_no"},
			{Key: "start_time", Value: "$cinema.start_time"},
			{Key: "end_time", Value: "$cinema.end_time"},
		}}},
	}

	// ── Dynamic $match stage ───────────────────────────────────────────────
	match := bson.D{}
	if f.Status != "" {
		match = append(match, bson.E{Key: "status", Value: f.Status})
	}
	if f.Date != "" {
		t, err := time.Parse("2006-01-02", f.Date)
		if err == nil {
			match = append(match, bson.E{Key: "created_at", Value: bson.D{
				{Key: "$gte", Value: t},
				{Key: "$lt", Value: t.Add(24 * time.Hour)},
			}})
		}
	}
	if f.Movie != "" {
		match = append(match, bson.E{Key: "movie_name", Value: bson.D{
			{Key: "$regex", Value: f.Movie},
			{Key: "$options", Value: "i"},
		}})
	}
	if len(match) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: match}})
	}

	// ── Count total (clone pipeline before adding sort/skip/limit) ─────────
	countPipeline := append(pipeline, bson.D{{Key: "$count", Value: "total"}})
	countCursor, err := r.col.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, fmt.Errorf("count bookings: %w", err)
	}
	defer countCursor.Close(ctx)
	var countResult []struct {
		Total int64 `bson:"total"`
	}
	_ = countCursor.All(ctx, &countResult)
	var total int64
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// ── Data with sort + pagination ────────────────────────────────────────
	pipeline = append(pipeline,
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		bson.D{{Key: "$skip", Value: int64((page - 1) * limit)}},
		bson.D{{Key: "$limit", Value: int64(limit)}},
	)
	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("find bookings: %w", err)
	}
	defer cursor.Close(ctx)

	data := make([]*repository.BookingWithCinema, 0)
	if err := cursor.All(ctx, &data); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	return &repository.AdminBookingResult{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
