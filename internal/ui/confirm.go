package ui

import (
	"dbear/internal/config"
	"fmt"

	"github.com/charmbracelet/huh"
)

func ConfirmTransfer(source, dest config.Connection, sourceVersion, destVersion string) (bool, error) {
	var confirmed bool

	summary := fmt.Sprintf(
		"Source: %s (%s %s)\nDestination: %s (%s %s)\n\nThis will overwrite all data in the destination database.",
		source.Name,
		source.Type,
		sourceVersion,
		dest.Name,
		dest.Type,
		destVersion,
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm Transfer").
				Description(summary).
				Value(&confirmed),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return false, err
	}

	return confirmed, nil
}

