package cinema

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
	domainusecase "github.com/n2pluto/cinema-booking-system/internal/domain/usecase"
)

var _ domainusecase.CinemaUseCase = (*UseCase)(nil)

type UseCase struct {
	cinemaRepo repository.CinemaRepository
}

func NewUseCase(cinemaRepo repository.CinemaRepository) *UseCase {
	return &UseCase{cinemaRepo: cinemaRepo}
}

func (uc *UseCase) ListCinemas(ctx context.Context, in domainusecase.ListCinemasInput) (*domainusecase.ListCinemasResult, error) {
	f := repository.CinemaFilter{
		StartDate: in.StartDate,
		EndDate:   in.EndDate,
		OrderBy:   in.OrderBy,
		Page:      in.Page,
		Limit:     in.Limit,
	}

	cinemas, total, err := uc.cinemaRepo.List(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("list cinemas: %w", err)
	}

	limit := in.Limit
	if limit < 1 {
		limit = 10
	}
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domainusecase.ListCinemasResult{
		Data:       cinemas,
		Total:      total,
		Page:       in.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (uc *UseCase) CreateShowtime(ctx context.Context, in domainusecase.CreateShowtimeInput) (*entity.Cinema, error) {
	seats := generateSeats(in.Rows, in.SeatsPerRow, in.Price)

	cinema := &entity.Cinema{
		MovieName: in.MovieName,
		TheaterNo: in.TheaterNo,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
		Price:     in.Price,
		Seats:     seats,
		CreatedAt: time.Now(),
	}

	return uc.cinemaRepo.Create(ctx, cinema)
}

func (uc *UseCase) GetCinema(ctx context.Context, id string) (*entity.Cinema, error) {
	c, err := uc.cinemaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("cinema not found")
	}
	return c, nil
}

// generateSeats builds a seat grid: rows A-Z, seats 1-N.
func generateSeats(rows, seatsPerRow, price int) []entity.Seat {
	seats := make([]entity.Seat, 0, rows*seatsPerRow)
	for r := 0; r < rows; r++ {
		rowLetter := string(rune('A' + r))
		for s := 1; s <= seatsPerRow; s++ {
			seats = append(seats, entity.Seat{
				SeatNo: fmt.Sprintf("%s%d", rowLetter, s),
				Status: entity.SeatAvailable,
				Price:  price,
			})
		}
	}
	return seats
}
