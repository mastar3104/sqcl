package highlight

import (
	"strings"
	"unicode"

	"github.com/mastar3104/sqcl/internal/db"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m" // Keywords
	colorYellow = "\033[33m" // Strings
	colorGreen  = "\033[32m" // Numbers
)

// SQLHighlighter implements syntax highlighting for SQL queries.
type SQLHighlighter struct {
	keywords map[string]struct{}
}

// NewSQLHighlighter creates a new SQL highlighter.
func NewSQLHighlighter(dialect db.Dialect) *SQLHighlighter {
	keywords := make(map[string]struct{})
	for _, kw := range dialect.Keywords() {
		keywords[strings.ToUpper(kw)] = struct{}{}
	}
	return &SQLHighlighter{
		keywords: keywords,
	}
}

// Paint implements readline.Painter interface.
func (h *SQLHighlighter) Paint(line []rune, pos int) []rune {
	if len(line) == 0 {
		return line
	}

	input := string(line)
	result := h.highlight(input)
	return []rune(result)
}

func (h *SQLHighlighter) highlight(input string) string {
	var result strings.Builder
	i := 0
	runes := []rune(input)

	for i < len(runes) {
		ch := runes[i]

		// String literal (single quote)
		if ch == '\'' {
			start := i
			i++
			for i < len(runes) && runes[i] != '\'' {
				if runes[i] == '\\' && i+1 < len(runes) {
					i += 2
				} else {
					i++
				}
			}
			if i < len(runes) {
				i++ // include closing quote
			}
			result.WriteString(colorYellow)
			result.WriteString(string(runes[start:i]))
			result.WriteString(colorReset)
			continue
		}

		// String literal (double quote)
		if ch == '"' {
			start := i
			i++
			for i < len(runes) && runes[i] != '"' {
				if runes[i] == '\\' && i+1 < len(runes) {
					i += 2
				} else {
					i++
				}
			}
			if i < len(runes) {
				i++ // include closing quote
			}
			result.WriteString(colorYellow)
			result.WriteString(string(runes[start:i]))
			result.WriteString(colorReset)
			continue
		}

		// Number
		if unicode.IsDigit(ch) || (ch == '.' && i+1 < len(runes) && unicode.IsDigit(runes[i+1])) {
			start := i
			hasDecimal := ch == '.'
			i++
			for i < len(runes) && (unicode.IsDigit(runes[i]) || (!hasDecimal && runes[i] == '.')) {
				if runes[i] == '.' {
					hasDecimal = true
				}
				i++
			}
			result.WriteString(colorGreen)
			result.WriteString(string(runes[start:i]))
			result.WriteString(colorReset)
			continue
		}

		// Identifier or keyword
		if unicode.IsLetter(ch) || ch == '_' {
			start := i
			i++
			for i < len(runes) && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) || runes[i] == '_') {
				i++
			}
			word := string(runes[start:i])
			if _, isKeyword := h.keywords[strings.ToUpper(word)]; isKeyword {
				result.WriteString(colorCyan)
				result.WriteString(word)
				result.WriteString(colorReset)
			} else {
				result.WriteString(word)
			}
			continue
		}

		// Other characters
		result.WriteRune(ch)
		i++
	}

	return result.String()
}
