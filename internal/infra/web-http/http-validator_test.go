package webhttp_test

import (
	"encoding/json"
	"testing"

	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/stretchr/testify/suite"
)

type HttpValidatorSuite struct {
	suite.Suite
}

func (h *HttpValidatorSuite) TestValidate_OnValidFields_ReturnsEmptyArray() {
	type Example struct {
		Field1 any `validate:"string"`
		Field2 any `validate:"integer"`
		Field3 any `validate:"notEmpty"`
	}
	var example Example
	err := json.Unmarshal([]byte(`
		{
			"field1": "abcdefg",
			"field2": 4,
			"field3": "abcdefg"
		}
	`), &example)
	h.Require().NoError(err)
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)

	errorMessages := validator.Validate(example)

	h.EqualValues([]string{}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnInvalidFieldWithTagString_ReturnErrors() {
	type Example struct {
		Field1 any `validate:"string"`
		Field2 any `validate:"string"`
		Field3 any `validate:"string"`
	}
	var example Example
	err := json.Unmarshal([]byte(`
		{
			"field1": 1,
			"field2": [],
			"field3": {}
		}
	`), &example)
	h.Require().NoError(err)
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field1 must be string", "field2 must be string", "field3 must be string"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnInvalidFieldWithTagInt_ReturnErrors() {
	type Example struct {
		Field1 any `validate:"integer"`
		Field2 any `validate:"integer"`
		Field3 any `validate:"integer"`
		Field4 any `validate:"integer"`
	}
	var example Example
	err := json.Unmarshal([]byte(`
		{
			"field1": 1,
			"field2": 1.5,
			"field3": -1,
			"field4": "1"
		}
	`), &example)
	h.Require().NoError(err)
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)

	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field2 must be integer", "field4 must be integer"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnInvalidFieldWithTagRequired_ReturnErrors() {
	type Example struct {
		Field1 any `validate:"required"`
		Field2 any `validate:"required"`
	}
	var example Example
	err := json.Unmarshal([]byte(`
		{
			"field2": null
		}
	`), &example)
	h.Require().NoError(err)
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field1 is required", "field2 is required"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnInvalidFieldWithTagEmpty_ReturnErrors() {
	type Example struct {
		Field1 any `validate:"notEmpty"`
		Field2 any `validate:"notEmpty"`
	}
	var example Example
	err := json.Unmarshal([]byte(`
		{
			"field1": "",
			"field2": " "
		}
	`), &example)
	h.Require().NoError(err)
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field1 must not be empty", "field2 must not be empty"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnInvalidFieldWithTagUuid4_ReturnErrors() {
	type Example struct {
		Field1 any `validate:"uuid4"`
	}
	var example Example
	err := json.Unmarshal([]byte(`
		{
			"field1": "abc"
		}
`), &example)
	h.Require().NoError(err)
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field1 must be uuidv4"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnInvalidFieldWithTagGte_ReturnErrors() {
	type Example struct {
		Field1 any `validate:"gte=8"`
	}
	var example Example
	err := json.Unmarshal([]byte(`
		{
			"field1": "abc"
		}
`), &example)
	h.Require().NoError(err)
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field1 must be greater than or equal to 8"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnInvalidFieldWithTagLt_ReturnErrors() {
	type Example struct {
		Field1 any `validate:"lt=4"`
	}
	var example Example
	err := json.Unmarshal([]byte(`
		{
			"field1": "abcdefgh"
		}
`), &example)
	h.Require().NoError(err)
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field1 must be less than 4"}, errorMessages)
}

func (h *HttpValidatorSuite) TestValidate_OnInvalidFieldWithTagPositive_ReturnErrors() {
	type Example struct {
		Field1 any `validate:"positive"`
		Field2 any `validate:"positive"`
		Field3 any `validate:"positive"`
	}
	var example Example
	err := json.Unmarshal([]byte(`
		{
			"field1": 0,
			"field2": -1,
			"field3": 1
		}
`), &example)
	h.Require().NoError(err)
	validator, err := webhttp.NewHttpValidator()
	h.Require().NoError(err)
	errorMessages := validator.Validate(example)

	h.EqualValues([]string{"field2 must be positive"}, errorMessages)
}

func TestHttpValidator(t *testing.T) {
	suite.Run(t, new(HttpValidatorSuite))
}
