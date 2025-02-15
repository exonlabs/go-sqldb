// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package pgsqldb

import (
	"database/sql"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	"github.com/exonlabs/go-utils/pkg/abc/dictx"
)

type Engine struct {
	// database config
	config *sqldb.Config

	sqlDB *sql.DB
}

func NewEngine(config dictx.Dict) (*Engine, error) {
	cfg, err := GetConfig(config)
	if err != nil {
		return nil, err
	}

	return &Engine{
		config: cfg,
	}, nil
}

func (dbe *Engine) Backend() sqldb.Backend {
	return sqldb.BACKEND_SQLITE
}

func (e *Engine) Config() *sqldb.Config {
	return e.config
}

func (e *Engine) SqlDB() *sql.DB {
	return e.sqlDB
}

// GenSchema generates table schema.
func (*Engine) GenSchema(tablename string, meta *sqldb.TableMeta) string {
	return ""
}
