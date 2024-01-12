package db

import (
	"net"
	"strconv"

	"github.com/go-sql-driver/mysql"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Host            string `json:"host" yaml:"host"`
	Port            int    `json:"port" yaml:"port"`
	Database        string `json:"database" yaml:"database"`
	Username        string `json:"username" yaml:"username"`
	Password        string `json:"password" yaml:"password"`
	Verbose         bool   `json:"verbose" yaml:"verbose"`
	MaxOpenConns    int    `json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int    `json:"conn_max_idle_time" yaml:"conn_max_idle_time"`
}

func (c *Config) DSN() string {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.User = c.Username
	mysqlConfig.Passwd = c.Password
	mysqlConfig.Addr = net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	mysqlConfig.DBName = c.Database
	mysqlConfig.ParseTime = true
	mysqlConfig.InterpolateParams = true
	mysqlConfig.Collation = "utf8mb4_unicode_ci"
	mysqlConfig.Net = "tcp"
	mysqlConfig.Params = map[string]string{
		"charset": "utf8mb4",
	}

	return mysqlConfig.FormatDSN()
}

func (c *Config) Dialector() gorm.Dialector {
	return gorm_mysql.Open(c.DSN())
}
