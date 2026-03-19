package money

import (
	"database/sql/driver"
	"testing"
)

func TestMoney_Value(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected string
	}{
		{1234, "USD", "1234 USD"},
		{-50, "EUR", "-50 EUR"},
		{0, "GBP", "0 GBP"},
		{999999, "JPY", "999999 JPY"},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		v, err := m.Value()
		if err != nil {
			t.Fatalf("Value(): unexpected error: %v", err)
		}
		str, ok := v.(string)
		if !ok {
			t.Fatalf("Value(): expected string, got %T", v)
		}
		if str != tc.expected {
			t.Errorf("Value(%d %s): expected %q got %q", tc.amount, tc.code, tc.expected, str)
		}
	}
}

func TestMoney_Scan_String(t *testing.T) {
	tcs := []struct {
		input    string
		amount   int64
		currency string
	}{
		{"1234 USD", 1234, "USD"},
		{"-50 EUR", -50, "EUR"},
		{"0 GBP", 0, "GBP"},
		{"999999 JPY", 999999, "JPY"},
	}

	for _, tc := range tcs {
		var m Money
		if err := m.Scan(tc.input); err != nil {
			t.Fatalf("Scan(%q): unexpected error: %v", tc.input, err)
		}
		if m.Amount() != tc.amount {
			t.Errorf("Scan(%q): expected amount %d got %d", tc.input, tc.amount, m.Amount())
		}
		if m.Currency().Code != tc.currency {
			t.Errorf("Scan(%q): expected currency %s got %s", tc.input, tc.currency, m.Currency().Code)
		}
	}
}

func TestMoney_Scan_Bytes(t *testing.T) {
	var m Money
	if err := m.Scan([]byte("1234 USD")); err != nil {
		t.Fatalf("Scan([]byte): unexpected error: %v", err)
	}
	if m.Amount() != 1234 {
		t.Errorf("expected 1234 got %d", m.Amount())
	}
}

func TestMoney_Scan_InvalidType(t *testing.T) {
	var m Money
	if err := m.Scan(42); err == nil {
		t.Error("Scan(int): expected error for unsupported type")
	}
}

func TestMoney_Scan_InvalidFormat(t *testing.T) {
	var m Money
	if err := m.Scan("1234"); err == nil {
		t.Error("Scan(no-space): expected error for invalid format")
	}
	if err := m.Scan("abc USD"); err == nil {
		t.Error("Scan(non-numeric): expected error for invalid amount")
	}
}

func TestMoney_SQL_RoundTrip(t *testing.T) {
	original := New(12345, "EUR")

	v, err := original.Value()
	if err != nil {
		t.Fatal(err)
	}

	var restored Money
	if err := restored.Scan(v.(driver.Value)); err != nil {
		t.Fatalf("Scan: %v", err)
	}

	eq, err := original.Equals(&restored)
	if err != nil || !eq {
		t.Errorf("SQL round-trip: values not equal (original=%d, restored=%d)",
			original.Amount(), restored.Amount())
	}
}
