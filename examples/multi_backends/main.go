package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	"github.com/exonlabs/go-utils/pkg/unix/xterm"
	"github.com/exonlabs/go-utils/pkg/xlog"
)

//////////////////////////////// models

type Foobar struct{ *sqldb.BaseModel }

func (*Foobar) TableName() string { return "foobar" }
func (*Foobar) TableMeta() *sqldb.TableMeta {
	return &sqldb.TableMeta{
		Columns: [][]string{
			{"guid", "TEXT NOT NULL", "PRIMARY"},
			{"col1", "VARCHAR(128) NOT NULL", "UNIQUE INDEX"},
			{"col2", "TEXT"},
			{"col3", "INTEGER"},
			{"col4", "BOOLEAN NOT NULL DEFAULT 0"},
		},
	}
}
func (*Foobar) DefaultOrders() []string {
	return []string{"col1 ASC"}
}
func (dbm *Foobar) InitializeData(dbs *sqldb.Session, tblname string) error {
	var err error
	for i := 0; i < 5; i++ {
		var num int64
		num, err = dbs.Query(dbm).Table(tblname).
			Filter("col1=$?", "foo_"+strconv.Itoa(i)).Count()
		if num == 0 {
			_, err = dbs.Query(dbm).Table(tblname).Insert(sqldb.Data{
				"col1": "foo_" + strconv.Itoa(i),
				"col2": "description_" + strconv.Itoa(i),
				"col3": i,
			})
		}
	}
	return err
}

//////////////////////////////// operations

func print_data(data []sqldb.Data) {
	if len(data) > 0 {
		keys := data[0].Keys()
		for _, item := range data {
			for _, k := range keys {
				fmt.Printf("%v: %v\n", k, item[k])
			}
		}
	}
}

func run_operations(dbh *sqldb.Handler) error {
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
	// 		print_data(items)
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
	// 		print_data(items)
	// 	}

	// 	// filtered entries
	// 	fmt.Println("\nGet filter columns entries:")
	// 	if items, err := dbs.Query(&Foobar{}).
	// 		Filter("col2 LIKE $? OR col3 IN ($?,$?)", "description_3", 1, 3).
	// 		All(); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	} else {
	// 		print_data(items)
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
	// 		print_data(items)
	// 	}

	// 	fmt.Println("\nDelete: first row")
	// 	if _, err := dbs.Query(&Foobar{}).FilterBy("col3", 1).
	// 		Delete(); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	}
	// 	fmt.Println("-- After Delete --")
	// 	if items, err := dbs.Query(&Foobar{}).All(); err != nil {
	// 		fmt.Println("ERROR:", err.Error())
	// 		return
	// 	} else {
	// 		print_data(items)
	// 	}

	return nil
}

////////////////////////////////

func main() {
	debug := flag.Int("x", 0, "set debug modes, (default: 0)")
	backend := flag.String("backend", "",
		fmt.Sprintf("select backend {%s}", strings.Join(sqldb.BACKENDS, "|")))
	setup := flag.Bool("setup", false, "perform database setup")
	flag.Parse()

	logger := xlog.NewStdoutLogger("main")

	switch {
	case *debug >= 5:
		logger.Level = xlog.TRACE4
	case *debug >= 3:
		logger.Level = xlog.TRACE2
	case *debug >= 1:
		logger.Level = xlog.DEBUG
	}

	// selecting backend
	if *backend == "" {
		fmt.Println()
		v, err := xterm.NewConsole().Required().
			SelectValue("Select Backend", sqldb.BACKENDS, nil)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
		*backend = v
	} else {
		if !slices.Contains(sqldb.BACKENDS, *backend) {
			fmt.Printf("Error: invalid backend '%s'\n", *backend)
			os.Exit(1)
		}
		fmt.Printf("\n* Using backend: %v\n", *backend)
	}

	// setting config
	opts, err := interactive_config(*backend)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			fmt.Print("\n--exit--\n\n")
		} else {
			fmt.Printf("Error: %s\n", err.Error())
		}
		os.Exit(1)
	}
	fmt.Println("\nUsing Options:")
	for _, k := range []string{"database", "host", "port",
		"username", "password", "extra_args"} {
		if opts.IsExist(k) {
			fmt.Printf(" - %-11v: %v\n", k, opts[k])
		}
	}

	// select engine and create db handler
	engine, err := select_engine(*backend, opts)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	dbh := sqldb.NewHandler(engine, opts, logger)

	// database setup
	if *setup {
		fmt.Println("\n* Running Database Setup:")
		if err := interactive_setup(*backend, opts); err != nil {
			if strings.Contains(err.Error(), "EOF") {
				fmt.Print("\n--exit--\n\n")
			} else {
				fmt.Printf("Error: %s\n", err.Error())
			}
			os.Exit(1)
		}

		fmt.Printf("Done\n\n")
		os.Exit(0)
	}

	if err := run_operations(dbh); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("\n* Done\n\n")
}

func interactive_config(backend string) (sqldb.Options, error) {
	switch backend {
	case "sqlite":
		return sqldb.SqliteBackend().InteractiveConfig(sqldb.Options{
			"database":   filepath.Join(os.TempDir(), "test.db"),
			"extra_args": "",
		})
	case "mysql":
		return sqldb.MysqlBackend().InteractiveConfig(sqldb.Options{
			"database":   "test",
			"host":       "localhost",
			"port":       3306,
			"username":   "admin",
			"password":   "admin",
			"extra_args": "",
		})
	case "pgsql":
		return sqldb.PgsqlBackend().InteractiveConfig(sqldb.Options{
			"database":   "test",
			"host":       "localhost",
			"port":       5432,
			"username":   "postgres",
			"password":   "",
			"extra_args": "",
		})
	case "mssql":
		return sqldb.MssqlBackend().InteractiveConfig(sqldb.Options{
			"database":   "test",
			"host":       "localhost",
			"port":       1433,
			"username":   "sa",
			"password":   "root@Root",
			"extra_args": "",
		})
	}
	return nil, fmt.Errorf("invalid backend '%s'", backend)
}

func interactive_setup(backend string, opts sqldb.Options) error {
	switch backend {
	case "sqlite":
		return sqldb.SqliteBackend().InteractiveSetup(opts)
	case "mysql":
		return sqldb.MysqlBackend().InteractiveSetup(opts)
	case "pgsql":
		return sqldb.PgsqlBackend().InteractiveSetup(opts)
	case "mssql":
		return sqldb.MssqlBackend().InteractiveSetup(opts)
	}
	return fmt.Errorf("invalid backend '%s'", backend)
}

func select_engine(backend string, opts sqldb.Options) (sqldb.Engine, error) {
	switch backend {
	case "sqlite":
		return sqldb.SqliteEngine(opts), nil
	case "mysql":
		return sqldb.MysqlEngine(opts), nil
	case "pgsql":
		return sqldb.PgsqlEngine(opts), nil
	case "mssql":
		return sqldb.MssqlEngine(opts), nil
	}
	return nil, fmt.Errorf("invalid backend '%s'", backend)
}
