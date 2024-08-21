package execute

import (
	"database/sql"
	"ryanclark532/migration-tool/internal/common"
)

type Server interface {
	Connect() (*sql.DB, error)
	GetDatabaseState(config common.Config) (*common.Database, error)
	GetLatestVersion() (int, error)
	Setup(migrationTable string) error
	Begin() (*sql.Tx, error)
	Close() error
	GetDB() *sql.DB
}
