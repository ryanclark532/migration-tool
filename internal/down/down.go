package down

import (
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/lexer"
	"ryanclark532/migration-tool/internal/paser"
	"strings"
)

func GenerateDown(querys []paser.Query) {
	var builder strings.Builder
	for _, query := range querys {

		switch query.Action.Type_ {
		case lexer.Create:
			handleCreate(query, &builder)
		case lexer.Alter:

			fmt.Println("world")
		}

	}

	os.WriteFile("./testing/out/down.sql", []byte(builder.String()), os.ModeAppend)

}

func handleCreate(query paser.Query, builder *strings.Builder) {
	switch query.Resource.Type_ {
	case lexer.Table:
		builder.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS %s;\n", query.ResourceName.Literal))
	}
}
