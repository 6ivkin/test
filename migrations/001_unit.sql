CREATE TABLE
	readers (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
		full_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
	);

CREATE TABLE
	books (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		isbn VARCHAR(32),
		inventory_number VARCHAR(100) NOT NULL UNIQUE,
		status VARCHAR(30) NOT NULL DEFAULT 'available',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		CONSTRAINT books_status_check CHECK (status IN ('available', 'borrowed', 'removed'))
	);

CREATE TABLE
	loans (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
		reader_id UUID NOT NULL REFERENCES readers (id),
		book_id UUID NOT NULL REFERENCES books (id),
		issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
		due_at TIMESTAMPTZ NOT NULL,
		returned_at TIMESTAMPTZ,
		renew_count INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
	);

CREATE INDEX idx_loans_reader_id ON loans (reader_id);

CREATE INDEX idx_loans_book_id ON loans (book_id);

CREATE UNIQUE INDEX idx_unique_active_book_loan ON loans (book_id)
WHERE
	returned_at IS NULL;