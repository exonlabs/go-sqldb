// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mssqldb

import (
	"database/sql"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	_ "github.com/microsoft/go-mssqldb"
)

// Config represents the database configuration params.
type Config struct {
	// database name
	Database string
	// database server host
	Host string
	// database server port number
	Port int
	// database access username
	Username string
	// database access password
	Password string
}

// InitConfig initializes configuration from configuration dict.
// it checks and returns error if not all options have valid values.
//
// The parsed options are:
//   - database: (string) the database name - REQUIRED
//   - host: (string) the database server IP or FQDN - REQUIRED
//   - port: (int) the database server port number - REQUIRED
//   - username: (string) the database  access username (if any)
//   - password: (string) the database access password (if any)
func (cfg *Config) InitConfig(d dictx.Dict) error {
	cfg.Database = dictx.GetString(d, "database", "")
	cfg.Host = dictx.GetString(d, "host", "")
	cfg.Port = dictx.GetInt(d, "port", 0)
	cfg.Username = dictx.GetString(d, "username", "")
	cfg.Password = dictx.GetString(d, "password", "")

	// validations
	if cfg.Database == "" {
		return sqldb.ErrDBName
	}
	if cfg.Host == "" {
		return sqldb.ErrDBHost
	}
	if cfg.Port == 0 {
		return sqldb.ErrDBPort
	}

	return nil
}

// Engine represents the backend engine structure.
type Engine struct {
	// engine config
	config *Config
	// sql driver
	sqlDB *sql.DB
}

// NewEngine creates new engine handler for the backend.
func NewEngine(opts dictx.Dict) (*Engine, error) {
	cfg := &Config{}
	if err := cfg.InitConfig(opts); err != nil {
		return nil, err
	}
	return &Engine{
		config: cfg,
	}, nil
}

// Backend returns the engine backend type.
func (e *Engine) Backend() string {
	return "mssql"
}

// SqlDB returns a backend driver handler.
func (e *Engine) SqlDB() (*sql.DB, error) {
	if e.sqlDB == nil {
		if e.config == nil {
			return nil, sqldb.ErrDBConfig
		}

		// TODO: create sqlDB
	}

	return e.sqlDB, nil
}

// Release the backend driver handler.
func (e *Engine) Release(_ *sql.DB) error {
	// nothing to do
	return nil
}

// CanRetryErr checks weather an operation error type can be retried.
func (e *Engine) CanRetryErr(err error) bool {
	return false
}

// SqlGenerator returns the engine SQL statment generator.
func (e *Engine) SqlGenerator() sqldb.SqlGenerator {
	return &sqldb.StdSqlGenerator{}
}
