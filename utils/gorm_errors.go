package utils

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var (
	// Data related errors
	ErrDuplicateEntry = errors.New("record already exists")
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidData    = errors.New("invalid data provided")
	ErrDatabaseError  = errors.New("database error")

	// Constraint errors
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrUniqueConstraint    = errors.New("unique constraint violation")
	ErrCheckConstraint     = errors.New("check constraint violation")
	ErrNotNullConstraint   = errors.New("not null constraint violation")

	// Transaction errors
	ErrTransactionConflict = errors.New("transaction conflict")
	ErrDeadlock            = errors.New("deadlock detected")
	ErrLockTimeout         = errors.New("lock timeout")

	// Connection errors
	ErrConnectionFailed   = errors.New("database connection failed")
	ErrConnectionTimeout  = errors.New("connection timeout")
	ErrTooManyConnections = errors.New("too many connections")

	// Query errors
	ErrInvalidQuery = errors.New("invalid query")
	ErrQueryTimeout = errors.New("query timeout")
	ErrRowScanError = errors.New("row scan error")
)

func HandleGormError(err error, entity string) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return fmt.Errorf("%s: %w", entity, ErrRecordNotFound)

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return fmt.Errorf("%s: %w", entity, ErrDuplicateEntry)

	case errors.Is(err, gorm.ErrInvalidData):
		return fmt.Errorf("%s: %w", entity, ErrInvalidData)

	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return fmt.Errorf("%s: %w", entity, ErrForeignKeyViolation)

	case errors.Is(err, gorm.ErrRegistered):
		return fmt.Errorf("%s: %w", entity, ErrDuplicateEntry)

	case errors.Is(err, gorm.ErrInvalidTransaction):
		return fmt.Errorf("%s: %w", entity, ErrTransactionConflict)

	case errors.Is(err, gorm.ErrNotImplemented):
		return fmt.Errorf("%s: operation not supported: %w", entity, ErrInvalidQuery)

	case errors.Is(err, gorm.ErrMissingWhereClause):
		return fmt.Errorf("%s: missing where clause: %w", entity, ErrInvalidQuery)

	case errors.Is(err, gorm.ErrUnsupportedRelation):
		return fmt.Errorf("%s: unsupported relation: %w", entity, ErrInvalidQuery)

	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return fmt.Errorf("%s: primary key required: %w", entity, ErrInvalidData)

	case errors.Is(err, gorm.ErrModelValueRequired):
		return fmt.Errorf("%s: model value required: %w", entity, ErrInvalidData)

	case errors.Is(err, gorm.ErrInvalidValueOfLength):
		return fmt.Errorf("%s: invalid value length: %w", entity, ErrInvalidData)

	case errors.Is(err, gorm.ErrInvalidDB):
		return fmt.Errorf("%s: %w", entity, ErrConnectionFailed)

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return fmt.Errorf("%s: %w", entity, ErrDuplicateEntry)

	default:
		// Check for specific error strings that might indicate other types of errors
		errStr := err.Error()
		switch {
		case containsAny(errStr, []string{"deadlock", "dead lock"}):
			return fmt.Errorf("%s: %w", entity, ErrDeadlock)

		case containsAny(errStr, []string{"lock timeout", "lock wait timeout"}):
			return fmt.Errorf("%s: %w", entity, ErrLockTimeout)

		case containsAny(errStr, []string{"too many connections"}):
			return fmt.Errorf("%s: %w", entity, ErrTooManyConnections)

		case containsAny(errStr, []string{"connection refused", "bad connection"}):
			return fmt.Errorf("%s: %w", entity, ErrConnectionFailed)

		case containsAny(errStr, []string{"context deadline exceeded", "query timeout"}):
			return fmt.Errorf("%s: %w", entity, ErrQueryTimeout)

		case containsAny(errStr, []string{"check constraint"}):
			return fmt.Errorf("%s: %w", entity, ErrCheckConstraint)

		case containsAny(errStr, []string{"not null constraint"}):
			return fmt.Errorf("%s: %w", entity, ErrNotNullConstraint)

		case containsAny(errStr, []string{"unique constraint", "UNIQUE constraint"}):
			return fmt.Errorf("%s: %w", entity, ErrUniqueConstraint)

		default:
			return fmt.Errorf("%s: %w: %v", entity, ErrDatabaseError, err)
		}
	}
}

// containsAny checks if the error string contains any of the given patterns
func containsAny(s string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(strings.ToLower(s), strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}
