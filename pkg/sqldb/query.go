package sqldb

const SQL_PLACEHOLDER = "$?"

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

	// // last runtime error
	// lasterr error
}

func newQuery(dbs *Session, dbm Model) *Query {
	return &Query{
		session:   dbs,
		model:     dbm,
		tablename: dbm.TableName(),
		orderby:   dbm.DefaultOrders(),
	}
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
	return nil, nil
}

// return first element matching select query
func (dbq *Query) First() (Data, error) {
	return nil, nil
}

// return one element matching filter params or nil
// there must be only one element or none
func (dbq *Query) One() (Data, error) {
	return nil, nil
}

// get model data defined by guid
func (dbq *Query) Get(guid string) (Data, error) {
	return nil, nil
}

// count number of entries matching defined filters
func (dbq *Query) Count() (int64, error) {
	return 0, nil
}

// insert data into table and return guid of new created entry
func (dbq *Query) Insert(data Data) (string, error) {
	return "", nil
}

// update table data and return number of affected entries
func (dbq *Query) Update(data Data) (int64, error) {
	return 0, nil
}

// delete from table data and return number of affected entries
func (dbq *Query) Delete() (int64, error) {
	return 0, nil
}

////////////////////////////////////////// helpers
