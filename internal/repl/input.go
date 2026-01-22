package repl

import (
	"strings"
)

// InputAccumulator handles multi-line SQL input accumulation.
type InputAccumulator struct {
	buffer strings.Builder
}

// NewInputAccumulator creates a new input accumulator.
func NewInputAccumulator() *InputAccumulator {
	return &InputAccumulator{}
}

// Add appends a line to the buffer.
func (a *InputAccumulator) Add(line string) {
	if a.buffer.Len() > 0 {
		a.buffer.WriteString("\n")
	}
	a.buffer.WriteString(line)
}

// IsComplete checks if the accumulated input is a complete SQL statement.
func (a *InputAccumulator) IsComplete() bool {
	content := strings.TrimSpace(a.buffer.String())
	if len(content) == 0 {
		return false
	}
	return strings.HasSuffix(content, ";")
}

// IsEmpty checks if the buffer is empty.
func (a *InputAccumulator) IsEmpty() bool {
	return a.buffer.Len() == 0
}

// Get returns the accumulated input.
func (a *InputAccumulator) Get() string {
	return a.buffer.String()
}

// GetTrimmed returns the accumulated input with trailing semicolon removed.
func (a *InputAccumulator) GetTrimmed() string {
	content := strings.TrimSpace(a.buffer.String())
	return strings.TrimSuffix(content, ";")
}

// Clear resets the buffer.
func (a *InputAccumulator) Clear() {
	a.buffer.Reset()
}

// IsInternalCommand checks if the input is an internal command (starts with :).
func IsInternalCommand(input string) bool {
	trimmed := strings.TrimSpace(input)
	return strings.HasPrefix(trimmed, ":")
}

// ParseCommand parses an internal command into name and arguments.
func ParseCommand(input string) (string, []string) {
	trimmed := strings.TrimSpace(input)
	if !strings.HasPrefix(trimmed, ":") {
		return "", nil
	}

	parts := strings.Fields(trimmed[1:])
	if len(parts) == 0 {
		return "", nil
	}

	return strings.ToLower(parts[0]), parts[1:]
}
