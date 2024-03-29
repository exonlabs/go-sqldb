package sqldb

const SQLITE_BACKEND = "sqlite"

type SqliteEngine struct{}

func NewSqliteEngine() *SqliteEngine {
	return &SqliteEngine{}
}

func (dbe *SqliteEngine) BackendName() string {
	return SQLITE_BACKEND
}
