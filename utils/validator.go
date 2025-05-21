package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type CustomValidator struct {
	Validator *validator.Validate
}

func BindAndValidateDecorator[T any](fn func(echo.Context, *T) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		input := new(T)
		if err := c.Bind(input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
		}

		if err := c.Validate(input); err != nil {
			log.Println(err)
			return err // This will be caught by our CustomHTTPErrorHandler
		}

		return fn(c, input)
	}
}

func FormatValidationErrors(err error) ([]ValidationError, bool) {
	var verrs validator.ValidationErrors
	if ok := errors.As(err, &verrs); !ok {
		return nil, false
	}

	formatted := make([]ValidationError, 0, len(verrs))
	for _, fe := range verrs {
		formatted = append(formatted, ValidationError{
			Field:   (fe.Field()),
			Message: getValidationMessage(fe),
		})
	}
	return formatted, true
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", (fe.Field()))
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", (fe.Field()), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", (fe.Field()))
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}

	return nil
}
