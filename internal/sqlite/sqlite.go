package sqlite

import (
	"database/sql"
	"errors"
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

func (s *SqLiteServer) Setup(migrationTable string) error {
	var exists string
	sqlBatch := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	err := s.Conn.QueryRow(sqlBatch, migrationTable).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		sqlBatch := fmt.Sprintf(`CREATE TABLE %s(
			EnterDateTime DATETIME2,
			Version INTEGER, 
			FileName VARCHAR(256)
		);`, migrationTable)
		_, err = s.Conn.Exec(sqlBatch)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (s *SqLiteServer) getTables() (map[string]common.Table, error) {
	sqlBatch := `SELECT name FROM sqlite_master WHERE type='table'`
	rows, err := s.Conn.Query(sqlBatch)
	if err != nil {
		return nil, err
	}

	tables := make(map[string]common.Table)
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		var table common.Table
		columns, err := s.GetTableColumns(tableName)
		if err != nil {
			return nil, err
		}
		table.Columns = columns

		tables[tableName] = table

	}
	return tables, nil
}

func (s *SqLiteServer) GetTableColumns(tableName string) (map[string]common.Column, error) {
	sqlBatch := fmt.Sprintf(`SELECT name, type FROM pragma_table_info('%s');`, tableName)
	rows, err := s.Conn.Query(sqlBatch)
	if err != nil {
		return nil, err
	}

	columns := make(map[string]common.Column)

	for rows.Next() {
		var colName string
		var col common.Column
		err = rows.Scan(&colName, &col.Type)
		if err != nil {
			return nil, err
		}
		columns[colName] = col

	}
	return columns, nil
}

func (s *SqLiteServer) GetDatabaseState() (*common.Database, error) {
	tables, err := s.getTables()
	procs:= make(map[string]common.Procedure)
	if err != nil {
		return nil, err
	}
	return &common.Database{
		Tables: tables,
		Procs: procs,
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

func (s *SqLiteServer) GetTableIndexes(tableName string) (map[string]common.Index, error) {
	sql := fmt.Sprintf(`
SELECT 
    idx.name AS index_name,
    'index' AS index_type,  
    info.name AS column_name,
    info.seqno + 1 AS column_position  
FROM 
    sqlite_master AS idx
JOIN 
    pragma_index_info(idx.name) AS info
WHERE 
    idx.type = 'index'
    AND idx.tbl_name = '%s'
ORDER BY 
    idx.name, info.seqno;
	`, tableName)

	rows, err := s.Conn.Query(sql)
	if err != nil {
		return nil, err
	}
	indexes := make(map[string]common.Index)
	for rows.Next() {
		var index common.Index
		var indexName string
		_ = rows.Scan(&indexName, &index.Type, &index.ColumnName, &index.Position)
		indexes[indexName] = index
	}
	return indexes, err
}
