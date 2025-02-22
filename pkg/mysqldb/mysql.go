// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mysqldb

import (
	"database/sql"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
)

// Engine represents the backend engine structure.
type Engine struct {
	// database config
	config *sqldb.Config
	// sql driver database handler
	sqlDB *sql.DB
}

// NewEngine creates new engine handler for the backend.
func NewEngine(config dictx.Dict) (*Engine, error) {
	cfg, err := GetConfig(config)
	if err != nil {
		return nil, err
	}

	e := &Engine{
		config: cfg,
	}
	if err := e.Open(); err != nil {
		return nil, err
	}

	return e, nil
}

// Backend returns the engine backend type.
func (e *Engine) Backend() sqldb.Backend {
	return sqldb.BACKEND_MYSQL
}

// Config returns the engine connection config.
func (e *Engine) Config() *sqldb.Config {
	return e.config
}

// SqlDB returns the driver database handler.
func (e *Engine) SqlDB() *sql.DB {
	return e.sqlDB
}

// Open the engine backend connection.
func (e *Engine) Open() error {
	return nil
}

// Close the engine backend connection.
func (e *Engine) Close() {
	if e.sqlDB != nil {
		e.sqlDB.Close()
	}
}

// CanRetryErr checks weather an operation error type can be retried.
func (e *Engine) CanRetryErr(err error) bool {
	return false
}

// GenSchema generates table schema.
func (*Engine) GenSchema(tablename string, meta *sqldb.TableMeta) string {
	return ""
}
