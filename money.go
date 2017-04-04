package money

import (
	"errors"
)

type amount struct {
	val int
}

// Money stores money information
type Money struct {
	amount   *amount
	currency *currency
}

var calc *calculator

// New creates and returns new instance of Money
func New(Amount int, code string) *Money {

	calc = new(calculator)

	return &Money{
		&amount{Amount},
		newCurrency(code),
	}
}

// SameCurrency check if given Money is equals by currency
func (m *Money) SameCurrency(om *Money) bool {
	return m.currency.equals(om.currency)
}

func (m *Money) assertSameCurrency(om *Money) error {
	if !m.SameCurrency(om) {
		return errors.New("Currencies don't match")
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

// Equals checks equality between two Money types
func (m *Money) Equals(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 0, nil
}

// GreaterThan checks whether the value of Money is greater than the other
func (m *Money) GreaterThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 1, nil
}

// GreaterThanOrEqual checks whether the value of Money is greater or equal than the other
func (m *Money) GreaterThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) >= 0, nil
}

// LessThan checks whether the value of Money is less than the other
func (m *Money) LessThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == -1, nil
}

// LessThanOrEqual checks whether the value of Money is less or equal than the other
func (m *Money) LessThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) <= 0, nil
}

// IsZero returns boolean of whether the value of Money is equals to zero
func (m *Money) IsZero() bool {
	return m.amount.val == 0
}

// IsPositive returns boolean of whether the value of Money is positive
func (m *Money) IsPositive() bool {
	return m.amount.val > 0
}

// IsNegative returns boolean of whether the value of Money is negative
func (m *Money) IsNegative() bool {
	return m.amount.val < 0
}

// Absolute returns new Money struct from given Money using absolute monetary value
func (m *Money) Absolute() *Money {
	return &Money{calc.absolute(m.amount), m.currency}
}

// Negative returns new Money struct from given Money using negative monetary value
func (m *Money) Negative() *Money {
	return &Money{calc.negative(m.amount), m.currency}
}

// Add returns new Money struct with value representing sum of Self and Other Money
func (m *Money) Add(om *Money) *Money {
	m.assertSameCurrency(om)
	return &Money{calc.add(m.amount, om.amount), m.currency}
}

// Subtract returns new Money struct with value representing difference of Self and Other Money
func (m *Money) Subtract(om *Money) *Money {
	m.assertSameCurrency(om)
	return &Money{calc.subtract(m.amount, om.amount), m.currency}
}

// Multiply returns new Money struct with value representing Self multiplied value by multiplier
func (m *Money) Multiply(mul int) *Money {
	return &Money{calc.multiply(m.amount, mul), m.currency}
}

// Divide returns new Money struct with value representing Self division value by given divider
func (m *Money) Divide(div int) *Money {
	return &Money{calc.divide(m.amount, div), m.currency}
}

// Round returns new Money struct with value rounded to nearest zero
func (m *Money) Round() *Money {
	return &Money{calc.round(m.amount), m.currency}
}

// Split returns slice of Money structs with split Self value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first will likely receive more pennies than ones that are listed later
func (m *Money) Split(n int) ([]*Money, error) {
	if n <= 0 {
		return nil, errors.New("Split must be higher than zero")
	}

	a := calc.divide(m.amount, n)
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{a, m.currency}
	}

	l := calc.modulus(m.amount, n).val

	// Add leftovers to the first parties
	for p := 0; l != 0; p++ {
		ms[p].amount = calc.add(ms[p].amount, &amount{1})
		l--
	}

	return ms, nil
}

// Allocate returns slice of Money structs with split Self value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
func (m *Money) Allocate(rs []int) ([]*Money, error) {
	if len(rs) == 0 {
		return nil, errors.New("No ratios specified")
	}

	// Calculate sum of ratios
	var sum int
	for _, r := range rs {
		sum += r
	}

	var total int
	var ms []*Money
	for _, r := range rs {
		party := &Money{
			calc.allocate(m.amount, r, sum),
			m.currency,
		}

		ms = append(ms, party)
		total += party.amount.val
	}

	// Calculate leftover value and divide to first parties
	lo := m.amount.val - total
	sub := 1
	if lo < 0 {
		sub = -1
	}

	for p := 0; lo != 0; p++ {
		ms[p].amount = calc.add(ms[p].amount, &amount{sub})
		lo -= sub
	}

	return ms, nil
}

// Display lets represent Money struct as string in given Currency value
func (m *Money) Display() (string, error) {
	c, err := m.currency.get()

	if err != nil {
		return "", err
	}

	f := NewFormatter(c.Fraction, ".", ",", c.Grapheme, c.Template)
	return f.Format(m.amount.val), nil
}
