package gocash

import (
	"log"
)

type Money struct {
	Number   *Number
	Currency *Currency
}

func New(a int, c string) *Money {
	return &Money{
		new(Number).New(a),
		new(Currency).Get(c),
	}
}

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
	case m.Number.Amount > om.Number.Amount:
		return 1
	case m.Number.Amount < om.Number.Amount:
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
	return m.Number.Amount == 0
}

func (m *Money) IsPositive() bool {
	return m.Number.Amount > 0
}

func (m *Money) IsNegative() bool {
	return m.Number.Amount < 0
}

func (m *Money) Absolute() *Money {
	m.Number.Absolute()
	return m
}

func (m *Money) Negative() *Money {
	m.Number.Negative()
	return m
}

func (m *Money) Add(om *Money) *Money {
	m.assertSameCurrency(om)
	return New(m.Number.Amount+om.Number.Amount, m.Currency.Code)
}

func (m *Money) Subtract(om *Money) *Money {
	m.assertSameCurrency(om)
	return New(m.Number.Amount-om.Number.Amount, m.Currency.Code)
}

func (m *Money) Multiply(mul int) *Money {
	return New(m.Number.Amount*mul, m.Currency.Code)
}

func (m *Money) Divide(div int) *Money {
	return New(m.Number.Amount/div, m.Currency.Code)
}

//func (m *Money) Multiply(om *Money) *Money {}
//func (m *Money) Divide(om *Money) *Money   {}
//func (m *Money) Allocate(r []int) []*Money {}
//func (m *Money) Split(n int) []*Money      {}
//func (m *Money) Round() {}
