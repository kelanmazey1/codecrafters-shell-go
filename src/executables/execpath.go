package executables

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Returns all names of files in $PATH
func Names() ([]string, error) {
	var commands []string
	dirs := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))

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
			commands = append(commands, d.Name())
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	if len(commands) == 0 {
		return nil, errors.New("no executables found")
	}

	return commands, nil
}

// Given a command name returns the executable path
func LocateExecutablePath(c string) (string, error) {
	var commandPath string
	dirs := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))

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

				info, err := os.Stat(path)

				// Checks file permissions
				if m := info.Mode(); m&0111 == 0 {
					return fs.ErrPermission
				}

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
		return "", fmt.Errorf("%s: command not found", c)
	}

	return commandPath, nil
}
