package main

import (
	"ryanclark532/migration-tool/internal/differ"
	"ryanclark532/migration-tool/internal/migrate"
	"ryanclark532/migration-tool/internal/sqlserver"
)

func main() {
	//Get State of the Database before migration
	server := &sqlserver.SqlServer{
		Server:   "localhost",
		Port:     1433,
		User:     "sa",
		Password: "yourStrong(!)Password",
		Database: "master",
	}

	err := server.Connect()
	if err != nil {
		panic(err)
	}

	version,err := server.GetLatestVersion()
	if err != nil{
		panic(err)
	}

	original, err := server.GetDatabaseState()
	if err != nil {
		panic(err)
	}

	migrate.DoMigration(server.Conn)

	post, err:= server.GetDatabaseState()
	if err != nil {
		panic(err)
	}


	differ.GetDiff(original, post, version)

	server.Close()
}
