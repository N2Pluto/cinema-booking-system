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
	var bookings []*entity.Booking
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
	var bookings []*entity.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}
