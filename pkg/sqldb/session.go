// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/exonlabs/go-utils/pkg/events"
)

// Session represents the database session object.
type Session struct {
	// parent database handler
	db *Database
	// sql transactioin handler
	sqlTX *sql.Tx
	// ctxBreak signals a break operation.
	ctxBreak   context.CancelFunc
	breakEvent *events.Event
}

// creates new database session object.
func newSession(db *Database) (*Session, error) {
	if db.DBLog != nil {
		db.DBLog.Debug("new session")
	}
	return &Session{
		db:         db,
		breakEvent: events.New(),
	}, nil
}

// Query creates new query object on current session.
func (s *Session) Query(model Model) *Query {
	return NewQuery(s, model)
}

// Breaks active session operation.
func (s *Session) Break() {
	s.breakEvent.Set()
	if s.ctxBreak != nil {
		s.ctxBreak()
	}
}

// Begin starts a new transactional scope.
func (s *Session) Begin() error {
	if s.db == nil {
		return ErrDBHandler
	}
	if s.db.DBLog != nil {
		s.db.DBLog.Trace("begin new transaction")
	}

	if tx, err := s.db.engine.SqlDB().Begin(); err != nil {
		return fmt.Errorf("%w - %v", ErrOperation, err)
	} else {
		s.sqlTX = tx
	}
	return nil
}

// Commit runs a commit action in transactional scope.
func (s *Session) Commit() error {
	// not in transaction
	if s.sqlTX == nil {
		return nil
	}
	if s.db.DBLog != nil {
		s.db.DBLog.Trace("transaction commit")
	}

	if err := s.sqlTX.Commit(); err != nil {
		return fmt.Errorf("%w - %v", ErrOperation, err)
	}
	s.sqlTX = nil
	return nil
}

// RollBack runs a rollback action in transactional scope.
func (s *Session) RollBack() error {
	// not in transaction
	if s.sqlTX == nil {
		return nil
	}
	if s.db.DBLog != nil {
		s.db.DBLog.Trace("transaction rollback")
	}

	if err := s.sqlTX.Rollback(); err != nil {
		return fmt.Errorf("%w - %v", ErrOperation, err)
	}
	s.sqlTX = nil
	return nil
}

// Executes runs a query without returning any rows. it takes the statment
// to run and the args are for any placeholder parameters in the query.
func (s *Session) Execute(stmt string, params ...any) (int, error) {
	if s.db == nil {
		return 0, ErrDBHandler
	}
	if s.db.DBLog != nil {
		s.db.DBLog.Debug("SQL:\n---\n%s\nPARAMS: %v\n---", stmt, params)
	}

	var err error
	var ctx context.Context
	var res sql.Result

	s.breakEvent.Clear()
	ctx, s.ctxBreak = context.WithDeadline(s.db.ctx, time.Now().Add(
		time.Duration(s.db.OperationTimeout*float64(time.Second))))
	defer s.ctxBreak()

	for {
		if s.sqlTX != nil {
			res, err = s.sqlTX.ExecContext(ctx, stmt, params...)
		} else {
			res, err = s.db.engine.SqlDB().ExecContext(ctx, stmt, params...)
		}
		if err == nil {
			n, err := res.RowsAffected()
			return int(n), err
		} else if !s.db.engine.CanRetryErr(err) {
			break
		}

		s.breakEvent.Wait(s.db.RetryInterval)

		if errors.Is(ctx.Err(), context.Canceled) {
			return 0, ErrBreak
		} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return 0, ErrTimeout
		}
	}

	return 0, fmt.Errorf("%w - %v", ErrOperation, err)
}

// FetchAll runs a query that returns rows. it takes the statment
// to run and the args are for any placeholder parameters in the query.
func (s *Session) FetchAll(stmt string, params ...any) ([]Data, error) {
	if s.db == nil {
		return nil, ErrDBHandler
	}
	if s.db.DBLog != nil {
		s.db.DBLog.Debug("SQL:\n---\n%s\nPARAMS: %v\n---", stmt, params)
	}

	var err error
	var ctx context.Context
	var rows *sql.Rows
	var colNames []string

	s.breakEvent.Clear()
	ctx, s.ctxBreak = context.WithDeadline(s.db.ctx, time.Now().Add(
		time.Duration(s.db.OperationTimeout*float64(time.Second))))
	defer s.ctxBreak()

	for {
		rows, err = s.db.engine.SqlDB().QueryContext(ctx, stmt, params...)
		if err == nil {
			defer rows.Close()
			colNames, err = rows.Columns()
			if err != nil {
				return nil, fmt.Errorf("%w - %v", ErrOperation, err)
			}
			break
		} else if !s.db.engine.CanRetryErr(err) {
			return nil, fmt.Errorf("%w - %v", ErrOperation, err)
		}

		s.breakEvent.Wait(s.db.RetryInterval)

		if errors.Is(ctx.Err(), context.Canceled) {
			return nil, ErrBreak
		} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, ErrTimeout
		}
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
