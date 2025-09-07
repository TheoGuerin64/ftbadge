package ftvalidator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	validator *validator.Validate
}

func fieldErrorMessage(fieldError validator.FieldError) string {
	field := fieldError.Field()
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters", field, fieldError.Param())
	default:
		return fmt.Sprintf("Invalid value for %s", field)
	}
}

func (v *Validator) Validate(i any) error {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	var sb strings.Builder
	for i, fieldError := range validationErrors {
		if i > 0 {
			sb.WriteString(" and ")
		}
		sb.WriteString(fieldErrorMessage(fieldError))
	}
	message := sb.String()

	return echo.NewHTTPError(http.StatusBadRequest, message).SetInternal(err)
}

func New() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}
