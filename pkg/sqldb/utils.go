// Copyright (c) 2024 ExonLabs, All rights reserved.
// Use of this source code is governed by a BSD 3-Clause
// license that can be found in the LICENSE file.

package sqldb

import (
	uuid "github.com/satori/go.uuid"
)

// NewGuid generates a new string guid in hex format.
func NewGuid() string {
	u := uuid.NewV5(uuid.NamespaceOID, uuid.NewV1().String())
	return u.String()
}

// FormatData applies data adaptors functions to data map.
// adaptor function matching is done by adaptor name and data keys.
// data values are modified in place in the original data map.
func FormatData(adaptors map[string]DataAdaptor, data ...Data) error {
	for i := range data {
		// loop adaptors and match adaptor name in data
		for name, fn := range adaptors {
			if oldVal, ok := data[i][name]; ok {
				newVal, err := fn(oldVal)
				if err != nil {
					return err
				}
				data[i][name] = newVal
			}
		}
	}
	return nil
}
