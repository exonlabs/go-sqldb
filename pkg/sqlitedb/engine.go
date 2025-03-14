// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	sqlite3 "github.com/mattn/go-sqlite3"
)

// Config represents the database configuration params.
type Config struct {
	// database path
	Database string
	// disable Foreign Keys restriction
	NoForeignKeys bool
	// disable Auto Vacuum mode
	NoAutoVacuum bool
	// disable WAL Journal mode
	NoJournalWAL bool
}

// InitConfig initializes configuration from configuration dict.
// it checks and returns error if not all options have valid values.
//
// The parsed options are:
//   - database: (string) the database file path - REQUIRED
//   - no_foreign_keys: (bool) disable Foreign Keys restriction
//   - no_auto_vacuum: (bool) disable Auto Vacuum mode
//   - no_journal_wal: (bool) disable WAL Journal mode
func (cfg *Config) InitConfig(d dictx.Dict) error {
	cfg.Database = dictx.GetString(d, "database", "")
	cfg.NoForeignKeys = dictx.Fetch(d, "no_foreign_keys", false)
	cfg.NoAutoVacuum = dictx.Fetch(d, "no_auto_vacuum", false)
	cfg.NoJournalWAL = dictx.Fetch(d, "no_journal_wal", false)

	// validations
	if cfg.Database == "" {
		return sqldb.ErrDBPath
	}

	return nil
}

// Engine represents the backend engine structure.
type Engine struct {
	// engine config
	config *Config
	// sql driver
	sqlDB *sql.DB
}

// NewEngine creates new engine handler for the backend.
func NewEngine(opts dictx.Dict) (*Engine, error) {
	cfg := &Config{}
	if err := cfg.InitConfig(opts); err != nil {
		return nil, err
	}
	return &Engine{
		config: cfg,
	}, nil
}

// Backend returns the engine backend type.
func (e *Engine) Backend() string {
	return "sqlite"
}

// SqlDB returns a backend driver handler.
func (e *Engine) SqlDB() (*sql.DB, error) {
	if e.sqlDB == nil {
		if e.config == nil {
			return nil, sqldb.ErrDBConfig
		}

		args := []string{}
		if e.config.NoForeignKeys {
			args = append(args, "_foreign_keys=0")
		} else {
			args = append(args, "_foreign_keys=1")
		}
		if !e.config.NoAutoVacuum {
			args = append(args, "_auto_vacuum=1")
		}
		// if !e.config.NoJournalWAL {
		// 	args = append(args, "_journal_mode=WAL")
		// }

		dsn := fmt.Sprintf("%s?%s",
			e.config.Database, strings.Join(args, "&"))

		if v, err := sql.Open("sqlite3", dsn); err != nil {
			return nil, fmt.Errorf("%w - %v", sqldb.ErrOpen, err)
		} else {
			e.sqlDB = v
		}
	}

	return e.sqlDB, nil
}

// Release the backend driver handler.
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
