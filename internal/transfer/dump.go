package transfer

import (
	"fmt"

	"dbear/internal/config"
)

func Dump(conn config.Connection, schemas []string) ([]byte, error) {
	switch conn.Type {
	case config.TypePostgreSQL, config.TypeMySQL:
		return dumpWithDocker(conn, schemas)
	case config.TypeSQLite:
		return DumpSQLite(conn)
	default:
		return nil, fmt.Errorf("unsupported database type for dump: %s", conn.Type)
	}
}

func dumpWithDocker(conn config.Connection, schemas []string) ([]byte, error) {
	version, err := DetectVersion(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to detect database version: %w", err)
	}

	image := GetDockerImage(conn.Type, version)
	if image == "" {
		return nil, fmt.Errorf("failed to determine docker image for database")
	}

	return DumpDatabase(conn, image, schemas)
}

func DumpFileExtension(connType string) string {
	if connType == config.TypePostgreSQL {
		return ".dump"
	}
	if connType == config.TypeMySQL || connType == config.TypeSQLite {
		return ".sql"
	}
	return ".sql"
}
