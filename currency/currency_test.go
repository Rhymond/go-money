package currency

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
		c := NewCurrency(tc.code).Get()

		if c.Code != tc.expected {
			t.Errorf("Expected %s got %s", tc.expected, c.Code)
		}
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
		c := NewCurrency(tc.code).Get()
		oc := NewCurrency(tc.other).Get()

		if !c.Equals(oc) {
			t.Errorf("Expected that %v is not equal %v", c, oc)
		}
	}
}

func TestAddCurrency(t *testing.T) {
	tcs := []struct {
		code     string
		template string
	}{
		{"GOLD", "1$"},
	}

	for _, tc := range tcs {
		AddCurrency(tc.code, "", tc.template, "", "", 0)
		c := NewCurrency(tc.code).Get()

		if c.Template != tc.template {
			t.Errorf("Expected currency template %v got %v", tc.template, c.Template)
		}
	}
}
