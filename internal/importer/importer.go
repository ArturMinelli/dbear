package importer

import "dbear/internal/config"

// ImportOptions contains configuration for importing connections
type ImportOptions struct {
	ConnectionName string
	DatabaseType   string
}

// Importer defines the interface for all import types
type Importer interface {
	Import(filePath string, options ImportOptions) ([]config.Connection, error)
}

