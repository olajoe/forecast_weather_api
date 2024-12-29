package validator

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	timeRFC3339Tag string = "rfc3339"
)

var rfc3339Validator validator.Func = func(fl validator.FieldLevel) bool {
	timeStr, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	_, err := time.Parse(time.RFC3339, timeStr)
	return err == nil
}

func NewValidator() *validator.Validate {
	validate := validator.New()

	if err := validate.RegisterValidation(timeRFC3339Tag, rfc3339Validator); err != nil {
		slog.Log(context.TODO(), slog.LevelError, fmt.Sprintf("cannot register custom RFC3339 validation function: %s", err.Error()))
	}

	return validate
}
