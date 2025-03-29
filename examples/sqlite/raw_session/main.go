// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/logging"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	sqlitedb "github.com/exonlabs/go-sqldb/pkg/sqlite_mattn"
)

var (
	db_path   = filepath.Join(os.TempDir(), "sample.db")
	db_config = dictx.Dict{
		"database": db_path,
		// "connect_args": "",
		// "operation_timeout": 3.0,
		// "retry_interval": 0.1,
	}
)

//////////////////////////////// operations

func run_initialize(db *sqldb.Database) error {
	if !db.Ping() {
		return errors.New("database connection down")
	}

	// create new session
	dbs := db.Session()

	schema := `
CREATE TABLE IF NOT EXISTS roles (
  guid VARCHAR(32) NOT NULL,
  title VARCHAR(128) NOT NULL,
  description TEXT,
  PRIMARY KEY (guid)
) WITHOUT ROWID;
CREATE UNIQUE INDEX IF NOT EXISTS ix_roles_guid ON roles (guid);
CREATE UNIQUE INDEX IF NOT EXISTS ix_roles_title ON roles (title);`
	if _, err := dbs.Exec(schema); err != nil {
		return err
	}

	return nil
}

func run_operations(db *sqldb.Database) error {
	// create new session
	dbs := db.Session()

	var stmt string
	var params []any

	// add new role
	fmt.Println("\n* Adding new role: role1")
	stmt = "INSERT INTO roles VALUES (?, ?, ?)"
	params = []any{"123456", "role1", "role1 description"}
	if _, err := dbs.Exec(stmt, params...); err != nil {
		return err
	}

	// listing all roles
	fmt.Println("\n* List all roles:")
	stmt = "SELECT * FROM roles;"
	roles, err := dbs.Fetch(stmt)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		for _, v := range roles {
			fmt.Println("  - " + dictx.String(v))
		}
		fmt.Printf("Total: %d\n", len(roles))
	}

	return nil
}

////////////////////////////////

func main() {
	log := logging.NewStdoutLogger("main")
	dblog := log.SubLogger("db")

	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			indx := bytes.Index(stack, []byte("panic({"))
			log.Panic("%s", r)
			log.Trace("\n----------\n%s----------", stack[indx:])
			os.Exit(1)
		}
	}()

	debug0 := flag.Bool("x", false, "\nenable debug logs")
	debug1 := flag.Bool("xx", false, "enable debug and trace logs")
	setup := flag.Bool("setup", false, "perform database setup")
	clean := flag.Bool("clean", false, "perform database cleanup")
	flag.Parse()

	switch {
	case *debug1:
		log.Level = logging.TRACE
		dblog.Level = logging.TRACE
	case *debug0:
		log.Level = logging.DEBUG
		dblog.Level = logging.DEBUG
	default:
		dblog = nil
	}

	log.Info("**** starting ****")

	log.Info("Using database: %s", db_path)
	fmt.Println()

	log.Info("Using Options:")
	log.Info("%s", dictx.String(db_config))
	fmt.Println()

	if *clean {
		fmt.Println("* Clean-Up database:")
		os.Remove(db_path)
		fmt.Println()
		log.Info("done")
		return
	}

	// create engine
	engine, err := sqlitedb.NewEngine(dblog, db_config)
	if err != nil {
		log.Error("create engine failed - %s", err)
		return
	}

	// create database handler
	db := sqldb.NewDatabase(dblog, engine, db_config)
	defer db.Shutdown()

	// setup database
	if *setup {
		fmt.Println("* Setup database:")

		// setup database
		if err := sqlitedb.InteractiveSetup(dblog, db_config); err != nil {
			if !strings.Contains(err.Error(), "EOF") {
				fmt.Printf("Error: %s\n", err)
			}
			fmt.Println()
			return
		}
		fmt.Println()

		// initialize database
		if err := run_initialize(db); err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Println()

		log.Info("done")
		return
	}

	if err := run_operations(db); err != nil {
		log.Info("Error: %s\n", err.Error())
		return
	}
	fmt.Println()

	log.Info("done")
}
