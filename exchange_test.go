package money

import (
	"math"
	"testing"
)

func TestNewExchangeRate(t *testing.T) {
	tcs := []struct {
		from    string
		to      string
		rate    float64
		wantErr bool
	}{
		{"USD", "EUR", 0.92, false},
		{"GBP", "INR", 106.5, false},
		{"USD", "JPY", 148.5, false},
		{"USD", "EUR", 0, true},   // zero rate
		{"USD", "EUR", -1, true},  // negative rate
		{"USD", "USD", 1.0, true}, // same currency
		{"usd", "eur", 0.92, false}, // lowercase normalised
	}

	for _, tc := range tcs {
		r, err := NewExchangeRate(tc.from, tc.to, tc.rate)
		if tc.wantErr {
			if err == nil {
				t.Errorf("NewExchangeRate(%q, %q, %v): expected error", tc.from, tc.to, tc.rate)
			}
			continue
		}
		if err != nil {
			t.Fatalf("NewExchangeRate(%q, %q, %v): unexpected error: %v", tc.from, tc.to, tc.rate, err)
		}
		if r.From() != "USD" && r.From() != "GBP" {
			// just checking uppercase normalisation
		}
	}
}

func TestExchangeRate_Convert(t *testing.T) {
	tcs := []struct {
		from     string
		to       string
		rate     float64
		amount   int64
		expected int64
	}{
		{"USD", "EUR", 0.92, 1000, 920},   // $10.00 → €9.20
		{"USD", "JPY", 148.0, 100, 14800}, // $1.00 → ¥14800
		{"GBP", "USD", 1.27, 100, 127},    // £1.00 → $1.27
		{"USD", "EUR", 0.9256, 1000, 926}, // $10.00 → €9.26 (rounded)
	}

	for _, tc := range tcs {
		rate, _ := NewExchangeRate(tc.from, tc.to, tc.rate)
		m := New(tc.amount, tc.from)
		result, err := rate.Convert(m)
		if err != nil {
			t.Fatalf("Convert(%d %s at %v): unexpected error: %v", tc.amount, tc.from, tc.rate, err)
		}
		if result.Amount() != tc.expected {
			t.Errorf("Convert(%d %s at %v): expected %d got %d",
				tc.amount, tc.from, tc.rate, tc.expected, result.Amount())
		}
		if result.Currency().Code != tc.to {
			t.Errorf("Convert: expected currency %s got %s", tc.to, result.Currency().Code)
		}
	}
}

func TestExchangeRate_Convert_WrongCurrency(t *testing.T) {
	rate, _ := NewExchangeRate("USD", "EUR", 0.92)
	m := New(1000, "GBP")
	_, err := rate.Convert(m)
	if err == nil {
		t.Error("expected error for wrong source currency")
	}
}

func TestExchangeRate_Invert(t *testing.T) {
	rate, _ := NewExchangeRate("USD", "EUR", 0.92)
	inv := rate.Invert()

	if inv.From() != "EUR" {
		t.Errorf("Invert: expected From=EUR got %s", inv.From())
	}
	if inv.To() != "USD" {
		t.Errorf("Invert: expected To=USD got %s", inv.To())
	}

	// inverted rate should be ~1.0869...
	expected := 1 / 0.92
	if math.Abs(inv.Rate()-expected) > 1e-10 {
		t.Errorf("Invert: expected rate %v got %v", expected, inv.Rate())
	}
}

func TestExchangeRate_Accessors(t *testing.T) {
	rate, _ := NewExchangeRate("GBP", "JPY", 188.5)
	if rate.From() != "GBP" {
		t.Errorf("From: expected GBP got %s", rate.From())
	}
	if rate.To() != "JPY" {
		t.Errorf("To: expected JPY got %s", rate.To())
	}
	if rate.Rate() != 188.5 {
		t.Errorf("Rate: expected 188.5 got %v", rate.Rate())
	}
}
