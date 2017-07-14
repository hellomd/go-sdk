package errors

import "testing"
import "encoding/json"
import "reflect"

// BasicInfo -
type BasicInfo struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
}

func TestValidator(t *testing.T) {
	validator := New()

	info := BasicInfo{
		"Felix",
		"",
	}

	expctedError := inError{
		"required",
		"last_name",
		"last_name is required",
	}

	err := validator.Struct(info)
	if err == ErrInvalidFields {
		resp := []inError{}
		err = json.Unmarshal([]byte(err.Error()), &resp)
		if err != nil {
			t.Error("Error message is not a json, got error: ", err.Error())
		}

		if !reflect.DeepEqual(expctedError, resp[0]) {
			t.Errorf("Expected error to be %v, go %v", expctedError, resp[0])
		}

	}

}
