// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

// Model defines the model interface.
type Model interface {
	// ModelMeta returns the model table metainfo.
	ModelMeta() *TableMeta

	// TableName returns the table name to use in statments.
	TableName() string
	// Columns returns the table columns to use in statments.
	Columns() []string
	// Orders returns the table order-by columns to use in statments.
	Orders() []string
	// IsAutoGuid returns true if the AutoGuid operations are enabled.
	IsAutoGuid() bool
	// DataEncode applies encoding to data before writing to database.
	DataEncode([]Data) error
	// DataDecode applies decoding on data after reading from database.
	DataDecode([]Data) error

	// CreateSchema creates the table schema in database.
	CreateSchema(dbs *Session, tablename string) error
	// AlterSchema modifies the table schema in database.
	AlterSchema(dbs *Session, tablename string) error
	// InitialData creates the initial data in table.
	InitialData(dbs *Session, tablename string) error
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

// Orders returns the table order by columns to use in statments.
func (m *BaseModel) Orders() []string {
	return m.DefaultOrders
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

// CreateSchema creates the table schema in database.
func (m *BaseModel) CreateSchema(dbs *Session, tablename string) error {
	// if s.db.Engine == nil {
	// 	return ErrDBEngine
	// }

	if dbs.db.Log != nil {
		dbs.db.Log.Debug("creating schema for table: %s", tablename)
	}
	// schema := s.db.Engine.GenSchema(tablename, (*TableMeta)(m))

	// _, err := s.Execute(schema)
	// return err
	return nil
}

// AlterSchema modifies the table schema in database.
func (m *BaseModel) AlterSchema(dbs *Session, tablename string) error {
	return nil
}

// InitialData creates the initial data in table.
func (m *BaseModel) InitialData(dbs *Session, tablename string) error {
	return nil
}

////////////////////////////////////////////////////

// ModelsMeta represents link between table name and model metainfo.
type ModelsMeta struct {
	Table string
	Model Model
}

// InitializeModels creates and alter the database models schema,
// then adds the models intial data.
func InitializeModels(db *Database, metainfo []ModelsMeta) error {
	if db == nil {
		return ErrDBHandler
	} else if db.engine == nil {
		return ErrDBEngine
	}

	// create new session
	dbs := db.Session()

	// create and alter schema
	if db.Log != nil {
		db.Log.Debug("creating tables schema")
	}
	for _, v := range metainfo {
		if err := v.Model.CreateSchema(dbs, v.Table); err != nil {
			return err
		}
	}
	for _, v := range metainfo {
		if err := v.Model.AlterSchema(dbs, v.Table); err != nil {
			return err
		}
	}

	// add intial data to tables
	if db.Log != nil {
		db.Log.Debug("adding tables initial data")
	}
	for _, v := range metainfo {
		if err := v.Model.InitialData(dbs, v.Table); err != nil {
			return err
		}
	}

	return nil
}
