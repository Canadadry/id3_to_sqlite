package model

import (
	"fmt"
	"strings"
)

func CreateUpsertQuery(files []File) (string, []interface{}) {
	if len(files) == 0 {
		return "", nil
	}

	columns := []string{"name"}
	columnSet := make(map[string]struct{})

	for _, file := range files {
		for col := range file.Fields {
			if _, exists := columnSet[col]; !exists {
				columnSet[col] = struct{}{}
				columns = append(columns, col)
			}
		}
	}

	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	quotedColumns := make([]string, len(columns))
	for i, col := range columns {
		quotedColumns[i] = fmt.Sprintf("\"%s\"", col)
	}

	insertColumns := strings.Join(quotedColumns, ", ")
	valuePlaceholders := "(" + strings.Join(placeholders, ", ") + ")"

	updateClauses := []string{}
	for _, col := range columns[1:] {
		updateClauses = append(updateClauses, fmt.Sprintf("\"%s\" = excluded.\"%s\"", col, col))
	}
	updateClause := strings.Join(updateClauses, ", ")

	query := fmt.Sprintf(
		"INSERT INTO files (%s) VALUES %s ON CONFLICT(name) DO UPDATE SET %s;",
		insertColumns,
		strings.Repeat(valuePlaceholders+", ", len(files)-1)+valuePlaceholders,
		updateClause,
	)

	var values []interface{}
	for _, file := range files {
		values = append(values, file.Name)
		for _, col := range columns[1:] {
			values = append(values, file.Fields[col])
		}
	}

	return query, values
}
