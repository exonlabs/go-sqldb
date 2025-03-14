// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import "strings"

// Model defines the model interface.
type Model interface {
	// TableMeta returns the model table metainfo.
	TableMeta() *TableMeta

	// TableName returns the table name to use in statments.
	TableName() string
	// Columns returns the table columns to use in statments.
	Columns() []string
	// Orders returns the table order-by columns to use in statments.
	Orders() []string
	// Limit returns the number of rows to fetch in statments.
	Limit() int
	// IsAutoGuid returns true if the AutoGuid operations are enabled.
	IsAutoGuid() bool
	// DataEncode applies encoding to data before writing to database.
	DataEncode([]Data) error
	// DataDecode applies decoding on data after reading from database.
	DataDecode([]Data) error

	// PreSchema is called before creating the table schema in database.
	PreSchema(dbs *session, meta *ModelMeta) error
	// PostSchema is called after creating the table schema in database.
	PostSchema(dbs *session, meta *ModelMeta) error

	// InitialData creates the initial data in table.
	InitialData(dbs *session, tablename string) error
}

// ModelMeta represents link between tablename and table metainfo.
type ModelMeta struct {
	Table string
	Model Model
}

// BaseModel defines a base model structure.
type BaseModel struct {
	// DefaultTable defines the default tablename to use in statments.
	DefaultTable string
	// DefaultColumns defines the default set of column to use in statments.
	// leave empty to use all columns.
	DefaultColumns []string
	// DefaultOrders defines the default orders to use in statments.
	// columns in this list should be present in DefaultColumns.
	DefaultOrders []string
	// DefaultLimit defines the default limit to use in statments.
	DefaultLimit int
	// AutoGuid enables the auto guid operations: which are to create new guid
	// for inserts and prevent guid column change in updates.
	AutoGuid bool
}

// TableName returns the table name to use in statments.
func (m *BaseModel) TableName() string {
	return m.DefaultTable
}

// Columns returns the table columns names to use in statments.
func (m *BaseModel) Columns() []string {
	return m.DefaultColumns
}

// Orders returns the table order-by columns to use in statments.
func (m *BaseModel) Orders() []string {
	return m.DefaultOrders
}

// Limit returns the number of rows to fetch in statments.
func (m *BaseModel) Limit() int {
	return m.DefaultLimit
}

// IsAutoGuid returns true if the AutoGuid operations are enabled.
func (m *BaseModel) IsAutoGuid() bool {
	return m.AutoGuid
}

// DataEncode applies encoding to data before writing to database.
func (m *BaseModel) DataEncode([]Data) error {
	return nil
}

// DataDecode applies decoding on data after reading from database.
func (m *BaseModel) DataDecode([]Data) error {
	return nil
}

// PreSchema is called before creating the table schema in database.
func (m *BaseModel) PreSchema(dbs *session, meta *ModelMeta) error {
	return nil
}

// PostSchema is called after creating the table schema in database.
func (m *BaseModel) PostSchema(dbs *session, meta *ModelMeta) error {
	return nil
}

// InitialData creates the initial data in table.
func (m *BaseModel) InitialData(dbs *session, tablename string) error {
	return nil
}

////////////////////////////////////////////////////

// InitializeModels creates and alter the database models schema,
// then adds the models intial data.
func InitializeModels(db *Database, metainfo []ModelMeta) error {
	if db == nil {
		return ErrDBHandler
	} else if db.engine == nil {
		return ErrDBEngine
	}

	// create new session
	dbs := db.Session()

	// create and alter schema
	if db.Log != nil {
		db.Log.Debug("creating models schema")
	}
	for _, meta := range metainfo {
		if err := meta.Model.PreSchema(dbs, &meta); err != nil {
			return err
		}
	}
	for _, meta := range metainfo {
		stmts := dbs.db.engine.SqlGenerator().
			Schema(meta.Table, meta.Model.TableMeta())
		if _, err := dbs.Exec(strings.Join(stmts, "\n")); err != nil {
			return err
		}
	}
	for _, meta := range metainfo {
		if err := meta.Model.PostSchema(dbs, &meta); err != nil {
			return err
		}
	}

	// add intial data to tables
	if db.Log != nil {
		db.Log.Debug("adding models initial data")
	}
	for _, meta := range metainfo {
		if err := meta.Model.InitialData(dbs, meta.Table); err != nil {
			return err
		}
	}

	return nil
}
