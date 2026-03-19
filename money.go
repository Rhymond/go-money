// Package money provides types and methods for precise monetary value
// arithmetic. It follows Martin Fowler's Money pattern by representing
// amounts as integers in the currency's smallest unit (e.g. cents,
// pence, paise), thereby avoiding floating-point rounding errors.
//
// # Quick start
//
//	pound := money.New(100, "GBP")   // £1.00
//	twoPounds, _ := pound.Add(pound) // £2.00
//
//	parties, _ := twoPounds.Split(3)
//	parties[0].Display() // £0.67
//	parties[1].Display() // £0.67
//	parties[2].Display() // £0.66
//
// All arithmetic operations return a new Money value and leave the
// receiver unchanged (immutable style).
package money

import (
	"encoding/json"
	"errors"
	"math"
)

// Amount is a data structure that stores the amount being used for calculations.
type Amount struct {
	val int64
}

// Money represents a monetary value: an amount denominated in a specific currency.
// The amount is stored as an integer in the currency's smallest unit
// (e.g. 100 represents £1.00, $1.00, €1.00 etc.).
type Money struct {
	amount   *Amount
	currency *Currency
}

// New creates and returns a new Money instance.
// amount must be expressed in the currency's smallest unit
// (e.g. 100 for £1.00).
// code is the ISO 4217 currency code (e.g. "GBP", "USD", "EUR").
func New(amount int64, code string) *Money {
	return &Money{
		amount:   &Amount{val: amount},
		currency: newCurrency(code).get(),
	}
}

// NewFromFloat creates a Money value from a float64 amount expressed in the
// currency's main unit (e.g. 1.25 for $1.25 USD).
// The float is multiplied by 10^Fraction and rounded to the nearest integer.
//
//	money.NewFromFloat(1.25, "USD")  // 125 cents  → $1.25
//	money.NewFromFloat(9.999, "USD") // 1000 cents → $10.00 (rounded)
//	money.NewFromFloat(100, "JPY")   // 100 yen    → ¥100
func NewFromFloat(amount float64, code string) *Money {
	c := newCurrency(code).get()
	exp := math.Pow10(c.Fraction)
	return &Money{
		amount:   &Amount{val: int64(math.Round(amount * exp))},
		currency: c,
	}
}

// Currency returns the currency associated with this Money.
func (m *Money) Currency() *Currency {
	return m.currency
}

// Amount returns the monetary value in the currency's smallest unit (e.g. cents).
func (m *Money) Amount() int64 {
	return m.amount.val
}

// SameCurrency reports whether m and om use the same currency.
func (m *Money) SameCurrency(om *Money) bool {
	return m.currency.equals(om.currency)
}

func (m *Money) assertSameCurrency(om *Money) error {
	if !m.SameCurrency(om) {
		return errors.New("currencies don't match")
	}

	return nil
}

func (m *Money) compare(om *Money) int {
	switch {
	case m.amount.val > om.amount.val:
		return 1
	case m.amount.val < om.amount.val:
		return -1
	}

	return 0
}

// Compare compares m to om. Both must share the same currency.
// It returns -1 if m < om, 0 if m == om, or +1 if m > om.
// Returns an error when the currencies differ.
func (m *Money) Compare(om *Money) (int, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return 0, err
	}

	return m.compare(om), nil
}

// Equals reports whether m and om represent the same monetary value.
// Returns an error when the currencies differ.
func (m *Money) Equals(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 0, nil
}

// GreaterThan reports whether m is greater than om.
// Returns an error when the currencies differ.
func (m *Money) GreaterThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 1, nil
}

// GreaterThanOrEqual reports whether m is greater than or equal to om.
// Returns an error when the currencies differ.
func (m *Money) GreaterThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) >= 0, nil
}

// LessThan reports whether m is less than om.
// Returns an error when the currencies differ.
func (m *Money) LessThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == -1, nil
}

// LessThanOrEqual reports whether m is less than or equal to om.
// Returns an error when the currencies differ.
func (m *Money) LessThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) <= 0, nil
}

// IsZero reports whether the monetary value is zero.
func (m *Money) IsZero() bool {
	return m.amount.val == 0
}

// IsPositive reports whether the monetary value is greater than zero.
func (m *Money) IsPositive() bool {
	return m.amount.val > 0
}

// IsNegative reports whether the monetary value is less than zero.
func (m *Money) IsNegative() bool {
	return m.amount.val < 0
}

// Absolute returns a new Money with the absolute (non-negative) monetary value.
func (m *Money) Absolute() *Money {
	return &Money{amount: mutate.calc.absolute(m.amount), currency: m.currency}
}

// Negative returns a new Money with the negated monetary value.
// Positive amounts become negative; already-negative amounts are unchanged.
func (m *Money) Negative() *Money {
	return &Money{amount: mutate.calc.negative(m.amount), currency: m.currency}
}

// Add returns a new Money whose value is m + om.
// Returns an error when the currencies differ.
func (m *Money) Add(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	return &Money{amount: mutate.calc.add(m.amount, om.amount), currency: m.currency}, nil
}

// Subtract returns a new Money whose value is m − om.
// Returns an error when the currencies differ.
func (m *Money) Subtract(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	return &Money{amount: mutate.calc.subtract(m.amount, om.amount), currency: m.currency}, nil
}

// Multiply returns a new Money whose value is m × mul.
func (m *Money) Multiply(mul int64) *Money {
	return &Money{amount: mutate.calc.multiply(m.amount, mul), currency: m.currency}
}

// Divide returns a new Money whose value is m ÷ div.
// Integer division is used; fractional remainders are truncated.
// Use Split or Allocate when lossless distribution is required.
func (m *Money) Divide(div int64) *Money {
	return &Money{amount: mutate.calc.divide(m.amount, div), currency: m.currency}
}

// Percentage returns a new Money representing the given percentage of m.
// The result is truncated to the smallest currency unit.
//
//	money.New(10000, "USD").Percentage(8.5) // $8.50 (8.5% of $100.00)
func (m *Money) Percentage(p float64) *Money {
	return &Money{
		amount:   &Amount{val: int64(math.Round(float64(m.amount.val) * p / 100))},
		currency: m.currency,
	}
}

// Round returns a new Money whose value is rounded to the nearest whole
// currency unit (e.g. to the nearest dollar, pound, or euro).
// The currency's Fraction field determines the rounding magnitude.
func (m *Money) Round() *Money {
	return &Money{amount: mutate.calc.round(m.amount, m.currency.Fraction), currency: m.currency}
}

// AsParts returns the whole-unit and fractional-unit components of the
// monetary value separately.
//
//	money.New(1234, "USD").AsParts() // (12, 34) → $12.34
//	money.New(-550, "GBP").AsParts() // (-5, 50) → -£5.50
func (m *Money) AsParts() (whole int64, frac int64) {
	if m.currency.Fraction == 0 {
		return m.amount.val, 0
	}
	exp := int64(math.Pow10(m.currency.Fraction))
	whole = m.amount.val / exp
	frac = m.amount.val % exp
	if frac < 0 {
		frac = -frac
	}
	return whole, frac
}

// Split distributes m evenly among n parties.
// Any leftover pennies (due to integer division) are given to the first
// parties in a round-robin fashion so that no value is lost.
//
//	money.New(100, "GBP").Split(3)
//	// → [£0.34, £0.33, £0.33]
func (m *Money) Split(n int) ([]*Money, error) {
	if n <= 0 {
		return nil, errors.New("split must be higher than zero")
	}

	a := mutate.calc.divide(m.amount, int64(n))
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{amount: a, currency: m.currency}
	}

	l := mutate.calc.modulus(m.amount, int64(n)).val

	// Add leftovers to the first parties.
	for p := 0; l != 0; p++ {
		ms[p].amount = mutate.calc.add(ms[p].amount, &Amount{1})
		l--
	}

	return ms, nil
}

// Allocate distributes m according to the given ratios without losing
// any pennies. Any leftover value (from integer arithmetic) is spread
// across the first parties in round-robin order.
//
//	money.New(100, "GBP").Allocate(33, 33, 33)
//	// → [£0.34, £0.33, £0.33]
func (m *Money) Allocate(rs ...int) ([]*Money, error) {
	if len(rs) == 0 {
		return nil, errors.New("no ratios specified")
	}

	// Calculate sum of ratios.
	var sum int
	for _, r := range rs {
		sum += r
	}

	var total int64
	var ms []*Money
	for _, r := range rs {
		party := &Money{
			amount:   mutate.calc.allocate(m.amount, r, sum),
			currency: m.currency,
		}

		ms = append(ms, party)
		total += party.amount.val
	}

	// Calculate leftover value and distribute to first parties.
	lo := m.amount.val - total
	sub := int64(1)
	if lo < 0 {
		sub = -sub
	}

	for p := 0; lo != 0; p++ {
		ms[p].amount = mutate.calc.add(ms[p].amount, &Amount{sub})
		lo -= sub
	}

	return ms, nil
}

// Display returns the monetary value formatted as a human-readable string
// using the currency's symbol, decimal separator, and thousands separator.
//
//	money.New(123456789, "EUR").Display() // "€1,234,567.89"
//	money.New(100, "GBP").Display()       // "£1.00"
func (m *Money) Display() string {
	c := m.currency.get()
	return c.Formatter().Format(m.amount.val)
}

// String implements the fmt.Stringer interface and returns the same value as
// Display, making Money usable directly in fmt.Print / fmt.Sprintf.
func (m *Money) String() string {
	return m.Display()
}

// AsFloat64 returns the monetary value as a float64 in the currency's main
// unit (e.g. 125 cents → 1.25 for USD).
//
// WARNING: float64 cannot represent all decimal fractions exactly. This method
// is provided for display/logging purposes only. Never use the result in further
// monetary calculations — use the integer Amount() instead.
//
//	money.New(1234, "USD").AsFloat64() // 12.34
//	money.New(100,  "JPY").AsFloat64() // 100.0
func (m *Money) AsFloat64() float64 {
	if m.currency.Fraction == 0 {
		return float64(m.amount.val)
	}
	exp := math.Pow10(m.currency.Fraction)
	return float64(m.amount.val) / exp
}

// ToWords returns the monetary value expressed as English words.
// The currency must be registered in CountryCurrencyMeta; otherwise
// the raw numeric string is returned.
//
//	money.New(100, "PHP").ToWords() // "one hundred pesos only"
//	money.New(150, "USD").ToWords() // "one dollar and fifty cents only"
func (m *Money) ToWords() string {
	c := m.currency.get()
	return GetCurrencyAmountWords(float64(m.Amount()), c.Code)
}

// moneyJSON is the canonical JSON representation of a Money value.
type moneyJSON struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// MarshalJSON implements the json.Marshaler interface.
// The JSON representation is {"amount":<int64>,"currency":"<CODE>"}.
//
//	money.New(100, "USD") → {"amount":100,"currency":"USD"}
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(&moneyJSON{
		Amount:   m.amount.val,
		Currency: m.currency.Code,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It expects {"amount":<int64>,"currency":"<CODE>"}.
func (m *Money) UnmarshalJSON(b []byte) error {
	var mj moneyJSON
	if err := json.Unmarshal(b, &mj); err != nil {
		return err
	}
	m.amount = &Amount{val: mj.Amount}
	m.currency = newCurrency(mj.Currency).get()
	return nil
}
