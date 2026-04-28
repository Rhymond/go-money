package money

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
)

// Injection points for backward compatibility.
// If you need to keep your JSON marshal/unmarshal way, overwrite them like below.
//
//	money.UnmarshalJSON = func (m *Money, b []byte) error { ... }
//	money.MarshalJSON = func (m Money) ([]byte, error) { ... }
var (
	// UnmarshalJSON is injection point of json.Unmarshaller for money.Money
	UnmarshalJSON = defaultUnmarshalJSON
	// MarshalJSON is injection point of json.Marshaller for money.Money
	MarshalJSON = defaultMarshalJSON

	// ErrCurrencyMismatch happens when two compared Money don't have the same currency.
	ErrCurrencyMismatch = errors.New("currencies don't match")

	// ErrInvalidJSONUnmarshal happens when the default money.UnmarshalJSON fails to unmarshal Money because of invalid data.
	ErrInvalidJSONUnmarshal = errors.New("invalid json unmarshal")
)

func defaultUnmarshalJSON(m *Money, b []byte) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	var amount float64
	if amountRaw, ok := data["amount"]; ok {
		amount, ok = amountRaw.(float64)
		if !ok {
			return ErrInvalidJSONUnmarshal
		}
	}

	var currency string
	if currencyRaw, ok := data["currency"]; ok {
		currency, ok = currencyRaw.(string)
		if !ok {
			return ErrInvalidJSONUnmarshal
		}
	}

	var ref *Money
	if amount == 0 && currency == "" {
		ref = &Money{}
	} else {
		ref = New(int64(amount), currency)
	}

	*m = *ref
	return nil
}

func defaultMarshalJSON(m Money) ([]byte, error) {
	if m == (Money{}) {
		m = *New(0, "")
	}

	buff := bytes.NewBufferString(fmt.Sprintf(`{"amount": %d, "currency": "%s"}`, m.Amount(), m.Currency().Code))
	return buff.Bytes(), nil
}

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

// NewFromFloat creates and returns new instance of Money from a float64.
// Always rounding trailing decimals down.
func NewFromFloat(amount float64, code string) *Money {
	currencyDecimals := math.Pow10(newCurrency(code).get().Fraction)
	return New(int64(amount*currencyDecimals), code)
}

// Currency returns the currency used by Money.
func (m *Money) Currency() *Currency {
	return m.currency
}

// Amount returns a copy of the internal monetary value as an int64.
func (m *Money) Amount() int64 {
	return m.amount.Int64()
}

// SameCurrency check if given Money is equals by currency.
func (m *Money) SameCurrency(om *Money) bool {
	return m.currency.equals(om.currency)
}

func (m *Money) assertSameCurrency(om *Money) error {
	if !m.SameCurrency(om) {
		return ErrCurrencyMismatch
	}

	return nil
}

func (m *Money) compare(om *Money) int {
	av, bv, _ := align(m.amount, om.amount)
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
	return m.amount.Sign() == 0
}

// IsPositive returns boolean of whether the value of Money is positive.
func (m *Money) IsPositive() bool {
	return m.amount.Sign() > 0
}

// IsNegative returns boolean of whether the value of Money is negative.
func (m *Money) IsNegative() bool {
	return m.amount.Sign() < 0
}

// Absolute returns new Money struct from given Money using absolute monetary value.
func (m *Money) Absolute() *Money {
	return &Money{amount: mutate.calc.absolute(m.amount), currency: m.currency}
}

// Negative returns new Money struct from given Money using negative monetary value.
func (m *Money) Negative() *Money {
	return &Money{amount: mutate.calc.negative(m.amount), currency: m.currency}
}

// Add returns new Money struct with value representing sum of Self and Other Money.
func (m *Money) Add(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return m, nil
	}

	k := New(0, m.currency.Code)

	for _, m2 := range ms {
		if err := m.assertSameCurrency(m2); err != nil {
			return nil, err
		}

		k.amount = mutate.calc.add(k.amount, m2.amount)
	}

	return &Money{amount: mutate.calc.add(m.amount, k.amount), currency: m.currency}, nil
}

// Subtract returns new Money struct with value representing difference of Self and Other Money.
func (m *Money) Subtract(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return m, nil
	}

	k := New(0, m.currency.Code)

	for _, m2 := range ms {
		if err := m.assertSameCurrency(m2); err != nil {
			return nil, err
		}

		k.amount = mutate.calc.add(k.amount, m2.amount)
	}

	return &Money{amount: mutate.calc.subtract(m.amount, k.amount), currency: m.currency}, nil
}

// Multiply returns a new Money struct whose value is Self multiplied by every
// supplied multiplier in order. Multipliers may be any type accepted by
// NewDecimal (int, float, string, *Decimal, ...). The result preserves full
// decimal precision: e.g. $0.01 * 1.5 yields 15 at exponent 3 ($0.015).
// Call .Round() if you need to collapse to the currency's smallest unit.
func (m *Money) Multiply(muls ...any) (*Money, error) {
	if len(muls) == 0 {
		return nil, errors.New("at least one multiplier is required")
	}

	result := m.amount
	for _, mul := range muls {
		d, err := NewDecimal(mul)
		if err != nil {
			return nil, err
		}
		result = mutate.calc.multiply(result, d)
	}

	return &Money{amount: result, currency: m.currency}, nil
}

// MustMultiply is like Multiply but panics if any multiplier is invalid.
// Use it for static multipliers that are known to be valid at the call site.
func (m *Money) MustMultiply(muls ...any) *Money {
	r, err := m.Multiply(muls...)
	if err != nil {
		panic(err)
	}
	return r
}

// Round returns new Money struct with value rounded to nearest zero.
func (m *Money) Round() *Money {
	return &Money{amount: mutate.calc.round(m.amount, m.currency.Fraction), currency: m.currency}
}

// Split returns slice of Money structs with split Self value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first will likely receive more pennies than ones that are listed later.
func (m *Money) Split(n int) ([]*Money, error) {
	if n <= 0 {
		return nil, errors.New("split must be higher than zero")
	}

	a := mutate.calc.divide(m.amount, int64(n))
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{amount: a, currency: m.currency}
	}

	l := new(big.Int).Abs(mutate.calc.modulus(m.amount, int64(n)).val)
	step := big.NewInt(1)
	if m.amount.Sign() < 0 {
		step.Neg(step)
	}
	one := &Decimal{val: new(big.Int).Set(step), exponent: m.amount.exponent}

	// Add leftovers to the first parties.
	for p := 0; l.Sign() != 0; p++ {
		ms[p].amount = mutate.calc.add(ms[p].amount, one)
		l.Sub(l, big.NewInt(1))
	}

	return ms, nil
}

// Allocate returns slice of Money structs with split Self value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
// Ratios may be any type accepted by NewDecimal (int, float, string, *Decimal, ...).
func (m *Money) Allocate(rs ...any) ([]*Money, error) {
	if len(rs) == 0 {
		return nil, errors.New("no ratios specified")
	}

	ratios := make([]*Decimal, len(rs))
	sum := &Decimal{val: new(big.Int)}
	for i, r := range rs {
		d, err := NewDecimal(r)
		if err != nil {
			return nil, err
		}
		if d.val.Sign() < 0 {
			return nil, errors.New("negative ratios not allowed")
		}
		ratios[i] = d
		sum = mutate.calc.add(sum, d)
	}

	total := new(big.Int)
	ms := make([]*Money, 0, len(rs))
	for _, r := range ratios {
		party := &Money{
			amount:   mutate.calc.allocate(m.amount, r, sum),
			currency: m.currency,
		}

		ms = append(ms, party)
		total.Add(total, party.amount.val)
	}

	// if the sum of all ratios is zero, then we just return zeros and don't do anything
	// with the leftover
	if sum.val.Sign() == 0 {
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
		ms[p].amount = mutate.calc.add(ms[p].amount, step)
		lo.Sub(lo, sub)
	}

	return ms, nil
}

// Display lets represent Money struct as string in given Currency value.
func (m *Money) Display() string {
	c := m.currency.get()
	return c.Formatter().Format(m.amount.Int64())
}

// AsMajorUnits lets represent Money struct as subunits (float64) in given Currency value
func (m *Money) AsMajorUnits() float64 {
	c := m.currency.get()
	return c.Formatter().ToMajorUnits(m.amount.Int64())
}

// UnmarshalJSON is implementation of json.Unmarshaller
func (m *Money) UnmarshalJSON(b []byte) error {
	return UnmarshalJSON(m, b)
}

// MarshalJSON is implementation of json.Marshaller
func (m Money) MarshalJSON() ([]byte, error) {
	return MarshalJSON(m)
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
