package execute

import (
	"database/sql"
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/up"
	"strings"
)

type Server interface {
	Connect() (*sql.DB, error)
	GetDatabaseState() (*common.Database, error)
	GetLatestVersion() (int, error)
	Close() error
	Setup(migrationTable string) error
}

func ExecuteUp(server Server, config common.Config, dryRun bool) error {
	conn, err := server.Connect()
	if err != nil {
		return err
	}

	err = server.Setup(config.MigrationTableName)
	if err != nil {
		return err
	}

	version, err := server.GetLatestVersion()
	if err != nil {
		return err
	}

	if !dryRun {
		original, err := server.GetDatabaseState()
		if err != nil {
			return err
		}

		errs := up.DoMigration(conn, version, config)
		if len(errs) > 0 {
			return err
		}

		post, err := server.GetDatabaseState()
		if err != nil {
			return err
		}

		var builder strings.Builder

		down.GetTableDiff(original.Tables, post.Tables, version, &builder)
		down.GetProcDiff(original.Procs, &builder)

		if builder.Len() != 0 {
			err := os.WriteFile(fmt.Sprintf("%s/down/%d.sql", config.OutputDir, version), []byte(builder.String()), os.ModeAppend)
			if err != nil {
				return err

			}
		}

	} else {
		errs := up.DoDryMigration(conn, version, config)
		if len(errs) > 0 {
			return errs[0]
		}
	}

	err = server.Close()
	if err != nil {
		return err
	}
	return nil
}
