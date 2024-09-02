package down

import (
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"strings"
)

func Down(server common.Server, config common.Config, dryRun bool, fileName string) error {
	//TODO Rethink "down" execution.
	/*
		Mode One
		1. Take in filename
		2. Has up script been run? Has Down script been generated
		3. Run down script and delete it, delete record from DB

		Mode two
		1. somehow store and calculate the last time a down migration has been run
		2. get all up scripts that have been run since then
		3. for each run the down script and delete record from DB
	*/

	files, err := common.CompletedFiles(server.GetDB())
	if err != nil {
		return err
	}

	ex := files[strings.Split(fileName, ".down")[0]]
	if !ex {
		return fmt.Errorf("Down File doesnt exist")
	}
	contents, err := os.ReadFile(fmt.Sprintf("%s/%s", config.OutputDir, fileName))
	if err != nil {
		return err
	}

	tx, err := server.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(contents))
	if err != nil {
		return err
	}

	if dryRun {
		return tx.Rollback()
	} else {
		return tx.Commit()
	}
}
