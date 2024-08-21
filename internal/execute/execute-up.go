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

func Up(server Server, config common.Config, dryRun bool) error {

	version, err := server.GetLatestVersion()
	if err != nil {
		return err
	}

	conn := server.GetDB()

	if !dryRun {
		return exec(server, config, version, conn)
	} else {
		return execDry(server, config, version)
	}
}

func exec(server Server, config common.Config, version int, conn *sql.DB) error {
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

	//TODO Unhappy with the current batching of up and down.
	//refactor to generate a seperate down script for each up script.
	//throw out version system and use random generated numbers
	// e.g mgt rollback 123456 will execute down for 123456

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

	errs := up.DoDryMigration(tx, version, config)
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
