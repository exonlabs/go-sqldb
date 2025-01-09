package sqldb

// all implemented backends names
func Backends() []string {
	return []string{
		SQLITE_BACKEND,
		MYSQL_BACKEND,
		PGSQL_BACKEND,
		MSSQL_BACKEND,
	}
}

type Engine interface {
	Backend() string
	FormatSql(string) string
	CanRetryErr(error) bool
}

type Backend interface {
	CreateSchema(string, Model) ([]string, error)
	InteractiveConfig(Options) (Options, error)
	InteractiveSetup(Options) error
}
