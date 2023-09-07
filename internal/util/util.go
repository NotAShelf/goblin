package util

import (
	"os"
)

// CreateDirectoryIfNotExists creates the specified directory if it doesn't exist.
func CreateDirectoryIfNotExists(directoryPath string) error {
	// Check if the directory already exists
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		// Directory doesn't exist, create it
		if err := os.MkdirAll(directoryPath, 0755); err != nil {
			return err
		}
	}
	return nil
}
