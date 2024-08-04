package main

import (
	"flag"
	"fmt"
	"os"
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

	switch os.Args[1] {
	case "up":
		if err := fsUp.Parse(os.Args[2:]); err == nil {
			fmt.Println("edit", dryRun)
		}
	case "down":
		if err := fsDown.Parse(os.Args[2:]); err == nil {
			fmt.Println("edit", dryRun)
		}
	default:
		_, err := fmt.Fprintln(os.Stderr, help)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
}

/*
func main() {

	help := slices.ContainsFunc(os.Args, func(s string) bool {
		return s == "-h" || s == "-help"
	})

	_ = slices.ContainsFunc(os.Args, func(s string) bool {
		return s == "-dry-run" || s == "-dry-run=true"
	})

	if help {
		fmt.Print(helpText)
		return
	}

	operation := os.Args[1]

	if operation != "up" && operation != "down" && operation != "check" {
		panic("Invalid operation")
	}
	c := loadJson()

	if c == nil {
		panic("Unable to load config, please provide a json file, config in an .env file, or cli flags, use -h for more information")
	}

	/*
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
