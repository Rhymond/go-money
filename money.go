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

func (m *Money) IsSameCurrency(om *Money) bool {
	return m.Currency.Equals(om.Currency)
}

func (m *Money) Compare(om *Money) int {

	if !m.IsSameCurrency(om) {
		log.Fatalf("Cannot compare different currencies")
	}

	switch {
	case m.Amount == om.Amount: return 0
	case m.Amount > om.Amount: return 1
	default: return -1
	}
}

func (m *Money) Equals(om *Money) bool {
	return m.IsSameCurrency(om) && m.Amount == om.Amount
}

func (m *Money)