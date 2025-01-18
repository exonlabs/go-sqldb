// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"database/sql"
	"sync"
)

type Session struct {
	// parent database handler
	handler *Handler

	// sql transactioin handler
	sqlTX *sql.Tx

	// operations sync mutex
	mu sync.Mutex
}

func newSession(dbh *Handler) *Session {
	return &Session{
		handler: dbh,
	}
}

func (dbs *Session) IsActive() bool {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	// if dbs.sqlDB != nil {
	// 	return dbs.sqlDB.Ping() == nil
	// }
	return false
}

func (dbs *Session) Connect() error {
	if dbs.handler == nil {
		return ErrDBHandler
	}

	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	return nil
}

func (dbs *Session) Close() error {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	return nil
}

func (dbs *Session) Begin() error {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	return nil
}

func (dbs *Session) Commit() error {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	return nil
}

func (dbs *Session) RollBack() error {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	return nil
}

func (dbs *Session) Execute(stmt string, params ...any) (int64, error) {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	return 0, nil
}

func (dbs *Session) FetchAll(stmt string, params ...any) ([]Data, error) {
	dbs.mu.Lock()
	defer dbs.mu.Unlock()

	return nil, nil
}
