package databasex

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func postgreDSN(cfg *Config) string {
	param := url.Values{}
	param.Add("connect_timeout", fmt.Sprint(cfg.DialTimeout.Seconds()))
	param.Add("user", url.QueryEscape(cfg.User))
	param.Add("password", url.QueryEscape(cfg.Password))
	param.Add("port", fmt.Sprint(cfg.Port))
	param.Add("timezone", cfg.TimeZone)
	param.Add("sslmode", "disable")

	connStr := fmt.Sprintf(connStringPostgresTemplate,
		cfg.Host,
		cfg.Name,
		param.Encode(),
	)

	return connStr
}

func NewPostgres(cfg *Config) (*sqlx.DB, error) {
	db, err := sqlx.Open(cfg.Driver, postgreDSN(cfg))
	if err != nil {
		return db, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return db, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	return db, nil
}
