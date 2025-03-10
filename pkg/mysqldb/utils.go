// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mysqldb

import (
	"fmt"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/console"
)

// InteractiveConfig gets the database configuration interactively from console.
//
// The parsed options are:
//   - database: (string) the database name - REQUIRED
//   - host: (string) the database server IP or FQDN - REQUIRED
//   - port: (int) the database server port number - REQUIRED
//   - username: (string) the database  access username (if any)
//   - password: (string) the database access password (if any)
func InteractiveConfig(d dictx.Dict) (dictx.Dict, error) {
	con, err := console.NewTermConsole()
	if err != nil {
		return nil, err
	}
	defer con.Close()

	// get database name
	db_name, err := con.Required().ReadValue(
		"Enter database name", dictx.GetString(d, "database", ""))
	if err != nil {
		return nil, err
	}

	// get database host
	db_host, err := con.Required().ReadValue(
		"Enter database host", dictx.GetString(d, "host", "localhost"))
	if err != nil {
		return nil, err
	}

	// get database port
	db_port, err := con.Required().ReadNumber(
		"Enter database port", int64(dictx.GetUint(d, "port", 3306)))
	if err != nil {
		return nil, err
	}

	// get database username
	db_user, err := con.ReadValue(
		"Enter database username", dictx.GetString(d, "username", ""))
	if err != nil {
		return nil, err
	}

	// get database password
	db_pass, err := con.Hidden().ReadValue(
		"Enter database password", dictx.GetString(d, "password", ""))
	if err != nil {
		return nil, err
	} else if db_pass != "" {
		if err = con.Hidden().ConfirmValue(
			"Confirm database password", db_pass); err != nil {
			return nil, err
		}
	}

	cfg, err := dictx.Clone(d)
	if err != nil {
		return nil, err
	}
	dictx.Merge(cfg, dictx.Dict{
		"database": db_name,
		"host":     db_host,
		"port":     int(db_port),
		"username": db_user,
		"password": db_pass,
	})

	return cfg, nil
}

// InteractiveSetup performs an interactive console based database setup.
//
// The parsed options are:
//   - database: (string) the database name - REQUIRED
//   - host: (string) the database server IP or FQDN - REQUIRED
//   - port: (int) the database server port number - REQUIRED
//   - username: (string) the database  access username (if any)
//   - password: (string) the database access password (if any)
func InteractiveSetup(d dictx.Dict) error {
	cfg := &Config{}
	if err := cfg.InitConfig(d); err != nil {
		return err
	}

	con, err := console.NewTermConsole()
	if err != nil {
		return err
	}
	defer con.Close()

	// get database admin username and password
	adm_user, err := con.Required().ReadValue(
		"Enter database admin username", "admin")
	if err != nil {
		return err
	}
	adm_pass, err := con.Hidden().ReadValue(
		"Enter database admin password", "")
	if err != nil {
		return err
	}

	// TODO:
	fmt.Println()
	fmt.Println("-- TODO -----------------")
	fmt.Println("database setup:")
	fmt.Printf("admin username: %s\n", adm_user)
	fmt.Printf("admin password: %s\n", adm_pass)
	fmt.Printf("database config: %s\n", cfg)
	fmt.Println("-------------------------")

	return nil
}
