// Package bootstrap
package bootstrap

import (
	"fmt"

	"auth-se/internal/appctx"
	"auth-se/pkg/databasex"
	"auth-se/pkg/logger"
)

// RegistryDatabase initialize database session
func RegistryDatabase(cfg *appctx.Database) *databasex.DB {

	db, err := databasex.CreateSession(&databasex.Config{
		Driver:       cfg.Driver,
		Host:         cfg.Host,
		Name:         cfg.Name,
		Password:     cfg.Pass,
		Port:         cfg.Port,
		User:         cfg.User,
		DialTimeout:  cfg.DialTimeout,
		MaxOpenConns: cfg.MaxOpen,
		MaxIdleConns: cfg.MaxIdle,
		MaxLifetime:  cfg.MaxLifeTime,
		Charset:      cfg.Charset,
		TimeZone:     cfg.Timezone,
	})
	if err != nil {
		logger.Fatal(
			err,
			logger.EventName("db"),
			logger.Any("host", cfg.Host),
			logger.Any("port", cfg.Port),
			logger.Any("driver", cfg.Driver),
			logger.Any("timezone", cfg.Timezone),
		)
	}

	//db := databasex.New()
	return databasex.New(db, false, cfg.Name)
}

// RegistryMultiDatabase initialize database session
func RegistryMultiDatabase(cfgWrite *appctx.Database, cfgRead *appctx.Database) databasex.Adapter {
	lf := logger.NewFields(
		logger.EventName("db"),
		logger.Any("host_read", cfgRead.Host),
		logger.Any("port_read", cfgRead.Port),
		logger.Any("host_write", cfgWrite.Host),
		logger.Any("port_write", cfgWrite.Port),
		logger.Any("driver_write", cfgWrite.Driver),
		logger.Any("timezone_write", cfgWrite.Timezone),
		logger.Any("driver_read", cfgRead.Driver),
		logger.Any("timezone_read", cfgRead.Timezone),
	)
	dbWrite, err := databasex.CreateSession(&databasex.Config{
		Driver:       cfgWrite.Driver,
		Host:         cfgWrite.Host,
		Name:         cfgWrite.Name,
		Password:     cfgWrite.Pass,
		Port:         cfgWrite.Port,
		User:         cfgWrite.User,
		DialTimeout:  cfgWrite.DialTimeout,
		MaxOpenConns: cfgWrite.MaxOpen,
		MaxIdleConns: cfgWrite.MaxIdle,
		MaxLifetime:  cfgWrite.MaxLifeTime,
		Charset:      cfgWrite.Charset,
		TimeZone:     cfgWrite.Timezone,
		WriteTimeout: cfgWrite.WriteTimeout,
		ReadTimeout:  cfgWrite.ReadTimeout,
	})

	if err != nil {
		logger.Fatal(fmt.Sprintf("db write %v", err), lf...)
	}

	dbRead, err := databasex.CreateSession(&databasex.Config{
		Driver:       cfgRead.Driver,
		Host:         cfgRead.Host,
		Name:         cfgRead.Name,
		Password:     cfgRead.Pass,
		Port:         cfgRead.Port,
		User:         cfgRead.User,
		DialTimeout:  cfgRead.DialTimeout,
		MaxOpenConns: cfgRead.MaxOpen,
		MaxIdleConns: cfgRead.MaxIdle,
		MaxLifetime:  cfgRead.MaxLifeTime,
		Charset:      cfgRead.Charset,
		TimeZone:     cfgRead.Timezone,
		WriteTimeout: cfgRead.WriteTimeout,
		ReadTimeout:  cfgRead.ReadTimeout,
	})

	if err != nil {
		logger.Fatal(fmt.Sprintf("db read %v", err), lf...)
	}

	return databasex.NewMulti(databasex.New(dbWrite, false, cfgRead.Name), databasex.New(dbRead, true, cfgRead.Name))
}
