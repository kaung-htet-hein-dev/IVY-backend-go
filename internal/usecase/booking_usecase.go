package usecase

import (
	"context"
	"time"

	"KaungHtetHein116/IVY-backend/api/v1/request"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/internal/repository"

	"github.com/google/uuid"
)

type BookingUsecase interface {
	CreateBooking(ctx context.Context, userID uuid.UUID, req *request.CreateBookingRequest) (*entity.Booking, error)
	GetBookingByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error)
	GetAllBookings(ctx context.Context) ([]entity.Booking, error)
	GetUserBookings(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error)
	UpdateBooking(ctx context.Context, id uuid.UUID, req *request.UpdateBookingRequest) (*entity.Booking, error)
	DeleteBooking(ctx context.Context, id uuid.UUID) error
}

type bookingUsecase struct {
	repo repository.BookingRepository
}

func NewBookingUsecase(repo repository.BookingRepository) BookingUsecase {
	return &bookingUsecase{repo: repo}
}

func (u *bookingUsecase) CreateBooking(ctx context.Context, userID uuid.UUID, req *request.CreateBookingRequest) (*entity.Booking, error) {
	booking := &entity.Booking{
		ID:         uuid.New(),
		UserID:     userID,
		ServiceID:  req.ServiceID,
		BranchID:   req.BranchID,
		BookedDate: req.BookedDate,
		BookedTime: req.BookedTime,
		Status:     "PENDING",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := u.repo.Create(ctx, booking)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (u *bookingUsecase) GetBookingByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *bookingUsecase) GetAllBookings(ctx context.Context) ([]entity.Booking, error) {
	return u.repo.GetAll(ctx)
}

func (u *bookingUsecase) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error) {
	return u.repo.GetByUserID(ctx, userID)
}

func (u *bookingUsecase) UpdateBooking(ctx context.Context, id uuid.UUID, req *request.UpdateBookingRequest) (*entity.Booking, error) {
	booking, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Status != "" {
		booking.Status = req.Status
	}
	booking.UpdatedAt = time.Now()

	err = u.repo.Update(ctx, booking)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (u *bookingUsecase) DeleteBooking(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
