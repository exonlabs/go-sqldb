package sqldb

const MYSQL_BACKEND = "mysql"

type MySqlEngine struct{}

func NewMySqlEngine() *MySqlEngine {
	return &MySqlEngine{}
}

func (dbe *MySqlEngine) BackendName() string {
	return MYSQL_BACKEND
}
