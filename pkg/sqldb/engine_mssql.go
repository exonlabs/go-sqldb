package sqldb

const MSSQL_BACKEND = "mssql"

type MsSqlEngine struct{}

func NewMsSqlEngine() *MsSqlEngine {
	return &MsSqlEngine{}
}

func (dbe *MsSqlEngine) BackendName() string {
	return MSSQL_BACKEND
}
