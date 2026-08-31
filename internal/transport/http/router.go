package httptransport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"

	api "github.com/6ivkin/test.git/internal/transport/http/api"
	"github.com/6ivkin/test.git/internal/transport/http/handler"
	httpmiddleware "github.com/6ivkin/test.git/internal/transport/http/middleware"
)

func NewRouter(
	readerHandler *handler.ReaderHandler,
	logger *slog.Logger,
) http.Handler {
	router := chi.NewRouter()

	// Middleware для ВСЕХ запросов.
	router.Use(
		httpmiddleware.RequestID,
	)

	router.Use(
		httpmiddleware.Logging(logger),
	)

	router.Use(
		httpmiddleware.Recover(logger),
	)

	// Swagger/OpenAPI документация.
	// На эти endpoints OpenAPI validator не распространяется.
	router.Get(
		"/docs",
		swaggerHandler,
	)

	router.Get(
		"/openapi.json",
		openAPIHandler,
	)

	// Загружаем OpenAPI specification.
	spec, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	validator :=
		oapimiddleware.OapiRequestValidatorWithOptions(
			spec,
			&oapimiddleware.Options{
				ErrorHandlerWithOpts: func(
					ctx context.Context,
					err error,
					w http.ResponseWriter,
					r *http.Request,
					opts oapimiddleware.ErrorHandlerOpts,
				) {
					w.Header().Set(
						"Content-Type",
						"application/json",
					)

					w.WriteHeader(
						opts.StatusCode,
					)

					_ = json.NewEncoder(w).Encode(
						api.ErrorResponse{
							Error: err.Error(),
						},
					)
				},
			},
		)

	strictHandler := api.NewStrictHandler(
		readerHandler,
		nil,
	)

	// Только API endpoints проходят OpenAPI validation.
	router.Group(func(r chi.Router) {
		r.Use(validator)

		api.HandlerFromMux(
			strictHandler,
			r,
		)
	})

	return router
}
