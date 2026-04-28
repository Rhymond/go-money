package money

import (
	"math/big"
	"testing"
)

func TestNewFromFixedPoint(t *testing.T) {
	tcs := []struct {
		units    int64
		fraction int64
		size     uint
		code     string
		expected int64
	}{
		{1, 250_000_000, FixedPointSizeNano, EUR, 125},
		{1, 5, FixedPointSizeDeci, USD, 150},
		{1, 25, FixedPointSizeCenti, USD, 125},
		{1, 250, FixedPointSizeMilli, USD, 125},
		{1, 250_000, FixedPointSizeMicro, USD, 125},
		{1, 250_000_000, FixedPointSizeNano, USD, 125},
		{100, 999_999_999, FixedPointSizeNano, JPY, 100},
		{1, 250_000_000, FixedPointSizeNano, BHD, 1250},
		{-1, -750_000_000, FixedPointSizeNano, USD, -175},
		{0, 1, FixedPointSizeNano, EUR, 0},
	}

	for _, tc := range tcs {
		m := NewFromFixedPoint(tc.units, tc.fraction, tc.size, tc.code)
		r := m.Amount()

		if r != tc.expected {
			t.Errorf("Expected NewFromFixedPoint(%d, %d, %d, %s) amount to be %d got %d",
				tc.units, tc.fraction, tc.size, tc.code, tc.expected, r)
		}
	}
}

// uint64 values above math.MaxInt64 must not be reinterpreted as negative.
func TestNewFromFixedPoint_LargeUint64(t *testing.T) {
	const big uint64 = 1 << 40
	m := NewFromFixedPoint(big, uint64(0), FixedPointSizeCenti, USD)
	expected := int64(big) * 100

	if r := m.Amount(); r != expected {
		t.Errorf("Expected %d got %d", expected, r)
	}
}

func TestNewFromMajorUnits(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected int64
	}{
		{5, USD, 500},
		{1, JPY, 1},
		{1, BHD, 1000},
		{-3, EUR, -300},
	}

	for _, tc := range tcs {
		m := NewFromMajorUnits(tc.amount, tc.code)
		r := m.Amount()

		if r != tc.expected {
			t.Errorf("Expected NewFromMajorUnits(%d, %s) = %d got %d",
				tc.amount, tc.code, tc.expected, r)
		}
	}
}

// NewFromMajorUnits should accept every integer width / signedness and produce
// equivalent Money.
func TestNewFromMajorUnits_AllIntegerTypes(t *testing.T) {
	expected := New(500, USD).Amount()
	tcs := []struct {
		name string
		got  int64
	}{
		{"int", NewFromMajorUnits(int(5), USD).Amount()},
		{"int8", NewFromMajorUnits(int8(5), USD).Amount()},
		{"int16", NewFromMajorUnits(int16(5), USD).Amount()},
		{"int32", NewFromMajorUnits(int32(5), USD).Amount()},
		{"int64", NewFromMajorUnits(int64(5), USD).Amount()},
		{"uint", NewFromMajorUnits(uint(5), USD).Amount()},
		{"uint8", NewFromMajorUnits(uint8(5), USD).Amount()},
		{"uint16", NewFromMajorUnits(uint16(5), USD).Amount()},
		{"uint32", NewFromMajorUnits(uint32(5), USD).Amount()},
		{"uint64", NewFromMajorUnits(uint64(5), USD).Amount()},
	}

	for _, tc := range tcs {
		if tc.got != expected {
			t.Errorf("Expected %s amount %d got %d", tc.name, expected, tc.got)
		}
	}
}

func TestNewFromFloat(t *testing.T) {
	tcs := []struct {
		amount   float64
		code     string
		expected int64
	}{
		{1.25, USD, 125},
		{-0.99, USD, -99},
		{0.014, USD, 1},
	}

	for _, tc := range tcs {
		m, err := NewFromFloat(tc.amount, tc.code)
		if err != nil {
			t.Errorf("Unexpected error for %v: %v", tc.amount, err)
			continue
		}
		r := m.Amount()

		if r != tc.expected {
			t.Errorf("Expected NewFromFloat(%v, %s) = %d got %d",
				tc.amount, tc.code, tc.expected, r)
		}
	}
}

// Pins exact values from issues #121 and #124 so a future change to the float
// path can't silently re-introduce the off-by-one-cent bug. Each value has a
// binary representation that undershoots the displayed decimal — the
// master-branch implementation (int64(amount * 100)) returned one cent low.
func TestNewFromFloat_IEEE754Pitfalls(t *testing.T) {
	tcs := []struct {
		amount   float64
		code     string
		expected int64
	}{
		{1.15, USD, 115},
		{136.98, USD, 13698},
		{18.99, USD, 1899},
		{73708.43, EUR, 7370843},
		{0.1 + 0.2, USD, 30},
	}

	for _, tc := range tcs {
		m, err := NewFromFloat(tc.amount, tc.code)
		if err != nil {
			t.Errorf("Unexpected error for %v: %v", tc.amount, err)
			continue
		}
		r := m.Amount()

		if r != tc.expected {
			t.Errorf("Expected NewFromFloat(%v, %s) = %d got %d",
				tc.amount, tc.code, tc.expected, r)
		}
	}
}

func TestNewFromFloat_Float32(t *testing.T) {
	m, err := NewFromFloat(float32(1.5), USD)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if r := m.Amount(); r != 150 {
		t.Errorf("Expected 150 got %d", r)
	}
}

func TestNewFromString(t *testing.T) {
	tcs := []struct {
		amount   string
		code     string
		expected int64
	}{
		{"1.25", USD, 125},
		{"100", USD, 10000},
		{"-0.50", EUR, -50},
		{"0.0001", USD, 0},
		{"1.235", BHD, 1235},
	}

	for _, tc := range tcs {
		m, err := NewFromString(tc.amount, tc.code)
		if err != nil {
			t.Errorf("Unexpected error for %q: %v", tc.amount, err)
			continue
		}
		r := m.Amount()

		if r != tc.expected {
			t.Errorf("Expected NewFromString(%q, %s) = %d got %d",
				tc.amount, tc.code, tc.expected, r)
		}
	}
}

func TestNewFromString_Invalid(t *testing.T) {
	_, err := NewFromString("not-a-number", USD)

	if err == nil {
		t.Error("Expected err")
	}
}

func TestNewFromDecimal(t *testing.T) {
	d1, _ := NewDecimalFromString("12.34")
	d2 := &Decimal{val: big.NewInt(7), exponent: 0}
	tcs := []struct {
		decimal  *Decimal
		code     string
		expected int64
	}{
		{d1, USD, 1234},
		{d2, USD, 700},
	}

	for _, tc := range tcs {
		m, err := NewFromDecimal(tc.decimal, tc.code)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			continue
		}
		r := m.Amount()

		if r != tc.expected {
			t.Errorf("Expected NewFromDecimal amount %d got %d", tc.expected, r)
		}
	}
}

func TestNewFromDecimal_Nil(t *testing.T) {
	_, err := NewFromDecimal(nil, USD)

	if err == nil {
		t.Error("Expected err")
	}
}
