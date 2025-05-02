package config

type WebServerConfig struct {
	Port int `env:"PORT, default=8080"`
	// More options can follow for middleware
}
