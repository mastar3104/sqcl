package completion

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/mastar3104/sqcl/internal/cache"
	"github.com/mastar3104/sqcl/internal/db"
)

// SQLCompleter implements readline.AutoCompleter for SQL input.
type SQLCompleter struct {
	cache   *cache.MetadataCache
	dialect db.Dialect
}

// NewSQLCompleter creates a new SQL completer.
func NewSQLCompleter(cache *cache.MetadataCache, dialect db.Dialect) *SQLCompleter {
	return &SQLCompleter{
		cache:   cache,
		dialect: dialect,
	}
}

// Do implements readline.AutoCompleter interface.
func (c *SQLCompleter) Do(line []rune, pos int) ([][]rune, int) {
	input := string(line[:pos])
	lastWord := GetLastWord(input)
	lastWordLower := strings.ToLower(lastWord)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	completionCtx := DetectContext(input)
	var candidates []string

	switch completionCtx {
	case ContextTable:
		candidates = c.getTableCandidates(ctx)
	case ContextColumn:
		candidates = c.getColumnCandidates(ctx, input)
		candidates = append(candidates, c.getKeywordCandidates()...)
	case ContextDatabase:
		candidates = c.getDatabaseCandidates(ctx)
	default:
		// Keywords + tables
		candidates = c.getKeywordCandidates()
		candidates = append(candidates, c.getTableCandidates(ctx)...)
	}

	// Filter candidates by prefix
	var matches []string
	for _, candidate := range candidates {
		candidateLower := strings.ToLower(candidate)
		if strings.HasPrefix(candidateLower, lastWordLower) {
			matches = append(matches, candidate)
		}
	}

	// Sort and deduplicate
	matches = uniqueSorted(matches)

	// Convert to readline format
	result := make([][]rune, len(matches))
	for i, match := range matches {
		suffix := match[len(lastWord):]
		result[i] = []rune(suffix)
	}

	return result, len(lastWord)
}

func (c *SQLCompleter) getTableCandidates(ctx context.Context) []string {
	tables, err := c.cache.GetTables(ctx)
	if err != nil {
		return nil
	}
	return tables
}

func (c *SQLCompleter) getColumnCandidates(ctx context.Context, input string) []string {
	// Get tables referenced in the query
	tables := GetTableContext(input)

	if len(tables) > 0 {
		// Get columns from specific tables
		columnSet := make(map[string]struct{})
		for _, table := range tables {
			cols, err := c.cache.GetColumns(ctx, table)
			if err != nil {
				continue
			}
			for _, col := range cols {
				columnSet[col.Name] = struct{}{}
			}
		}

		columns := make([]string, 0, len(columnSet))
		for col := range columnSet {
			columns = append(columns, col)
		}
		return columns
	}

	// No specific tables, get all columns
	columns, err := c.cache.GetAllColumns(ctx)
	if err != nil {
		return nil
	}
	return columns
}

func (c *SQLCompleter) getDatabaseCandidates(ctx context.Context) []string {
	dbs, err := c.cache.GetDatabases(ctx)
	if err != nil {
		return nil
	}
	return dbs
}

func (c *SQLCompleter) getKeywordCandidates() []string {
	return c.dialect.Keywords()
}

func uniqueSorted(items []string) []string {
	if len(items) == 0 {
		return items
	}

	seen := make(map[string]struct{})
	var result []string

	for _, item := range items {
		lower := strings.ToLower(item)
		if _, exists := seen[lower]; !exists {
			seen[lower] = struct{}{}
			result = append(result, item)
		}
	}

	sort.Strings(result)
	return result
}

// Ensure SQLCompleter implements readline.AutoCompleter
var _ readline.AutoCompleter = (*SQLCompleter)(nil)
