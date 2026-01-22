package render

// OutputFormat represents the output format type.
type OutputFormat int

const (
	FormatTable OutputFormat = iota
	FormatCSV
	FormatJSON
)

// FormatConfig holds configuration for output formatting.
type FormatConfig struct {
	Format      OutputFormat
	NullDisplay string
	MaxColWidth int
	ShowHeaders bool
}

// DefaultFormatConfig returns the default format configuration.
func DefaultFormatConfig() FormatConfig {
	return FormatConfig{
		Format:      FormatTable,
		NullDisplay: "NULL",
		MaxColWidth: 50,
		ShowHeaders: true,
	}
}
