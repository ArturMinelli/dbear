package config

const (
	TypeMySQL      = "mysql"
	TypePostgreSQL = "postgresql"
	TypeSQLite     = "sqlite"
)

type Connection struct {
	Name     string `json:"name" yaml:"name"`
	Type     string `json:"type" yaml:"type"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Database string `json:"database" yaml:"database"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Config struct {
	Connections []Connection `json:"connections" yaml:"connections"`
}

type Manager interface {
	Load() (*Config, error)
	Save(config *Config) error
}

func IsValidType(dbType string) bool {
	return dbType == TypeMySQL || dbType == TypePostgreSQL || dbType == TypeSQLite
}

