package execpath

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// TODO: Think about whether this should be a commands.Command
func LocateExecutablePath(c string) (string, error) {
	var commandPath string
	dirs := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))

	if c == "my_exe" {
		fmt.Println(dirs)
		fmt.Println(os.Getenv("PATH"))
	}
	for _, dir := range dirs {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {

			// If directory doesn't exist then continue
			if errors.Is(err, fs.ErrNotExist) {
				return filepath.SkipDir
			}

			if err != nil {
				return err
			}
			// We only want files
			if d.IsDir() {
				return nil
			}

			if d.Name() == c {
				if c == "my_exe" {
					fmt.Println(dir)
					fmt.Println(os.Getenv("PATH"))
					fmt.Println(path)
				}
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
