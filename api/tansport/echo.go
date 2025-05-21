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

func NewApiErrorResponse(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, echo.Map{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func NewApiSuccessResponse(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, echo.Map{
		"code":    code,
		"message": message,
		"data":    data,
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
