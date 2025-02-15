// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"errors"
	"fmt"
)

var (
	// ErrError indicates the parent error.
	ErrError = errors.New("")

	// ErrDBHandler indicates an invalid or not defined database handler.
	ErrDBHandler = fmt.Errorf("%winvalid database handler", ErrError)
	// ErrDBBackend indicates an invalid or not defined database backend.
	ErrDBBackend = fmt.Errorf("%winvalid database backend", ErrError)
	// ErrDBEngine indicates an invalid or not defined database engine.
	ErrDBEngine = fmt.Errorf("%winvalid database engine", ErrError)
	// ErrDBSession indicates an invalid or not defined database session.
	ErrDBSession = fmt.Errorf("%winvalid database session", ErrError)

	// ErrDBConfig indicates an invalid or not defined database config.
	ErrDBConfig = fmt.Errorf("%winvalid database config", ErrError)
	// ErrDBPath indicates an invalid or not defined database path.
	ErrDBPath = fmt.Errorf("%winvalid database path", ErrDBConfig)
	// ErrDBName indicates an invalid or not defined database name.
	ErrDBName = fmt.Errorf("%winvalid database name", ErrDBConfig)
	// ErrDBHost indicates an invalid or not defined database host.
	ErrDBHost = fmt.Errorf("%winvalid database host", ErrDBConfig)
	// ErrDBPort indicates an invalid or not defined database port number.
	ErrDBPort = fmt.Errorf("%winvalid database port", ErrDBConfig)

	// ErrOpen indicates the connection to database failed.
	ErrOpen = fmt.Errorf("%wconnection failed", ErrError)
	// ErrClosed indicates that the database connection is closed.
	ErrClosed = fmt.Errorf("%wconnection closed", ErrError)
	// ErrBreak indicates an operation interruption.
	ErrBreak = fmt.Errorf("%woperation break", ErrError)
	// ErrTimeout indicates that the database operation timed out.
	ErrTimeout = fmt.Errorf("%woperation timeout", ErrError)
	// ErrOperation indicates a database operation error.
	ErrOperation = fmt.Errorf("%woperation error", ErrError)
)
