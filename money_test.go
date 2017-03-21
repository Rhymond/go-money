package gocash

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New(1, "EUR")

	if m.Amount != 1 {
		t.Errorf("Expected %d got %d", 1, m.Amount)
	}

	if m.Currency.Code != "EUR" {
		t.Errorf("Expected currency %s got %s", "EUR", m.Currency.Code)
	}

	m = New(-100, "EUR")

	if m.Amount != -100 {
		t.Errorf("Expected %d got %d", -100, m.Amount)
	}
}

func TestMoney_SameCurrency(t *testing.T) {
	m := New(0, "EUR")
	om := New(0, "USD")

	if m.SameCurrency(om) {
		t.Errorf("Expected %s not to be same as %s", m.Currency.Code, om.Currency.Code)
	}

	om = New(0, "EUR")

	if !m.SameCurrency(om) {
		t.Errorf("Expected %s to be same as %s", m.Currency.Code, om.Currency.Code)
	}
}

func TestMoney_Equals(t *testing.T) {
	m := New(0, "EUR")

	tc := map[int]bool{
		-1: false,
		0:  true,
		1:  false,
	}

	for amount, expected := range tc {
		om := New(amount, "EUR")

		if m.Equals(om) != expected {
			t.Errorf("Expected %d Equals %d == %t got %t", m.Amount, om.Amount, expected, m.Equals(om))
		}
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	m := New(0, "EUR")

	tc := map[int]bool{
		-1: true,
		0:  false,
		1:  false,
	}

	for amount, expected := range tc {
		om := New(amount, "EUR")

		if m.GreaterThan(om) != expected {
			t.Errorf("Expected %d Greater Than %d == %t got %t", m.Amount, om.Amount, expected, m.GreaterThan(om))
		}
	}
}

func TestMoney_GreaterThanOrEqual(t *testing.T) {
	m := New(0, "EUR")

	tc := map[int]bool{
		-1: true,
		0:  true,
		1:  false,
	}

	for amount, expected := range tc {
		om := New(amount, "EUR")

		if m.GreaterThanOrEqual(om) != expected {
			t.Errorf("Expected %d Equals Or Greater Than %d == %t got %t", m.Amount, om.Amount, expected, m.GreaterThanOrEqual(om))
		}
	}
}

func TestMoney_LessThan(t *testing.T) {
	m := New(0, "EUR")

	tc := map[int]bool{
		-1: false,
		0:  false,
		1:  true,
	}

	for amount, expected := range tc {
		om := New(amount, "EUR")

		if m.LessThan(om) != expected {
			t.Errorf("Expected %d Less Than %d == %t got %t", m.Amount, om.Amount, expected, m.LessThan(om))
		}
	}
}

func TestMoney_LessThanOrEqual(t *testing.T) {
	m := New(0, "EUR")

	tc := map[int]bool{
		-1: false,
		0:  true,
		1:  true,
	}

	for amount, expected := range tc {
		om := New(amount, "EUR")

		if m.LessThanOrEqual(om) != expected {
			t.Errorf("Expected %d Equal Or Less Than %d == %t got %t", m.Amount, om.Amount, expected, m.LessThanOrEqual(om))
		}
	}
}

func TestMoney_IsZero(t *testing.T) {
	m := New(0, "EUR")

	if !m.IsZero() {
		t.Errorf("Expected zero got %d", m.Amount)
	}

	m = New(1, "EUR")

	if m.IsZero() {
		t.Errorf("Expected non zero got %d", m.Amount)
	}
}

func TestMoney_IsNegative(t *testing.T) {
	m := New(0, "EUR")

	if m.IsNegative() {
		t.Errorf("Expected not negative got %d", m.Amount)
	}

	m = New(-1, "EUR")

	if !m.IsNegative() {
		t.Errorf("Expected negative got %d", m.Amount)
	}
}

func TestMoney_IsPositive(t *testing.T) {
	m := New(0, "EUR")

	if m.IsPositive() {
		t.Errorf("Expected not positive got %d", m.Amount)
	}

	m = New(1, "EUR")

	if !m.IsPositive() {
		t.Errorf("Expected positive got %d", m.Amount)
	}
}
