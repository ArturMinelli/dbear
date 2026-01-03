package connection

import "dbear/internal/config"

type Connection = config.Connection

const (
	TypeMySQL      = config.TypeMySQL
	TypePostgreSQL = config.TypePostgreSQL
	TypeSQLite     = config.TypeSQLite
)

func IsValidType(dbType string) bool {
	return config.IsValidType(dbType)
}

