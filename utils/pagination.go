package utils

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"context"

	"gorm.io/gorm"
)

// CalculatePagination creates pagination response from total count, limit and offset
func CalculatePagination(total int64, limit, offset int) *transport.PaginationResponse {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	page := (offset / limit) + 1
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &transport.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
		HasNext:    total > int64(offset+limit),
		HasPrev:    offset > 0,
	}
}

// CountAndPaginate counts records and returns pagination
func CountAndPaginate(ctx context.Context, db *gorm.DB, model interface{}, limit, offset int) (*transport.PaginationResponse, error) {
	var total int64
	if err := db.WithContext(ctx).Model(model).Count(&total).Error; err != nil {
		return nil, err
	}
	return CalculatePagination(total, limit, offset), nil
}
