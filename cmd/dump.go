package cmd

import (
	"fmt"
	"os"
	"time"

	"dbear/internal/connection"
	"dbear/internal/transfer"
	"dbear/internal/ui"

	"github.com/spf13/cobra"
)

var dumpOutputPath string

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump a database to a file",
	Long:  "Dump the selected source database to a file. Uses Docker for PostgreSQL/MySQL and native sqlite3 for SQLite.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Loading connections...")
		manager := connection.NewManager(configManager)

		connections, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to load connections: %w", err)
		}

		if len(connections) < 1 {
			return fmt.Errorf("at least 1 connection is required for dump")
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

		outputPath := dumpOutputPath
		if outputPath == "" {
			timestamp := time.Now().Format("20060102_150405")
			extension := transfer.DumpFileExtension(sourceConn.Type)
			outputPath = fmt.Sprintf("dump_%s_%s%s", sourceName, timestamp, extension)
		}

		result, err := ui.RunWithSpinner("Dumping database...", func() (interface{}, error) {
			return transfer.Dump(*sourceConn)
		})
		if err != nil {
			return fmt.Errorf("dump failed: %w", err)
		}

		dumpData := result.([]byte)
		fmt.Println("Writing file...")
		if err := os.WriteFile(outputPath, dumpData, 0600); err != nil {
			return fmt.Errorf("failed to write dump file: %w", err)
		}

		fmt.Printf("Dump written to %s\n", outputPath)
		return nil
	},
}

func init() {
	dumpCmd.Flags().StringVarP(&dumpOutputPath, "output", "o", "", "output file path (default: dump_<connection>_<timestamp>.<ext>)")
	rootCmd.AddCommand(dumpCmd)
}
