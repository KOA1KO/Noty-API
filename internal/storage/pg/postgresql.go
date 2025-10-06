package pg

import (
	"Noty/internal/storage/migrate"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New(storageDSN string) (*Storage, error) {
	const op = "internal.storage.pg.New"

	db, err := sql.Open("pgx", storageDSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// hardcode: delete and run with CI/CD
	if err := migrate.Run(db); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
