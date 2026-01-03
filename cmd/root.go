package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dbear/internal/config"
	"github.com/spf13/cobra"
)

var configPath string
var configManager config.Manager

var rootCmd = &cobra.Command{
	Use:   "dbear",
	Short: "Database connection manager",
	Long:  "A CLI tool for managing database connections",
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	defaultConfigPath := filepath.Join(homeDir, ".config", "dbear", "config.json")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", defaultConfigPath, "config file path")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configManager = config.NewJSONManager(configPath)
}

func Execute() error {
	return rootCmd.Execute()
}

