package handlers

import (
	httpCodes "Noty/internal/api/http_codes"
	"Noty/internal/ctxKey"
	"Noty/internal/models"
	"Noty/pkg/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.api.http.middleware.jwtAuthMiddleware"

		log := h.Logger.With(
			sl.Operation(op),
		)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Info(`the "Authorization" header is missing`)
			httpCodes.WriteErrorResponse(w, http.StatusUnauthorized, `отсутствует заголовок "Authorization"`)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			log.Info("token error")
			httpCodes.WriteErrorResponse(w, http.StatusUnauthorized, "ошибка токена")
			return
		}

		claims := jwt.MapClaims{}
		err := h.TokenManager.Parse(token, claims)
		if err != nil {
			log.Error("parse token error", sl.Err(err))
			httpCodes.CheckOutAndWriteError(w, err)
			return
		}

		expTime, err := claims.GetExpirationTime()
		if err != nil {
			log.Error("cannot get expiration time", sl.Err(err))
			httpCodes.CheckOutAndWriteError(w, err)
			return
		}
		if time.Now().After(expTime.Time) {
			log.Info("invalid token")
			httpCodes.WriteErrorResponse(w, http.StatusUnauthorized, "Пользователь не авторизован")
			return
		}

		next.ServeHTTP(w, r.WithContext(ctxKey.WithUserID(r.Context(), claims)))
	})
}

func (h *Handler) GetMeta(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		limitStr := query.Get("limit")
		limit, _ := strconv.Atoi(limitStr)
		if limit < 1 {
			limit = 2
		}

		offsetStr := query.Get("offset")
		offset, _ := strconv.Atoi(offsetStr)
		if offset < 0 {
			offset = 0
		}

		sort := strings.ToUpper(query.Get("sort"))
		if sort != models.ASC && sort != models.DESC {
			sort = models.DESC
		}

		noteMeta := models.NoteMeta{
			Limit:  limit,
			Offset: offset,
			Sort:   sort,
		}

		next.ServeHTTP(w, r.WithContext(ctxKey.WithNoteMeta(r.Context(), noteMeta)))
	})
}

func (h *Handler) MWLog() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := h.Logger.With(
			slog.String("component", "middleware/logger"),
		)

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			// r.ProtoMajor = 1 => HTTP/1.1
			// r.ProtoMajor = 2 => HTTP/2

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
