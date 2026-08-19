package reader

import (
	"context"
	"errors"
	"net/mail"
	"strings"
)

var (
	ErrInvalidFullName = errors.New("full name is required")
	ErrInvalidEmail    = errors.New("invalid email")
)

type CreateInput struct {
	FullName string
	Email    string
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Reader, error) {
	fullName := strings.TrimSpace(input.FullName)
	email := strings.TrimSpace(input.Email)

	if fullName == "" {
		return Reader{}, ErrInvalidFullName
	}

	if email == "" {
		return Reader{}, ErrInvalidEmail
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return Reader{}, ErrInvalidEmail
	}

	return s.repository.Create(
		ctx,
		Reader{
			FullName: fullName,
			Email:    email,
		},
	)
}
