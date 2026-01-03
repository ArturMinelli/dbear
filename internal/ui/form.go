package ui

import (
	"dbear/internal/config"
	"fmt"
	"strconv"
	"github.com/charmbracelet/huh"
)

func CreateConnectionForm() (*config.Connection, error) {
	var name, dbType, host, portStr, database, username, password string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Connection Name").
				Description("A unique name for this connection").
				Value(&name).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("connection name is required")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("Database Type").
				Description("Select the database type").
				Options(
					huh.NewOption("PostgreSQL", config.TypePostgreSQL),
					huh.NewOption("MySQL", config.TypeMySQL),
					huh.NewOption("SQLite", config.TypeSQLite),
				).
				Value(&dbType),
			huh.NewInput().
				Title("Host").
				Description("Database host address").
				Value(&host).
				Placeholder("localhost"),
			huh.NewInput().
				Title("Port").
				Description("Database port number").
				Value(&portStr).
				Placeholder("5432"),
			huh.NewInput().
				Title("Database").
				Description("Database name").
				Value(&database),
			huh.NewInput().
				Title("Username").
				Description("Database username").
				Value(&username),
			huh.NewInput().
				Title("Password").
				Description("Database password").
				Value(&password).
				EchoMode(huh.EchoModePassword),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return nil, err
	}

	portValue := 5432
	if portStr != "" {
		if parsed, err := strconv.Atoi(portStr); err == nil {
			portValue = parsed
		}
	}

	if dbType == config.TypeSQLite {
		portValue = 0
		if host == "" {
			host = "localhost"
		}
	}

	conn := &config.Connection{
		Name:     name,
		Type:     dbType,
		Host:     host,
		Port:     portValue,
		Database: database,
		Username: username,
		Password: password,
	}

	return conn, nil
}

