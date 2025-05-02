package config

import (
	"context"
	"sync"

	"github.com/sethvargo/go-envconfig"
)

type AppConfig struct {
	DbConfig        DatabaseConfig  `env:", prefix=DB_"`
	WebServerConfig WebServerConfig `env:" ,prefix=WEB_"`
}

var (
	configInstance AppConfig
	configOnce     sync.Once
	configErr      error
)

func Instance() (*AppConfig, error) {
	configOnce.Do(
		func() {
			ctx := context.Background()
			if err := envconfig.Process(ctx, &configInstance); err != nil {
				configErr = err
			}
		},
	)

	return &configInstance, configErr
}
