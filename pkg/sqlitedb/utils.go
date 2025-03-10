// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"os"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/abc/fsx"
	"github.com/exonlabs/go-utils/pkg/console"
)

// InteractiveConfig gets the database configuration interactively from console.
//
// The parsed options are:
//   - database: (string) the database file path
func InteractiveConfig(d dictx.Dict) (dictx.Dict, error) {
	con, err := console.NewTermConsole()
	if err != nil {
		return nil, err
	}
	defer con.Close()

	// get database path
	db_path, err := con.Required().ReadValue(
		"Enter database path", dictx.GetString(d, "database", ""))
	if err != nil {
		return nil, err
	}

	cfg, err := dictx.Clone(d)
	if err != nil {
		return nil, err
	}
	dictx.Merge(cfg, dictx.Dict{
		"database": db_path,
	})

	return cfg, nil
}

// InteractiveSetup performs an interactive console based database setup.
//
// The parsed options are:
//   - database: (string) the database file path
func InteractiveSetup(d dictx.Dict) error {
	cfg := &Config{}
	if err := cfg.InitConfig(d); err != nil {
		return err
	}

	// create database file if not exist
	if !fsx.IsExist(cfg.Database) {
		file, err := os.OpenFile(cfg.Database, os.O_CREATE|os.O_WRONLY, 0o664)
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}
