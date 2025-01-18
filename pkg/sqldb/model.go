// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

// type TableMeta = struct {
// 	// custom options for special database uses
// 	Options Options

// 	// Table column definitions
// 	//	https://www.w3schools.com/sql/sql_create_table.asp
// 	// 	https://www.w3schools.com/sql/sql_datatypes.asp
// 	//
// 	// 	{column_name, datatype [, constraint]}
// 	//	constraint: is optional space seperated mix of {PRIMARY|UNIQUE|INDEX}.
// 	//	constraints can be added per column or globaly per table.
// 	//
// 	// Example:
// 	//
// 	//	{"col1", "VARCHAR(128) NOT NULL", "UNIQUE INDEX"}
// 	//	{"col2", "INTEGER", "INDEX"}
// 	//	{"col3", "TEXT"}
// 	//	{"col4", "BOOLEAN NOT NULL DEFAULT 0"}
// 	Columns [][]string

// 	// Table constraints definitions
// 	// 	https://www.w3schools.com/sql/sql_constraints.asp
// 	//
// 	// Example:
// 	//
// 	//	`PRIMARY KEY (col1)`
// 	//	`FOREIGN KEY (col1) REFERENCES table1 (col2) ON UPDATE CASCADE`
// 	//	`UNIQUE (col1,col2)`
// 	//	`CHECK (col1>=10 AND col2="val")`
// 	//
// 	// with Naming (for ALTER modifications):
// 	//
// 	//	`CONSTRAINT pk_name PRIMARY KEY (col1,col2)`
// 	//	`CONSTRAINT uc_name UNIQUE (col1)`
// 	//	`CONSTRAINT ck_name CHECK (col1 IN (0,1,2))`
// 	Constraints []string
// }

// type Model interface {
// 	TableName() string
// 	TableMeta() *TableMeta
// }
// type ModelAutoGuid interface{ AutoGuid() bool }
// type ModelDefaultOrders interface{ DefaultOrders() []string }
// type ModelDataReaders interface{ DataReaders() map[string]DataAdaptor }
// type ModelDataWriters interface{ DataWriters() map[string]DataAdaptor }
// type ModelAlterSchema interface{ AlterSchema(*Handler, string) error }
// type ModelInitialData interface{ InitialData(*Handler, string) error }

// func (dbh *Handler) CreateSchema(models map[string]Model) error {
// 	schema_sql := []string{}
// 	for tblname, model := range models {
// 		sql, err := dbh.engine.GenSchema(tblname, model)
// 		if err != nil {
// 			return err
// 		}
// 		schema_sql = append(schema_sql, sql...)
// 		dbh.Logger.Trace2("table '%s' schema:\n%s",
// 			tblname, strings.Join(sql, "\n"))
// 	}

// 	dbs := dbh.Session()
// 	defer dbs.Close()

// 	// create schema in session transaction
// 	if err := dbs.Begin(); err != nil {
// 		return err
// 	}
// 	// create tables schema
// 	for _, sql := range schema_sql {
// 		if _, err := dbs.Execute(sql); err != nil {
// 			return err
// 		}
// 	}
// 	// run tables schema upgrades
// 	for tblname, model := range models {
// 		if mdl, ok := model.(ModelUpgradeSchema); ok {
// 			if err := mdl.UpgradeSchema(dbs, tblname); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	if err := dbs.Commit(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (dbh *Handler) InitializeData(models map[string]Model) error {
// 	dbs := dbh.Session()
// 	defer dbs.Close()

// 	// initialze data in session transaction
// 	if err := dbs.Begin(); err != nil {
// 		return err
// 	}
// 	for tblname, model := range models {
// 		if mdl, ok := model.(ModelInitializeData); ok {
// 			if err := mdl.InitializeData(dbs, tblname); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	if err := dbs.Commit(); err != nil {
// 		return err
// 	}
// 	return nil
// }
