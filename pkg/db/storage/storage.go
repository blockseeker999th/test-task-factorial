package storage

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveCalculations(valueA int, valueB int) (int64, error) {
	const op = "storage.SaveCalculations"

	stmt, err := s.db.Prepare("INSERT INTO calculations (a, b) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	var id *int64
	err = stmt.QueryRow(valueA, valueB).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return *id, nil
}
