package transport

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FieldValidationError struct {
	Namespace string      `json:"namespace"`
	Field     string      `json:"field"`
	Error     string      `json:"error"`
	Kind      string      `json:"kind"`
	Value     interface{} `json:"value"`
	Message   string      `json:"message"`
}

type PaginationResponse struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

func NewApiErrorResponse(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, echo.Map{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func NewApiSuccessResponse(c echo.Context, code int, message string, data interface{}, opts ...interface{}) error {
	var pagination PaginationResponse
	if len(opts) > 0 {
		if p, ok := opts[0].(PaginationResponse); ok {
			pagination = p
		}
	}

	return c.JSON(code, echo.Map{
		"code":       code,
		"message":    message,
		"data":       data,
		"pagination": pagination,
	})
}

func NewApiCreateSuccessResponse(c echo.Context, message string, data interface{}) error {
	formattedMessage := message
	if message == "" {
		formattedMessage = http.StatusText(http.StatusCreated)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"code":    http.StatusCreated,
		"message": formattedMessage,
		"data":    data,
	})
}

func NewBindingErrorResponse(c echo.Context, err *json.UnmarshalTypeError) error {
	return c.JSON(http.StatusBadRequest, echo.Map{
		"message": "Bad Request",
		"data": []FieldValidationError{{
			Namespace: err.Struct,
			Field:     err.Field,
			Error:     err.Error(),
			Kind:      err.Type.Kind().String(),
			Value:     err.Value,
			Message:   err.Error(),
		}},
	})
}
