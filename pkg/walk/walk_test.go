package walk

import (
	"os"
	"testing"
	"time"
)

type mockFileInfo struct {
	name  string
	isDir bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() os.FileMode  { return 0 }
func (m *mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

func TestListFilesFunc(t *testing.T) {
	tests := map[string]struct {
		path      string
		info      os.FileInfo
		err1      error
		err2      error
		expected  []string
		expectErr bool
	}{
		"file path, no error": {
			path:      "/path/to/file.txt",
			info:      &mockFileInfo{name: "file.txt", isDir: false},
			err1:      nil,
			err2:      nil,
			expected:  []string{"/path/to/file.txt"},
			expectErr: false,
		},
		"directory path, no error": {
			path:      "/path/to/dir",
			info:      &mockFileInfo{name: "dir", isDir: true},
			err1:      nil,
			err2:      nil,
			expected:  []string{},
			expectErr: false,
		},
		"error 1 case": {
			path:      "/path/to/file.txt",
			info:      &mockFileInfo{name: "file.txt", isDir: false},
			err1:      os.ErrPermission,
			err2:      nil,
			expected:  []string{},
			expectErr: true,
		},
		"error 2 case": {
			path:      "/path/to/file.txt",
			info:      &mockFileInfo{name: "file.txt", isDir: false},
			err2:      os.ErrPermission,
			err1:      nil,
			expected:  []string{},
			expectErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var files []string
			err := ListFilesFunc(tt.path, tt.info, tt.err1, func(path string) error {
				files = append(files, path)
				return tt.err2
			})

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			} else if err != nil {
				t.Errorf("did not expect an error but got: %v", err)
			}

			if len(files) != len(tt.expected) {
				t.Errorf("expected %v files, got %v", len(tt.expected), len(files))
			}
			for i, expectedFile := range tt.expected {
				if files[i] != expectedFile {
					t.Errorf("expected file at index %d to be %s, got %s", i, expectedFile, files[i])
				}
			}
		})
	}
}
