package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func DeleteAllData(path *string) error {
	dirEntries, err := os.ReadDir(*path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range dirEntries {
		entryPath := filepath.Join(*path, entry.Name())

		if err := os.RemoveAll(entryPath); err != nil {
			return fmt.Errorf("failed to remove %s: %w", entryPath, err)
		}
	}

	return nil
}
