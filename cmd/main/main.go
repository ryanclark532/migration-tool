package main

import (
	"fmt"
	"ryanclark532/migration-tool/internal/sqlserver"
)

func main() {
	//Get State of the Database before migration
	server := sqlserver.SqlServer{
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

	db, err := server.GetDatabaseState()
	if err != nil {
		panic(err)
	}

	fmt.Println(db)

	//Execute update scripts and increment version

	//get post migration database state

	//calculate diff between states

	//Use dif to produce a "down script for the version"

	server.Close()
}

/*
func main() {

	files := utils.CrawlDir("./testing/in")

	for _, file := range files {
		content, err := os.ReadFile(fmt.Sprintf("./testing/in/%s", file))
		if err != nil {
			panic(err)
		}

		tokenizer := lexer.NewTokenizer(string(content))
		parser := paser.CreateParser(&tokenizer)

		var queries []paser.Query

		for {
			query := parser.GetNextQuery()
			if query.Action.Type_ == lexer.Eof {
				break
			}

			if query.Action.Type_ == lexer.Illegal {
				continue
			}

			queries = append(queries, query)
		}

		down.GenerateDown(queries, file[:len(file)-4]+".down.sql")
	}
}
*/
