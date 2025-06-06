package repository

import (
	"github.com/google/uuid"
)

// QueryFilter holds common query parameters
type QueryFilter struct {
	Limit     int
	Offset    int
	SortBy    string
	SortOrder string
}

// BookingQueryFilter specific filter for booking queries
type BookingQueryFilter struct {
	QueryFilter
	UserID     uuid.UUID
	Status     string
	BookedDate string
}

func NewBookingFilter() *BookingQueryFilter {
	return &BookingQueryFilter{
		QueryFilter: QueryFilter{
			Limit:  10, // default limit
			Offset: 0,  // default offset
		},
	}
}
