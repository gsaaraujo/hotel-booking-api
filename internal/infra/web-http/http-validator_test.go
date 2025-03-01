package webhttp_test

import (
	"testing"

	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/stretchr/testify/suite"
)

type HttpValidatorSuite struct {
	suite.Suite
}

func (h *HttpValidatorSuite) TestValidate_OnTagsStringNotEmpty_ReturnsEmptyArray() {
	type Example struct {
		Field any `validate:"string,notEmpty"`
	}
	example := Example{
		Field: "abcdefg",
	}
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnTagRequired_ReturnsError() {
	type Example struct {
		Field1 any `validate:"required"`
		Field2 any `validate:"required"`
		Field3 any `validate:"required"`
	}
	example := Example{
		Field2: 0,
		Field3: nil,
	}
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field1 is required", "field3 is required"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnTagEmpty_ReturnsError() {
	type Example struct {
		Field1 any `validate:"notEmpty"`
		Field2 any `validate:"notEmpty"`
	}
	example := Example{
		Field1: "",
		Field2: " ",
	}
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field1 must not be empty", "field2 must not be empty"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnTagString_ReturnsError() {
	type Example struct {
		Field1 any `validate:"string"`
		Field2 any `validate:"string"`
	}
	example := Example{
		Field1: 1,
		Field2: []string{},
	}
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field1 must be string", "field2 must be string"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnTagUuid4_ReturnsError() {
	type Example struct {
		Field any `validate:"uuid4"`
	}
	example := Example{
		Field: "abc",
	}
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field must be uuidv4"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnTagGte_ReturnsError() {
	type Example struct {
		Field any `validate:"gte=8"`
	}
	example := Example{
		Field: "abc",
	}
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field must be greater than or equal to 8"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnTagLt_ReturnsError() {
	type Example struct {
		Field any `validate:"lt=4"`
	}
	example := Example{
		Field: "abcdefgh",
	}
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field must be less than 4"}, errorMessages)
}

func TestHttpValidator(t *testing.T) {
	suite.Run(t, new(HttpValidatorSuite))
}
