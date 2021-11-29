package money

import (
	"database/sql/driver"
	"fmt"
)

// Value implements driver.Valuer to serialize a Currency code into a string for saving to a database
func (c Currency) Value() (driver.Value, error) {
	return c.Code, nil
}

// Scan implements sql.Scanner to deserialize a Currency from a string value read from a database
func (c *Currency) Scan(src interface{}) error {
	var val *Currency
	// let's support string only
	switch src.(type) {
	case string:
		val = GetCurrency(src.(string))
	default:
		return fmt.Errorf("%T is not a supported type for a Currency (store the Currency.Code value as a string only)", src)
	}

	if val == nil {
		return fmt.Errorf("something went wrong; getCurrency(%#v) returned nil", src)
	}

	// copy the value
	*c = *val

	return nil
}
