package reader

import (
	"context"
	"errors"
)

var ErrEmailAlreadyExsists = errors.New("reader with this email already exsists")

type Repository interface {
	Create(
		ctx context.Context,
		reader Reader,
	) (Reader, error)
}
