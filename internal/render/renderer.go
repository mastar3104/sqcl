package render

import "github.com/mastar3104/sqcl/internal/db"

// Renderer defines the interface for rendering query results.
type Renderer interface {
	Render(result *db.QueryResult) string
}
