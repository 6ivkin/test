package reader

import (
	"context"
	"errors"
)

var (
	ErrEmailAlreadyExsists = errors.New("reader with this email already exsists")
	ErrReaderNotFound      = errors.New("reader not found")
	ErrInvalidReaderID     = errors.New("invalid reader id")
)

type Repository interface {
	Create(
		ctx context.Context,
		reader Reader,
	) (Reader, error)

	GetByID(
		ctx context.Context,
		id string,
	) (Reader, error)

	Deactivate(
		ctx context.Context,
		id string,
	) error
}
