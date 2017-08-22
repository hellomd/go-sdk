package errors

import (
	"encoding/json"
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

const (
	validationErrorCode = "invalid_entity"
	validationErrorMsg  = "Entity Validation Failed"
)

// NewValidator -
func NewValidator() *ValidateMD {
	validate := &ValidateMD{validator.New()}
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return validate
}

// ValidateMD is HelloMD APIs validator
type ValidateMD struct {
	*validator.Validate
}

// Struct -
func (mv *ValidateMD) Struct(s interface{}) error {
	err := mv.Validate.Struct(s)
	if err != nil {
		return &validationError{err.(validator.ValidationErrors)}
	}
	return nil
}

type validationError struct {
	errors []validator.FieldError
}

// ValidationJSONError -
type ValidationJSONError struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *validationError) Error() string {

	a := []ValidationJSONError{}

	for _, fe := range v.errors {
		a = append(a, ValidationJSONError{fe.Tag(), fe.Field(), fe.Field() + " is " + fe.Tag()})
	}

	b, _ := json.Marshal(JSONError{
		Code:    validationErrorCode,
		Message: validationErrorMsg,
		Errors:  a,
	})
	return string(b)

}
