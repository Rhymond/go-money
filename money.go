package money

import (
	"errors"
	"math"
	"math/big"
)

// ErrCurrencyMismatch happens when two compared Money don't have the same currency.
var ErrCurrencyMismatch = errors.New("currencies don't match")

// Money represents monetary value information, stores
// currency and amount value.
type Money struct {
	amount   *Decimal  `db:"amount"`
	currency *Currency `db:"currency"`
}

// New creates and returns new instance of Money.
func New(amount int64, code string) *Money {
	return &Money{
		amount:   &Decimal{val: big.NewInt(amount)},
		currency: newCurrency(code).get(),
	}
}

// safe returns a Money with both fields non-nil. The receiver is returned
// unchanged when already populated; otherwise a normalized copy is returned
// (zero amount and a default Currency for any unset field). The original is
// never mutated. Every method that touches m.amount or m.currency calls
// safe() first so the zero value Money{} doesn't panic.
func (m *Money) safe() *Money {
	if m.amount != nil && m.currency != nil {
		return m
	}
	fixed := *m
	if fixed.amount == nil {
		fixed.amount = &Decimal{val: new(big.Int)}
	}
	if fixed.currency == nil {
		fixed.currency = (&Currency{}).getDefault()
	}
	return &fixed
}

// Currency returns the currency used by Money. The returned *Currency is a
// defensive copy — mutating its fields does not affect the global registry
// or any other Money instance.
func (m *Money) Currency() *Currency {
	c := m.safe().currency
	cp := *c
	return &cp
}

// Amount returns the monetary value in the currency's smallest unit as an
// int64. Sub-smallest-unit precision (introduced by Multiply with a
// fractional Decimal) is truncated toward zero — call Round first if you
// want to snap to the smallest unit.
//
// The internal value is held in a big.Int, so arithmetic (Add, Subtract,
// Multiply, Allocate, Split) preserves precision even when intermediate
// values exceed int64. Amount itself truncates to int64 — values whose
// smallest-unit magnitude exceeds ~9.2e18 cannot be read back faithfully
// through Amount, Display, AsMajorUnits, MarshalJSON, MarshalXML, or Value.
// Callers working anywhere near that ceiling should keep their working
// magnitudes inside int64 or read the Decimal directly via
// NewDecimalFromMoney.
func (m *Money) Amount() int64 {
	if m.amount == nil {
		return 0
	}
	if m.amount.exponent <= 0 {
		return m.amount.val.Int64()
	}
	return new(big.Int).Quo(m.amount.val, pow10(int64(m.amount.exponent))).Int64()
}

// SameCurrency check if given Money is equals by currency.
func (m *Money) SameCurrency(om *Money) bool {
	return m.safe().currency.equals(om.safe().currency)
}

func (m *Money) assertSameCurrency(om *Money) error {
	if !m.SameCurrency(om) {
		return ErrCurrencyMismatch
	}

	return nil
}

func (m *Money) compare(om *Money) int {
	av, bv, _ := align(m.safe().amount, om.safe().amount)
	return av.Cmp(bv)
}

// Equals checks equality between two Money types.
func (m *Money) Equals(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 0, nil
}

// GreaterThan checks whether the value of Money is greater than the other.
func (m *Money) GreaterThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 1, nil
}

// GreaterThanOrEqual checks whether the value of Money is greater or equal than the other.
func (m *Money) GreaterThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) >= 0, nil
}

// LessThan checks whether the value of Money is less than the other.
func (m *Money) LessThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == -1, nil
}

// LessThanOrEqual checks whether the value of Money is less or equal than the other.
func (m *Money) LessThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) <= 0, nil
}

// IsZero returns boolean of whether the value of Money is equals to zero.
func (m *Money) IsZero() bool {
	return m.safe().amount.Sign() == 0
}

// IsPositive returns boolean of whether the value of Money is positive.
func (m *Money) IsPositive() bool {
	return m.safe().amount.Sign() > 0
}

// IsNegative returns boolean of whether the value of Money is negative.
func (m *Money) IsNegative() bool {
	return m.safe().amount.Sign() < 0
}

// Absolute returns new Money struct from given Money using absolute monetary value.
func (m *Money) Absolute() *Money {
	m = m.safe()
	return &Money{amount: m.amount.absolute(), currency: m.currency}
}

// Negative returns new Money struct from given Money using negative monetary value.
func (m *Money) Negative() *Money {
	m = m.safe()
	return &Money{amount: m.amount.negative(), currency: m.currency}
}

// Add returns new Money struct with value representing sum of Self and Other Money.
func (m *Money) Add(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return m, nil
	}

	m = m.safe()
	sum := m.amount

	for _, m2 := range ms {
		if err := m.assertSameCurrency(m2); err != nil {
			return nil, err
		}

		sum = sum.add(m2.safe().amount)
	}

	return &Money{amount: sum, currency: m.currency}, nil
}

// Subtract returns new Money struct with value representing difference of Self and Other Money.
func (m *Money) Subtract(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return m, nil
	}

	m = m.safe()
	diff := m.amount

	for _, m2 := range ms {
		if err := m.assertSameCurrency(m2); err != nil {
			return nil, err
		}

		diff = diff.subtract(m2.safe().amount)
	}

	return &Money{amount: diff, currency: m.currency}, nil
}

// Multiply returns a new Money struct with value representing Self multiplied
// by every supplied Decimal in order. The result preserves full decimal
// precision — e.g. $0.01 * 1.5 yields 15 at exponent 1 ($0.015). Display() and
// Amount() truncate sub-smallest-unit precision toward zero; call Round to
// snap to nearest smallest unit.
//
// With no multipliers it returns Self unchanged. A nil *Decimal is treated as
// "multiply by nothing": the running result collapses to zero (currency
// preserved), and any remaining multipliers are ignored.
func (m *Money) Multiply(muls ...*Decimal) *Money {
	m = m.safe()
	result := m.amount
	for _, d := range muls {
		if d == nil {
			return &Money{
				amount:   &Decimal{val: new(big.Int)},
				currency: m.currency,
			}
		}
		result = result.multiply(d)
	}
	return &Money{amount: result, currency: m.currency}
}

// Round returns a new Money struct with sub-smallest-unit precision removed
// using half-away-from-zero. After Round, the underlying Decimal sits at
// exponent 0, so Amount() and Display() reflect the rounded value exactly.
//
// e.g. New(1, USD).Multiply(NewDecimalFromString("1.5")).Round() yields
// $0.02 (rounded up from $0.015). Money created via New() (already at
// smallest-unit precision) is returned unchanged.
func (m *Money) Round() *Money {
	m = m.safe()
	return &Money{amount: m.amount.roundToSmallestUnit(), currency: m.currency}
}

// RoundToMajor returns a new Money rounded to the nearest major-unit
// boundary using half-away-from-zero. e.g. for USD (Fraction=2), a value of
// $1.25 rounds to $1.00 and $1.75 rounds to $2.00. Useful for summary or
// presentation views that want whole-currency totals; for the common
// "drop sub-cent precision" case use Round instead.
func (m *Money) RoundToMajor() *Money {
	m = m.safe()
	return &Money{amount: m.amount.round(m.currency.Fraction), currency: m.currency}
}

// Split returns slice of Money structs with split Self value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first will likely receive more pennies than ones that are listed later.
func (m *Money) Split(n int) ([]*Money, error) {
	if n <= 0 {
		return nil, errors.New("split must be higher than zero")
	}

	m = m.safe()
	a := m.amount.divide(int64(n))
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{amount: a, currency: m.currency}
	}

	l := new(big.Int).Abs(m.amount.modulus(int64(n)).val)
	stepVal := big.NewInt(1)
	if m.amount.Sign() < 0 {
		stepVal.Neg(stepVal)
	}
	step := &Decimal{val: stepVal, exponent: m.amount.exponent}
	oneBig := big.NewInt(1)

	// Add leftovers to the first parties.
	for p := 0; l.Sign() != 0; p++ {
		ms[p].amount = ms[p].amount.add(step)
		l.Sub(l, oneBig)
	}

	return ms, nil
}

// Allocate returns slice of Money structs with split Self value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
func (m *Money) Allocate(rs ...int) ([]*Money, error) {
	if len(rs) == 0 {
		return nil, errors.New("no ratios specified")
	}

	m = m.safe()

	// Calculate sum of ratios.
	var sum int64
	for _, r := range rs {
		if r < 0 {
			return nil, errors.New("negative ratios not allowed")
		}
		if int64(r) > (math.MaxInt64 - sum) {
			return nil, errors.New("sum of given ratios exceeds max int")
		}
		sum += int64(r)
	}

	total := new(big.Int)
	ms := make([]*Money, 0, len(rs))
	for _, r := range rs {
		party := &Money{
			amount:   m.amount.allocate(int64(r), sum),
			currency: m.currency,
		}

		ms = append(ms, party)
		total.Add(total, party.amount.val)
	}

	// if the sum of all ratios is zero, then we just return zeros and don't do anything
	// with the leftover
	if sum == 0 {
		return ms, nil
	}

	// Calculate leftover value and divide to first parties.
	lo := new(big.Int).Sub(m.amount.val, total)
	sub := big.NewInt(1)
	if lo.Sign() < 0 {
		sub.Neg(sub)
	}
	step := &Decimal{val: new(big.Int).Set(sub), exponent: m.amount.exponent}

	for p := 0; lo.Sign() != 0; p++ {
		ms[p].amount = ms[p].amount.add(step)
		lo.Sub(lo, sub)
	}

	return ms, nil
}

// Display lets represent Money struct as string in given Currency value.
func (m *Money) Display() string {
	m = m.safe()
	c := m.currency.get()
	return c.Formatter().Format(m.Amount())
}

// AsMajorUnits lets represent Money struct as subunits (float64) in given Currency value
func (m *Money) AsMajorUnits() float64 {
	m = m.safe()
	c := m.currency.get()
	return c.Formatter().ToMajorUnits(m.Amount())
}

// Compare function compares two money of the same type
//
//	if m.amount > om.amount returns (1, nil)
//	if m.amount == om.amount returns (0, nil)
//	if m.amount < om.amount returns (-1, nil)
//
// If compare moneys from distinct currency, return (0, ErrCurrencyMismatch).
func (m *Money) Compare(om *Money) (int, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return 0, err
	}

	return m.compare(om), nil
}
