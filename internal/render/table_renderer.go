package render

import (
	"fmt"
	"strings"

	"github.com/mastar3104/sqcl/internal/db"
)

// TableRenderer renders query results as formatted tables.
type TableRenderer struct {
	NullDisplay string
	MaxColWidth int
}

// NewTableRenderer creates a new table renderer with default settings.
func NewTableRenderer() *TableRenderer {
	return &TableRenderer{
		NullDisplay: "NULL",
		MaxColWidth: 50,
	}
}

// Render formats a query result as a table string.
func (r *TableRenderer) Render(result *db.QueryResult) string {
	var sb strings.Builder

	if result.IsSelect {
		r.renderTable(&sb, result)
	} else {
		r.renderExecResult(&sb, result)
	}

	// Add execution time
	sb.WriteString(fmt.Sprintf("\nTime: %v", result.Duration))

	return sb.String()
}

func (r *TableRenderer) renderTable(sb *strings.Builder, result *db.QueryResult) {
	if len(result.Columns) == 0 {
		sb.WriteString("Empty result set")
		return
	}

	// Calculate column widths
	widths := make([]int, len(result.Columns))
	for i, col := range result.Columns {
		widths[i] = len(col)
	}

	for _, row := range result.Rows {
		for i, val := range row {
			str := r.formatValue(val)
			if len(str) > widths[i] {
				widths[i] = len(str)
			}
		}
	}

	// Build separator line
	separator := r.buildSeparator(widths)

	// Render header
	sb.WriteString(separator)
	sb.WriteString("\n")
	r.renderRow(sb, result.Columns, widths, true)
	sb.WriteString(separator)
	sb.WriteString("\n")

	// Render data rows
	for _, row := range result.Rows {
		strRow := make([]string, len(row))
		for i, val := range row {
			strRow[i] = r.formatValue(val)
		}
		r.renderRow(sb, strRow, widths, false)
	}

	sb.WriteString(separator)
	sb.WriteString("\n")

	// Row count
	sb.WriteString(fmt.Sprintf("%d row(s) in set", len(result.Rows)))
}

func (r *TableRenderer) renderRow(sb *strings.Builder, values []string, widths []int, isHeader bool) {
	sb.WriteString("|")
	for i, val := range values {
		display := val

		// Pad the value
		padding := widths[i] - len(display)
		sb.WriteString(" ")
		sb.WriteString(display)
		sb.WriteString(strings.Repeat(" ", padding))
		sb.WriteString(" |")
	}
	sb.WriteString("\n")
}

func (r *TableRenderer) buildSeparator(widths []int) string {
	var sb strings.Builder
	sb.WriteString("+")
	for _, w := range widths {
		sb.WriteString(strings.Repeat("-", w+2))
		sb.WriteString("+")
	}
	return sb.String()
}

func (r *TableRenderer) formatValue(val interface{}) string {
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

func (r *TableRenderer) renderExecResult(sb *strings.Builder, result *db.QueryResult) {
	sb.WriteString(fmt.Sprintf("Query OK, %d row(s) affected", result.RowsAffected))
	if result.LastInsertID > 0 {
		sb.WriteString(fmt.Sprintf(" (last insert ID: %d)", result.LastInsertID))
	}
}
