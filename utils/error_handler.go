// Package utils provides utility functions for the task manager application
package utils

import (
	transport "KaungHtetHein116/IVY-backend/api/tansport"
	"KaungHtetHein116/IVY-backend/pkg/constants"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ErrorResponse represents the structure of error responses sent to clients
// Message contains the error description and Data can hold additional error details
type ErrorResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CustomHTTPErrorHandler is a centralized error handler for the Echo framework
// It processes different types of errors and converts them into appropriate HTTP responses
func CustomHTTPErrorHandler(err error, c echo.Context) {
	// Prevent handling errors for already committed responses
	if c.Response().Committed {
		return
	}

	// Handle validation errors from the validator package
	// These occur when request payload validation fails
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		validationErrors, _ := FormatValidationErrors(err)

		transport.NewApiErrorResponse(c,
			http.StatusBadRequest, constants.ErrValidationFailed,
			validationErrors)

		return
	}

	// Handle database record not found errors
	// This covers both GORM's native not found error and custom not found error
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrRecordNotFound) {
		transport.NewApiErrorResponse(c,
			http.StatusNotFound, constants.ErrErrorRecordNotFound,
			nil)
		return
	}

	// Handle duplicate entry errors
	// This occurs when trying to create a record that violates unique constraints
	if errors.Is(err, ErrDuplicateEntry) {
		transport.NewApiErrorResponse(c,
			http.StatusBadRequest, constants.ErrDuplicatedData,
			nil)
		return
	}

	// Handle Echo's built-in HTTP errors
	if he, ok := err.(*echo.HTTPError); ok {
		transport.NewApiErrorResponse(c,
			he.Code, he.Message.(string),
			nil)
		return
	}

	// Handle all other unspecified errors as internal server errors
	transport.NewApiErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
}
