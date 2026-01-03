package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"dbear/internal/connection"
	"dbear/internal/ui"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a database",
	Long:  "Connect to a database using the specified connection name, or select from a list if no name is provided",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var connectionName string
		var err error

		manager := connection.NewManager(configManager)

		if len(args) == 0 {
			connections, err := manager.List()
			if err != nil {
				return fmt.Errorf("failed to load connections: %w", err)
			}

			connectionName, err = ui.SelectConnection(connections)
			if err != nil {
				return err
			}
		} else {
			connectionName = args[0]
		}

		conn, err := manager.Get(connectionName)
		if err != nil {
			return fmt.Errorf("failed to load connection: %w", err)
		}

		if conn == nil {
			return fmt.Errorf("connection '%s' not found", connectionName)
		}

		connString, err := connection.BuildConnectionString(*conn)
		if err != nil {
			return fmt.Errorf("failed to build connection string: %w", err)
		}

		usqlCmd := exec.Command("usql", connString)
		usqlCmd.Stdin = os.Stdin
		usqlCmd.Stdout = os.Stdout
		usqlCmd.Stderr = os.Stderr

		if err := usqlCmd.Run(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
