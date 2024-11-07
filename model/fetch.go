package model

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type File struct {
	Name   string
	Fields map[string]string
}

func Fetch(db DBTX, ctx context.Context, limit, offset int) ([]File, error) {
	query := fmt.Sprintf("SELECT * FROM files LIMIT %d OFFSET %d", limit, offset)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var files []File

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		file := File{
			Fields: make(map[string]string),
		}
		for i, col := range columns {
			val := values[i]
			if str, ok := val.(string); ok {
				if col == "name" {
					file.Name = str
				} else {
					file.Fields[col] = str
				}
			} else if val != nil {
				file.Fields[col] = fmt.Sprintf("%v", val)
			}
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
