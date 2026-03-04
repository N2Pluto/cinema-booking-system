package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type cinemaRepository struct {
	col *mongo.Collection
}

func NewCinemaRepository(col *mongo.Collection) repository.CinemaRepository {
	return &cinemaRepository{col: col}
}

func (r *cinemaRepository) List(ctx context.Context, f repository.CinemaFilter) ([]*entity.Cinema, int64, error) {
	filter := bson.D{}
	if f.StartDate != nil || f.EndDate != nil {
		dateFilter := bson.D{}
		if f.StartDate != nil {
			dateFilter = append(dateFilter, bson.E{Key: "$gte", Value: *f.StartDate})
		}
		if f.EndDate != nil {
			dateFilter = append(dateFilter, bson.E{Key: "$lte", Value: *f.EndDate})
		}
		filter = append(filter, bson.E{Key: "start_time", Value: dateFilter})
	}

	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sortDir := -1
	if f.OrderBy == "asc" {
		sortDir = 1
	}

	page := f.Page
	if page < 1 {
		page = 1
	}
	limit := f.Limit
	if limit < 1 {
		limit = 10
	}
	skip := int64((page - 1) * limit)

	opts := options.Find().
		SetSort(bson.D{{Key: "start_time", Value: sortDir}}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var cinemas []*entity.Cinema
	if err := cursor.All(ctx, &cinemas); err != nil {
		return nil, 0, err
	}
	return cinemas, total, nil
}

func (r *cinemaRepository) FindByID(ctx context.Context, id string) (*entity.Cinema, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid cinema id: %w", err)
	}
	var c entity.Cinema
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &c, err
}

func (r *cinemaRepository) Create(ctx context.Context, cinema *entity.Cinema) (*entity.Cinema, error) {
	cinema.ID = bson.NewObjectID()
	cinema.CreatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, cinema)
	if err != nil {
		return nil, err
	}
	return cinema, nil
}

func (r *cinemaRepository) LockSeats(ctx context.Context, cinemaID bson.ObjectID, seatNos []string, userID bson.ObjectID) error {
	// Atomically lock seats only when ALL requested seats are AVAILABLE.
	filter := bson.D{
		{Key: "_id", Value: cinemaID},
		{Key: "$expr", Value: bson.D{
			{Key: "$eq", Value: bson.A{
				bson.D{{Key: "$size", Value: bson.D{
					{Key: "$filter", Value: bson.D{
						{Key: "input", Value: "$seats"},
						{Key: "as", Value: "s"},
						{Key: "cond", Value: bson.D{
							{Key: "$and", Value: bson.A{
								bson.D{{Key: "$in", Value: bson.A{"$$s.seat_no", seatNos}}},
								bson.D{{Key: "$eq", Value: bson.A{"$$s.status", string(entity.SeatAvailable)}}},
							}},
						}},
					}},
				}}},
				len(seatNos),
			}},
		}},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "seats.$[elem].status", Value: string(entity.SeatLocked)},
			{Key: "seats.$[elem].user_id", Value: userID},
		}},
	}

	opts := options.UpdateOne().SetArrayFilters([]interface{}{
		bson.D{
			{Key: "elem.seat_no", Value: bson.D{{Key: "$in", Value: seatNos}}},
			{Key: "elem.status", Value: string(entity.SeatAvailable)},
		},
	})

	result, err := r.col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("seats not available or cinema not found")
	}
	return nil
}

func (r *cinemaRepository) ReleaseSeats(ctx context.Context, cinemaID bson.ObjectID, seatNos []string) error {
	filter := bson.D{{Key: "_id", Value: cinemaID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "seats.$[elem].status", Value: string(entity.SeatAvailable)},
		}},
		{Key: "$unset", Value: bson.D{
			{Key: "seats.$[elem].user_id", Value: ""},
		}},
	}
	opts := options.UpdateOne().SetArrayFilters([]interface{}{
		bson.D{{Key: "elem.seat_no", Value: bson.D{{Key: "$in", Value: seatNos}}}},
	})
	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *cinemaRepository) ConfirmSeats(ctx context.Context, cinemaID bson.ObjectID, seatNos []string) error {
	filter := bson.D{{Key: "_id", Value: cinemaID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "seats.$[elem].status", Value: string(entity.SeatBooked)},
		}},
	}
	opts := options.UpdateOne().SetArrayFilters([]interface{}{
		bson.D{{Key: "elem.seat_no", Value: bson.D{{Key: "$in", Value: seatNos}}}},
	})
	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}
