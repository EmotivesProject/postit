package model

import (
	"database/sql/driver"
	"encoding/json"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)

	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &j)
}

// https://gist.github.com/yanmhlv/d00aa61082d3b8d71bed
