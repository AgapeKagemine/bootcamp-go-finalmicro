package database

import (
	"database/sql"
	"fmt"
	"orchestrator/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB() (db *sql.DB, err error) {
	config := config.Database{
		Driver:   "pgx",
		Username: "training",
		Password: 1234,
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "gomicro-orchestrator",
	}

	connString := fmt.Sprintf("postgres://%s:%d@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database)

	db, err = sql.Open(config.Driver, connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
