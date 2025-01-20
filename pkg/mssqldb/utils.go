// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package mssqldb

import (
	"fmt"
	"strings"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/console"
)

func PrepareConfig(cfg *sqldb.DBConfig) error {
	cfg.Database = strings.TrimSpace(cfg.Database)
	cfg.Host = strings.TrimSpace(cfg.Host)
	cfg.Username = strings.TrimSpace(cfg.Username)
	cfg.Password = strings.TrimSpace(cfg.Password)

	if cfg.Database == "" {
		return sqldb.ErrDBName
	}
	if cfg.Host == "" {
		return sqldb.ErrDBHost
	}
	if cfg.Port == 0 {
		return sqldb.ErrDBPort
	}
	return nil
}

func InteractiveConfig(defaults dictx.Dict) (*sqldb.DBConfig, error) {
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

	return &sqldb.DBConfig{
		Database: db_name,
		Host:     db_host,
		Port:     int(db_port),
		Username: db_user,
		Password: db_pass,
	}, nil
}

func InteractiveSetup(cfg *sqldb.DBConfig) error {
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
	fmt.Printf("database config:\n%s\n\n", cfg)
	fmt.Printf("admin user: %s\n", adm_user)
	fmt.Printf("admin pass: %s\n", adm_pass)

	return nil
}
