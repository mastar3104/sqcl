package repl

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/mastar3104/sqcl/internal/cache"
	"github.com/mastar3104/sqcl/internal/completion"
	"github.com/mastar3104/sqcl/internal/db"
	"github.com/mastar3104/sqcl/internal/highlight"
	"github.com/mastar3104/sqcl/internal/render"
)

// REPL represents the read-eval-print loop.
type REPL struct {
	connector      db.Connector
	cache          *cache.MetadataCache
	dialect        db.Dialect
	renderer       render.Renderer
	outputFormat   render.OutputFormat
	commandHandler *CommandHandler
	historyFile    string
	readline       *readline.Instance
}

// Config holds REPL configuration.
type Config struct {
	Connector   db.Connector
	Cache       *cache.MetadataCache
	Dialect     db.Dialect
	HistoryFile string
}

// New creates a new REPL instance.
func New(cfg Config) (*REPL, error) {
	completer := completion.NewSQLCompleter(cfg.Cache, cfg.Dialect)
	highlighter := highlight.NewSQLHighlighter(cfg.Dialect)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "sqcl> ",
		HistoryFile:     cfg.HistoryFile,
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
		Painter:             highlighter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize readline: %w", err)
	}

	r := &REPL{
		connector:    cfg.Connector,
		cache:        cfg.Cache,
		dialect:      cfg.Dialect,
		renderer:     render.NewTableRenderer(),
		outputFormat: render.FormatTable,
		historyFile:  cfg.HistoryFile,
		readline:     rl,
	}
	r.commandHandler = NewCommandHandler(r)

	return r, nil
}

// Run starts the REPL loop.
func (r *REPL) Run() error {
	defer r.readline.Close()

	fmt.Println("Welcome to sqcl - Terminal SQL Client")
	fmt.Println("Type :help for help, :quit to exit")
	fmt.Println()

	accumulator := NewInputAccumulator()

	for {
		prompt := r.getPrompt()
		if !accumulator.IsEmpty() {
			prompt = "   -> "
		}
		r.readline.SetPrompt(prompt)

		line, err := r.readline.Readline()
		if err == readline.ErrInterrupt {
			if accumulator.IsEmpty() {
				continue
			}
			accumulator.Clear()
			fmt.Println("Query cancelled")
			continue
		}
		if err == io.EOF {
			fmt.Println("\nGoodbye!")
			return nil
		}
		if err != nil {
			return fmt.Errorf("readline error: %w", err)
		}

		line = strings.TrimSpace(line)

		// Handle empty input
		if line == "" {
			if accumulator.IsEmpty() {
				continue
			}
			accumulator.Add("")
			continue
		}

		// Handle internal commands (only when not accumulating)
		if accumulator.IsEmpty() && IsInternalCommand(line) {
			cmd, args := ParseCommand(line)
			result := r.commandHandler.Execute(cmd, args)

			if result.Error != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", result.Error)
			} else if result.Output != "" {
				fmt.Println(result.Output)
			}

			if result.ShouldQuit {
				fmt.Println("Goodbye!")
				return nil
			}
			continue
		}

		// Accumulate SQL input
		accumulator.Add(line)

		// Check if statement is complete
		if accumulator.IsComplete() {
			query := accumulator.Get()
			accumulator.Clear()

			r.executeQuery(query)
		}
	}
}

func (r *REPL) executeQuery(query string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	result, err := r.connector.Execute(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	output := r.renderer.Render(result)
	fmt.Println(output)
}

// Close cleans up REPL resources.
func (r *REPL) Close() error {
	if r.readline != nil {
		return r.readline.Close()
	}
	return nil
}

// SetFormat changes the output format.
func (r *REPL) SetFormat(format render.OutputFormat) {
	r.outputFormat = format
	switch format {
	case render.FormatCSV:
		r.renderer = render.NewCSVRenderer()
	case render.FormatJSON:
		r.renderer = render.NewJSONRenderer()
	default:
		r.renderer = render.NewTableRenderer()
	}
}

// GetFormat returns the current output format.
func (r *REPL) GetFormat() render.OutputFormat {
	return r.outputFormat
}

// Cache returns the metadata cache.
func (r *REPL) Cache() *cache.MetadataCache {
	return r.cache
}

// Connector returns the database connector.
func (r *REPL) Connector() db.Connector {
	return r.connector
}

// Renderer returns the current renderer.
func (r *REPL) Renderer() render.Renderer {
	return r.renderer
}

func (r *REPL) getPrompt() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dbName, err := r.connector.GetCurrentDatabase(ctx)
	if err != nil || dbName == "" {
		return "sqcl> "
	}
	return fmt.Sprintf("sqcl(%s)> ", dbName)
}

func filterInput(r rune) (rune, bool) {
	switch r {
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
