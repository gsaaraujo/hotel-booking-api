package webhttp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type HttpValidator struct {
	validate *validator.Validate
}

func NewHttpValidator() (HttpValidator, error) {
	newValidator := validator.New(validator.WithRequiredStructEnabled())
	err := newValidator.RegisterValidation("string", isString)

	if err != nil {
		return HttpValidator{}, err
	}

	err = newValidator.RegisterValidation("integer", isInteger)

	if err != nil {
		return HttpValidator{}, err
	}

	err = newValidator.RegisterValidation("notEmpty", isNotEmpty)

	if err != nil {
		return HttpValidator{}, err
	}

	err = newValidator.RegisterValidation("positive", isPositive)

	if err != nil {
		return HttpValidator{}, err
	}

	HttpValidator := HttpValidator{
		validate: newValidator,
	}

	return HttpValidator, nil
}

func isString(fieldLevel validator.FieldLevel) bool {
	return fieldLevel.Field().Kind() == reflect.String
}

func isInteger(fieldLevel validator.FieldLevel) bool {
	if fieldLevel.Field().Kind() != reflect.Float64 {
		return false
	}

	value := fieldLevel.Field().Float()
	return value == float64(int(value))
}

func isNotEmpty(fieldLevel validator.FieldLevel) bool {
	field := fieldLevel.Field()

	if field.Kind() != reflect.String {
		return false
	}

	return strings.TrimSpace(field.String()) != ""
}
func isPositive(fieldLevel validator.FieldLevel) bool {
	if fieldLevel.Field().Kind() != reflect.Float64 {
		return false
	}

	value := fieldLevel.Field().Float()
	return value >= 0
}

func (h *HttpValidator) Validate(body any) []string {
	err := h.validate.Struct(body)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := []string{}

		for _, validationError := range validationErrors {
			tag := validationError.Tag()
			param := validationError.Param()
			field := strings.ToLower(validationError.Field()[:1]) + validationError.Field()[1:]

			switch tag {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", field))
			case "uuid4":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be uuidv4", field))
			case "gte":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be greater than or equal to %s", field, param))
			case "lt":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be less than %s", field, param))
			case "string":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be string", field))
			case "integer":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be integer", field))
			case "notEmpty":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must not be empty", field))
			case "positive":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be positive", field))
			}
		}

		return errorMessages
	}

	return []string{}
}
