// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mysqldb

import (
	"database/sql"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
)

type Engine struct {
	sqlDB *sql.DB

	// database config
	config *sqldb.Config
}

func initConfig(d dictx.Dict) (*sqldb.Config, error) {
	for _, k := range []string{"database", "host", "username", "password"} {
		if dictx.IsExist(d, k) {
			s := strings.TrimSpace(dictx.GetString(d, k, ""))
			dictx.Set(d, k, s)
		}
	}

	if dictx.GetString(d, "database", "") == "" {
		return sqldb.ErrDBName
	}
	if dictx.GetString(d, "host", "") == "" {
		return sqldb.ErrDBHost
	}
	if dictx.GetInt(d, "port", 0) == 0 {
		return sqldb.ErrDBPort
	}
	return nil
}

func NewEngine(cfg *sqldb.Config) (*Engine, error) {
	if err := PrepareConfig(cfg); err != nil {
		return nil, err
	}

	return &Engine{
		config: cfg,
	}, nil
}

func (dbe *Engine) Backend() int {
	return sqldb.BACKEND_SQLITE
}

func (dbe *Engine) Config() *sqldb.Config {
	return dbe.config
}

func (dbe *Engine) SqlDB() *sql.DB {
	return dbe.sqlDB
}
