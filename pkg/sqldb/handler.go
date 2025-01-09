package sqldb

import (
	"database/sql"
	"strings"
)

const (
	defaultConnectTimeout = float64(5)
	defaultRetryInterval  = float64(0.2)
)

type Handler struct {
	Logger *Logger

	// database backend engine handler
	sqlDB *sql.DB
	sqlTX *sql.Tx
	// engine Engine
	// dsn    string

	// session connection params
	ConnectTimeout float64
	RetryInterval  float64
}

func NewHandler(dbe Engine, opts Options, logger *Logger) *Handler {
	h := &Handler{
		Logger: logger,
		engine: dbe,
		ConnectTimeout: opts.GetFloat64(
			"connect_timeout", defaultConnectTimeout),
		RetryInterval: opts.GetFloat64(
			"retry_interval", defaultRetryInterval),
	}
	h.SetConfig(opts)
	return h
}

func (dbh *Handler) Engine() Engine {
	return dbh.engine
}

func (dbh *Handler) Session() *Session {
	return newSession(dbh)
}

func (dbh *Handler) SetConfig(opts Options) error {

	return nil
}

// func (dbs *Session) Query(dbm Model) *Query {
// 	return newQuery(dbs, dbm)
// }

// func (dbs *Session) IsActive() bool {
// 	if dbs.sqlDB != nil {
// 		return dbs.sqlDB.Ping() == nil
// 	}
// 	return false
// }

// func (dbs *Session) InTransaction() bool {
// 	return dbs.sqlTX != nil
// }

// func (dbs *Session) Open() error {
// 	return nil
// }

// func (dbs *Session) Close() error {
// 	return nil
// }

// func (dbs *Session) Begin() error {
// 	return nil
// }

// func (dbs *Session) Commit() error {
// 	return nil
// }

// func (dbs *Session) RollBack() error {
// 	return nil
// }

// func (dbs *Session) Execute(sql string, params ...any) (int64, error) {

// 	return 0, nil
// }

// func (dbs *Session) FetchAll(sql string, params ...any) ([]Data, error) {
// 	return nil, nil
// }


func (dbh *Handler) CreateSchema(models map[string]Model) error {
	schema_sql := []string{}
	for tblname, model := range models {
		sql, err := dbh.engine.GenSchema(tblname, model)
		if err != nil {
			return err
		}
		schema_sql = append(schema_sql, sql...)
		dbh.Logger.Trace2("table '%s' schema:\n%s",
			tblname, strings.Join(sql, "\n"))
	}

	dbs := dbh.Session()
	defer dbs.Close()

	// create schema in session transaction
	if err := dbs.Begin(); err != nil {
		return err
	}
	// create tables schema
	for _, sql := range schema_sql {
		if _, err := dbs.Execute(sql); err != nil {
			return err
		}
	}
	// run tables schema upgrades
	for tblname, model := range models {
		if mdl, ok := model.(ModelUpgradeSchema); ok {
			if err := mdl.UpgradeSchema(dbs, tblname); err != nil {
				return err
			}
		}
	}
	if err := dbs.Commit(); err != nil {
		return err
	}
	return nil
}

func (dbh *Handler) InitializeData(models map[string]Model) error {
	dbs := dbh.Session()
	defer dbs.Close()

	// initialze data in session transaction
	if err := dbs.Begin(); err != nil {
		return err
	}
	for tblname, model := range models {
		if mdl, ok := model.(ModelInitializeData); ok {
			if err := mdl.InitializeData(dbs, tblname); err != nil {
				return err
			}
		}
	}
	if err := dbs.Commit(); err != nil {
		return err
	}
	return nil
}
