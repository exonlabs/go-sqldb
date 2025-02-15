// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mssqldb

import (
	"fmt"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/console"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
)

// GetConfig creates configuration object from configuration dict.
// it checks and returns error if not all options have valid values.
//
// The parsed config options are:
//   - database: (string) the database name - REQUIRED
//   - host: (string) the database server IP or FQDN - REQUIRED
//   - port: (int) the database server port number - REQUIRED
//   - username: (string) the database  access username (if any)
//   - password: (string) the database access password (if any)
//   - args: (string) the database extra connection params.
func GetConfig(config dictx.Dict) (*sqldb.Config, error) {
	cfg := sqldb.NewConfig(config)

	// validations
	if cfg.Database == "" {
		return nil, sqldb.ErrDBName
	}
	if cfg.Host == "" {
		return nil, sqldb.ErrDBHost
	}
	if cfg.Port == 0 {
		return nil, sqldb.ErrDBPort
	}

	return cfg, nil
}

// InteractiveConfig gets the database configuration interactively from console.
// The database default options are detailed in [GetConfig]
func InteractiveConfig(defaults dictx.Dict) (dictx.Dict, error) {
	con, err := console.NewTermConsole()
	if err != nil {
		return nil, err
	}
	defer con.Close()

	// get database name
	db_name, err := con.Required().ReadValue("Enter database name",
		dictx.GetString(defaults, "database", ""))
	if err != nil {
		return nil, err
	}

	// get database host
	db_host, err := con.Required().ReadValue("Enter database host",
		dictx.GetString(defaults, "host", "localhost"))
	if err != nil {
		return nil, err
	}

	// get database port
	db_port, err := con.Required().ReadNumber("Enter database port",
		int64(dictx.GetUint(defaults, "port", 1433)))
	if err != nil {
		return nil, err
	}

	// get database username
	db_user, err := con.ReadValue("Enter database username",
		dictx.GetString(defaults, "username", ""))
	if err != nil {
		return nil, err
	}

	// get database password
	db_pass, err := con.Hidden().ReadValue("Enter database password",
		dictx.GetString(defaults, "password", ""))
	if err != nil {
		return nil, err
	} else if db_pass != "" {
		if err = con.Hidden().ConfirmValue(
			"Confirm database password", db_pass); err != nil {
			return nil, err
		}
	}

	cfg, err := dictx.Clone(defaults)
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
// The database config options are detailed in [GetConfig]
func InteractiveSetup(config dictx.Dict) error {
	cfg, err := GetConfig(config)
	if err != nil {
		return err
	}

	con, err := console.NewTermConsole()
	if err != nil {
		return err
	}
	defer con.Close()

	// get database admin username and password
	adm_user, err := con.Required().ReadValue(
		"Enter database admin username", "sa")
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
