package sqldb

import (
	"encoding/hex"

	uuid "github.com/satori/go.uuid"
)

// generate new GUID
func NewGuid() string {
	return hex.EncodeToString(
		uuid.NewV5(uuid.NewV1(), uuid.NewV4().String()).Bytes())
}

// apply data adapters to set of data, data is modified in place
func FormatData(adapters map[string]DataAdapter, data []Data) error {
	for key, fn := range adapters {
		// ignore adapters not in data list
		if !data[0].IsExist(key) {
			continue
		}
		// loop on data and modify
		for i := 0; i < len(data); i++ {
			val, err := fn(data[i].Get(key, nil))
			if err != nil {
				return err
			}
			data[i].Set(key, val)
		}
	}
	return nil
}
