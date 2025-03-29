// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mssqldb

import (
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
)

// Config represents the database configuration params.
//
// https://github.com/microsoft/go-mssqldb?tab=readme-ov-file#connection-parameters-and-dsn
type Config struct {
	// Database name
	Database string
	// Database server address
	Address string
	// Database access username
	Username string
	// Database access password
	Password string
	// ConnectArgs holds connection params
	ConnectArgs string
}

// NewConfig creates configuration object from options dict.
// it checks and returns error if not all options have valid values.
//
// The parsed options are:
//   - database: (string) the database name - REQUIRED
//   - address: (string) the database server address - REQUIRED
//   - username: (string) the database  access username (if any)
//   - password: (string) the database access password (if any)
//   - connect_args: (string) holds connection params
func NewConfig(opts dictx.Dict) (*Config, error) {
	cfg := &Config{
		Database:    dictx.GetString(opts, "database", ""),
		Address:     dictx.GetString(opts, "address", ""),
		Username:    dictx.GetString(opts, "username", ""),
		Password:    dictx.GetString(opts, "password", ""),
		ConnectArgs: dictx.GetString(opts, "connect_args", ""),
	}

	// validations
	if cfg.Database == "" {
		return nil, sqldb.ErrDBName
	}
	if cfg.Address == "" {
		return nil, sqldb.ErrDBAddr
	}

	return cfg, nil
}

// DSN returns the driver-specific data source name.
//
// format: sqlserver://[username[:password]@]address[?param1=value1&...&paramN=valueN]
func (cfg *Config) DSN() string {
	args := []string{
		"database=" + cfg.Database,
	}

	conn_args := strings.TrimSpace(cfg.ConnectArgs)
	if len(conn_args) > 0 {
		args = append(args, conn_args)
	}

	var auth string
	if cfg.Username != "" {
		if cfg.Password == "" {
			auth = fmt.Sprintf("%s@", cfg.Username)
		} else {
			auth = fmt.Sprintf("%s:%s@", cfg.Username, cfg.Password)
		}
	}

	return fmt.Sprintf("sqlserver://%s%s?%s",
		auth, cfg.Address, strings.Join(args, "&"))
}
