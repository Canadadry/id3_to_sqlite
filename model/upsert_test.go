package model

import (
	"reflect"
	"testing"
)

func TestGenerateUpsertQuery(t *testing.T) {
	tests := map[string]struct {
		input          []File
		expectedQuery  string
		expectedValues []interface{}
	}{
		"no files": {
			input:          []File{},
			expectedQuery:  "",
			expectedValues: nil,
		},
		"single file, no fields": {
			input: []File{
				{Name: "file1", Fields: map[string]string{}},
			},
			expectedQuery:  "INSERT INTO files (\"name\") VALUES (?) ON CONFLICT(name) DO UPDATE SET ;",
			expectedValues: []interface{}{"file1"},
		},
		"single file, with fields": {
			input: []File{
				{Name: "file1", Fields: map[string]string{"size": "123", "created_at": "2024-01-01"}},
			},
			expectedQuery:  "INSERT INTO files (\"name\", \"size\", \"created_at\") VALUES (?, ?, ?) ON CONFLICT(name) DO UPDATE SET \"size\" = excluded.\"size\", \"created_at\" = excluded.\"created_at\";",
			expectedValues: []interface{}{"file1", "123", "2024-01-01"},
		},
		"multiple files with different fields": {
			input: []File{
				{Name: "file1", Fields: map[string]string{"size": "123"}},
				{Name: "file2", Fields: map[string]string{"created_at": "2024-01-01", "modified_at": "2024-02-01"}},
			},
			expectedQuery:  "INSERT INTO files (\"name\", \"size\", \"created_at\", \"modified_at\") VALUES (?, ?, ?, ?), (?, ?, ?, ?) ON CONFLICT(name) DO UPDATE SET \"size\" = excluded.\"size\", \"created_at\" = excluded.\"created_at\", \"modified_at\" = excluded.\"modified_at\";",
			expectedValues: []interface{}{"file1", "123", "", "", "file2", "", "2024-01-01", "2024-02-01"},
		},
		"single file, column with spaces": {
			input: []File{
				{Name: "file1", Fields: map[string]string{"column with spaces": "value with spaces"}},
			},
			expectedQuery:  "INSERT INTO files (\"name\", \"column with spaces\") VALUES (?, ?) ON CONFLICT(name) DO UPDATE SET \"column with spaces\" = excluded.\"column with spaces\";",
			expectedValues: []interface{}{"file1", "value with spaces"},
		},
		"single file, column with slashes": {
			input: []File{
				{Name: "file1", Fields: map[string]string{"column/with/slashes": "value/with/slashes"}},
			},
			expectedQuery:  "INSERT INTO files (\"name\", \"column/with/slashes\") VALUES (?, ?) ON CONFLICT(name) DO UPDATE SET \"column/with/slashes\" = excluded.\"column/with/slashes\";",
			expectedValues: []interface{}{"file1", "value/with/slashes"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			query, values := CreateUpsertQuery(tc.input)

			if query != tc.expectedQuery {
				t.Errorf("expected query %q, got %q", tc.expectedQuery, query)
			}

			if !reflect.DeepEqual(values, tc.expectedValues) {
				t.Errorf("expected values %v, got %v", tc.expectedValues, values)
			}
		})
	}
}
