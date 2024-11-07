package walk

import (
	"os"
	"path/filepath"
)

func ListFilesFunc(path string, info os.FileInfo, err error, action func(path string) error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if action != nil {
		return action(path)
	}
	return nil
}

func Walk(entryPoint string, action func(path string) error) error {
	return filepath.Walk(entryPoint, func(path string, info os.FileInfo, err error) error {
		return ListFilesFunc(path, info, err, action)
	})
}
