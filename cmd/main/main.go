package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"ryanclark532/migration-tool/internal/execute"
	"ryanclark532/migration-tool/internal/sqlite"
	"ryanclark532/migration-tool/internal/sqlserver"
)

var help = "No command or an invalid command specified.\nAvailable commands are: 'up', 'down, 'check'. Run with valid command and -h to see available flags"

func main() {
	var dryRun bool

	fsUp := flag.NewFlagSet("up", flag.ExitOnError)
	fsUp.BoolVar(&dryRun, "dry-run", false, "Execute in dry run mode, changes will be validated but no actions will be commited to the database")

	var versionNo int

	fsDown := flag.NewFlagSet("down", flag.ExitOnError)
	fsDown.BoolVar(&dryRun, "dry-run", false, "Execute in dry run mode, changes will be validated but no actions will be commited to the database")
	fsDown.IntVar(&versionNo, "version", 0, "Version number to revert to. This will execute downwards migrations in order to get to the desired version")

	if len(os.Args) < 2 {
		_, err := fmt.Fprintln(os.Stderr, help)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	c := loadJson()

	if c == nil {
		panic("Unable to load config, please provide a json file, config in an .env file, or cli flags, use -h for more information")
	}

	server := getServer(c)

	switch os.Args[1] {
	case "up":
		if err := fsUp.Parse(os.Args[2:]); err != nil {
			panic(err)
		}
		err := execute.ExecuteUp(server, *c, dryRun)
		panic(err)

	case "down":
		if err := fsDown.Parse(os.Args[2:]); err != nil {
			panic(err)
		}

		err := execute.ExecuteDown(server, *c, dryRun)
		panic(err)
	case "help":
		fmt.Println("Help me")
	default:
		_, err := fmt.Fprintln(os.Stderr, help)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
}

func getServer(config *common.Config) execute.Server {
	switch config.DbType {
	case "Sqlite":
		return &sqlite.SqLiteServer{
			FilePath: config.FilePath,
		}
	case "sqlserver":
		return &sqlserver.SqlServer{
			Server:   config.Database,
			Port:     config.Port,
			User:     config.User,
			Password: config.Password,
			Database: config.Database,
		}

	default:
		panic(fmt.Sprintf("Unsupported database type: %s", config.DbType))
	}
}

func loadJson() *common.Config {
	jsonFile, err := os.ReadFile("migration-settings.json")
	if err != nil {
		panic(err)
	}
	var c common.Config
	err = json.Unmarshal(jsonFile, &c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded config from migration-settings.json")
	return &c
}
