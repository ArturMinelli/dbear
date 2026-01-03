package importer

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"dbear/internal/config"
)

// EnvImporter implements the Importer interface for .env files
type EnvImporter struct{}

// EnvImportOptions contains options specific to env file imports
type EnvImportOptions struct {
	ImportOptions
	ConnectionStringKey string
	HostKey             string
	PortKey             string
	DatabaseKey         string
	UsernameKey         string
	PasswordKey         string
}

// NewEnvImporter creates a new EnvImporter instance
func NewEnvImporter() *EnvImporter {
	return &EnvImporter{}
}

// Import imports connections from an .env file
func (e *EnvImporter) Import(filePath string, options ImportOptions) ([]config.Connection, error) {
	envMap, err := ParseEnvFile(filePath)
	if err != nil {
		return nil, err
	}

	envOptions := EnvImportOptions{
		ImportOptions: options,
		HostKey:       "DB_HOST",
		PortKey:       "DB_PORT",
		DatabaseKey:   "DB_DATABASE",
		UsernameKey:   "DB_USERNAME",
		PasswordKey:   "DB_PASSWORD",
	}

	// Check if connection string mode is being used
	if envOptions.ConnectionStringKey != "" {
		return e.importFromConnectionString(envMap, envOptions)
	}

	return e.importFromVariables(envMap, envOptions)
}

// ImportWithOptions imports connections from an .env file with full env-specific options
func (e *EnvImporter) ImportWithOptions(filePath string, options EnvImportOptions) ([]config.Connection, error) {
	envMap, err := ParseEnvFile(filePath)
	if err != nil {
		return nil, err
	}

	// Check if connection string mode is being used
	if options.ConnectionStringKey != "" {
		return e.importFromConnectionString(envMap, options)
	}

	return e.importFromVariables(envMap, options)
}

// importFromConnectionString imports a connection from a connection string in an env variable
func (e *EnvImporter) importFromConnectionString(envMap map[string]string, options EnvImportOptions) ([]config.Connection, error) {
	connString := GetEnvValue(envMap, options.ConnectionStringKey)
	if connString == "" {
		return nil, fmt.Errorf("connection string key '%s' not found or empty in env file", options.ConnectionStringKey)
	}

	conn, err := parseConnectionString(connString, options.DatabaseType)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	conn.Name = options.ConnectionName
	if conn.Type == "" {
		conn.Type = options.DatabaseType
	}

	return []config.Connection{conn}, nil
}

// importFromVariables imports a connection from multiple env variables
func (e *EnvImporter) importFromVariables(envMap map[string]string, options EnvImportOptions) ([]config.Connection, error) {
	if options.DatabaseType == "" {
		return nil, fmt.Errorf("database type is required when importing from variables")
	}

	host := GetEnvValue(envMap, options.HostKey)
	database := GetEnvValue(envMap, options.DatabaseKey)
	username := GetEnvValue(envMap, options.UsernameKey)
	password := GetEnvValue(envMap, options.PasswordKey)

	portStr := GetEnvValue(envMap, options.PortKey)
	port := 0
	if portStr != "" {
		parsedPort, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid port value '%s': %w", portStr, err)
		}
		port = parsedPort
	} else {
		// Set default ports based on database type
		switch options.DatabaseType {
		case config.TypePostgreSQL:
			port = 5432
		case config.TypeMySQL:
			port = 3306
		case config.TypeSQLite:
			port = 0
		}
	}

	// SQLite doesn't need host/port/username/password
	if options.DatabaseType == config.TypeSQLite {
		if database == "" {
			return nil, fmt.Errorf("database key '%s' is required for SQLite", options.DatabaseKey)
		}
		host = "localhost"
	}

	conn := config.Connection{
		Name:     options.ConnectionName,
		Type:     options.DatabaseType,
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}

	return []config.Connection{conn}, nil
}

// parseConnectionString parses a standard database URL connection string
func parseConnectionString(connString string, defaultType string) (config.Connection, error) {
	parsedURL, err := url.Parse(connString)
	if err != nil {
		return config.Connection{}, fmt.Errorf("invalid connection string format: %w", err)
	}

	dbType := defaultType
	scheme := strings.ToLower(parsedURL.Scheme)

	// Determine database type from scheme
	switch {
	case scheme == "postgres" || scheme == "postgresql":
		dbType = config.TypePostgreSQL
	case scheme == "mysql":
		dbType = config.TypeMySQL
	case scheme == "sqlite" || scheme == "sqlite3":
		dbType = config.TypeSQLite
	default:
		if defaultType == "" {
			return config.Connection{}, fmt.Errorf("unsupported database scheme: %s", scheme)
		}
		dbType = defaultType
	}

	conn := config.Connection{
		Type: dbType,
	}

	// Extract user info
	if parsedURL.User != nil {
		conn.Username = parsedURL.User.Username()
		if password, ok := parsedURL.User.Password(); ok {
			conn.Password = password
		}
	}

	// Extract host and port
	if parsedURL.Host != "" {
		hostParts := strings.Split(parsedURL.Host, ":")
		conn.Host = hostParts[0]
		if len(hostParts) > 1 {
			port, err := strconv.Atoi(hostParts[1])
			if err != nil {
				return config.Connection{}, fmt.Errorf("invalid port in connection string: %w", err)
			}
			conn.Port = port
		} else {
			// Set default ports
			switch dbType {
			case config.TypePostgreSQL:
				conn.Port = 5432
			case config.TypeMySQL:
				conn.Port = 3306
			case config.TypeSQLite:
				conn.Port = 0
			}
		}
	}

	// Extract database name (path for SQLite, database name for others)
	if parsedURL.Path != "" {
		// Remove leading slash
		dbPath := strings.TrimPrefix(parsedURL.Path, "/")
		conn.Database = dbPath
	}

	// SQLite specific handling
	if dbType == config.TypeSQLite {
		if conn.Host == "" {
			conn.Host = "localhost"
		}
		conn.Port = 0
	}

	return conn, nil
}

