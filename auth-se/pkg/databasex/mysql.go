package databasex

import (
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/go-sql-driver/mysql"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"auth-se/pkg/util"
)

func mysqlDSN(cfg *Config) string {
	if len(strings.Trim(cfg.Charset, "")) == 0 {
		cfg.Charset = "UTF8"
	}

	param := url.Values{}
	param.Add("timeout", fmt.Sprintf("%v", cfg.DialTimeout))
	param.Add("charset", cfg.Charset)
	param.Add("parseTime", "True")
	param.Add("loc", cfg.TimeZone)

	connStr := fmt.Sprintf(connStringMysqlTemplate,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		param.Encode(),
	)

	return connStr
}

func NewMySQLSession(cfg *Config) (*sqlx.DB, error) {
	conf := mysql.NewConfig()
	conf.User = cfg.User
	conf.Passwd = cfg.Password
	conf.Passwd = cfg.Password
	conf.DBName = cfg.Name
	conf.Net = "tcp"
	conf.ParseTime = true
	conf.Timeout = cfg.DialTimeout // dial timeout
	conf.ReadTimeout = cfg.ReadTimeout
	conf.WriteTimeout = cfg.WriteTimeout

	conf.Addr = net.JoinHostPort(cfg.Host, util.ToString(cfg.Port))
	driver, err := mysql.NewConnector(conf)
	if err != nil {
		return nil, err
	}

	drv := sql.OpenDB(driver)

	drv.SetMaxOpenConns(cfg.MaxOpenConns)
	drv.SetMaxIdleConns(cfg.MaxIdleConns)
	drv.SetConnMaxLifetime(cfg.MaxLifetime)
	db := sqlx.NewDb(drv, "mysql")
	return db, nil

}
