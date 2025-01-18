package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/logging"
)

var DB_OPTS = dictx.Dict{
	"database":   filepath.Join(os.TempDir(), "sample.db"),
	"extra_args": "",
}

//////////////////////////////// models

var Role *role

type role struct{}

func (*role) TableName() string { return "roles" }
func (*role) AutoGuid() bool    { return true }
func (*role) DefaultOrders() []string {
	return []string{"title ASC"}
}
func (*role) TableMeta() *sqldb.TableMeta {
	return &sqldb.TableMeta{
		Columns: [][]string{
			{"guid", "VARCHAR(32) NOT NULL", "PRIMARY"},
			{"title", "TEXT NOT NULL", "UNIQUE INDEX"},
			{"description", "TEXT"},
			{"access_level", "INTEGER"},
			{"builtin", "BOOLEAN DEFAULT 0"},
		},
		Constraints: []string{
			"CHECK (access_level>=0 AND access_level<=5)",
		},
	}
}
func (*role) InitialData(dbh *sqldb.Handler, _ string) error {
	// if dbs == nil {
	// 	return errors.New("invalid database session")
	// }

	// // check if default 'Administrator' role already exist
	// num, err := dbs.Query(Role).Filter("title=$?", "Administrator").Count()
	// if err != nil || num > 0 {
	// 	return err
	// }

	// // create default 'Administrator' role
	// _, err = dbs.Query(Role).Insert(sqldb.Data{
	// 	"title":        "Administrator",
	// 	"description":  "Administrator Full Access",
	// 	"access_level": 5,
	// 	"builtin":      true,
	// })
	// return err
	return nil
}

var User *user

type user struct{}

func (*user) TableName() string { return "users" }
func (*user) AutoGuid() bool    { return true }
func (*user) DefaultOrders() []string {
	return []string{"username ASC"}
}
func (*user) TableMeta() *sqldb.TableMeta {
	return &sqldb.TableMeta{
		Columns: [][]string{
			{"guid", "VARCHAR(32) NOT NULL", "PRIMARY"},
			{"username", "TEXT NOT NULL", "UNIQUE INDEX"},
			{"password", "TEXT NOT NULL"},
			{"enabled", "BOOLEAN DEFAULT 1"},
			{"role_guid", "VARCHAR(32) NOT NULL"},
		},
		Constraints: []string{
			"FOREIGN KEY (role_guid) REFERENCES roles (guid)" +
				" ON UPDATE CASCADE ON DELETE RESTRICT",
		},
	}
}

// func (*user) InitialData(dbs *sqldb.Session, _ string) error {
// 	// if dbs == nil {
// 	// 	return errors.New("invalid database session")
// 	// }

// 	// // check if default 'Admin' user already exist
// 	// num, err := dbs.Query(User).Filter("username=$?", "admin").Count()
// 	// if err != nil || num > 0 {
// 	// 	return err
// 	// }

// 	// // get default 'Administrator' role
// 	// role, err := dbs.Query(Role).Filter("title=$?", "Administrator").One()
// 	// if err != nil {
// 	// 	return err
// 	// } else if role == nil {
// 	// 	return errors.New("default 'Administrator' role not found")
// 	// }
// 	// role_guid := role.GetString("guid", "")
// 	// if len(role_guid) == 0 {
// 	// 	return errors.New("invalid empty 'Administrator' role guid")
// 	// }

// 	// // create default 'Admin' user
// 	// _, err = dbs.Query(User).Insert(sqldb.Data{
// 	// 	"username":  "admin",
// 	// 	"password":  "12345",
// 	// 	"enabled":   true,
// 	// 	"role_guid": role_guid,
// 	// })
// 	// return err
// 	return nil
// }

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

func run_initialization(dbh *sqldb.Handler) error {
	models := map[string]sqldb.Model{
		Role.TableName(): Role,
		User.TableName(): User,
	}
	if err := dbh.CreateSchema(models); err != nil {
		return err
	}
	// return dbh.InitialData(models)
	return nil
}

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
// 	// 		PrintModelData(items)
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
// 	// 		PrintModelData(items)
// 	// 	}

// 	// 	// filtered entries
// 	// 	fmt.Println("\nGet filter columns entries:")
// 	// 	if items, err := dbs.Query(&Foobar{}).
// 	// 		Filter("col2 LIKE $? OR col3 IN ($?,$?)", "description_3", 1, 3).
// 	// 		All(); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	} else {
// 	// 		PrintModelData(items)
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
// 	// 		PrintModelData(items)
// 	// 	}

// 	// fmt.Println("\nDelete: first row")
// 	// if _, err := dbs.Query(&Foobar{}).FilterBy("col3", 1).

// 	// 		Delete(); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	}

// 	// fmt.Println("-- After Delete --")

// 	// 	if items, err := dbs.Query(&Foobar{}).All(); err != nil {
// 	// 		fmt.Println("ERROR:", err.Error())
// 	// 		return
// 	// 	} else {

// 	// 		PrintModelData(items)
// 	// 	}

// 	return nil
// }

////////////////////////////////

func main() {
	debug := flag.Int("x", 0, "set debug modes, (default: 0)")
	setup := flag.Bool("setup", false, "perform database setup")
	flag.Parse()

	logger := logging.NewStdoutLogger("main")

	switch {
	case *debug >= 5:
		logger.Level = logging.TRACE3
	case *debug >= 3:
		logger.Level = logging.TRACE2
	case *debug >= 1:
		logger.Level = logging.DEBUG
	}

	// config
	database := dictx.Fetch(DB_OPTS, "database", "")
	if database == "" {
		fmt.Printf("Error: invalid database path\n")
		os.Exit(1)
	}
	fmt.Printf("\n* Using database: %v\n", database)
	fmt.Println("\nUsing Options:")
	for _, k := range []string{"database", "extra_args"} {
		// if DB_OPTS.IsExist(k) {
		// 	fmt.Printf(" - %-11v: %v\n", k, DB_OPTS[k])
		// }
	}

	// select engine and create db handler
	engine := sqldb.SqliteEngine(DB_OPTS)
	// dbh := sqldb.NewHandler(engine, DB_OPTS, logger)

	// database setup
	if *setup {
		fmt.Println("\n* Running Database Setup:")
		if err := sqldb.SqliteBackend().InteractiveSetup(DB_OPTS); err != nil {
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
