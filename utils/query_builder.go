package utils

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// QueryBuilder holds common query building functionality
type QueryBuilder struct {
	query *gorm.DB
}

// NewQueryBuilder creates a new QueryBuilder instance
func NewQueryBuilder(db *gorm.DB, ctx context.Context) *QueryBuilder {
	return &QueryBuilder{
		query: db.WithContext(ctx),
	}
}

// ApplyUUIDFilter applies a UUID filter if the ID string is valid
func (qb *QueryBuilder) ApplyUUIDFilter(field, idStr string) *QueryBuilder {
	if id, err := uuid.Parse(idStr); err == nil && id != uuid.Nil {
		qb.query = qb.query.Where(field+" = ?", id)
	}
	return qb
}

// ApplyStringFilters applies multiple string filters from a map
func (qb *QueryBuilder) ApplyStringFilters(filters map[string]string) *QueryBuilder {
	for field, value := range filters {
		if value != "" {
			qb.query = qb.query.Where(field+" = ?", value)
		}
	}
	return qb
}

// ApplyInFilter applies an IN filter for comma-separated values
func (qb *QueryBuilder) ApplyInFilter(field, value string) *QueryBuilder {
	if value != "" {
		values := strings.Split(value, ",")
		// Trim whitespace from each value
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}
		// Only apply filter if we have non-empty values
		if len(values) > 0 && values[0] != "" {
			qb.query = qb.query.Where(field+" IN ?", values)
		}
	}
	return qb
}

// ApplySorting applies sorting based on the provided parameters
func (qb *QueryBuilder) ApplySorting(sortBy, sortOrder string) *QueryBuilder {
	if sortBy != "" {
		order := "ASC"
		if sortOrder == "desc" {
			order = "DESC"
		}
		qb.query = qb.query.Order(sortBy + " " + order)
	}
	return qb
}

// ApplyPagination applies limit and offset for pagination
func (qb *QueryBuilder) ApplyPagination(limit, offset int) *QueryBuilder {
	qb.query = qb.query.Limit(limit).Offset(offset)
	return qb
}

// ApplyPreloads adds preload relations
func (qb *QueryBuilder) ApplyPreloads(preloads ...string) *QueryBuilder {
	for _, preload := range preloads {
		qb.query = qb.query.Preload(preload)
	}
	return qb
}

// Build returns the final query
func (qb *QueryBuilder) Build() *gorm.DB {
	return qb.query
}
