// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mysqldb

import (
	"fmt"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/console"
)

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
		int64(dictx.GetUint(defaults, "port", 3306)))
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

func InteractiveSetup(cfg dictx.Dict) error {
	if err := PrepareConfig(cfg); err != nil {
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
	// DATABASE SETUP

	fmt.Println("database setup:")
	fmt.Printf("admin user: %s\n", adm_user)
	fmt.Printf("admin pass: %s\n", adm_pass)

	return nil
}
