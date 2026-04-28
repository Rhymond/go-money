package money

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		res, err := m.Multiply(tc.multiplier)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		r := res.amount.val.Int64()

		if r != tc.expected {
			t.Errorf("Expected %d * %d = %d got %d", tc.amount, tc.multiplier, tc.expected, r)
		}
	}
}

func TestMoney_MultiplyFloat(t *testing.T) {
	m := New(1000, USD) // $10.00, exponent 0
	r, err := m.Multiply(1.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.amount.val.Int64() != 15000 || r.amount.exponent != 1 {
		t.Errorf("expected 15000 e-1, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplySubUnit(t *testing.T) {
	// $0.01 * 1.5 = $0.015 — must preserve the extra precision instead of rounding.
	m := New(1, USD)
	r, err := m.Multiply(1.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.amount.val.Int64() != 15 || r.amount.exponent != 1 {
		t.Errorf("expected 15 e-1, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplyString(t *testing.T) {
	m := New(10000, USD) // $100.00
	r, err := m.Multiply("0.075")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 10000 * 75 = 750000, exponent 0 + 3 = 3 → 750.000 (in cents)
	if r.amount.val.Int64() != 750000 || r.amount.exponent != 3 {
		t.Errorf("expected 750000 e-3, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplyChained(t *testing.T) {
	m := New(1000, USD) // $10.00
	r, err := m.Multiply(2, "0.5", 1.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 1000 * 2 * 5 * 15, exponents 0 + 0 + 1 + 1 = 2 → 150000 e-2 = 1500.00 cents = $15.00
	if r.amount.val.Int64() != 150000 || r.amount.exponent != 2 {
		t.Errorf("expected 150000 e-2, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MultiplyDecimal(t *testing.T) {
	m := New(200, USD)
	d, _ := NewDecimal("2.5")
	r, err := m.Multiply(d)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.amount.val.Int64() != 5000 || r.amount.exponent != 1 {
		t.Errorf("expected 5000 e-1, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MustMultiply(t *testing.T) {
	m := New(1000, USD)
	r := m.MustMultiply(1.5)
	if r.amount.val.Int64() != 15000 || r.amount.exponent != 1 {
		t.Errorf("expected 15000 e-1, got %s e-%d", r.amount.val.String(), r.amount.exponent)
	}
}

func TestMoney_MustMultiplyPanics(t *testing.T) {
	m := New(100, USD)

	cases := []struct {
		name string
		args []any
	}{
		{"no multipliers", nil},
		{"unsupported type", []any{[]int{1, 2}}},
		{"invalid string", []any{"not-a-number"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for %s", tc.name)
				}
			}()
			m.MustMultiply(tc.args...)
		})
	}
}

func TestMoney_MultiplyErrors(t *testing.T) {
	m := New(100, USD)

	if _, err := m.Multiply(); err == nil {
		t.Error("expected error when no multipliers supplied")
	}

	if _, err := m.Multiply([]int{1, 2}); err == nil {
		t.Error("expected error for unsupported multiplier type")
	}

	if _, err := m.Multiply(2, "not-a-number"); err == nil {
		t.Error("expected error for invalid string multiplier")
	}
}

func TestMoney_Round(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
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
		m := New(tc.amount, EUR)
		r := m.Round().amount.val.Int64()

		if r != tc.expected {
			t.Errorf("Expected rounded %d to be %d got %d", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_RoundWithExponential(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
	}{
		{12555, 13000},
	}

	for _, tc := range tcs {
		AddCurrency("CUR", "*", "$1", ".", ",", 3)
		m := New(tc.amount, "CUR")
		r := m.Round().amount.val.Int64()

		if r != tc.expected {
			t.Errorf("Expected rounded %d to be %d got %d", tc.amount, tc.expected, r)
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
	}

	for _, tc := range tcs {
		m := New(tc.amount, EUR)
		var rs []int64
		ratios := make([]any, len(tc.ratios))
		for i, r := range tc.ratios {
			ratios[i] = r
		}
		split, _ := m.Allocate(ratios...)

		for _, party := range split {
			rs = append(rs, party.amount.val.Int64())
		}

		if !reflect.DeepEqual(tc.expected, rs) {
			t.Errorf("Expected allocation of %d for ratios %v to be %v got %v", tc.amount, tc.ratios,
				tc.expected, rs)
		}
	}
}

func TestMoney_AllocateFloatRatios(t *testing.T) {
	// $10.00 split 1.5 : 2.5 : 1.0 = 30% : 50% : 20% → $3.00, $5.00, $2.00
	m := New(1000, USD)
	parties, err := m.Allocate(1.5, 2.5, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int64{300, 500, 200}
	got := make([]int64, len(parties))
	for i, p := range parties {
		got[i] = p.amount.val.Int64()
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}

	// Allocations must sum back to the original amount.
	var sum int64
	for _, v := range got {
		sum += v
	}
	if sum != 1000 {
		t.Errorf("expected parties to sum to 1000, got %d", sum)
	}
}

func TestMoney_AllocateMixedRatios(t *testing.T) {
	// Mix int, float, string ratios — should behave the same as if all were
	// expressed as the same numeric type.
	m := New(1000, USD)
	parties, err := m.Allocate(1, 1.5, "2.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// sum = 5.0, shares = 0.2 / 0.3 / 0.5 → 200, 300, 500
	want := []int64{200, 300, 500}
	got := make([]int64, len(parties))
	for i, p := range parties {
		got[i] = p.amount.val.Int64()
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMoney_AllocateFloatLeftover(t *testing.T) {
	// Penny that doesn't divide evenly by float ratios — leftover must still
	// land on the first party so the total is preserved.
	m := New(5, USD)
	parties, err := m.Allocate(0.5, 0.25, 0.25)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var sum int64
	for _, p := range parties {
		sum += p.amount.val.Int64()
	}
	if sum != 5 {
		t.Errorf("expected sum 5, got %d", sum)
	}
}

func TestMoney_AllocateInvalidRatio(t *testing.T) {
	m := New(100, USD)
	if _, err := m.Allocate(1, "not-a-number"); err == nil {
		t.Error("expected error for invalid ratio")
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

func TestDefaultMarshal(t *testing.T) {
	given := New(12345, IQD)
	expected := `{"amount":12345,"currency":"IQD"}`

	b, err := json.Marshal(given)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s got %s", expected, string(b))
	}
}

func TestCustomMarshal(t *testing.T) {
	given := New(12345, IQD)
	expected := `{"amount":12345,"currency_code":"IQD","currency_fraction":3}`
	MarshalJSON = func(m Money) ([]byte, error) {
		buff := bytes.NewBufferString(fmt.Sprintf(`{"amount": %d, "currency_code": "%s", "currency_fraction": %d}`, m.Amount(), m.Currency().Code, m.Currency().Fraction))
		return buff.Bytes(), nil
	}

	b, err := json.Marshal(given)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s got %s", expected, string(b))
	}
}

func TestDefaultUnmarshal(t *testing.T) {
	given := `{"amount": 10012, "currency":"USD"}`
	expected := "$100.12"
	var m Money
	err := json.Unmarshal([]byte(given), &m)
	if err != nil {
		t.Error(err)
	}

	if m.Display() != expected {
		t.Errorf("Expected %s got %s", expected, m.Display())
	}
}

func TestCustomUnmarshal(t *testing.T) {
	given := `{"amount": 10012, "currency_code":"USD", "currency_fraction":2}`
	expected := "$100.12"
	UnmarshalJSON = func(m *Money, b []byte) error {
		data := make(map[string]interface{})
		err := json.Unmarshal(b, &data)
		if err != nil {
			return err
		}
		ref := New(int64(data["amount"].(float64)), data["currency_code"].(string))
		*m = *ref
		return nil
	}

	var m Money
	err := json.Unmarshal([]byte(given), &m)
	if err != nil {
		t.Error(err)
	}

	if m.Display() != expected {
		t.Errorf("Expected %s got %s", expected, m.Display())
	}
}
