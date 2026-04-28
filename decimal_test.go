package money

import (
	"math"
	"math/big"
	"strings"
	"testing"
)

func TestNewDecimalFromInt(t *testing.T) {
	checks := []struct {
		name    string
		got     *Decimal
		wantVal int64
	}{
		{"int", NewDecimalFromInt(int(5)), 5},
		{"int8", NewDecimalFromInt(int8(-3)), -3},
		{"int16", NewDecimalFromInt(int16(1234)), 1234},
		{"int32", NewDecimalFromInt(int32(-99)), -99},
		{"int64", NewDecimalFromInt(int64(1 << 40)), 1 << 40},
		{"uint", NewDecimalFromInt(uint(7)), 7},
		{"uint8", NewDecimalFromInt(uint8(255)), 255},
		{"uint16", NewDecimalFromInt(uint16(65535)), 65535},
		{"uint32", NewDecimalFromInt(uint32(1 << 30)), 1 << 30},
		{"uint64", NewDecimalFromInt(uint64(1 << 40)), 1 << 40},
	}

	for _, tc := range checks {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got.exponent != 0 {
				t.Errorf("expected exponent 0, got %d", tc.got.exponent)
			}
			if tc.got.val.Int64() != tc.wantVal {
				t.Errorf("expected val %d, got %s", tc.wantVal, tc.got.val.String())
			}
		})
	}
}

func TestNewDecimalFromFloat(t *testing.T) {
	d, err := NewDecimalFromFloat(5.25)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.val.Int64() != 525 || d.exponent != 2 {
		t.Errorf("expected 525e-2, got %se-%d", d.val.String(), d.exponent)
	}

	d, err = NewDecimalFromFloat(float32(-1.5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.val.Int64() != -15 || d.exponent != 1 {
		t.Errorf("expected -15e-1, got %se-%d", d.val.String(), d.exponent)
	}
}

func TestNewDecimalFromString(t *testing.T) {
	cases := []struct {
		input   string
		wantVal string
		wantExp int
	}{
		{"3.14159", "314159", 5},
		{"-0.001", "-1", 3},
		{"42", "42", 0},
		{"+10.5", "105", 1},
		{".5", "5", 1},
		{"1000000000000000000000", "1000000000000000000000", 0},
	}
	for _, tc := range cases {
		d, err := NewDecimalFromString(tc.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.input, err)
		}
		if d.val.String() != tc.wantVal || d.exponent != tc.wantExp {
			t.Errorf("input %q: expected %s e-%d, got %s e-%d",
				tc.input, tc.wantVal, tc.wantExp, d.val.String(), d.exponent)
		}
	}
}

func TestNewDecimalFromFloat_NaNInf(t *testing.T) {
	type tc struct {
		name    string
		fn      func() (*Decimal, error)
		wantSub string
	}
	cases := []tc{
		{"NaN float64", func() (*Decimal, error) { return NewDecimalFromFloat(math.NaN()) }, "NaN"},
		{"+Inf float64", func() (*Decimal, error) { return NewDecimalFromFloat(math.Inf(1)) }, "infinity"},
		{"-Inf float64", func() (*Decimal, error) { return NewDecimalFromFloat(math.Inf(-1)) }, "infinity"},
		{"NaN float32", func() (*Decimal, error) { return NewDecimalFromFloat(float32(math.NaN())) }, "NaN"},
		{"+Inf float32", func() (*Decimal, error) { return NewDecimalFromFloat(float32(math.Inf(1))) }, "infinity"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := c.fn()
			if err == nil {
				t.Fatalf("expected error")
			}
			if !strings.Contains(err.Error(), c.wantSub) {
				t.Errorf("expected error to mention %q, got %q", c.wantSub, err.Error())
			}
		})
	}
}

func TestNewDecimalFromString_Invalid(t *testing.T) {
	for _, s := range []string{"", "   ", "abc", "1.2.3", "--5"} {
		if _, err := NewDecimalFromString(s); err == nil {
			t.Errorf("expected error for %q", s)
		}
	}
}

func TestNewDecimalFromBigInt(t *testing.T) {
	src := big.NewInt(42)
	d, err := NewDecimalFromBigInt(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	src.SetInt64(99)
	if d.val.Int64() != 42 {
		t.Errorf("expected deep copy, got mutated: %s", d.val.String())
	}

	if _, err := NewDecimalFromBigInt(nil); err == nil {
		t.Error("expected error for nil *big.Int")
	}
}

func TestNewDecimalFromMoney(t *testing.T) {
	m := New(500, EUR)
	d, err := NewDecimalFromMoney(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.val.Int64() != 500 || d.exponent != 0 {
		t.Errorf("expected 500 e-0, got %s e-%d", d.val.String(), d.exponent)
	}
	m.amount.val.SetInt64(1)
	if d.val.Int64() != 500 {
		t.Error("expected deep copy from Money, got mutation")
	}

	if _, err := NewDecimalFromMoney(nil); err == nil {
		t.Error("expected error for nil *Money")
	}
}

func TestCalculator_AddMismatchedExponents(t *testing.T) {
	a, _ := NewDecimalFromString("1.5")   // 15 e-1
	b, _ := NewDecimalFromString("1.005") // 1005 e-3
	sum := mutate.calc.add(a, b)          // expect 2505 e-3

	if sum.val.String() != "2505" || sum.exponent != 3 {
		t.Errorf("expected 2505 e-3, got %s e-%d", sum.val.String(), sum.exponent)
	}
}

func TestCalculator_SubtractMismatchedExponents(t *testing.T) {
	a, _ := NewDecimalFromString("2")
	b, _ := NewDecimalFromString("0.25")
	diff := mutate.calc.subtract(a, b)

	if diff.val.String() != "175" || diff.exponent != 2 {
		t.Errorf("expected 175 e-2, got %s e-%d", diff.val.String(), diff.exponent)
	}
}
