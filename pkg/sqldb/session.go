// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/exonlabs/go-utils/pkg/events"
)

// Session interface
type Session interface {
	// Database returns the associated database handler.
	Database() *Database
	// Query returns a new query handler.
	Query(model Model) Query

	// // Open the backend connection.
	// Open() error
	// // Close the backend connection.
	// Close()
	// // Ping checks if the backend connection is active.
	// Ping() bool
	// Cancel breaks all active session operation.
	Cancel()

	// Begin starts a new transactional scope.
	Begin() error
	// Commit runs a commit action in transactional scope.
	Commit() error
	// RollBack runs a rollback action in transactional scope.
	RollBack() error
	// Exec runs a query without returning any rows. it takes the statment
	// to run and the args are for any placeholder parameters in the query.
	Exec(stmt string, params ...any) (int, error)
	// Fetch runs a query that returns rows. it takes the statment
	// to run and the args are for any placeholder parameters in the query.
	Fetch(stmt string, params ...any) ([]Data, error)
}

// session represents a standard database session.
type session struct {
	// database handler
	db *Database

	// sql driver and transactioin handlers
	sqlDB *sql.DB
	sqlTX *sql.Tx

	// break event and context-cancel
	breakEvent *events.Event
	ctxBreak   context.CancelFunc
}

// newSession creates new database session.
func newSession(db *Database) (*session, error) {
	if db == nil {
		return nil, ErrDBHandler
	}
	return &session{
		db:         db,
		breakEvent: events.New(),
	}, nil
}

// Database returns the associated database handler
func (s *session) Database() *Database {
	return s.db
}

// Query returns a database query handler.
func (s *session) Query(model Model) Query {
	return newQuery(s, model)

}

// Ping checks if the backend connection is active.
func (s *session) Ping() bool {
	if s.db != nil {
		return s.db.Ping()
	}
	return false
}

// Cancel breaks all active session operation.
func (s *session) Cancel() {
	s.breakEvent.Set()
	if s.ctxBreak != nil {
		s.ctxBreak()
	}
}

// Begin starts a new transactional scope.
func (s *session) Begin() error {
	// already in transaction
	if s.sqlDB != nil && s.sqlTX != nil {
		return nil
	}

	if s.db == nil {
		return ErrDBHandler
	} else if s.db.engine == nil {
		return ErrDBEngine
	}

	if sdb, err := s.db.engine.SqlDB(); err != nil {
		return fmt.Errorf("%w - %v", ErrOperation, err)
	} else {
		s.sqlDB = sdb
	}

	if s.db.Log != nil {
		s.db.Log.Trace("begin new transaction")
	}

	if tx, err := s.sqlDB.Begin(); err != nil {
		return fmt.Errorf("%w - %v", ErrOperation, err)
	} else {
		s.sqlTX = tx
	}
	return nil
}

// Commit runs a commit action in transactional scope.
func (s *session) Commit() error {
	// not in transaction
	if s.sqlDB == nil || s.sqlTX == nil {
		return fmt.Errorf("%w - not in transaction", ErrOperation)
	}

	defer func() {
		s.db.engine.Release(s.sqlDB)
		s.sqlDB = nil
		s.sqlTX = nil
	}()

	if s.db.Log != nil {
		s.db.Log.Trace("transaction commit")
	}

	if err := s.sqlTX.Commit(); err != nil {
		return fmt.Errorf("%w - %v", ErrOperation, err)
	}
	return nil
}

// RollBack runs a rollback action in transactional scope.
func (s *session) RollBack() error {
	// not in transaction
	if s.sqlDB == nil || s.sqlTX == nil {
		return fmt.Errorf("%w - not in transaction", ErrOperation)
	}

	defer func() {
		s.db.engine.Release(s.sqlDB)
		s.sqlDB = nil
		s.sqlTX = nil
	}()

	if s.db.Log != nil {
		s.db.Log.Trace("transaction rollback")
	}

	if err := s.sqlTX.Rollback(); err != nil {
		return fmt.Errorf("%w - %v", ErrOperation, err)
	}
	return nil
}

// Exec runs a query without returning any rows. it takes the statment
// to run and the args are for any placeholder parameters in the query.
func (s *session) Exec(stmt string, params ...any) (int, error) {
	if s.db == nil {
		return 0, ErrDBHandler
	} else if s.db.engine == nil {
		return 0, ErrDBEngine
	}

	if s.db.Log != nil {
		s.db.Log.Debug("SQL:\n---\n%s\nPARAMS: %v\n---", stmt, params)
	}

	var err error
	var sqldb *sql.DB
	var ctx context.Context
	var res sql.Result

	// not in transaction
	if s.sqlDB == nil || s.sqlTX == nil {
		if sqldb, err = s.db.engine.SqlDB(); err != nil {
			return 0, fmt.Errorf("%w - %v", ErrOperation, err)
		}
		defer s.db.engine.Release(sqldb)
	}

	s.breakEvent.Clear()
	if s.db.OperationTimeout > 0 {
		ctx, s.ctxBreak = context.WithDeadline(s.db.ctx, time.Now().Add(
			time.Duration(s.db.OperationTimeout*float64(time.Second))))
	} else {
		ctx, s.ctxBreak = context.WithCancel(s.db.ctx)
	}
	defer s.ctxBreak()

	for {
		if s.sqlDB != nil && s.sqlTX != nil {
			res, err = s.sqlTX.ExecContext(ctx, stmt, params...)
		} else {
			res, err = sqldb.ExecContext(ctx, stmt, params...)
		}
		if err == nil {
			n, err := res.RowsAffected()
			return int(n), err
		} else if err == context.Canceled {
			return 0, ErrBreak
		} else if err == context.DeadlineExceeded {
			return 0, ErrTimeout
		} else if !s.db.engine.CanRetryErr(err) {
			break
		}
		s.breakEvent.Wait(s.db.RetryInterval)
	}

	return 0, fmt.Errorf("%w - %v", ErrOperation, err)
}

// Fetch runs a query that returns rows. it takes the statment
// to run and the args are for any placeholder parameters in the query.
func (s *session) Fetch(stmt string, params ...any) ([]Data, error) {
	if s.db == nil {
		return nil, ErrDBHandler
	} else if s.db.engine == nil {
		return nil, ErrDBEngine
	}

	if s.db.Log != nil {
		s.db.Log.Debug("SQL:\n---\n%s\nPARAMS: %v\n---", stmt, params)
	}

	var err error
	var sqldb *sql.DB
	var ctx context.Context
	var rows *sql.Rows
	var colNames []string

	if sqldb, err = s.db.engine.SqlDB(); err != nil {
		return nil, fmt.Errorf("%w - %v", ErrOperation, err)
	}
	defer s.db.engine.Release(sqldb)

	s.breakEvent.Clear()
	if s.db.OperationTimeout > 0 {
		ctx, s.ctxBreak = context.WithDeadline(s.db.ctx, time.Now().Add(
			time.Duration(s.db.OperationTimeout*float64(time.Second))))
	} else {
		ctx, s.ctxBreak = context.WithCancel(s.db.ctx)
	}
	defer s.ctxBreak()

	for {
		rows, err = sqldb.QueryContext(ctx, stmt, params...)
		if err == nil {
			defer rows.Close()
			colNames, err = rows.Columns()
			if err != nil {
				return nil, fmt.Errorf("%w - %v", ErrOperation, err)
			}
			break
		} else if err == context.Canceled {
			return nil, ErrBreak
		} else if err == context.DeadlineExceeded {
			return nil, ErrTimeout
		} else if !s.db.engine.CanRetryErr(err) {
			return nil, fmt.Errorf("%w - %v", ErrOperation, err)
		}
		s.breakEvent.Wait(s.db.RetryInterval)
	}

	result := []Data{}
	lenCols := len(colNames)

	// create empty slice to represent cols data, and second
	// slice containing pointers to items in cols data slice.
	colsData := make([]any, lenCols)
	colsPtrs := make([]any, lenCols)
	for k := range colsData {
		colsPtrs[k] = &colsData[k]
	}
	for rows.Next() {
		if err := rows.Scan(colsPtrs...); err != nil {
			return nil, fmt.Errorf("%w - %v", ErrOperation, err)
		}
		rowData := Data{}
		// retrieve value for each column from data slice,
		for k, colName := range colNames {
			rowData[colName] = colsData[k]
		}
		result = append(result, rowData)
	}

	return result, nil
}
