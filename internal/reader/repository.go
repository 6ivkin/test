package reader

import (
	"context"
	"errors"
)

var (
	ErrEmailAlreadyExsists = errors.New("reader with this email already exsists")
	ErrReaderNotFound      = errors.New("reader not found")
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
}
