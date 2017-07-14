package errors

import (
	"encoding/json"
)

type JSONError struct {
	Code        string                `json:"code"`
	Message     string                `json:"message"`
	Description string                `json:"description,omitempty"`
	Erros       []validationJSONError `json:"errors,omitempty"`
}

func (jError *JSONError) Error() string {
	b, _ := json.Marshal(jError)
	return string(b)
}
