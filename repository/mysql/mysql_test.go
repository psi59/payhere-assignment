package mysql

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	defaultTestDBDSN = "root:1234@tcp(127.0.0.1:3306)/mysql?multiStatements=true&readTimeout=30s&writeTimeout=30s&charset=utf8mb4%2Cutf8"
)

var conn *gorm.DB

func TestMain(m *testing.M) {
	dbName := fmt.Sprintf("payhere_%s", xid.New())
	if err := createTestDatabase(dbName); err != nil {
		log.Fatal().Err(err).Send()
	}
	defer func() {
		if err := deleteTestDatabase(dbName); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()
	code := m.Run()
	os.Exit(code)
}

func createTestDatabase(dbName string) error {
	dsn := getEnv("TEST_DATABASE_DSN", defaultTestDBDSN)
	mysqlConfig, err := mysql.ParseDSN(dsn)
	if err != nil {
		return errors.Wrap(err, "failed to parse dsn")
	}
	mysqlConfig.DBName = ""
	mysqlConfig.ParseTime = true
	mysqlConfig.InterpolateParams = true

	dialect := gorm_mysql.Open(mysqlConfig.FormatDSN())
	conn, err = gorm.Open(dialect, &gorm.Config{DisableNestedTransaction: true})
	if err != nil {
		return errors.Wrap(err, "failed to connect database")
	}
	conn = conn.Debug()

	//goDB, err := conn.DB()
	if err != nil {
		return errors.Wrap(err, "failed to get golang db")
	}
	if err := conn.Exec("CREATE DATABASE " + dbName + " COLLATE utf8_general_ci").Error; err != nil {
		return errors.Wrap(err, "failed to create database")
	}

	mysqlConfig.DBName = dbName
	dialect = gorm_mysql.Open(mysqlConfig.FormatDSN())
	log.Debug().Msgf("DSN string : %s", mysqlConfig.FormatDSN())

	conn, err = gorm.Open(dialect, &gorm.Config{DisableNestedTransaction: true})
	if err != nil {
		return errors.Wrap(err, "failed to connect database")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get current working directory for createTestDatabase")
	}
	query, err := ioutil.ReadFile(fmt.Sprintf("%s/tables.sql", currentDir))
	if err != nil {
		return errors.Wrap(err, "failed to read table sql file")
	}

	if err := conn.Exec(string(query)).Error; err != nil {
		return errors.Wrap(err, "failed to execute table creates query")
	}
	conn = conn.Debug()

	return nil
}

func deleteTestDatabase(dbName string) error {
	goDB, err := conn.DB()
	if err != nil {
		return errors.Wrap(err, "failed to get golang db")
	}

	if _, err := goDB.Exec("DROP DATABASE ", dbName); err != nil {
		return errors.Wrapf(err, "failed to delete database(%s)", dbName)
	}

	return nil
}

func getEnv(k, defaultValue string) string {
	v := os.Getenv(k)
	if v == "" {
		return defaultValue
	}

	return v
}
