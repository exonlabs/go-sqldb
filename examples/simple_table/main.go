package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	"github.com/exonlabs/go-utils/pkg/xlog"
)

var (
	database = path.Join(os.TempDir(), "sample.db")

	logger   = xlog.NewStdoutLogger("main")
	dblogger = logger.ChildLogger("db")
)

//////////////////////////////// models

type Parent struct {
	*sqldb.BaseModel
}

func (dbm *Parent) TableName() string {
	return "parent"
}
func (dbm *Parent) TableSchema() *sqldb.TableSchema {
	return nil
}

type Child struct {
	*sqldb.BaseModel
}

func (dbm *Child) TableName() string {
	return "child"
}
func (dbm *Child) TableSchema() *sqldb.TableSchema {
	return nil
}

////////////////////////////////

func run_setup() error {
	return nil
}

func run_operations() error {
	// 	// define tables
	// 	tables := map[db.TableName]db.IModel{
	// 		"foobar": &Foobar{},
	// 	}

	// 	if err := dbh.InitDatabase(tables); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	}
	// 	fmt.Println("\nDB initialize: Done")

	// 	dbs := dbh.Session()
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
	// 		Filter("col2 LIKE $? OR col3 IN ($?,$?)", "description_3", 1, 3).
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

func main() {
	debug := flag.Int("x", 0, "set debug modes, (default: 0)")
	setup := flag.Bool("setup", false, "perform database setup")
	flag.Parse()

	switch {
	case *debug >= 5:
		logger.Level = xlog.TRACE4
		dblogger.Level = xlog.TRACE4
	case *debug >= 3:
		logger.Level = xlog.TRACE2
		dblogger.Level = xlog.TRACE2
	case *debug > 0:
		logger.Level = xlog.DEBUG
		dblogger.Level = xlog.DEBUG
	}

	fmt.Printf("\n* Using database: %v\n", database)

	if *setup {
		if err := run_setup(); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		} else {
			fmt.Println("Done")
		}
		return
	}

	if err := run_operations(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Println("Done")
	}
}
