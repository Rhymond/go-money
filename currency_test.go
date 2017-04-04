package money

import (
	"testing"
)

func TestCurrency_Get(t *testing.T) {
	tc := map[string]string{
		"EUR": "EUR",
		"Eur": "EUR",
	}

	for code, expected := range tc {
		c, _ := newCurrency(code).get()

		if c.Code != expected {
			t.Errorf("Expected %s got %s", expected, c.Code)
		}
	}
}

func TestCurrency_Equals(t *testing.T) {
	tc := map[string]string{
		"EUR": "EUR",
		"Eur": "EUR",
		"usd": "USD",
	}

	for code, other := range tc {
		c, _ := newCurrency(code).get()
		oc, _ := newCurrency(other).get()

		if !c.equals(oc) {
			t.Errorf("Expected that %v is not equal %v", c, oc)
		}
	}
}
