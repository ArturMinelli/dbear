package cmd

import (
	"fmt"

	"dbear/internal/config"
	"dbear/internal/connection"
	"dbear/internal/ui"

	"github.com/spf13/cobra"
)

var connectionsCmd = &cobra.Command{
	Use:   "connections",
	Short: "Manage database connections",
	Long:  "Create, list, and manage database connections",
}

var connectionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new database connection",
	Long:  "Create a new database connection through an interactive form",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := ui.CreateConnectionForm()
		if err != nil {
			return fmt.Errorf("failed to create connection: %w", err)
		}

		if conn.Name == "" {
			return fmt.Errorf("connection name is required")
		}

		if !config.IsValidType(conn.Type) {
			return fmt.Errorf("invalid database type: %s", conn.Type)
		}

		manager := connection.NewManager(configManager)
		if err := manager.Create(*conn); err != nil {
			return fmt.Errorf("failed to save connection: %w", err)
		}

		fmt.Printf("Connection '%s' created successfully.\n", conn.Name)
		return nil
	},
}

var connectionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all database connections",
	Long:  "Display all saved database connections in a table",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := connection.NewManager(configManager)
		connections, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to load connections: %w", err)
		}

		if err := ui.DisplayConnectionsList(connections); err != nil {
			return fmt.Errorf("failed to display connections: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectionsCmd)
	connectionsCmd.AddCommand(connectionsCreateCmd)
	connectionsCmd.AddCommand(connectionsListCmd)
	connectionsCmd.AddCommand(importCmd)
}
