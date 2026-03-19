package money

import "testing"

func TestCurrencyForCountry(t *testing.T) {
	tcs := []struct {
		code     string
		expected string
		ok       bool
	}{
		{"US", "USD", true},
		{"us", "USD", true}, // case-insensitive
		{"GB", "GBP", true},
		{"JP", "JPY", true},
		{"DE", "EUR", true},
		{"AU", "AUD", true},
		{"IN", "INR", true},
		{"XX", "", false},
		{"ZZ", "", false},
	}
	for _, tc := range tcs {
		got, ok := CurrencyForCountry(tc.code)
		if ok != tc.ok || got != tc.expected {
			t.Errorf("CurrencyForCountry(%q): expected (%q, %v) got (%q, %v)",
				tc.code, tc.expected, tc.ok, got, ok)
		}
	}
}

func TestNewForCountry(t *testing.T) {
	// known country
	m, err := NewForCountry(1000, "US")
	if err != nil {
		t.Fatalf("NewForCountry(US): unexpected error: %v", err)
	}
	if m.Currency().Code != "USD" {
		t.Errorf("NewForCountry(US): expected USD got %s", m.Currency().Code)
	}
	if m.Amount() != 1000 {
		t.Errorf("NewForCountry(US): expected amount 1000 got %d", m.Amount())
	}

	// case-insensitive
	m2, err := NewForCountry(500, "gb")
	if err != nil {
		t.Fatalf("NewForCountry(gb): unexpected error: %v", err)
	}
	if m2.Currency().Code != "GBP" {
		t.Errorf("NewForCountry(gb): expected GBP got %s", m2.Currency().Code)
	}

	// unknown country
	_, err = NewForCountry(100, "XX")
	if err == nil {
		t.Error("NewForCountry(XX): expected error for unknown country")
	}
}
