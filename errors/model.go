package errors

import (
	"encoding/json"
)

// JSONError -
type JSONError struct {
	Code    string                `json:"code"`
	Message string                `json:"message"`
	Errors  []ValidationJSONError `json:"errors,omitempty"`
}

func (jError *JSONError) Error() string {
	b, _ := json.Marshal(jError)
	return string(b)
}
