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

func NewHttpValidator() HttpValidator {
	newValidator := validator.New(validator.WithRequiredStructEnabled())
	newValidator.RegisterValidation("string", isString)
	newValidator.RegisterValidation("notEmpty", isNotEmpty)

	return HttpValidator{
		validate: newValidator,
	}
}

func isString(fieldLevel validator.FieldLevel) bool {
	_, ok := fieldLevel.Field().Interface().(string)
	return ok
}

func isNotEmpty(fieldLevel validator.FieldLevel) bool {
	field := fieldLevel.Field()

	if field.Kind() != reflect.String {
		return false
	}

	return strings.TrimSpace(field.String()) != ""
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
			case "notEmpty":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must not be empty", field))
			}
		}

		return errorMessages
	}

	return []string{}
}
