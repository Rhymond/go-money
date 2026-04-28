package money

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// ============================================================================
// database/sql
// ============================================================================

func TestMoney_Value(t *testing.T) {
	tcs := []struct {
		have      *Money
		separator string
		want      string
	}{
		{New(10, CAD), "|", "10|CAD"},
		{New(-10, USD), "+-+", "-10+-+USD"},
	}
	for _, tc := range tcs {
		DBMoneyValueSeparator = tc.separator
		got, err := tc.have.Value()
		if err != nil {
			t.Errorf("Value() error = %v", err)
			continue
		}
		if !reflect.DeepEqual(got, driver.Value(tc.want)) {
			t.Errorf("Expected %v got %v", tc.want, got)
		}
	}
	DBMoneyValueSeparator = DefaultDBMoneyValueSeparator
}

func TestMoney_Scan(t *testing.T) {
	tcs := []struct {
		src       interface{}
		separator string
		want      *Money
		wantErr   bool
	}{
		{src: "10|CAD", want: New(10, CAD)},
		{src: "20|USD", want: New(20, USD)},
		{src: "30000,IDR", separator: ",", want: New(30000, IDR)},
		{src: "10|", wantErr: true},
		{src: "|SAR", wantErr: true},
		{src: "10", wantErr: true},
		{src: "USD", wantErr: true},
		{src: "USD|10", wantErr: true},
		{src: "", wantErr: true},
		{src: "a|b|c", wantErr: true},
	}
	for _, tc := range tcs {
		if tc.separator != "" {
			DBMoneyValueSeparator = tc.separator
		} else {
			DBMoneyValueSeparator = DefaultDBMoneyValueSeparator
		}
		got := &Money{}
		err := got.Scan(tc.src)
		if (err != nil) != tc.wantErr {
			t.Errorf("Scan(%#v) error = %v, wantErr %v", tc.src, err, tc.wantErr)
			continue
		}
		if tc.wantErr {
			continue
		}
		eq, _ := tc.want.Equals(got)
		if !eq {
			t.Errorf("Expected %s %s got %s %s", tc.want.Display(), tc.want.Currency().Code, got.Display(), got.Currency().Code)
		}
	}
	DBMoneyValueSeparator = DefaultDBMoneyValueSeparator
}

func TestCurrency_Value(t *testing.T) {
	for code, cc := range currencies {
		got, err := cc.Value()
		if err != nil {
			t.Errorf("Value() error = %v", err)
			continue
		}
		if !reflect.DeepEqual(got, driver.Value(code)) {
			t.Errorf("Expected %v got %v", code, got)
		}
	}
}

func TestCurrency_Scan(t *testing.T) {
	for code, want := range currencies {
		got := &Currency{}
		if err := got.Scan(code); err != nil {
			t.Errorf("Scan(%s) error = %v", code, err)
			continue
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %#v got %#v", want, got)
		}
	}
}

// ============================================================================
// JSON
// ============================================================================

func TestMoney_MarshalJSON(t *testing.T) {
	tcs := []struct {
		given    *Money
		expected string
	}{
		{New(12345, IQD), `{"amount":12345,"currency":"IQD"}`},
		{&Money{}, `{"amount":0,"currency":""}`},
	}

	for _, tc := range tcs {
		b, err := json.Marshal(tc.given)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			continue
		}
		if string(b) != tc.expected {
			t.Errorf("Expected %s got %s", tc.expected, string(b))
		}
	}
}

func TestMoney_MarshalJSON_Custom(t *testing.T) {
	defer func() { MarshalJSON = defaultMarshalJSON }()

	given := New(12345, IQD)
	expected := `{"amount":12345,"currency_code":"IQD","currency_fraction":3}`
	MarshalJSON = func(m Money) ([]byte, error) {
		buff := bytes.NewBufferString(fmt.Sprintf(`{"amount": %d, "currency_code": "%s", "currency_fraction": %d}`, m.Amount(), m.Currency().Code, m.Currency().Fraction))
		return buff.Bytes(), nil
	}

	b, err := json.Marshal(given)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s got %s", expected, string(b))
	}
}

func TestMoney_UnmarshalJSON(t *testing.T) {
	tcs := []struct {
		given    string
		expected string
	}{
		{`{"amount": 10012, "currency":"USD"}`, "$100.12"},
	}

	for _, tc := range tcs {
		var m Money
		if err := json.Unmarshal([]byte(tc.given), &m); err != nil {
			t.Errorf("Unexpected error: %v", err)
			continue
		}
		if m.Display() != tc.expected {
			t.Errorf("Expected %s got %s", tc.expected, m.Display())
		}
	}
}

func TestMoney_UnmarshalJSON_Zero(t *testing.T) {
	tcs := []string{
		`{"amount": 0, "currency":""}`,
		`{}`,
	}

	for _, given := range tcs {
		var m Money
		if err := json.Unmarshal([]byte(given), &m); err != nil {
			t.Errorf("Unexpected error: %v", err)
			continue
		}
		if m != (Money{}) {
			t.Errorf("Expected zero value got %+v", m)
		}
	}
}

func TestMoney_UnmarshalJSON_Invalid(t *testing.T) {
	tcs := []string{
		`{"amount": "foo", "currency": "USD"}`,
		`{"amount": 1234, "currency": 1234}`,
	}

	for _, given := range tcs {
		var m Money
		err := json.Unmarshal([]byte(given), &m)
		if !errors.Is(err, ErrInvalidJSONUnmarshal) {
			t.Errorf("Expected ErrInvalidJSONUnmarshal got %+v", err)
		}
	}
}

func TestMoney_UnmarshalJSON_Custom(t *testing.T) {
	defer func() { UnmarshalJSON = defaultUnmarshalJSON }()

	given := `{"amount": 10012, "currency_code":"USD", "currency_fraction":2}`
	expected := "$100.12"
	UnmarshalJSON = func(m *Money, b []byte) error {
		data := make(map[string]interface{})
		if err := json.Unmarshal(b, &data); err != nil {
			return err
		}
		*m = *New(int64(data["amount"].(float64)), data["currency_code"].(string))
		return nil
	}

	var m Money
	if err := json.Unmarshal([]byte(given), &m); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if m.Display() != expected {
		t.Errorf("Expected %s got %s", expected, m.Display())
	}
}

// ============================================================================
// XML
// ============================================================================

func TestMoney_MarshalXML(t *testing.T) {
	tcs := []struct {
		given    *Money
		expected string
	}{
		{New(12345, IQD), `<Money><amount>12345</amount><currency>IQD</currency></Money>`},
		{&Money{}, `<Money><amount>0</amount><currency></currency></Money>`},
	}

	for _, tc := range tcs {
		b, err := xml.Marshal(tc.given)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			continue
		}
		if string(b) != tc.expected {
			t.Errorf("Expected %s got %s", tc.expected, string(b))
		}
	}
}

func TestMoney_UnmarshalXML(t *testing.T) {
	tcs := []struct {
		given    string
		expected string
	}{
		{`<Money><amount>10012</amount><currency>USD</currency></Money>`, "$100.12"},
	}

	for _, tc := range tcs {
		var m Money
		if err := xml.Unmarshal([]byte(tc.given), &m); err != nil {
			t.Errorf("Unexpected error: %v", err)
			continue
		}
		if m.Display() != tc.expected {
			t.Errorf("Expected %s got %s", tc.expected, m.Display())
		}
	}
}

func TestMoney_UnmarshalXML_Zero(t *testing.T) {
	given := `<Money><amount>0</amount><currency></currency></Money>`
	var m Money

	if err := xml.Unmarshal([]byte(given), &m); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if m != (Money{}) {
		t.Errorf("Expected zero value got %+v", m)
	}
}

func TestMoney_XML_Roundtrip(t *testing.T) {
	given := New(-12345, EUR)

	b, err := xml.Marshal(given)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var got Money
	if err := xml.Unmarshal(b, &got); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if got.Amount() != given.Amount() || got.Currency().Code != given.Currency().Code {
		t.Errorf("Expected %d %s got %d %s", given.Amount(), given.Currency().Code, got.Amount(), got.Currency().Code)
	}
}

