// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mysqldb

import (
	"database/sql"
	"sync"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	mysql "github.com/exonlabs/mysql"
)

// Engine represents the backend engine structure.
type Engine struct {
	// engine config
	config *Config
	// sql driver
	sqlDB *sql.DB
	// state numtex
	muState sync.Mutex
}

// NewEngine creates new engine handler for the backend.
func NewEngine(opts dictx.Dict) (*Engine, error) {
	cfg, err := NewConfig(opts)
	if err != nil {
		return nil, err
	}
	return &Engine{
		config: cfg,
	}, nil
}

// Backend returns the engine backend type.
func (e *Engine) Backend() string {
	return "mysql"
}

// SqlDB create or return existing backend driver handler.
func (e *Engine) SqlDB() (*sql.DB, error) {
	e.muState.Lock()
	defer e.muState.Unlock()

	// return existing driver handler
	if e.sqlDB != nil {
		return e.sqlDB, nil
	}

	if e.config == nil {
		return nil, sqldb.ErrDBConfig
	}
	// TODO
	return e.sqlDB, nil
}

// Close shutsdown the backend driver handler and free resources.
func (e *Engine) Close(sqldb *sql.DB) error {
	e.muState.Lock()
	defer e.muState.Unlock()

	// do nothing if already closed
	if e.sqlDB == nil {
		return nil
	}

	if err := e.sqlDB.Close(); err != nil {
		return err
	}
	e.sqlDB = nil
	return nil
}

// Release frees the backend driver resources between sessions.
func (e *Engine) Release(_ *sql.DB) error {
	// nothing to do
	return nil
}

// CanRetryErr checks weather an operation error type can be retried.
func (e *Engine) CanRetryErr(err error) bool {
	switch err {
	case mysql.ErrBusyBuffer:
		return true
	}
	return false
}

// SqlGenerator returns the engine SQL statment generator.
func (e *Engine) SqlGenerator() sqldb.SqlGenerator {
	return &sqldb.StdSqlGenerator{}
}
