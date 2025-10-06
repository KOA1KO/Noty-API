package httpx

import (
	httpCodes "Noty/internal/api/http_codes"
	"Noty/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func DecodeAndValidate(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var err error

	if err = decoder.Decode(dst); err != nil {
		httpCodes.WriteErrorResponse(w, http.StatusInternalServerError, service.ErrInternal.Error())
		return err
	}
	validate := validator.New()
	if err = validate.Struct(dst); err != nil {
		httpCodes.WriteErrorResponse(w, http.StatusBadRequest, service.ErrBadRequest.Error())
		return err
	}
	return nil
}
