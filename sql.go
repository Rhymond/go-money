package money

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// Value implements the driver.Valuer interface, allowing Money to be stored
// directly in any database/sql-compatible database.
//
// The stored format is "<amount> <CURRENCY>" where amount is the integer value
// in the currency's smallest unit.
//
//	money.New(1234, "USD").Value() // "1234 USD"
//	money.New(-50, "EUR").Value()  // "-50 EUR"
func (m Money) Value() (driver.Value, error) {
	return fmt.Sprintf("%d %s", m.amount.val, m.currency.Code), nil
}

// Scan implements the sql.Scanner interface, allowing Money to be populated
// directly from a database/sql query result.
//
// It accepts string or []byte values in the format "<amount> <CURRENCY>"
// (e.g. "1234 USD") as produced by Value().
//
//	var m money.Money
//	row.Scan(&m)
func (m *Money) Scan(src interface{}) error {
	var str string
	switch v := src.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("go-money: unsupported Scan source type %T (expected string)", src)
	}

	parts := strings.SplitN(str, " ", 2)
	if len(parts) != 2 {
		return fmt.Errorf(
			"go-money: invalid money string %q — expected format \"<amount> <CURRENCY>\"",
			str,
		)
	}

	amount, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("go-money: invalid amount %q: %w", parts[0], err)
	}

	m.amount = &Amount{val: amount}
	m.currency = newCurrency(parts[1]).get()
	return nil
}
