// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
)

// Config represents the database configuration params.
//
// https://www.sqlite.org/pragma.html
type Config struct {
	// Database path
	Database string
	// ConnectArgs holds connection params
	ConnectArgs string
}

// NewConfig creates configuration object from options dict.
// it checks and returns error if not all options have valid values.
//
// The parsed options are:
//   - database: (string) the database file path - REQUIRED
//   - connect_args: (string) holds connection params
func NewConfig(opts dictx.Dict) (*Config, error) {
	cfg := &Config{
		Database:    dictx.GetString(opts, "database", ""),
		ConnectArgs: dictx.GetString(opts, "connect_args", ""),
	}

	// validations
	if cfg.Database == "" {
		return nil, sqldb.ErrDBPath
	}

	return cfg, nil
}

// DSN returns the driver-specific data source name.
//
// format: dbpath[?param1=value1&...&paramN=valueN]
func (cfg *Config) DSN() string {
	args := []string{}

	conn_args := strings.TrimSpace(cfg.ConnectArgs)
	if len(conn_args) > 0 {
		args = append(args, conn_args)
	}

	if !strings.Contains(conn_args, "_pragma=foreign_keys(") {
		args = append(args, "_pragma=foreign_keys(1)")
	}
	if !strings.Contains(conn_args, "_pragma=busy_timeout(") {
		args = append(args, "_pragma=busy_timeout(100)")
	}

	return fmt.Sprintf("%s?%s", cfg.Database, strings.Join(args, "&"))
}
