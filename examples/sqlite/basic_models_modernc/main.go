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
	"github.com/exonlabs/go-utils/pkg/logging"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	sqlitedb "github.com/exonlabs/go-sqldb/pkg/sqlite_modernc"
)

var (
	db_path   = filepath.Join(os.TempDir(), "sample.db")
	db_config = dictx.Dict{
		"database": db_path,
		// "connect_args": "",
		// "operation_timeout": 3.0,
		// "retry_interval": 0.1,
	}
	meta_args = dictx.Dict{
		"sqlite_without_rowid": true,
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

func (m *role) InitialData(dbs *sqldb.Session, _ string) error {
	buff := []sqldb.Data{
		{
			"title":        "managers",
			"description":  "company managers",
			"access_level": 5,
			"public_join":  false,
		}, {
			"title":        "employees",
			"description":  "all company employees",
			"access_level": 3,
			"public_join":  false,
		}, {
			"title":        "visitors",
			"description":  "company guests and visitors",
			"access_level": 1,
			"public_join":  true,
		},
	}

	for _, v := range buff {
		// check if already exists
		if n, err := dbs.Query(Role).
			FilterBy("title", v["title"]).Count(); err != nil {
			return err
		} else if n > 0 {
			continue
		}
		// create new entry
		_, err := dbs.Query(Role).Insert(v)
		if err != nil {
			return err
		}
	}

	return nil
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

func (m *person) InitialData(dbs *sqldb.Session, _ string) error {
	persons := []sqldb.Data{
		{
			"name":   "Manager",
			"email":  "manager@company.com",
			"active": true,
			"role":   "managers",
		}, {
			"name":   "Employee",
			"email":  "employee@company.com",
			"active": true,
			"role":   "employees",
		}, {
			"name":   "Guest",
			"email":  "",
			"active": false,
			"role":   "visitors",
		},
	}

	for _, data := range persons {
		// check if already exists
		if name := dictx.GetString(data, "name", ""); name == "" {
			return errors.New("invalid data, empty 'name' value")
		} else {
			if n, err := dbs.Query(Person).
				FilterBy("name", name).Count(); err != nil {
				return err
			} else if n > 0 {
				continue
			}
		}

		// check and get role
		role_title := dictx.GetString(data, "role", "")
		if role_title == "" {
			return errors.New("invalid data, empty 'role' value")
		}
		role, err := dbs.Query(Role).FilterBy("title", role_title).One()
		if err != nil {
			return err
		} else if role == nil {
			return fmt.Errorf("role not found: %s", role_title)
		}

		// create new entry
		data["role_guid"] = role["guid"]
		delete(data, "role")
		_, err = dbs.Query(Person).Insert(data)
		if err != nil {
			return err
		}
	}

	return nil
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
	person_guid, err := dbs.Query(Person).Insert(sqldb.Data{
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

	// update person info
	fmt.Println("\n* Update person1 email")
	if err := dbs.Query(Person).UpdateGuid(person_guid, sqldb.Data{
		"email":  "person1@domain",
		"active": false,
	}); err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	if persons, err := dbs.Query(Person).Columns("name", "email").
		All(); err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		for _, v := range persons {
			fmt.Println("  - " + print_data(v))
		}
	}

	// update all persons info
	fmt.Println("\n* Update all persons info (email and active values)")
	if _, err = dbs.Query(Person).Update(sqldb.Data{
		"email":  "",
		"active": false,
	}); err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	if persons, err := dbs.Query(Person).Columns("name", "email", "active").
		All(); err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		for _, v := range persons {
			fmt.Println("  - " + print_data(v))
		}
		count, _ := dbs.Query(Person).Count()
		fmt.Printf("Total: %d\n", count)
	}

	// delete all persons
	fmt.Println("\n* Delete all persons then roles")
	if _, err := dbs.Query(Person).Delete(); err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	if count, err := dbs.Query(Person).Count(); err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		fmt.Printf("Total persons: %d\n", count)
	}
	if _, err := dbs.Query(Role).Delete(); err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	if count, err := dbs.Query(Role).Count(); err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		fmt.Printf("Total roles: %d\n", count)
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
	log.Info("%s", print_data(db_config))
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
