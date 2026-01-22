package completion

import (
	"strings"
)

// CompletionContext represents the type of completion expected.
type CompletionContext int

const (
	ContextKeywordOrTable CompletionContext = iota
	ContextTable
	ContextColumn
	ContextDatabase
)

// tableContextKeywords are keywords that expect a table name to follow.
var tableContextKeywords = map[string]bool{
	"FROM":     true,
	"JOIN":     true,
	"UPDATE":   true,
	"INTO":     true,
	"TABLE":    true,
	"DESC":     true,
	"DESCRIBE": true,
	"TRUNCATE": true,
	"DROP":     true,
}

// columnContextKeywords are keywords that expect column names to follow.
var columnContextKeywords = map[string]bool{
	"SELECT":  true,
	"WHERE":   true,
	"AND":     true,
	"OR":      true,
	"ON":      true,
	"SET":     true,
	"ORDER":   true,
	"GROUP":   true,
	"HAVING":  true,
	"BY":      true,
	"BETWEEN": true,
}

// DetectContext analyzes the SQL input to determine what type of completion is needed.
func DetectContext(input string) CompletionContext {
	tokens := Tokenize(input)
	nonWS := GetNonWhitespaceTokens(tokens)

	if len(nonWS) == 0 {
		return ContextKeywordOrTable
	}

	// Get the last meaningful token(s)
	lastToken := nonWS[len(nonWS)-1]
	lastUpper := strings.ToUpper(lastToken.Value)

	// Check if the last token is a keyword that expects a table
	if tableContextKeywords[lastUpper] {
		return ContextTable
	}

	// Check for USE -> database context
	if lastUpper == "USE" {
		return ContextDatabase
	}

	// Check if we're after SELECT but before FROM (column context)
	if isInSelectClause(nonWS) {
		return ContextColumn
	}

	// Check if last token is a column context keyword
	if columnContextKeywords[lastUpper] {
		return ContextColumn
	}

	// Check previous token for context
	if len(nonWS) >= 2 {
		prevToken := nonWS[len(nonWS)-2]
		prevUpper := strings.ToUpper(prevToken.Value)

		// If previous token expects a table and current is a partial word
		if tableContextKeywords[prevUpper] {
			return ContextTable
		}

		if prevUpper == "USE" {
			return ContextDatabase
		}

		// If previous token expects a column (WHERE, AND, OR, etc.) and current is a partial word
		if columnContextKeywords[prevUpper] {
			return ContextColumn
		}

		// ORDER BY, GROUP BY context
		if prevUpper == "BY" && len(nonWS) >= 3 {
			thirdLast := strings.ToUpper(nonWS[len(nonWS)-3].Value)
			if thirdLast == "ORDER" || thirdLast == "GROUP" {
				return ContextColumn
			}
		}
	}

	// Default: keywords and tables
	return ContextKeywordOrTable
}

// isInSelectClause checks if we're in the SELECT clause (before FROM).
func isInSelectClause(tokens []Token) bool {
	foundSelect := false
	foundFrom := false

	for _, t := range tokens {
		upper := strings.ToUpper(t.Value)
		if upper == "SELECT" {
			foundSelect = true
		}
		if upper == "FROM" {
			foundFrom = true
		}
	}

	return foundSelect && !foundFrom
}

// GetTableContext extracts context about which table's columns should be suggested.
func GetTableContext(input string) []string {
	tokens := Tokenize(input)
	nonWS := GetNonWhitespaceTokens(tokens)

	var tables []string
	for i, t := range nonWS {
		upper := strings.ToUpper(t.Value)
		if tableContextKeywords[upper] && i+1 < len(nonWS) {
			nextToken := nonWS[i+1]
			if nextToken.Type == TokenWord {
				tableName := strings.Trim(nextToken.Value, "`")
				tables = append(tables, tableName)
			}
		}
	}

	return tables
}
