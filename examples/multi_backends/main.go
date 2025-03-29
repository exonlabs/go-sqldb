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
	"sort"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/abc/slicex"
	"github.com/exonlabs/go-utils/pkg/logging"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"

	mssqldb "github.com/exonlabs/go-sqldb/pkg/mssql_microsoft"
	mysqldb "github.com/exonlabs/go-sqldb/pkg/mysql_sqldriver"
	pgsqldb "github.com/exonlabs/go-sqldb/pkg/pgsql_libpq"
	sqlitedb "github.com/exonlabs/go-sqldb/pkg/sqlite_mattn"
)

var (
	BACKENDS = []string{"sqlite", "mysql", "pgsql", "mssql"}

	db_name   = "testdb"
	db_path   = os.TempDir()
	db_config = dictx.Dict{
		"database": db_name,
		"username": "test123",
		"password": "test123",
		// "connect_args": "",
		// "operation_timeout": 3.0,
		// "retry_interval": 0.1,
	}
	meta_args = dictx.Dict{
		"sqlite_without_rowid": true,
		// "mysql_storage_engine": "InnoDB",
		// "disable_table_exists": true,
		// "disable_index_exists": true,
	}
)

//////////////////////////////// models

type role struct{ sqldb.BaseModel }

var Role = &role{sqldb.BaseModel{
	DefaultTable:  "roles",
	DefaultOrders: []string{"title ASC"},
	AutoGuid:      true,
}}

func (m *role) TableMeta() *sqldb.TableMeta {
	return &sqldb.TableMeta{
		Columns: []sqldb.ColumnMeta{
			{Name: "title", Type: "VARCHAR(128) NOT NULL",
				Unique: true, Index: true},
			{Name: "description", Type: "TEXT"},
			{Name: "access_level", Type: "INTEGER"},
			{Name: "public_join", Type: "BOOLEAN DEFAULT false"},
		},
		Constraints: []sqldb.ConstraintMeta{
			{Definition: "CHECK (access_level>=1 AND access_level<=5)"},
		},
		AutoGuid: true,
		Args:     meta_args,
	}
}

type person struct{ sqldb.BaseModel }

var Person *person = &person{sqldb.BaseModel{
	DefaultTable:  "persons",
	DefaultOrders: []string{"name ASC"},
	AutoGuid:      true,
}}

func (m *person) TableMeta() *sqldb.TableMeta {
	return &sqldb.TableMeta{
		Columns: []sqldb.ColumnMeta{
			{Name: "name", Type: "VARCHAR(128) NOT NULL",
				Unique: true, Index: true},
			{Name: "email", Type: "VARCHAR(256)"},
			{Name: "active", Type: "BOOLEAN DEFAULT true"},
			{Name: "role_guid", Type: "VARCHAR(32) NOT NULL"},
		},
		Constraints: []sqldb.ConstraintMeta{
			{Definition: "FOREIGN KEY (role_guid) REFERENCES roles (guid) " +
				"ON UPDATE CASCADE ON DELETE RESTRICT"},
		},
		AutoGuid: true,
		Args:     meta_args,
	}
}

//////////////////////////////// operations

func print_data(d dictx.Dict) string {
	s := ""
	keys := []string{}
	for k := range d {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		switch d[k].(type) {
		case byte, []byte, string, []string, rune, []rune:
			s += fmt.Sprintf("%s: %s, ", k, d[k])
		default:
			s += fmt.Sprintf("%s: %v, ", k, d[k])
		}
	}
	if len(s) > 0 {
		s = s[:len(s)-2] // Remove the trailing ", "
	}
	return "{" + s + "}"
}

func run_initialize(db *sqldb.Database) error {
	metainfo := []sqldb.ModelMeta{
		{Table: Role.DefaultTable, Model: Role},
		{Table: Person.DefaultTable, Model: Person},
	}
	return sqldb.InitializeModels(db, metainfo)
}

func run_operations(db *sqldb.Database) error {
	if !db.Ping() {
		return errors.New("database connection down")
	}

	// create new session
	dbs := db.Session()

	// cleanup
	fmt.Println("\n* cleanup")
	if _, err := dbs.Query(Person).Delete(); err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	if _, err := dbs.Query(Role).Delete(); err != nil {
		fmt.Println("ERROR:", err.Error())
	}

	// add new role
	fmt.Println("\n* Adding new role: role1")
	role_guid, err := dbs.Query(Role).Insert(sqldb.Data{
		"title":       "role1",
		"description": "role number 1",
	})
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	if role, err := dbs.Query(Role).GetGuid(role_guid); err != nil {
		fmt.Println("ERROR:", err.Error())
	} else if role == nil {
		fmt.Println("adding role failed")
	} else {
		fmt.Println("  - " + print_data(role))
	}

	// Add new person to role1
	fmt.Println("\n* Adding new person: person1")
	_, err = dbs.Query(Person).Insert(sqldb.Data{
		"name":      "person1",
		"active":    true,
		"role_guid": role_guid,
	})
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	person, err := dbs.Query(Person).FilterBy("name", "person1").One()
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	} else if person == nil {
		fmt.Println("adding person failed")
	} else {
		fmt.Println("  - " + print_data(person))
	}

	// listing all roles
	fmt.Println("\n* List all roles:")
	if roles, err := dbs.Query(Role).All(); err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		for _, v := range roles {
			fmt.Println("  - " + print_data(v))
		}
		fmt.Printf("Total: %d\n", len(roles))
	}

	// listing all persons
	fmt.Println("\n* List all persons:")
	if persons, err := dbs.Query(Person).All(); err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		for _, v := range persons {
			fmt.Println("  - " + print_data(v))
		}
		fmt.Printf("Total: %d\n", len(persons))
	}

	return nil
}

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
	backend := flag.String("backend", "",
		fmt.Sprintf("select backend {%s}", strings.Join(BACKENDS, "|")))
	setup := flag.Bool("setup", false, "perform database setup")
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

	// check backend
	if slicex.Index(BACKENDS, *backend) < 0 {
		fmt.Printf("Error: invalid backend '%s'\n", *backend)
		return
	}

	log.Info("**** starting ****")

	log.Info("Using Backend: %s", *backend)
	fmt.Println()

	var err error

	// setting backend config
	fmt.Println("* Configure database:")
	switch *backend {
	case "sqlite":
		dictx.Set(db_config, "database", filepath.Join(db_path, db_name+".db"))
		db_config, err = sqlitedb.InteractiveConfig(db_config)
	case "mysql":
		db_config, err = mysqldb.InteractiveConfig(db_config)
	case "pgsql":
		db_config, err = pgsqldb.InteractiveConfig(db_config)
	case "mssql":
		db_config, err = mssqldb.InteractiveConfig(db_config)
	}
	if err != nil {
		if !strings.Contains(err.Error(), "EOF") {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Println()
		return
	}
	fmt.Println()

	log.Info("Using Options:")
	log.Info("%s", print_data(db_config))
	fmt.Println()

	// create engine
	var engine sqldb.Engine
	switch *backend {
	case "sqlite":
		engine, err = sqlitedb.NewEngine(dblog, db_config)
	case "mysql":
		engine, err = mysqldb.NewEngine(dblog, db_config)
	case "pgsql":
		engine, err = pgsqldb.NewEngine(dblog, db_config)
	case "mssql":
		engine, err = mssqldb.NewEngine(dblog, db_config)
	}
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

		switch *backend {
		case "sqlite":
			err = sqlitedb.InteractiveSetup(dblog, db_config)
		case "mysql":
			err = mysqldb.InteractiveSetup(dblog, db_config)
		case "pgsql":
			err = pgsqldb.InteractiveSetup(dblog, db_config)
		case "mssql":
			err = mssqldb.InteractiveSetup(dblog, db_config)
		}
		if err != nil {
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
