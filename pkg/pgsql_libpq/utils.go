// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package pgsqldb

import (
	"fmt"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/console"
	"github.com/exonlabs/go-utils/pkg/logging"
)

// InteractiveConfig gets the database configuration interactively from console.
// it takes default options and return new input options.
//
// The parsed default options are:
//   - database: (string) the database name - REQUIRED
//   - address: (string) the database server address - REQUIRED
//   - username: (string) the database  access username (if any)
//   - password: (string) the database access password (if any)
//   - connect_args: (string) holds connection params
func InteractiveConfig(defaults dictx.Dict) (dictx.Dict, error) {
	con, err := console.NewTermConsole()
	if err != nil {
		return nil, err
	}
	defer con.Close()

	// get database name
	db_name, err := con.Required().ReadValue(
		"Enter database name", dictx.GetString(defaults, "database", ""))
	if err != nil {
		return nil, err
	}

	// get database address
	db_addr, err := con.Required().ReadValue(
		"Enter database address",
		dictx.GetString(defaults, "address", "localhost:5432"))
	if err != nil {
		return nil, err
	}

	// get database username
	db_user, err := con.ReadValue(
		"Enter database username", dictx.GetString(defaults, "username", ""))
	if err != nil {
		return nil, err
	}

	// get database password
	default_pass := dictx.GetString(defaults, "password", "")
	db_pass, err := con.Hidden().ReadValue(
		"Enter database password", default_pass)
	if err != nil {
		return nil, err
	} else if db_pass != "" && db_pass != default_pass {
		if err = con.Hidden().ConfirmValue(
			"Confirm database password", db_pass); err != nil {
			return nil, err
		}
	}

	// get connect args
	connect_args, err := con.ReadValue(
		"Enter connect args", dictx.GetString(defaults, "connect_args", ""))
	if err != nil {
		return nil, err
	}

	cfg, err := dictx.Clone(defaults)
	if err != nil {
		return nil, err
	}
	dictx.Merge(cfg, dictx.Dict{
		"database":     db_name,
		"address":      db_addr,
		"username":     db_user,
		"password":     db_pass,
		"connect_args": connect_args,
	})

	return cfg, nil
}

// InteractiveSetup performs an interactive console based database setup.
// it takes database options and makes config validation.
//
// The parsed options are:
//   - database: (string) the database name - REQUIRED
//   - address: (string) the database server address - REQUIRED
//   - username: (string) the database  access username (if any)
//   - password: (string) the database access password (if any)
//   - connect_args: (string) holds connection params
func InteractiveSetup(log *logging.Logger, opts dictx.Dict) error {
	con, err := console.NewTermConsole()
	if err != nil {
		return err
	}
	defer con.Close()

	// get database admin username and password
	adm_user, err := con.Required().ReadValue(
		"Enter database admin username", "postgres")
	if err != nil {
		return err
	}
	adm_pass, err := con.Hidden().ReadValue(
		"Enter database admin password", "")
	if err != nil {
		return err
	}

	// create engine
	engine, err := NewEngine(log, opts)
	if err != nil {
		return err
	}

	// store connection access and replace with admin access
	conn_user, conn_pass := engine.cfg.Username, engine.cfg.Password
	engine.cfg.Username = adm_user
	engine.cfg.Password = adm_pass

	// store connection database
	db_name := engine.cfg.Database
	engine.cfg.Database = "postgres"

	// create database handler
	db := sqldb.NewDatabase(log, engine, opts)
	defer db.Shutdown()

	// create new session
	dbs := db.Session()

	var stmt string

	// create role
	stmt = fmt.Sprintf("SELECT * FROM pg_roles WHERE rolname='%s';",
		conn_user)
	if res, err := dbs.Fetch(stmt); err != nil {
		return err
	} else if len(res) <= 0 {
		stmt = fmt.Sprintf("CREATE ROLE %s LOGIN PASSWORD '%s';",
			conn_user, conn_pass)
	} else {
		stmt = fmt.Sprintf("ALTER ROLE %s WITH PASSWORD '%s';",
			conn_user, conn_pass)
	}
	if _, err := dbs.Exec(stmt); err != nil {
		return err
	}

	// create database
	stmt = fmt.Sprintf("SELECT * FROM pg_database WHERE datname='%s';",
		db_name)
	if res, err := dbs.Fetch(stmt); err != nil {
		return err
	} else if len(res) <= 0 {
		stmt = fmt.Sprintf("CREATE DATABASE %s OWNER %s ENCODING 'UTF8';",
			db_name, conn_user)
		if _, err := dbs.Exec(stmt); err != nil {
			return err
		}
	}

	// grant privileges
	stmt = fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;",
		db_name, conn_user)
	if _, err := dbs.Exec(stmt); err != nil {
		return err
	}

	return nil
}
