package repository

import (
	"context"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
)

type UserRepository interface {
	FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	Upsert(ctx context.Context, user *entity.User) (*entity.User, error)
}
