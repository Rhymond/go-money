package money

import (
	"testing"
)

func Testcurrency_Get(t *testing.T) {
	tc := map[string]string{
		"EUR": "EUR",
		"Eur": "EUR",
	}

	for code, expected := range tc {
		c := new(currency).Get(code)

		if c.Code != expected {
			t.Errorf("Expected %s got %s", expected, c.Code)
		}
	}
}

func Testcurrency_Equals(t *testing.T) {
	tc := map[string]string{
		"EUR": "EUR",
		"Eur": "EUR",
		"usd": "USD",
	}

	for code, other := range tc {
		c := new(currency).Get(code)
		oc := new(currency).Get(other)

		if !c.Equals(oc) {
			t.Errorf("Expected that %v is not equal %v", c, oc)
		}
	}
}
