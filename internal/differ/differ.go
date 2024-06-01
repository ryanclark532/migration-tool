package differ

import (
	"fmt"
	"ryanclark532/migration-tool/internal/sqlserver"
	"slices"
)


type processedColumn struct {
	Name      string
	TableName string
}

func GetDiff(original *sqlserver.Database, post *sqlserver.Database) {
	processedTables := make(map[string]bool)
	processedColumns := make(map[string]bool)

	// Process original tables
	for _, table := range original.Tables {
		processedTables[table.Name] = true
		postTableIndex := slices.IndexFunc(post.Tables, func(t sqlserver.Table) bool { return t.Name == table.Name })
		if postTableIndex == -1 {
			fmt.Printf("Drop Table %s\n", table.Name)
			continue
		}
		diffTable(&table, &post.Tables[postTableIndex], processedColumns)
	}

	// Process post tables
	for _, table := range post.Tables {
		if _, exists := processedTables[table.Name]; exists {
			for _, column := range table.Columns {
				key := fmt.Sprintf("%s.%s", table.Name, column.Name)
				if _, exists := processedColumns[key]; !exists {
					fmt.Printf("Create column %s on %s\n", column.Name, table.Name)
				}
			}
		} else {
			fmt.Printf("Create Table %s\n", table.Name)
		}
	}
}

func diffTable(old *sqlserver.Table, post *sqlserver.Table, processedColumns map[string]bool) {
	for _, column := range old.Columns {
		key := fmt.Sprintf("%s.%s", old.Name, column.Name)
		processedColumns[key] = true
		postColumnIndex := slices.IndexFunc(post.Columns, func(c sqlserver.Column) bool { return column.Name == c.Name })
		if postColumnIndex == -1 {
			fmt.Printf("Drop column %s on %s\n", column.Name, old.Name)
			continue
		}
		diffColumn(&column, &post.Columns[postColumnIndex], old.Name)
	}
}

func diffColumn(old *sqlserver.Column, post *sqlserver.Column, tableName string) {
	if old.Type == post.Type {
		fmt.Printf("%s unchanged on %s\n", old.Name, tableName)
	} else {
		fmt.Printf("Alter Column %s on %s\n", old.Name, tableName)
	}
}

