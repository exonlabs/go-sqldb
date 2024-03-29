package sqldb

const PGSQL_BACKEND = "pgsql"

type PgSqlEngine struct{}

func NewPgSqlEngine() *PgSqlEngine {
	return &PgSqlEngine{}
}

func (dbe *PgSqlEngine) BackendName() string {
	return PGSQL_BACKEND
}
