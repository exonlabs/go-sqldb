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
	"github.com/exonlabs/go-sqldb/pkg/sqlitedb"
)

var (
	db_path   = filepath.Join(os.TempDir(), "sample.db")
	db_config = dictx.Dict{
		"database": db_path,
		// "connect_args": "",
		// "operation_timeout": 5.0,
		// "retry_interval": 0.1,
	}
)

//////////////////////////////// models

type group struct{ sqldb.BaseModel }

var Group = &group{sqldb.BaseModel{
	DefaultTable:  "groups",
	DefaultOrders: []string{"title ASC"},
	AutoGuid:      true,
}}

func (m *group) TableMeta() *sqldb.TableMeta {
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
		Args: dictx.Dict{
			"without_rowid": true,
		},
	}
}

func (m *group) InitialData(dbs *sqldb.Session, _ string) error {
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
		if n, err := dbs.Query(Group).
			FilterBy("title", v["title"]).Count(); err != nil {
			return err
		} else if n > 0 {
			continue
		}
		// create new entry
		_, err := dbs.Query(Group).Insert(v)
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
			{Name: "group_guid", Type: "VARCHAR(32) NOT NULL"},
		},
		Constraints: []sqldb.ConstraintMeta{
			{Definition: "FOREIGN KEY (group_guid) REFERENCES groups (guid) " +
				"ON UPDATE CASCADE ON DELETE RESTRICT"},
		},
		AutoGuid: true,
		Args: dictx.Dict{
			"without_rowid": true,
		},
	}
}

func (m *person) InitialData(dbs *sqldb.Session, _ string) error {
	persons := []sqldb.Data{
		{
			"name":   "Manager",
			"email":  "manager@company.com",
			"active": true,
			"group":  "managers",
		}, {
			"name":   "Employee",
			"email":  "employee@company.com",
			"active": true,
			"group":  "employees",
		}, {
			"name":   "Guest",
			"email":  "",
			"active": false,
			"group":  "visitors",
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

		// check and get group
		grp_title := dictx.GetString(data, "group", "")
		if grp_title == "" {
			return errors.New("invalid data, empty 'group' value")
		}
		grp, err := dbs.Query(Group).FilterBy("title", grp_title).One()
		if err != nil {
			return err
		} else if grp == nil {
			return fmt.Errorf("group not found: %s", grp_title)
		}

		// create new entry
		data["group_guid"] = grp["guid"]
		delete(data, "group")
		_, err = dbs.Query(Person).Insert(data)
		if err != nil {
			return err
		}
	}

	return nil
}

//////////////////////////////// operations

// func print_data(data []sqldb.Data) {
// 	if len(data) > 0 {
// 		keys := data[0].Keys()
// 		for _, item := range data {
// 			for _, k := range keys {
// 				fmt.Printf("%v: %v\n", k, item[k])
// 			}
// 		}
// 	}
// }

func run_initialize(db *sqldb.Database) error {
	metainfo := []sqldb.ModelMeta{
		{Table: Group.DefaultTable, Model: Group},
		{Table: Person.DefaultTable, Model: Person},
	}
	return sqldb.InitializeModels(db, metainfo)
}

func run_operations(_ *sqldb.Database) error {
	// 	// define tables
	// 	tables := map[db.TableName]db.IModel{
	// 		"foobar": &Foobar{},
	// 	}

	// 	if err := db.InitDatabase(tables); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	}
	// 	fmt.Println("\nDB initialize: Done")

	// 	dbs := db.Session()
	// 	defer dbs.Close()

	// 	// checking DB
	// 	fmt.Println("\nAll entries:")
	// 	if items, err := dbs.Query(&Foobar{}).All(); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	} else {
	// 		PrintModelData(items)
	// 	}
	// 	if total, err := dbs.Query(&Foobar{}).Count(); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	} else {
	// 		fmt.Println("\nTotal Items:", total)
	// 	}

	// 	// custom columns
	// 	fmt.Println("\nGet custom columns entries:")
	// 	if items, err := dbs.Query(&Foobar{}).Columns("col1", "col2").
	// 		Limit(2).OrderBy("col1 DESC").All(); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	} else {
	// 		PrintModelData(items)
	// 	}

	// 	// filtered entries
	// 	fmt.Println("\nGet filter columns entries:")
	// 	if items, err := dbs.Query(&Foobar{}).
	// 		Filter("col2 LIKE ? OR col3 IN (?,?)", "description_3", 1, 3).
	// 		All(); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	} else {
	// 		PrintModelData(items)
	// 	}

	// 	// update entries
	// 	fmt.Println("\nModify: first row")
	// 	if _, err := dbs.Query(&Foobar{}).FilterBy("col3", 1).
	// 		Update(db.ModelData{
	// 			"col1": "boo_1", "col2": "boo_2", "col4": 1,
	// 		}); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	}
	// 	fmt.Println("-- After Modify --")
	// 	if items, err := dbs.Query(&Foobar{}).All(); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	} else {
	// 		PrintModelData(items)
	// 	}

	// fmt.Println("\nDelete: first row")
	// if _, err := dbs.Query(&Foobar{}).FilterBy("col3", 1).

	// 		Delete(); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	}

	// fmt.Println("-- After Delete --")

	// 	if items, err := dbs.Query(&Foobar{}).All(); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	} else {

	// 		PrintModelData(items)
	// 	}

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

	// create engine
	engine, err := sqlitedb.NewEngine(db_config)
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
		if err := sqlitedb.InteractiveSetup(db_config); err != nil {
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

	// cleanup
	defer os.Remove(db_path)

	if err := run_operations(db); err != nil {
		log.Info("Error: %s\n", err.Error())
		return
	}

	log.Info("done")
}
