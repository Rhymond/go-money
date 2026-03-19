package money

import (
	"database/sql/driver"
	"encoding/json"
)

// NullMoney represents a Money that may be NULL in a database column or
// absent in a JSON payload. It mirrors the pattern of sql.NullString.
//
// Use Valid to distinguish a zero-value Money from a missing value.
type NullMoney struct {
	Money Money
	Valid bool // Valid is true when Money is not NULL / not absent.
}

// NewNullMoney returns a valid NullMoney wrapping m.
func NewNullMoney(m *Money) NullMoney {
	return NullMoney{Money: *m, Valid: true}
}

// Value implements driver.Valuer for SQL writes.
// Returns nil when Valid is false (stored as SQL NULL).
func (n NullMoney) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Money.Value()
}

// Scan implements sql.Scanner for SQL reads.
// Sets Valid=false when the source value is nil (SQL NULL).
func (n *NullMoney) Scan(value interface{}) error {
	if value == nil {
		n.Money = Money{}
		n.Valid = false
		return nil
	}
	n.Valid = true
	return n.Money.Scan(value)
}

// MarshalJSON implements json.Marshaler.
// Encodes as JSON null when Valid is false.
func (n NullMoney) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Money)
}

// UnmarshalJSON implements json.Unmarshaler.
// Sets Valid=false for JSON null; otherwise decodes normally.
func (n *NullMoney) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(data, &n.Money)
}

// MarshalText implements encoding.TextMarshaler.
// Returns an empty slice when Valid is false.
func (n NullMoney) MarshalText() ([]byte, error) {
	if !n.Valid {
		return []byte{}, nil
	}
	return n.Money.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
// Sets Valid=false for an empty input; otherwise decodes normally.
func (n *NullMoney) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return n.Money.UnmarshalText(data)
}
