package transfer

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"dbear/internal/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func DetectVersion(conn config.Connection) (string, error) {
	switch conn.Type {
	case config.TypePostgreSQL:
		return detectPostgreSQLVersion(conn)
	case config.TypeMySQL:
		return detectMySQLVersion(conn)
	case config.TypeSQLite:
		return detectSQLiteVersion(conn)
	default:
		return "unknown", fmt.Errorf("unsupported database type: %s", conn.Type)
	}
}

func detectPostgreSQLVersion(conn config.Connection) (string, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conn.Host,
		conn.Port,
		conn.Username,
		conn.Password,
		conn.Database,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		return "", fmt.Errorf("failed to query version: %w", err)
	}

	majorVersion := extractPostgreSQLMajorVersion(version)
	if majorVersion == "" {
		return "latest", nil
	}

	return majorVersion, nil
}

func detectMySQLVersion(conn config.Connection) (string, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		conn.Username,
		conn.Password,
		conn.Host,
		conn.Port,
		conn.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		return "", fmt.Errorf("failed to query version: %w", err)
	}

	majorVersion := extractMySQLMajorVersion(version)
	if majorVersion == "" {
		return "latest", nil
	}

	return majorVersion, nil
}

func detectSQLiteVersion(conn config.Connection) (string, error) {
	return "native", nil
}

func extractPostgreSQLMajorVersion(version string) string {
	re := regexp.MustCompile(`PostgreSQL (\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) < 2 {
		return ""
	}

	major := matches[1]

	if major == "15" || major == "16" {
		return major
	}

	if major == "14" {
		return "14"
	}

	if major == "13" {
		return "13"
	}

	if major == "12" {
		return "12"
	}

	return major
}

func extractMySQLMajorVersion(version string) string {
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return ""
	}

	major := parts[0]
	minor := parts[1]

	if major == "8" {
		return "8.0"
	}

	if major == "5" {
		if strings.HasPrefix(minor, "7") {
			return "5.7"
		}
		return "5.6"
	}

	return fmt.Sprintf("%s.%s", major, minor)
}

func GetDockerImage(dbType, version string) string {
	if version == "latest" || version == "" {
		if dbType == config.TypePostgreSQL {
			return "postgres:latest"
		}
		if dbType == config.TypeMySQL {
			return "mysql:latest"
		}
		return ""
	}

	if dbType == config.TypePostgreSQL {
		return fmt.Sprintf("postgres:%s", version)
	}

	if dbType == config.TypeMySQL {
		return fmt.Sprintf("mysql:%s", version)
	}

	return ""
}

