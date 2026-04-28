package money

import (
	"math/big"
	"testing"
)

func TestNewFromFixedPoint_Issue151Example(t *testing.T) {
	// €1.25 expressed as nanos — the example from issue #151.
	m := NewFromFixedPoint(1, 250_000_000, FixedPointSizeNano, EUR)
	if got := m.Amount(); got != 125 {
		t.Errorf("expected 125 cents, got %d", got)
	}
	if got := m.Display(); got != "€1.25" {
		t.Errorf("expected €1.25, got %s", got)
	}
}

func TestNewFromFixedPoint_AllSizes(t *testing.T) {
	cases := []struct {
		name     string
		units    int64
		fraction int64
		size     uint
		code     string
		want     int64 // smallest currency unit
	}{
		{"deci",  1, 5, FixedPointSizeDeci,  USD, 150},
		{"centi", 1, 25, FixedPointSizeCenti, USD, 125},
		{"milli", 1, 250, FixedPointSizeMilli, USD, 125},
		{"micro", 1, 250_000, FixedPointSizeMicro, USD, 125},
		{"nano",  1, 250_000_000, FixedPointSizeNano, USD, 125},
		{"jpy_no_fraction_truncates", 100, 999_999_999, FixedPointSizeNano, JPY, 100},
		{"bhd_3_fraction", 1, 250_000_000, FixedPointSizeNano, BHD, 1250},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewFromFixedPoint(tc.units, tc.fraction, tc.size, tc.code)
			if got := m.Amount(); got != tc.want {
				t.Errorf("expected %d, got %d", tc.want, got)
			}
		})
	}
}

func TestNewFromFixedPoint_Negative(t *testing.T) {
	// $-1.75 per Google's spec: units=-1, fraction=-750_000_000.
	m := NewFromFixedPoint(-1, -750_000_000, FixedPointSizeNano, USD)
	if got := m.Amount(); got != -175 {
		t.Errorf("expected -175 cents, got %d", got)
	}
}

func TestNewFromFixedPoint_SubUnitTruncates(t *testing.T) {
	// One nano-euro is well below a cent — should round to zero.
	m := NewFromFixedPoint(0, 1, FixedPointSizeNano, EUR)
	if got := m.Amount(); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestNewFromMajorUnits(t *testing.T) {
	cases := []struct {
		amount int64
		code   string
		want   int64
	}{
		{5, USD, 500},     // $5.00 → 500 cents
		{1, JPY, 1},       // ¥1 → 1 (no fraction)
		{1, BHD, 1000},    // 1 BHD → 1000 (3 fraction digits)
		{-3, EUR, -300},
	}

	for _, tc := range cases {
		m := NewFromMajorUnits(tc.amount, tc.code)
		if got := m.Amount(); got != tc.want {
			t.Errorf("NewFromMajorUnits(%d, %s) = %d, want %d", tc.amount, tc.code, got, tc.want)
		}
	}
}

// Generic NewFromMajorUnits should accept every integer width / signedness
// and produce equivalent Money.
func TestNewFromMajorUnits_AllIntegerTypes(t *testing.T) {
	want := New(500, USD).Amount()

	if got := NewFromMajorUnits(int(5), USD).Amount(); got != want {
		t.Errorf("int: got %d, want %d", got, want)
	}
	if got := NewFromMajorUnits(int8(5), USD).Amount(); got != want {
		t.Errorf("int8: got %d, want %d", got, want)
	}
	if got := NewFromMajorUnits(int16(5), USD).Amount(); got != want {
		t.Errorf("int16: got %d, want %d", got, want)
	}
	if got := NewFromMajorUnits(int32(5), USD).Amount(); got != want {
		t.Errorf("int32: got %d, want %d", got, want)
	}
	if got := NewFromMajorUnits(int64(5), USD).Amount(); got != want {
		t.Errorf("int64: got %d, want %d", got, want)
	}
	if got := NewFromMajorUnits(uint(5), USD).Amount(); got != want {
		t.Errorf("uint: got %d, want %d", got, want)
	}
	if got := NewFromMajorUnits(uint8(5), USD).Amount(); got != want {
		t.Errorf("uint8: got %d, want %d", got, want)
	}
	if got := NewFromMajorUnits(uint16(5), USD).Amount(); got != want {
		t.Errorf("uint16: got %d, want %d", got, want)
	}
	if got := NewFromMajorUnits(uint32(5), USD).Amount(); got != want {
		t.Errorf("uint32: got %d, want %d", got, want)
	}
	if got := NewFromMajorUnits(uint64(5), USD).Amount(); got != want {
		t.Errorf("uint64: got %d, want %d", got, want)
	}
}

// uint64 values above math.MaxInt64 must not be reinterpreted as negative.
// We use NewFromFixedPoint so the result still fits in an int64 amount after
// truncation by the size exponent.
func TestNewFromFixedPoint_LargeUint64(t *testing.T) {
	// units = 2^63 (one above MaxInt64) at nano size, truncated to whole
	// dollars: expected cents = (2^63 * 100) / 1 = 2^63 * 100 — overflows
	// int64 cents. Use a smaller value that exercises the unsigned path
	// without overflowing the final Amount().
	const big uint64 = 1 << 40 // ~1.1 trillion units, fits in int64 cents
	m := NewFromFixedPoint(big, uint64(0), FixedPointSizeCenti, USD)
	want := int64(big) * 100
	if got := m.Amount(); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestNewFromFloat(t *testing.T) {
	m, err := NewFromFloat(1.25, USD)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Amount(); got != 125 {
		t.Errorf("expected 125, got %d", got)
	}

	// Negative.
	m, err = NewFromFloat(-0.99, USD)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Amount(); got != -99 {
		t.Errorf("expected -99, got %d", got)
	}

	// Sub-cent truncates.
	m, _ = NewFromFloat(0.014, USD)
	if got := m.Amount(); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

// TestNewFromFloat_IEEE754Pitfalls pins the exact values from issues #121 and
// #124 so a future change to the float path can't silently re-introduce the
// off-by-one-cent bug. Each of these values has a binary representation that
// undershoots the displayed decimal — the master-branch implementation
// (int64(amount * 100)) returned one cent low because it multiplied the
// binary float and truncated. The v2 path goes float → shortest-round-trip
// string → decimal, which preserves the user's intended value exactly.
func TestNewFromFloat_IEEE754Pitfalls(t *testing.T) {
	cases := []struct {
		input float64
		code  string
		want  int64 // smallest currency unit
		issue string
	}{
		{1.15, USD, 115, "#121"},
		{136.98, USD, 13698, "#124"},
		{18.99, USD, 1899, "#124 (kylebragger comment)"},
		{73708.43, EUR, 7370843, "#124 (vaihtovirta comment)"},
		{0.1 + 0.2, USD, 30, "classic float-add pitfall"},
	}

	for _, tc := range cases {
		m, err := NewFromFloat(tc.input, tc.code)
		if err != nil {
			t.Fatalf("%v (%s): unexpected error: %v", tc.input, tc.issue, err)
		}
		if got := m.Amount(); got != tc.want {
			t.Errorf("%v (%s): got %d, want %d", tc.input, tc.issue, got, tc.want)
		}
	}
}

func TestNewFromFloat_Float32(t *testing.T) {
	m, err := NewFromFloat(float32(1.5), USD)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Amount(); got != 150 {
		t.Errorf("expected 150, got %d", got)
	}
}

func TestNewFromString(t *testing.T) {
	cases := []struct {
		s    string
		code string
		want int64
	}{
		{"1.25", USD, 125},
		{"100", USD, 10000},
		{"-0.50", EUR, -50},
		{"0.0001", USD, 0}, // sub-cent truncates
		{"1.235", BHD, 1235}, // BHD has 3 fraction digits
	}

	for _, tc := range cases {
		m, err := NewFromString(tc.s, tc.code)
		if err != nil {
			t.Fatalf("NewFromString(%q, %s) unexpected error: %v", tc.s, tc.code, err)
		}
		if got := m.Amount(); got != tc.want {
			t.Errorf("NewFromString(%q, %s) = %d, want %d", tc.s, tc.code, got, tc.want)
		}
	}
}

func TestNewFromString_Invalid(t *testing.T) {
	if _, err := NewFromString("not-a-number", USD); err == nil {
		t.Error("expected error")
	}
}

func TestNewFromDecimal(t *testing.T) {
	d, _ := NewDecimal("12.34")
	m, err := NewFromDecimal(d, USD)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Amount(); got != 1234 {
		t.Errorf("expected 1234, got %d", got)
	}

	// Decimal with exponent smaller than currency.Fraction → scale up.
	d2 := &Decimal{val: big.NewInt(7), exponent: 0} // = 7 (major units)
	m, err = NewFromDecimal(d2, USD)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Amount(); got != 700 {
		t.Errorf("expected 700, got %d", got)
	}
}

func TestNewFromDecimal_Nil(t *testing.T) {
	if _, err := NewFromDecimal(nil, USD); err == nil {
		t.Error("expected error for nil *Decimal")
	}
}

// Constructors are interchangeable: every path that produces "$1.25" should
// give an identical Money.
func TestConstructors_Equivalence(t *testing.T) {
	want := New(125, USD)

	mFromString, err := NewFromString("1.25", USD)
	if err != nil {
		t.Fatalf("NewFromString: %v", err)
	}
	mFromFloat, err := NewFromFloat(1.25, USD)
	if err != nil {
		t.Fatalf("NewFromFloat: %v", err)
	}
	d, err := NewDecimal("1.25")
	if err != nil {
		t.Fatalf("NewDecimal: %v", err)
	}
	mFromDecimal, err := NewFromDecimal(d, USD)
	if err != nil {
		t.Fatalf("NewFromDecimal: %v", err)
	}

	candidates := map[string]*Money{
		"NewFromFixedPoint": NewFromFixedPoint(1, 250_000_000, FixedPointSizeNano, USD),
		"NewFromString":     mFromString,
		"NewFromFloat":      mFromFloat,
		"NewFromDecimal":    mFromDecimal,
	}

	for name, m := range candidates {
		eq, err := want.Equals(m)
		if err != nil || !eq {
			t.Errorf("%s produced amount %d, expected %d", name, m.Amount(), want.Amount())
		}
	}
}
