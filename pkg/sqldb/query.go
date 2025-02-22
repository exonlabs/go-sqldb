// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
)

// Query represents the query object.
type Query struct {
	dbs   *Session
	model Model

	// table name to use other than default, this allows for mapping
	// same model to multiple tables at runtime
	tablename string

	// columns list to fetch from table
	columns []string

	// query filters and args
	filter      string
	filterby    []string
	filtersArgs []any

	// query grouping, having and ordering
	groupby    []string
	orderby    []string
	having     string
	havingArgs []any

	// query data set limit and offset
	offset int
	limit  int
}

// NewQuery creates a new query object
func NewQuery(dbs *Session, model Model) *Query {
	return &Query{
		dbs:       dbs,
		model:     model,
		tablename: model.TableName(),
		columns:   model.Columns(),
		orderby:   model.Orders(),
	}
}

// TableName sets the runtime table name in statment.
func (q *Query) TableName(name string) *Query {
	name = strings.TrimSpace(name)
	if name != "" {
		q.tablename = name
	}
	return q
}

// Columns sets the columns in statment.
func (q *Query) Columns(columns ...string) *Query {
	if len(columns) > 0 {
		q.columns = columns
		q.orderby = []string{columns[0] + " ASC"}
	}
	return q
}

// Filter adds filter expresion to the statment with args. using this
// function overides any filters added using the Query.FilterBy function.
func (q *Query) Filter(expr string, args ...any) *Query {
	if expr != "" {
		q.filterby = nil
		q.filter = expr
		q.filtersArgs = args
	}
	return q
}

// FilterBy adds multiple AND related filters to the statment. using this
// function overides any filters added using the Query.Filter function.
func (q *Query) FilterBy(column string, value any) *Query {
	if column != "" {
		if q.filter != "" {
			q.filter = ""
			q.filtersArgs = nil
		}
		q.filterby = append(q.filterby,
			fmt.Sprintf("%s=%s", column, SQL_PLACEHOLDER))
		q.filtersArgs = append(q.filtersArgs, value)
	}
	return q
}

// GroupBy adds grouping expresion to the statment.
func (q *Query) GroupBy(columns ...string) *Query {
	if len(columns) > 0 {
		q.groupby = columns
	}
	return q
}

// OrderBy adds ordering expresion to the statment.
// order_expr has the format: "column_name ASC|DESC"
func (q *Query) OrderBy(order_expr ...string) *Query {
	if len(order_expr) > 0 {
		q.orderby = order_expr
	}
	return q
}

// Having adds having expr in statment
func (q *Query) Having(expr string, args ...any) *Query {
	if expr != "" {
		q.having = expr
		q.havingArgs = args
	}
	return q
}

// Offset add offset in statment
func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

// Limit adds limit in statment
func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

// All returns all data entries from table matching defined filters.
func (q *Query) All() ([]Data, error) {
	// create the statment
	stmt := "SELECT "
	if q.dbs.db.Backend() == BACKEND_MSSQL &&
		q.limit > 0 && len(q.orderby) == 0 {
		stmt += fmt.Sprintf("TOP(%d) ", q.limit)
	}
	if len(q.columns) > 0 {
		stmt += strings.Join(q.columns, ", ")
	} else {
		stmt += "*"
	}
	stmt += " FROM " + q.tablename

	if q.filter != "" {
		stmt += "\nWHERE " + q.filter
	} else if len(q.filterby) > 0 {
		stmt += "\nWHERE " + strings.Join(q.filterby, " AND ")
	}
	if len(q.groupby) > 0 {
		stmt += "\nGROUP BY " + strings.Join(q.groupby, ", ")
	}
	if q.having != "" {
		stmt += "\nHAVING " + q.having
	}
	if len(q.orderby) > 0 {
		stmt += "\nORDER BY " + strings.Join(q.orderby, ", ")
	}
	if q.dbs.db.Backend() == BACKEND_MSSQL {
		if q.offset > 0 || q.limit > 0 {
			stmt += fmt.Sprintf("\nOFFSET %d ROWS", q.offset)
		}
		if q.limit > 0 {
			stmt += fmt.Sprintf("\nFETCH NEXT %d ROWS ONLY", q.limit)
		}
	} else {
		if q.offset > 0 {
			stmt += fmt.Sprintf("\nOFFSET %d", q.offset)
		}
		if q.limit > 0 {
			stmt += fmt.Sprintf("\nLIMIT %d", q.limit)
		}
	}
	stmt += ";"

	// create the params for statment placeholders
	params := append(q.filtersArgs, q.havingArgs...)

	// run query and fetch data
	result, err := q.dbs.FetchAll(stmt, params...)
	if err != nil {
		return nil, err
	}

	// apply decoding on result data
	if len(result) > 0 {
		if err := q.model.DataDecode(result); err != nil {
			return nil, fmt.Errorf(
				"%w - decoding data error, %v", ErrOperation, err)
		}
	}

	return result, nil
}

// First returns the first data entry from table matching defined filters.
func (q *Query) First() (Data, error) {
	q.offset, q.limit = 0, 1
	result, err := q.All()
	if len(result) >= 1 {
		return result[0], nil
	}
	return nil, err
}

// One returns only one data entry from table matching defined filters.
// there must be none or only one element matched else an error is returned.
func (q *Query) One() (Data, error) {
	q.offset, q.limit = 0, 2
	result, err := q.All()
	if len(result) >= 2 {
		return nil, fmt.Errorf("%w - multiple entries found", ErrOperation)
	} else if len(result) >= 1 {
		return result[0], nil
	}
	return nil, err
}

// Get is a short form to fetch only one element by column named guid.
func (q *Query) Get(guid string) (Data, error) {
	return q.FilterBy("guid", guid).One()
}

// Counts the number of table entries matching defined filters.
func (q *Query) Count() (int, error) {
	// create the statment
	stmt := "SELECT count(*) as count FROM " + q.tablename
	if q.filter != "" {
		stmt += "\nWHERE " + q.filter
	} else if len(q.filterby) > 0 {
		stmt += "\nWHERE " + strings.Join(q.filterby, " AND ")
	}
	if len(q.groupby) > 0 {
		stmt += "\nGROUP BY " + strings.Join(q.groupby, ", ")
	}
	if q.having != "" {
		stmt += "\nHAVING " + q.having
	}
	stmt += ";"

	// create the params for statment placeholders
	params := append(q.filtersArgs, q.havingArgs...)

	// run query and fetch data
	result, err := q.dbs.FetchAll(stmt, params...)
	if err != nil {
		return 0, err
	}
	if len(result) > 0 {
		if count, ok := result[0]["count"]; ok {
			n, err := strconv.ParseInt(fmt.Sprint(count), 10, 64)
			if err == nil {
				return int(n), nil
			}
		}
	}

	return 0, fmt.Errorf("%w - invalid query result", ErrOperation)
}

// Inserts data into table and returns the guid for new entry.
// If Model AutoGuid is enabled, a new guid value is generated when the
// insert data have empty or no guid value.
func (q *Query) Insert(data Data) (string, error) {
	if data == nil {
		return "", fmt.Errorf("%w - empty insert data", ErrOperation)
	}

	// apply encoding on insert data
	if err := q.model.DataEncode([]Data{data}); err != nil {
		return "", fmt.Errorf(
			"%w - encoding data error, %v", ErrOperation, err)
	}

	// check and create guid in data
	guid := dictx.Fetch(data, "guid", "")
	if q.model.IsAutoGuid() && guid == "" {
		guid = NewGuid()
		dictx.Set(data, "guid", guid)
	}

	columns, holders, params := []string{}, []string{}, []any{}
	for k, v := range data {
		columns = append(columns, k)
		holders = append(holders, SQL_PLACEHOLDER)
		params = append(params, v)
	}

	// create the statment
	stmt := "INSERT INTO " + q.tablename
	stmt += fmt.Sprintf("\n(%v)", strings.Join(columns, ", "))
	stmt += fmt.Sprintf("\nVALUES (%v)", strings.Join(holders, ", "))
	stmt += ";"

	_, err := q.dbs.Execute(stmt, params...)
	if err != nil {
		return "", err
	}

	return guid, nil
}

// Updates data in table matching defined filters and returns the number
// of affected entries.
func (q *Query) Update(data Data) (int, error) {
	if data == nil {
		return 0, fmt.Errorf("%w - empty update data", ErrOperation)
	}

	// apply encoding on update data
	if err := q.model.DataEncode([]Data{data}); err != nil {
		return 0, fmt.Errorf(
			"%w - encoding data error, %v", ErrOperation, err)
	}

	columns, params := []string{}, []any{}
	for k, v := range data {
		columns = append(columns, k+"="+SQL_PLACEHOLDER)
		params = append(params, v)
	}

	stmt := "UPDATE " + q.tablename
	stmt += "\nSET " + strings.Join(columns, ", ")
	if q.filter != "" {
		stmt += "\nWHERE " + q.filter
	} else if len(q.filterby) > 0 {
		stmt += "\nWHERE " + strings.Join(q.filterby, " AND ")
	}
	params = append(params, q.filtersArgs...)

	return q.dbs.Execute(stmt, params...)
}

// Deletes data from table matching defined filters and returns the number
// of affected entries.
func (q *Query) Delete() (int, error) {
	stmt := "DELETE FROM " + q.tablename
	if q.filter != "" {
		stmt += "\nWHERE " + q.filter
	} else if len(q.filterby) > 0 {
		stmt += "\nWHERE " + strings.Join(q.filterby, " AND ")
	}

	return q.dbs.Execute(stmt, q.filtersArgs...)
}
