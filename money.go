package gocash

import "log"

type Money struct {
	Amount   int
	Currency *Currency
}

func New(a int, c string) *Money {
	return &Money{
		Amount:   a,
		Currency: new(Currency).Get(c),
	}
}

func (m *Money) SameCurrency(om *Money) bool {
	return m.Currency.Equals(om.Currency)
}

func (m *Money) assertSameCurrency(om *Money) {
	if !m.SameCurrency(om) {
		log.Fatalf("Currencies don't match")
	}
}

func (m *Money) compare(om *Money) int {

	m.assertSameCurrency(om)

	switch {
	case m.Amount > om.Amount:
		return 1
	case m.Amount < om.Amount:
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
	return m.compare(om) >= 1
}

func (m *Money) LessThan(om *Money) bool {
	return m.compare(om) == -1
}

func (m *Money) LessThanOrEqual(om *Money) bool {
	return m.compare(om) <= 0
}

//func (m *Money) IsZero() bool     {}
//func (m *Money) IsPositive() bool {}
//func (m *Money) IsNegative() bool {}
//
//func (m *Money) Negative() {}
//func (m *Money) Absolute() {}
//func (m *Money) Round()    {}
//
//func (m *Money) Add(om *Money) *Money      {}
//func (m *Money) Subtract(om *Money) *Money {}
//func (m *Money) Multiply(om *Money) *Money {}
//func (m *Money) Divide(om *Money) *Money   {}
//func (m *Money) Allocate(r []int) []*Money {}
//func (m *Money) Split(n int) []*Money      {}
