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

var validate = ValidateMD{validator.New()}

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}

// New -
func New() *ValidateMD {
	return &validate
}

// ValidateMD is HelloMD APIs validator
type ValidateMD struct {
	*validator.Validate
}

// Struct -
func (mv *ValidateMD) Struct(s interface{}) error {
	err := mv.Validate.Struct(s)
	if err != nil {
		ErrInvalidFields.errors = err.(validator.ValidationErrors)
		return ErrInvalidFields
	}
	return nil
}

type validationError struct {
	errors []validator.FieldError
}

type validationJSONError struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *validationError) Error() string {

	a := []validationJSONError{}

	for _, fe := range v.errors {
		a = append(a, validationJSONError{fe.Tag(), fe.Field(), fe.Field() + " is " + fe.Tag()})
	}

	b, _ := json.Marshal(JSONError{
		Code:    validationErrorCode,
		Message: validationErrorMsg,
		Erros:   a,
	})
	return string(b)

}
