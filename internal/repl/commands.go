package repl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mastar3104/sqcl/internal/cache"
	"github.com/mastar3104/sqcl/internal/db"
	"github.com/mastar3104/sqcl/internal/render"
)

// CommandHandler handles internal REPL commands.
type CommandHandler struct {
	connector db.Connector
	cache     *cache.MetadataCache
	renderer  *render.TableRenderer
}

// NewCommandHandler creates a new command handler.
func NewCommandHandler(connector db.Connector, cache *cache.MetadataCache, renderer *render.TableRenderer) *CommandHandler {
	return &CommandHandler{
		connector: connector,
		cache:     cache,
		renderer:  renderer,
	}
}

// CommandResult holds the result of a command execution.
type CommandResult struct {
	Output   string
	ShouldQuit bool
	Error    error
}

// Execute runs an internal command.
func (h *CommandHandler) Execute(cmd string, args []string) CommandResult {
	switch cmd {
	case "help", "h", "?":
		return h.helpCommand()
	case "quit", "q", "exit":
		return CommandResult{ShouldQuit: true}
	case "reload", "refresh":
		return h.reloadCommand()
	case "tables":
		return h.tablesCommand()
	case "columns", "cols":
		return h.columnsCommand(args)
	case "databases", "dbs":
		return h.databasesCommand()
	case "status":
		return h.statusCommand()
	default:
		return CommandResult{Error: fmt.Errorf("unknown command: %s (type :help for available commands)", cmd)}
	}
}

func (h *CommandHandler) helpCommand() CommandResult {
	help := `Available commands:
  :help, :h, :?       Show this help message
  :quit, :q, :exit    Exit the client
  :reload, :refresh   Reload metadata cache
  :tables             List all tables
  :columns <table>    Show columns for a table
  :databases, :dbs    List all databases
  :status             Show connection status

SQL queries must end with a semicolon (;)
Use TAB for auto-completion
Use Ctrl+C to cancel current input, Ctrl+D to exit`
	return CommandResult{Output: help}
}

func (h *CommandHandler) reloadCommand() CommandResult {
	h.cache.Reload()
	return CommandResult{Output: "Metadata cache reloaded"}
}

func (h *CommandHandler) tablesCommand() CommandResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tables, err := h.cache.GetTables(ctx)
	if err != nil {
		return CommandResult{Error: fmt.Errorf("failed to get tables: %w", err)}
	}

	if len(tables) == 0 {
		return CommandResult{Output: "No tables found"}
	}

	result := &db.QueryResult{
		Columns:  []string{"Table"},
		IsSelect: true,
	}
	for _, t := range tables {
		result.Rows = append(result.Rows, []interface{}{t})
	}

	output := h.renderer.Render(result)
	return CommandResult{Output: output}
}

func (h *CommandHandler) columnsCommand(args []string) CommandResult {
	if len(args) == 0 {
		return CommandResult{Error: fmt.Errorf("usage: :columns <table_name>")}
	}

	tableName := args[0]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	columns, err := h.cache.GetColumns(ctx, tableName)
	if err != nil {
		return CommandResult{Error: fmt.Errorf("failed to get columns: %w", err)}
	}

	if len(columns) == 0 {
		return CommandResult{Output: fmt.Sprintf("No columns found for table '%s'", tableName)}
	}

	result := &db.QueryResult{
		Columns:  []string{"Column", "Type", "Nullable", "Key", "Default"},
		IsSelect: true,
	}
	for _, col := range columns {
		nullable := "NO"
		if col.IsNullable {
			nullable = "YES"
		}
		key := ""
		if col.IsPrimary {
			key = "PRI"
		}
		def := "NULL"
		if col.Default != nil {
			def = *col.Default
		}
		result.Rows = append(result.Rows, []interface{}{
			col.Name, col.DataType, nullable, key, def,
		})
	}

	output := h.renderer.Render(result)
	return CommandResult{Output: output}
}

func (h *CommandHandler) databasesCommand() CommandResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	databases, err := h.cache.GetDatabases(ctx)
	if err != nil {
		return CommandResult{Error: fmt.Errorf("failed to get databases: %w", err)}
	}

	if len(databases) == 0 {
		return CommandResult{Output: "No databases found"}
	}

	result := &db.QueryResult{
		Columns:  []string{"Database"},
		IsSelect: true,
	}
	for _, d := range databases {
		result.Rows = append(result.Rows, []interface{}{d})
	}

	output := h.renderer.Render(result)
	return CommandResult{Output: output}
}

func (h *CommandHandler) statusCommand() CommandResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := "Connected"
	if err := h.connector.Ping(ctx); err != nil {
		status = fmt.Sprintf("Disconnected (%v)", err)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Connection Status: %s\n", status))

	return CommandResult{Output: sb.String()}
}
