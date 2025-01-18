// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

// type Query struct {
// 	handler *Handler
// 	model   Model

// 	// runtime table name to use, this allows for mapping
// 	// same model to multiple tables
// 	tablename string
// 	// query columns
// 	columns []string
// 	// query sql filters and args
// 	filters  []string
// 	execargs []any
// 	// query grouping and ordering
// 	groupby []string
// 	orderby []string
// 	having  string
// 	// query data set limit and offset
// 	limit  int
// 	offset int

// 	// last runtime error
// 	lasterr error
// }

// func newQuery(dbh *Handler, model Model) *Query {
// 	q := &Query{
// 		handler:   dbh,
// 		model:     model,
// 		tablename: model.TableName(),
// 	}
// 	if mdl, ok := model.(ModelDefaultOrders); ok {
// 		q.orderby = mdl.DefaultOrders()
// 	}
// 	return q
// }

// ////////////////////////////////////////// Creation

// // set runtime table name instead of default
// func (q *Query) Table(tblname string) *Query {
// 	if tblname != "" {
// 		q.tablename = tblname
// 	}
// 	return q
// }

// // set columns
// func (q *Query) Columns(columns ...string) *Query {
// 	q.columns = columns
// 	return q
// }

// // add filters
// func (q *Query) Filter(expr string, params ...any) *Query {
// 	if expr != "" {
// 		q.filters = append(q.filters, expr)
// 		q.execargs = append(q.execargs, params...)
// 	}
// 	return q
// }

// // add filters
// func (q *Query) FilterBy(column string, value any) *Query {
// 	if column != "" {
// 		cond := ""
// 		if len(q.filters) > 0 {
// 			cond = "AND "
// 		}
// 		q.filters = append(
// 			q.filters, cond+column+"="+SQL_PLACEHOLDER)
// 		q.execargs = append(q.execargs, value)
// 	}
// 	return q
// }

// // add grouping
// func (q *Query) GroupBy(groupby ...string) *Query {
// 	q.groupby = groupby
// 	return q
// }

// // add ordering: "colname ASC|DESC"
// func (q *Query) OrderBy(orderby ...string) *Query {
// 	q.orderby = orderby
// 	return q
// }

// // add having expr
// func (q *Query) Having(expr string, val any) *Query {
// 	if expr != "" {
// 		q.having = expr
// 		q.execargs = append(q.execargs, val)
// 	}
// 	return q
// }

// // add limit
// func (q *Query) Limit(limit int) *Query {
// 	q.limit = limit
// 	return q
// }

// // add offset
// func (q *Query) Offset(offset int) *Query {
// 	q.offset = offset
// 	return q
// }

// ////////////////////////////////////////// Operations

// // return all elements matching select query
// func (q *Query) All() ([]Data, error) {
// 	if q.lasterr != nil {
// 		return nil, q.lasterr
// 	}

// 	backend := q.handler.engine.BackendName()
// 	limitPrefix := ""
// 	if backend == MSSQL_BACKEND {
// 		if q.limit > 0 && len(q.orderby) == 0 {
// 			limitPrefix = fmt.Sprintf("TOP(%v) ", q.limit)
// 		}
// 	}

// 	sql := "SELECT " + limitPrefix
// 	if len(q.columns) > 0 {
// 		sql += strings.Join(q.columns, ", ")
// 	} else {
// 		sql += "*"
// 	}
// 	sql += " FROM " + q.tablename

// 	if len(q.filters) > 0 {
// 		sql += "\nWHERE " + strings.Join(q.filters, " ")
// 	}
// 	if len(q.groupby) > 0 {
// 		sql += "\nGROUP BY " + strings.Join(q.groupby, ", ")
// 	}
// 	if len(q.having) > 0 {
// 		sql += "\nHAVING " + q.having
// 	}
// 	if len(q.orderby) > 0 {
// 		sql += "\nORDER BY " + strings.Join(q.orderby, ", ")
// 	}
// 	if backend == MSSQL_BACKEND {
// 		if q.offset > 0 || q.limit > 0 {
// 			sql += fmt.Sprintf(
// 				"\nOFFSET %v ROWS", q.offset)
// 		}
// 		if q.limit > 0 {
// 			sql += fmt.Sprintf(
// 				"\nFETCH NEXT %v ROWS ONLY", q.limit)
// 		}
// 	} else {
// 		if q.limit > 0 {
// 			sql += fmt.Sprintf("\nLIMIT %v", q.limit)
// 		}
// 		if q.offset > 0 {
// 			sql += fmt.Sprintf("\nOFFSET %v", q.offset)
// 		}
// 	}
// 	sql += ";"

// 	// run query and fetch data
// 	result, err := q.handler.FetchAll(sql, q.execargs...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// apply DataReaders adapters
// 	if mdl, ok := q.model.(ModelDataReaders); ok {
// 		if err := FormatData(mdl.DataReaders(), result); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return result, nil
// }

// // return first element matching select query
// func (dbq *Query) First() (Data, error) {
// 	dbq.limit, dbq.offset = 1, 0
// 	result, err := dbq.All()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(result) > 0 {
// 		return result[0], nil
// 	}
// 	return nil, nil
// }

// // return one element matching filter params or nil
// // there must be only one element or none
// func (dbq *Query) One() (Data, error) {
// 	dbq.limit, dbq.offset = 2, 0
// 	result, err := dbq.All()
// 	if err != nil {
// 		return nil, err
// 	}
// 	switch len(result) {
// 	case 0:
// 		return nil, nil
// 	case 1:
// 		return result[0], nil
// 	}
// 	return nil, errors.New("multiple entries found")
// }

// // get model data defined by guid
// func (dbq *Query) Get(guid string) (Data, error) {
// 	dbq.filters = []string{"guid=" + SQL_PLACEHOLDER}
// 	dbq.execargs = []any{guid}
// 	dbq.groupby = []string{}
// 	dbq.orderby = []string{}
// 	dbq.having = ""
// 	return dbq.One()
// }

// // count number of entries matching defined filters
// func (dbq *Query) Count() (int64, error) {
// 	if dbq.lasterr != nil {
// 		return 0, dbq.lasterr
// 	}

// 	sql := "SELECT count(*) as count FROM " + dbq.tablename
// 	if len(dbq.filters) > 0 {
// 		sql += "\nWHERE " + strings.Join(dbq.filters, " ")
// 	}
// 	if len(dbq.groupby) > 0 {
// 		sql += "\nGROUP BY " + strings.Join(dbq.groupby, ", ")
// 	}
// 	sql += ";"

// 	result, err := dbq.handler.FetchAll(sql, dbq.execargs...)
// 	if err != nil {
// 		return 0, err
// 	}
// 	if len(result) > 0 {
// 		if count, ok := result[0]["count"].(int64); ok {
// 			return count, nil
// 		}
// 	}

// 	return 0, fmt.Errorf("invalid query result")
// }

// // insert data into table and return guid of new created entry
// func (dbq *Query) Insert(data Data) (string, error) {
// 	if dbq.lasterr != nil {
// 		return "", dbq.lasterr
// 	}

// 	// apply DataWriters adapters
// 	if mdl, ok := dbq.model.(ModelDataWriters); ok {
// 		if err := FormatData(mdl.DataWriters(), []Data{data}); err != nil {
// 			return "", err
// 		}
// 	}

// 	// check and create guid in data
// 	guid := ""
// 	if _, ok := dbq.model.(ModelAutoGuid); ok {
// 		guid = dictx.Fetch(data, "guid", "")
// 		if guid == "" {
// 			guid = NewGuid()
// 			dictx.Set(data, "guid", guid)
// 		}
// 	}

// 	columns, holders, execargs := []string{}, []string{}, []any{}
// 	for k, v := range data {
// 		columns = append(columns, k)
// 		holders = append(holders, SQL_PLACEHOLDER)
// 		execargs = append(execargs, v)
// 	}

// 	sql := "INSERT INTO " + dbq.tablename
// 	sql += fmt.Sprintf("\n(%v)", strings.Join(columns, ", "))
// 	sql += fmt.Sprintf("\nVALUES (%v)", strings.Join(holders, ", "))
// 	sql += ";"

// 	err := dbq.handler.Execute(sql, execargs...)
// 	if err != nil {
// 		return "", err
// 	}

// 	return guid, nil
// }

// // update table data and return number of affected entries
// func (dbq *Query) Update(data Data) (int64, error) {
// 	if dbq.lasterr != nil {
// 		return 0, dbq.lasterr
// 	}

// 	// check and remove guid from data
// 	if _, ok := dbq.model.(ModelAutoGuid); ok {
// 		dictx.Delete(data, "guid")
// 	}

// 	// apply DataWriters adapters
// 	if mdl, ok := dbq.model.(ModelDataWriters); ok {
// 		if err := FormatData(mdl.DataWriters(), []Data{data}); err != nil {
// 			return 0, err
// 		}
// 	}

// 	columns, execargs := []string{}, []any{}
// 	for k, v := range data {
// 		columns = append(columns, k+"="+SQL_PLACEHOLDER)
// 		execargs = append(execargs, v)
// 	}

// 	sql := "UPDATE " + dbq.tablename
// 	sql += "\nSET " + strings.Join(columns, ", ")
// 	if len(dbq.filters) > 0 {
// 		sql += "\nWHERE " + strings.Join(dbq.filters, " ")
// 	}

// 	execargs = append(execargs, dbq.execargs...)

// 	num_rows, err := dbq.handler.Execute(sql, execargs...)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return num_rows, nil
// }

// // delete from table data and return number of affected entries
// func (dbq *Query) Delete() (int64, error) {
// 	if dbq.lasterr != nil {
// 		return 0, dbq.lasterr
// 	}

// 	sql := "DELETE FROM " + dbq.tablename
// 	if len(dbq.filters) > 0 {
// 		sql += "\nWHERE " + strings.Join(dbq.filters, " ")
// 	}

// 	err := dbq.handler.Execute(sql, dbq.execargs...)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return dbq.handler.RowsAffected(), nil
// }

// ////////////////////////////////////////// helpers

// // check valid SQL identifier string
// func (dbq *Query) Ident(name string) string {
// 	match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", name)
// 	if match {
// 		return name
// 	}
// 	dbq.lasterr = fmt.Errorf("invalid sql identifier '%v'", name)
// 	return ""
// }
