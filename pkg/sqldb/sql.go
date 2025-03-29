// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
)

// SQL Statment placeholder for variables.
const SQL_PLACEHOLDER = "?"

// Data represents the table data. where each row is represented
// into a map for columns as keys and data as values.
type Data = map[string]any

////////////////////////////////////////////////////

// StmtAttrs represents the SQL statment attributes.
type StmtAttrs struct {
	Tablename   string
	Columns     []string
	Filters     string
	FiltersArgs []any
	Groupby     []string
	Orderby     []string
	Having      string
	HavingArgs  []any
	Offset      int
	Limit       int
}

// ColumnMeta represents column definition.
//
// References:
//   - https://www.w3schools.com/sql/sql_create_table.asp
//   - https://www.w3schools.com/sql/sql_datatypes.asp
type ColumnMeta struct {
	// the column name, should be unique per table.
	Name string
	// the column data type as defined in SQL syntax.
	// ex. "VARCHAR(128) NOT NULL", "BOOLEAN NOT NULL DEFAULT false"
	Type string
	// set column primary key constraint.
	Primary bool
	// set column unique value constraint.
	Unique bool
	// set to create column index.
	Index bool
}

// ConstraintMeta represents constraint definitions.
//
// References:
//   - https://www.w3schools.com/sql/sql_constraints.asp
type ConstraintMeta struct {
	// the constraint name, should be unique per table.
	Name string
	// the constraint definition as defined in SQL syntax.
	// ex. "PRIMARY KEY (col1,col2)"
	//     "FOREIGN KEY (col1) REFERENCES table1 (col2) ON UPDATE CASCADE"
	//     "UNIQUE (col1,col2)"
	//     "CHECK (col1 IN (0,1,2))"
	//     "CHECK (col1>=10 AND col2="val")"
	Definition string
}

// TableMeta represents table definition, columns and constraints.
//
// References:
//   - https://www.w3schools.com/sql/sql_create_table.asp
//   - https://www.w3schools.com/sql/sql_datatypes.asp
//   - https://www.w3schools.com/sql/sql_constraints.asp
type TableMeta struct {
	// Table Columns meta
	Columns []ColumnMeta
	// Table Constraints as defined in SQL syntax. constraints are appended to
	// table after auto generated columns constraints.
	Constraints []ConstraintMeta
	// AutoGuid sets weather to enable AutoGuid operations, which is to
	// create a first primary guid column for table.
	// guid column is created with schema "guid VARCHAR(32) NOT NULL"
	AutoGuid bool
	// Extra options for backends.
	Args dictx.Dict
}

////////////////////////////////////////////////////

// SqlGenerator interface defines SQL statments generator.
type SqlGenerator interface {
	// FormatStmt prepares the statment placeholders format
	FormatStmt(stmt string) string

	// Select generates a SELECT statment
	Select(attrs *StmtAttrs) (string, []any)
	// Count generates a SELECT count(*) statment
	Count(attrs *StmtAttrs) (string, []any)
	// Insert generates an INSERT statment
	Insert(attrs *StmtAttrs, data Data) (string, []any)
	// Update generates an UPDATE statment
	Update(attrs *StmtAttrs, data Data) (string, []any)
	// Delete generates a DELETE statment
	Delete(attrs *StmtAttrs) (string, []any)

	// Schema generates table schema statments from metainfo
	Schema(tablename string, meta *TableMeta) []string
}

// StdSqlGenerator represents a standard SQL statment generator.
type StdSqlGenerator struct{}

// FormatStmt prepares the statment placeholders format
func (*StdSqlGenerator) FormatStmt(stmt string) string {
	return stmt
}

// Select generates a SELECT statment from attrs.
func (*StdSqlGenerator) Select(attrs *StmtAttrs) (string, []any) {
	// create the statment
	stmt := "SELECT "
	if len(attrs.Columns) > 0 {
		stmt += strings.Join(attrs.Columns, ", ")
	} else {
		stmt += "*"
	}
	stmt += " FROM " + attrs.Tablename

	if attrs.Filters != "" {
		stmt += " WHERE " + attrs.Filters
	}
	if len(attrs.Groupby) > 0 {
		stmt += " GROUP BY " + strings.Join(attrs.Groupby, ", ")
	}
	if attrs.Having != "" {
		stmt += " HAVING " + attrs.Having
	}
	if len(attrs.Orderby) > 0 {
		stmt += " ORDER BY " + strings.Join(attrs.Orderby, ", ")
	}
	if attrs.Offset > 0 {
		stmt += fmt.Sprintf(" OFFSET %d", attrs.Offset)
	}
	if attrs.Limit > 0 {
		stmt += fmt.Sprintf(" LIMIT %d", attrs.Limit)
	}
	stmt += ";"

	// create the params for statment placeholders
	params := append(attrs.FiltersArgs, attrs.HavingArgs...)

	return stmt, params
}

// Count generates a SELECT count(*) statment
func (*StdSqlGenerator) Count(attrs *StmtAttrs) (string, []any) {
	// create the statment
	stmt := "SELECT count(*) as count FROM " + attrs.Tablename

	if attrs.Filters != "" {
		stmt += " WHERE " + attrs.Filters
	}
	if len(attrs.Groupby) > 0 {
		stmt += " GROUP BY " + strings.Join(attrs.Groupby, ", ")
	}
	if attrs.Having != "" {
		stmt += " HAVING " + attrs.Having
	}
	stmt += ";"

	// create the params for statment placeholders
	params := append(attrs.FiltersArgs, attrs.HavingArgs...)

	return stmt, params
}

// Insert generates an INSERT statment
func (*StdSqlGenerator) Insert(attrs *StmtAttrs, data Data) (string, []any) {
	// create the statment
	columns, holders, params := []string{}, []string{}, []any{}
	for k, v := range data {
		columns = append(columns, k)
		holders = append(holders, SQL_PLACEHOLDER)
		params = append(params, v)
	}
	stmt := "INSERT INTO " + attrs.Tablename
	stmt += fmt.Sprintf(" (%v)", strings.Join(columns, ", "))
	stmt += fmt.Sprintf(" VALUES (%v)", strings.Join(holders, ", "))
	stmt += ";"

	return stmt, params
}

// Update generates an UPDATE statment
func (*StdSqlGenerator) Update(attrs *StmtAttrs, data Data) (string, []any) {
	// create the statment
	columns, params := []string{}, []any{}
	for k, v := range data {
		columns = append(columns, k+"="+SQL_PLACEHOLDER)
		params = append(params, v)
	}
	stmt := "UPDATE " + attrs.Tablename
	stmt += " SET " + strings.Join(columns, ", ")
	if attrs.Filters != "" {
		stmt += " WHERE " + attrs.Filters
	}

	// create the params for statment placeholders
	params = append(params, attrs.FiltersArgs...)

	return stmt, params
}

// Delete generates a DELETE statment
func (*StdSqlGenerator) Delete(attrs *StmtAttrs) (string, []any) {
	// create the statment
	stmt := "DELETE FROM " + attrs.Tablename
	if attrs.Filters != "" {
		stmt += " WHERE " + attrs.Filters
	}

	return stmt, attrs.FiltersArgs
}

// Schema generates table schema from table metainfo
func (*StdSqlGenerator) Schema(tablename string, meta *TableMeta) []string {
	var buff, constraints, indexes []string

	// if AutoGuid, add guid column if not exist as first column
	if meta.AutoGuid && meta.Columns[0].Name != "guid" {
		meta.Columns = append([]ColumnMeta{
			{Name: "guid", Type: "VARCHAR(32) NOT NULL", Primary: true},
		}, meta.Columns...)
	}

	table_exists := " IF NOT EXISTS"
	if dictx.Fetch(meta.Args, "disable_table_exists", false) {
		table_exists = ""
	}
	index_exists := " IF NOT EXISTS"
	if dictx.Fetch(meta.Args, "disable_index_exists", false) {
		index_exists = ""
	}

	// loop and parse columns meta
	for _, c := range meta.Columns {
		buff = append(buff, fmt.Sprintf("%s %s", c.Name, c.Type))

		// add constraints and indexes
		if c.Primary {
			constraints = append(constraints,
				fmt.Sprintf("PRIMARY KEY (%s)", c.Name))
		} else if c.Unique {
			constraints = append(constraints,
				fmt.Sprintf("UNIQUE (%s)", c.Name))
		}
		if c.Primary || c.Index {
			indexes = append(indexes, fmt.Sprintf(
				"CREATE INDEX%s ix_%s_%s ON %s (%s);",
				index_exists, tablename, c.Name, tablename, c.Name))
		}
	}

	// append column constraints
	buff = append(buff, constraints...)

	// add explicit table constraints
	for _, c := range meta.Constraints {
		if c.Name != "" {
			buff = append(buff, fmt.Sprintf(
				"CONSTRAINT %s %s", c.Name, c.Definition))
		} else {
			buff = append(buff, c.Definition)
		}
	}

	stmt := fmt.Sprintf(
		"CREATE TABLE%s %s (\n  %s\n);",
		table_exists, tablename, strings.Join(buff, ",\n  "))

	return append([]string{stmt}, indexes...)
}

////////////////////////////////////////////////////

// SqlIdent checks for a valid SQL identifier string.
func SqlIdent(s string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", s)
	return matched
}
