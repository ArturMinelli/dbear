package connection

import (
	"fmt"
	"net/url"

	"dbear/internal/config"
)

func BuildConnectionString(conn config.Connection) (string, error) {
	switch conn.Type {
	case TypePostgreSQL:
		return buildPostgreSQLString(conn), nil
	case TypeMySQL:
		return buildMySQLString(conn), nil
	case TypeSQLite:
		return buildSQLiteString(conn), nil
	default:
		return "", fmt.Errorf("unsupported database type: %s", conn.Type)
	}
}

func buildPostgreSQLString(conn config.Connection) string {
	user := url.UserPassword(conn.Username, conn.Password)
	u := &url.URL{
		Scheme: "postgres",
		User:   user,
		Host:   fmt.Sprintf("%s:%d", conn.Host, conn.Port),
		Path:   conn.Database,
	}
	return u.String()
}

func buildMySQLString(conn config.Connection) string {
	user := url.UserPassword(conn.Username, conn.Password)
	u := &url.URL{
		Scheme: "mysql",
		User:   user,
		Host:   fmt.Sprintf("%s:%d", conn.Host, conn.Port),
		Path:   conn.Database,
	}
	return u.String()
}

func buildSQLiteString(conn config.Connection) string {
	return fmt.Sprintf("sqlite://%s", conn.Database)
}

