package sqlserver

import "database/sql"

type SchemaObject struct {
	Name       string
	ObjectName string
	ObjectType string
}

type SqlServer struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
	conn     *sql.DB
}

type Database struct {
	tables []Table
	procs  []Procedure
}

type Procedure struct {
	Name       string
	Definition string
}

type Table struct {
	Name       string
	Columns    []Column
	Contraints []Constraint
	Indexes    []Index
}

type Constraint struct {
	Name       string
	ColumnName string
	Type       string
}

type Index struct {
	Name       string
	Type       string
	ColumnName string
	Position   string
}

type Column struct {
	Name string
	Type string
}
