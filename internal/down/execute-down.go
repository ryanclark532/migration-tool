package down

import (
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/up"
)

func Down(server common.Server, config common.Config, dryRun bool, fileName string) []error {
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
	completedFiles, err := up.CompletedFiles(server.GetDB())
	if err != nil {
		return []error{err}
	}

}
