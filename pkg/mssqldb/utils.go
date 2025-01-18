// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mssqldb

import (
	"errors"
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

	cfg := dictx.Dict{}

	// set database name
	if val, err := con.Required().ReadValue("Enter database name",
		dictx.GetString(defaults, "database", "")); err != nil {
		return nil, err
	} else {
		dictx.Set(cfg, "database", val)
	}

	// set database host
	if val, err := con.Required().ReadValue("Enter database host",
		dictx.GetString(defaults, "host", "localhost")); err != nil {
		return nil, err
	} else {
		dictx.Set(cfg, "host", val)
	}

	// set database port
	if val, err := con.Required().ReadNumber("Enter database port",
		int64(dictx.GetUint(defaults, "port", 1433))); err != nil {
		return nil, err
	} else {
		dictx.Set(cfg, "port", val)
	}

	// set database username
	if val, err := con.ReadValue("Enter database username",
		dictx.GetString(defaults, "username", "")); err != nil {
		return nil, err
	} else {
		dictx.Set(cfg, "username", val)
	}

	// set database password
	if val, err := con.Hidden().ReadValue("Enter database password",
		dictx.GetString(defaults, "password", "")); err != nil {
		return nil, err
	} else if val != "" {
		if err = con.Hidden().ConfirmValue(
			"Confirm database password", val); err != nil {
			return nil, err
		}
		dictx.Set(cfg, "password", val)
	}

	return cfg, nil
}

func InteractiveSetup(cfg dictx.Dict) error {
	// check required params
	db_name := dictx.GetString(cfg, "database", "")
	if db_name == "" {
		return errors.New("empty database name")
	}
	db_host := dictx.GetString(cfg, "host", "")
	if db_host == "" {
		return errors.New("empty database host")
	}
	db_port := dictx.GetUint(cfg, "port", 0)
	if db_port == 0 {
		return errors.New("empty database port")
	}
	db_user := dictx.GetString(cfg, "username", "")
	db_pass := dictx.GetString(cfg, "password", "")

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
	// DATABASE SETUP

	fmt.Println()
	fmt.Printf("database name: %s\n", db_name)
	fmt.Printf("database host: %s\n", db_host)
	fmt.Printf("database port: %d\n", db_port)
	fmt.Printf("database user: %s\n", db_user)
	fmt.Printf("database pass: %s\n", db_pass)
	fmt.Println()
	fmt.Printf("admin user: %s\n", adm_user)
	fmt.Printf("admin pass: %s\n", adm_pass)

	return nil
}
