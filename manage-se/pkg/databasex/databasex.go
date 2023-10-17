package databasex

import "time"

const (
	connStringMysqlTemplate = "%s:%s@tcp(%s:%d)/%s?%s"
	//connStringPostgresTemplate = "%s:%s@%s:%d/%s?%s"
	connStringPostgresTemplate = "postgres://%s/%s?%s"
)

var (
	dsn = map[string]func(*Config) string{
		"mysql":    mysqlDSN,
		"postgres": postgreDSN,
	}
)

type (
	Config struct {
		Host         string
		Port         int
		User         string
		Password     string
		Name         string
		Charset      string
		MaxOpenConns int
		MaxIdleConns int
		MaxLifetime  time.Duration
		Type         string
		TimeZone     string
		Driver       string
		DialTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
)
