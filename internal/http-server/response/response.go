package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func Ok() Response {
	return Response{Status: StatusOk}
}

func Error(err string) Response {
	return Response{Status: StatusError, Error: err}
}

func ValidateError(errors validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errors {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid Url", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Error:  fmt.Sprintf("validation error: %v", errMsgs),
	}
}
