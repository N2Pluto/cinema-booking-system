package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainusecase "github.com/n2pluto/cinema-booking-system/internal/domain/usecase"
)

type BookingHandler struct {
	bookingUC domainusecase.BookingUseCase
}

func NewBookingHandler(bookingUC domainusecase.BookingUseCase) *BookingHandler {
	return &BookingHandler{bookingUC: bookingUC}
}

// POST /api/booking/lock
func (h *BookingHandler) Lock(c *gin.Context) {
	var body struct {
		CinemaID    string   `json:"cinema_id"    binding:"required"`
		SeatNumbers []string `json:"seat_numbers" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	booking, err := h.bookingUC.LockSeats(c.Request.Context(), domainusecase.LockSeatsInput{
		UserID:      userID,
		CinemaID:    body.CinemaID,
		SeatNumbers: body.SeatNumbers,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, booking)
}

// POST /api/booking/confirm
func (h *BookingHandler) Confirm(c *gin.Context) {
	var body struct {
		BookingID string `json:"booking_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	booking, err := h.bookingUC.ConfirmBooking(c.Request.Context(), body.BookingID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, booking)
}

// POST /api/booking/cancel
func (h *BookingHandler) Cancel(c *gin.Context) {
	var body struct {
		BookingID string `json:"booking_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if err := h.bookingUC.CancelBooking(c.Request.Context(), body.BookingID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}

// GET /api/booking/me
func (h *BookingHandler) GetMine(c *gin.Context) {
	userID := c.GetString("user_id")
	bookings, err := h.bookingUC.GetMyBookings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

// GET /api/admin/bookings?movie=&date=YYYY-MM-DD&status=&page=&limit=
func (h *BookingHandler) ListAll(c *gin.Context) {
	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	result, err := h.bookingUC.ListBookings(c.Request.Context(), domainusecase.AdminListBookingsInput{
		Movie:  c.Query("movie"),
		Date:   c.Query("date"),
		Status: c.Query("status"),
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
