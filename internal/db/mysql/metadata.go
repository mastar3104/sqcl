package mysql

import (
	"context"
	"database/sql"

	"github.com/mastar3104/sqcl/internal/db"
)

// MetadataProvider implements db.MetadataProvider for MySQL.
type MetadataProvider struct {
	database *sql.DB
	dialect  *Dialect
}

// NewMetadataProvider creates a new MySQL metadata provider.
func NewMetadataProvider(database *sql.DB) *MetadataProvider {
	return &MetadataProvider{
		database: database,
		dialect:  NewDialect(),
	}
}

// GetTables returns a list of table names in the current database.
func (m *MetadataProvider) GetTables(ctx context.Context) ([]string, error) {
	rows, err := m.database.QueryContext(ctx, m.dialect.GetTablesQuery())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, rows.Err()
}

// GetColumns returns column information for a specific table.
func (m *MetadataProvider) GetColumns(ctx context.Context, tableName string) ([]db.ColumnInfo, error) {
	rows, err := m.database.QueryContext(ctx, m.dialect.GetColumnsQuery(tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []db.ColumnInfo
	for rows.Next() {
		var col db.ColumnInfo
		var isNullable string
		var columnKey string
		var defaultVal sql.NullString

		if err := rows.Scan(&col.Name, &col.DataType, &isNullable, &columnKey, &defaultVal); err != nil {
			return nil, err
		}

		col.IsNullable = isNullable == "YES"
		col.IsPrimary = columnKey == "PRI"
		if defaultVal.Valid {
			col.Default = &defaultVal.String
		}

		columns = append(columns, col)
	}

	return columns, rows.Err()
}

// GetDatabases returns a list of database names.
func (m *MetadataProvider) GetDatabases(ctx context.Context) ([]string, error) {
	rows, err := m.database.QueryContext(ctx, m.dialect.GetDatabasesQuery())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		databases = append(databases, dbName)
	}

	return databases, rows.Err()
}
