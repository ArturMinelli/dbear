package transfer

import (
	"fmt"

	"dbear/internal/config"
)

func Transfer(source, dest config.Connection) error {
	if source.Type != dest.Type {
		return fmt.Errorf("source and destination databases must be of the same type")
	}

	switch source.Type {
	case config.TypePostgreSQL, config.TypeMySQL:
		return transferWithDocker(source, dest)
	case config.TypeSQLite:
		return transferSQLite(source, dest)
	default:
		return fmt.Errorf("unsupported database type: %s", source.Type)
	}
}

func transferWithDocker(source, dest config.Connection) error {
	sourceVersion, err := DetectVersion(source)
	if err != nil {
		return fmt.Errorf("failed to detect source version: %w", err)
	}

	destVersion, err := DetectVersion(dest)
	if err != nil {
		return fmt.Errorf("failed to detect destination version: %w", err)
	}

	sourceImage := GetDockerImage(source.Type, sourceVersion)
	if sourceImage == "" {
		return fmt.Errorf("failed to determine docker image for source database")
	}

	destImage := GetDockerImage(dest.Type, destVersion)
	if destImage == "" {
		return fmt.Errorf("failed to determine docker image for destination database")
	}

	dumpData, err := DumpDatabase(source, sourceImage)
	if err != nil {
		return fmt.Errorf("failed to dump source database: %w", err)
	}

	if err := RestoreDatabase(dest, destImage, dumpData); err != nil {
		return fmt.Errorf("failed to restore destination database: %w", err)
	}

	return nil
}

func transferSQLite(source, dest config.Connection) error {
	dumpData, err := DumpSQLite(source)
	if err != nil {
		return fmt.Errorf("failed to dump source database: %w", err)
	}

	if err := RestoreSQLite(dest, dumpData); err != nil {
		return fmt.Errorf("failed to restore destination database: %w", err)
	}

	return nil
}

