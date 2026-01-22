package cache

import (
	"context"
	"sync"
	"time"

	"github.com/mastar3104/sqcl/internal/db"
)

// MetadataCache provides TTL-based caching for database metadata.
type MetadataCache struct {
	provider db.MetadataProvider
	ttl      time.Duration

	mu            sync.RWMutex
	tables        []string
	tablesExpiry  time.Time
	columns       map[string][]db.ColumnInfo
	columnsExpiry map[string]time.Time
	databases     []string
	databasesExp  time.Time
}

// NewMetadataCache creates a new metadata cache with the specified TTL.
func NewMetadataCache(provider db.MetadataProvider, ttl time.Duration) *MetadataCache {
	return &MetadataCache{
		provider:      provider,
		ttl:           ttl,
		columns:       make(map[string][]db.ColumnInfo),
		columnsExpiry: make(map[string]time.Time),
	}
}

// GetTables returns cached table names or fetches them if expired.
func (c *MetadataCache) GetTables(ctx context.Context) ([]string, error) {
	c.mu.RLock()
	if time.Now().Before(c.tablesExpiry) && c.tables != nil {
		tables := c.tables
		c.mu.RUnlock()
		return tables, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if time.Now().Before(c.tablesExpiry) && c.tables != nil {
		return c.tables, nil
	}

	tables, err := c.provider.GetTables(ctx)
	if err != nil {
		return nil, err
	}

	c.tables = tables
	c.tablesExpiry = time.Now().Add(c.ttl)
	return tables, nil
}

// GetColumns returns cached column info or fetches if expired.
func (c *MetadataCache) GetColumns(ctx context.Context, tableName string) ([]db.ColumnInfo, error) {
	c.mu.RLock()
	if expiry, ok := c.columnsExpiry[tableName]; ok && time.Now().Before(expiry) {
		cols := c.columns[tableName]
		c.mu.RUnlock()
		return cols, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if expiry, ok := c.columnsExpiry[tableName]; ok && time.Now().Before(expiry) {
		return c.columns[tableName], nil
	}

	cols, err := c.provider.GetColumns(ctx, tableName)
	if err != nil {
		return nil, err
	}

	c.columns[tableName] = cols
	c.columnsExpiry[tableName] = time.Now().Add(c.ttl)
	return cols, nil
}

// GetDatabases returns cached database names or fetches if expired.
func (c *MetadataCache) GetDatabases(ctx context.Context) ([]string, error) {
	c.mu.RLock()
	if time.Now().Before(c.databasesExp) && c.databases != nil {
		dbs := c.databases
		c.mu.RUnlock()
		return dbs, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if time.Now().Before(c.databasesExp) && c.databases != nil {
		return c.databases, nil
	}

	dbs, err := c.provider.GetDatabases(ctx)
	if err != nil {
		return nil, err
	}

	c.databases = dbs
	c.databasesExp = time.Now().Add(c.ttl)
	return dbs, nil
}

// GetAllColumns returns all cached columns for all known tables.
func (c *MetadataCache) GetAllColumns(ctx context.Context) ([]string, error) {
	tables, err := c.GetTables(ctx)
	if err != nil {
		return nil, err
	}

	columnSet := make(map[string]struct{})
	for _, table := range tables {
		cols, err := c.GetColumns(ctx, table)
		if err != nil {
			continue
		}
		for _, col := range cols {
			columnSet[col.Name] = struct{}{}
		}
	}

	columns := make([]string, 0, len(columnSet))
	for col := range columnSet {
		columns = append(columns, col)
	}
	return columns, nil
}

// Reload clears all cached data, forcing refresh on next access.
func (c *MetadataCache) Reload() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.tables = nil
	c.tablesExpiry = time.Time{}
	c.columns = make(map[string][]db.ColumnInfo)
	c.columnsExpiry = make(map[string]time.Time)
	c.databases = nil
	c.databasesExp = time.Time{}
}
