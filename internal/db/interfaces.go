package db

import (
	"context"
	"database/sql"
)

// Connector defines the interface for database connection and query execution.
type Connector interface {
	Connect(ctx context.Context, dsn string) error
	Close() error
	Execute(ctx context.Context, query string) (*QueryResult, error)
	ExecuteWithParams(ctx context.Context, query string, args []interface{}) (*QueryResult, error)
	Ping(ctx context.Context) error
	DB() *sql.DB
	GetCurrentDatabase(ctx context.Context) (string, error)
}

// MetadataProvider defines the interface for retrieving database metadata.
type MetadataProvider interface {
	GetTables(ctx context.Context) ([]string, error)
	GetColumns(ctx context.Context, tableName string) ([]ColumnInfo, error)
	GetDatabases(ctx context.Context) ([]string, error)
}

// Dialect defines database-specific SQL syntax and keywords.
type Dialect interface {
	Name() string
	Keywords() []string
	QuoteIdentifier(name string) string
	GetTablesQuery() string
	GetColumnsQuery(tableName string) string
	GetDatabasesQuery() string
}
