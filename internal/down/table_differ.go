package down

import (
	"fmt"
	"ryanclark532/migration-tool/internal/common"
	"strings"
)

func GetTableDiff(original map[string]common.Table, post map[string]common.Table, version int,  builder *strings.Builder)  {
	processedTables := make(map[string]bool)
	processedColumns := make(map[string]bool)

	for key, table := range original {
		if key == "Migrations" {
			continue
		}

		if prev, ex := post[key]; ex {
			//table exists both in original and post proceed to diff table
			diffTable(builder, table, prev, processedColumns)
			processedTables[key] = true
		} else {
			//table exists in original but not in post, therefore dropped during migration. Create table
			common.CreateTable(builder, key, table.Columns)
		}
	}

	for key, table := range post {
		if key == "Migrations" {
			continue
		}

		if _, exists := processedTables[key]; exists {
			//we  have already touched this table and can assume it exists in both old and new, therefore continue to columns
			for colKey := range table.Columns {
				if _, exists := processedColumns[fmt.Sprintf("%s.%s", key, colKey)]; !exists {
					//we havent touched this column so can assume it exists in post but not pre, therefor drop column
					common.DropColumn(builder, key, colKey)
				}
			}
		} else {
			//we havent touched this table and can assume it exists in post but not pre migration, therefore  drop table
			common.DropTable(builder, key)
		}
	}

}

func diffTable(builder *strings.Builder, old common.Table, post common.Table, processedColumns map[string]bool) {
	for key, column := range old.Columns {
		if prev, ex := post.Columns[key]; ex {
			//column exists on table both pre and post migration, continue to diff column
			diffColumn(builder, column, prev, key, key)
			processedColumns[fmt.Sprintf("%s.%s", prev, key)] = true
		} else {
			//column exists on the old table but not the new. Add column
			common.AddColumn(builder, key, column, key)
		}
	}
}

func diffColumn(builder *strings.Builder, old common.Column, post common.Column, tableName string, colName string) {
	if old.Type != post.Type {
		common.AlterColumn(builder, tableName, post, colName)
	}
}
