package common

import (
	"fmt"
	"strings"
)

func CreateTable(builder *strings.Builder, tableName string, columns map[string]Column) {
	builder.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))
	for key, col := range columns {
		builder.WriteString(fmt.Sprintf("%s %s,\n", key, col.Type))
	}
	builder.WriteString(");\n")
}

func DropTable(builder *strings.Builder, tableName string) {
	builder.WriteString(fmt.Sprintf("DROP TABLE %s;\n", tableName))
}

func DropColumn(builder *strings.Builder, tableName string, columnName string) {
	builder.WriteString(fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;\n", tableName, columnName))
}

func AlterColumn(builder *strings.Builder, tableName string, column Column, colName string) {
	builder.WriteString(fmt.Sprintf("ALTER TABLE %s ALTER %s %s;\n", tableName, colName, column.Type))
}

func AddColumn(builder *strings.Builder, tableName string, column Column, colName string) {
	builder.WriteString(fmt.Sprintf("ALTER TABLE %s ADD %s %s;\n", tableName, colName, column.Type))
}

func DropProc(builder *strings.Builder, procName string) {
	//thing
}
