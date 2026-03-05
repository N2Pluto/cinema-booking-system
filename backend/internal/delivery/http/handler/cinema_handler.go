package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	domainusecase "github.com/n2pluto/cinema-booking-system/internal/domain/usecase"
)

type CinemaHandler struct {
	cinemaUC domainusecase.CinemaUseCase
}

func NewCinemaHandler(cinemaUC domainusecase.CinemaUseCase) *CinemaHandler {
	return &CinemaHandler{cinemaUC: cinemaUC}
}

// GET /api/cinema  or  GET /api/admin/cinema
func (h *CinemaHandler) List(c *gin.Context) {
	in := domainusecase.ListCinemasInput{
		OrderBy: c.DefaultQuery("orderBy", "desc"),
		Page:    queryInt(c, "page", 1),
		Limit:   queryInt(c, "limit", 10),
	}

	if s := c.Query("startDate"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err == nil {
			in.StartDate = &t
		}
	}
	if s := c.Query("endDate"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err == nil {
			in.EndDate = &t
		}
	}

	result, err := h.cinemaUC.ListCinemas(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /api/seats/:cinemaId
func (h *CinemaHandler) GetSeats(c *gin.Context) {
	cinema, err := h.cinemaUC.GetCinema(c.Request.Context(), c.Param("cinemaId"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cinema)
}

// POST /api/admin/cinema
func (h *CinemaHandler) CreateShowtime(c *gin.Context) {
	var body struct {
		MovieName   string    `json:"movie_name"    binding:"required"`
		TheaterNo   int       `json:"theater_no"    binding:"required"`
		StartTime   time.Time `json:"start_time"    binding:"required"`
		EndTime     time.Time `json:"end_time"      binding:"required"`
		Price       int       `json:"price"         binding:"required"`
		Rows        int       `json:"rows"          binding:"required,min=1,max=26"`
		SeatsPerRow int       `json:"seats_per_row" binding:"required,min=1,max=50"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cinema, err := h.cinemaUC.CreateShowtime(c.Request.Context(), domainusecase.CreateShowtimeInput{
		MovieName:   body.MovieName,
		TheaterNo:   body.TheaterNo,
		StartTime:   body.StartTime,
		EndTime:     body.EndTime,
		Price:       body.Price,
		Rows:        body.Rows,
		SeatsPerRow: body.SeatsPerRow,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cinema)
}

func queryInt(c *gin.Context, key string, def int) int {
	s := c.Query(key)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 {
		return def
	}
	return v
}
