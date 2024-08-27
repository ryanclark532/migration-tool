package common

import (
	"database/sql"
)

type Server interface {
	Connect() (*sql.DB, error)
	GetDatabaseState(config Config) (*Database, error)
	GetLatestVersion() (int, error)
	Setup(migrationTable string) error
	Begin() (*sql.Tx, error)
	Close() error
	GetDB() *sql.DB
}
