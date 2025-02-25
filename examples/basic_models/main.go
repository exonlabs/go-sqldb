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
		// "operation_timeout": 5.0,
		// "retry_interval": 0.1,
	}
)

//////////////////////////////// models

type job struct{ sqldb.BaseModel }

var Job *job = &job{sqldb.BaseModel{
	DefaultTable:  "jobs",
	DefaultOrders: []string{"title ASC"},
	AutoGuid:      true,
}}

type jobMeta struct{ sqldb.BaseModelMeta }

var JobMeta *jobMeta = &jobMeta{sqldb.BaseModelMeta{
	Columns: []sqldb.ColumnMeta{
		{Name: "title", Type: "VARCHAR(128) NOT NULL",
			Unique: true, Index: true},
		{Name: "description", Type: "TEXT"},
		{Name: "access_level", Type: "INTEGER"},
		{Name: "high_management", Type: "BOOLEAN DEFAULT false"},
	},
	Constraints: []sqldb.ConstraintMeta{
		{Definition: "CHECK (access_level>=0 AND access_level<=5)"},
	},
	AutoGuid: true,
	Args: dictx.Dict{
		"sqlite_without_rowid": true,
	},
}}

func (*jobMeta) InitialData(dbs *sqldb.Session, _ string) error {
	jobs := []sqldb.Data{{
		"title":           "Default_Employee",
		"description":     "Default Employee Position",
		"access_level":    1,
		"high_management": false,
	}, {
		"title":           "General_Manager",
		"description":     "General Manager Position",
		"access_level":    5,
		"high_management": true,
	}}

	for _, data := range jobs {
		// check if already exists
		job, err := dbs.Query(Job).FilterBy("title", data["title"]).One()
		if err != nil {
			return err
		} else if job != nil { // already exists
			continue
		}

		// create new job
		if _, err = dbs.Query(Job).Insert(data); err != nil {
			return err
		}
	}

	return nil
}

type employee struct{ sqldb.BaseModel }

var Employee *employee = &employee{sqldb.BaseModel{
	DefaultTable:  "employees",
	DefaultOrders: []string{"fullname ASC"},
	AutoGuid:      true,
}}

type employeeMeta struct{ sqldb.BaseModelMeta }

var EmployeeMeta *employeeMeta = &employeeMeta{sqldb.BaseModelMeta{
	Columns: []sqldb.ColumnMeta{
		{Name: "fullname", Type: "VARCHAR(128) NOT NULL",
			Unique: true, Index: true},
		{Name: "email", Type: "VARCHAR(256)"},
		{Name: "active", Type: "BOOLEAN DEFAULT true"},
		{Name: "job_guid", Type: "VARCHAR(32) NOT NULL"},
	},
	Constraints: []sqldb.ConstraintMeta{
		{Definition: "FOREIGN KEY (job_guid) REFERENCES jobs (guid) " +
			"ON UPDATE CASCADE ON DELETE RESTRICT"},
	},
	AutoGuid: true,
	Args: dictx.Dict{
		"sqlite_without_rowid": true,
	},
}}

func (*employeeMeta) InitialData(dbs *sqldb.Session, _ string) error {
	jobs_guids := map[string]string{}
	if jobs, err := dbs.Query(Job).
		Columns("guid", "title").All(); err != nil {
		return err
	} else {
		for _, j := range jobs {
			jobs_guids[j["title"].(string)] = j["guid"].(string)
		}
	}

	employees := []sqldb.Data{{
		"fullname":  "Employee 001",
		"email":     "employee.001@company.com",
		"active":    true,
		"job_title": "General_Manager",
	}, {
		"fullname":  "Employee 002",
		"email":     "employee.002@company.com",
		"active":    true,
		"job_title": "Default_Employee",
	}, {
		"fullname":  "Employee 003",
		"email":     "",
		"active":    false,
		"job_title": "Default_Employee",
	}}

	for _, data := range employees {
		// check if already exists
		empl, err := dbs.Query(Employee).
			FilterBy("fullname", data["fullname"]).One()
		if err != nil {
			return err
		} else if empl != nil { // already exists
			continue
		}

		// check job exists
		if job_guid, ok := jobs_guids[data["job_title"].(string)]; ok {
			data["job_guid"] = job_guid
			delete(data, "job_title")
		} else {
			return fmt.Errorf("job not found: %v", data["job_title"])
		}

		// create new employee
		if _, err = dbs.Query(Employee).Insert(data); err != nil {
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
	metainfo := []sqldb.TableModelMeta{
		{TableName: Job.DefaultTable, ModelMeta: JobMeta},
		{TableName: Employee.DefaultTable, ModelMeta: EmployeeMeta},
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
