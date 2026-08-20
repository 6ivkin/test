package reader

import (
	"context"
	"errors"
	"testing"
)

type repositoryMock struct {
	createFn func(
		ctx context.Context,
		reader Reader,
	) (Reader, error)

	getByIDFn func(
		ctx context.Context,
		id string,
	) (Reader, error)

	deactivateFn func(
		ctx context.Context,
		id string,
	) error
}

func (m *repositoryMock) Create(ctx context.Context, reader Reader) (Reader, error) {
	if m.createFn == nil {
		panic("Create was not excepted to be called")
	}

	return m.createFn(ctx, reader)
}

func (m *repositoryMock) GetByID(ctx context.Context, id string) (Reader, error) {
	if m.getByIDFn == nil {
		panic("GetByID was not excepted to be called")
	}

	return m.getByIDFn(ctx, id)
}

func (m *repositoryMock) Deactivate(ctx context.Context, id string) error {
	if m.deactivateFn == nil {
		panic("Deactivate was not excepted to be callled")
	}

	return m.deactivateFn(ctx, id)
}

func TestService_Deactivate_Success(t *testing.T) {
	const readerID = "d819891f-8d3f-4c2a-9d93-c71a6441aebe"

	repository := &repositoryMock{
		deactivateFn: func(ctx context.Context, id string) error {
			if id != readerID {
				t.Fatalf("id = %q, want %q", id, readerID)
			}
			return nil
		},
	}

	service := NewService(repository)

	err := service.Deactivate(context.Background(), readerID)
	if err != nil {
		t.Fatalf("Deactivate() error = %v, want nil", err)
	}
}

func TestService_Deactivate_InvalidID(t *testing.T) {
	repository := &repositoryMock{}

	service := NewService(repository)

	err := service.Deactivate(context.Background(), "banana")
	if !errors.Is(err, ErrInvalidReaderID) {
		t.Fatalf("Deactivate() error = %v, want %v", err, ErrInvalidReaderID)
	}
}

func TestService_Deactivate_ReaderNotFound(
	t *testing.T,
) {
	const readerID = "11111111-1111-4111-8111-111111111111"

	repository := &repositoryMock{
		deactivateFn: func(
			ctx context.Context,
			id string,
		) error {
			return ErrReaderNotFound
		},
	}

	service := NewService(repository)

	err := service.Deactivate(
		context.Background(),
		readerID,
	)

	if !errors.Is(err, ErrReaderNotFound) {
		t.Fatalf(
			"Deactivate() error = %v, want %v",
			err,
			ErrReaderNotFound,
		)
	}
}
