// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
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
		// "connect_timeout":  5.0,
		// "connect_interval": 0.5,
	}
)

//////////////////////////////// models

type role struct{ sqldb.BaseModel }

func Role() *role {
	return &role{sqldb.BaseModel{
		DefaultTable:  "roles",
		DefaultOrders: []string{"title ASC"},
		AutoGuid:      true,
	}}
}

type roleMeta struct{ sqldb.BaseModelMeta }

func RoleMeta() *roleMeta {
	return &roleMeta{sqldb.BaseModelMeta{
		Columns: []sqldb.ColumnMeta{
			{Name: "title", Type: "VARCHAR(128) NOT NULL", Unique: true, Index: true},
			{Name: "description", Type: "TEXT"},
			{Name: "access_level", Type: "INTEGER"},
			{Name: "builtin", Type: "BOOLEAN DEFAULT 0"},
		},
		Constraints: []sqldb.ConstraintMeta{
			{Definition: "CHECK (access_level>=0 AND access_level<=5)"},
		},
		AutoGuid: true,
		Args: dictx.Dict{
			"sqlite_without_rowid": true,
		},
	}}
}

func (*roleMeta) InitialData(dbs *sqldb.Session, _ string) error {
	// check if default 'Administrator' role already exist
	num, err := sqldb.NewQuery(dbs, Role()).Filter("title=?", "Administrator").Count()
	if err != nil || num > 0 {
		return err
	}

	// create default 'Administrator' role
	_, err = sqldb.NewQuery(dbs, Role()).Insert(sqldb.Data{
		"title":        "Administrator",
		"description":  "Administrator Full Access",
		"access_level": 5,
		"builtin":      true,
	})
	return err
}

type user struct{ sqldb.BaseModel }

func User() *user {
	return &user{sqldb.BaseModel{
		DefaultTable:  "users",
		DefaultOrders: []string{"username ASC"},
		AutoGuid:      true,
	}}
}

type userMeta struct{ sqldb.BaseModelMeta }

func UserMeta() *userMeta {
	return &userMeta{sqldb.BaseModelMeta{
		Columns: []sqldb.ColumnMeta{
			{Name: "username", Type: "VARCHAR(128) NOT NULL",
				Unique: true, Index: true},
			{Name: "password", Type: "VARCHAR(128) NOT NULL"},
			{Name: "enabled", Type: "BOOLEAN DEFAULT 1"},
			{Name: "role_guid", Type: "VARCHAR(32) NOT NULL"},
		},
		Constraints: []sqldb.ConstraintMeta{
			{Definition: "FOREIGN KEY (role_guid) REFERENCES roles (guid) " +
				"ON UPDATE CASCADE ON DELETE RESTRICT"},
		},
		AutoGuid: true,
		Args: dictx.Dict{
			"sqlite_without_rowid": true,
		},
	}}
}

func (*userMeta) InitialData(dbs *sqldb.Session, _ string) error {
	// if dbs == nil {
	// 	return errors.New("invalid database session")
	// }

	// // check if default 'Admin' user already exist
	// num, err := dbs.Query(User).Filter("username=?", "admin").Count()
	// if err != nil || num > 0 {
	// 	return err
	// }

	// // get default 'Administrator' role
	// role, err := sqldb.NewQuery(dbs, Role()).Filter("title=?", "Administrator").One()
	// if err != nil {
	// 	return err
	// } else if role == nil {
	// 	return errors.New("default 'Administrator' role not found")
	// }
	// role_guid := role.GetString("guid", "")
	// if len(role_guid) == 0 {
	// 	return errors.New("invalid empty 'Administrator' role guid")
	// }

	// // create default 'Admin' user
	// _, err = dbs.Query(User).Insert(sqldb.Data{
	// 	"username":  "admin",
	// 	"password":  "12345",
	// 	"enabled":   true,
	// 	"role_guid": role_guid,
	// })
	// return err
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
	metainfo := map[string]sqldb.ModelMeta{
		Role().DefaultTable: RoleMeta(),
		User().DefaultTable: UserMeta(),
	}
	return sqldb.InitializeDatabase(db, metainfo)
}

func run_operations(db *sqldb.Database) error {
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
	dbLog := log.SubLogger("db")

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
		dbLog.Level = logging.TRACE
	case *debug0:
		log.Level = logging.DEBUG
		dbLog.Level = logging.DEBUG
	default:
		dbLog = nil
	}

	log.Info("**** starting ****")

	log.Info("Using database: %s", db_path)
	fmt.Println()

	log.Info("Using Options:")
	log.Info("%s", dictx.String(db_config))
	fmt.Println()

	// create engine
	db_engine, err := sqlitedb.NewEngine(db_config)
	if err != nil {
		log.Error("create engine failed - %s", err)
		return
	}

	// create database handler
	db, err := sqldb.NewDatabase(db_engine, dbLog, db_config)
	if err != nil {
		log.Error("create database handler failed - %s", err)
		return
	}

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
