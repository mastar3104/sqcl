package mysql

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/mastar3104/sqcl/internal/db"

	_ "github.com/go-sql-driver/mysql"
)

// Connector implements db.Connector for MySQL.
type Connector struct {
	db *sql.DB
}

// NewConnector creates a new MySQL connector.
func NewConnector() *Connector {
	return &Connector{}
}

// Connect establishes a connection to the MySQL database.
func (c *Connector) Connect(ctx context.Context, dsn string) error {
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(time.Minute * 5)

	if err := database.PingContext(ctx); err != nil {
		database.Close()
		return err
	}

	c.db = database
	return nil
}

// Close closes the database connection.
func (c *Connector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Ping checks if the database connection is alive.
func (c *Connector) Ping(ctx context.Context) error {
	if c.db == nil {
		return sql.ErrConnDone
	}
	return c.db.PingContext(ctx)
}

// DB returns the underlying *sql.DB.
func (c *Connector) DB() *sql.DB {
	return c.db
}

// Execute runs a SQL query and returns the result.
func (c *Connector) Execute(ctx context.Context, query string) (*db.QueryResult, error) {
	start := time.Now()

	trimmed := strings.TrimSpace(strings.ToUpper(query))
	isSelect := strings.HasPrefix(trimmed, "SELECT") ||
		strings.HasPrefix(trimmed, "SHOW") ||
		strings.HasPrefix(trimmed, "DESCRIBE") ||
		strings.HasPrefix(trimmed, "DESC") ||
		strings.HasPrefix(trimmed, "EXPLAIN")

	if isSelect {
		return c.executeQuery(ctx, query, start)
	}
	return c.executeExec(ctx, query, start)
}

func (c *Connector) executeQuery(ctx context.Context, query string, start time.Time) (*db.QueryResult, error) {
	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var resultRows [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make([]interface{}, len(columns))
		for i, v := range values {
			if b, ok := v.([]byte); ok {
				row[i] = string(b)
			} else {
				row[i] = v
			}
		}
		resultRows = append(resultRows, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &db.QueryResult{
		Columns:  columns,
		Rows:     resultRows,
		IsSelect: true,
		Duration: time.Since(start),
	}, nil
}

func (c *Connector) executeExec(ctx context.Context, query string, start time.Time) (*db.QueryResult, error) {
	result, err := c.db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	lastInsertID, _ := result.LastInsertId()

	return &db.QueryResult{
		RowsAffected: rowsAffected,
		LastInsertID: lastInsertID,
		IsSelect:     false,
		Duration:     time.Since(start),
	}, nil
}

// GetCurrentDatabase returns the current database name.
func (c *Connector) GetCurrentDatabase(ctx context.Context) (string, error) {
	var name sql.NullString
	err := c.db.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&name)
	if err != nil {
		return "", err
	}
	return name.String, nil
}
