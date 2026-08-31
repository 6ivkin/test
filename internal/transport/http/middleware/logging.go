package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter

	status int
	bytes  int
}

func (w *responseWriter) WriteHeader(
	status int,
) {
	if w.status != 0 {
		return
	}

	w.status = status

	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(
	data []byte,
) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}

	n, err := w.ResponseWriter.Write(data)

	w.bytes += n

	return n, err
}

func Logging(
	logger *slog.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			start := time.Now()

			writer := &responseWriter{
				ResponseWriter: w,
			}

			next.ServeHTTP(writer, r)

			status := writer.status
			if status == 0 {
				status = http.StatusOK
			}

			logger.InfoContext(
				r.Context(),
				"http request",
				slog.String(
					"request_id",
					RequestIDFromContext(r.Context()),
				),
				slog.String(
					"method",
					r.Method,
				),
				slog.String(
					"path",
					r.URL.Path,
				),
				slog.Int(
					"status",
					status,
				),
				slog.Int(
					"bytes",
					writer.bytes,
				),
				slog.Duration(
					"duration",
					time.Since(start),
				),
			)
		})
	}
}
