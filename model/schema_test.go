package model

import (
	"testing"
)

func TestGenerateCreateTableQuery(t *testing.T) {
	tests := map[string]struct {
		input    []string
		expected string
	}{
		"no columns": {
			input:    []string{},
			expected: "CREATE TABLE files (name TEXT);",
		},
		"single column": {
			input:    []string{"size"},
			expected: "CREATE TABLE files (name TEXT, \"size\" TEXT);",
		},
		"multiple columns": {
			input:    []string{"size", "created_at", "modified_at"},
			expected: "CREATE TABLE files (name TEXT, \"size\" TEXT, \"created_at\" TEXT, \"modified_at\" TEXT);",
		},
		"column with spaces": {
			input:    []string{"column with spaces"},
			expected: "CREATE TABLE files (name TEXT, \"column with spaces\" TEXT);",
		},
		"column with slashes": {
			input:    []string{"column/with/slashes"},
			expected: "CREATE TABLE files (name TEXT, \"column/with/slashes\" TEXT);",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := CreateTableQuery(tc.input)
			if result != tc.expected {
				t.Errorf("expected %q but got %q", tc.expected, result)
			}
		})
	}
}
