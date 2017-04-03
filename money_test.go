package money

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	m := New(1, "EUR")

	if m.amount.val != 1 {
		t.Errorf("Expected %d got %d", 1, m.amount.val)
	}

	if m.currency.Code != "EUR" {
		t.Errorf("Expected currency %s got %s", "EUR", m.currency.Code)
	}

	m = New(-100, "EUR")

	if m.amount.val != -100 {
		t.Errorf("Expected %d got %d", -100, m.amount.val)
	}
}

func TestMoney_Samecurrency(t *testing.T) {
	m := New(0, "EUR")
	om := New(0, "USD")

	if m.Samecurrency(om) {
		t.Errorf("Expected %s not to be same as %s", m.currency.Code, om.currency.Code)
	}

	om = New(0, "EUR")

	if !m.Samecurrency(om) {
		t.Errorf("Expected %s to be same as %s", m.currency.Code, om.currency.Code)
	}
}

func TestMoney_Equals(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		amount   int
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.amount, "EUR")
		r := m.Equals(om)

		if r != tc.expected {
			t.Errorf("Expected %d Equals %d == %t got %t", m.amount.val,
				om.amount.val, tc.expected, r)
		}
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		amount   int
		expected bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.amount, "EUR")
		r := m.GreaterThan(om)

		if r != tc.expected {
			t.Errorf("Expected %d Greater Than %d == %t got %t", m.amount.val,
				om.amount.val, tc.expected, r)
		}
	}
}

func TestMoney_GreaterThanOrEqual(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		amount   int
		expected bool
	}{
		{-1, true},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.amount, "EUR")
		r := m.GreaterThanOrEqual(om)

		if r != tc.expected {
			t.Errorf("Expected %d Equals Or Greater Than %d == %t got %t", m.amount.val,
				om.amount.val, tc.expected, r)
		}
	}
}

func TestMoney_LessThan(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		amount   int
		expected bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}

	for _, tc := range tcs {
		om := New(tc.amount, "EUR")
		r := m.LessThan(om)

		if r != tc.expected {
			t.Errorf("Expected %d Less Than %d == %t got %t", m.amount.val,
				om.amount.val, tc.expected, r)
		}
	}
}

func TestMoney_LessThanOrEqual(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		amount   int
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, true},
	}

	for _, tc := range tcs {
		om := New(tc.amount, "EUR")
		r := m.LessThanOrEqual(om)

		if r != tc.expected {
			t.Errorf("Expected %d Equal Or Less Than %d == %t got %t", m.amount.val,
				om.amount.val, tc.expected, r)
		}
	}
}

func TestMoney_IsZero(t *testing.T) {
	tcs := []struct {
		amount   int
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		m := New(tc.amount, "EUR")
		r := m.IsZero()

		if r != tc.expected {
			t.Errorf("Expected %d to be zero == %t got %t", m.amount.val, tc.expected, r)
		}
	}
}

func TestMoney_IsNegative(t *testing.T) {
	tcs := []struct {
		amount   int
		expected bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}

	for _, tc := range tcs {
		m := New(tc.amount, "EUR")
		r := m.IsNegative()

		if r != tc.expected {
			t.Errorf("Expected %d to be negative == %t got %t", m.amount.val,
				tc.expected, r)
		}
	}
}

func TestMoney_IsPositive(t *testing.T) {
	tcs := []struct {
		amount   int
		expected bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}

	for _, tc := range tcs {
		m := New(tc.amount, "EUR")
		r := m.IsPositive()

		if r != tc.expected {
			t.Errorf("Expected %d to be positive == %t got %t", m.amount.val,
				tc.expected, r)
		}
	}
}

func TestMoney_Absolute(t *testing.T) {
	tcs := []struct {
		amount   int
		expected int
	}{
		{-1, 1},
		{0, 0},
		{1, 1},
	}

	for _, tc := range tcs {
		m := New(tc.amount, "EUR")
		r := m.Absolute().amount.val

		if r != tc.expected {
			t.Errorf("Expected absolute %d to be %d got %d", m.amount.val,
				tc.expected, r)
		}
	}
}

func TestMoney_Negative(t *testing.T) {
	tcs := []struct {
		amount   int
		expected int
	}{
		{-1, -1},
		{0, -0},
		{1, -1},
	}

	for _, tc := range tcs {
		m := New(tc.amount, "EUR")
		r := m.Negative().amount.val

		if r != tc.expected {
			t.Errorf("Expected absolute %d to be %d got %d", m.amount.val,
				tc.expected, r)
		}
	}
}

func TestMoney_Add(t *testing.T) {
	tcs := []struct {
		amount1  int
		amount2  int
		expected int
	}{
		{5, 5, 10},
		{10, 5, 15},
		{1, -1, 0},
	}

	for _, tc := range tcs {
		m := New(tc.amount1, "EUR")
		om := New(tc.amount2, "EUR")
		r := m.Add(om).amount.val

		if r != tc.expected {
			t.Errorf("Expected %d + %d = %d got %d", tc.amount1, tc.amount2, tc.expected, r)
		}
	}
}

func TestMoney_Subtract(t *testing.T) {
	tcs := []struct {
		amount1  int
		amount2  int
		expected int
	}{
		{5, 5, 0},
		{10, 5, 5},
		{1, -1, 2},
	}

	for _, tc := range tcs {
		m := New(tc.amount1, "EUR")
		om := New(tc.amount2, "EUR")
		r := m.Subtract(om).amount.val

		if r != tc.expected {
			t.Errorf("Expected %d - %d = %d got %d", tc.amount1, tc.amount2, tc.expected, r)
		}
	}
}

func TestMoney_Multiply(t *testing.T) {
	tcs := []struct {
		amount     int
		multiplier int
		expected   int
	}{
		{5, 5, 25},
		{10, 5, 50},
		{1, -1, -1},
		{1, 0, 0},
	}

	for _, tc := range tcs {
		m := New(tc.amount, "EUR")
		r := m.Multiply(tc.multiplier).amount.val

		if r != tc.expected {
			t.Errorf("Expected %d * %d = %d got %d", tc.amount, tc.multiplier, tc.expected, r)
		}
	}
}

func TestMoney_Divide(t *testing.T) {
	tcs := []struct {
		amount   int
		divisor  int
		expected int
	}{
		{5, 5, 1},
		{10, 5, 2},
		{1, -1, -1},
		{10, 3, 3},
	}

	for _, tc := range tcs {
		m := New(tc.amount, "EUR")
		r := m.Divide(tc.divisor).amount.val

		if r != tc.expected {
			t.Errorf("Expected %d * %d = %d got %d", tc.amount, tc.divisor, tc.expected, r)
		}
	}
}

func TestMoney_Round(t *testing.T) {
	tcs := []struct {
		amount   int
		expected int
	}{
		{125, 100},
		{175, 200},
		{349, 300},
		{351, 400},
		{0, 0},
		{-1, 0},
		{-75, -100},
	}

	for _, tc := range tcs {
		m := New(tc.amount, "EUR")
		r := m.Round().amount.val

		if r != tc.expected {
			t.Errorf("Expected rounded %d to be %d got %d", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_Split(t *testing.T) {
	tcs := []struct {
		amount   int
		split    int
		expected []int
	}{
		{100, 3, []int{34, 33, 33}},
		{100, 4, []int{25, 25, 25, 25}},
		{5, 3, []int{2, 2, 1}},
	}

	for _, tc := range tcs {
		m := New(tc.amount, "EUR")
		var rs []int

		for _, party := range m.Split(tc.split) {
			rs = append(rs, party.amount.val)
		}

		if !reflect.DeepEqual(tc.expected, rs) {
			t.Errorf("Expected split of %d to be %v got %v", tc.amount, tc.expected, rs)
		}
	}
}

func TestMoney_Allocate(t *testing.T) {
	tcs := []struct {
		amount   int
		ratios   []int
		expected []int
	}{
		{100, []int{50, 50}, []int{50, 50}},
		{100, []int{30, 30, 30}, []int{34, 33, 33}},
		{200, []int{25, 25, 50}, []int{50, 50, 100}},
		{5, []int{50, 25, 25}, []int{3, 1, 1}},
	}

	for _, tc := range tcs {
		m := New(tc.amount, "EUR")
		var rs []int

		for _, party := range m.Allocate(tc.ratios) {
			rs = append(rs, party.amount.val)
		}

		if !reflect.DeepEqual(tc.expected, rs) {
			t.Errorf("Expected allocation of %d for ratios %v to be %v got %v", tc.amount, tc.ratios,
				tc.expected, rs)
		}
	}
}

func TestMoney_Chain(t *testing.T) {
	m := New(10, "EUR")
	om := New(5, "EUR")

	// 10 + 5 = 15 / 5 = 3 * 4 = 12 - 5 = 7
	r := m.Add(om).Divide(5).Multiply(4).Subtract(New(5, "EUR")).amount.val
	e := 7

	if r != e {
		t.Errorf("Expected %d got %d", e, r)
	}
}

func TestMoney_Format(t *testing.T) {
	tcs := []struct {
		amount   int
		code     string
		expected string
	}{
		{100, "GBP", "£1.00"},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		r := m.Display()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.amount, tc.expected, r)
		}
	}

}
