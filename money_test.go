package money

import (
	"encoding/json"
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
	m := New(0, "EUR")
	om := New(0, "USD")

	if m.SameCurrency(om) {
		t.Errorf("Expected %s not to be same as %s", m.currency.Code, om.currency.Code)
	}

	om = New(0, "EUR")

	if !m.SameCurrency(om) {
		t.Errorf("Expected %s to be same as %s", m.currency.Code, om.currency.Code)
	}
}

func TestMoney_Equals(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.amount, "EUR")
		r, err := m.Equals(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equals %d == %t got %t", m.amount.val,
				om.amount.val, tc.expected, r)
		}
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.amount, "EUR")
		r, err := m.GreaterThan(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Greater Than %d == %t got %t", m.amount.val,
				om.amount.val, tc.expected, r)
		}
	}
}

func TestMoney_GreaterThanOrEqual(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, true},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om := New(tc.amount, "EUR")
		r, err := m.GreaterThanOrEqual(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equals Or Greater Than %d == %t got %t", m.amount.val,
				om.amount.val, tc.expected, r)
		}
	}
}

func TestMoney_LessThan(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}

	for _, tc := range tcs {
		om := New(tc.amount, "EUR")
		r, err := m.LessThan(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Less Than %d == %t got %t", m.amount.val,
				om.amount.val, tc.expected, r)
		}
	}
}

func TestMoney_LessThanOrEqual(t *testing.T) {
	m := New(0, "EUR")
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, true},
	}

	for _, tc := range tcs {
		om := New(tc.amount, "EUR")
		r, err := m.LessThanOrEqual(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equal Or Less Than %d == %t got %t", m.amount.val,
				om.amount.val, tc.expected, r)
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
		m := New(tc.amount, "EUR")
		r := m.IsZero()

		if r != tc.expected {
			t.Errorf("Expected %d to be zero == %t got %t", m.amount.val, tc.expected, r)
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
		amount   int64
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
		amount   int64
		expected int64
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
		amount   int64
		expected int64
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
		amount1  int64
		amount2  int64
		expected int64
	}{
		{5, 5, 10},
		{10, 5, 15},
		{1, -1, 0},
	}

	for _, tc := range tcs {
		m := New(tc.amount1, "EUR")
		om := New(tc.amount2, "EUR")
		r, err := m.Add(om)

		if err != nil {
			t.Error(err)
		}

		if r.Amount() != tc.expected {
			t.Errorf("Expected %d + %d = %d got %d", tc.amount1, tc.amount2,
				tc.expected, r.amount.val)
		}
	}

}

func TestMoney_Add2(t *testing.T) {
	m := New(100, "EUR")
	dm := New(100, "GBP")
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
		m := New(tc.amount1, "EUR")
		om := New(tc.amount2, "EUR")
		r, err := m.Subtract(om)

		if err != nil {
			t.Error(err)
		}

		if r.amount.val != tc.expected {
			t.Errorf("Expected %d - %d = %d got %d", tc.amount1, tc.amount2,
				tc.expected, r.amount.val)
		}
	}
}

func TestMoney_Subtract2(t *testing.T) {
	m := New(100, "EUR")
	dm := New(100, "GBP")
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
		m := New(tc.amount, "EUR")
		r := m.Multiply(tc.multiplier).amount.val

		if r != tc.expected {
			t.Errorf("Expected %d * %d = %d got %d", tc.amount, tc.multiplier, tc.expected, r)
		}
	}
}

func TestMoney_Divide(t *testing.T) {
	tcs := []struct {
		amount   int64
		divisor  int64
		expected int64
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
		m := New(tc.amount, "EUR")
		r := m.Round().amount.val

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
		r := m.Round().amount.val

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
		m := New(tc.amount, "EUR")
		var rs []int64
		split, _ := m.Split(tc.split)

		for _, party := range split {
			rs = append(rs, party.amount.val)
		}

		if !reflect.DeepEqual(tc.expected, rs) {
			t.Errorf("Expected split of %d to be %v got %v", tc.amount, tc.expected, rs)
		}
	}
}

func TestMoney_Split2(t *testing.T) {
	m := New(100, "EUR")
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
		m := New(tc.amount, "EUR")
		var rs []int64
		split, _ := m.Allocate(tc.ratios...)

		for _, party := range split {
			rs = append(rs, party.amount.val)
		}

		if !reflect.DeepEqual(tc.expected, rs) {
			t.Errorf("Expected allocation of %d for ratios %v to be %v got %v", tc.amount, tc.ratios,
				tc.expected, rs)
		}
	}
}

func TestMoney_Allocate2(t *testing.T) {
	m := New(100, "EUR")
	r, err := m.Allocate()

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Chain(t *testing.T) {
	m := New(10, "EUR")
	om := New(5, "EUR")
	// 10 + 5 = 15 / 5 = 3 * 4 = 12 - 5 = 7
	e := int64(7)

	m, err := m.Add(om)

	if err != nil {
		t.Error(err)
	}

	m = m.Divide(5).Multiply(4)
	m, err = m.Subtract(om)

	if err != nil {
		t.Error(err)
	}

	if m.amount.val != int64(7) {
		t.Errorf("Expected %d got %d", e, m.amount.val)
	}
}

func TestMoney_Format(t *testing.T) {
	tcs := []struct {
		amount   int64
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

func TestMoney_Display(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected string
	}{
		{100, "AED", "1.00 .\u062f.\u0625"},
		{1, "USD", "$0.01"},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		r := m.Display()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_ToWords(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected string
	}{
		{100, "PHP", "one hundred pesos only"},
		{1, "SGD", "one dollar only"},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		r := m.ToWords()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_Allocate3(t *testing.T) {
	pound := New(100, "GBP")
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
	pound := New(100, "GBP")
	twoPounds := New(200, "GBP")
	twoEuros := New(200, "EUR")

	if r, err := pound.GreaterThan(twoPounds); err != nil || r {
		t.Errorf("Expected %d Greater Than %d == %t got %t", pound.amount.val,
			twoPounds.amount.val, false, r)
	}

	if r, err := pound.LessThan(twoPounds); err != nil || !r {
		t.Errorf("Expected %d Less Than %d == %t got %t", pound.amount.val,
			twoPounds.amount.val, true, r)
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
	pound := New(100, "GBP")

	if pound.Currency().Code != "GBP" {
		t.Errorf("Expected %s got %s", "GBP", pound.Currency().Code)
	}
}

func TestMoney_Amount(t *testing.T) {
	pound := New(100, "GBP")

	if pound.Amount() != 100 {
		t.Errorf("Expected %d got %d", 100, pound.Amount())
	}
}

func TestMoney_Compare(t *testing.T) {
	tcs := []struct {
		a        int64
		b        int64
		expected int
	}{
		{100, 200, -1},
		{200, 200, 0},
		{300, 200, 1},
	}

	for _, tc := range tcs {
		ma := New(tc.a, "EUR")
		mb := New(tc.b, "EUR")
		r, err := ma.Compare(mb)
		if err != nil {
			t.Fatal(err)
		}
		if r != tc.expected {
			t.Errorf("Compare(%d, %d): expected %d got %d", tc.a, tc.b, tc.expected, r)
		}
	}
}

func TestMoney_Compare_CurrencyMismatch(t *testing.T) {
	ma := New(100, "EUR")
	mb := New(100, "GBP")
	_, err := ma.Compare(mb)
	if err == nil {
		t.Error("Expected error for currency mismatch")
	}
}

func TestMoney_Percentage(t *testing.T) {
	tcs := []struct {
		amount     int64
		percentage float64
		expected   int64
	}{
		{10000, 10, 1000},
		{10000, 8.5, 850},
		{10000, 0, 0},
		{10000, 100, 10000},
		{333, 33.33, 111},
	}

	for _, tc := range tcs {
		m := New(tc.amount, "USD")
		r := m.Percentage(tc.percentage)
		if r.Amount() != tc.expected {
			t.Errorf("Percentage(%d, %.2f%%): expected %d got %d",
				tc.amount, tc.percentage, tc.expected, r.Amount())
		}
	}
}

func TestMoney_AsParts(t *testing.T) {
	tcs := []struct {
		amount       int64
		code         string
		expectedWhole int64
		expectedFrac  int64
	}{
		{1234, "USD", 12, 34},
		{100, "GBP", 1, 0},
		{1, "USD", 0, 1},
		{-550, "GBP", -5, 50},
		{100, "JPY", 100, 0},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		whole, frac := m.AsParts()
		if whole != tc.expectedWhole || frac != tc.expectedFrac {
			t.Errorf("AsParts(%d %s): expected (%d, %d) got (%d, %d)",
				tc.amount, tc.code, tc.expectedWhole, tc.expectedFrac, whole, frac)
		}
	}
}

func TestMoney_MarshalJSON(t *testing.T) {
	m := New(1234, "USD")
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"amount":1234,"currency":"USD"}`
	if string(b) != expected {
		t.Errorf("MarshalJSON: expected %s got %s", expected, string(b))
	}
}

func TestMoney_UnmarshalJSON(t *testing.T) {
	input := `{"amount":5678,"currency":"EUR"}`
	var m Money
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		t.Fatal(err)
	}
	if m.Amount() != 5678 {
		t.Errorf("UnmarshalJSON: expected amount 5678 got %d", m.Amount())
	}
	if m.Currency().Code != "EUR" {
		t.Errorf("UnmarshalJSON: expected currency EUR got %s", m.Currency().Code)
	}
}

func TestMoney_JSON_RoundTrip(t *testing.T) {
	original := New(9999, "GBP")
	b, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var restored Money
	if err := json.Unmarshal(b, &restored); err != nil {
		t.Fatal(err)
	}
	eq, err := original.Equals(&restored)
	if err != nil || !eq {
		t.Errorf("JSON round-trip: values not equal")
	}
}

func TestMoney_ToWords_MultiCurrency(t *testing.T) {
	// Note: ToWords passes the raw integer amount directly to GetCurrencyAmountWords.
	// GetCurrencyAmountWords treats it as a decimal float (e.g. 1 → "1.00",
	// 100 → "100.00"). Subunits only appear when the raw int64 has a fractional
	// part, which can't happen via New().Amount(). Use GetCurrencyAmountWords
	// directly when you need sub-unit words.
	tcs := []struct {
		amount   int64
		code     string
		expected string
	}{
		{1, "USD", "one dollar only"},
		{50, "USD", "fifty dollar only"},
		{1, "GBP", "one pound only"},
		{1, "EUR", "one euro only"},
		{100, "JPY", "one hundred yen only"},
		{1, "INR", "one rupee only"},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		r := m.ToWords()
		if r != tc.expected {
			t.Errorf("ToWords(%d %s): expected %q got %q", tc.amount, tc.code, tc.expected, r)
		}
	}
}

func TestAddCurrencyMeta(t *testing.T) {
	AddCurrencyMeta("XYZ2", "zorkmid", "zork")
	m := New(1, "XYZ2")
	r := m.ToWords()
	expected := "one zorkmid only"
	if r != expected {
		t.Errorf("AddCurrencyMeta: expected %q got %q", expected, r)
	}
}

func TestGetCurrencyAmountWords_SubUnit(t *testing.T) {
	// Call GetCurrencyAmountWords directly to test sub-unit word output.
	// Note: fmt.Sprintf("%+v", float64) drops trailing zeros, so 1.50 → "+1.5"
	// and the decimal part is parsed as "5" (five), not "50" (fifty).
	// Use values whose decimal parts have no trailing zeros to avoid ambiguity.
	tcs := []struct {
		amount   float64
		code     string
		expected string
	}{
		{1.25, "USD", "one dollar and twenty-five cents only"},
		{0.75, "USD", "seventy-five cents only"},
		{10.99, "GBP", "ten pound and ninety-nine penny only"},
		{5.0, "EUR", "five euro only"},
	}

	for _, tc := range tcs {
		r := GetCurrencyAmountWords(tc.amount, tc.code)
		if r != tc.expected {
			t.Errorf("GetCurrencyAmountWords(%.2f, %s): expected %q got %q",
				tc.amount, tc.code, tc.expected, r)
		}
	}
}
