package response

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Message string `json:"errors,omitempty"`
}

func NewError(
	w http.ResponseWriter, r *http.Request, log *slog.Logger, err error, errStatus int, message string,
) {
	log.Error(message, "error", err.Error())
	w.WriteHeader(errStatus)
	render.JSON(w, r, makeResponse(message))
}

func NewValidateError(
	w http.ResponseWriter, r *http.Request, log *slog.Logger, errStatus int, message string,
	err error,
) {
	var validateErr validator.ValidationErrors
	errors.As(err, &validateErr)

	log.Error(message, "error", err.Error())
	w.WriteHeader(errStatus)
	render.JSON(w, r, validationError(validateErr))
}

func makeResponse(message string) Response {
	return Response{
		Message: message,
	}
}

func validationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Message: strings.Join(errMsgs, ", "),
	}

}
