package app

import (
	"os"
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
