// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
)

// SQL statment variable placeholder
const SQL_PLACEHOLDER = "?"

// Data type defines the table column data. where each column data is
// represented into a map for columns as keys and data as values.
type Data = map[string]any

// Backend represents the databases type.
type Backend uint8

const (
	// Defined backends.
	BACKEND_NONE   = Backend(0)
	BACKEND_SQLITE = Backend(1)
	BACKEND_MYSQL  = Backend(2)
	BACKEND_PGSQL  = Backend(3)
	BACKEND_MSSQL  = Backend(4)
)

// represent the [Backend] type string format
func (b Backend) String() string {
	switch b {
	case BACKEND_SQLITE:
		return "sqlite"
	case BACKEND_MYSQL:
		return "mysql"
	case BACKEND_PGSQL:
		return "pgsql"
	case BACKEND_MSSQL:
		return "mssql"
	}
	return "vv"
}

// Config represents the database configuration params.
type Config struct {
	// database name or path
	Database string
	// database host for client/server type databases.
	Host string
	// database port number for client/server type databases.
	Port int
	// database access username
	Username string
	// database access password
	Password string
	// connection options for backends.
	ConnectArgs string
}

// NewConfig creates a new database [Config] object.
//
// The parsed options are:
//   - database: (string) the database name or path.
//   - host: (string) host for client/server type databases.
//   - port: (int) port number for client/server type databases.
//   - username: (string) database access username.
//   - password: (string) database access password.
//   - connect_args: (string) connection options for backends.
func NewConfig(d dictx.Dict) *Config {
	return &Config{
		Database:    strings.TrimSpace(dictx.GetString(d, "database", "")),
		Host:        strings.TrimSpace(dictx.GetString(d, "host", "")),
		Port:        dictx.GetInt(d, "port", 0),
		Username:    strings.TrimSpace(dictx.GetString(d, "username", "")),
		Password:    strings.TrimSpace(dictx.GetString(d, "password", "")),
		ConnectArgs: strings.TrimSpace(dictx.GetString(d, "connect_args", "")),
	}
}

// represent the [Config] in string format
func (cfg *Config) String() string {
	s := fmt.Sprintf("Database: %s", cfg.Database)
	if cfg.Host != "" {
		s += fmt.Sprintf(", Host: %s, Port: %d", cfg.Host, cfg.Port)
	}
	if cfg.Username != "" || cfg.Password != "" {
		pw := ""
		if cfg.Password != "" {
			pw = "*****"
		}
		s += fmt.Sprintf(", Username: %s, Password: %s", cfg.Username, pw)
	}
	if cfg.ConnectArgs != "" {
		s += fmt.Sprintf(", Args: %s", cfg.ConnectArgs)
	}
	return fmt.Sprintf("{%s}", s)
}

type Engine interface {
	// Backend returns the engine backend type.
	Backend() Backend
	// Config returns the engine connection config.
	Config() *Config
	// SqlDB returns the driver database handler.
	SqlDB() *sql.DB
	// Open the engine backend connection.
	Open() error
	// Close the engine backend connection.
	Close()
	// Checks weather an operation error type can be retried.
	CanRetryErr(err error) bool
	// GenSchema generates table schema.
	GenSchema(tablename string, meta *TableMeta) string
}
