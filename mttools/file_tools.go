package mttools

import (
	"errors"
	"os"
	"path/filepath"
)

// Returns true if file exists and access is not denied.
func IsFileExists(path string) bool {
	file_info, err := os.Stat(path)

	if os.IsNotExist(err) || os.IsPermission(err) {
		return false
	}

	if file_info.IsDir() {
		return false
	}

	return true
}

// Returns true if directory exists and access is not denied.
func IsDirExists(path string) bool {
	file_info, err := os.Stat(path)

	if os.IsNotExist(err) || os.IsPermission(err) {
		return false
	}

	if !file_info.IsDir() {
		return false
	}

	return true
}

// Returns absolute path for directory. If `path` is "" current working directory is used.
func GetDirAbsolutePath(path string) (abs_path string, err error) {
	abs_path = path

	if abs_path == "" {
		abs_path = "." //current directory
	}

	if !filepath.IsAbs(abs_path) {
		abs_path, err = filepath.Abs(path)
		if err != nil {
			return
		}
	}

	if !IsDirExists(abs_path) {
		return abs_path, errors.New("\"" + path + "\" directory does not exists")
	}

	return
}
