package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func NewValidator() echo.Validator {
	validate := validator.New()
	return &CustomValidator{
		validator: validate,
	}
}
