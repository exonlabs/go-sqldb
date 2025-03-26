// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"database/sql"
	"strings"
	"sync"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/logging"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// Engine represents the backend engine structure.
type Engine struct {
	// Log is the logger instance for database logging.
	Log *logging.Logger

	// engine config
	cfg *Config
	// driver handler
	sdb *sql.DB

	// muState defines mutex for state change operations (open/close).
	muState sync.Mutex
}

// NewEngine creates new engine handler for the backend.
func NewEngine(log *logging.Logger, opts dictx.Dict) (*Engine, error) {
	cfg, err := NewConfig(opts)
	if err != nil {
		return nil, err
	}
	return &Engine{
		Log: log,
		cfg: cfg,
	}, nil
}

// Backend returns the engine backend type.
func (e *Engine) Backend() string {
	return "sqlite"
}

// SqlDB create or return existing backend driver handler.
func (e *Engine) SqlDB() (*sql.DB, error) {
	if e.cfg == nil {
		return nil, sqldb.ErrDBConfig
	}

	e.muState.Lock()
	defer e.muState.Unlock()

	// create new backend driver handler
	if e.sdb == nil {
		dsn := e.cfg.DSN()
		if e.Log != nil {
			e.Log.Trace("Open SqlDB: %s", dsn)
		}
		sdb, err := sql.Open("sqlite", dsn)
		if err != nil {
			return nil, err
		}
		e.sdb = sdb
	}

	return e.sdb, nil
}

// Release frees the backend driver resources between sessions.
func (e *Engine) Release(_ *sql.DB) error {
	// nothing to do
	return nil
}

// Close shutsdown the backend driver handler and free resources.
func (e *Engine) Close(_ *sql.DB) error {
	e.muState.Lock()
	defer e.muState.Unlock()

	if e.sdb != nil {
		if e.Log != nil {
			e.Log.Trace("Close SqlDB")
		}
		if err := e.sdb.Close(); err != nil {
			return err
		}
		e.sdb = nil
	}

	return nil
}

// CanRetryErr checks weather an operation error type can be retried.
func (e *Engine) CanRetryErr(err error) bool {
	if err, ok := err.(*sqlite.Error); ok {
		switch err.Code() {
		case sqlite3.SQLITE_BUSY, sqlite3.SQLITE_LOCKED,
			sqlite3.SQLITE_IOERR, sqlite3.SQLITE_CANTOPEN:
			return true
		}
	}
	return false
}

// SqlGenerator represents sqlite SQL statment generator.
type SqlGenerator struct {
	sqldb.StdSqlGenerator
}

// Schema generates table schema statments from metainfo
func (g *SqlGenerator) Schema(tablename string, meta *sqldb.TableMeta) []string {
	stmts := g.StdSqlGenerator.Schema(tablename, meta)

	if dictx.Fetch(meta.Args, "sqlite_without_rowid", false) {
		s := strings.TrimSuffix(strings.TrimSpace(stmts[0]), ";")
		stmts[0] = s + " WITHOUT ROWID;"
	}

	return stmts
}

// SqlGenerator returns the engine SQL statment generator.
func (e *Engine) SqlGenerator() sqldb.SqlGenerator {
	return &SqlGenerator{}
}
