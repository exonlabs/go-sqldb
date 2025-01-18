// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"database/sql"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/logging"
)

type Logger = logging.Logger
type Options = dictx.Dict

type Data = map[string]any
type DataAdaptor = func(any) (any, error)

type DbInfo struct {
	Database string
	Host     string
	Port     int
	Username string
	Password string
	Options  dictx.Dict
}

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

type Engine interface {
	SqlDB() *sql.DB
	Backend() int
	// FormatSql(string) string
	// CanRetryErr(error) bool
}

// type Backend interface {
// 	CreateSchema(string, Model) ([]string, error)
// 	InteractiveConfig(Options) (Options, error)
// 	InteractiveSetup(Options) error
// }
