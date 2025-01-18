// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mysqldb

// type mysql_engine struct {
// 	Database  string
// 	Host      string
// 	Port      int
// 	Username  string
// 	Password  string
// 	ExtraArgs string
// }

// func MysqlEngine(opts Options) *mysql_engine {
// 	return &mysql_engine{}
// }

// func (dbe *mysql_engine) BackendName() string {
// 	return MYSQL_BACKEND
// }

// func (dbe *mysql_engine) GenSchema(
// 	tblname string, model Model) ([]string, error) {

// 	return []string{}, nil
// }

// /////////////////////////////////////////////////////////

// type mysql_backend struct{}

// func MysqlBackend() *mysql_backend { return &mysql_backend{} }

// // interactive database configuration
// func (*mysql_backend) InteractiveConfig(opts Options) (Options, error) {
// 	// con := xterm.NewConsole()

// 	// if v, err := con.Required().ReadValue(
// 	// 	"Enter database name",
// 	// 	opts.GetString("database", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	opts.Set("database", v)
// 	// }

// 	// if v, err := con.Required().ReadValue(
// 	// 	"Enter database host IP/FQDN",
// 	// 	opts.GetString("host", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	opts.Set("host", v)
// 	// }

// 	// if v, err := con.Required().ReadNumberWLimit(
// 	// 	"Enter database port number",
// 	// 	opts.GetInt("port", 0), 0, 65536); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	opts.Set("port", v)
// 	// }

// 	// if v, err := con.ReadValue(
// 	// 	"Enter database access username",
// 	// 	opts.GetString("username", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	opts.Set("username", v)
// 	// }

// 	// if v, err := con.Hidden().ReadValue(
// 	// 	"Enter database access password",
// 	// 	opts.GetString("password", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	opts.Set("password", v)
// 	// }

// 	// if v, err := con.ReadValue(
// 	// 	"Enter connection extra args",
// 	// 	opts.GetString("extra_args", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	opts.Set("extra_args", v)
// 	// }

// 	return opts, nil
// }

// // interactive database setup
// func (*mysql_backend) InteractiveSetup(opts Options) error {

// 	///////////////////////////
// 	// TODO
// 	//////////////////////////

// 	return nil
// }
