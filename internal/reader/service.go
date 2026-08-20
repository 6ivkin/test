package reader

import (
	"context"
	"net/mail"
	"strings"

	"github.com/google/uuid"
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

func (s *Service) GetByID(ctx context.Context, id string) (Reader, error) {
	id = strings.TrimSpace(id)

	if err := validReaderID(id); err != nil {
		return Reader{}, err
	}

	return s.repository.GetByID(ctx, id)
}

func (s *Service) Deactivate(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)

	if err := validReaderID(id); err != nil {
		return err
	}

	return s.repository.Deactivate(ctx, id)
}

func validReaderID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidReaderID
	}

	return nil
}
