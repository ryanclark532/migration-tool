package common

import "database/sql"

type commonDb interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func CompletedFiles(conn commonDb) (map[string]bool, error) {
	sqlBatch := ("SELECT FileName From Migrations")
	fileNames := make(map[string]bool)
	rows, err := conn.Query(sqlBatch)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var fileName string
		err = rows.Scan(&fileName)
		if err != nil {
			return fileNames, err
		}
		fileNames[fileName] = true
	}

	return fileNames, nil
}
