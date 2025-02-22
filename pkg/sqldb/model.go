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
	// ex. "VARCHAR(128) NOT NULL", "BOOLEAN NOT NULL DEFAULT false"
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
	// AutoGuid sets weather to enable AutoGuid operations, which is to
	// create a first primary guid column for table.
	// guid column is created with schema "guid VARCHAR(32) NOT NULL"
	AutoGuid bool
	// Extra options for specific backends. each option is prefixed with the
	// backend name plus underscore. ex: "sqlite_<OPTION_NAME>"
	Args dictx.Dict
}

// Model represents the model interface.
type Model interface {
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
// does nothing by default.
func (m *BaseModel) DataEncode([]Data) error {
	return nil
}

// DataDecode applies decoding on data after reading from database.
// does nothing by default.
func (m *BaseModel) DataDecode([]Data) error {
	return nil
}

// ModelMeta represents the model metainfo interface.
type ModelMeta interface {
	// CreateSchema creates the table schema in database.
	CreateSchema(s *Session, tablename string) error
	// AlterSchema modifies the table schema in database.
	AlterSchema(s *Session, tablename string) error
	// InitialData creates the initial data in table.
	InitialData(s *Session, tablename string) error
}

// TableModelMeta represents link between table name and model metainfo.
type TableModelMeta struct {
	TableName string
	ModelMeta ModelMeta
}

// BaseModelMeta defines a base model metainfo structure.
type BaseModelMeta TableMeta

// CreateSchema creates the table schema in database.
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

// AlterSchema modifies the table schema in database. default nothing.
func (m *BaseModelMeta) AlterSchema(s *Session, tablename string) error {
	return nil
}

// InitialData creates the initial data in table. default none.
func (m *BaseModelMeta) InitialData(s *Session, tablename string) error {
	return nil
}
