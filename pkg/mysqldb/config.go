// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mysqldb

import (
	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
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
	// connection character set
	CharSet string
	// connection collate
	Collate string
}

// NewConfig creates configuration object from options dict.
// it checks and returns error if not all options have valid values.
//
// The parsed options are:
//   - database: (string) the database name - REQUIRED
//   - host: (string) the database server IP or FQDN - REQUIRED
//   - port: (int) the database server port number - REQUIRED
//   - username: (string) the database  access username (if any)
//   - password: (string) the database access password (if any)
//   - charset: (string) connection character set
//   - collate: (string) connection collate
func NewConfig(opts dictx.Dict) (*Config, error) {
	cfg := &Config{
		Database: dictx.GetString(opts, "database", ""),
		Host:     dictx.GetString(opts, "host", ""),
		Port:     dictx.GetInt(opts, "port", 0),
		Username: dictx.GetString(opts, "username", ""),
		Password: dictx.GetString(opts, "password", ""),
		CharSet:  dictx.GetString(opts, "charset", ""),
		Collate:  dictx.GetString(opts, "collate", ""),
	}

	// validations
	if cfg.Database == "" {
		return nil, sqldb.ErrDBName
	}
	if cfg.Host == "" {
		return nil, sqldb.ErrDBHost
	}
	if cfg.Port == 0 {
		return nil, sqldb.ErrDBPort
	}

	return cfg, nil
}

// DSN returns the driver-specific data source name.
func (cfg *Config) DSN() string {
	// TODO

	return ""
}
