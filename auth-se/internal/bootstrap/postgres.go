package bootstrap

import (
	"time"

	"auth-se/pkg/logger"
	"auth-se/pkg/postgres"

	config "auth-se/internal/appctx"
)

func RegistryPostgres(cfg *config.Database) postgres.Adapter {
	db, err := postgres.NewAdapter(&postgres.Config{
		Host:         cfg.Host,
		Name:         cfg.Name,
		Password:     cfg.Pass,
		Port:         cfg.Port,
		User:         cfg.User,
		Timeout:      time.Duration(cfg.TimeoutSecond) * time.Second,
		MaxOpenConns: cfg.MaxOpen,
		MaxIdleConns: cfg.MaxIdle,
		MaxLifetime:  time.Duration(cfg.MaxLifeTimeMS) * time.Millisecond,
		Timezone:     cfg.Timezone,
	})

	if err != nil {
		logger.Fatal(
			err,
			logger.EventName("db"),
			logger.Any("host", cfg.Host),
			logger.Any("port", cfg.Port),
		)
	}

	return db
}
