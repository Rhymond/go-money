package money

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// ErrInvalidJSONUnmarshal happens when the default money.UnmarshalJSON fails
// to unmarshal Money because of invalid data.
var ErrInvalidJSONUnmarshal = errors.New("invalid json unmarshal")

// Injection points for backward compatibility. Overwrite to keep your own
// JSON / XML shape:
//
//	money.UnmarshalJSON = func(m *Money, b []byte) error { ... }
//	money.MarshalJSON   = func(m Money) ([]byte, error) { ... }
//
// These globals are not safe to mutate concurrently with use — set them once
// during package initialization.
var (
	UnmarshalJSON = defaultUnmarshalJSON
	MarshalJSON   = defaultMarshalJSON

	UnmarshalXML = defaultUnmarshalXML
	MarshalXML   = defaultMarshalXML
)

// ============================================================================
// database/sql
// ============================================================================

const DefaultDBMoneyValueSeparator = "|"

// DBMoneyValueSeparator joins amount and currency when storing Money via
// driver.Valuer / sql.Scanner — e.g. "amount|currency_code".
//
// Not safe to mutate concurrently with use — set once during package init.
var DBMoneyValueSeparator = DefaultDBMoneyValueSeparator

// Value implements driver.Valuer.
func (m Money) Value() (driver.Value, error) {
	return fmt.Sprintf("%d%s%s", m.Amount(), DBMoneyValueSeparator, m.Currency().Code), nil
}

// Scan implements sql.Scanner. NULL values produce an error; wrap Money in
// your own nullable type if a column is nullable. Accepts both string and
// []byte sources — drivers vary on which they hand back for text columns.
func (m *Money) Scan(src interface{}) error {
	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("don't know how to scan %T into Money; update your query to return a money.DBMoneyValueSeparator-separated pair of \"amount%scurrency_code\"", src, DBMoneyValueSeparator)
	}

	parts := strings.Split(s, DBMoneyValueSeparator)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("%#v is not valid to scan into Money; update your query to return a money.DBMoneyValueSeparator-separated pair of \"amount%scurrency_code\"", s, DBMoneyValueSeparator)
	}

	amount, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("scanning %#v into an amount: %v", parts[0], err)
	}

	*m = Money{
		amount:   &Decimal{val: big.NewInt(amount)},
		currency: newCurrency(parts[1]).get(),
	}
	return nil
}

// Value implements driver.Valuer for Currency.
func (c Currency) Value() (driver.Value, error) {
	return c.Code, nil
}

// Scan implements sql.Scanner for Currency. Unknown codes are accepted with
// default formatting, matching the lenient behavior of New(). Accepts both
// string and []byte sources.
func (c *Currency) Scan(src interface{}) error {
	var code string
	switch v := src.(type) {
	case string:
		code = v
	case []byte:
		code = string(v)
	default:
		return fmt.Errorf("%T is not a supported type for a Currency (store the Currency.Code value as a string only)", src)
	}
	*c = *newCurrency(code).get()
	return nil
}

// ============================================================================
// JSON
// ============================================================================

type jsonMoney struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// MarshalJSON implements json.Marshaler.
func (m Money) MarshalJSON() ([]byte, error) {
	return MarshalJSON(m)
}

// UnmarshalJSON implements json.Unmarshaler.
func (m *Money) UnmarshalJSON(b []byte) error {
	return UnmarshalJSON(m, b)
}

func defaultMarshalJSON(m Money) ([]byte, error) {
	sm := m.safe()
	return json.Marshal(jsonMoney{
		Amount:   sm.Amount(),
		Currency: sm.currency.Code,
	})
}

func defaultUnmarshalJSON(m *Money, b []byte) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()

	data := make(map[string]interface{})
	if err := dec.Decode(&data); err != nil {
		return err
	}

	var amount int64
	if amountRaw, ok := data["amount"]; ok {
		n, ok := amountRaw.(json.Number)
		if !ok {
			return ErrInvalidJSONUnmarshal
		}
		v, err := n.Int64()
		if err != nil {
			return ErrInvalidJSONUnmarshal
		}
		amount = v
	}

	var currency string
	if currencyRaw, ok := data["currency"]; ok {
		currency, ok = currencyRaw.(string)
		if !ok {
			return ErrInvalidJSONUnmarshal
		}
	}

	*m = *New(amount, currency)
	return nil
}

// ============================================================================
// XML
// ============================================================================

type xmlMoney struct {
	Amount   int64  `xml:"amount"`
	Currency string `xml:"currency"`
}

// MarshalXML implements xml.Marshaler.
func (m Money) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return MarshalXML(m, e, start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (m *Money) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return UnmarshalXML(m, d, start)
}

func defaultMarshalXML(m Money, e *xml.Encoder, start xml.StartElement) error {
	sm := m.safe()
	return e.EncodeElement(xmlMoney{
		Amount:   sm.Amount(),
		Currency: sm.currency.Code,
	}, start)
}

func defaultUnmarshalXML(m *Money, d *xml.Decoder, start xml.StartElement) error {
	var aux xmlMoney
	if err := d.DecodeElement(&aux, &start); err != nil {
		return err
	}
	*m = *New(aux.Amount, aux.Currency)
	return nil
}

