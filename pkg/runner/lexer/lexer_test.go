package lexer

import (
	"reflect"
	"testing"
)

func TestLex(t *testing.T) {
	tests := map[string]struct {
		in  string
		out []string
	}{
		"single cmd": {
			in:  "./cmd",
			out: []string{"./cmd"},
		},
		"with args": {
			in:  "./cmd -s -d",
			out: []string{"./cmd", "-s", "-d"},
		},
		"with quoted args": {
			in:  "./cmd -s -d \"test\"",
			out: []string{"./cmd", "-s", "-d", "test"},
		},
		"with quoted args that contain a space": {
			in:  "./cmd -s -d \"with a space\"",
			out: []string{"./cmd", "-s", "-d", "with a space"},
		},
		"with unbalanced quoted args": {
			in:  "./cmd -s -d \"unbalanced",
			out: []string{"./cmd", "-s", "-d", "unbalanced"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := Lex(tt.in)
			if !reflect.DeepEqual(got, tt.out) {
				t.Fatalf("want %#v got %#v", tt.out, got)
			}
		})
	}
}
