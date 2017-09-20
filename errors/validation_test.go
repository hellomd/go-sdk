package errors

import "testing"
import "encoding/json"
import "reflect"
import "fmt"

// BasicInfo -
type BasicInfo struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
}

func TestValidator(t *testing.T) {
	validator := NewValidator()

	info := BasicInfo{
		"Felix",
		"",
	}

	expctedError := JSONError{
		Code:    validationErrorCode,
		Message: validationErrorMsg,
		Errors: []ValidationJSONError{ValidationJSONError{
			"required",
			"last_name",
			"last_name is required",
		}},
	}

	err := validator.Struct(info)
	resp := JSONError{}
	err = json.Unmarshal([]byte(err.Error()), &resp)
	if err != nil {
		t.Error("Error message is not a json, got error: ", err.Error())
	}

	if !reflect.DeepEqual(expctedError, resp) {
		t.Errorf("Expected error to be %v, go %v", expctedError, resp)
	}
}

func TestIsValidationError(t *testing.T) {
	if !IsValidationError(new(validationError)) {
		t.Error("expected to be validation error")
	}

	if IsValidationError(fmt.Errorf("anything")) {
		t.Error("expected to not be validation error")
	}
}
