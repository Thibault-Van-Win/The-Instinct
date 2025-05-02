package config

import "fmt"

type DataBaseType string

const (
	MongoDB DataBaseType = "mongodb"
)

type DatabaseConfig struct {
	Type DataBaseType `env:"TYPE,default=mongodb"`
	User string       `env:"USER, default=user"`
	Pass string       `env:"PASS, default=secret"`
	Host string       `env:"HOST, default=localhost"`
	Port int          `env:"PORT, default=27017"`
}

func (dbc *DatabaseConfig) ConnString() (string, error) {
	switch dbc.Type {
	case MongoDB:
		return fmt.Sprintf(
			"mongodb://%s:%s@%s:%d",
			dbc.User,
			dbc.Pass,
			dbc.Host,
			dbc.Port,
		), nil
	default:
		return "", fmt.Errorf("unsupported database type: %s", dbc.Type)
	}
}
