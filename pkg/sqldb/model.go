// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"github.com/exonlabs/go-utils/pkg/abc/dictx"
)

// ColumnMeta defines strcture holding column definitions and constraints.
//
// References:
//   - https://www.w3schools.com/sql/sql_create_table.asp
//   - https://www.w3schools.com/sql/sql_datatypes.asp
type ColumnMeta struct {
	// the column name, should be unique per table.
	Name string
	// the column data type as defined in SQL syntax.
	// ex. "INTEGER(10)", "VARCHAR(128) NOT NULL", "BOOLEAN NOT NULL DEFAULT 0"
	Type string
	// set column primary key constraint.
	Primary bool
	// set column unique value constraint.
	Unique bool
	// set to create column index.
	Index bool
}

// ConstraintMeta defines strcture holding constraints definitions.
//
// References:
//   - https://www.w3schools.com/sql/sql_constraints.asp
type ConstraintMeta struct {
	// the constraint name, should be unique per table.
	Name string
	// the constraint definition as defined in SQL syntax.
	// ex. "PRIMARY KEY (col1,col2)"
	//     "FOREIGN KEY (col1) REFERENCES table1 (col2) ON UPDATE CASCADE"
	//     "UNIQUE (col1,col2)"
	//     "CHECK (col1 IN (0,1,2))"
	//     "CHECK (col1>=10 AND col2="val")"
	Definition string
}

// TableMeta defines strcture holding table definitions and constraints.
//
// References:
//   - https://www.w3schools.com/sql/sql_create_table.asp
//   - https://www.w3schools.com/sql/sql_datatypes.asp
//   - https://www.w3schools.com/sql/sql_constraints.asp
type TableMeta struct {
	// Table Columns meta
	Columns []ColumnMeta
	// Table Constraints as defined in SQL syntax. constraints are appended to
	// table after auto generated columns constraints.
	Constraints []ConstraintMeta
	// AutoGuid sets weather to create a guid primary column for table.
	// guid column is created with name "guid" and type "VARCHAR(32) NOT NULL"
	AutoGuid bool
	// Extra options for specific backends. each option is prefixed with the
	// backend name plus underscore. ex: "sqlite_<OPTION_NAME>"
	Args dictx.Dict
}

// Model represents the model interface.
type Model interface {
	TableName() string
	Columns() []string
	Orders() []string
	IsAutoGuid() bool
	DataReaders() map[string]DataAdaptor
	DataWriters() map[string]DataAdaptor
}

type BaseModel struct {
	DefaultTable   string
	DefaultColumns []string
	DefaultOrders  []string
	AutoGuid       bool
	Readers        map[string]DataAdaptor
	Writers        map[string]DataAdaptor
}

func (m *BaseModel) TableName() string {
	return m.DefaultTable
}

func (m *BaseModel) Columns() []string {
	return m.DefaultColumns
}

func (m *BaseModel) Orders() []string {
	return m.DefaultOrders
}

func (m *BaseModel) IsAutoGuid() bool {
	return m.AutoGuid
}

func (m *BaseModel) DataReaders() map[string]DataAdaptor {
	return nil
}

func (m *BaseModel) DataWriters() map[string]DataAdaptor {
	return nil
}

// ModelMeta represents the model metainfo interface.
type ModelMeta interface {
	CreateSchema(s *Session, tablename string) error
	AlterSchema(s *Session, tablename string) error
	InitialData(s *Session, tablename string) error
}

type BaseModelMeta TableMeta

func (m *BaseModelMeta) CreateSchema(s *Session, tablename string) error {
	if s.db.engine == nil {
		return ErrDBEngine
	}

	if s.db.DBLog != nil {
		s.db.DBLog.Debug("creating schema for table: %s", tablename)
	}
	schema := s.db.engine.GenSchema(tablename, (*TableMeta)(m))

	_, err := s.Execute(schema)
	return err
}

func (m *BaseModelMeta) AlterSchema(s *Session, tablename string) error {
	return nil
}

func (m *BaseModelMeta) InitialData(s *Session, tablename string) error {
	return nil
}
