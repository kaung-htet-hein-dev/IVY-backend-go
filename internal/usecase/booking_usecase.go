package usecase

import (
	"context"

	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/internal/repository"
	"KaungHtetHein116/IVY-backend/utils"

	"github.com/google/uuid"
)

type BookingUsecase interface {
	CreateBooking(ctx context.Context, userID string, req *request.CreateBookingRequest) (*entity.Booking, error)
	GetBookingByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error)
	GetAllBookings(ctx context.Context, filter *params.BookingQueryParams) ([]entity.Booking, *transport.PaginationResponse, error)
	GetUserBookings(ctx context.Context, userID string) ([]entity.Booking, error)
	UpdateBooking(ctx context.Context, id uuid.UUID, req *request.UpdateBookingRequest) (*entity.Booking, error)
	DeleteBooking(ctx context.Context, id uuid.UUID) error
	GetTimeSlotsByBranchIDAndDate(ctx context.Context, branchID uuid.UUID, bookedDate string) ([]Slot, error)
}

type bookingUsecase struct {
	repo repository.BookingRepository
}

func NewBookingUsecase(repo repository.BookingRepository) BookingUsecase {
	return &bookingUsecase{repo: repo}
}

func (u *bookingUsecase) CreateBooking(ctx context.Context, userID string, req *request.CreateBookingRequest) (*entity.Booking, error) {
	booking := &entity.Booking{
		ID:         uuid.New(),
		UserID:     userID,
		ServiceID:  req.ServiceID,
		BranchID:   req.BranchID,
		BookedDate: req.BookedDate,
		BookedTime: req.BookedTime,
		Note:       req.Note,
		Status:     "PENDING",
	}

	// Check if the user already has a booking for this service at this time
	err := u.repo.CheckSameUserBooking(ctx, userID, req.BookedDate, req.BookedTime)
	if err != nil {
		if err == utils.ErrUserHadBooking {
			return nil, utils.ErrUserHadBooking
		}
		return nil, err
	}

	err = u.repo.Create(ctx, booking)

	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (u *bookingUsecase) GetBookingByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *bookingUsecase) GetAllBookings(ctx context.Context, filter *params.BookingQueryParams) ([]entity.Booking, *transport.PaginationResponse, error) {
	return u.repo.GetAll(ctx, filter)
}

func (u *bookingUsecase) GetUserBookings(ctx context.Context, userID string) ([]entity.Booking, error) {
	return u.repo.GetByUserID(ctx, userID)
}

func (u *bookingUsecase) UpdateBooking(ctx context.Context, id uuid.UUID, req *request.UpdateBookingRequest) (*entity.Booking, error) {
	// Check if booking exists
	_, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := u.repo.Update(ctx, id, req); err != nil {
		return nil, err
	}

	// Get updated booking
	return u.repo.GetByID(ctx, id)
}

func (u *bookingUsecase) DeleteBooking(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func (u *bookingUsecase) GetTimeSlotsByBranchIDAndDate(ctx context.Context, branchID uuid.UUID, bookedDate string) ([]Slot, error) {
	takenTimeSlots := u.repo.GetBookingTimeSlotByDateAndBranch(ctx, branchID, bookedDate)

	return getAvailableTimeSlots(takenTimeSlots), nil
}

type Slot struct {
	Slot        string `json:"slot"`
	IsAvailable bool   `json:"is_available"`
}

func getAvailableTimeSlots(takenTimeSlots []string) []Slot {
	allSlots := map[string]int8{
		"09:00 AM": 0,
		"09:30 AM": 0,
		"10:00 AM": 0,
		"10:30 AM": 0,
		"11:00 AM": 0,
		"11:30 AM": 0,
		"12:00 PM": 0,
		"12:30 PM": 0,
		"01:00 PM": 0,
		"01:30 PM": 0,
		"02:00 PM": 0,
		"02:30 PM": 0,
		"03:00 PM": 0,
		"03:30 PM": 0,
		"04:00 PM": 0,
		"04:30 PM": 0,
		"05:00 PM": 0,
	}

	available := make([]Slot, 0)

	for _, slot := range takenTimeSlots {
		_, exist := allSlots[slot]
		if exist {
			allSlots[slot] = allSlots[slot] + 1
		}
	}

	for slot, count := range allSlots {
		isAvailable := false
		if count < 2 {
			isAvailable = true
		}
		available = append(available, Slot{
			Slot:        slot,
			IsAvailable: isAvailable,
		})
	}

	return available
}
