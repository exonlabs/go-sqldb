package sqldb

import "github.com/exonlabs/go-utils/pkg/types"

// example:
//
//	ColumnSchema{"col1", "VARCHAR(128) NOT NULL", "UNIQUE INDEX"}
//	ColumnSchema{"col2", "INTEGER", "INDEX"}
//	ColumnSchema{"col3", "TEXT", ""}
//	ColumnSchema{"col4", "BOOLEAN NOT NULL DEFAULT 0", ""}
type ColumnSchema struct {
	Name string
	// SQL column type name and definition
	Definition string
	// space seperated mix of: PRIMARY, UNIQUE, INDEX
	Constraint string
}

// example:
//
//	ConstraintSchema{"c1", "FOREIGN KEY ("col1") REFERENCES "table1" ("col1")"}
//	ConstraintSchema{"c2", "CHECK ("col2" IN (0,1,2))"}
type ConstraintSchema struct {
	Name       string
	Definition string
}

type TableSchema = struct {
	Options     Options
	Columns     []ColumnSchema
	Constraints []ConstraintSchema
}

type Data = types.Dict
type DataAdapter = func(any) (any, error)

type Model interface {
	TableName() string
	TableSchema() *TableSchema
	DefaultOrders() []string
	DataReaders() map[string]DataAdapter
	DataWriters() map[string]DataAdapter
	UpgradeSchema(*Session, ...string) error
	InitializeData(*Session, ...string) error
}

type BaseModel struct {
}

func (dbm *BaseModel) DefaultOrders() []string {
	return nil
}
func (dbm *BaseModel) DataReaders() map[string]DataAdapter {
	return nil
}
func (dbm *BaseModel) DataWriters() map[string]DataAdapter {
	return nil
}
func (dbm *BaseModel) UpgradeSchema(dbs *Session, tbl ...string) error {
	return nil
}
func (dbm *BaseModel) InitializeData(dbs *Session, tbl ...string) error {
	return nil
}
