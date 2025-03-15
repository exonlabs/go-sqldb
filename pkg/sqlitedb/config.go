// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqlitedb

import (
	"fmt"
	"strings"

	"github.com/exonlabs/go-utils/pkg/abc/dictx"

	"github.com/exonlabs/go-sqldb/pkg/sqldb"
)

// Config represents the database configuration params.
type Config struct {
	// database path
	Database string
	// disable Foreign Keys restriction
	NoForeignKeys bool
	// disable Auto Vacuum mode
	NoAutoVacuum bool
	// disable WAL Journal mode
	NoJournalWAL bool
}

// NewConfig creates configuration object from options dict.
// it checks and returns error if not all options have valid values.
//
// The parsed options are:
//   - database: (string) the database file path - REQUIRED
//   - no_foreign_keys: (bool) disable Foreign Keys restriction
//   - no_auto_vacuum: (bool) disable Auto Vacuum mode
//   - no_journal_wal: (bool) disable WAL Journal mode
func NewConfig(opts dictx.Dict) (*Config, error) {
	cfg := &Config{
		Database:      dictx.GetString(opts, "database", ""),
		NoForeignKeys: dictx.Fetch(opts, "no_foreign_keys", false),
		NoAutoVacuum:  dictx.Fetch(opts, "no_auto_vacuum", false),
		NoJournalWAL:  dictx.Fetch(opts, "no_journal_wal", false),
	}

	// validations
	if cfg.Database == "" {
		return nil, sqldb.ErrDBPath
	}

	return cfg, nil
}

// DSN returns the driver-specific data source name.
func (cfg *Config) DSN() string {
	args := []string{}

	if cfg.NoForeignKeys {
		args = append(args, "_foreign_keys=0")
	} else {
		args = append(args, "_foreign_keys=1")
	}

	if !cfg.NoAutoVacuum {
		args = append(args, "_auto_vacuum=1")
	}

	// if !cfg.NoJournalWAL {
	// 	args = append(args, "_journal_mode=WAL")
	// }

	return fmt.Sprintf("%s?%s", cfg.Database, strings.Join(args, "&"))
}
