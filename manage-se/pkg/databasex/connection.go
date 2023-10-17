package databasex

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// CreateSession create new session maria db
func CreateSession(cfg *Config) (*sqlx.DB, error) {

	switch cfg.Driver {
	case `mysql`:
		return NewMySQLSession(cfg)
	case `postgres`:
		return NewPostgres(cfg)

	}

	return nil, fmt.Errorf(`not support database driver %s`, cfg.Driver)
}
