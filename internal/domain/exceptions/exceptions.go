package exceptions

import "net/http"

func NewInternalServerError(messages ...string) *ErrorType {
	return &ErrorType{
		StatusCode: http.StatusInternalServerError,
		Messages:   messages,
		Type:       "Internal Server Error",
	}
}

func NewValidationError(messages ...string) *ErrorType {
	return &ErrorType{
		StatusCode: http.StatusBadRequest,
		Messages:   messages,
		Type:       "Validation Error",
	}
}

func NewNotFoundError(messages ...string) *ErrorType {
	return &ErrorType{
		StatusCode: http.StatusNotFound,
		Messages:   messages,
		Type:       "Not Found Error",
	}
}

func NewNotesQueueError(messages ...string) *ErrorType {
	return &ErrorType{
		StatusCode: http.StatusInternalServerError,
		Messages:   messages,
		Type:       "Notes queue error",
	}
}
