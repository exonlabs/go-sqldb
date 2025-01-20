// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldbutils

import (
	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/mssqldb"
	"github.com/exonlabs/go-sqldb/pkg/mysqldb"
	"github.com/exonlabs/go-sqldb/pkg/pgsqldb"
	"github.com/exonlabs/go-sqldb/pkg/sqlitedb"
)

func CreateEngine(backend string, cfg *sqldb.DBConfig) (sqldb.Engine, error) {
	switch backend {
	case sqldb.BACKEND(sqldb.BACKEND_SQLITE):
		return sqlitedb.NewEngine(cfg)
	case sqldb.BACKEND(sqldb.BACKEND_MYSQL):
		return mysqldb.NewEngine(cfg)
	case sqldb.BACKEND(sqldb.BACKEND_PGSQL):
		return pgsqldb.NewEngine(cfg)
	case sqldb.BACKEND(sqldb.BACKEND_MSSQL):
		return mssqldb.NewEngine(cfg)
	}
	return nil, sqldb.ErrDBBackend
}

func InteractiveConfig(backend string, defaults dictx.Dict) (*sqldb.DBConfig, error) {
	switch backend {
	case sqldb.BACKEND(sqldb.BACKEND_SQLITE):
		return sqlitedb.InteractiveConfig(defaults)
	case sqldb.BACKEND(sqldb.BACKEND_MYSQL):
		return mysqldb.InteractiveConfig(defaults)
	case sqldb.BACKEND(sqldb.BACKEND_PGSQL):
		return pgsqldb.InteractiveConfig(defaults)
	case sqldb.BACKEND(sqldb.BACKEND_MSSQL):
		return mssqldb.InteractiveConfig(defaults)
	}
	return nil, sqldb.ErrDBBackend
}

func InteractiveSetup(backend string, cfg *sqldb.DBConfig) error {
	switch backend {
	case sqldb.BACKEND(sqldb.BACKEND_SQLITE):
		return sqlitedb.InteractiveSetup(cfg)
	case sqldb.BACKEND(sqldb.BACKEND_MYSQL):
		return mysqldb.InteractiveSetup(cfg)
	case sqldb.BACKEND(sqldb.BACKEND_PGSQL):
		return pgsqldb.InteractiveSetup(cfg)
	case sqldb.BACKEND(sqldb.BACKEND_MSSQL):
		return mssqldb.InteractiveSetup(cfg)
	}
	return sqldb.ErrDBBackend
}
