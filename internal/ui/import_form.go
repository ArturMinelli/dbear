package ui

import (
	"dbear/internal/config"
	"dbear/internal/importer"
	"fmt"

	"github.com/charmbracelet/huh"
)

const (
	importMethodConnectionString = "connection_string"
	importMethodMultipleVars     = "multiple_vars"
)

// ImportEnvForm presents an interactive form to collect import configuration
func ImportEnvForm() (*importer.EnvImportOptions, error) {
	var connectionName string
	var importMethod string
	var databaseType string
	var connectionStringKey string
	var hostKey string
	var portKey string
	var databaseKey string
	var usernameKey string
	var passwordKey string

	// Step 1: Connection name
	nameForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Connection Name").
				Description("Enter a unique name for this connection").
				Value(&connectionName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("connection name is required")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm())

	if err := nameForm.Run(); err != nil {
		return nil, fmt.Errorf("failed to get connection name: %w", err)
	}

	// Step 2: Import method selection
	methodForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Import Method").
				Description("How is the connection configured in your .env file?").
				Options(
					huh.NewOption("Connection String (single env variable)", importMethodConnectionString),
					huh.NewOption("Multiple Variables (separate keys for host, port, etc.)", importMethodMultipleVars),
				).
				Value(&importMethod),
		),
	).WithTheme(huh.ThemeCharm())

	if err := methodForm.Run(); err != nil {
		return nil, fmt.Errorf("failed to get import method: %w", err)
	}

	// Step 3: Configuration based on method
	if importMethod == importMethodConnectionString {
		// Connection string method
		connStringForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Connection String Key").
					Description("The environment variable name containing the connection string (e.g., DATABASE_URL)").
					Value(&connectionStringKey).
					Placeholder("DATABASE_URL").
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("connection string key is required")
						}
						return nil
					}),
				huh.NewSelect[string]().
					Title("Database Type").
					Description("Database type (optional - can be inferred from connection string)").
					Options(
						huh.NewOption("Auto-detect from connection string", ""),
						huh.NewOption("PostgreSQL", config.TypePostgreSQL),
						huh.NewOption("MySQL", config.TypeMySQL),
						huh.NewOption("SQLite", config.TypeSQLite),
					).
					Value(&databaseType),
			),
		).WithTheme(huh.ThemeCharm())

		if err := connStringForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to get connection string configuration: %w", err)
		}
	} else {
		// Multiple variables method
		// First get database type (required)
		dbTypeForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Database Type").
					Description("Select the database type").
					Options(
						huh.NewOption("PostgreSQL", config.TypePostgreSQL),
						huh.NewOption("MySQL", config.TypeMySQL),
						huh.NewOption("SQLite", config.TypeSQLite),
					).
					Value(&databaseType).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("database type is required")
						}
						return nil
					}),
			),
		).WithTheme(huh.ThemeCharm())

		if err := dbTypeForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to get database type: %w", err)
		}

		// Then get the env variable keys
		keysForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Host Key").
					Description("Environment variable name for database host").
					Value(&hostKey).
					Placeholder("DB_HOST"),
				huh.NewInput().
					Title("Port Key").
					Description("Environment variable name for database port").
					Value(&portKey).
					Placeholder("DB_PORT"),
				huh.NewInput().
					Title("Database Key").
					Description("Environment variable name for database name").
					Value(&databaseKey).
					Placeholder("DB_DATABASE"),
				huh.NewInput().
					Title("Username Key").
					Description("Environment variable name for database username").
					Value(&usernameKey).
					Placeholder("DB_USERNAME"),
				huh.NewInput().
					Title("Password Key").
					Description("Environment variable name for database password").
					Value(&passwordKey).
					Placeholder("DB_PASSWORD"),
			),
		).WithTheme(huh.ThemeCharm())

		if err := keysForm.Run(); err != nil {
			return nil, fmt.Errorf("failed to get environment variable keys: %w", err)
		}

		// Set defaults if empty
		if hostKey == "" {
			hostKey = "DB_HOST"
		}
		if portKey == "" {
			portKey = "DB_PORT"
		}
		if databaseKey == "" {
			databaseKey = "DB_DATABASE"
		}
		if usernameKey == "" {
			usernameKey = "DB_USERNAME"
		}
		if passwordKey == "" {
			passwordKey = "DB_PASSWORD"
		}
	}

	// Build and return options
	options := importer.EnvImportOptions{
		ImportOptions: importer.ImportOptions{
			ConnectionName: connectionName,
			DatabaseType:   databaseType,
		},
		ConnectionStringKey: connectionStringKey,
		HostKey:             hostKey,
		PortKey:             portKey,
		DatabaseKey:         databaseKey,
		UsernameKey:         usernameKey,
		PasswordKey:         passwordKey,
	}

	return &options, nil
}

