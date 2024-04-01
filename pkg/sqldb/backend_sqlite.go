package sqldb

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/exonlabs/go-utils/pkg/unix/xterm"
)

const SQLITE_BACKEND = "sqlite"

type sqlite_engine struct {
	Database string
}

func SqliteEngine(opts Options) *sqlite_engine {
	return &sqlite_engine{}
}

func (dbe *sqlite_engine) BackendName() string {
	return SQLITE_BACKEND
}

func (dbe *sqlite_engine) GenSchema(
	tblname string, model Model) ([]string, error) {

	if tblname == "" {
		tblname = model.TableName()
	}

	meta := model.TableMeta()
	auto_guid := meta.Options.GetBool("sqlite_without_rowid", false)

	columns := meta.Columns
	// add guid column if not exist as first column
	if _, ok := model.(ModelSetAutoGuid); ok {
		if columns[0][0] != "guid" {
			columns = append([][]string{
				{"guid", "VARCHAR(32) NOT NULL", "PRIMARY"},
			}, columns...)
		}
		auto_guid = true
	}

	var expr, constraints, indexes []string

	for _, c := range columns {
		expr = append(expr, c[0]+" "+c[1])

		// add check constraint for bool datatype
		if strings.Contains(c[1], "BOOLEAN") {
			constraints = append(constraints,
				fmt.Sprintf("CHECK (%v IN (0,1))", c[0]))
		}

		// no column constraint
		if len(c) < 3 {
			continue
		}

		if strings.Contains(c[2], "PRIMARY") {
			// add primary_key constraint
			constraints = append(constraints,
				fmt.Sprintf("PRIMARY KEY (%v)", c[0]))
		} else if strings.Contains(c[2], "UNIQUE") &&
			!strings.Contains(c[2], "INDEX") {
			// add unique constraint if not indexed column
			constraints = append(constraints,
				fmt.Sprintf("UNIQUE (%v)", c[0]))
		}

		if strings.Contains(c[2], "PRIMARY") ||
			strings.Contains(c[2], "INDEX") {
			u := ""
			if strings.Contains(c[2], "PRIMARY") ||
				strings.Contains(c[2], "UNIQUE") {
				u = "UNIQUE "
			}
			indexes = append(indexes, fmt.Sprintf(
				"CREATE %vINDEX IF NOT EXISTS ix_%v_%v ON %v (%v);",
				u, tblname, c[0], tblname, c[0]))
		}
	}

	// add column constraints
	expr = append(expr, constraints...)
	// add explicit table constraints
	expr = append(expr, meta.Constraints...)

	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tblname)
	sql += "\n   " + strings.Join(expr, ",\n   ")
	if auto_guid {
		sql += "\n) WITHOUT ROWID;"
	} else {
		sql += "\n);"
	}

	result := []string{sql}
	result = append(result, indexes...)
	return result, nil
}

/////////////////////////////////////////////////////////

type sqlite_backend struct{}

func SqliteBackend() *sqlite_backend { return &sqlite_backend{} }

// interactive database configuration
func (*sqlite_backend) InteractiveConfig(opts Options) (Options, error) {
	con := xterm.NewConsole()

	if v, err := con.Required().ReadValue(
		"Enter database path",
		opts.GetString("database", "")); err != nil {
		return nil, err
	} else {
		opts.Set("database", v)
	}

	if v, err := con.ReadValue(
		"Enter connection extra args",
		opts.GetString("extra_args", "")); err != nil {
		return nil, err
	} else {
		opts.Set("extra_args", v)
	}

	return opts, nil
}

// interactive database setup
func (*sqlite_backend) InteractiveSetup(opts Options) error {
	database := opts.GetString("database", "")
	if database == "" {
		return fmt.Errorf("%w - invalid database path", ErrOperation)
	}

	if _, err := os.Stat(database); os.IsNotExist(err) {
		syscall.Umask(0)
		f, err := os.OpenFile(
			database, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o664)
		if err != nil {
			return fmt.Errorf("%w - %s", ErrOperation, err.Error())
		}
		defer f.Close()
	}
	return nil
}
