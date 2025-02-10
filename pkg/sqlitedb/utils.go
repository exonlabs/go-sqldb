// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"os"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/abc/fsx"
	"github.com/exonlabs/go-utils/pkg/console"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
)

func PrepareConfig(cfg dictx.Dict) error {
	for _, k := range []string{"database"} {
		if dictx.IsExist(cfg, k) {
			s := strings.TrimSpace(dictx.GetString(cfg, k, ""))
			dictx.Set(cfg, k, s)
		}
	}

	if dictx.GetString(cfg, "database", "") == "" {
		return sqldb.ErrDBPath
	}
	return nil
}

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

func InteractiveSetup(cfg dictx.Dict) error {
	if err := PrepareConfig(cfg); err != nil {
		return err
	}

	dp_path := dictx.GetString(cfg, "database", "")
	if !fsx.IsExist(dp_path) {
		file, err := os.OpenFile(dp_path, os.O_CREATE|os.O_WRONLY, 0o664)
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}
