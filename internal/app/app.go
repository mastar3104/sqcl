package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mastar3104/sqcl/internal/cache"
	"github.com/mastar3104/sqcl/internal/db"
	"github.com/mastar3104/sqcl/internal/db/mysql"
	"github.com/mastar3104/sqcl/internal/history"
	"github.com/mastar3104/sqcl/internal/repl"
)

var (
	ErrDSNRequired      = errors.New("DSN is required")
	ErrUnsupportedDriver = errors.New("unsupported database driver")
)

// App represents the main application.
type App struct {
	config    Config
	connector db.Connector
	metadata  db.MetadataProvider
	dialect   db.Dialect
	cache     *cache.MetadataCache
	repl      *repl.REPL
}

// New creates a new application instance.
func New(config Config) (*App, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &App{config: config}, nil
}

// Run starts the application.
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	// Initialize components
	if err := a.initialize(ctx); err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}
	defer a.cleanup()

	// Run the REPL
	return a.repl.Run()
}

func (a *App) initialize(ctx context.Context) error {
	// Create connector based on driver
	var connector db.Connector
	var dialect db.Dialect

	switch a.config.Driver {
	case "mysql":
		connector = mysql.NewConnector()
		dialect = mysql.NewDialect()
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedDriver, a.config.Driver)
	}

	// Connect to database
	connectCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := connector.Connect(connectCtx, a.config.DSN); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	a.connector = connector

	// Create metadata provider
	switch a.config.Driver {
	case "mysql":
		a.metadata = mysql.NewMetadataProvider(connector.DB())
	}
	a.dialect = dialect

	// Create metadata cache
	a.cache = cache.NewMetadataCache(a.metadata, a.config.CacheTTL)

	// Ensure history directory exists
	if err := history.EnsureHistoryDir(a.config.HistoryFile); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not create history directory: %v\n", err)
	}

	// Create REPL
	r, err := repl.New(repl.Config{
		Connector:   a.connector,
		Cache:       a.cache,
		Dialect:     a.dialect,
		HistoryFile: a.config.HistoryFile,
	})
	if err != nil {
		return fmt.Errorf("failed to create REPL: %w", err)
	}
	a.repl = r

	// Preload metadata cache in background
	go func() {
		preloadCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _ = a.cache.GetTables(preloadCtx)
	}()

	return nil
}

func (a *App) cleanup() {
	if a.repl != nil {
		a.repl.Close()
	}
	if a.connector != nil {
		a.connector.Close()
	}
}
