package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/6ivkin/test.git/internal/reader"
	api "github.com/6ivkin/test.git/internal/transport/http/api"
)

var _ api.StrictServerInterface = (*ReaderHandler)(nil)

func toAPIReader(r reader.Reader) (api.Reader, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return api.Reader{}, fmt.Errorf("parse reader id: %w", err)
	}

	return api.Reader{
		Id:        openapi_types.UUID(id),
		FullName:  r.FullName,
		Email:     openapi_types.Email(r.Email),
		IsActive:  r.IsActive,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func (h *ReaderHandler) CreateReader(
	ctx context.Context,
	request api.CreateReaderRequestObject,
) (api.CreateReaderResponseObject, error) {
	if request.Body == nil {
		return api.CreateReader400JSONResponse{
			Error: "invalid request body",
		}, nil
	}

	result, err := h.service.Create(
		ctx,
		reader.CreateInput{
			FullName: request.Body.FullName,
			Email:    string(request.Body.Email),
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, reader.ErrInvalidFullName),
			errors.Is(err, reader.ErrInvalidEmail):

			return api.CreateReader400JSONResponse{
				Error: err.Error(),
			}, nil

		case errors.Is(err, reader.ErrEmailAlreadyExists):

			return api.CreateReader409JSONResponse{
				Error: err.Error(),
			}, nil

		default:
			return api.CreateReader500JSONResponse{
				Error: "internal server error",
			}, nil
		}
	}

	response, err := toAPIReader(result)
	if err != nil {
		return api.CreateReader500JSONResponse{
			Error: "internal server error",
		}, nil
	}

	return api.CreateReader201JSONResponse(response), nil
}

func (h *ReaderHandler) GetReaderByID(
	ctx context.Context,
	request api.GetReaderByIDRequestObject,
) (api.GetReaderByIDResponseObject, error) {
	result, err := h.service.GetByID(
		ctx,
		request.Id.String(),
	)

	if err != nil {
		switch {
		case errors.Is(err, reader.ErrInvalidReaderID):
			return api.GetReaderByID400JSONResponse{
				Error: err.Error(),
			}, nil

		case errors.Is(err, reader.ErrReaderNotFound):
			return api.GetReaderByID404JSONResponse{
				Error: err.Error(),
			}, nil

		default:
			return api.GetReaderByID500JSONResponse{
				Error: "internal server error",
			}, nil
		}
	}

	response, err := toAPIReader(result)
	if err != nil {
		return api.GetReaderByID500JSONResponse{
			Error: "internal server error",
		}, nil
	}

	return api.GetReaderByID200JSONResponse(response), nil
}

func (h *ReaderHandler) DeactivateReader(
	ctx context.Context,
	request api.DeactivateReaderRequestObject,
) (api.DeactivateReaderResponseObject, error) {
	err := h.service.Deactivate(
		ctx,
		request.Id.String(),
	)

	if err != nil {
		switch {
		case errors.Is(err, reader.ErrInvalidReaderID):
			return api.DeactivateReader400JSONResponse{
				Error: err.Error(),
			}, nil

		case errors.Is(err, reader.ErrReaderNotFound):
			return api.DeactivateReader404JSONResponse{
				Error: err.Error(),
			}, nil

		default:
			return api.DeactivateReader500JSONResponse{
				Error: "internal server error",
			}, nil
		}
	}

	return api.DeactivateReader204Response{}, nil
}
