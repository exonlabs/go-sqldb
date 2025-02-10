// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"database/sql"
	"fmt"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
)

type Data = dictx.Dict
type DataAdaptor = func(any) (any, error)

// const SQL_PLACEHOLDER = "$?"

// const (
// 	// CONNECT_TIMEOUT defines the default timeout for connect in seconds.
// 	CONNECT_TIMEOUT = float64(5)
// 	// RETRY_INTERVAL defines the number of connection retries.
// 	RETRY_INTERVAL = int(3)
// )

const (
	BACKEND_NONE = int(iota)
	BACKEND_SQLITE
	BACKEND_MYSQL
	BACKEND_PGSQL
	BACKEND_MSSQL
)

func BACKEND(backend int) string {
	switch backend {
	case BACKEND_SQLITE:
		return "sqlite"
	case BACKEND_MYSQL:
		return "mysql"
	case BACKEND_PGSQL:
		return "pgsql"
	case BACKEND_MSSQL:
		return "mssql"
	}
	return ""
}

type Config struct {
	// database name or path
	Database string
	// database server host and port
	Host string
	Port int
	// database access username and password
	Username string
	Password string
	// extra backend or specific options
	Options dictx.Dict
}

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
	if cfg.Options != nil {
		s += fmt.Sprintf(", Options: %s", dictx.String(cfg.Options))
	}
	return "{" + s + "}"
}

type Engine interface {
	Backend() int
	SqlDB() *sql.DB

	// FormatSql(string) string
	// CanRetryErr(error) bool
}

// type Backend interface {
// 	CreateSchema(string, Model) ([]string, error)
// 	InteractiveConfig(Options) (Options, error)
// 	InteractiveSetup(Options) error
// }
