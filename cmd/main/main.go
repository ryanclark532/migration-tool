package main

import (
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/down"
	"ryanclark532/migration-tool/internal/lexer"
	"ryanclark532/migration-tool/internal/paser"
	"ryanclark532/migration-tool/internal/utils"
)

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
