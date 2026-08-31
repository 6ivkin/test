package middleware

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	api "github.com/6ivkin/test.git/internal/transport/http/api"
)

func Recover(
	logger *slog.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}

				logger.ErrorContext(
					r.Context(),
					"panic recovered",
					slog.String(
						"request_id",
						RequestIDFromContext(r.Context()),
					),
					slog.String(
						"panic",
						fmt.Sprint(recovered),
					),
					slog.String(
						"stack",
						string(debug.Stack()),
					),
				)

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				w.WriteHeader(
					http.StatusInternalServerError,
				)

				_ = json.NewEncoder(w).Encode(
					api.ErrorResponse{
						Error: "internal server error",
					},
				)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
