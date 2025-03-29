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

// Session represents a standard database Session.
type Session struct {
	// database handler
	db *Database

	// sql driver and transactioin handlers
	sdb *sql.DB
	stx *sql.Tx

	// break event and context-cancel
	breakEvent *events.Event
	ctxBreak   context.CancelFunc

	// OperationTimeout sets the timeout in seconds for operations.
	// use 0 or negative value to disable operation timeout. (default 5.0 sec)
	OperationTimeout float64
	// RetryInterval sets the interval in seconds between operation retries.
	// trials are done untill operation is done or timeout is reached.
	// retry interval value must be > 0. (default 0.1 sec)
	RetryInterval float64
}

// NewSession creates new database session.
func NewSession(db *Database) *Session {
	s := &Session{
		db:         db,
		breakEvent: events.New(),
	}
	if db != nil {
		s.OperationTimeout = db.OperationTimeout
		s.RetryInterval = db.RetryInterval
	}
	return s
}

// Database returns the associated database handler
func (s *Session) Database() *Database {
	return s.db
}

// Query returns a database query handler.
func (s *Session) Query(model Model) *Query {
	return NewQuery(s, model)
}

// Ping checks if the backend connection is active.
func (s *Session) Ping() bool {
	if s.sdb != nil {
		return s.sdb.Ping() == nil
	}
	return s.db.Ping()
}

// Cancel breaks all active session operation.
func (s *Session) Cancel() {
	s.breakEvent.Set()
	if s.ctxBreak != nil {
		s.ctxBreak()
	}
}

// check attrs before running query
func (s *Session) check_run() error {
	if s.db == nil {
		return ErrDBHandler
	}
	return s.db.check_run()
}

// Begin starts a new transactional scope.
func (s *Session) Begin() error {
	// already in transaction
	if s.sdb != nil && s.stx != nil {
		return nil
	}

	if err := s.check_run(); err != nil {
		return err
	}

	if sdb, err := s.db.engine.SqlDB(); err != nil {
		return fmt.Errorf("%w - %v", ErrOpen, err)
	} else {
		s.sdb = sdb
	}

	if s.db.Log != nil {
		s.db.Log.Trace("begin new transaction")
	}
	if stx, err := s.sdb.Begin(); err != nil {
		return fmt.Errorf("%w - %v", ErrOperation, err)
	} else {
		s.stx = stx
	}
	return nil
}

// Commit runs a commit action in transactional scope.
func (s *Session) Commit() error {
	// not in transaction
	if s.sdb == nil || s.stx == nil {
		return fmt.Errorf("%w - not in transaction", ErrOperation)
	}

	if err := s.check_run(); err != nil {
		return err
	}

	defer func() {
		s.db.engine.Release(s.sdb)
		s.sdb, s.stx = nil, nil
	}()

	if s.db.Log != nil {
		s.db.Log.Trace("transaction commit")
	}
	if err := s.stx.Commit(); err != nil {
		return fmt.Errorf("%w - %v", ErrOperation, err)
	}
	return nil
}

// RollBack runs a rollback action in transactional scope.
func (s *Session) RollBack() error {
	// not in transaction
	if s.sdb == nil || s.stx == nil {
		return fmt.Errorf("%w - not in transaction", ErrOperation)
	}

	if err := s.check_run(); err != nil {
		return err
	}

	defer func() {
		s.db.engine.Release(s.sdb)
		s.sdb, s.stx = nil, nil
	}()

	if s.db.Log != nil {
		s.db.Log.Trace("transaction rollback")
	}
	if err := s.stx.Rollback(); err != nil {
		return fmt.Errorf("%w - %v", ErrOperation, err)
	}
	return nil
}

// Exec runs a query without returning any rows. it takes the statment
// to run and the args are for any placeholder parameters in the query.
func (s *Session) Exec(stmt string, params ...any) (int, error) {
	if err := s.check_run(); err != nil {
		return 0, err
	}

	if len(params) > 0 {
		stmt = s.db.engine.SqlGenerator().FormatStmt(stmt)
	}

	if s.db.Log != nil {
		if len(params) > 0 {
			s.db.Log.Trace("SQL: %s %v", stmt, params)
		} else {
			s.db.Log.Trace("SQL: %s", stmt)
		}
	}

	var err error
	var sdb *sql.DB
	var ctx context.Context
	var res sql.Result

	// not in transaction
	if s.sdb == nil || s.stx == nil {
		if sdb, err = s.db.engine.SqlDB(); err != nil {
			return 0, fmt.Errorf("%w - %v", ErrOpen, err)
		}
		defer s.db.engine.Release(sdb)
	}

	s.breakEvent.Clear()
	if s.OperationTimeout > 0 {
		ctx, s.ctxBreak = context.WithDeadline(s.db.ctx, time.Now().Add(
			time.Duration(s.OperationTimeout*float64(time.Second))))
	} else {
		ctx, s.ctxBreak = context.WithCancel(s.db.ctx)
	}
	defer s.ctxBreak()

	var lastErr error
	for {
		if s.sdb != nil && s.stx != nil {
			res, err = s.stx.ExecContext(ctx, stmt, params...)
		} else {
			res, err = sdb.ExecContext(ctx, stmt, params...)
		}
		if err == nil {
			n, err := res.RowsAffected()
			return int(n), err
		} else if err == context.Canceled {
			return 0, ErrBreak
		} else if err == context.DeadlineExceeded {
			return 0, fmt.Errorf("%w - %v", ErrTimeout, lastErr)
		} else {
			lastErr = err
			if !s.db.engine.CanRetryErr(err) {
				break
			}
		}
		s.breakEvent.Wait(s.RetryInterval)
	}

	return 0, fmt.Errorf("%w - %v", ErrOperation, err)
}

// Fetch runs a query that returns rows. it takes the statment
// to run and the args are for any placeholder parameters in the query.
func (s *Session) Fetch(stmt string, params ...any) ([]Data, error) {
	if err := s.check_run(); err != nil {
		return nil, err
	}

	if len(params) > 0 {
		stmt = s.db.engine.SqlGenerator().FormatStmt(stmt)
	}

	if s.db.Log != nil {
		if len(params) > 0 {
			s.db.Log.Trace("SQL: %s %v", stmt, params)
		} else {
			s.db.Log.Trace("SQL: %s", stmt)
		}
	}

	var err error
	var sdb *sql.DB
	var ctx context.Context
	var rows *sql.Rows
	var colNames []string

	if sdb, err = s.db.engine.SqlDB(); err != nil {
		return nil, fmt.Errorf("%w - %v", ErrOpen, err)
	}
	defer s.db.engine.Release(sdb)

	s.breakEvent.Clear()
	if s.OperationTimeout > 0 {
		ctx, s.ctxBreak = context.WithDeadline(s.db.ctx, time.Now().Add(
			time.Duration(s.OperationTimeout*float64(time.Second))))
	} else {
		ctx, s.ctxBreak = context.WithCancel(s.db.ctx)
	}
	defer s.ctxBreak()

	var lastErr error
	for {
		rows, err = sdb.QueryContext(ctx, stmt, params...)
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
			return nil, fmt.Errorf("%w - %v", ErrTimeout, lastErr)
		} else {
			lastErr = err
			if !s.db.engine.CanRetryErr(err) {
				return nil, fmt.Errorf("%w - %v", ErrOperation, err)
			}
		}
		s.breakEvent.Wait(s.RetryInterval)
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
