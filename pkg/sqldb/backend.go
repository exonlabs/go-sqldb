package sqldb

// all implemented backend engines
var BACKENDS = []string{
	SQLITE_BACKEND,
	MYSQL_BACKEND,
	PGSQL_BACKEND,
	MSSQL_BACKEND,
}

type Engine interface {
	BackendName() string
	GenSchema(string, Model) ([]string, error)
	// FormatStatment(string) string
}

type BackendManager interface {
	InteractiveConfig(Options) (Options, error)
	InteractiveSetup(Options) error
}
