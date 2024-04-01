package sqldb

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Query struct {
	session *Session
	model   Model

	// runtime table name to use, this allows for mapping
	// same model to multiple tables
	tablename string
	// query columns
	columns []string
	// query sql filters and args
	filters  []string
	execargs []any
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

func newQuery(dbs *Session, dbm Model) *Query {
	q := &Query{
		session:   dbs,
		model:     dbm,
		tablename: dbm.TableName(),
	}
	if mdl, ok := dbm.(ModelDefaultOrders); ok {
		q.orderby = mdl.DefaultOrders()
	}
	return q
}

////////////////////////////////////////// Creation

// set runtime table name instead of default
func (dbq *Query) Table(name string) *Query {
	dbq.tablename = name
	return dbq
}

// set columns
func (dbq *Query) Columns(columns ...string) *Query {
	dbq.columns = columns
	return dbq
}

// add filters
func (dbq *Query) Filter(expr string, params ...any) *Query {
	dbq.filters = append(dbq.filters, expr)
	dbq.execargs = append(dbq.execargs, params...)
	return dbq
}

// add filters
func (dbq *Query) FilterBy(column string, value any) *Query {
	cond := ""
	if len(dbq.filters) > 0 {
		cond = "AND "
	}
	dbq.filters = append(
		dbq.filters, cond+column+"="+SQL_PLACEHOLDER)
	dbq.execargs = append(dbq.execargs, value)
	return dbq
}

// add grouping
func (dbq *Query) GroupBy(groupby ...string) *Query {
	dbq.groupby = groupby
	return dbq
}

// add ordering: "colname ASC|DESC"
func (dbq *Query) OrderBy(orderby ...string) *Query {
	dbq.orderby = orderby
	return dbq
}

// add having expr
func (dbq *Query) Having(expr string, val any) *Query {
	dbq.having = expr
	dbq.execargs = append(dbq.execargs, val)
	return dbq
}

// add limit
func (dbq *Query) Limit(limit int) *Query {
	dbq.limit = limit
	return dbq
}

// add offset
func (dbq *Query) Offset(offset int) *Query {
	dbq.offset = offset
	return dbq
}

////////////////////////////////////////// Operations

// return all elements matching select query
func (dbq *Query) All() ([]Data, error) {
	if dbq.lasterr != nil {
		return nil, dbq.lasterr
	}

	backend := dbq.session.handler.Engine.BackendName()
	limitPrefix := ""
	if backend == MSSQL_BACKEND {
		if dbq.limit > 0 && len(dbq.orderby) == 0 {
			limitPrefix = fmt.Sprintf("TOP(%v) ", dbq.limit)
		}
	}

	sql := "SELECT " + limitPrefix
	if len(dbq.columns) > 0 {
		sql += strings.Join(dbq.columns, ", ")
	} else {
		sql += "*"
	}
	sql += " FROM " + dbq.tablename

	if len(dbq.filters) > 0 {
		sql += "\nWHERE " + strings.Join(dbq.filters, " ")
	}
	if len(dbq.groupby) > 0 {
		sql += "\nGROUP BY " + strings.Join(dbq.groupby, ", ")
	}
	if len(dbq.having) > 0 {
		sql += "\nHAVING " + dbq.having
	}
	if len(dbq.orderby) > 0 {
		sql += "\nORDER BY " + strings.Join(dbq.orderby, ", ")
	}
	if backend == MSSQL_BACKEND {
		if dbq.offset > 0 || dbq.limit > 0 {
			sql += fmt.Sprintf(
				"\nOFFSET %v ROWS", dbq.offset)
		}
		if dbq.limit > 0 {
			sql += fmt.Sprintf(
				"\nFETCH NEXT %v ROWS ONLY", dbq.limit)
		}
	} else {
		if dbq.limit > 0 {
			sql += fmt.Sprintf("\nLIMIT %v", dbq.limit)
		}
		if dbq.offset > 0 {
			sql += fmt.Sprintf("\nOFFSET %v", dbq.offset)
		}
	}
	sql += ";"

	// run query and fetch data
	result, err := dbq.session.FetchAll(sql, dbq.execargs...)
	if err != nil {
		return nil, err
	}

	// apply DataReaders adapters
	if mdl, ok := dbq.model.(ModelDataReaders); ok {
		if err := FormatData(mdl.DataReaders(), result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// return first element matching select query
func (dbq *Query) First() (Data, error) {
	dbq.limit, dbq.offset = 1, 0
	result, err := dbq.All()
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}

// return one element matching filter params or nil
// there must be only one element or none
func (dbq *Query) One() (Data, error) {
	dbq.limit, dbq.offset = 2, 0
	result, err := dbq.All()
	if err != nil {
		return nil, err
	}
	switch len(result) {
	case 0:
		return nil, nil
	case 1:
		return result[0], nil
	}
	return nil, errors.New("multiple entries found")
}

// get model data defined by guid
func (dbq *Query) Get(guid string) (Data, error) {
	dbq.filters = []string{"guid=" + SQL_PLACEHOLDER}
	dbq.execargs = []any{guid}
	dbq.groupby = []string{}
	dbq.orderby = []string{}
	dbq.having = ""
	return dbq.One()
}

// count number of entries matching defined filters
func (dbq *Query) Count() (int64, error) {
	if dbq.lasterr != nil {
		return 0, dbq.lasterr
	}

	sql := "SELECT count(*) as count FROM " + dbq.tablename
	if len(dbq.filters) > 0 {
		sql += "\nWHERE " + strings.Join(dbq.filters, " ")
	}
	if len(dbq.groupby) > 0 {
		sql += "\nGROUP BY " + strings.Join(dbq.groupby, ", ")
	}
	sql += ";"

	result, err := dbq.session.FetchAll(sql, dbq.execargs...)
	if err != nil {
		return 0, err
	}
	if len(result) > 0 {
		if count, ok := result[0]["count"].(int64); ok {
			return count, nil
		}
	}

	return 0, fmt.Errorf("invalid query result")
}

// insert data into table and return guid of new created entry
func (dbq *Query) Insert(data Data) (string, error) {
	if dbq.lasterr != nil {
		return "", dbq.lasterr
	}

	// apply DataWriters adapters
	if mdl, ok := dbq.model.(ModelDataWriters); ok {
		if err := FormatData(mdl.DataWriters(), []Data{data}); err != nil {
			return "", err
		}
	}

	// check and create guid in data
	guid := ""
	if _, ok := dbq.model.(ModelSetAutoGuid); ok {
		guid = data.GetString("guid", "")
		if guid == "" {
			guid = NewGuid()
			data.Set("guid", guid)
		}
	}

	columns, holders, execargs := []string{}, []string{}, []any{}
	for k, v := range data {
		columns = append(columns, k)
		holders = append(holders, SQL_PLACEHOLDER)
		execargs = append(execargs, v)
	}

	sql := "INSERT INTO " + dbq.tablename
	sql += fmt.Sprintf("\n(%v)", strings.Join(columns, ", "))
	sql += fmt.Sprintf("\nVALUES (%v)", strings.Join(holders, ", "))
	sql += ";"

	err := dbq.session.Execute(sql, execargs...)
	if err != nil {
		return "", err
	}

	return guid, nil
}

// update table data and return number of affected entries
func (dbq *Query) Update(data Data) (int64, error) {
	if dbq.lasterr != nil {
		return 0, dbq.lasterr
	}

	// check and remove guid from data
	if _, ok := dbq.model.(ModelSetAutoGuid); ok {
		data.Del("guid")
	}

	// apply DataWriters adapters
	if mdl, ok := dbq.model.(ModelDataWriters); ok {
		if err := FormatData(mdl.DataWriters(), []Data{data}); err != nil {
			return 0, err
		}
	}

	columns, execargs := []string{}, []any{}
	for k, v := range data {
		columns = append(columns, k+"="+SQL_PLACEHOLDER)
		execargs = append(execargs, v)
	}

	sql := "UPDATE " + dbq.tablename
	sql += "\nSET " + strings.Join(columns, ", ")
	if len(dbq.filters) > 0 {
		sql += "\nWHERE " + strings.Join(dbq.filters, " ")
	}

	execargs = append(execargs, dbq.execargs...)

	err := dbq.session.Execute(sql, execargs...)
	if err != nil {
		return 0, err
	}

	return dbq.session.RowsAffected(), nil
}

// delete from table data and return number of affected entries
func (dbq *Query) Delete() (int64, error) {
	if dbq.lasterr != nil {
		return 0, dbq.lasterr
	}

	sql := "DELETE FROM " + dbq.tablename
	if len(dbq.filters) > 0 {
		sql += "\nWHERE " + strings.Join(dbq.filters, " ")
	}

	err := dbq.session.Execute(sql, dbq.execargs...)
	if err != nil {
		return 0, err
	}

	return dbq.session.RowsAffected(), nil
}

////////////////////////////////////////// helpers

// check valid SQL identifier string
func (dbq *Query) Ident(name string) string {
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", name)
	if match {
		return name
	}
	dbq.lasterr = fmt.Errorf("invalid sql identifier '%v'", name)
	return ""
}
