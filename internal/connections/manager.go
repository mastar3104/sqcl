package connections

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
	ErrConnectionExists   = errors.New("connection already exists")
)

// Connection represents a saved database connection.
type Connection struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}

// ConnectionsFile represents the structure of the connections file.
type ConnectionsFile struct {
	Connections []Connection `json:"connections"`
}

// Manager handles saving and loading connection configurations.
type Manager struct {
	filePath string
}

// NewManager creates a new connection manager.
func NewManager() *Manager {
	return &Manager{
		filePath: GetConnectionsFilePath(),
	}
}

// GetConnectionsFilePath returns the path to the connections file.
func GetConnectionsFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".sqcl/connections.json"
	}
	return filepath.Join(home, ".sqcl", "connections.json")
}

// EnsureConnectionsDir ensures the directory for the connections file exists.
func EnsureConnectionsDir() error {
	path := GetConnectionsFilePath()
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}

// Load reads all connections from the file.
func (m *Manager) Load() ([]Connection, error) {
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Connection{}, nil
		}
		return nil, err
	}

	var file ConnectionsFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}

	return file.Connections, nil
}

// save writes all connections to the file.
func (m *Manager) save(connections []Connection) error {
	if err := EnsureConnectionsDir(); err != nil {
		return err
	}

	file := ConnectionsFile{Connections: connections}
	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.filePath, data, 0600)
}

// Save adds or updates a connection.
func (m *Manager) Save(conn Connection) error {
	connections, err := m.Load()
	if err != nil {
		return err
	}

	// Check if connection with same name exists
	for i, c := range connections {
		if c.Name == conn.Name {
			// Update existing connection
			connections[i] = conn
			return m.save(connections)
		}
	}

	// Add new connection
	connections = append(connections, conn)
	return m.save(connections)
}

// Get retrieves a connection by name.
func (m *Manager) Get(name string) (Connection, error) {
	connections, err := m.Load()
	if err != nil {
		return Connection{}, err
	}

	for _, c := range connections {
		if c.Name == name {
			return c, nil
		}
	}

	return Connection{}, ErrConnectionNotFound
}

// Remove deletes a connection by name.
func (m *Manager) Remove(name string) error {
	connections, err := m.Load()
	if err != nil {
		return err
	}

	for i, c := range connections {
		if c.Name == name {
			connections = append(connections[:i], connections[i+1:]...)
			return m.save(connections)
		}
	}

	return ErrConnectionNotFound
}

// List returns all saved connections.
func (m *Manager) List() ([]Connection, error) {
	return m.Load()
}
