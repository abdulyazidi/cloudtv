package main

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var ErrValidation = errors.New("validation failed")

type SignupParams struct {
	Username        string `validate:"required,min=3,max=32,alphanum"`
	Email           string `validate:"required,email"`
	Password        string `validate:"required,min=8,max=128"`
	ConfirmPassword string `validate:"required,eqfield=Password"`
}

func main() {
	validate := validator.New(validator.WithRequiredStructEnabled())
	ins := SignupParams{
		Username:        "asd",
		Email:           "asd@asd.com",
		Password:        "1234567",
		ConfirmPassword: "12345678",
	}
	if err := validate.Struct(ins); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMessages []string
			for _, e := range validationErrors {
				errMessages = append(errMessages, formatValidationError(e))
			}
			fmt.Println(errMessages)

		}
	}
	fmt.Println("end")

}

func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", e.Field())
	case "eqfield":
		return fmt.Sprintf("%s must match %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s failed %s validation", e.Field(), e.Tag())
	}
}
