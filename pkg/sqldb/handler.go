package sqldb

import (
	"strings"
)

const (
	defaultConnectTimeout = float64(5)
	defaultRetryInterval  = float64(0.2)
)

type Handler struct {
	Logger *Logger

	// database backend Engine handler
	Engine Engine

	// session connection params
	ConnectTimeout float64
	RetryInterval  float64
}

func NewHandler(dbe Engine, opts Options, logger *Logger) *Handler {
	return &Handler{
		Logger: logger,
		Engine: dbe,
		ConnectTimeout: opts.GetFloat64(
			"connect_timeout", defaultConnectTimeout),
		RetryInterval: opts.GetFloat64(
			"retry_interval", defaultRetryInterval),
	}
}

func (dbh *Handler) Session() *Session {
	return newSession(dbh)
}

func (dbh *Handler) CreateSchema(models map[string]Model) error {
	schema_sql := []string{}
	for tblname, model := range models {
		sql, err := dbh.Engine.GenSchema(tblname, model)
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
		if err := dbs.Execute(sql); err != nil {
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
