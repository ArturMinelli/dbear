package cmd

import (
	"fmt"

	"dbear/internal/connection"
	"dbear/internal/importer"
	"dbear/internal/ui"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import database connections from files",
	Long:  "Import database connections from various file formats (env, json, yaml, etc.)",
}

var importEnvCmd = &cobra.Command{
	Use:   "env [filepath]",
	Short: "Import connection from .env file",
	Long: `Import a database connection from a .env file through an interactive form.

The form will guide you through:
1. Naming your connection
2. Choosing your import method (connection string or multiple variables)
3. Configuring the environment variable keys

By default, when using multiple variables, it follows Laravel conventions:
- DB_HOST, DB_PORT, DB_DATABASE, DB_USERNAME, DB_PASSWORD`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		// Get import configuration from interactive form
		envOptions, err := ui.ImportEnvForm()
		if err != nil {
			return fmt.Errorf("failed to get import configuration: %w", err)
		}

		// Create env importer
		envImporter := importer.NewEnvImporter()

		// Import connections
		connections, err := envImporter.ImportWithOptions(filePath, *envOptions)
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
	importCmd.AddCommand(importEnvCmd)
}

