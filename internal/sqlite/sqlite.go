package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"ryanclark532/migration-tool/internal/common"
	"strings"
)

type SqLiteServer struct {
	FilePath string
	Conn     *sql.DB
}

func (s *SqLiteServer) Connect() (*sql.DB, error) {
	//get serve file from options
	db, err := sql.Open("sqlite3", s.FilePath)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	s.Conn = db
	return db, nil
}

func (s *SqLiteServer) Close() error {
	return s.Conn.Close()
}

func (s *SqLiteServer) getTables() ([]common.Table, error) {
	sqlBatch := `SELECT name FROM sqlite_master WHERE type='table'`
	rows, err := s.Conn.Query(sqlBatch)
	if err != nil {
		return nil, err
	}

	var tables []common.Table
	for rows.Next() {
		var table common.Table
		err = rows.Scan(&table.Name)
		if err != nil {
			return nil, err
		}
		columns, err := s.GetTableColumns(table.Name)
		if err != nil {
			return nil, err
		}
		table.Columns = columns

		tables = append(tables, table)

	}
	return tables, nil
}

func (s *SqLiteServer) GetTableColumns(tableName string) ([]common.Column, error) {
	sqlBatch := fmt.Sprintf(`SELECT name, type FROM pragma_table_info('%s');`, tableName)
	rows, err := s.Conn.Query(sqlBatch)
	if err != nil {
		return nil, err
	}

	var columns []common.Column

	for rows.Next() {
		var col common.Column
		err = rows.Scan(&col.Name, &col.Type)
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)

	}
	return columns, nil
}

func (s *SqLiteServer) GetDatabaseState() (*common.Database, error) {
	tables, err := s.getTables()
	if err != nil {
		return nil, err
	}
	return &common.Database{
		Tables: tables,
	}, nil
}

func (s *SqLiteServer) GetLatestVersion() (int, error) {
	sqlBatch := `SELECT MAX(Version) FROM Migrations`

	rows := s.Conn.QueryRow(sqlBatch)

	version := 0
	err := rows.Scan(&version)
	if err != nil && strings.Contains(err.Error(), "converting NULL to int is unsupported") {
		err = nil
	}
	return version + 1, err
}
