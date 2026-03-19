package money

import (
	"encoding/json"
	"testing"
)

// ---------- NewNullMoney ----------

func TestNewNullMoney(t *testing.T) {
	m := New(100, "USD")
	n := NewNullMoney(m)
	if !n.Valid {
		t.Error("NewNullMoney: expected Valid=true")
	}
	if n.Money.Amount() != 100 {
		t.Errorf("NewNullMoney: expected amount 100 got %d", n.Money.Amount())
	}
}

// ---------- Value ----------

func TestNullMoney_Value_Valid(t *testing.T) {
	n := NewNullMoney(New(1234, "USD"))
	v, err := n.Value()
	if err != nil {
		t.Fatalf("NullMoney.Value: unexpected error: %v", err)
	}
	if v == nil {
		t.Error("NullMoney.Value: expected non-nil driver.Value")
	}
}

func TestNullMoney_Value_Null(t *testing.T) {
	n := NullMoney{Valid: false}
	v, err := n.Value()
	if err != nil {
		t.Fatalf("NullMoney.Value: unexpected error: %v", err)
	}
	if v != nil {
		t.Errorf("NullMoney.Value(null): expected nil got %v", v)
	}
}

// ---------- Scan ----------

func TestNullMoney_Scan_Nil(t *testing.T) {
	var n NullMoney
	if err := n.Scan(nil); err != nil {
		t.Fatalf("NullMoney.Scan(nil): unexpected error: %v", err)
	}
	if n.Valid {
		t.Error("NullMoney.Scan(nil): expected Valid=false")
	}
}

func TestNullMoney_Scan_String(t *testing.T) {
	var n NullMoney
	if err := n.Scan("1234 USD"); err != nil {
		t.Fatalf("NullMoney.Scan: unexpected error: %v", err)
	}
	if !n.Valid {
		t.Error("NullMoney.Scan: expected Valid=true")
	}
	if n.Money.Amount() != 1234 {
		t.Errorf("NullMoney.Scan: expected amount 1234 got %d", n.Money.Amount())
	}
}

// ---------- MarshalJSON / UnmarshalJSON ----------

func TestNullMoney_MarshalJSON_Valid(t *testing.T) {
	n := NewNullMoney(New(1234, "USD"))
	b, err := json.Marshal(n)
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	if string(b) != `{"amount":1234,"currency":"USD"}` {
		t.Errorf("MarshalJSON valid: expected object got %s", string(b))
	}
}

func TestNullMoney_MarshalJSON_Null(t *testing.T) {
	n := NullMoney{Valid: false}
	b, err := json.Marshal(n)
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	if string(b) != "null" {
		t.Errorf("MarshalJSON null: expected 'null' got %s", string(b))
	}
}

func TestNullMoney_UnmarshalJSON_Valid(t *testing.T) {
	var n NullMoney
	if err := json.Unmarshal([]byte(`{"amount":5678,"currency":"EUR"}`), &n); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}
	if !n.Valid || n.Money.Amount() != 5678 || n.Money.Currency().Code != "EUR" {
		t.Errorf("UnmarshalJSON valid: unexpected result %+v", n)
	}
}

func TestNullMoney_UnmarshalJSON_Null(t *testing.T) {
	var n NullMoney
	if err := json.Unmarshal([]byte("null"), &n); err != nil {
		t.Fatalf("UnmarshalJSON null: %v", err)
	}
	if n.Valid {
		t.Error("UnmarshalJSON null: expected Valid=false")
	}
}

func TestNullMoney_JSON_RoundTrip(t *testing.T) {
	original := NewNullMoney(New(9999, "JPY"))
	b, _ := json.Marshal(original)
	var restored NullMoney
	if err := json.Unmarshal(b, &restored); err != nil {
		t.Fatal(err)
	}
	if !restored.Valid || restored.Money.Amount() != 9999 {
		t.Errorf("JSON round-trip: unexpected result %+v", restored)
	}
}

// ---------- MarshalText / UnmarshalText ----------

func TestNullMoney_MarshalText_Valid(t *testing.T) {
	n := NewNullMoney(New(100, "USD"))
	b, err := n.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "100 USD" {
		t.Errorf("MarshalText valid: expected '100 USD' got %q", string(b))
	}
}

func TestNullMoney_MarshalText_Null(t *testing.T) {
	n := NullMoney{Valid: false}
	b, err := n.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 0 {
		t.Errorf("MarshalText null: expected empty got %q", string(b))
	}
}

func TestNullMoney_UnmarshalText_Valid(t *testing.T) {
	var n NullMoney
	if err := n.UnmarshalText([]byte("200 GBP")); err != nil {
		t.Fatal(err)
	}
	if !n.Valid || n.Money.Amount() != 200 || n.Money.Currency().Code != "GBP" {
		t.Errorf("UnmarshalText valid: unexpected %+v", n)
	}
}

func TestNullMoney_UnmarshalText_Empty(t *testing.T) {
	var n NullMoney
	if err := n.UnmarshalText([]byte{}); err != nil {
		t.Fatal(err)
	}
	if n.Valid {
		t.Error("UnmarshalText empty: expected Valid=false")
	}
}
