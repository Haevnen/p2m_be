// Package gormdb build connection to db with gorm and ddtrace
package gormdb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"go.uber.org/multierr"
	gormMySQLDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Haevnen/p2m_be/pkg/logger"
)

const (
	defaultMySQLPort = 3306
)

// Config the least enough config for gorm
type Config struct {
	DBUser              string
	DBPass              string
	DBPort              int
	DBHost              string
	DBName              string
	DBCollation         string
	DBInterpolateParams bool
	DBLocation          string

	LogSQL      bool
	NotifyQuery string

	MaxOpenConn       int
	MaxLifetimeSecond int
}

// New build gorm DB
func New(cfg *Config) (db *gorm.DB, close func() error, err error) {
	conn, err := newDBConn(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("new db conn: %w", err)
	}

	conn.SetMaxOpenConns(cfg.MaxOpenConn)
	conn.SetConnMaxLifetime(time.Duration(cfg.MaxLifetimeSecond) * time.Second)

	gormDB, err := gorm.Open(
		gormMySQLDriver.New(gormMySQLDriver.Config{Conn: conn}),
		&gorm.Config{
			Logger:      logger.NewGormLogger(&logger.Config{LogSQL: cfg.LogSQL, NotifyQuery: cfg.NotifyQuery}),
			PrepareStmt: true,
		},
	)
	if err != nil {
		if closeErr := conn.Close(); closeErr != nil {
			err = multierr.Append(err, closeErr)
		}
		return nil, nil, err
	}

	return gormDB, conn.Close, nil
}

func newDBConn(cfg *Config) (conn *sql.DB, err error) {
	dns, err := BuildMySQLConnectionString(cfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := sql.Open("mysql", dns)
	return sqlDB, err
}

// BuildMySQLConnectionString build mysql connection string
func BuildMySQLConnectionString(c *Config) (string, error) {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"

	if c.DBUser == "" {
		return "", errors.New("db user is not set")
	}
	cfg.User = c.DBUser
	cfg.Passwd = c.DBPass

	port := defaultMySQLPort
	if c.DBPort != 0 {
		port = c.DBPort
	}
	if strings.Contains(c.DBHost, ":") {
		cfg.Addr = c.DBHost
	} else {
		cfg.Addr = fmt.Sprintf("%s:%d", c.DBHost, port)
	}

	if c.DBName == "" {
		return "", errors.New("db name is not set")
	}
	cfg.DBName = c.DBName

	cfg.ParseTime = true
	cfg.InterpolateParams = c.DBInterpolateParams
	if c.DBCollation != "" {
		cfg.Collation = c.DBCollation
	}

	if c.DBLocation != "" {
		loc, err := time.LoadLocation(c.DBLocation)
		if err != nil {
			return "", err
		}
		cfg.Loc = loc
	}

	return cfg.FormatDSN(), nil
}
