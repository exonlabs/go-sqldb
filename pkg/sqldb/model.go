package sqldb

type TableMeta = struct {
	// custom options for special database uses
	Options Options

	// Table column definitions
	//	https://www.w3schools.com/sql/sql_create_table.asp
	// 	https://www.w3schools.com/sql/sql_datatypes.asp
	//
	// 	{column_name, datatype [, constraint]}
	//	constraint: is optional space seperated mix of {PRIMARY|UNIQUE|INDEX}.
	//	constraints can be added per column or globaly per table.
	//
	// Example:
	//
	//	{"col1", "VARCHAR(128) NOT NULL", "UNIQUE INDEX"}
	//	{"col2", "INTEGER", "INDEX"}
	//	{"col3", "TEXT"}
	//	{"col4", "BOOLEAN NOT NULL DEFAULT 0"}
	Columns [][]string

	// Table constraints definitions
	// 	https://www.w3schools.com/sql/sql_constraints.asp
	//
	// Example:
	//
	//	`PRIMARY KEY (col1)`
	//	`FOREIGN KEY (col1) REFERENCES table1 (col2) ON UPDATE CASCADE`
	//	`UNIQUE (col1,col2)`
	//	`CHECK (col1>=10 AND col2="val")`
	//
	// with Naming (for ALTER modifications):
	//
	//	`CONSTRAINT pk_name PRIMARY KEY (col1,col2)`
	//	`CONSTRAINT uc_name UNIQUE (col1)`
	//	`CONSTRAINT ck_name CHECK (col1 IN (0,1,2))`
	Constraints []string
}

type Model interface {
	TableName() string
	TableMeta() *TableMeta
}
type ModelAutoGuid interface{ AutoGuid() bool }
type ModelDefaultOrders interface{ DefaultOrders() []string }
type ModelDataReaders interface{ DataReaders() map[string]DataAdapter }
type ModelDataWriters interface{ DataWriters() map[string]DataAdapter }
type ModelAlterSchema interface{ AlterSchema(*Handler, string) error }
type ModelInitialData interface{ InitialData(*Handler, string) error }
