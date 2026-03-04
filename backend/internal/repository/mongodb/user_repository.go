package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type userRepository struct {
	col *mongo.Collection
}

func NewUserRepository(col *mongo.Collection) *userRepository {
	return &userRepository{col: col}
}

func (r *userRepository) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	var user entity.User
	err := r.col.FindOne(ctx, bson.M{"google_id": googleID}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user entity.User
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) Upsert(ctx context.Context, user *entity.User) (*entity.User, error) {
	now := time.Now()

	filter := bson.M{"google_id": user.GoogleID}
	update := bson.M{
		"$set": bson.M{
			"email":        user.Email,
			"display_name": user.DisplayName,
			"updated_at":   now,
		},
		"$setOnInsert": bson.M{
			"google_id":  user.GoogleID,
			"role":       entity.RoleUser,
			"created_at": now,
		},
	}

	opts := options.UpdateOne().SetUpsert(true)
	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return r.FindByGoogleID(ctx, user.GoogleID)
}
