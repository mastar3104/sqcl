package render

import (
	"encoding/json"
	"fmt"

	"github.com/mastar3104/sqcl/internal/db"
)

// JSONRenderer renders query results as JSON format.
type JSONRenderer struct {
	Pretty bool
}

// NewJSONRenderer creates a new JSON renderer.
func NewJSONRenderer() *JSONRenderer {
	return &JSONRenderer{
		Pretty: true,
	}
}

// Render formats a query result as JSON string.
func (r *JSONRenderer) Render(result *db.QueryResult) string {
	if !result.IsSelect {
		return r.renderExecResult(result)
	}

	// Convert rows to array of maps
	rows := make([]map[string]interface{}, 0, len(result.Rows))
	for _, row := range result.Rows {
		record := make(map[string]interface{})
		for i, col := range result.Columns {
			record[col] = row[i]
		}
		rows = append(rows, record)
	}

	var jsonBytes []byte
	var err error
	if r.Pretty {
		jsonBytes, err = json.MarshalIndent(rows, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(rows)
	}

	if err != nil {
		return fmt.Sprintf("Error encoding JSON: %v", err)
	}

	output := string(jsonBytes)
	output += fmt.Sprintf("\n\n%d row(s) in set", len(result.Rows))
	output += fmt.Sprintf("\nTime: %v", result.Duration)

	return output
}

func (r *JSONRenderer) renderExecResult(result *db.QueryResult) string {
	output := fmt.Sprintf("Query OK, %d row(s) affected", result.RowsAffected)
	if result.LastInsertID > 0 {
		output += fmt.Sprintf(" (last insert ID: %d)", result.LastInsertID)
	}
	output += fmt.Sprintf("\nTime: %v", result.Duration)
	return output
}
