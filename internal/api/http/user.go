package handlers

import (
	httpCodes "Noty/internal/api/http_codes"
	"Noty/internal/api/httpx"
	"Noty/internal/service"
	"Noty/pkg/logger/sl"
	"log/slog"
	"net/http"
)

func (h *Handler) PostSignUp(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handlers.user.PostSignUp"

	ctx := r.Context()

	log := h.Logger.With(
		slog.String("op", op),
	)

	var req InputUserRequest
	if err := httpx.DecodeAndValidate(w, r, &req); err != nil {
		log.Error("Decode error", sl.Err(err))
		return
	}

	log.Debug("request body decoded", slog.Any("username", req.Username))

	input := service.SignInput{
		Username: req.Username,
		Password: req.Password,
	}

	response, err := h.UserService.SignIn(ctx, input)
	if err != nil {
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	log.Info("user added succesfully")
	httpCodes.WriteJSONResponse(w, http.StatusCreated, response)
}

func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handlers.user.PostLogin"

	ctx := r.Context()

	log := h.Logger.With(
		slog.String("op", op),
	)

	var req InputUserRequest
	if err := httpx.DecodeAndValidate(w, r, &req); err != nil {
		log.Error("Decode error", sl.Err(err))
		return
	}

	log.Debug("request body decoded", slog.Any("username", req.Username))

	input := service.SignInput{
		Username: req.Username,
		Password: req.Password,
	}

	response, err := h.UserService.Login(ctx, input)
	if err != nil {
		httpCodes.CheckOutAndWriteError(w, err)
		return
	}

	log.Info("user added succesfully")
	httpCodes.WriteJSONResponse(w, http.StatusOK, response)
}

// func (h *Handler) GetRefresh(w http.ResponseWriter, r *http.Request) {
// 	const op = "internal.handlers.user.GetRefresh"

// 	ctx := r.Context()

// 	log := h.Logger.With(
// 		slog.String("op", op),
// 	)

// 	var req InputRefresh
// 	if err := httpx.DecodeAndValidate(w, r, &req); err != nil {
// 		log.Error("Decode error", sl.Err(err))
// 		return
// 	}

// 	log.Info("request body decoded")

// 	input := service.RefreshInput{
// 		RefreshToken: req.RefreshToken,
// 	}

// 	response, err := h.UserService.RefreshJWT(ctx, input)
// 	if err != nil {
// 		httpCodes.CheckOutAndWriteError(w, err)
// 		return
// 	}

// 	log.Info("user added succesfully")
// 	httpCodes.WriteJSONResponse(w, http.StatusOK, response)
// }
