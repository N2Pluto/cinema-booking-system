package usecase

import (
	"context"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
)

type ListCinemasInput struct {
	StartDate *time.Time
	EndDate   *time.Time
	OrderBy   string // "asc" | "desc"  (default "desc")
	Page      int    // 1-based
	Limit     int
}

type ListCinemasResult struct {
	Data       []*entity.Cinema
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

type CreateShowtimeInput struct {
	MovieName    string
	TheaterNo    int
	StartTime    time.Time
	EndTime      time.Time
	Price        int
	Rows         int // e.g. 5 → rows A-E
	SeatsPerRow  int // e.g. 10 → seats 1-10
}

type CinemaUseCase interface {
	ListCinemas(ctx context.Context, in ListCinemasInput) (*ListCinemasResult, error)
	CreateShowtime(ctx context.Context, in CreateShowtimeInput) (*entity.Cinema, error)
	GetCinema(ctx context.Context, id string) (*entity.Cinema, error)
}
