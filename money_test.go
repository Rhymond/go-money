package gocash

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New(1, "EUR")

	if m.Number.Amount != 1 {
		t.Errorf("Expected %d got %d", 1, m.Number.Amount)
	}

	if m.Currency.Code != "EUR" {
		t.Errorf("Expected currency %s got %s", "EUR", m.Currency.Code)
	}

	m = New(-100, "EUR")

	if m.Number.Amount != -100 {
		t.Errorf("Expected %d got %d", -100, m.Number.Amount)
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
			t.Errorf("Expected %d Equals %d == %t got %t", m.Number.Amount, om.Number.Amount, expected, m.Equals(om))
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
			t.Errorf("Expected %d Greater Than %d == %t got %t", m.Number.Amount, om.Number.Amount, expected, m.GreaterThan(om))
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
			t.Errorf("Expected %d Equals Or Greater Than %d == %t got %t", m.Number.Amount, om.Number.Amount, expected, m.GreaterThanOrEqual(om))
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
			t.Errorf("Expected %d Less Than %d == %t got %t", m.Number.Amount, om.Number.Amount, expected, m.LessThan(om))
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
			t.Errorf("Expected %d Equal Or Less Than %d == %t got %t", m.Number.Amount, om.Number.Amount, expected, m.LessThanOrEqual(om))
		}
	}
}

func TestMoney_IsZero(t *testing.T) {
	tc := map[int]bool{
		-1: false,
		0:  true,
		1:  false,
	}

	for amount, expected := range tc {
		m := New(amount, "EUR")

		if m.IsZero() != expected {
			t.Errorf("Expected %d to be zero == %t got %t", m.Number.Amount, expected, m.IsZero())
		}
	}
}

func TestMoney_IsNegative(t *testing.T) {
	tc := map[int]bool{
		-1: true,
		0:  false,
		1:  false,
	}

	for amount, expected := range tc {
		m := New(amount, "EUR")

		if m.IsNegative() != expected {
			t.Errorf("Expected %d to be negative == %t got %t", m.Number.Amount, expected, m.IsNegative())
		}
	}
}

func TestMoney_IsPositive(t *testing.T) {
	tc := map[int]bool{
		-1: false,
		0:  false,
		1:  true,
	}

	for amount, expected := range tc {
		m := New(amount, "EUR")

		if m.IsPositive() != expected {
			t.Errorf("Expected %d to be positive == %t got %t", m.Number.Amount, expected, m.IsPositive())
		}
	}
}

func TestMoney_Absolute(t *testing.T) {
	tc := map[int]int{
		-1: 1,
		0:  0,
		1:  1,
	}

	for amount, expected := range tc {
		m := New(amount, "EUR")

		if m.Absolute().Number.Amount != expected {
			t.Errorf("Expected absolute %d to be %d got %d", m.Number.Amount, expected, m.Absolute().Number.Amount)
		}
	}
}

func TestMoney_Negative(t *testing.T) {
	tc := map[int]int{
		-1: -1,
		0:  0,
		1:  -1,
	}

	for amount, expected := range tc {
		m := New(amount, "EUR")

		if m.Negative().Number.Amount != expected {
			t.Errorf("Expected absolute %d to be %d got %d", m.Number.Amount, expected, m.Negative().Number.Amount)
		}
	}
}

func TestMoney_Add(t *testing.T) {
	tc := map[int][2]int{
		10: {5, 5},
		15: {10, 5},
		0:  {-1, 1},
	}

	for e, a := range tc {
		m := New(a[0], "EUR")
		om := New(a[1], "EUR")

		if m.Add(om).Number.Amount != e {
			t.Errorf(
				"Expected %d + %d = %d got %d",
				a[0],
				a[1],
				e,
				m.Add(om).Number.Amount,
			)
		}
	}
}

func TestMoney_Subtract(t *testing.T) {
	tc := map[int][2]int{
		0:  {5, 5},
		5:  {10, 5},
		-2: {-1, 1},
	}

	for e, a := range tc {
		m := New(a[0], "EUR")
		om := New(a[1], "EUR")

		if m.Subtract(om).Number.Amount != e {
			t.Errorf(
				"Expected %d - %d = %d got %d",
				a[0],
				a[1],
				e,
				m.Subtract(om).Number.Amount,
			)
		}
	}
}

func TestMoney_Multiply(t *testing.T) {
	tc := map[int][2]int{
		25: {5, 5},
		50: {10, 5},
		-1: {-1, 1},
		0:  {1, 0},
	}

	for e, a := range tc {
		m := New(a[0], "EUR")

		if m.Multiply(a[1]).Number.Amount != e {
			t.Errorf(
				"Expected %d * %d = %d got %d",
				a[0],
				a[1],
				e,
				m.Multiply(a[1]).Number.Amount,
			)
		}
	}
}

func TestMoney_Divide(t *testing.T) {
	tc := map[int][2]int{
		1:  {5, 5},
		2:  {10, 5},
		-1: {-1, 1},
		3:  {10, 3},
	}

	for e, a := range tc {
		m := New(a[0], "EUR")

		if m.Divide(a[1]).Number.Amount != e {
			t.Errorf(
				"Expected %d * %d = %d got %d",
				a[0],
				a[1],
				e,
				m.Divide(a[1]).Number.Amount,
			)
		}
	}
}
