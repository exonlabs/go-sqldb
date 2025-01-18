package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	"github.com/exonlabs/go-sqldb/pkg/sqldbutils"
	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/logging"
)

//////////////////////////////// models

// type Foobar struct{ *sqldb.BaseModel }

// func (*Foobar) TableName() string { return "foobar" }
// func (*Foobar) TableMeta() *sqldb.TableMeta {
// 	return &sqldb.TableMeta{
// 		Columns: [][]string{
// 			{"guid", "TEXT NOT NULL", "PRIMARY"},
// 			{"col1", "VARCHAR(128) NOT NULL", "UNIQUE INDEX"},
// 			{"col2", "TEXT"},
// 			{"col3", "INTEGER"},
// 			{"col4", "BOOLEAN NOT NULL DEFAULT 0"},
// 		},
// 	}
// }
// func (*Foobar) DefaultOrders() []string {
// 	return []string{"col1 ASC"}
// }
// func (dbm *Foobar) InitializeData(dbs *sqldb.Session, tblname string) error {
// 	var err error
// 	for i := 0; i < 5; i++ {
// 		var num int64
// 		num, err = dbs.Query(dbm).Table(tblname).
// 			Filter("col1=$?", "foo_"+strconv.Itoa(i)).Count()
// 		if num == 0 {
// 			_, err = dbs.Query(dbm).Table(tblname).Insert(sqldb.Data{
// 				"col1": "foo_" + strconv.Itoa(i),
// 				"col2": "description_" + strconv.Itoa(i),
// 				"col3": i,
// 			})
// 		}
// 	}
// 	return err
// }

//////////////////////////////// operations

// func print_data(data []sqldb.Data) {
// 	if len(data) > 0 {
// 		keys := data[0].Keys()
// 		for _, item := range data {
// 			for _, k := range keys {
// 				log.Info("%v: %v\n", k, item[k])
// 			}
// 		}
// 	}
// }

// func run_operations(dbh *sqldb.Handler) error {
// 	// 	// define tables
// 	// 	tables := map[db.TableName]db.IModel{
// 	// 		"foobar": &Foobar{},
// 	// 	}

// 	// 	if err := dbh.InitDatabase(tables); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	}
// 	// 	fmt.Println("\nDB initialize: Done")

// 	// 	dbs := dbh.Session()
// 	// 	defer dbs.Close()

// 	// 	// checking DB
// 	// 	fmt.Println("\nAll entries:")
// 	// 	if items, err := dbs.Query(&Foobar{}).All(); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	} else {
// 	// 		print_data(items)
// 	// 	}
// 	// 	if total, err := dbs.Query(&Foobar{}).Count(); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	} else {
// 	// 		fmt.Println("\nTotal Items:", total)
// 	// 	}

// 	// 	// custom columns
// 	// 	fmt.Println("\nGet custom columns entries:")
// 	// 	if items, err := dbs.Query(&Foobar{}).Columns("col1", "col2").
// 	// 		Limit(2).OrderBy("col1 DESC").All(); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	} else {
// 	// 		print_data(items)
// 	// 	}

// 	// 	// filtered entries
// 	// 	fmt.Println("\nGet filter columns entries:")
// 	// 	if items, err := dbs.Query(&Foobar{}).
// 	// 		Filter("col2 LIKE $? OR col3 IN ($?,$?)", "description_3", 1, 3).
// 	// 		All(); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	} else {
// 	// 		print_data(items)
// 	// 	}

// 	// 	// update entries
// 	// 	fmt.Println("\nModify: first row")
// 	// 	if _, err := dbs.Query(&Foobar{}).FilterBy("col3", 1).
// 	// 		Update(db.ModelData{
// 	// 			"col1": "boo_1", "col2": "boo_2", "col4": 1,
// 	// 		}); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	}
// 	// 	fmt.Println("-- After Modify --")
// 	// 	if items, err := dbs.Query(&Foobar{}).All(); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	} else {
// 	// 		print_data(items)
// 	// 	}

// 	// 	fmt.Println("\nDelete: first row")
// 	// 	if _, err := dbs.Query(&Foobar{}).FilterBy("col3", 1).
// 	// 		Delete(); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	}
// 	// 	fmt.Println("-- After Delete --")
// 	// 	if items, err := dbs.Query(&Foobar{}).All(); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	} else {
// 	// 		print_data(items)
// 	// 	}

// 	return nil
// }

func run_operations(dbh *sqldb.Handler) error {

	return nil
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

	backends := []string{
		sqldb.BACKEND(sqldb.BACKEND_SQLITE),
		sqldb.BACKEND(sqldb.BACKEND_MYSQL),
		sqldb.BACKEND(sqldb.BACKEND_PGSQL),
		sqldb.BACKEND(sqldb.BACKEND_MSSQL),
	}

	debug0 := flag.Bool("x", false, "\nenable debug logs")
	debug1 := flag.Bool("xx", false, "enable debug and trace1 logs")
	debug2 := flag.Bool("xxx", false, "enable debug and trace2 logs")
	debug3 := flag.Bool("xxxx", false, "enable debug and trace3 logs")
	backend := flag.String("backend", "",
		fmt.Sprintf("select backend {%s}", strings.Join(backends, "|")))
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

	// selecting backend
	if *backend == "" || !slices.Contains(backends, *backend) {
		fmt.Printf("Error: invalid backend '%s'\n\n", *backend)
		return
	}

	database := "test"
	if *backend == sqldb.BACKEND(sqldb.BACKEND_SQLITE) {
		database = filepath.Join(os.TempDir(), database+".db")
	}

	log.Info("**** starting ****")

	log.Info("Using backend: %s", *backend)
	fmt.Println()

	// setting backend config
	fmt.Println("* Configure database:")
	cfg, err := sqldbutils.InteractiveConfig(
		*backend, dictx.Dict{"database": database})
	if err != nil {
		if !strings.Contains(err.Error(), "EOF") {
			fmt.Printf("Error: %s\n", err.Error())
		}
		fmt.Println()
		return
	}
	fmt.Println()

	log.Info("Database Options:")
	for _, k := range []string{
		"database", "host", "port", "username", "password"} {
		if dictx.IsExist(cfg, k) {
			log.Info(" - %-8v :  %v", k, dictx.Get(cfg, k, nil))
		}
	}
	fmt.Println()

	// setup database
	if *setup {
		fmt.Println("* Setup database:")
		err := sqldbutils.InteractiveSetup(*backend, cfg)
		if err != nil {
			if !strings.Contains(err.Error(), "EOF") {
				fmt.Printf("Error: %s\n", err.Error())
			}
		}
		fmt.Println()
		return
	}

	// select engine and create db handler
	engine, err := sqldbutils.CreateEngine(*backend, cfg)
	if err != nil {
		log.Error("create engine failed - %s", err.Error())
		return
	}
	dbh := sqldb.NewHandler(engine, log, cfg)

	if err := run_operations(dbh); err != nil {
		log.Info("Error: %s\n", err.Error())
		return
	}

}
