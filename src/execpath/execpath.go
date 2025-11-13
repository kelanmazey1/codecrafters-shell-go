package execpath

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Executable struct {
	Path string
	Info fs.FileInfo
}

// Find executable for c
// TODO: Think about whether this should be a commands.Command
func LocateExecutablePath(c string) (string, error) {
	dirs := GetPathDirs()
	var commandPath string

	for _, dir := range dirs {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// We only want files
			if d.IsDir() {
				return nil
			}

			if d.Name() == c {
				_, err := os.Stat(path) // Checks for permissions

				if err != nil {
					return err
				}

				commandPath = path

				return filepath.SkipAll
			}
			return nil

		})

		// Continue if error is permissions
		if errors.Is(err, fs.ErrPermission) {
			continue
		}
		// Return any error when walking dirs
		if err != nil && err != filepath.SkipAll {
			return "", err
		}
	}
	if commandPath == "" {
		// commandPath is never set ie. c executable doesn't exist
		return "", fmt.Errorf("executable for '%s' not found in PATH", c)
	}

	return commandPath, nil
}

func GetPathDirs() []string {
	path := os.Getenv("PATH")
	return strings.Split(path, string(os.PathListSeparator))
}
