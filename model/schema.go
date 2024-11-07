package model

import (
	"fmt"
)

func CreateTableQuery(columns []string) string {
	query := "CREATE TABLE files (name TEXT"
	for _, column := range columns {
		query += fmt.Sprintf(", \"%s\" TEXT", column)
	}
	query += ");"

	return query
}
