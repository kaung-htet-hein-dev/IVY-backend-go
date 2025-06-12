package handler

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/usecase"
	"KaungHtetHein116/IVY-backend/utils"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type BookingHandler struct {
	usecase usecase.BookingUsecase
}

func NewBookingHandler(u usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{usecase: u}
}

func (h *BookingHandler) CreateBooking(c echo.Context, req *request.CreateBookingRequest) error {
	userID := c.Get("user_id").(string)

	booking, err := h.usecase.CreateBooking(c.Request().Context(), userID, req)

	if err == utils.ErrUserHadBooking {
		return transport.NewApiErrorResponse(c, http.StatusConflict,
			"You already have a booking for this service at this time. If you wish to book, please cancel the first booking.",
			err)
	}

	if errors.Is(err, utils.ErrServiceNotFound) || errors.Is(err, utils.ErrCategoryNotFound) {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "Service or Branch not found", nil)
	}

	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to create booking", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusCreated, "Booking created successfully", booking)
}

func (h *BookingHandler) GetBookingByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid booking ID", err)
	}

	booking, err := h.usecase.GetBookingByID(c.Request().Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Booking not found", err)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to get booking", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Booking retrieved successfully", booking)
}

func (h *BookingHandler) GetAllBookings(c echo.Context) error {
	filter := params.NewBookingQueryParams()
	err := c.Bind(filter)

	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err)
	}

	bookings, err := h.usecase.GetAllBookings(c.Request().Context(), filter)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to get bookings", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Bookings retrieved successfully", bookings)
}

func (h *BookingHandler) GetUserBookings(c echo.Context) error {
	userID := c.Get("user_id").(string)

	bookings, err := h.usecase.GetUserBookings(c.Request().Context(), userID)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to get user bookings", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "User bookings retrieved successfully", bookings)
}

func (h *BookingHandler) UpdateBooking(c echo.Context, req *request.UpdateBookingRequest) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid booking ID", err)
	}

	booking, err := h.usecase.UpdateBooking(c.Request().Context(), id, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Booking not found", err)
		}

		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to update booking", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Booking updated successfully", booking)
}

func (h *BookingHandler) DeleteBooking(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid booking ID", err)
	}

	err = h.usecase.DeleteBooking(c.Request().Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Booking not found", err)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to delete booking", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusNoContent, "Booking deleted successfully", nil)
}

func (h *BookingHandler) GetAvailableSlots(c echo.Context) error {
	branchID := c.QueryParam("branch_id")
	bookedDate := c.QueryParam("booked_date")

	if branchID == "" || bookedDate == "" {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Branch ID and booked date are required", nil)
	}

	timeSlots, err := h.usecase.GetTimeSlotsByBranchIDAndDate(
		c.Request().Context(),
		uuid.MustParse(branchID),
		bookedDate,
	)

	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to get available slots", err)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Available slots retrieved successfully", timeSlots)
}
