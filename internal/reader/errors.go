package reader

import "errors"

var (
	ErrInvalidFullName = errors.New("full name is required")
	ErrInvalidEmail    = errors.New("invalid email")

	ErrEmailAlreadyExists = errors.New(
		"reader with this email already exists",
	)

	ErrReaderNotFound  = errors.New("reader not found")
	ErrInvalidReaderID = errors.New("invalid reader id")
)
