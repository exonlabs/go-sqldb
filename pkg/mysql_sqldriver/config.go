// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mysqldb

import (
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
)

// Config represents the database configuration params.
//
// https://github.com/go-sql-driver/mysql?tab=readme-ov-file#dsn-data-source-name
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
// format: [username[:password]@]address/dbname[?param1=value1&...&paramN=valueN]
func (cfg *Config) DSN() string {
	args := []string{}

	conn_args := strings.TrimSpace(cfg.ConnectArgs)
	if len(conn_args) > 0 {
		args = append(args, conn_args)
	}

	if !strings.Contains(conn_args, "charset=") {
		args = append(args, "charset=utf8mb4,utf8")
	}
	if !strings.Contains(conn_args, "collation=") {
		args = append(args, "collation=utf8mb4_unicode_ci")
	}
	if !strings.Contains(conn_args, "multiStatements=") {
		args = append(args, "multiStatements=true")
	}
	if !strings.Contains(conn_args, "interpolateParams=") {
		args = append(args, "interpolateParams=true")
	}

	var auth string
	if cfg.Username != "" {
		if cfg.Password == "" {
			auth = fmt.Sprintf("%s@", cfg.Username)
		} else {
			auth = fmt.Sprintf("%s:%s@", cfg.Username, cfg.Password)
		}
	}

	return fmt.Sprintf("%s%s/%s?%s",
		auth, cfg.Address, cfg.Database, strings.Join(args, "&"))
}
