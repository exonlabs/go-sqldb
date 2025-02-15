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
	// termEvent signals a termination operation.
	termEvent *events.Event

	// ConnectTimeout defines the timeout in seconds for database connection.
	// use 0 or negative value to disable timeout.
	ConnectTimeout float64
	// ConnectInterval defines the time interval in seconds between connect
	// retries. trials is done untill connection opens or timeout is reached.
	// connect interval value must be > 0.
	ConnectInterval float64
}

// NewDatabase creates a new database handler.
//
// The parsed options are:
//   - connect_timeout: (float64) the timeout in seconds for database connection.
//     use 0 or negative value to disable timeout.
//   - connect_interval: (float64) the time interval in seconds between connect
//     retries. trials is done untill connection opens or timeout is reached.
//     connect interval value must be > 0.
func NewDatabase(engine Engine, dblog *logging.Logger, opts dictx.Dict) (*Database, error) {
	if engine == nil {
		return nil, ErrDBEngine
	}

	db := &Database{
		engine:          engine,
		DBLog:           dblog,
		breakEvent:      events.New(),
		termEvent:       events.New(),
		ConnectTimeout:  5.0,
		ConnectInterval: 0.2,
	}

	if v := dictx.GetFloat(opts, "connect_timeout", 0); v > 0 {
		db.ConnectTimeout = v
	} else {
		db.ConnectTimeout = -1
	}
	if v := dictx.GetFloat(opts, "connect_interval", 0); v > 0 {
		db.ConnectInterval = v
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

// Closes all the database sessions and operation.
func (db *Database) Close() {
	db.termEvent.Set()
}

// InitializeDatabase creates the database schema and adds intial data.
// It first creates and alter the tables schema in a transactional scope, then
// adds the intial tables data in second transaction.
func InitializeDatabase(db *Database, metainfo map[string]ModelMeta) error {
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

	// create and alter schema in transaction
	if db.DBLog != nil {
		db.DBLog.Debug("creating tables schema")
	}
	if err = dbs.Begin(); err == nil {
		for tablename, meta := range metainfo {
			if err = meta.CreateSchema(dbs, tablename); err != nil {
				return err
			}
		}
		for tablename, meta := range metainfo {
			if err = meta.AlterSchema(dbs, tablename); err != nil {
				return err
			}
		}
		err = dbs.Commit()
	}
	if err != nil {
		return err
	}

	// add intial data to tables in transaction
	if db.DBLog != nil {
		db.DBLog.Debug("adding tables initial data")
	}
	if err = dbs.Begin(); err == nil {
		for tablename, meta := range metainfo {
			if err = meta.InitialData(dbs, tablename); err != nil {
				return err
			}
		}
		err = dbs.Commit()
	}

	return err
}
