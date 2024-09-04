package common

type SchemaObject struct {
	ObjectName string
	ObjectType string
}

type Database struct {
	Tables map[string]Table
	Procs  map[string]Procedure
}

type Procedure struct {
	Definition string
}

type Table struct {
	Columns    map[string]Column
	Contraints map[string]Constraint
	Indexes    map[string]Index
}

type Constraint struct {
	ColumnName string
	Type       string
}

type Index struct {
	Type       string
	ColumnName string
	Position   string
}

type Column struct {
	Type string
}
