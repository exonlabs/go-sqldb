package sqldb

import "github.com/exonlabs/go-utils/pkg/unix/xterm"

const MSSQL_BACKEND = "mssql"

type mssql_engine struct{}

func MssqlEngine(opts Options) *mssql_engine {
	return &mssql_engine{}
}

func (dbe *mssql_engine) BackendName() string {
	return MSSQL_BACKEND
}

func (dbe *mssql_engine) GenSchema(
	tblname string, model Model) ([]string, error) {

	return []string{}, nil
}

/////////////////////////////////////////////////////////

type mssql_backend struct{}

func MssqlBackend() *mssql_backend { return &mssql_backend{} }

// interactive database configuration
func (*mssql_backend) InteractiveConfig(opts Options) (Options, error) {
	con := xterm.NewConsole()

	if v, err := con.Required().ReadValue(
		"Enter database name",
		opts.GetString("database", "")); err != nil {
		return nil, err
	} else {
		opts.Set("database", v)
	}

	if v, err := con.Required().ReadValue(
		"Enter database host IP/FQDN",
		opts.GetString("host", "")); err != nil {
		return nil, err
	} else {
		opts.Set("host", v)
	}

	if v, err := con.Required().ReadNumberWLimit(
		"Enter database port number",
		opts.GetInt("port", 0), 0, 65536); err != nil {
		return nil, err
	} else {
		opts.Set("port", v)
	}

	if v, err := con.ReadValue(
		"Enter database access username",
		opts.GetString("username", "")); err != nil {
		return nil, err
	} else {
		opts.Set("username", v)
	}

	if v, err := con.Hidden().ReadValue(
		"Enter database access password",
		opts.GetString("password", "")); err != nil {
		return nil, err
	} else {
		opts.Set("password", v)
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
func (*mssql_backend) InteractiveSetup(opts Options) error {

	///////////////////////////
	// TODO
	//////////////////////////

	return nil
}
