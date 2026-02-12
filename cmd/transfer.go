package cmd

import (
	"fmt"

	"dbear/internal/connection"
	"dbear/internal/transfer"
	"dbear/internal/ui"

	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer data between databases",
	Long:  "Transfer data from a source database to a destination database using Docker containers or native tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Loading connections...")
		manager := connection.NewManager(configManager)

		connections, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to load connections: %w", err)
		}

		if len(connections) < 2 {
			return fmt.Errorf("at least 2 connections are required for transfer")
		}

		sourceName, err := ui.SelectConnectionWithTitle(connections, "Select source connection")
		if err != nil {
			return fmt.Errorf("failed to select source connection: %w", err)
		}

		sourceConn, err := manager.Get(sourceName)
		if err != nil {
			return fmt.Errorf("failed to load source connection: %w", err)
		}

		if sourceConn == nil {
			return fmt.Errorf("source connection '%s' not found", sourceName)
		}

		filteredConnections := []connection.Connection{}
		for _, conn := range connections {
			if conn.Name != sourceName {
				filteredConnections = append(filteredConnections, conn)
			}
		}

		destName, err := ui.SelectConnectionWithTitle(filteredConnections, "Select destination connection")
		if err != nil {
			return fmt.Errorf("failed to select destination connection: %w", err)
		}

		destConn, err := manager.Get(destName)
		if err != nil {
			return fmt.Errorf("failed to load destination connection: %w", err)
		}

		if destConn == nil {
			return fmt.Errorf("destination connection '%s' not found", destName)
		}

		if sourceConn.Type != destConn.Type {
			return fmt.Errorf("source and destination databases must be of the same type (source: %s, destination: %s)", sourceConn.Type, destConn.Type)
		}

		sourceVersion, err := transfer.DetectVersion(*sourceConn)
		if err != nil {
			return fmt.Errorf("failed to detect source database version: %w", err)
		}

		destVersion, err := transfer.DetectVersion(*destConn)
		if err != nil {
			return fmt.Errorf("failed to detect destination database version: %w", err)
		}

		confirmed, err := ui.ConfirmTransfer(*sourceConn, *destConn, sourceVersion, destVersion)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		if !confirmed {
			return fmt.Errorf("transfer cancelled by user")
		}

		_, err = ui.RunWithSpinner("Transferring data...", func() (interface{}, error) {
			return nil, transfer.Transfer(*sourceConn, *destConn)
		})
		if err != nil {
			return fmt.Errorf("transfer failed: %w", err)
		}

		fmt.Printf("Transfer completed successfully from '%s' to '%s'\n", sourceName, destName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
}

