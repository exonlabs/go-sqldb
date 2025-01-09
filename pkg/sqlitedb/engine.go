// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	sqlite3 "github.com/mattn/go-sqlite3"
)

// open creates new backend driver handler.
func open(cfg *Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("empty config")
	}
	return sql.Open("sqlite3", cfg.DSN())
}

// close shutsdown the backend driver handler.
func close(sdb *sql.DB) error {
	if sdb != nil {
		return sdb.Close()
	}
	return nil
}

// Engine represents the backend engine structure.
type Engine struct {
	cfg *Config
	sdb *sql.DB
}

// NewEngine creates new engine handler for the backend.
func NewEngine(opts dictx.Dict) (*Engine, error) {
	cfg, err := NewConfig(opts)
	if err != nil {
		return nil, err
	}
	sdb, err := open(cfg)
	if err != nil {
		return nil, err
	}
	return &Engine{
		cfg: cfg,
		sdb: sdb,
	}, nil
}

// Backend returns the engine backend type.
func (e *Engine) Backend() string {
	return "sqlite"
}

// SqlDB create or return existing backend driver handler.
func (e *Engine) SqlDB() (*sql.DB, error) {
	if e.sdb == nil {
		return nil, errors.New("no engine driver handler")
	}
	return e.sdb, nil
}

// Release frees the backend driver resources between sessions.
func (e *Engine) Release(_ *sql.DB) error {
	// nothing to do
	return nil
}

// CanRetryErr checks weather an operation error type can be retried.
func (e *Engine) CanRetryErr(err error) bool {
	switch err {
	case sqlite3.ErrBusy, sqlite3.ErrLocked, sqlite3.ErrProtocol,
		sqlite3.ErrIoErr, sqlite3.ErrCantOpen:
		return true
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
	if dictx.Fetch(meta.Args, "without_rowid", false) {
		s := strings.TrimSuffix(strings.TrimSpace(stmts[0]), ";")
		stmts[0] = s + " WITHOUT ROWID;"
	}
	return stmts
}

// SqlGenerator returns the engine SQL statment generator.
func (e *Engine) SqlGenerator() sqldb.SqlGenerator {
	return &SqlGenerator{}
}
