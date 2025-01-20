// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"database/sql"
	"fmt"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
)

type Data = map[string]any
type DataAdaptor = func(any) (any, error)

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

const SQL_PLACEHOLDER = "$?"

type DBConfig struct {
	Database string
	Host     string
	Port     int
	Username string
	Password string
	Options  dictx.Dict
}

func (cfg *DBConfig) String() string {
	pw := ""
	if cfg.Password != "" {
		pw = "*****"
	}
	return fmt.Sprintf(
		"database: %s, host: %s, port: %d, username: %s, password: %s, options: %s",
		cfg.Database, cfg.Host, cfg.Port, cfg.Username, pw, cfg.Options,
	)
}

type Engine interface {
	Backend() int
	Config() *DBConfig
	SqlDB() *sql.DB

	// FormatSql(string) string
	// CanRetryErr(error) bool
}

// type Backend interface {
// 	CreateSchema(string, Model) ([]string, error)
// 	InteractiveConfig(Options) (Options, error)
// 	InteractiveSetup(Options) error
// }
