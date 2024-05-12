package main

import (
	"os"
	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/lexer"
	"ryanclark532/migration-tool/internal/paser"
)

func main() {
	content, err := os.ReadFile("./testing/in/example.sql")
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

	down.GenerateDown(queries)

}
