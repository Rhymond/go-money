package money

import (
	"math"
	"math/big"
	"strings"
	"testing"
)

func TestNewDecimal_Integers(t *testing.T) {
	cases := []struct {
		name    string
		input   interface{}
		wantVal int64
	}{
		{"int", int(5), 5},
		{"int8", int8(-3), -3},
		{"int16", int16(1234), 1234},
		{"int32", int32(-99), -99},
		{"int64", int64(1 << 40), 1 << 40},
		{"uint", uint(7), 7},
		{"uint8", uint8(255), 255},
		{"uint16", uint16(65535), 65535},
		{"uint32", uint32(1 << 30), 1 << 30},
		{"uint64", uint64(1 << 40), 1 << 40},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := NewDecimal(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if d.exponent != 0 {
				t.Errorf("expected exponent 0, got %d", d.exponent)
			}
			if d.val.Int64() != tc.wantVal {
				t.Errorf("expected val %d, got %s", tc.wantVal, d.val.String())
			}
		})
	}
}

func TestNewDecimal_Float(t *testing.T) {
	d, err := NewDecimal(5.25)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.val.Int64() != 525 || d.exponent != 2 {
		t.Errorf("expected 525e-2, got %se-%d", d.val.String(), d.exponent)
	}

	d, err = NewDecimal(float32(-1.5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.val.Int64() != -15 || d.exponent != 1 {
		t.Errorf("expected -15e-1, got %se-%d", d.val.String(), d.exponent)
	}
}

func TestNewDecimal_String(t *testing.T) {
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
		d, err := NewDecimal(tc.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.input, err)
		}
		if d.val.String() != tc.wantVal || d.exponent != tc.wantExp {
			t.Errorf("input %q: expected %s e-%d, got %s e-%d",
				tc.input, tc.wantVal, tc.wantExp, d.val.String(), d.exponent)
		}
	}
}

func TestNewDecimal_FloatNaNInf(t *testing.T) {
	cases := []struct {
		name    string
		input   interface{}
		wantSub string
	}{
		{"NaN float64", math.NaN(), "NaN"},
		{"+Inf float64", math.Inf(1), "infinity"},
		{"-Inf float64", math.Inf(-1), "infinity"},
		{"NaN float32", float32(math.NaN()), "NaN"},
		{"+Inf float32", float32(math.Inf(1)), "infinity"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewDecimal(tc.input)
			if err == nil {
				t.Fatalf("expected error")
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Errorf("expected error to mention %q, got %q", tc.wantSub, err.Error())
			}
		})
	}
}

func TestNewDecimal_StringInvalid(t *testing.T) {
	for _, s := range []string{"", "   ", "abc", "1.2.3", "--5"} {
		if _, err := NewDecimal(s); err == nil {
			t.Errorf("expected error for %q", s)
		}
	}
}

func TestNewDecimal_BigInt(t *testing.T) {
	src := big.NewInt(42)
	d, err := NewDecimal(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	src.SetInt64(99)
	if d.val.Int64() != 42 {
		t.Errorf("expected deep copy, got mutated: %s", d.val.String())
	}
}

func TestNewDecimal_FromDecimal(t *testing.T) {
	src := &Decimal{val: big.NewInt(123), exponent: 2}
	d, err := NewDecimal(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	src.val.SetInt64(999)
	src.exponent = 5
	if d.val.Int64() != 123 || d.exponent != 2 {
		t.Errorf("expected deep copy, got %s e-%d", d.val.String(), d.exponent)
	}
}

func TestNewDecimal_FromMoney(t *testing.T) {
	m := New(500, EUR)
	d, err := NewDecimal(m)
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
}

func TestNewDecimal_Unsupported(t *testing.T) {
	if _, err := NewDecimal([]int{1, 2}); err == nil {
		t.Error("expected error for unsupported type")
	}
}

func TestCalculator_AddMismatchedExponents(t *testing.T) {
	a, _ := NewDecimal("1.5")     // 15 e-1
	b, _ := NewDecimal("1.005")   // 1005 e-3
	sum := mutate.calc.add(a, b)  // expect 2505 e-3

	if sum.val.String() != "2505" || sum.exponent != 3 {
		t.Errorf("expected 2505 e-3, got %s e-%d", sum.val.String(), sum.exponent)
	}
}

func TestCalculator_SubtractMismatchedExponents(t *testing.T) {
	a, _ := NewDecimal("2")
	b, _ := NewDecimal("0.25")
	diff := mutate.calc.subtract(a, b)

	if diff.val.String() != "175" || diff.exponent != 2 {
		t.Errorf("expected 175 e-2, got %s e-%d", diff.val.String(), diff.exponent)
	}
}
