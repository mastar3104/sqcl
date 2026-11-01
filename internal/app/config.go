package app

import (
	"time"

	"github.com/mastar3104/sqcl/internal/history"
)

// Config holds the application configuration.
type Config struct {
	DSN          string
	Driver       string
	HistoryFile  string
	CacheTTL     time.Duration
	Query        string // 指定時はREPLを起動せずこのSQLを実行して終了する
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Driver:      "mysql",
		HistoryFile: history.GetHistoryFilePath(),
		CacheTTL:    60 * time.Second,
	}
}

// Validate checks if the configuration is valid.
func (c Config) Validate() error {
	if c.DSN == "" {
		return ErrDSNRequired
	}
	return nil
}
