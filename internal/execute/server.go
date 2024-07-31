package execute

import (
	"database/sql"
	"ryanclark532/migration-tool/internal/common"
)

type Server interface {
	Connect() (*sql.DB, error)
	GetDatabaseState(config common.Config) (*common.Database, error)
	GetLatestVersion() (int, error)
	Close() error
	Setup(migrationTable string) error
}
