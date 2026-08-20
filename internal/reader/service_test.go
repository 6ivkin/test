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

func TestService_GetByID(t *testing.T) {
	const readerID = "d819891f-8d3f-4c2a-9d93-c71a6441aebe"

	expectedReader := Reader{
		ID:       readerID,
		FullName: "Egor Ivanov",
		Email:    "egor@example.com",
		IsActive: true,
	}

	tests := []struct {
		name       string
		id         string
		repository Repository
		want       Reader
		wantErr    error
	}{
		{
			name: "success",
			id:   readerID,
			repository: &repositoryMock{
				getByIDFn: func(
					ctx context.Context,
					id string,
				) (Reader, error) {
					return expectedReader, nil
				},
			},
			want: expectedReader,
		},
		{
			name: "reader not found",
			id:   readerID,
			repository: &repositoryMock{
				getByIDFn: func(
					ctx context.Context,
					id string,
				) (Reader, error) {
					return Reader{}, ErrReaderNotFound
				},
			},
			wantErr: ErrReaderNotFound,
		},
		{
			name:       "invalid reader id",
			id:         "banana",
			repository: &repositoryMock{},
			wantErr:    ErrInvalidReaderID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.repository)

			got, err := service.GetByID(
				context.Background(),
				tt.id,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"GetByID() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if got != tt.want {
				t.Fatalf(
					"GetByID() = %+v, want %+v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	createdReader := Reader{
		ID:       "d819891f-8d3f-4c2a-9d93-c71a6441aebe",
		FullName: "Egor Ivanov",
		Email:    "egor@example.com",
		IsActive: true,
	}

	tests := []struct {
		name       string
		input      CreateInput
		repository Repository
		want       Reader
		wantErr    error
	}{
		{
			name: "success",
			input: CreateInput{
				FullName: "  Egor Ivanov  ",
				Email:    "  egor@example.com  ",
			},
			repository: &repositoryMock{
				createFn: func(
					ctx context.Context,
					reader Reader,
				) (Reader, error) {
					expectedReader := Reader{
						FullName: "Egor Ivanov",
						Email:    "egor@example.com",
					}

					if reader != expectedReader {
						t.Fatalf(
							"Create() repository reader = %+v, want %+v",
							reader,
							expectedReader,
						)
					}

					return createdReader, nil
				},
			},
			want: createdReader,
		},
		{
			name: "empty full name",
			input: CreateInput{
				FullName: "   ",
				Email:    "egor@example.com",
			},
			repository: &repositoryMock{},
			wantErr:    ErrInvalidFullName,
		},
		{
			name: "empty email",
			input: CreateInput{
				FullName: "Egor Ivanov",
				Email:    "   ",
			},
			repository: &repositoryMock{},
			wantErr:    ErrInvalidEmail,
		},
		{
			name: "invalid email",
			input: CreateInput{
				FullName: "Egor Ivanov",
				Email:    "not-an-email",
			},
			repository: &repositoryMock{},
			wantErr:    ErrInvalidEmail,
		},
		{
			name: "email already exists",
			input: CreateInput{
				FullName: "Egor Ivanov",
				Email:    "egor@example.com",
			},
			repository: &repositoryMock{
				createFn: func(
					ctx context.Context,
					reader Reader,
				) (Reader, error) {
					return Reader{}, ErrEmailAlreadyExists
				},
			},
			wantErr: ErrEmailAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.repository)

			got, err := service.Create(
				context.Background(),
				tt.input,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"Create() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if got != tt.want {
				t.Fatalf(
					"Create() = %+v, want %+v",
					got,
					tt.want,
				)
			}
		})
	}
}
