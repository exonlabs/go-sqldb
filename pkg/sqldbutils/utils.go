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

func InteractiveConfig(backend string, defaults dictx.Dict) (dictx.Dict, error) {
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

func InteractiveSetup(backend string, cfg dictx.Dict) error {
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

func CreateEngine(backend string, cfg dictx.Dict) (sqldb.Engine, error) {
	// switch backend {
	// case sqldb.BACKEND(sqldb.BACKEND_SQLITE):
	// 	return sqlitedb.Engine(cfg)
	// case sqldb.BACKEND(sqldb.BACKEND_MYSQL):
	// 	return mysqldb.Engine(cfg)
	// case sqldb.BACKEND(sqldb.BACKEND_PGSQL):
	// 	return pgsqldb.Engine(cfg)
	// case sqldb.BACKEND(sqldb.BACKEND_MSSQL):
	// 	return mssqldb.Engine(cfg)
	// }
	return nil, sqldb.ErrDBBackend
}
