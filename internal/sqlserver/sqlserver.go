package sqlserver

import (
	"database/sql"
	"fmt"
	"ryanclark532/migration-tool/internal/common"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

type SqlServer struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
	Conn     *sql.DB
}

func (s *SqlServer) Connect() (*sql.DB, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		s.Server, s.User, s.Password, s.Port, s.Database)

	db, err := sql.Open("sqlserver", connString)
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

func (s *SqlServer) Close() error {
	return s.Conn.Close()
}

func (s *SqlServer) getServerObjects() ([]common.SchemaObject, error) {
	sqlContent := `
	SELECT 
    schema_name(schema_id) AS schema_name,
    name AS object_name,
    type_desc AS object_type
	FROM 
 	   sys.objects
	WHERE 
 	   schema_id = SCHEMA_ID('dbo');
	`
	rows, err := s.Conn.Query(sqlContent)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemaObjects []common.SchemaObject

	for rows.Next() {
		var t common.SchemaObject
		_ = rows.Scan(&t.Name, &t.ObjectName, &t.ObjectType)

		if strings.HasPrefix(t.ObjectName, "spt_") || t.ObjectName == "MSreplication_options" || t.ObjectName == "Migrations" {
			continue
		}
		schemaObjects = append(schemaObjects, t)
	}
	return schemaObjects, nil
}

func (s *SqlServer) getTableColumns(tableName string) ([]common.Column, error) {
	sqlContent := fmt.Sprintf(`SELECT COLUMN_NAME, DATA_TYPE
	FROM INFORMATION_SCHEMA.COLUMNS
	WHERE TABLE_SCHEMA = 'dbo' AND TABLE_NAME = '%s';
	`, tableName)

	rows, err := s.Conn.Query(sqlContent)
	if err != nil {
		return nil, err
	}

	var columns []common.Column
	for rows.Next() {
		var t common.Column
		_ = rows.Scan(&t.Name, &t.Type)
		columns = append(columns, t)
	}
	return columns, err
}

func (s *SqlServer) getTableContrains(tablename string) ([]common.Constraint, error) {
	sql := fmt.Sprintf(`
	SELECT 
    tc.constraint_name AS constraint_name,
    kc.column_name AS column_name,
    tc.constraint_type AS constraint_type
	FROM 
 	   INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
	JOIN 
 	   INFORMATION_SCHEMA.KEY_COLUMN_USAGE kc ON tc.constraint_name = kc.constraint_name
	WHERE 
  	  kc.table_name = '%s'
	ORDER BY 
 	   tc.constraint_name, kc.ordinal_position;
	`, tablename)

	rows, err := s.Conn.Query(sql)
	var constraints []common.Constraint
	for rows.Next() {
		var t common.Constraint
		_ = rows.Scan(&t.Name, &t.ColumnName, &t.Type)
		constraints = append(constraints, t)
	}
	return constraints, err
}

func (s *SqlServer) getTableIndexes(tableName string) ([]common.Index, error) {
	sql := fmt.Sprintf(`
	SELECT 
    idx.name AS index_name,
    idx.type_desc AS index_type,
    col.name AS column_name,
    ic.key_ordinal AS column_position
	FROM 
	    sys.indexes idx
	JOIN 
 	   sys.index_columns ic ON idx.object_id = ic.object_id AND idx.index_id = ic.index_id
	JOIN 
 	   sys.columns col ON ic.object_id = col.object_id AND ic.column_id = col.column_id
	WHERE 
  	  idx.object_id = OBJECT_ID('%s')
	ORDER BY 
 	   idx.name, ic.key_ordinal;
	`, tableName)

	rows, err := s.Conn.Query(sql)
	if err != nil {
		return nil, err
	}
	var indexes []common.Index
	for rows.Next() {
		var t common.Index
		_ = rows.Scan(&t.Name, &t.Type, &t.ColumnName, &t.Position)
		indexes = append(indexes, t)
	}
	return indexes, err
}

func (s *SqlServer) getProcedureDetails(procName string) (common.Procedure, error) {
	sql := fmt.Sprintf(`
	SELECT 
    definition
	FROM 
	    sys.sql_modules
	WHERE 
	    object_id = OBJECT_ID('%s');
	`, procName)

	rows, err := s.Conn.Query(sql)
	if err != nil {
		return common.Procedure{}, err
	}

	var description string
	for rows.Next() {
		_ = rows.Scan(&description)
	}
	return common.Procedure{Name: procName, Definition: strings.TrimSpace(description)}, nil
}

func (s *SqlServer) GetDatabaseState() (*common.Database, error) {
	objects, err := s.getServerObjects()
	if err != nil {
		panic(err)
	}

	var tables []common.Table
	var procedures []common.Procedure

	for _, object := range objects {
		switch object.ObjectType {
		case "USER_TABLE":
			columns, err := s.getTableColumns(object.ObjectName)
			if err != nil {
				return nil, err
			}

			constraints, err := s.getTableContrains(object.ObjectName)
			if err != nil {
				return nil, err
			}

			indexes, err := s.getTableIndexes(object.ObjectName)
			if err != nil {
				return nil, err
			}
			t := common.Table{
				Columns:    columns,
				Name:       object.ObjectName,
				Contraints: constraints,
				Indexes:    indexes,
			}
			tables = append(tables, t)

		case "SQL_STORED_PROCEDURE":
			proc, err := s.getProcedureDetails(object.ObjectName)
			if err != nil {
				return nil, err
			}
			procedures = append(procedures, proc)
		}

	}
	return &common.Database{
		Tables: tables,
		Procs:  procedures,
	}, nil
}
func (s *SqlServer) GetLatestVersion() (int, error) {
	sql := `SELECT MAX(Version) FROM Migrations`

	rows := s.Conn.QueryRow(sql)

	version := 0
	err := rows.Scan(&version)
	if err != nil && strings.Contains(err.Error(), "converting NULL to int is unsupported") {
		err = nil
	}
	return version + 1, err
}
