package down

import (
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"slices"
	"strings"
)


type processedColumn struct {
	Name      string
	TableName string
}

func GetDiff(original *common.Database, post *common.Database, version int) error {
	processedTables := make(map[string]bool)
	processedColumns := make(map[string]bool)

	var builder strings.Builder

	// Process original tables
	for _, table := range original.Tables {
		processedTables[table.Name] = true
		postTableIndex := slices.IndexFunc(post.Tables, func(t common.Table) bool { return t.Name == table.Name })
		if postTableIndex == -1 {
			common.DropTable(&builder, table.Name)
			continue
		}
		diffTable(&builder, &table, &post.Tables[postTableIndex], processedColumns)
	}

	// Process post tables
	for _, table := range post.Tables {
		if _, exists := processedTables[table.Name]; exists {
			for _, column := range table.Columns {
				key := fmt.Sprintf("%s.%s", table.Name, column.Name)
				if _, exists := processedColumns[key]; !exists {
					common.AddColumn(&builder, table.Name, column)
				}
			}
		} else {
			common.CreateTable(&builder, table.Name, table.Columns)
		}
	}

	if builder.Len() != 0 {
	err := os.WriteFile(fmt.Sprintf("./output/down/%d-down.sql", version), []byte(builder.String()), os.ModeAppend)
	return err
	}
	return nil
}

func diffTable(builder *strings.Builder, old *common.Table, post *common.Table, processedColumns map[string]bool) {
	for _, column := range old.Columns {
		key := fmt.Sprintf("%s.%s", old.Name, column.Name)
		processedColumns[key] = true
		postColumnIndex := slices.IndexFunc(post.Columns, func(c common.Column) bool { return column.Name == c.Name })
		if postColumnIndex == -1 {
			common.DropColumn(builder, old.Name, column.Name)
			continue
		}
		diffColumn(builder, &column, &post.Columns[postColumnIndex], old.Name)
	}
}

func diffColumn(builder *strings.Builder, old *common.Column, post *common.Column, tableName string) {
	if old.Type != post.Type {
		common.AlterColumn(builder, tableName, *old)
	}
}

