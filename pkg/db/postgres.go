package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/blockseeker999th/test-task-factorial/pkg/config"

	_ "github.com/lib/pq"
)

type PostgreSQLStorage struct {
	db *sql.DB
}

func ConnectDB(config *config.Config) *PostgreSQLStorage {
	const op = "db.storage.ConnectDB"

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("error path: %s, error: %v", op, err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("error pinging PostgreSQL: %v", err)
	}

	return &PostgreSQLStorage{db: db}
}

func (s *PostgreSQLStorage) InitNewPostgreSQLStorage() (*sql.DB, error) {
	if err := s.createCalculationsTable(); err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *PostgreSQLStorage) createCalculationsTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS calculations (
	id SERIAL PRIMARY KEY,
	a INT NOT NULL CHECK (a >= 0),
	b INT NOT NULL CHECK (a >= 0),
	createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);
	`)

	if err != nil {
		return err
	}

	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_calc ON calculations(a,b)`)

	if err != nil {
		return err
	}

	return nil
}
