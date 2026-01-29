package placeholder

import (
	"testing"
)

func TestCountPlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "no placeholders",
			query:    "SELECT * FROM users",
			expected: 0,
		},
		{
			name:     "single placeholder",
			query:    "SELECT * FROM users WHERE id = ?",
			expected: 1,
		},
		{
			name:     "multiple placeholders",
			query:    "SELECT * FROM users WHERE status = ? AND age > ?",
			expected: 2,
		},
		{
			name:     "placeholder in single quoted string excluded",
			query:    "SELECT * FROM users WHERE name = '?'",
			expected: 0,
		},
		{
			name:     "placeholder in double quoted string excluded",
			query:    `SELECT * FROM users WHERE name = "?"`,
			expected: 0,
		},
		{
			name:     "placeholder in backtick excluded",
			query:    "SELECT * FROM `table?` WHERE id = ?",
			expected: 1,
		},
		{
			name:     "mixed string and real placeholder",
			query:    "SELECT * FROM users WHERE name = 'test?' AND id = ?",
			expected: 1,
		},
		{
			name:     "escaped quote does not close string",
			query:    `SELECT * FROM users WHERE name = 'it\'s a test' AND id = ?`,
			expected: 1,
		},
		{
			name:     "INSERT with multiple placeholders",
			query:    "INSERT INTO logs (level, message) VALUES (?, ?)",
			expected: 2,
		},
		{
			name:     "complex query with strings and placeholders",
			query:    "SELECT * FROM users WHERE status = 'active' AND role = ? AND age > ?",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountPlaceholders(tt.query)
			if result != tt.expected {
				t.Errorf("CountPlaceholders(%q) = %d, want %d", tt.query, result, tt.expected)
			}
		})
	}
}

func TestParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "empty string returns nil",
			input:    "",
			expected: nil,
		},
		{
			name:     "NULL returns nil",
			input:    "NULL",
			expected: nil,
		},
		{
			name:     "null lowercase returns nil",
			input:    "null",
			expected: nil,
		},
		{
			name:     "Null mixed case returns nil",
			input:    "Null",
			expected: nil,
		},
		{
			name:     "integer value",
			input:    "123",
			expected: int64(123),
		},
		{
			name:     "negative integer",
			input:    "-42",
			expected: int64(-42),
		},
		{
			name:     "float value",
			input:    "3.14",
			expected: float64(3.14),
		},
		{
			name:     "negative float",
			input:    "-2.5",
			expected: float64(-2.5),
		},
		{
			name:     "string value",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "string with spaces trimmed",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "string that looks like number but has letters",
			input:    "123abc",
			expected: "123abc",
		},
		{
			name:     "zero integer",
			input:    "0",
			expected: int64(0),
		},
		{
			name:     "zero float",
			input:    "0.0",
			expected: float64(0.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseValue(tt.input)
			if result != tt.expected {
				t.Errorf("ParseValue(%q) = %v (%T), want %v (%T)",
					tt.input, result, result, tt.expected, tt.expected)
			}
		})
	}
}
