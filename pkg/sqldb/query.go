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
	// database session
	dbs *Session
	// the query attached model
	model Model
	// SQL statment attributes
	attrs StmtAttrs
}

// NewQuery creates a new query object
func NewQuery(dbs *Session, model Model) *Query {
	q := &Query{
		dbs:   dbs,
		model: model,
	}
	if model != nil {
		// set the initial values
		q.attrs.Tablename = model.TableName()
		q.attrs.Columns = model.Columns()
		q.attrs.Orderby = model.Orders()
		q.attrs.Limit = model.Limit()
	}
	return q
}

// Session returns the associated session handler
func (q *Query) Session() *Session {
	return q.dbs
}

// TableName sets the table name in statment.
func (q *Query) TableName(name string) *Query {
	name = strings.TrimSpace(name)
	if name != "" {
		q.attrs.Tablename = name
	}
	return q
}

// Columns sets the columns in statment.
func (q *Query) Columns(columns ...string) *Query {
	if len(columns) > 0 {
		q.attrs.Columns = columns
		if len(q.attrs.Orderby) == 0 {
			q.attrs.Orderby = []string{columns[0] + " ASC"}
		}
	} else {
		q.attrs.Columns = nil
		q.attrs.Orderby = nil
	}
	return q
}

// Filters sets filtering expresion to the statment with args.
func (q *Query) Filters(expr string, args ...any) *Query {
	q.attrs.Filters = strings.TrimSpace(expr)
	q.attrs.FiltersArgs = args
	return q
}

// FilterBy adds AND related filter to the statment.
func (q *Query) FilterBy(column string, value any) *Query {
	if column != "" {
		if q.attrs.Filters != "" {
			q.attrs.Filters += " AND "
		}
		q.attrs.Filters = fmt.Sprintf("%s=%s", column, SQL_PLACEHOLDER)
		q.attrs.FiltersArgs = append(q.attrs.FiltersArgs, value)
	}
	return q
}

// GroupBy adds grouping expresion to the statment.
func (q *Query) GroupBy(columns ...string) *Query {
	q.attrs.Groupby = columns
	return q
}

// OrderBy adds ordering expresion to the statment.
// orders has the format: "column ASC|DESC"
func (q *Query) OrderBy(orders ...string) *Query {
	q.attrs.Orderby = orders
	return q
}

// Having adds having expr in the statment.
func (q *Query) Having(expr string, args ...any) *Query {
	q.attrs.Having = expr
	q.attrs.HavingArgs = args
	return q
}

// Offset add offset in the statment.
func (q *Query) Offset(offset int) *Query {
	q.attrs.Offset = offset
	return q
}

// Limit adds limit in the statment.
func (q *Query) Limit(limit int) *Query {
	q.attrs.Limit = limit
	return q
}

// check attrs before running query
func (q *Query) check_run() error {
	if q.attrs.Tablename == "" {
		return fmt.Errorf("%w - empty table name", ErrOperation)
	}
	if q.dbs == nil {
		return ErrDBSession
	}
	return q.dbs.check_run()
}

// All returns all data entries matching defined filters.
func (q *Query) All() ([]Data, error) {
	if err := q.check_run(); err != nil {
		return nil, err
	}

	// generate and run query
	stmt, params := q.dbs.db.engine.SqlGenerator().Select(&q.attrs)
	result, err := q.dbs.Fetch(stmt, params...)
	if err != nil {
		return nil, err
	}

	// apply decoding on result data
	if len(result) > 0 {
		if err := q.model.DataDecode(result); err != nil {
			return nil, fmt.Errorf(
				"%w - decoding data failed, %v", ErrOperation, err)
		}
	}

	return result, nil
}

// First returns the first data entry matching defined filters.
func (q *Query) First() (Data, error) {
	q.attrs.Offset, q.attrs.Limit = 0, 1
	result, err := q.All()
	if len(result) >= 1 {
		return result[0], nil
	}
	return nil, err
}

// One returns and check that only one data entry matches the defined filters.
// there must be only one element matched or none, else an error is returned.
func (q *Query) One() (Data, error) {
	q.attrs.Offset, q.attrs.Limit = 0, 2
	result, err := q.All()
	if len(result) >= 2 {
		return nil, fmt.Errorf("%w - multiple entries found", ErrOperation)
	} else if len(result) >= 1 {
		return result[0], nil
	}
	return nil, err
}

// Get is a short form to fetch only one element by guid.
// there must be a guid primary column in model.
func (q *Query) Get(guid string) (Data, error) {
	return q.FilterBy("guid", guid).One()
}

// Counts the number of table entries matching defined filters.
func (q *Query) Count() (int, error) {
	if err := q.check_run(); err != nil {
		return 0, err
	}

	// generate and run query
	stmt, params := q.dbs.db.engine.SqlGenerator().Count(&q.attrs)
	result, err := q.dbs.Fetch(stmt, params...)
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

// Inserts new data entry and returns the guid for new entry.
// If Model AutoGuid is enabled, a new guid value is generated when the
// insert data have empty or no guid value.
func (q *Query) Insert(data Data) (string, error) {
	if data == nil {
		return "", fmt.Errorf("%w - empty insert data", ErrOperation)
	}
	if err := q.check_run(); err != nil {
		return "", err
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

	// generate and run query
	stmt, params := q.dbs.db.engine.SqlGenerator().Insert(&q.attrs, data)
	_, err := q.dbs.Exec(stmt, params...)
	if err != nil {
		return "", err
	}

	return guid, nil
}

// Updates data entries matching defined filters and returns the number
// of affected entries.
func (q *Query) Update(data Data) (int, error) {
	if data == nil {
		return 0, fmt.Errorf("%w - empty update data", ErrOperation)
	}
	if err := q.check_run(); err != nil {
		return 0, err
	}

	// apply encoding on update data
	if err := q.model.DataEncode([]Data{data}); err != nil {
		return 0, fmt.Errorf(
			"%w - encoding data error, %v", ErrOperation, err)
	}

	// generate and run query
	stmt, params := q.dbs.db.engine.SqlGenerator().Update(&q.attrs, data)
	return q.dbs.Exec(stmt, params...)
}

// Deletes data entries matching defined filters and returns the number
// of affected entries.
func (q *Query) Delete() (int, error) {
	if err := q.check_run(); err != nil {
		return 0, err
	}

	// generate and run query
	stmt, params := q.dbs.db.engine.SqlGenerator().Delete(&q.attrs)
	return q.dbs.Exec(stmt, params...)
}
