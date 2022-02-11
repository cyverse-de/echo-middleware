package params

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Define a single validator to do all of the validations for us.
var v = validator.New()

// ValidatedQueryParam extracts a query parameter, validates it, and returns the parameter value as a string. The
// validationTag parameter may contain any validation tag supported by github.com/go-playground/validator/v10.
func ValidatedQueryParam(ctx echo.Context, name, validationTag string) (string, error) {
	value := ctx.QueryParam(name)

	// Validate the value.
	if err := v.Var(value, validationTag); err != nil {
		return "", err
	}

	return value, nil
}

// ValidatedPathParam extracts a path parameter, validates it, and returns the parameter value as a string. The
// validationTag parameter may contain any validation tag supported by github.com/go-playground-validator/v10.
func ValidatedPathParam(ctx echo.Context, name, validationTag string) (string, error) {
	value := ctx.Param(name)

	// Validate the value.
	if err := v.Var(value, validationTag); err != nil {
		return "", err
	}

	return value, nil
}
