package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/6ivkin/test.git/internal/reader"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	postgresUniqueViolation = "23505"
	readersEmailConstraint  = "uq_readers_email"
)

type ReaderRepository struct {
	pool *pgxpool.Pool
}

func NewReaderRepository(pool *pgxpool.Pool) *ReaderRepository {
	return &ReaderRepository{
		pool: pool,
	}
}

func (r *ReaderRepository) Create(ctx context.Context, entity reader.Reader) (reader.Reader, error) {
	query, args, err := goqu.
		Dialect("postgres").
		Insert("readers").
		Rows(
			goqu.Record{
				"full_name": entity.FullName,
				"email":     entity.Email,
			},
		).
		Returning(
			"id",
			"full_name",
			"email",
			"is_active",
			"created_at",
			"updated_at",
		).
		Prepared(true).
		ToSQL()

	if err != nil {
		return reader.Reader{}, fmt.Errorf("build create reader query: %w", err)
	}

	var result reader.Reader

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&result.ID,
		&result.FullName,
		&result.Email,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == postgresUniqueViolation &&
			pgErr.ConstraintName == readersEmailConstraint {
			return reader.Reader{},
				reader.ErrEmailAlreadyExsists
		}

		return reader.Reader{}, fmt.Errorf("create reader: %w", err)
	}

	return result, nil
}

func (r *ReaderRepository) GetByID(ctx context.Context, id string) (reader.Reader, error) {
	query, args, err := goqu.Dialect("postgres").From("readers").
		Select(
			"id",
			"full_name",
			"email",
			"is_active",
			"created_at",
			"updated_at",
		).Where(
		goqu.Ex{
			"id": id,
		},
	).Prepared(true).
		ToSQL()

	if err != nil {
		return reader.Reader{},
			fmt.Errorf("build get reader query: %w", err)
	}

	var result reader.Reader

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&result.ID,
		&result.FullName,
		&result.Email,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return reader.Reader{}, reader.ErrReaderNotFound
		}

		return reader.Reader{}, fmt.Errorf("get ready by id: %w", err)
	}

	return result, nil
}

func (r *ReaderRepository) Deactivate(ctx context.Context, id string) error {
	query, args, err := goqu.
		Dialect("postgres").
		Update("readers").
		Set(
			goqu.Record{
				"is_active":  false,
				"updated_at": goqu.L("NOW()"),
			},
		).
		Where(
			goqu.Ex{
				"id": id,
			},
		).
		Prepared(true).
		ToSQL()

	if err != nil {
		return fmt.Errorf("build deactivate reader query: %w", err)
	}

	commandTag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("deactiveted reader: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return reader.ErrReaderNotFound
	}

	return nil
}
