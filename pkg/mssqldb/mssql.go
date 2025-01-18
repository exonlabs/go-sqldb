// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mssqldb

// const MSSQL_BACKEND = "mssql"

// type mssql_engine struct{}

// func MssqlEngine(opts Options) *mssql_engine {
// 	return &mssql_engine{}
// }

// func (dbe *mssql_engine) BackendName() string {
// 	return MSSQL_BACKEND
// }

// func (dbe *mssql_engine) GenSchema(
// 	tblname string, model Model) ([]string, error) {

// 	return []string{}, nil
// }

// /////////////////////////////////////////////////////////

// type mssql_backend struct{}

// func MssqlBackend() *mssql_backend { return &mssql_backend{} }

// // interactive database configuration
// func (*mssql_backend) InteractiveConfig(opts Options) (Options, error) {
// 	// con := xterm.NewConsole()

// 	// if v, err := con.Required().ReadValue(
// 	// 	"Enter database name",
// 	// 	dictx.Fetch(opts, "database", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	dictx.Set(opts, "database", v)
// 	// }

// 	// if v, err := con.Required().ReadValue(
// 	// 	"Enter database host IP/FQDN",
// 	// 	dictx.Fetch(opts, "host", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	dictx.Set(opts, "host", v)
// 	// }

// 	// if v, err := con.Required().ReadNumberWLimit(
// 	// 	"Enter database port number",
// 	// 	dictx.GetInt(opts, "port", 0), 0, 65536); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	dictx.Set(opts, "port", v)
// 	// }

// 	// if v, err := con.ReadValue(
// 	// 	"Enter database access username",
// 	// 	dictx.Fetch(opts, "username", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	dictx.Set(opts, "username", v)
// 	// }

// 	// if v, err := con.Hidden().ReadValue(
// 	// 	"Enter database access password",
// 	// 	dictx.Fetch(opts, "password", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	dictx.Set(opts, "password", v)
// 	// }

// 	// if v, err := con.ReadValue(
// 	// 	"Enter connection extra args",
// 	// 	dictx.Fetch(opts, "extra_args", "")); err != nil {
// 	// 	return nil, err
// 	// } else {
// 	// 	dictx.Set(opts, "extra_args", v)
// 	// }

// 	return opts, nil
// }

// // interactive database setup
// func (*mssql_backend) InteractiveSetup(opts Options) error {

// 	///////////////////////////
// 	// TODO
// 	//////////////////////////

// 	return nil
// }
