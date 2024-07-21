package down

import (
	"fmt"
	"os"
	"ryanclark532/migration-tool/internal/common"
	"slices"
	"strings"
)

func GetTableDiff(original []common.Table, post []common.Table, version int, config common.Config) error {
	processedTables := make(map[string]bool)
	processedColumns := make(map[string]bool)

	var builder strings.Builder

	for _, table := range original {
		if table.Name == "Migrations" {
			continue
		}

		postIndex := slices.IndexFunc(post, func(t common.Table) bool { return t.Name == table.Name })

		if postIndex != -1 {
			//table exists both in original and post proceed to diff table
			diffTable(&builder, table, post[postIndex], processedColumns)
			processedTables[table.Name] = true
		} else {
			//table exists in original but not in post, therefore dropped during migration. Create table
			common.CreateTable(&builder, table.Name, table.Columns)
		}
	}


	for _, table:= range post {
		if table.Name == "Migrations" {
			continue
		}

		if _, exists := processedTables[table.Name]; exists {
			//we  have already touched this table and can assume it exists in both old and new, therefore continue to columns 
			for _, column := range table.Columns {
				if _, exists := processedColumns[fmt.Sprintf("%s.%s", table.Name, column.Name)]; !exists {
					//we havent touched this column so can assume it exists in post but not pre, therefor drop column
					common.DropColumn(&builder,table.Name, column.Name)
				}
			}
		} else {
			//we havent touched this table and can assume it exists in post but not pre migration, therefore  drop table
			common.DropTable(&builder, table.Name )
		}
	}

	if builder.Len() != 0 {
		err := os.WriteFile(fmt.Sprintf("%s/down/%d-tables-down.sql",config.OutputDir, version), []byte(builder.String()), os.ModeAppend)
		return err
	}
	return nil
}

func diffTable(builder *strings.Builder, old common.Table, post common.Table, processedColumns map[string]bool) {
	for _, column := range old.Columns {
		postIndex := slices.IndexFunc(post.Columns, func(c common.Column) bool { return c.Name == column.Name })
		if postIndex != -1 {
			//column exists on table both pre and post migration, continue to diff column
			diffColumn(builder,column, post.Columns[postIndex], old.Name)
			processedColumns[fmt.Sprintf("%s.%s", old.Name, column.Name)] = true
		} else {
			//column exists on the old table but not the new. Add column
			common.AddColumn(builder, old.Name, column)
		}
	}
}

func diffColumn(builder *strings.Builder, old common.Column, post common.Column, tableName string) {
	if old.Type != post.Type {
		common.AlterColumn(builder, tableName, post)
	}
}
