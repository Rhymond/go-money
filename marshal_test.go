package money

import (
	"bytes"
	"encoding/xml"
	"testing"
)

// ---------- MarshalBinary / UnmarshalBinary ----------

func TestMoney_MarshalBinary(t *testing.T) {
	m := New(1234, "USD")
	b, err := m.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary: unexpected error: %v", err)
	}
	if len(b) != 8+3 { // 8 bytes int64 + "USD"
		t.Errorf("MarshalBinary: expected len 11 got %d", len(b))
	}

	var m2 Money
	if err := m2.UnmarshalBinary(b); err != nil {
		t.Fatalf("UnmarshalBinary: unexpected error: %v", err)
	}
	if m2.Amount() != 1234 || m2.Currency().Code != "USD" {
		t.Errorf("Binary round-trip: expected 1234 USD got %d %s", m2.Amount(), m2.Currency().Code)
	}
}

func TestMoney_MarshalBinary_Negative(t *testing.T) {
	m := New(-500, "EUR")
	b, err := m.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var m2 Money
	if err := m2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if m2.Amount() != -500 || m2.Currency().Code != "EUR" {
		t.Errorf("Binary negative round-trip: expected -500 EUR got %d %s", m2.Amount(), m2.Currency().Code)
	}
}

func TestMoney_UnmarshalBinary_TooShort(t *testing.T) {
	var m Money
	if err := m.UnmarshalBinary([]byte{0x01, 0x02}); err == nil {
		t.Error("UnmarshalBinary: expected error for too-short data")
	}
}

// ---------- MarshalText / UnmarshalText ----------

func TestMoney_MarshalText(t *testing.T) {
	m := New(1234, "USD")
	b, err := m.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText: unexpected error: %v", err)
	}
	if string(b) != "1234 USD" {
		t.Errorf("MarshalText: expected %q got %q", "1234 USD", string(b))
	}
}

func TestMoney_UnmarshalText(t *testing.T) {
	var m Money
	if err := m.UnmarshalText([]byte("5678 GBP")); err != nil {
		t.Fatalf("UnmarshalText: unexpected error: %v", err)
	}
	if m.Amount() != 5678 || m.Currency().Code != "GBP" {
		t.Errorf("UnmarshalText: expected 5678 GBP got %d %s", m.Amount(), m.Currency().Code)
	}
}

func TestMoney_UnmarshalText_BadFormat(t *testing.T) {
	var m Money
	if err := m.UnmarshalText([]byte("nodollar")); err == nil {
		t.Error("UnmarshalText: expected error for missing separator")
	}
}

func TestMoney_UnmarshalText_BadAmount(t *testing.T) {
	var m Money
	if err := m.UnmarshalText([]byte("abc USD")); err == nil {
		t.Error("UnmarshalText: expected error for non-numeric amount")
	}
}

func TestMoney_Text_RoundTrip(t *testing.T) {
	original := New(-9999, "JPY")
	b, _ := original.MarshalText()
	var restored Money
	if err := restored.UnmarshalText(b); err != nil {
		t.Fatal(err)
	}
	if restored.Amount() != -9999 || restored.Currency().Code != "JPY" {
		t.Errorf("Text round-trip: expected -9999 JPY got %d %s", restored.Amount(), restored.Currency().Code)
	}
}

// ---------- MarshalXML / UnmarshalXML ----------

func TestMoney_MarshalXML(t *testing.T) {
	m := New(1234, "USD")

	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	start := xml.StartElement{Name: xml.Name{Local: "Money"}}
	if err := enc.EncodeElement(m, start); err != nil {
		t.Fatalf("MarshalXML: %v", err)
	}
	enc.Flush()

	got := buf.String()
	if got == "" {
		t.Error("MarshalXML: empty output")
	}
	// Should contain amount and currency
	if !bytes.Contains(buf.Bytes(), []byte("1234")) {
		t.Errorf("MarshalXML: expected amount 1234 in output: %s", got)
	}
	if !bytes.Contains(buf.Bytes(), []byte("USD")) {
		t.Errorf("MarshalXML: expected currency USD in output: %s", got)
	}
}

func TestMoney_UnmarshalXML(t *testing.T) {
	input := `<Money><amount>5678</amount><currency>EUR</currency></Money>`
	var m Money
	if err := xml.Unmarshal([]byte(input), &m); err != nil {
		t.Fatalf("UnmarshalXML: %v", err)
	}
	if m.Amount() != 5678 || m.Currency().Code != "EUR" {
		t.Errorf("UnmarshalXML: expected 5678 EUR got %d %s", m.Amount(), m.Currency().Code)
	}
}

func TestMoney_UnmarshalXML_Error(t *testing.T) {
	// amount field is not an integer — DecodeElement should fail
	input := `<Money><amount>notanint</amount><currency>USD</currency></Money>`
	var m Money
	if err := xml.Unmarshal([]byte(input), &m); err == nil {
		t.Error("UnmarshalXML: expected error for non-integer amount")
	}
}

func TestMoney_XML_RoundTrip(t *testing.T) {
	original := New(-12345, "GBP")

	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	start := xml.StartElement{Name: xml.Name{Local: "Money"}}
	enc.EncodeElement(original, start)
	enc.Flush()

	var restored Money
	if err := xml.Unmarshal(buf.Bytes(), &restored); err != nil {
		t.Fatalf("XML round-trip unmarshal: %v", err)
	}
	if restored.Amount() != -12345 || restored.Currency().Code != "GBP" {
		t.Errorf("XML round-trip: expected -12345 GBP got %d %s", restored.Amount(), restored.Currency().Code)
	}
}
