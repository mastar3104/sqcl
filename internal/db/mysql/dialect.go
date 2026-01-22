package mysql

import "fmt"

// Dialect implements db.Dialect for MySQL.
type Dialect struct {
	keywords []string
}

// NewDialect creates a new MySQL dialect.
func NewDialect() *Dialect {
	return &Dialect{
		keywords: []string{
			"SELECT", "FROM", "WHERE", "AND", "OR", "NOT", "IN", "LIKE",
			"ORDER", "BY", "ASC", "DESC", "LIMIT", "OFFSET", "GROUP",
			"HAVING", "JOIN", "LEFT", "RIGHT", "INNER", "OUTER", "CROSS",
			"ON", "AS", "DISTINCT", "ALL", "UNION", "INTERSECT", "EXCEPT",
			"INSERT", "INTO", "VALUES", "UPDATE", "SET", "DELETE",
			"CREATE", "TABLE", "DATABASE", "INDEX", "VIEW", "DROP", "ALTER",
			"ADD", "COLUMN", "PRIMARY", "KEY", "FOREIGN", "REFERENCES",
			"UNIQUE", "CHECK", "DEFAULT", "NULL", "NOT", "AUTO_INCREMENT",
			"IF", "EXISTS", "CASE", "WHEN", "THEN", "ELSE", "END",
			"COUNT", "SUM", "AVG", "MIN", "MAX", "COALESCE", "NULLIF",
			"CAST", "CONVERT", "DATE", "TIME", "DATETIME", "TIMESTAMP",
			"YEAR", "MONTH", "DAY", "HOUR", "MINUTE", "SECOND",
			"TRUE", "FALSE", "BETWEEN", "IS", "ESCAPE",
			"SHOW", "TABLES", "DATABASES", "COLUMNS", "DESCRIBE", "EXPLAIN",
			"USE", "TRUNCATE", "RENAME", "TO", "GRANT", "REVOKE",
			"BEGIN", "COMMIT", "ROLLBACK", "TRANSACTION", "SAVEPOINT",
			"INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT",
			"FLOAT", "DOUBLE", "DECIMAL", "NUMERIC",
			"VARCHAR", "CHAR", "TEXT", "BLOB", "BINARY", "VARBINARY",
			"BOOLEAN", "BOOL", "ENUM", "JSON",
		},
	}
}

// Name returns the dialect name.
func (d *Dialect) Name() string {
	return "mysql"
}

// Keywords returns the list of SQL keywords.
func (d *Dialect) Keywords() []string {
	return d.keywords
}

// QuoteIdentifier quotes an identifier for MySQL.
func (d *Dialect) QuoteIdentifier(name string) string {
	return "`" + name + "`"
}

// GetTablesQuery returns the query to list tables.
func (d *Dialect) GetTablesQuery() string {
	return "SHOW TABLES"
}

// GetColumnsQuery returns the query to get column information for a table.
func (d *Dialect) GetColumnsQuery(tableName string) string {
	return fmt.Sprintf(`
		SELECT
			COLUMN_NAME,
			DATA_TYPE,
			IS_NULLABLE,
			COLUMN_KEY,
			COLUMN_DEFAULT
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_NAME = '%s'
		AND TABLE_SCHEMA = DATABASE()
		ORDER BY ORDINAL_POSITION
	`, tableName)
}

// GetDatabasesQuery returns the query to list databases.
func (d *Dialect) GetDatabasesQuery() string {
	return "SHOW DATABASES"
}
