package common

import "database/sql"

type TxStarter interface {
	Begin() (*sql.Tx, error)
}

type CommonDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
