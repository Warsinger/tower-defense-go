//go:build !js

package components

import (
	"os"
	"path/filepath"
)

func readStats() ([]byte, error) { return os.ReadFile(statsFile) }

func writeStats(data []byte) error {
	dir, _ := filepath.Split(statsFile)
	if err := ensureDir(dir); err != nil {
		return err
	}
	return os.WriteFile(statsFile, data, 0644)
}
