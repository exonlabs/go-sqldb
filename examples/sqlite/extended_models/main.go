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

//////////////////////////////// models

type dataset struct{ sqldb.BaseModel }

var DataSet = &dataset{sqldb.BaseModel{
	DefaultTable: "datasets",
	AutoGuid:     false,
}}

func (m *dataset) TableMeta() *sqldb.TableMeta {
	return &sqldb.TableMeta{
		Columns: []sqldb.ColumnMeta{
			{Name: "data", Type: "TEXT"},
			{Name: "csv", Type: "TEXT"}, // 2-way mapping
			{Name: "hex", Type: "TEXT"}, // 1-way mapping
		},
		AutoGuid: false,
		Args: dictx.Dict{
			"sqlite_without_rowid": false,
		},
	}
}

func (m *dataset) DataEncode(data []sqldb.Data) error {
	for i := 0; i < len(data); i++ {
		if _, ok := data[i]["csv"]; ok {
			if v, ok := data[i]["csv"].([]string); ok {
				data[i]["csv"] = strings.Join(v, ",")
			} else {
				data[i]["csv"] = ""
			}
		}
		if _, ok := data[i]["hex"]; ok {
			if v, ok := data[i]["hex"].([]byte); ok {
				data[i]["hex"] = fmt.Sprintf("%X", v)
			} else {
				data[i]["hex"] = ""
			}
		}
	}
	return nil
}

func (m *dataset) DataDecode(data []sqldb.Data) error {
	for i := 0; i < len(data); i++ {
		if _, ok := data[i]["csv"]; ok {
			if v, ok := data[i]["csv"].(string); ok {
				data[i]["csv"] = strings.Split(v, ",")
			} else {
				data[i]["csv"] = ""
			}
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
		{Table: DataSet.DefaultTable, Model: DataSet},
	}
	return sqldb.InitializeModels(db, metainfo)
}

func run_operations(db *sqldb.Database) error {
	if !db.Ping() {
		return errors.New("database connection down")
	}

	// create new session
	dbs := db.Session()

	// add new data
	data := sqldb.Data{
		"data": "normal data",
		"csv":  []string{"1", "2", "3", "4"},
		"hex":  []byte("hex data 1 2 3 4"),
	}
	fmt.Printf("\n* Adding new data:\n  - %v\n", data)
	_, err := dbs.Query(DataSet).Insert(data)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}

	// listing all data
	fmt.Println("\n* List all data (high level from model)")
	if data, err := dbs.Query(DataSet).All(); err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		for _, v := range data {
			fmt.Println("  - " + print_data(v))
		}
		fmt.Printf("Total: %d\n", len(data))
	}

	// listing all data
	fmt.Println("\n* List all data (low level by session)")
	if data, err := dbs.Fetch("SELECT * FROM datasets;"); err != nil {
		fmt.Println("ERROR:", err.Error())
	} else {
		for _, v := range data {
			fmt.Println("  - " + dictx.String(v))
		}
		fmt.Printf("Total: %d\n", len(data))
	}

	// delete all data
	fmt.Println("\n* Delete all data")
	if _, err := dbs.Query(DataSet).Delete(); err != nil {
		fmt.Println("ERROR:", err.Error())
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
