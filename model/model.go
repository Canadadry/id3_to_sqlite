package model

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Close() error
}

type File struct {
	Name   string
	Fields map[string]string
}

func Open(ctx context.Context, path string) (DBTX, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("cannot open database : %w", err)
	}

	return db, nil
}

func Migrate(db DBTX, ctx context.Context, columns []string) error {
	_, err := db.ExecContext(ctx, CreateTableQuery(columns))
	if err != nil {
		return fmt.Errorf("cannot create database structure : %w", err)
	}
	return nil
}

func Upsert(db DBTX, ctx context.Context, files []File) error {
	query, values := CreateUpsertQuery(files)

	_, err := db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("cannot upsert : %w", err)
	}
	return nil
}
