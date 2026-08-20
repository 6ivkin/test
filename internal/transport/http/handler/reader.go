package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/6ivkin/test.git/internal/reader"
)

type ReaderHandler struct {
	service *reader.Service
}

func NewReaderHandler(
	service *reader.Service,
) *ReaderHandler {
	return &ReaderHandler{
		service: service,
	}
}

type createReaderRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *ReaderHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request createReaderRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeJSON(
			w,
			http.StatusBadRequest,
			errorResponse{
				Error: "invalid request body",
			},
		)

		return
	}

	result, err := h.service.Create(
		r.Context(),
		reader.CreateInput{
			FullName: request.FullName,
			Email:    request.Email,
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, reader.ErrInvalidFullName),
			errors.Is(err, reader.ErrInvalidEmail):

			writeJSON(
				w,
				http.StatusBadRequest,
				errorResponse{
					Error: err.Error(),
				},
			)

		case errors.Is(err, reader.ErrEmailAlreadyExsists):

			writeJSON(
				w,
				http.StatusConflict,
				errorResponse{
					Error: err.Error(),
				},
			)

		default:
			writeJSON(
				w,
				http.StatusInternalServerError,
				errorResponse{
					Error: "internal server error",
				},
			)
		}

		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		result,
	)
}

func (h *ReaderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, reader.ErrReaderNotFound):
			writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *ReaderHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.Deactivate(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, reader.ErrInvalidReaderID):
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		case errors.Is(err, reader.ErrReaderNotFound):
			writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}
