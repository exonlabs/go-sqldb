// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"context"
	"database/sql"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/logging"
)

// Engine defines the database engine interface.
type Engine interface {
	// Backend returns the engine backend type.
	Backend() string

	// SqlDB create or return existing backend driver handler.
	SqlDB() (*sql.DB, error)
	// Release frees the backend driver resources between sessions.
	Release(*sql.DB) error

	// CanRetryErr checks weather an operation error type can be retried.
	CanRetryErr(err error) bool

	// SqlGenerator returns the engine SQL statment generator.
	SqlGenerator() SqlGenerator
}

// Database represents the database object.
type Database struct {
	// Log is the logger instance for database logging.
	Log *logging.Logger

	// engine represents the database backend engine
	engine Engine

	// database context
	ctx       context.Context
	ctxCancel context.CancelFunc

	// OperationTimeout defines the timeout in seconds for database operation.
	// use 0 or negative value to disable operation timeout. (default 5.0 sec)
	OperationTimeout float64
	// RetryInterval defines the interval in seconds between operation retries.
	// trials are done untill operation is done or timeout is reached.
	// retry interval value must be > 0. (default 0.1 sec)
	RetryInterval float64
}

// NewDatabase creates a new database handler.
//
// The parsed options are:
//   - operation_timeout: (float64) the timeout in seconds for database operation.
//     use 0 or negative value to disable operation timeout. (default 5.0 sec)
//   - retry_interval: (float64) the interval in seconds between operation retries.
//     trials are done untill operation is done or timeout is reached.
//     retry interval value must be > 0. (default 0.1 sec)
func NewDatabase(log *logging.Logger, engine Engine, opts dictx.Dict) *Database {
	db := &Database{
		Log:    log,
		engine: engine,
	}
	db.ctx, db.ctxCancel = context.WithCancel(context.Background())

	if v := dictx.GetFloat(opts, "operation_timeout", 5.0); v > 0 {
		db.OperationTimeout = v
	} else {
		db.OperationTimeout = -1
	}
	if v := dictx.GetFloat(opts, "retry_interval", 0.1); v > 0 {
		db.RetryInterval = v
	}

	return db
}

// Backend returns the database backend type.
func (db *Database) Backend() string {
	if db.engine != nil {
		return db.engine.Backend()
	}
	return ""
}

// check attrs before running query
func (db *Database) check_run() error {
	if db.engine == nil {
		return ErrDBEngine
	}
	return nil
}

// Session returns a new session handler.
func (db *Database) Session() *Session {
	return NewSession(db)
}

// Query returns a new query handler.
func (db *Database) Query(model Model) *Query {
	return NewQuery(NewSession(db), model)
}

// Ping checks if database connection is active.
func (db *Database) Ping() bool {
	if db.engine != nil {
		if sdb, err := db.engine.SqlDB(); err == nil {
			defer db.engine.Release(sdb)
			return sdb.Ping() == nil
		}
	}
	return false
}

// Shutdown closes all the database sessions.
func (db *Database) Shutdown() {
	if db.ctxCancel != nil {
		db.ctxCancel()
	}
}
