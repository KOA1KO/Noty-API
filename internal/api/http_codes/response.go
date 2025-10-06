package httpCodes

import (
	"Noty/internal/ctxKey"
	"Noty/internal/service"
	"Noty/pkg/auth"
	"encoding/json"
	"errors"
	"net/http"
)

func ResponseError(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ResponseOK() Response {
	return Response{
		Status: StatusOK,
	}
}

const contentTypeJSON = "application/json"

func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-type", contentTypeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-type", contentTypeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ResponseError(message)); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

// compare errors by errors.Is and write it with WriteErrorResponse func
func CheckOutAndWriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrBadRequest):
		WriteErrorResponse(w, http.StatusBadRequest, service.ErrBadRequest.Error())

	case errors.Is(err, service.ErrInvalidCredentials):
		WriteErrorResponse(w, http.StatusBadRequest, service.ErrInvalidCredentials.Error())

	case errors.Is(err, service.ErrUserExists):
		WriteErrorResponse(w, http.StatusConflict, service.ErrUserExists.Error())

	case errors.Is(err, service.ErrUserNotFound):
		WriteErrorResponse(w, http.StatusNotFound, service.ErrUserNotFound.Error())

	case errors.Is(err, ctxKey.ErrUnauthorized) || errors.Is(err, ctxKey.ErrInvalidSubject):
		WriteErrorResponse(w, http.StatusUnauthorized, ctxKey.ErrUnauthorized.Error())

	case errors.Is(err, service.ErrNoteNotFound):
		WriteErrorResponse(w, http.StatusNotFound, service.ErrNoteNotFound.Error())

	case errors.Is(err, auth.ErrInvalidToken):
		WriteErrorResponse(w, http.StatusUnauthorized, auth.ErrInvalidToken.Error())

	case errors.Is(err, service.ErrForbidden):
		WriteErrorResponse(w, http.StatusForbidden, service.ErrForbidden.Error())

	default:
		WriteErrorResponse(w, http.StatusInternalServerError, "internal server error")
	}
}
