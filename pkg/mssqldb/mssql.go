// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mssqldb

import (
	"database/sql"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
)

type Engine struct {
	sqlDB *sql.DB

	// database config
	config *sqldb.Config
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
