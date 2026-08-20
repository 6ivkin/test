package httptransport

import (
	"net/http"

	"github.com/6ivkin/test.git/internal/transport/http/handler"
)

func NewRouter(readerHandler *handler.ReaderHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /api/v1/readers",
		readerHandler.Create,
	)

	mux.HandleFunc(
		"GET /api/v1/readers/{id}",
		readerHandler.GetByID,
	)

	mux.HandleFunc(
		"DELETE /api/v1/readers/{id}",
		readerHandler.Deactivate,
	)

	return mux
}
