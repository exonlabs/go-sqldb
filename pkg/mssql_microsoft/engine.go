// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mssqldb

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/logging"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	_ "github.com/microsoft/go-mssqldb"
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
	return "mssql"
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
		sdb, err := sql.Open("sqlserver", dsn)
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
	return false
}

// SqlGenerator represents mssql SQL statment generator.
type SqlGenerator struct {
	sqldb.StdSqlGenerator
}

// FormatStmt prepares the statment placeholders format
func (*SqlGenerator) FormatStmt(stmt string) string {
	n := strings.Count(stmt, sqldb.SQL_PLACEHOLDER)
	for i := 0; i <= n; i++ {
		stmt = strings.Replace(
			stmt, sqldb.SQL_PLACEHOLDER, "@p"+strconv.Itoa(i+1), 1)
	}
	return stmt
}

// Select generates a SELECT statment from attrs.
func (*SqlGenerator) Select(attrs *sqldb.StmtAttrs) (string, []any) {
	// create the statment
	stmt := "SELECT "
	if attrs.Limit > 0 && len(attrs.Orderby) == 0 {
		stmt += fmt.Sprintf(" TOP(%d)", attrs.Limit)
	}
	if len(attrs.Columns) > 0 {
		stmt += strings.Join(attrs.Columns, ", ")
	} else {
		stmt += "*"
	}
	stmt += " FROM " + attrs.Tablename

	if attrs.Filters != "" {
		stmt += " WHERE " + attrs.Filters
	}
	if len(attrs.Groupby) > 0 {
		stmt += " GROUP BY " + strings.Join(attrs.Groupby, ", ")
	}
	if attrs.Having != "" {
		stmt += " HAVING " + attrs.Having
	}
	if len(attrs.Orderby) > 0 {
		stmt += " ORDER BY " + strings.Join(attrs.Orderby, ", ")
		if attrs.Offset > 0 || attrs.Limit > 0 {
			stmt += fmt.Sprintf(" OFFSET %d ROWS", attrs.Offset)
		}
		if attrs.Limit > 0 {
			stmt += fmt.Sprintf(" FETCH NEXT %d ROWS ONLY", attrs.Limit)
		}
	}
	stmt += ";"

	// create the params for statment placeholders
	params := append(attrs.FiltersArgs, attrs.HavingArgs...)

	return stmt, params
}

// Schema generates table schema statments from metainfo
func (*SqlGenerator) Schema(tablename string, meta *sqldb.TableMeta) []string {
	var buff, constraints, indexes []string

	// if AutoGuid, add guid column if not exist as first column
	if meta.AutoGuid && meta.Columns[0].Name != "guid" {
		meta.Columns = append([]sqldb.ColumnMeta{
			{Name: "guid", Type: "VARCHAR(32) NOT NULL", Primary: true},
		}, meta.Columns...)
	}

	// loop and parse columns meta
	for _, c := range meta.Columns {
		col_type := c.Type
		if strings.Contains(col_type, "BOOLEAN") {
			col_type = strings.ReplaceAll(col_type, "BOOLEAN", "BIT")
			col_type = strings.ReplaceAll(col_type, "false", "0")
			col_type = strings.ReplaceAll(col_type, "true", "1")
		}
		buff = append(buff, fmt.Sprintf("%s %s", c.Name, col_type))

		// add constraints and indexes
		if c.Primary {
			constraints = append(constraints,
				fmt.Sprintf("PRIMARY KEY (%s)", c.Name))
		} else if c.Unique {
			constraints = append(constraints,
				fmt.Sprintf("UNIQUE (%s)", c.Name))
		}
		if c.Primary || c.Index {
			indexes = append(indexes, fmt.Sprintf(
				"IF NOT EXISTS (SELECT * FROM sys.indexes "+
					"WHERE name='ix_%s_%s')\n"+
					"CREATE INDEX ix_%s_%s ON %s (%s);",
				tablename, c.Name, tablename, c.Name, tablename, c.Name))
		}
	}

	// append column constraints
	buff = append(buff, constraints...)

	// add explicit table constraints
	for _, c := range meta.Constraints {
		c_def := c.Definition
		if strings.Contains(c_def, "RESTRICT") {
			c_def = strings.ReplaceAll(c_def, "RESTRICT", "NO ACTION")
		}
		if c.Name != "" {
			buff = append(buff, fmt.Sprintf(
				"CONSTRAINT %s %s", c.Name, c_def))
		} else {
			buff = append(buff, c_def)
		}
	}

	stmt := fmt.Sprintf(
		"IF OBJECT_ID(N'%s', N'U') IS NULL\n"+
			"CREATE TABLE %s (\n  %s\n);",
		tablename, tablename, strings.Join(buff, ",\n  "))

	return append([]string{stmt}, indexes...)
}

// SqlGenerator returns the engine SQL statment generator.
func (e *Engine) SqlGenerator() sqldb.SqlGenerator {
	return &SqlGenerator{}
}
