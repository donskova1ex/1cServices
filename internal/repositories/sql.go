package repositories

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewSQLDB(dbDSN string) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", dbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to sql: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sql: %w", err)
	}
	return db, nil
}
