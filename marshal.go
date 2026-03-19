package money

import (
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// MarshalBinary implements encoding.BinaryMarshaler.
// The format is: 8 bytes little-endian int64 amount followed by the ASCII
// currency code. This is also the encoding used when storing Money in MongoDB
// via the BSON binary type (mongo-driver calls MarshalBinary automatically).
//
//	b, _ := money.New(1234, "USD").MarshalBinary()
//	// b = [0xD2 0x04 0x00 0x00 0x00 0x00 0x00 0x00 0x55 0x53 0x44]
func (m Money) MarshalBinary() ([]byte, error) {
	code := m.currency.Code
	buf := make([]byte, 8+len(code))
	binary.LittleEndian.PutUint64(buf[:8], uint64(m.amount.val))
	copy(buf[8:], code)
	return buf, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
// It expects data produced by MarshalBinary.
func (m *Money) UnmarshalBinary(data []byte) error {
	if len(data) < 9 {
		return errors.New("go-money: binary data too short (need ≥ 9 bytes)")
	}
	m.amount = &Amount{val: int64(binary.LittleEndian.Uint64(data[:8]))}
	m.currency = newCurrency(string(data[8:])).get()
	return nil
}

// MarshalText implements encoding.TextMarshaler.
// The format is "<amount> <CURRENCY>" (the same format used by the SQL Valuer).
//
//	b, _ := money.New(1234, "USD").MarshalText() // "1234 USD"
func (m Money) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%d %s", m.amount.val, m.currency.Code)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It expects data produced by MarshalText.
func (m *Money) UnmarshalText(data []byte) error {
	parts := strings.SplitN(string(data), " ", 2)
	if len(parts) != 2 {
		return fmt.Errorf("go-money: invalid text %q (want \"<amount> <CURRENCY>\")", string(data))
	}
	n, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("go-money: invalid amount %q: %w", parts[0], err)
	}
	m.amount = &Amount{val: n}
	m.currency = newCurrency(strings.TrimSpace(parts[1])).get()
	return nil
}

// moneyXMLDoc is the on-wire XML representation of a Money value.
type moneyXMLDoc struct {
	Amount   int64  `xml:"amount"`
	Currency string `xml:"currency"`
}

// MarshalXML implements xml.Marshaler.
// Each Money value is encoded as:
//
//	<Money><amount>1234</amount><currency>USD</currency></Money>
func (m Money) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(moneyXMLDoc{Amount: m.amount.val, Currency: m.currency.Code}, start)
}

// UnmarshalXML implements xml.Unmarshaler.
// It expects the format produced by MarshalXML.
func (m *Money) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc moneyXMLDoc
	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}
	m.amount = &Amount{val: doc.Amount}
	m.currency = newCurrency(doc.Currency).get()
	return nil
}
