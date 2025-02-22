// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"database/sql"
	"fmt"
)

// Session represents the database session object.
type Session struct {
	// parent database handler
	db *Database
	// sql transactioin handler
	sqlTX *sql.Tx
}

// creates new database session object.
func newSession(db *Database) (*Session, error) {
	if db.DBLog != nil {
		db.DBLog.Debug("new session")
	}
	return &Session{
		db: db,
	}, nil
}

// Query creates new query object on current session.
func (s *Session) Query(model Model) *Query {
	return NewQuery(s, model)
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

// Commit executes a commit action in transactional scope.
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

// RollBack executes a rollback action in transactional scope.
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

// Execute executes a query without returning any rows. it takes the statment
// to run and the args are for any placeholder parameters in the query.
func (s *Session) Execute(stmt string, params ...any) (int, error) {
	if s.db == nil {
		return 0, ErrDBHandler
	}
	if s.db.DBLog != nil {
		s.db.DBLog.Debug("SQL:\n---\n%s\nPARAMS: %v\n---", stmt, params)
	}

	var err error
	var res sql.Result

	// TODO: add retries

	if s.sqlTX != nil {
		res, err = s.sqlTX.Exec(stmt, params...)
	} else {
		res, err = s.db.engine.SqlDB().Exec(stmt, params...)
	}
	if err == nil {
		n, err := res.RowsAffected()
		return int(n), err
	}

	return 0, fmt.Errorf("%w - %v", ErrOperation, err)
}

// FetchAll executes a query that returns rows. it takes the statment
// to run and the args are for any placeholder parameters in the query.
func (s *Session) FetchAll(stmt string, params ...any) ([]Data, error) {
	if s.db == nil {
		return nil, ErrDBHandler
	}
	if s.db.DBLog != nil {
		s.db.DBLog.Debug("SQL:\n---\n%s\nPARAMS: %v\n---", stmt, params)
	}

	var err error
	var rows *sql.Rows
	var colNames []string

	// TODO: add retries

	rows, err = s.db.engine.SqlDB().Query(stmt, params...)
	if err == nil {
		defer rows.Close()
		colNames, err = rows.Columns()
	}
	if err != nil {
		return nil, fmt.Errorf("%w - %v", ErrOperation, err)
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
