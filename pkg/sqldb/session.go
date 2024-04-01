package sqldb

type Session struct {
	handler *Handler
}

func newSession(dbh *Handler) *Session {
	return &Session{
		handler: dbh,
	}
}

func (dbs *Session) Query(dbm Model) *Query {
	return newQuery(dbs, dbm)
}

func (dbs *Session) IsActive() bool {
	return false
}

func (dbs *Session) InTransaction() bool {
	return false
}

func (dbs *Session) Open() error {
	return nil
}

func (dbs *Session) Close() error {
	return nil
}

func (dbs *Session) Begin() error {
	return nil
}

func (dbs *Session) Commit() error {
	return nil
}

func (dbs *Session) RollBack() error {
	return nil
}

func (dbs *Session) Execute(sql string, params ...any) error {
	return nil
}

func (dbs *Session) FetchAll(sql string, params ...any) ([]Data, error) {
	return nil, nil
}

func (dbs *Session) RowsAffected() int64 {
	return 0
}
