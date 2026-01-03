package cmd

import (
	"fmt"

	"dbear/internal/config"
	"dbear/internal/connection"
	"dbear/internal/importer"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	importConnectionNameFlag      string
	importDatabaseTypeFlag        string
	importConnectionStringKeyFlag string
	importHostKeyFlag             string
	importPortKeyFlag             string
	importDatabaseKeyFlag         string
	importUsernameKeyFlag         string
	importPasswordKeyFlag         string
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import database connections from files",
	Long:  "Import database connections from various file formats (env, json, yaml, etc.)",
}

var importEnvCmd = &cobra.Command{
	Use:   "env [filepath]",
	Short: "Import connection from .env file",
	Long: `Import a database connection from a .env file.

You can import in two ways:
1. From a connection string: Use --connection-string-key to specify which env variable contains the connection string
2. From multiple variables: Use --type and optionally customize the env keys with --host-key, --port-key, etc.

By default, when using multiple variables, it follows Laravel conventions:
- DB_HOST, DB_PORT, DB_DATABASE, DB_USERNAME, DB_PASSWORD`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		// Determine connection name
		connectionName := importConnectionNameFlag
		if connectionName == "" {
			var nameInput string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Connection Name").
						Description("Enter a name for this connection").
						Value(&nameInput).
						Validate(func(s string) error {
							if s == "" {
								return fmt.Errorf("connection name is required")
							}
							return nil
						}),
				),
			).WithTheme(huh.ThemeCharm())

			if err := form.Run(); err != nil {
				return fmt.Errorf("failed to get connection name: %w", err)
			}
			connectionName = nameInput
		}

		// Validate database type if provided
		if importDatabaseTypeFlag != "" && !config.IsValidType(importDatabaseTypeFlag) {
			return fmt.Errorf("invalid database type: %s", importDatabaseTypeFlag)
		}

		// Create env importer
		envImporter := importer.NewEnvImporter()

		// Build env-specific options
		envOptions := importer.EnvImportOptions{
			ImportOptions: importer.ImportOptions{
				ConnectionName: connectionName,
				DatabaseType:   importDatabaseTypeFlag,
			},
			ConnectionStringKey: importConnectionStringKeyFlag,
			HostKey:             importHostKeyFlag,
			PortKey:             importPortKeyFlag,
			DatabaseKey:         importDatabaseKeyFlag,
			UsernameKey:         importUsernameKeyFlag,
			PasswordKey:         importPasswordKeyFlag,
		}

		// Set defaults for variable keys if not provided
		if envOptions.ConnectionStringKey == "" {
			if envOptions.HostKey == "" {
				envOptions.HostKey = "DB_HOST"
			}
			if envOptions.PortKey == "" {
				envOptions.PortKey = "DB_PORT"
			}
			if envOptions.DatabaseKey == "" {
				envOptions.DatabaseKey = "DB_DATABASE"
			}
			if envOptions.UsernameKey == "" {
				envOptions.UsernameKey = "DB_USERNAME"
			}
			if envOptions.PasswordKey == "" {
				envOptions.PasswordKey = "DB_PASSWORD"
			}
		}

		// Import connections
		connections, err := envImporter.ImportWithOptions(filePath, envOptions)
		if err != nil {
			return fmt.Errorf("failed to import connection: %w", err)
		}

		if len(connections) == 0 {
			return fmt.Errorf("no connections found in env file")
		}

		// Save connections
		manager := connection.NewManager(configManager)
		for _, conn := range connections {
			if err := manager.Create(conn); err != nil {
				return fmt.Errorf("failed to save connection '%s': %w", conn.Name, err)
			}
			fmt.Printf("Connection '%s' imported successfully.\n", conn.Name)
		}

		return nil
	},
}

func init() {
	// Parent import command flags
	importCmd.PersistentFlags().StringVar(&importConnectionNameFlag, "name", "", "Connection name (will prompt if not provided)")

	// Env subcommand flags
	importEnvCmd.Flags().StringVar(&importDatabaseTypeFlag, "type", "", "Database type (postgresql, mysql, sqlite)")
	importEnvCmd.Flags().StringVar(&importConnectionStringKeyFlag, "connection-string-key", "", "Env variable key containing the connection string")
	importEnvCmd.Flags().StringVar(&importHostKeyFlag, "host-key", "", "Env variable key for host (default: DB_HOST)")
	importEnvCmd.Flags().StringVar(&importPortKeyFlag, "port-key", "", "Env variable key for port (default: DB_PORT)")
	importEnvCmd.Flags().StringVar(&importDatabaseKeyFlag, "database-key", "", "Env variable key for database (default: DB_DATABASE)")
	importEnvCmd.Flags().StringVar(&importUsernameKeyFlag, "username-key", "", "Env variable key for username (default: DB_USERNAME)")
	importEnvCmd.Flags().StringVar(&importPasswordKeyFlag, "password-key", "", "Env variable key for password (default: DB_PASSWORD)")

	importCmd.AddCommand(importEnvCmd)
}

