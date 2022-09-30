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

// Returns number of `value` values found in `values_list`
func CountValues(value interface{}, values_list ...interface{}) (count int) {
	count = 0

	for _, element := range values_list {
		if element == value {
			count++
		}
	}
	return count
}
