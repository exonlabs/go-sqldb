// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"os"
	"strings"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/abc/fsx"
	"github.com/exonlabs/go-utils/pkg/console"
)

func PrepareConfig(cfg *sqldb.DBConfig) error {
	cfg.Database = strings.TrimSpace(cfg.Database)

	if cfg.Database == "" {
		return sqldb.ErrDBPath
	}
	return nil
}

func InteractiveConfig(defaults dictx.Dict) (*sqldb.DBConfig, error) {
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

	return &sqldb.DBConfig{
		Database: db_path,
	}, nil
}

func InteractiveSetup(cfg *sqldb.DBConfig) error {
	if err := PrepareConfig(cfg); err != nil {
		return err
	}

	if !fsx.IsExist(cfg.Database) {
		file, err := os.OpenFile(
			cfg.Database, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o664)
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}
