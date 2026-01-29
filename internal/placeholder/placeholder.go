package placeholder

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

// CountPlaceholders counts the number of ? placeholders in a SQL query,
// excluding those inside string literals (single quotes, double quotes, backticks).
func CountPlaceholders(query string) int {
	count := 0
	inSingleQuote := false
	inDoubleQuote := false
	inBacktick := false
	escaped := false

	for i := 0; i < len(query); i++ {
		c := query[i]

		if escaped {
			escaped = false
			continue
		}

		if c == '\\' {
			escaped = true
			continue
		}

		switch c {
		case '\'':
			if !inDoubleQuote && !inBacktick {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote && !inBacktick {
				inDoubleQuote = !inDoubleQuote
			}
		case '`':
			if !inSingleQuote && !inDoubleQuote {
				inBacktick = !inBacktick
			}
		case '?':
			if !inSingleQuote && !inDoubleQuote && !inBacktick {
				count++
			}
		}
	}

	return count
}

// ParseValue converts user input to an appropriate Go type for database parameters.
// - Empty string or "NULL" (case-insensitive) -> nil
// - Valid integer -> int64
// - Valid float -> float64
// - Otherwise -> string
func ParseValue(input string) interface{} {
	trimmed := strings.TrimSpace(input)

	// Empty or NULL
	if trimmed == "" || strings.EqualFold(trimmed, "NULL") {
		return nil
	}

	// Try integer
	if i, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
		return i
	}

	// Try float
	if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return f
	}

	// Return as string
	return trimmed
}

// PromptForValues prompts the user to enter values for each placeholder.
// Returns the values slice, a cancelled flag, and any error.
func PromptForValues(rl *readline.Instance, count int, query string) ([]interface{}, bool, error) {
	if count == 0 {
		return nil, false, nil
	}

	fmt.Println()
	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Enter values for %d placeholder(s):\n", count)
	fmt.Println("  (Press Enter for NULL, Ctrl+C to cancel)")
	fmt.Println()

	values := make([]interface{}, count)

	// Save original prompt and restore later
	originalPrompt := rl.Config.Prompt
	defer rl.SetPrompt(originalPrompt)

	for i := 0; i < count; i++ {
		rl.SetPrompt(fmt.Sprintf("  [%d]> ", i+1))

		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			fmt.Println()
			return nil, true, nil
		}
		if err == io.EOF {
			fmt.Println()
			return nil, true, nil
		}
		if err != nil {
			return nil, false, fmt.Errorf("readline error: %w", err)
		}

		values[i] = ParseValue(line)
	}

	fmt.Println()
	return values, false, nil
}
