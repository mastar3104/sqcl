package render

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/mastar3104/sqcl/internal/db"
)

// CSVRenderer renders query results as CSV format.
type CSVRenderer struct {
	NullDisplay string
}

// NewCSVRenderer creates a new CSV renderer.
func NewCSVRenderer() *CSVRenderer {
	return &CSVRenderer{
		NullDisplay: "NULL",
	}
}

// Render formats a query result as CSV string.
func (r *CSVRenderer) Render(result *db.QueryResult) string {
	if !result.IsSelect {
		return r.renderExecResult(result)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	writer.Write(result.Columns)

	// Write rows
	for _, row := range result.Rows {
		record := make([]string, len(row))
		for i, val := range row {
			record[i] = r.formatValue(val)
		}
		writer.Write(record)
	}

	writer.Flush()

	output := buf.String()
	output += fmt.Sprintf("\n%d row(s) in set", len(result.Rows))
	output += fmt.Sprintf("\nTime: %v", result.Duration)

	return output
}

func (r *CSVRenderer) formatValue(val interface{}) string {
	if val == nil {
		return r.NullDisplay
	}

	switch v := val.(type) {
	case []byte:
		return string(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (r *CSVRenderer) renderExecResult(result *db.QueryResult) string {
	output := fmt.Sprintf("Query OK, %d row(s) affected", result.RowsAffected)
	if result.LastInsertID > 0 {
		output += fmt.Sprintf(" (last insert ID: %d)", result.LastInsertID)
	}
	output += fmt.Sprintf("\nTime: %v", result.Duration)
	return output
}
