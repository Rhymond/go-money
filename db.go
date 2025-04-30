package money

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

var (
	// DBMoneyValueSeparator is used to join together the Amount and Currency components of money.Money instances
	// allowing them to be stored as strings (via the driver.Valuer interface) and unmarshalled as strings (via
	// the sql.Scanner interface); set this value to use a different separator.
	DBMoneyValueSeparator = DefaultDBMoneyValueSeparator
)

const (
	// DefaultDBMoneyValueSeparator is the default value for DBMoneyValueSeparator; can be used to reset the
	// active separator value
	DefaultDBMoneyValueSeparator = "|"
)

// Value implements driver.Valuer to serialise a Money instance into a delimited string using the DBMoneyValueSeparator
// for example: "amount|currency_code"
func (m *Money) Value() (driver.Value, error) {
	return fmt.Sprintf("%d%s%s", m.amount, DBMoneyValueSeparator, m.Currency().Code), nil
}

// Scan implements sql.Scanner to deserialize a Money instance from a DBMoneyValueSeparator-separated string
// for example: "amount|currency_code"
func (m *Money) Scan(src interface{}) error {
	var amount Amount
	currency := &Currency{}

	// let's support string and int64
	switch src.(type) {
	case string:
		parts := strings.Split(src.(string), DBMoneyValueSeparator)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("%#v is not valid to scan into Money; update your query to return a money.DBMoneyValueSeparator-separated pair of \"amount%scurrency_code\"", src.(string), DBMoneyValueSeparator)
		}

		if a, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
			amount = a
		} else {
			return fmt.Errorf("scanning %#v into an Amount: %v", parts[0], err)
		}

		if err := currency.Scan(parts[1]); err != nil {
			return fmt.Errorf("scanning %#v into a Currency: %v", parts[1], err)
		}
	default:
		return fmt.Errorf("don't know how to scan %T into Money; update your query to return a money.DBMoneyValueSeparator-separated pair of \"amount%scurrency_code\"", src, DBMoneyValueSeparator)
	}

	// allocate new Money with the scanned amount and currency
	*m = Money{
		amount:   amount,
		currency: currency,
	}

	return nil
}

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
		return fmt.Errorf("GetCurrency(%#v) returned nil", src)
	}

	// copy the value
	*c = *val

	return nil
}
