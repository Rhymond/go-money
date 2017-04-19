package money

import (
	"testing"
)

func TestCurrency_Get(t *testing.T) {
	tcs := []struct {
		code     string
		expected string
	}{
		{"EUR", "EUR"},
		{"Eur", "EUR"},
	}

	for _, tc := range tcs {
		c := newCurrency(tc.code).get()

		if c.Code != tc.expected {
			t.Errorf("Expected %s got %s", tc.expected, c.Code)
		}
	}
}

func TestCurrency_Get1(t *testing.T) {
	code := "RANDOM"
	c := newCurrency(code).get()

	if c.Grapheme != code {
		t.Errorf("Expected %s got %s", c.Grapheme, code)
	}
}

func TestCurrency_Equals(t *testing.T) {
	tcs := []struct {
		code  string
		other string
	}{
		{"EUR", "EUR"},
		{"Eur", "EUR"},
		{"usd", "USD"},
	}

	for _, tc := range tcs {
		c := newCurrency(tc.code).get()
		oc := newCurrency(tc.other).get()

		if !c.equals(oc) {
			t.Errorf("Expected that %v is not equal %v", c, oc)
		}
	}
}

func TestCurrency_AddCurrency(t *testing.T) {
	tcs := []struct {
		code     string
		template string
	}{
		{"GOLD", "1$"},
	}

	for _, tc := range tcs {
		AddCurrency(tc.code, "", tc.template, "", "", 0)
		c := newCurrency(tc.code).get()

		if c.Template != tc.template {
			t.Errorf("Expected currency template %v got %v", tc.template, c.Template)
		}
	}
}
