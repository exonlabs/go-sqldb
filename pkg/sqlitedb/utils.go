// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"errors"
	"os"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"
	"github.com/exonlabs/go-utils/pkg/abc/fsx"
	"github.com/exonlabs/go-utils/pkg/console"
)

func InteractiveConfig(defaults dictx.Dict) (dictx.Dict, error) {
	con, err := console.NewTermConsole()
	if err != nil {
		return nil, err
	}
	defer con.Close()

	cfg := dictx.Dict{}

	// set database path
	if val, err := con.Required().ReadValue("Enter database path",
		dictx.Fetch(defaults, "database", "")); err != nil {
		return nil, err
	} else {
		dictx.Set(cfg, "database", val)
	}

	return cfg, nil
}

func InteractiveSetup(cfg dictx.Dict) error {
	// check required params
	db_path := strings.TrimSpace(dictx.Fetch(cfg, "database", ""))
	if db_path == "" {
		return errors.New("empty database path")
	}

	if !fsx.IsExist(db_path) {
		file, err := os.OpenFile(
			db_path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}
