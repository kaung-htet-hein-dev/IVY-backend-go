package repository

import (
	"KaungHtetHein116/IVY-backend/api/v1/params"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"KaungHtetHein116/IVY-backend/utils"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *entity.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error)
	GetAll(ctx context.Context, filter *params.ServiceQueryParams) ([]entity.Booking, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error)
	Update(ctx context.Context, id uuid.UUID, updates interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	CheckSameUserBooking(ctx context.Context, userID uuid.UUID, bookedDate string, bookedTime string) error
	GetBookingTimeSlotByDateAndBranch(ctx context.Context, branchID uuid.UUID, bookedDate string) []string
	BuildQuery(ctx context.Context, params *params.ServiceQueryParams, preloads ...string) *gorm.DB
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *entity.Booking) error {
	// Check if referenced entities exist
	var service entity.Service
	var branch entity.Branch
	var user entity.User

	if err := r.db.WithContext(ctx).First(&service, "id = ?", booking.ServiceID).Error; err != nil {
		return utils.ErrServiceNotFound
	}
	if err := r.db.WithContext(ctx).First(&branch, "id = ?", booking.BranchID).Error; err != nil {
		return utils.ErrBranchNotFound
	}
	if err := r.db.WithContext(ctx).First(&user, "id = ?", booking.UserID).Error; err != nil {
		return utils.ErrUserNotFound
	}

	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error) {
	var booking entity.Booking
	err := r.db.WithContext(ctx).
		First(&booking, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetAll(ctx context.Context, filter *params.ServiceQueryParams) ([]entity.Booking, error) {
	var bookings []entity.Booking

	query := r.BuildQuery(ctx, filter, "Service", "Branch")
	err := query.Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error) {
	var bookings []entity.Booking
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) Update(ctx context.Context, id uuid.UUID, updates interface{}) error {
	return r.db.WithContext(ctx).Model(&entity.Booking{}).Where("id = ?", id).Updates(updates).Error
}

func (r *bookingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entity.Booking{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *bookingRepository) CheckSameUserBooking(ctx context.Context, userID uuid.UUID,
	bookedDate string, bookedTime string) error {

	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Booking{}).
		Where("user_id = ? AND booked_date = ? AND booked_time = ?", userID, bookedDate, bookedTime).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return utils.ErrUserHadBooking
	}

	return nil
}

func (r *bookingRepository) GetBookingTimeSlotByDateAndBranch(ctx context.Context,
	branchID uuid.UUID, bookedDate string) []string {

	timeSlots := make([]string, 0)

	var bookings []entity.Booking
	err := r.db.WithContext(ctx).
		Where("branch_id = ? AND booked_date = ?", branchID, bookedDate).
		Find(&bookings).Error
	if err != nil {
		return timeSlots
	}

	for _, booking := range bookings {
		timeSlots = append(timeSlots, booking.BookedTime)
	}

	return timeSlots
}

func (r *bookingRepository) BuildQuery(ctx context.Context, params *params.ServiceQueryParams, preloads ...string) *gorm.DB {
	builder := utils.NewQueryBuilder(r.db, ctx)

	// Apply UUID filters
	builder.ApplyUUIDFilter("user_id", params.UserID).
		ApplyUUIDFilter("branch_id", params.BranchID).
		ApplyUUIDFilter("category_id", params.CategoryID)

	// Apply string filters
	builder.ApplyStringFilters(map[string]string{
		"status":      params.Status,
		"booked_date": params.BookedDate,
		"booked_time": params.BookedTime,
	})

	// Apply sorting
	if params.SortBy != "" {
		builder.ApplySorting(params.SortBy, params.SortOrder)
	} else {
		builder.ApplySorting("created_at", "desc")
	}

	// Apply pagination and preloads
	builder.ApplyPagination(params.Limit, params.Offset).
		ApplyPreloads(preloads...)

	return builder.Build()
}
