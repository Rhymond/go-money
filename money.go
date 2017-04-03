package money

import (
	"log"
)

type amount struct {
	val int
}

type Money struct {
	amount   *amount
	Currency *Currency
}

var calc *calculator

// New creates and returns new instance of Money
func New(value int, currency string) *Money {

	calc = new(calculator)

	return &Money{
		&amount{value},
		new(Currency).Get(currency),
	}
}

// SameCurrency check if given Money is equals by currency
func (m *Money) SameCurrency(om *Money) bool {
	return m.Currency.Equals(om.Currency)
}

func (m *Money) assertSameCurrency(om *Money) {
	if !m.SameCurrency(om) {
		log.Fatalf("Currencies %s and %s don't match", m.Currency.Code, om.Currency.Code)
	}
}

func (m *Money) compare(om *Money) int {

	m.assertSameCurrency(om)

	switch {
	case m.amount.val > om.amount.val:
		return 1
	case m.amount.val < om.amount.val:
		return -1
	}

	return 0
}

func (m *Money) Equals(om *Money) bool {
	return m.compare(om) == 0
}

func (m *Money) GreaterThan(om *Money) bool {
	return m.compare(om) == 1
}

func (m *Money) GreaterThanOrEqual(om *Money) bool {
	return m.compare(om) >= 0
}

func (m *Money) LessThan(om *Money) bool {
	return m.compare(om) == -1
}

func (m *Money) LessThanOrEqual(om *Money) bool {
	return m.compare(om) <= 0
}

func (m *Money) IsZero() bool {
	return m.amount.val == 0
}

func (m *Money) IsPositive() bool {
	return m.amount.val > 0
}

func (m *Money) IsNegative() bool {
	return m.amount.val < 0
}

func (m *Money) Absolute() *Money {
	return &Money{calc.absolute(m.amount), m.Currency}
}

func (m *Money) Negative() *Money {
	return &Money{calc.negative(m.amount), m.Currency}
}

func (m *Money) Add(om *Money) *Money {
	m.assertSameCurrency(om)
	return &Money{calc.add(m.amount, om.amount), m.Currency}
}

func (m *Money) Subtract(om *Money) *Money {
	m.assertSameCurrency(om)
	return &Money{calc.subtract(m.amount, om.amount), m.Currency}
}

func (m *Money) Multiply(mul int) *Money {
	return &Money{calc.multiply(m.amount, mul), m.Currency}
}

func (m *Money) Divide(div int) *Money {
	return &Money{calc.divide(m.amount, div), m.Currency}
}

func (m *Money) Round() *Money {
	return &Money{calc.round(m.amount), m.Currency}
}

func (m *Money) Split(n int) []*Money {
	if n <= 0 {
		log.Fatalf("Split must be higher than zero")
	}

	a := calc.divide(m.amount, n)
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{a, m.Currency}
	}

	l := calc.modulus(m.amount, n).val

	// Add leftovers to the first parties
	for p := 0; l != 0; p++ {
		ms[p].amount = calc.add(ms[p].amount, &amount{1})
		l -= 1
	}

	return ms
}

func (m *Money) Allocate(rs []int) []*Money {
	if len(rs) == 0 {
		log.Fatalf("No ratios specified")
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
			m.Currency,
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

	return ms
}

func (m *Money) Display() string {
	f := NewFormatter(m.Currency.Fraction, ".", ",",
		m.Currency.Grapheme, m.Currency.Template)

	return f.Format(m.amount.val)
}
