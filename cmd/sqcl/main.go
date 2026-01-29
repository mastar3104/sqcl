package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mastar3104/sqcl/internal/app"
	"github.com/mastar3104/sqcl/internal/connections"
	"github.com/mastar3104/sqcl/internal/history"
)

var (
	version = "dev"
)

func main() {
	// Check for subcommands first
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "save":
			runSave(os.Args[2:])
			return
		case "list":
			runList()
			return
		case "remove":
			runRemove(os.Args[2:])
			return
		}
	}

	// Default: run connect
	runConnect()
}

func runSave(args []string) {
	fs := flag.NewFlagSet("save", flag.ExitOnError)
	dsn := fs.String("dsn", "", "Database connection string")
	driver := fs.String("driver", "mysql", "Database driver")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: sqcl save <name> -dsn 'connection_string' [-driver mysql]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}

	if len(args) < 1 {
		fs.Usage()
		os.Exit(1)
	}

	name := args[0]
	if err := fs.Parse(args[1:]); err != nil {
		os.Exit(1)
	}

	if *dsn == "" {
		fmt.Fprintln(os.Stderr, "Error: -dsn flag is required")
		fs.Usage()
		os.Exit(1)
	}

	manager := connections.NewManager()
	conn := connections.Connection{
		Name:   name,
		Driver: *driver,
		DSN:    *dsn,
	}

	if err := manager.Save(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving connection: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Connection '%s' saved successfully.\n", name)
}

func runList() {
	manager := connections.NewManager()
	conns, err := manager.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing connections: %v\n", err)
		os.Exit(1)
	}

	if len(conns) == 0 {
		fmt.Println("No saved connections.")
		return
	}

	fmt.Println("Saved connections:")
	for _, c := range conns {
		fmt.Printf("  - %s (%s)\n", c.Name, c.Driver)
	}
}

func runRemove(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: sqcl remove <name>\n")
		os.Exit(1)
	}

	name := args[0]
	manager := connections.NewManager()

	if err := manager.Remove(name); err != nil {
		if err == connections.ErrConnectionNotFound {
			fmt.Fprintf(os.Stderr, "Error: connection '%s' not found\n", name)
		} else {
			fmt.Fprintf(os.Stderr, "Error removing connection: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("Connection '%s' removed successfully.\n", name)
}

func runConnect() {
	// Parse command line flags
	dsn := flag.String("dsn", "", "Database connection string (e.g., 'user:pass@tcp(host:port)/dbname')")
	driver := flag.String("driver", "mysql", "Database driver (mysql)")
	connectionName := flag.String("c", "", "Use saved connection by name")
	historyFile := flag.String("history", history.GetHistoryFilePath(), "History file path")
	cacheTTL := flag.Duration("cache-ttl", 60*time.Second, "Metadata cache TTL")
	showVersion := flag.Bool("version", false, "Show version and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "sqcl - Terminal SQL Client\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  sqcl -dsn 'user:pass@tcp(host:port)/dbname'\n")
		fmt.Fprintf(os.Stderr, "  sqcl -c <connection_name>\n")
		fmt.Fprintf(os.Stderr, "  sqcl save <name> -dsn 'connection_string'\n")
		fmt.Fprintf(os.Stderr, "  sqcl list\n")
		fmt.Fprintf(os.Stderr, "  sqcl remove <name>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("sqcl version %s\n", version)
		os.Exit(0)
	}

	// Resolve connection from saved name if specified
	actualDSN := *dsn
	actualDriver := *driver

	if *connectionName != "" {
		manager := connections.NewManager()
		conn, err := manager.Get(*connectionName)
		if err != nil {
			if err == connections.ErrConnectionNotFound {
				fmt.Fprintf(os.Stderr, "Error: connection '%s' not found\n", *connectionName)
				fmt.Fprintln(os.Stderr, "Use 'sqcl list' to see saved connections.")
			} else {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
			os.Exit(1)
		}
		actualDSN = conn.DSN
		actualDriver = conn.Driver
	}

	if actualDSN == "" {
		fmt.Fprintln(os.Stderr, "Error: -dsn flag or -c flag is required")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	// Create and run application
	config := app.Config{
		DSN:         actualDSN,
		Driver:      actualDriver,
		HistoryFile: *historyFile,
		CacheTTL:    *cacheTTL,
	}

	application, err := app.New(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
