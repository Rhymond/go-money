package money

import (
	"log"
)

type Amount struct {
	val int
}

type Money struct {
	Amount   *Amount
	Currency *Currency
}

var calc *Calculator

// New creates and returns new instance of Money
func New(am int, curr string) *Money {

	calc = new(Calculator)

	return &Money{
		&Amount{am},
		new(Currency).Get(curr),
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
	case m.Amount.val > om.Amount.val:
		return 1
	case m.Amount.val < om.Amount.val:
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
	return m.Amount.val == 0
}

func (m *Money) IsPositive() bool {
	return m.Amount.val > 0
}

func (m *Money) IsNegative() bool {
	return m.Amount.val < 0
}

func (m *Money) Absolute() *Money {
	return &Money{calc.absolute(m.Amount), m.Currency}
}

func (m *Money) Negative() *Money {
	return &Money{calc.negative(m.Amount), m.Currency}
}

func (m *Money) Add(om *Money) *Money {
	m.assertSameCurrency(om)
	return &Money{calc.add(m.Amount, om.Amount), m.Currency}
}

func (m *Money) Subtract(om *Money) *Money {
	m.assertSameCurrency(om)
	return &Money{calc.subtract(m.Amount, om.Amount), m.Currency}
}

func (m *Money) Multiply(mul int) *Money {
	return &Money{calc.multiply(m.Amount, mul), m.Currency}
}

func (m *Money) Divide(div int) *Money {
	return &Money{calc.divide(m.Amount, div), m.Currency}
}

func (m *Money) Round() *Money {
	return &Money{calc.round(m.Amount), m.Currency}
}

func (m *Money) Split(n int) []*Money {
	if n <= 0 {
		log.Fatalf("Split must be higher than zero")
	}

	a := calc.divide(m.Amount, n)
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{a, m.Currency}
	}

	l := calc.modulus(m.Amount, n).val

	// Add leftovers to the first parties
	for p := 0; l != 0; p++ {
		ms[p].Amount = calc.add(ms[p].Amount, &Amount{1})
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
			calc.allocate(m.Amount, r, sum),
			m.Currency,
		}

		ms = append(ms, party)
		total += party.Amount.val
	}

	// Calculate leftover value and divide to first parties
	lo := m.Amount.val - total
	sub := 1
	if lo < 0 {
		sub = -1
	}

	for p := 0; lo != 0; p++ {
		ms[p].Amount = calc.add(ms[p].Amount, &Amount{sub})
		lo -= sub
	}

	return ms
}

func unit() {

}