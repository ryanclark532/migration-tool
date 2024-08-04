package execute

import (
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/up"
	"strings"
)

func Up(server Server, config common.Config, dryRun bool) error {
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
		return exec(server, config, version, conn)
	} else {
		return execDry(server, config, version)
	}
}

func exec(server Server, config common.Config, version int, conn common.CommonDB) error {
	original, err := server.GetDatabaseState(config)
	if err != nil {
		return err
	}

	errs := up.DoMigration(conn, version, config)
	if len(errs) > 0 {
		return errs[0]
	}

	post, err := server.GetDatabaseState(config)
	if err != nil {
		return err
	}

	var builder strings.Builder

	down.GetTableDiff(original.Tables, post.Tables, version, &builder)
	down.GetProcDiff(original.Procs, &builder)

	if builder.Len() != 0 {
		return os.WriteFile(fmt.Sprintf("%s/down/%d.sql", config.OutputDir, version), []byte(builder.String()), os.ModeAppend)
	}
	return nil
}
func execDry(server Server, config common.Config, version int) error {
	tx, err := server.Begin()
	if err != nil {
		return err
	}

	original, err := server.GetDatabaseState(config)
	if err != nil {
		return err
	}

	errs := up.DoMigration(tx, version, config)
	if len(errs) > 0 {
		return errs[0]
	}

	post, err := server.GetDatabaseState(config)
	if err != nil {
		return err
	}

	var builder strings.Builder

	down.GetTableDiff(original.Tables, post.Tables, version, &builder)
	down.GetProcDiff(original.Procs, &builder)

	//TODO pretty print differences

	err = tx.Rollback()
	return err
}
