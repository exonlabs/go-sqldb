// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"os"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/abc/fsx"
	"github.com/exonlabs/go-utils/pkg/console"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
)

// GetConfig creates configuration object from configuration dict.
// it checks and returns error if not all options have valid values.
//
// The parsed config options are:
//   - database: (string) the database file path - REQUIRED
func GetConfig(config dictx.Dict) (*sqldb.Config, error) {
	cfg := sqldb.NewConfig(config)

	// validations
	if cfg.Database == "" {
		return nil, sqldb.ErrDBPath
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

	// get database path
	db_path, err := con.Required().ReadValue("Enter database path",
		dictx.GetString(defaults, "database", ""))
	if err != nil {
		return nil, err
	}

	cfg, err := dictx.Clone(defaults)
	if err != nil {
		return nil, err
	}
	dictx.Merge(cfg, dictx.Dict{
		"database": db_path,
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
