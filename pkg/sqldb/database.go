// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/events"
	"github.com/exonlabs/go-utils/pkg/logging"
)

// Database represents the database object.
type Database struct {
	// engine represents the database backend
	engine Engine

	// DBLog is the logger instance for database logging.
	DBLog *logging.Logger

	// breakEvent signals a break operation.
	breakEvent *events.Event
	// termEvent signals a termination event.
	termEvent *events.Event

	// OperationTimeout defines the timeout in seconds for database operation.
	// use 0 or negative value to disable operation timeout.
	OperationTimeout float64
	// RetryInterval defines the time interval in seconds between operation
	// retries. trials are done untill operation is done or timeout is reached.
	// retry interval value must be > 0.
	RetryInterval float64
}

// NewDatabase creates a new database handler.
//
// The parsed options are:
//   - operation_timeout: (float64) the timeout in seconds for database operation.
//     use 0 or negative value to disable operation timeout.
//   - retry_interval: (float64) the time interval in seconds between operation
//     retries. trials are done untill operation is done or timeout is reached.
//     retry interval value must be > 0.
func NewDatabase(engine Engine, dblog *logging.Logger, opts dictx.Dict) (*Database, error) {
	if engine == nil {
		return nil, ErrDBEngine
	}

	db := &Database{
		engine:           engine,
		DBLog:            dblog,
		breakEvent:       events.New(),
		termEvent:        events.New(),
		OperationTimeout: 5.0,
		RetryInterval:    0.2,
	}

	if v := dictx.GetFloat(opts, "operation_timeout", 0); v > 0 {
		db.OperationTimeout = v
	} else {
		db.OperationTimeout = -1
	}
	if v := dictx.GetFloat(opts, "retry_interval", 0); v > 0 {
		db.RetryInterval = v
	}

	return db, nil
}

// Backend returns the database backend type.
func (db *Database) Backend() Backend {
	if db.engine != nil {
		return db.engine.Backend()
	}
	return BACKEND_NONE
}

// NewSession creates a new session object.
func (db *Database) NewSession() (*Session, error) {
	if db.engine == nil {
		return nil, ErrDBEngine
	}
	return newSession(db)
}

// Checks if database connection is active.
func (db *Database) IsActive() bool {
	if db.engine != nil {
		return db.engine.SqlDB().Ping() == nil
	}
	return false
}

// Breaks any active database operation.
func (db *Database) Break() {
	db.breakEvent.Set()
}

// Closes all the database sessions and operations.
func (db *Database) Close() {
	db.breakEvent.Set()
	db.termEvent.Set()
}

// InitializeDatabase first creates and alter the models table schema,
// then add the intial tables data.
func InitializeDatabase(db *Database, metainfo []TableModelMeta) error {
	if db == nil {
		return ErrDBHandler
	}
	if db.engine == nil {
		return ErrDBEngine
	}

	// create new session
	dbs, err := db.NewSession()
	if err != nil {
		return err
	}

	// create and alter schema
	if db.DBLog != nil {
		db.DBLog.Debug("creating tables schema")
	}
	for _, v := range metainfo {
		if err = v.ModelMeta.CreateSchema(dbs, v.TableName); err != nil {
			return err
		}
	}
	for _, v := range metainfo {
		if err = v.ModelMeta.AlterSchema(dbs, v.TableName); err != nil {
			return err
		}
	}

	// add intial data to tables
	if db.DBLog != nil {
		db.DBLog.Debug("adding tables initial data")
	}
	for _, v := range metainfo {
		if err = v.ModelMeta.InitialData(dbs, v.TableName); err != nil {
			return err
		}
	}

	return nil
}
