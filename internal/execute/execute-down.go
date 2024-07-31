package execute

import (
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"strings"
)

func ExecuteDown(server Server, config common.Config, dryRun bool) error {
	conn, err := server.Connect()
	if err != nil {
		return err
	}

	version, err := server.GetLatestVersion()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(config.OutputDir+"/down/%d.sql", version-1)
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(strings.TrimSpace(string(fileContent)))
	if err != nil {
		return err
	}

	if !dryRun {
		err = tx.Commit()
		if err != nil {
			return err
		}
		err = os.Remove(filename)
		if err != nil {
			return err
		}
		sql := fmt.Sprintf("DELETE FROM %s WHERE Version = %d", config.MigrationTableName, version-1)
		_, err := conn.Exec(sql)
		if err != nil {
			return err
		}
	} else {
		err = tx.Rollback()
		if err != nil {
			return err
		}
	}

	err = server.Close()
	if err != nil {
		return err
	}
	return nil
}
