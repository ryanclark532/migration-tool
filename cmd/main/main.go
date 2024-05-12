package main

import (
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/lexer"
	"ryanclark532/migration-tool/internal/paser"
)

func main() {
	content, err := os.ReadFile("./example-update.sql")
	if err != nil {
		panic(err)
	}

	tokenizer := lexer.NewTokenizer(string(content))
	parser := paser.CreateParser(&tokenizer)

	for {
		query := parser.GetNextQuery()
		if query.Action.Type_ == lexer.Eof {
			break
		}

		if query.Action.Type_ == lexer.Illegal {
			continue
		}

		fmt.Println(query)
	}

	q := parser.GetNextQuery()

	fmt.Println(q)
}
