package money

import (
	"math"
	"testing"
)

func TestMoney_Multiply(t *testing.T) {
	tcs := []struct {
		amount     int64
		multiplier int64
		expected   int64
	}{
		{5, 5, 25},
		{10, 5, 50},
		{1, -1, -1},
		{1, 0, 0},
	}

	for _, tc := range tcs {
		m := New(tc.amount, EUR)
		r := m.Multiply(NewMultiplier(tc.multiplier)).amount.val.Int64()

		if r != tc.expected {
			t.Errorf("Expected %d * %d = %d got %d", tc.amount, tc.multiplier, tc.expected, r)
		}
	}
}

func TestMoney_MultiplyFloat(t *testing.T) {
	m := New(1000, USD) // $10.00, exponent 0
	r := m.Multiply(NewMultiplier(1.5))
	if r.amount.val.Int64() != 15000 || r.amount.exponent != 1 {
		t.Errorf("expected 15000 e-1, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplySubUnit(t *testing.T) {
	// $0.01 * 1.5 = $0.015 — must preserve the extra precision instead of rounding.
	m := New(1, USD)
	r := m.Multiply(NewMultiplier(1.5))
	if r.amount.val.Int64() != 15 || r.amount.exponent != 1 {
		t.Errorf("expected 15 e-1, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplyString(t *testing.T) {
	m := New(10000, USD) // $100.00
	mul, err := NewMultiplierFromString("0.075")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.Multiply(mul)
	// 10000 * 75 = 750000, exponent 0 + 3 = 3 → 750.000 (in cents)
	if r.amount.val.Int64() != 750000 || r.amount.exponent != 3 {
		t.Errorf("expected 750000 e-3, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplyChained(t *testing.T) {
	m := New(1000, USD) // $10.00
	half, err := NewMultiplierFromString("0.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.Multiply(NewMultiplier(2), half, NewMultiplier(1.5))
	// 1000 * 2 * 5 * 15, exponents 0 + 0 + 1 + 1 = 2 → 150000 e-2 = 1500.00 cents = $15.00
	if r.amount.val.Int64() != 150000 || r.amount.exponent != 2 {
		t.Errorf("expected 150000 e-2, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplyDecimal(t *testing.T) {
	m := New(200, USD)
	d, _ := NewDecimal("2.5")
	r := m.Multiply(NewMultiplier(d))
	if r.amount.val.Int64() != 5000 || r.amount.exponent != 1 {
		t.Errorf("expected 5000 e-1, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplyNoArgs(t *testing.T) {
	m := New(100, USD)
	r := m.Multiply()
	if r.amount.val.Int64() != 100 || r.amount.exponent != 0 {
		t.Errorf("expected unchanged 100 e-0, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplyZeroMultiplier(t *testing.T) {
	// A zero-valued Multiplier{} multiplies by zero, collapsing the result to
	// zero Money. This matches ordinary arithmetic and avoids quietly hiding
	// a caller's mistake (e.g. ignoring an error from NewMultiplierFromString).
	m := New(100, USD)
	r := m.Multiply(Multiplier{})
	if r.amount.val.Sign() != 0 {
		t.Errorf("expected 0, got %s", r.amount.val.String())
	}

	// Subsequent multipliers after a zero stay at zero.
	r = m.Multiply(NewMultiplier(2), Multiplier{}, NewMultiplier(3))
	if r.amount.val.Sign() != 0 {
		t.Errorf("expected 0 after chain with zero multiplier, got %s", r.amount.val.String())
	}
}

func TestMoney_MultiplyGenericIntegerTypes(t *testing.T) {
	m := New(100, USD)

	if r := m.Multiply(NewMultiplier(int8(2))).amount.val.Int64(); r != 200 {
		t.Errorf("int8: got %d", r)
	}
	if r := m.Multiply(NewMultiplier(uint64(3))).amount.val.Int64(); r != 300 {
		t.Errorf("uint64: got %d", r)
	}
	// float32(2.5) parses as 25 e-1, so 100 * (25 e-1) = 2500 e-1 = 250.
	r := m.Multiply(NewMultiplier(float32(2.5)))
	if r.amount.val.Int64() != 2500 || r.amount.exponent != 1 {
		t.Errorf("float32: got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestNewMultiplierFromString_Invalid(t *testing.T) {
	if _, err := NewMultiplierFromString("not-a-number"); err == nil {
		t.Error("expected error for invalid string")
	}
}

func TestNewMultiplierPanics(t *testing.T) {
	cases := []struct {
		name string
		fn   func()
	}{
		{"NaN float64", func() { NewMultiplier(math.NaN()) }},
		{"+Inf float64", func() { NewMultiplier(math.Inf(1)) }},
		{"-Inf float64", func() { NewMultiplier(math.Inf(-1)) }},
		{"NaN float32", func() { NewMultiplier(float32(math.NaN())) }},
		{"nil decimal", func() { NewMultiplier((*Decimal)(nil)) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for %s", tc.name)
				}
			}()
			tc.fn()
		})
	}
}
