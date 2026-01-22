package history

import (
	"os"
	"path/filepath"
)

const (
	DefaultHistoryFile = ".sqlc_history"
)

// GetHistoryFilePath returns the path to the history file.
func GetHistoryFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return DefaultHistoryFile
	}
	return filepath.Join(home, DefaultHistoryFile)
}

// EnsureHistoryDir ensures the directory for the history file exists.
func EnsureHistoryDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}
