package db

import "time"

// QueryResult holds the result of a SQL query execution.
type QueryResult struct {
	Columns      []string
	Rows         [][]interface{}
	RowsAffected int64
	LastInsertID int64
	IsSelect     bool
	Duration     time.Duration
}

// ColumnInfo holds metadata about a table column.
type ColumnInfo struct {
	Name       string
	DataType   string
	IsNullable bool
	IsPrimary  bool
	Default    *string
}

// TableInfo holds metadata about a database table.
type TableInfo struct {
	Name    string
	Columns []ColumnInfo
}
