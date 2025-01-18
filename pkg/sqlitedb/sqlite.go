// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	// 	"github.com/exonlabs/go-utils/pkg/abc/dictx"
)

type Engine struct {
	sqlDB *sql.DB

	Database string
	// ExtraArgs string
}

// func SqliteDB(opts Options) (*sqlite_engine, error) {
// 	// // params
// 	// database, _ := options["database"].(string)
// 	// if len(database) == 0 {
// 	// 	return nil, fmt.Errorf("invalid database configuration")
// 	// }
// 	// extargs, _ := options["extargs"].(string)
// 	// if !strings.Contains(extargs, "_foreign_keys=") {
// 	// 	extargs = "_foreign_keys=1&" + extargs
// 	// }

// 	// // create data source name
// 	// dsn := fmt.Sprintf("%v?%v", database, extargs)

// 	// sqlDB, err := sql.Open("sqlite3", dsn)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// return sqlDB, nil

// 	return nil, nil
// }

// // return backend name
// func (*sqlite_engine) Backend() string { return SQLITE_BACKEND }

// // format args placeholders in sql statment
// func (*sqlite_engine) FormatSql(sql string) string {
// 	return strings.Replace(sql, SQL_PLACEHOLDER, "?", -1)
// }

// // check if retry operation is practical for certain error type
// func (*sqlite_engine) CanRetryErr(err error) bool {
// 	return false
// }

// /////////////////////////////////////////////////////////

// type sqlite_backend struct{}

// func SqliteBackend() *sqlite_backend {
// 	return &sqlite_backend{}
// }

// func (*sqlite_engine) CreateSchema(
// 	tblname string, model Model) ([]string, error) {

// 	if tblname == "" {
// 		tblname = model.TableName()
// 	}

// 	meta := model.TableMeta()
// 	auto_guid := dictx.Fetch(meta.Options, "sqlite_without_rowid", false)

// 	columns := meta.Columns
// 	// add guid column if not exist as first column
// 	if _, ok := model.(ModelAutoGuid); ok {
// 		if columns[0][0] != "guid" {
// 			columns = append([][]string{
// 				{"guid", "VARCHAR(32) NOT NULL", "PRIMARY"},
// 			}, columns...)
// 		}
// 		auto_guid = true
// 	}

// 	var expr, constraints, indexes []string

// 	for _, c := range columns {
// 		expr = append(expr, c[0]+" "+c[1])

// 		// add check constraint for bool datatype
// 		if strings.Contains(c[1], "BOOLEAN") {
// 			constraints = append(constraints,
// 				fmt.Sprintf("CHECK (%v IN (0,1))", c[0]))
// 		}

// 		// no column constraint
// 		if len(c) < 3 {
// 			continue
// 		}

// 		if strings.Contains(c[2], "PRIMARY") {
// 			// add primary_key constraint
// 			constraints = append(constraints,
// 				fmt.Sprintf("PRIMARY KEY (%v)", c[0]))
// 		} else if strings.Contains(c[2], "UNIQUE") &&
// 			!strings.Contains(c[2], "INDEX") {
// 			// add unique constraint if not indexed column
// 			constraints = append(constraints,
// 				fmt.Sprintf("UNIQUE (%v)", c[0]))
// 		}

// 		if strings.Contains(c[2], "PRIMARY") ||
// 			strings.Contains(c[2], "INDEX") {
// 			u := ""
// 			if strings.Contains(c[2], "PRIMARY") ||
// 				strings.Contains(c[2], "UNIQUE") {
// 				u = "UNIQUE "
// 			}
// 			indexes = append(indexes, fmt.Sprintf(
// 				"CREATE %vINDEX IF NOT EXISTS ix_%v_%v ON %v (%v);",
// 				u, tblname, c[0], tblname, c[0]))
// 		}
// 	}

// 	// add column constraints
// 	expr = append(expr, constraints...)
// 	// add explicit table constraints
// 	expr = append(expr, meta.Constraints...)

// 	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tblname)
// 	sql += "\n   " + strings.Join(expr, ",\n   ")
// 	if auto_guid {
// 		sql += "\n) WITHOUT ROWID;"
// 	} else {
// 		sql += "\n);"
// 	}

// 	result := []string{sql}
// 	result = append(result, indexes...)
// 	return result, nil
// }

// // interactive database configuration
// func (*sqlite_backend) InteractiveConfig(opts Options) (Options, error) {
// 	// con := xterm.NewConsole()

// 	// if v, err := con.Required().ReadValue(
// 	// 	"Enter database path",
// 	// 	dictx.Fetch(opts, "database", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	dictx.Set(opts, "database", v)
// 	// }

// 	// if v, err := con.ReadValue(
// 	// 	"Enter connection extra args",
// 	// 	dictx.Fetch(opts, "extra_args", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	dictx.Set(opts, "extra_args", v)
// 	// }

// 	return opts, nil
// }

// // interactive database setup
// func (*sqlite_backend) InteractiveSetup(opts Options) error {
// 	// database := dictx.Fetch(opts, "database", "")
// 	// if database == "" {
// 	// 	return fmt.Errorf("%w - invalid database path", ErrOperation)
// 	// }

// 	// if _, err := os.Stat(database); os.IsNotExist(err) {
// 	// 	syscall.Umask(0)
// 	// 	f, err := os.OpenFile(
// 	// 		database, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o664)
// 	// 	if err != nil {
// 	// 		return fmt.Errorf("%w - %s", ErrOperation, err.Error())
// 	// 	}
// 	// 	defer f.Close()
// 	// }
// 	return nil
// }
