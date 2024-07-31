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

func main() {

	c := loadJson()

	operation := flag.String("operation", "up", "Migration Operation, available operations: up, down")

	dryRun := flag.Bool("dry-run", false, "Run in dry-run mode")

	flag.Parse()

	if c == nil {
		panic("Unable to load config, please provide a json file, config in an .env file, or cli flags, use -h for more information")
	}

	server := getServer(c)

	switch *operation {
	case "up":
		err := execute.ExecuteUp(server, *c, *dryRun)
		panic(err)
	case "down":
		err := execute.ExecuteDown(server, *c, *dryRun)
		panic(err)
	default:
		panic(fmt.Sprintf("Unsupported operation: %s", *operation))
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

/*
func loanEnv() *config {
	p, err := strconv.ParseInt(os.Getenv("port"), 0, 8)
	if err != nil {
		panic(err)
	}
	return &config{
		DbType:   os.Getenv("dbType"),
		FilePath: os.Getenv("filePath"),
		Name:     os.Getenv("database"),
		Port:     int(p),
		user:     os.Getenv("user"),
		password: os.Getenv("password"),
		database: os.Getenv("database"),
	}

}

func loadFlags() *config {
	c := config{}

	flag.StringVar(&c.DbType, "type", "Sqlite", "The type of the Database, e.g Sqlite, SqlServer")

	flag.StringVar(&c.FilePath, "path", "./database.db", "If type if Sqlite, path to database file")

	flag.StringVar(&c.Name, "server", "database", "The FQDN of the Server")

	flag.IntVar(&c.Port, "port", 0, "The port the database is listening on")

	flag.StringVar(&c.user, "user", "user", "The username to authenticate against the database")

	flag.StringVar(&c.password, "password", "password", "The password to authenticate against the database")

	flag.StringVar(&c.database, "database", "database", "The name of the database of the server")

	return &c
}
*/
