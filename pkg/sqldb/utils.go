// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	"fmt"
	"regexp"

	uuid "github.com/satori/go.uuid"
)

// NewGuid generates a new string guid in hex format.
func NewGuid() string {
	u := uuid.NewV5(uuid.NamespaceOID, string(uuid.NewV1().Bytes()))
	return fmt.Sprintf("%x", u.Bytes())
}

// SqlIdent checks for a valid SQL identifier string.
func SqlIdent(s string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", s)
	return matched
}
