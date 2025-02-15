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
	"github.com/mattn/go-sqlite3"
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
	return sqldb.BACKEND_SQLITE
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
	if e.config == nil {
		return sqldb.ErrDBConfig
	}

	e.Close()

	extargs := e.config.ConnectArgs
	if !strings.Contains(extargs, "_foreign_keys=") {
		extargs += "&_foreign_keys=1"
	}
	if !strings.Contains(extargs, "_auto_vacuum=") {
		extargs += "&_auto_vacuum=1"
	}
	// if !strings.Contains(extargs, "_journal_mode=") {
	// 	extargs += "&_journal_mode=WAL"
	// }

	dsn := fmt.Sprintf("%s?%s", e.config.Database, extargs)
	if db, err := sql.Open("sqlite3", dsn); err != nil {
		return err
	} else {
		e.sqlDB = db
	}
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
	switch err {
	case sqlite3.ErrBusy, sqlite3.ErrLocked:
		return true
	}
	return false
}

// GenSchema generates table schema.
func (*Engine) GenSchema(tablename string, meta *sqldb.TableMeta) string {
	var buff, constraints, indexes []string

	// if AutoGuid, add guid column if not exist as first column
	if meta.AutoGuid && meta.Columns[0].Name != "guid" {
		meta.Columns = append([]sqldb.ColumnMeta{
			{Name: "guid", Type: "VARCHAR(32) NOT NULL", Primary: true},
		}, meta.Columns...)
	}

	// loop and parse columns meta
	for _, c := range meta.Columns {
		buff = append(buff, c.Name+" "+c.Type)

		// add check constraint for bool datatype
		if strings.Contains(c.Type, "BOOLEAN") {
			constraints = append(constraints,
				fmt.Sprintf("CHECK (%v IN (0,1))", c.Name))
		}

		// add constraints and indexes
		if c.Primary {
			constraints = append(constraints,
				fmt.Sprintf("PRIMARY KEY (%v)", c.Name))
			indexes = append(indexes, fmt.Sprintf(
				"CREATE UNIQUE INDEX IF NOT EXISTS ix_%v_%v ON %v (%v);",
				tablename, c.Name, tablename, c.Name))
		} else if c.Unique && c.Index {
			indexes = append(indexes, fmt.Sprintf(
				"CREATE UNIQUE INDEX IF NOT EXISTS ix_%v_%v ON %v (%v);",
				tablename, c.Name, tablename, c.Name))
		} else {
			if c.Unique {
				constraints = append(constraints,
					fmt.Sprintf("UNIQUE (%v)", c.Name))
			}
			if c.Index {
				indexes = append(indexes, fmt.Sprintf(
					"CREATE INDEX IF NOT EXISTS ix_%v_%v ON %v (%v);",
					tablename, c.Name, tablename, c.Name))
			}
		}
	}

	// append column constraints
	buff = append(buff, constraints...)

	// add explicit table constraints
	for _, c := range meta.Constraints {
		if c.Name != "" {
			buff = append(buff, fmt.Sprintf(
				"CONSTRAINT %s %s", c.Name, c.Definition))
		} else {
			buff = append(buff, c.Definition)
		}
	}

	schema := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (\n  %s\n)",
		tablename, strings.Join(buff, ",\n  "))
	if dictx.Fetch(meta.Args, "sqlite_without_rowid", false) {
		schema += " WITHOUT ROWID;"
	} else {
		schema += ";"
	}
	if len(indexes) > 0 {
		schema += "\n" + strings.Join(indexes, "\n")
	}

	return schema
}
