package persistence

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Roles is a custom type for handling JSONB roles in PostgreSQL
type Roles []string

// Value implements the driver.Valuer interface
func (r Roles) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}
	return json.Marshal(r)
}

// Scan implements the sql.Scanner interface
func (r *Roles) Scan(value interface{}) error {
	if value == nil {
		*r = Roles{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, &r)
}
