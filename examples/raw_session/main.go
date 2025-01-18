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

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/logging"
)

var (
	DB_PATH   = filepath.Join(os.TempDir(), "sample.db")
	DB_CONFIG = dictx.Dict{
		"database": DB_PATH,
	}
)

func run_operations(dbh *sqldb.Handler) error {

	return nil
}

func backend_setup(dbh *sqldb.Handler) error {

	return nil
}

func backend_clean(dbh *sqldb.Handler) error {
	return os.Remove(DB_PATH)
}

func main() {
	log := logging.NewStdoutLogger("main")

	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			indx := bytes.Index(stack, []byte("panic({"))
			log.Panic("%s", r)
			log.Trace1("\n----------\n%s----------", stack[indx:])
			os.Exit(1)
		} else {
			log.Info("exit")
			os.Exit(0)
		}
	}()

	debug0 := flag.Bool("x", false, "\nenable debug logs")
	debug1 := flag.Bool("xx", false, "enable debug and trace1 logs")
	debug2 := flag.Bool("xxx", false, "enable debug and trace2 logs")
	debug3 := flag.Bool("xxxx", false, "enable debug and trace3 logs")
	setup := flag.Bool("setup", false, "perform database setup")
	flag.Parse()

	switch {
	case *debug3:
		log.Level = logging.TRACE3
	case *debug2:
		log.Level = logging.TRACE2
	case *debug1:
		log.Level = logging.TRACE1
	case *debug0:
		log.Level = logging.DEBUG
	}

	log.Info("**** starting ****")

	fmt.Printf("\n* Using database: %s\n", DB_PATH)
	// fmt.Println("\nUsing Options:")
	// for _, k := range []string{"database", "extra_args"} {
	// 	// if DB_CONFIG.IsExist(k) {
	// 	// 	fmt.Printf(" - %-11v: %v\n", k, DB_CONFIG[k])
	// 	// }
	// }

	// select engine and create db handler
	engine := sqldb.SqliteEngine(DB_CONFIG)
	// dbh := sqldb.NewHandler(engine, DB_CONFIG, logger)

	// database setup
	if *setup {
		fmt.Println("\n* Running Database Setup:")
		if err := sqldb.SqliteBackend().InteractiveSetup(DB_CONFIG); err != nil {
			if strings.Contains(err.Error(), "EOF") {
				fmt.Print("\n--exit--\n\n")
			} else {
				fmt.Printf("Error: %s\n", err.Error())
			}
			os.Exit(1)
		}
		if err := run_initialization(dbh); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Done\n\n")
		os.Exit(0)
	}

	// run database operations
	// if err := run_operations(dbh); err != nil {
	// 	fmt.Printf("Error: %s\n", err.Error())
	// 	os.Exit(1)
	// }
	fmt.Printf("\n* Done\n\n")
}
