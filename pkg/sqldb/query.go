// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
)

type Query struct {
	dbs   *Session
	model Model

	// runtime table name to use, this allows for mapping
	// same model to multiple tables
	tablename string
	// query columns
	columns []string
	// query sql filters and args
	filters  []string
	exprargs []any
	// query grouping and ordering
	groupby []string
	orderby []string
	having  string
	// query data set limit and offset
	limit  int
	offset int

	// last runtime error
	lasterr error
}

func NewQuery(dbs *Session, model Model) *Query {
	return &Query{
		dbs:       dbs,
		model:     model,
		tablename: model.TableName(),
		columns:   model.Columns(),
		orderby:   model.Orders(),
	}
}

// TableName sets the runtime table name in statment instead of default
func (q *Query) TableName(name string) *Query {
	name = strings.TrimSpace(name)
	if name != "" {
		q.tablename = name
	}
	return q
}

// Columns sets the columns in statment
func (q *Query) Columns(columns ...string) *Query {
	if len(columns) > 0 {
		q.columns = columns
	}
	return q
}

// Filter adds filter expresion in statment
func (q *Query) Filter(expr string, params ...any) *Query {
	if expr != "" {
		q.filters = append(q.filters, expr)
		if len(params) > 0 {
			q.exprargs = append(q.exprargs, params...)
		}
	}
	return q
}

// GroupBy adds grouping in statment
func (q *Query) GroupBy(groupby ...string) *Query {
	q.groupby = groupby
	return q
}

// OrderBy add ordering "colname ASC|DESC" in statment
func (q *Query) OrderBy(orderby ...string) *Query {
	q.orderby = orderby
	return q
}

// Having adds having expr in statment
func (q *Query) Having(expr string, val any) *Query {
	if expr != "" {
		q.having = expr
		q.exprargs = append(q.exprargs, val)
	}
	return q
}

// Limit adds limit in statment
func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

// Offset add offset in statment
func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

// // add filters
// func (q *Query) FilterBy(column string, value any) *Query {
// 	if column != "" {
// 		cond := ""
// 		if len(q.filters) > 0 {
// 			cond = "AND "
// 		}
// 		q.filters = append(
// 			q.filters, cond+column+"="+SQL_PLACEHOLDER)
// 		q.exprargs = append(q.exprargs, value)
// 	}
// 	return q
// }

// All return all elements matching select query
func (q *Query) All() ([]Data, error) {
	if q.lasterr != nil {
		return nil, q.lasterr
	}

	limitPrefix := ""
	if q.dbs.db.Backend() == BACKEND_MSSQL {
		if q.limit > 0 && len(q.orderby) == 0 {
			limitPrefix = fmt.Sprintf("TOP(%v) ", q.limit)
		}
	}

	sql := "SELECT " + limitPrefix
	if len(q.columns) > 0 {
		sql += strings.Join(q.columns, ", ")
	} else {
		sql += "*"
	}
	sql += " FROM " + q.tablename

	if len(q.filters) > 0 {
		sql += "\nWHERE " + strings.Join(q.filters, " ")
	}
	if len(q.groupby) > 0 {
		sql += "\nGROUP BY " + strings.Join(q.groupby, ", ")
	}
	if len(q.having) > 0 {
		sql += "\nHAVING " + q.having
	}
	if len(q.orderby) > 0 {
		sql += "\nORDER BY " + strings.Join(q.orderby, ", ")
	}
	if q.dbs.db.Backend() == BACKEND_MSSQL {
		if q.offset > 0 || q.limit > 0 {
			sql += fmt.Sprintf(
				"\nOFFSET %v ROWS", q.offset)
		}
		if q.limit > 0 {
			sql += fmt.Sprintf(
				"\nFETCH NEXT %v ROWS ONLY", q.limit)
		}
	} else {
		if q.limit > 0 {
			sql += fmt.Sprintf("\nLIMIT %v", q.limit)
		}
		if q.offset > 0 {
			sql += fmt.Sprintf("\nOFFSET %v", q.offset)
		}
	}
	sql += ";"

	// run query and fetch data
	result, err := q.dbs.FetchAll(sql, q.exprargs...)
	if err != nil {
		return nil, err
	}

	// // apply DataReaders adapters
	// if mdl, ok := q.model.(ModelDataReaders); ok {
	// 	if err := FormatData(mdl.DataReaders(), result); err != nil {
	// 		return nil, err
	// 	}
	// }

	return result, nil
}

// First return first element matching select query
func (q *Query) First() (Data, error) {
	q.limit, q.offset = 1, 0
	result, err := q.All()
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}

// One return one element matching filter params or nil
// there must be only one element or none.
func (q *Query) One() (Data, error) {
	q.limit, q.offset = 2, 0
	result, err := q.All()
	if err != nil {
		return nil, err
	}
	switch len(result) {
	case 0:
		return nil, nil
	case 1:
		return result[0], nil
	}
	return nil, fmt.Errorf("%w - multiple entries found", ErrOperation)
}

// Gets model data defined by guid
func (q *Query) Get(guid string) (Data, error) {
	if !q.model.IsAutoGuid() {
		return nil, fmt.Errorf("%w - AutoGuid is disabled", ErrOperation)
	}
	q.filters = []string{"guid=?"}
	q.exprargs = []any{guid}
	q.groupby = []string{}
	q.orderby = []string{}
	q.having = ""
	return q.One()
}

// Counts the number of entries matching defined filters.
func (q *Query) Count() (int64, error) {
	if q.lasterr != nil {
		return 0, q.lasterr
	}

	sql := "SELECT count(*) as count FROM " + q.tablename
	if len(q.filters) > 0 {
		sql += "\nWHERE " + strings.Join(q.filters, " ")
	}
	if len(q.groupby) > 0 {
		sql += "\nGROUP BY " + strings.Join(q.groupby, ", ")
	}
	sql += ";"

	result, err := q.dbs.FetchAll(sql, q.exprargs...)
	if err != nil {
		return 0, err
	}
	if len(result) > 0 {
		if count, ok := result[0]["count"].(int64); ok {
			return count, nil
		}
	}

	return 0, fmt.Errorf("%w - invalid query result", ErrOperation)
}

// Insert data into table and return guid of new created entry if AutoGuid enabled.
func (q *Query) Insert(data Data) (string, error) {
	if q.lasterr != nil {
		return "", q.lasterr
	}

	// // apply DataWriters adapters
	// if mdl, ok := q.model.(ModelDataWriters); ok {
	// 	if err := FormatData(mdl.DataWriters(), []Data{data}); err != nil {
	// 		return "", err
	// 	}
	// }

	// check and create guid in data
	guid := ""
	if q.model.IsAutoGuid() {
		guid = dictx.Fetch(data, "guid", "")
		if guid == "" {
			guid = NewGuid()
			dictx.Set(data, "guid", guid)
		}
	}

	columns, holders, exprargs := []string{}, []string{}, []any{}
	for k, v := range data {
		columns = append(columns, k)
		holders = append(holders, "?")
		exprargs = append(exprargs, v)
	}

	sql := "INSERT INTO " + q.tablename
	sql += fmt.Sprintf("\n(%v)", strings.Join(columns, ", "))
	sql += fmt.Sprintf("\nVALUES (%v)", strings.Join(holders, ", "))
	sql += ";"

	_, err := q.dbs.Execute(sql, exprargs...)
	if err != nil {
		return "", err
	}

	return guid, nil
}

// update table data and return number of affected entries
func (q *Query) Update(data Data) (int64, error) {
	if q.lasterr != nil {
		return 0, q.lasterr
	}

	// check and remove guid from data
	if q.model.IsAutoGuid() {
		dictx.Delete(data, "guid")
	}

	// // apply DataWriters adapters
	// if mdl, ok := q.model.(ModelDataWriters); ok {
	// 	if err := FormatData(mdl.DataWriters(), []Data{data}); err != nil {
	// 		return 0, err
	// 	}
	// }

	columns, exprargs := []string{}, []any{}
	for k, v := range data {
		columns = append(columns, k+"=?")
		exprargs = append(exprargs, v)
	}

	sql := "UPDATE " + q.tablename
	sql += "\nSET " + strings.Join(columns, ", ")
	if len(q.filters) > 0 {
		sql += "\nWHERE " + strings.Join(q.filters, " ")
	}

	exprargs = append(exprargs, q.exprargs...)

	return q.dbs.Execute(sql, exprargs...)
}

// Deletes from table data and return number of affected entries
func (q *Query) Delete() (int64, error) {
	if q.lasterr != nil {
		return 0, q.lasterr
	}

	sql := "DELETE FROM " + q.tablename
	if len(q.filters) > 0 {
		sql += "\nWHERE " + strings.Join(q.filters, " ")
	}

	return q.dbs.Execute(sql, q.exprargs...)
}

// ////////////////////////////////////////// helpers

// // check valid SQL identifier string
// func (q *Query) Ident(name string) string {
// 	match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", name)
// 	if match {
// 		return name
// 	}
// 	q.lasterr = fmt.Errorf("invalid sql identifier '%v'", name)
// 	return ""
// }
