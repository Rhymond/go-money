package money

import (
	"errors"
	"math"
	"math/big"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	m := New(1, EUR)

	if m.amount.val.Int64() != 1 {
		t.Errorf("Expected %d got %d", 1, m.amount.val.Int64())
	}

	if m.currency.Code != EUR {
		t.Errorf("Expected currency %s got %s", EUR, m.currency.Code)
	}

	m = New(-100, EUR)

	if m.amount.val.Int64() != -100 {
		t.Errorf("Expected %d got %d", -100, m.amount.val.Int64())
	}
}

func TestNew_WithUnregisteredCurrency(t *testing.T) {
	const currencyFooCode = "FOO"
	const expectedAmount = 100
	const expectedDisplay = "1.00FOO"

	m := New(100, currencyFooCode)

	if m.amount.val.Int64() != expectedAmount {
		t.Errorf("Expected amount %d got %d", expectedAmount, m.amount.val.Int64())
	}

	if m.currency.Code != currencyFooCode {
		t.Errorf("Expected currency code %s got %s", currencyFooCode, m.currency.Code)
	}

	if m.Display() != expectedDisplay {
		t.Errorf("Expected display %s got %s", expectedDisplay, m.Display())
	}
}

func TestCurrency(t *testing.T) {
	code := "MOCK"
	decimals := 5
	AddCurrency(code, "M$", "1 $", ".", ",", decimals)
	m := New(1, code)
	c := m.Currency().Code
	if c != code {
		t.Errorf("Expected %s got %s", code, c)
	}
	f := m.Currency().Fraction
	if f != decimals {
		t.Errorf("Expected %d got %d", decimals, f)
	}
}

func TestMoney_SameCurrency(t *testing.T) {
	m := New(0, EUR)
	om := New(0, USD)

	if m.SameCurrency(om) {
		t.Errorf("Expected %s not to be same as %s", m.currency.Code, om.currency.Code)
	}

	om = New(0, EUR)

	if !m.SameCurrency(om) {
		t.Errorf("Expected %s to be same as %s", m.currency.Code, om.currency.Code)
	}
}

func TestMoney_Equals(t *testing.T) {
	m := New(0, EUR)
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.amount, EUR)
		r, err := m.Equals(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equals %d == %t got %t", m.amount.val.Int64(),
				om.amount.val.Int64(), tc.expected, r)
		}
	}
}

func TestMoney_Equals_DifferentCurrencies(t *testing.T) {
	t.Parallel()

	eur := New(0, EUR)
	usd := New(0, USD)

	_, err := eur.Equals(usd)
	if err == nil || !errors.Is(ErrCurrencyMismatch, err) {
		t.Errorf("Expected Equals to return %q, got %v", ErrCurrencyMismatch.Error(), err)
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	m := New(0, EUR)
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.amount, EUR)
		r, err := m.GreaterThan(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Greater Than %d == %t got %t", m.amount.val.Int64(),
				om.amount.val.Int64(), tc.expected, r)
		}
	}
}

func TestMoney_GreaterThanOrEqual(t *testing.T) {
	m := New(0, EUR)
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, true},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.amount, EUR)
		r, err := m.GreaterThanOrEqual(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equals Or Greater Than %d == %t got %t", m.amount.val.Int64(),
				om.amount.val.Int64(), tc.expected, r)
		}
	}
}

func TestMoney_LessThan(t *testing.T) {
	m := New(0, EUR)
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}

	for _, tc := range tcs {
		om := New(tc.amount, EUR)
		r, err := m.LessThan(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Less Than %d == %t got %t", m.amount.val.Int64(),
				om.amount.val.Int64(), tc.expected, r)
		}
	}
}

func TestMoney_LessThanOrEqual(t *testing.T) {
	m := New(0, EUR)
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, true},
	}

	for _, tc := range tcs {
		om := New(tc.amount, EUR)
		r, err := m.LessThanOrEqual(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equal Or Less Than %d == %t got %t", m.amount.val.Int64(),
				om.amount.val.Int64(), tc.expected, r)
		}
	}
}

func TestMoney_IsZero(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		m := New(tc.amount, EUR)
		r := m.IsZero()

		if r != tc.expected {
			t.Errorf("Expected %d to be zero == %t got %t", m.amount.val.Int64(), tc.expected, r)
		}
	}
}

func TestMoney_IsNegative(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}

	for _, tc := range tcs {
		m := New(tc.amount, EUR)
		r := m.IsNegative()

		if r != tc.expected {
			t.Errorf("Expected %d to be negative == %t got %t", m.amount.val.Int64(),
				tc.expected, r)
		}
	}
}

func TestMoney_IsPositive(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}

	for _, tc := range tcs {
		m := New(tc.amount, EUR)
		r := m.IsPositive()

		if r != tc.expected {
			t.Errorf("Expected %d to be positive == %t got %t", m.amount.val.Int64(),
				tc.expected, r)
		}
	}
}

func TestMoney_Absolute(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
	}{
		{-1, 1},
		{0, 0},
		{1, 1},
	}

	for _, tc := range tcs {
		m := New(tc.amount, EUR)
		r := m.Absolute().amount.val.Int64()

		if r != tc.expected {
			t.Errorf("Expected absolute %d to be %d got %d", m.amount.val.Int64(),
				tc.expected, r)
		}
	}
}

func TestMoney_Negative(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
	}{
		{-1, -1},
		{0, -0},
		{1, -1},
	}

	for _, tc := range tcs {
		m := New(tc.amount, EUR)
		r := m.Negative().amount.val.Int64()

		if r != tc.expected {
			t.Errorf("Expected absolute %d to be %d got %d", m.amount.val.Int64(),
				tc.expected, r)
		}
	}
}

func TestMoney_Add(t *testing.T) {
	tcs := []struct {
		amount1  int64
		amount2  int64
		expected int64
	}{
		{5, 5, 10},
		{10, 5, 15},
		{1, -1, 0},
	}

	for _, tc := range tcs {
		m := New(tc.amount1, EUR)
		om := New(tc.amount2, EUR)
		r, err := m.Add(om)
		if err != nil {
			t.Error(err)
		}

		if r.Amount() != tc.expected {
			t.Errorf("Expected %d + %d = %d got %d", tc.amount1, tc.amount2,
				tc.expected, r.amount.val.Int64())
		}
	}
}

func TestMoney_Add2(t *testing.T) {
	m := New(100, EUR)
	dm := New(100, GBP)
	r, err := m.Add(dm)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Add3(t *testing.T) {
	tcs := []struct {
		amount1  int64
		amount2  int64
		amount3  int64
		expected int64
	}{
		{5, 5, 3, 13},
		{10, 5, 4, 19},
		{1, -1, 2, 2},
		{3, -1, -4, -2},
	}

	for _, tc := range tcs {
		mon1 := New(tc.amount1, EUR)
		mon2 := New(tc.amount2, EUR)
		mon3 := New(tc.amount3, EUR)
		r, err := mon1.Add(mon2, mon3)

		if err != nil {
			t.Error(err)
		}

		if r.Amount() != tc.expected {
			t.Errorf("Expected %d + %d + %d = %d got %d", tc.amount1, tc.amount2, tc.amount3,
				tc.expected, r.amount.val.Int64())
		}
	}
}

func TestMoney_Add4(t *testing.T) {
	m := New(100, EUR)
	r, err := m.Add()

	if err != nil {
		t.Error(err)
	}

	if r.amount.val.Int64() != 100 {
		t.Error("Expected amount to be 100")
	}
}

func TestMoney_Subtract(t *testing.T) {
	tcs := []struct {
		amount1  int64
		amount2  int64
		expected int64
	}{
		{5, 5, 0},
		{10, 5, 5},
		{1, -1, 2},
	}

	for _, tc := range tcs {
		m := New(tc.amount1, EUR)
		om := New(tc.amount2, EUR)
		r, err := m.Subtract(om)
		if err != nil {
			t.Error(err)
		}

		if r.amount.val.Int64() != tc.expected {
			t.Errorf("Expected %d - %d = %d got %d", tc.amount1, tc.amount2,
				tc.expected, r.amount.val.Int64())
		}
	}
}

func TestMoney_Subtract2(t *testing.T) {
	m := New(100, EUR)
	dm := New(100, GBP)
	r, err := m.Subtract(dm)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Subtract3(t *testing.T) {
	tcs := []struct {
		amount1  int64
		amount2  int64
		amount3  int64
		expected int64
	}{
		{5, 5, 3, -3},
		{10, -5, 4, 11},
		{1, -1, 2, 0},
		{7, 1, -4, 10},
	}

	for _, tc := range tcs {
		mon1 := New(tc.amount1, EUR)
		mon2 := New(tc.amount2, EUR)
		mon3 := New(tc.amount3, EUR)
		r, err := mon1.Subtract(mon2, mon3)

		if err != nil {
			t.Error(err)
		}

		if r.Amount() != tc.expected {
			t.Errorf("Expected (%d) - (%d) - (%d) = %d got %d", tc.amount1, tc.amount2, tc.amount3,
				tc.expected, r.amount.val.Int64())
		}
	}
}

func TestMoney_Subtract4(t *testing.T) {
	m := New(100, EUR)
	r, err := m.Subtract()

	if err != nil {
		t.Error(err)
	}

	if r.amount.val.Int64() != 100 {
		t.Error("Expected amount to be 100")
	}
}

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
		r := m.Multiply(NewDecimalFromInt(tc.multiplier)).amount.val.Int64()

		if r != tc.expected {
			t.Errorf("Expected %d * %d = %d got %d", tc.amount, tc.multiplier, tc.expected, r)
		}
	}
}

func TestMoney_MultiplyFloat(t *testing.T) {
	m := New(1000, USD) // $10.00, exponent 0
	mul, err := NewDecimalFromFloat(1.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.Multiply(mul)
	if r.amount.val.Int64() != 15000 || r.amount.exponent != 1 {
		t.Errorf("expected 15000 e-1, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

// Amount() and Display() must read through the exponent so that callers see
// $15.00 (not $150.00) after a fractional Multiply.
func TestMoney_AmountAfterMultiply(t *testing.T) {
	m := New(1000, USD) // $10.00
	mul, err := NewDecimalFromFloat(1.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.Multiply(mul)

	if got := r.Amount(); got != 1500 {
		t.Errorf("Amount(): expected 1500, got %d", got)
	}
	if got := r.Display(); got != "$15.00" {
		t.Errorf("Display(): expected $15.00, got %s", got)
	}
	if got := r.AsMajorUnits(); got != 15.00 {
		t.Errorf("AsMajorUnits(): expected 15.00, got %v", got)
	}
}

// Sub-smallest-unit precision truncates toward zero on read.
func TestMoney_AmountTruncatesSubCent(t *testing.T) {
	m := New(1, USD) // $0.01
	mul, _ := NewDecimalFromFloat(1.5)
	r := m.Multiply(mul) // val=15, exp=1 → 1.5 cents

	if got := r.Amount(); got != 1 {
		t.Errorf("Amount(): expected 1 (truncated), got %d", got)
	}
	if got := r.Display(); got != "$0.01" {
		t.Errorf("Display(): expected $0.01, got %s", got)
	}
}

func TestMoney_MultiplySubUnit(t *testing.T) {
	// $0.01 * 1.5 = $0.015 — must preserve the extra precision instead of rounding.
	m := New(1, USD)
	mul, err := NewDecimalFromFloat(1.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.Multiply(mul)
	if r.amount.val.Int64() != 15 || r.amount.exponent != 1 {
		t.Errorf("expected 15 e-1, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplyString(t *testing.T) {
	m := New(10000, USD) // $100.00
	mul, err := NewDecimalFromString("0.075")
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
	half, err := NewDecimalFromString("0.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	oneAndHalf, err := NewDecimalFromFloat(1.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.Multiply(NewDecimalFromInt(2), half, oneAndHalf)
	// 1000 * 2 * 5 * 15, exponents 0 + 0 + 1 + 1 = 2 → 150000 e-2 = 1500.00 cents = $15.00
	if r.amount.val.Int64() != 150000 || r.amount.exponent != 2 {
		t.Errorf("expected 150000 e-2, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplyDecimal(t *testing.T) {
	m := New(200, USD)
	d, _ := NewDecimalFromString("2.5")
	r := m.Multiply(d)
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

func TestMoney_MultiplyGenericIntegerTypes(t *testing.T) {
	m := New(100, USD)

	if r := m.Multiply(NewDecimalFromInt(int8(2))).amount.val.Int64(); r != 200 {
		t.Errorf("int8: got %d", r)
	}
	if r := m.Multiply(NewDecimalFromInt(uint64(3))).amount.val.Int64(); r != 300 {
		t.Errorf("uint64: got %d", r)
	}
	// float32(2.5) parses as 25 e-1, so 100 * (25 e-1) = 2500 e-1 = 250.
	mul, err := NewDecimalFromFloat(float32(2.5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.Multiply(mul)
	if r.amount.val.Int64() != 2500 || r.amount.exponent != 1 {
		t.Errorf("float32: got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

// A nil *Decimal multiplies to zero — "multiply by nothing returns nothing".
// Currency is preserved; subsequent multipliers are not applied.
func TestMoney_MultiplyNilDecimalReturnsZero(t *testing.T) {
	m := New(100, USD)
	r := m.Multiply((*Decimal)(nil))

	if got := r.Amount(); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
	if got := r.Currency().Code; got != USD {
		t.Errorf("expected currency USD, got %s", got)
	}

	// Nil short-circuits — the trailing 999 must not run.
	r = m.Multiply(NewDecimalFromInt(2), nil, NewDecimalFromInt(999))
	if got := r.Amount(); got != 0 {
		t.Errorf("expected 0 from short-circuit, got %d", got)
	}
}

// Round on Money created via New() (already at smallest-unit precision) is a
// no-op. Sub-smallest-unit precision introduced by Multiply collapses to the
// nearest smallest unit using half-away-from-zero.
func TestMoney_Round(t *testing.T) {
	noop := []int64{125, 175, 0, -1, -75}
	for _, amount := range noop {
		m := New(amount, EUR)
		r := m.Round().amount.val.Int64()
		if r != amount {
			t.Errorf("Round(%d, EUR) at exp=0: expected %d (no-op), got %d", amount, amount, r)
		}
	}

	subCent := []struct {
		val      int64
		exp      int
		expected int64
	}{
		{15, 1, 2},   // $0.015 → $0.02 (round up)
		{14, 1, 1},   // $0.014 → $0.01 (round down)
		{5, 1, 1},    // $0.005 → $0.01 (boundary, half away from zero)
		{4, 1, 0},    // $0.004 → $0.00
		{-15, 1, -2}, // -$0.015 → -$0.02
		{-5, 1, -1},  // -$0.005 → -$0.01 (half away from zero, negative)
		{0, 1, 0},
		{1234, 3, 1}, // 1.234 cents → 1 cent
		{1500, 3, 2}, // 1.500 cents → 2 cents (boundary)
	}
	for _, tc := range subCent {
		m := &Money{
			amount:   &Decimal{val: big.NewInt(tc.val), exponent: tc.exp},
			currency: newCurrency(USD).get(),
		}
		r := m.Round()
		if r.amount.val.Int64() != tc.expected || r.amount.exponent != 0 {
			t.Errorf("Round(val=%d exp=%d): expected %d e-0, got %s e-%d",
				tc.val, tc.exp, tc.expected, r.amount.val.String(), r.amount.exponent)
		}
	}
}

// RoundToMajor rounds val to the nearest major-unit boundary using
// half-away-from-zero. This is what Round did before v2's exponent-aware
// refactor; preserved for callers that want presentation-grade summaries.
func TestMoney_RoundToMajor(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
	}{
		{125, 100},
		{175, 200},
		{349, 300},
		{351, 400},
		{50, 100}, // boundary: 0.50 rounds away from zero
		{150, 200},
		{0, 0},
		{-1, 0},
		{-50, -100}, // boundary: -0.50 rounds away from zero
		{-75, -100},
	}

	for _, tc := range tcs {
		m := New(tc.amount, EUR)
		r := m.RoundToMajor().amount.val.Int64()

		if r != tc.expected {
			t.Errorf("Expected RoundToMajor(%d) to be %d got %d", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_RoundToMajorWithExponential(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
	}{
		{12555, 13000},
		{12500, 13000}, // boundary
	}

	for _, tc := range tcs {
		AddCurrency("CUR", "*", "$1", ".", ",", 3)
		m := New(tc.amount, "CUR")
		r := m.RoundToMajor().amount.val.Int64()

		if r != tc.expected {
			t.Errorf("Expected RoundToMajor(%d) to be %d got %d", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_Split(t *testing.T) {
	tcs := []struct {
		amount   int64
		split    int
		expected []int64
	}{
		{100, 3, []int64{34, 33, 33}},
		{100, 4, []int64{25, 25, 25, 25}},
		{5, 3, []int64{2, 2, 1}},
		{-101, 4, []int64{-26, -25, -25, -25}},
		{-101, 4, []int64{-26, -25, -25, -25}},
		{-2, 3, []int64{-1, -1, 0}},
	}

	for _, tc := range tcs {
		m := New(tc.amount, EUR)
		var rs []int64
		split, _ := m.Split(tc.split)

		for _, party := range split {
			rs = append(rs, party.amount.val.Int64())
		}

		if !reflect.DeepEqual(tc.expected, rs) {
			t.Errorf("Expected split of %d to be %v got %v", tc.amount, tc.expected, rs)
		}
	}
}

func TestMoney_Split2(t *testing.T) {
	m := New(100, EUR)
	r, err := m.Split(-10)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Allocate(t *testing.T) {
	tcs := []struct {
		amount   int64
		ratios   []int
		expected []int64
	}{
		{100, []int{50, 50}, []int64{50, 50}},
		{100, []int{30, 30, 30}, []int64{34, 33, 33}},
		{200, []int{25, 25, 50}, []int64{50, 50, 100}},
		{5, []int{50, 25, 25}, []int64{3, 1, 1}},
		{0, []int{0, 0, 0, 0}, []int64{0, 0, 0, 0}},
		{0, []int{50, 10}, []int64{0, 0}},
		{10, []int{0, 100}, []int64{0, 10}},
		{10, []int{0, 0}, []int64{0, 0}},
	}

	for _, tc := range tcs {
		m := New(tc.amount, EUR)
		var rs []int64
		split, _ := m.Allocate(tc.ratios...)

		for _, party := range split {
			rs = append(rs, party.amount.val.Int64())
		}

		if !reflect.DeepEqual(tc.expected, rs) {
			t.Errorf("Expected allocation of %d for ratios %v to be %v got %v", tc.amount, tc.ratios,
				tc.expected, rs)
		}
	}
}

func TestMoney_Allocate2(t *testing.T) {
	m := New(100, EUR)
	r, err := m.Allocate()

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Format(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected string
	}{
		{100, GBP, "£1.00"},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		r := m.Display()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_Display(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected string
	}{
		{100, AED, "1.00 .\u062f.\u0625"},
		{1, USD, "$0.01"},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		r := m.Display()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_AsMajorUnits(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected float64
	}{
		{100, AED, 1.00},
		{1, USD, 0.01},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		r := m.AsMajorUnits()

		if r != tc.expected {
			t.Errorf("Expected value as major units of %d to be %f got %f", tc.amount, tc.expected, r)
		}
	}
}

func TestAllocateOverflow(t *testing.T) {
	m := New(math.MaxInt64, EUR)
	_, err := m.Allocate(math.MaxInt, 1)
	if err == nil {
		t.Fatalf("expected an error, but got nil")
	}

	expectedErrorMessage := "sum of given ratios exceeds max int"
	if err.Error() != expectedErrorMessage {
		t.Fatalf("expected error message %q, but got %q", expectedErrorMessage, err.Error())
	}
}

func TestMoney_Allocate3(t *testing.T) {
	pound := New(100, GBP)
	parties, err := pound.Allocate(33, 33, 33)
	if err != nil {
		t.Error(err)
	}

	if parties[0].Display() != "£0.34" {
		t.Errorf("Expected %s got %s", "£0.34", parties[0].Display())
	}

	if parties[1].Display() != "£0.33" {
		t.Errorf("Expected %s got %s", "£0.33", parties[1].Display())
	}

	if parties[2].Display() != "£0.33" {
		t.Errorf("Expected %s got %s", "£0.33", parties[2].Display())
	}
}

func TestMoney_Comparison(t *testing.T) {
	pound := New(100, GBP)
	twoPounds := New(200, GBP)
	twoEuros := New(200, EUR)

	if r, err := pound.GreaterThan(twoPounds); err != nil || r {
		t.Errorf("Expected %d Greater Than %d == %t got %t", pound.amount.val.Int64(),
			twoPounds.amount.val.Int64(), false, r)
	}

	if r, err := pound.LessThan(twoPounds); err != nil || !r {
		t.Errorf("Expected %d Less Than %d == %t got %t", pound.amount.val.Int64(),
			twoPounds.amount.val.Int64(), true, r)
	}

	if r, err := pound.LessThan(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.GreaterThan(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.Equals(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.LessThanOrEqual(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.GreaterThanOrEqual(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := twoPounds.Compare(pound); r != 1 && err != nil {
		t.Errorf("Expected %d Greater Than %d == %d got %d", pound.amount,
			twoPounds.amount, 1, r)
	}

	if r, err := pound.Compare(twoPounds); r != -1 && err != nil {
		t.Errorf("Expected %d Less Than %d == %d got %d", pound.amount,
			twoPounds.amount, -1, r)
	}

	if _, err := pound.Compare(twoEuros); err != ErrCurrencyMismatch {
		t.Error("Expected err")
	}

	anotherTwoEuros := New(200, EUR)
	if r, err := twoEuros.Compare(anotherTwoEuros); r != 0 && err != nil {
		t.Errorf("Expected %d Equals to %d == %d got %d", anotherTwoEuros.amount,
			twoEuros.amount, 0, r)
	}
}

func TestMoney_Currency(t *testing.T) {
	pound := New(100, GBP)

	if pound.Currency().Code != GBP {
		t.Errorf("Expected %s got %s", GBP, pound.Currency().Code)
	}
}

func TestMoney_Amount(t *testing.T) {
	pound := New(100, GBP)

	if pound.Amount() != 100 {
		t.Errorf("Expected %d got %d", 100, pound.Amount())
	}
}

// Every read/write/comparison method on the zero-value Money{} must succeed
// without panicking. Pre-fix, only Amount() and the marshalers handled
// nil-field state; everything else nil-derefed.
func TestMoney_ZeroValue_NoPanic(t *testing.T) {
	other := New(100, USD)

	checks := []struct {
		name string
		fn   func(*Money)
	}{
		{"Amount", func(m *Money) { _ = m.Amount() }},
		{"Currency", func(m *Money) { _ = m.Currency().Code }},
		{"Display", func(m *Money) { _ = m.Display() }},
		{"AsMajorUnits", func(m *Money) { _ = m.AsMajorUnits() }},
		{"IsZero", func(m *Money) { _ = m.IsZero() }},
		{"IsPositive", func(m *Money) { _ = m.IsPositive() }},
		{"IsNegative", func(m *Money) { _ = m.IsNegative() }},
		{"SameCurrency", func(m *Money) { _ = m.SameCurrency(other) }},
		{"Absolute", func(m *Money) { _ = m.Absolute() }},
		{"Negative", func(m *Money) { _ = m.Negative() }},
		{"Round", func(m *Money) { _ = m.Round() }},
		{"RoundToMajor", func(m *Money) { _ = m.RoundToMajor() }},
		{"Multiply", func(m *Money) { _ = m.Multiply(NewDecimalFromInt(2)) }},
		{"Add empty", func(m *Money) { _, _ = m.Add() }},
		{"Subtract empty", func(m *Money) { _, _ = m.Subtract() }},
		{"Equals zero", func(m *Money) { _, _ = m.Equals(&Money{}) }},
		{"Compare zero", func(m *Money) { _, _ = m.Compare(&Money{}) }},
		{"Split", func(m *Money) { _, _ = m.Split(2) }},
		{"Allocate", func(m *Money) { _, _ = m.Allocate(1, 1) }},
		{"Value", func(m *Money) { _, _ = m.Value() }},
	}

	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Money{}.%s panicked: %v", c.name, r)
				}
			}()
			var z Money
			c.fn(&z)
		})
	}

	// The semantic checks: zero Money should report zero/empty.
	var z Money
	if !z.IsZero() {
		t.Errorf("Money{}.IsZero() = false, want true")
	}
	if z.Amount() != 0 {
		t.Errorf("Money{}.Amount() = %d, want 0", z.Amount())
	}
	if got := z.Currency().Code; got != "" {
		t.Errorf("Money{}.Currency().Code = %q, want \"\"", got)
	}
}

// Money.Currency() must return a defensive copy — mutating fields on the
// returned *Currency must not bleed through to the global registry or to
// other Money instances created with the same code.
func TestMoney_Currency_DefensiveCopy(t *testing.T) {
	m := New(100, USD)
	c := m.Currency()
	originalFraction := c.Fraction
	c.Fraction = 99
	c.Code = "MUTATED"

	again := New(100, USD)
	if got := again.Currency().Fraction; got != originalFraction {
		t.Errorf("Mutation of caller's *Currency leaked to registry: Fraction=%d, want %d", got, originalFraction)
	}
	if got := again.Currency().Code; got != USD {
		t.Errorf("Mutation of caller's *Currency leaked to registry: Code=%q, want %q", got, USD)
	}
}

// GetCurrency must return a defensive copy too.
func TestGetCurrency_DefensiveCopy(t *testing.T) {
	c := GetCurrency(USD)
	if c == nil {
		t.Fatal("expected non-nil USD")
	}
	originalFraction := c.Fraction
	c.Fraction = 99

	again := GetCurrency(USD)
	if again.Fraction != originalFraction {
		t.Errorf("Mutation of GetCurrency() result leaked: Fraction=%d, want %d", again.Fraction, originalFraction)
	}
}
