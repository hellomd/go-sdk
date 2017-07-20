package errors

const (
	unexpectedCode  = "unexpected_error"
	unexpectedMsg   = "An unexpected error ocurred"
	notFoundCode    = "entity_not_found"
	notFoundMessage = "The entity could not be found"
)

// HelloMD API's domain errors -
var (
	ErrInvalidFields = &validationError{}

	ErrUnexptectedError = &JSONError{Code: unexpectedCode, Message: unexpectedMsg}

	ErrNotFound = &JSONError{Code: notFoundCode, Message: notFoundMessage}
)
