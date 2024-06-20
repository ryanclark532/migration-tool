package common



type SchemaObject struct {
	Name       string
	ObjectName string
	ObjectType string
}

type Database struct {
	Tables []Table
	Procs  []Procedure
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
