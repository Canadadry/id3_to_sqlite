package model

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFetchAllFiles(t *testing.T) {
	tests := map[string]struct {
		limit    int
		offset   int
		mockRows []string
		mockCols []string
		expected []File
	}{
		"no rows": {
			limit:    10,
			offset:   0,
			mockRows: []string{},
			mockCols: []string{"name"},
			expected: []File{},
		},
		"single row": {
			limit:    10,
			offset:   0,
			mockRows: []string{"file1"},
			mockCols: []string{"name"},
			expected: []File{
				{
					Name:   "file1",
					Fields: map[string]string{},
				},
			},
		},
		"multiple rows with columns": {
			limit:  10,
			offset: 0,
			mockRows: []string{
				"file1", "value1",
				"file2", "value2",
			},
			mockCols: []string{"name", "extra"},
			expected: []File{
				{
					Name:   "file1",
					Fields: map[string]string{"extra": "value1"},
				},
				{
					Name:   "file2",
					Fields: map[string]string{"extra": "value2"},
				},
			},
		},
		"column with spaces": {
			limit:  10,
			offset: 0,
			mockRows: []string{
				"file1", "value with spaces",
			},
			mockCols: []string{"name", "column with spaces"},
			expected: []File{
				{
					Name:   "file1",
					Fields: map[string]string{"column with spaces": "value with spaces"},
				},
			},
		},
		"column with slashes": {
			limit:  10,
			offset: 0,
			mockRows: []string{
				"file1", "value/with/slashes",
			},
			mockCols: []string{"name", "column/with/slashes"},
			expected: []File{
				{
					Name:   "file1",
					Fields: map[string]string{"column/with/slashes": "value/with/slashes"},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows(tc.mockCols)
			for i := 0; i < len(tc.mockRows); i += len(tc.mockCols) {
				values := make([]driver.Value, len(tc.mockCols))
				for j := range tc.mockCols {
					values[j] = tc.mockRows[i+j]
				}
				rows.AddRow(values...)
			}

			mock.ExpectQuery("SELECT \\* FROM files LIMIT \\d+ OFFSET \\d+").
				WillReturnRows(rows)

			ctx := context.Background()
			files, err := Fetch(db, ctx, tc.limit, tc.offset)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(files) != len(tc.expected) {
				t.Fatalf("expected %d files, got %d", len(tc.expected), len(files))
			}
			for i, file := range files {
				if file.Name != tc.expected[i].Name {
					t.Errorf("expected name %q, got %q", tc.expected[i].Name, file.Name)
				}
				for k, v := range tc.expected[i].Fields {
					if file.Fields[k] != v {
						t.Errorf("expected field %q to have value %q, got %q", k, v, file.Fields[k])
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
