package money

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Scan is an implementation the database/sql scanner interface
func (c *Currency) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	data, ok := value.(string)
	if !ok {
		return errors.New("type assertion .(string) failed.")
	}
	*c = *newCurrency(data).get()
	return nil
}

// Value is an implementation of driver.Value
func (c *Currency) Value() (driver.Value, error) {
	return c.Code, nil
}

// Scan is an implementation the database/sql scanner interface
func (m *Money) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed.")
	}
	nm := &Money{}
	err := json.Unmarshal(data, &nm)
	if nm == nil || err != nil {
		return err
	}
	if nm.Currency_ == nil || nm.Currency_.Code == "" {
		return errors.New("invalid currency")
	}
	*m = *nm
	return nil
}

// Value is an implementation of driver.Value
func (m *Money) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}
