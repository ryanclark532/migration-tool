package common

import (
	"fmt"
	"strings"
)

func CreateTable(builder *strings.Builder, tableName string, columns []Column){
	builder.WriteString(fmt.Sprintf(`CREATE TABLE %s (\n`, tableName))
	for _, col := range columns{
		builder.WriteString(fmt.Sprintf(`%s %s,\n`, col.Name, col.Type))
	}
	builder.WriteString(`);\n`)
}


func DropTable(builder *strings.Builder, tableName string){
	builder.WriteString(fmt.Sprintf(`DROP TABLE %s;\n`, tableName))
}

func DropColumn(builder *strings.Builder, tableName string, columnName string){
	builder.WriteString(fmt.Sprintf(`ALTER TABLE %s DROP COLUMN %s;\n`, tableName, columnName))
}

func AlterColumn(builder *strings.Builder, tableName string, column Column){
	builder.WriteString(fmt.Sprintf(`ALTER TABLE %s\n ALTER COLUMN %s %s;\n`, tableName, column.Name, column.Type))
}

func AddColumn(builder *strings.Builder, tableName string, column Column){
	builder.WriteString(fmt.Sprintf(`ALTER TABLE %s\n ADD COLUMN %s %s;\n`, tableName, column.Name, column.Type))
}
