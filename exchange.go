package money

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// ExchangeRate represents the rate at which one currency converts to another.
// It is used to convert a Money value from its source currency to a target
// currency without losing precision in the source representation.
//
// Rates are stored as float64 and are precise enough for all practical
// foreign-exchange use cases. The converted amount is rounded to the target
// currency's smallest unit using RoundHalfUp.
//
//	rate, _ := money.NewExchangeRate("USD", "EUR", 0.92)
//	euros, _ := rate.Convert(money.New(1000, "USD")) // €9.20
type ExchangeRate struct {
	from string
	to   string
	rate float64
}

// NewExchangeRate creates an ExchangeRate converting from → to at the given
// rate (positive, non-zero).
//
//	NewExchangeRate("USD", "EUR", 0.92)  // 1 USD = 0.92 EUR
//	NewExchangeRate("GBP", "INR", 106.5) // 1 GBP = 106.5 INR
func NewExchangeRate(from, to string, rate float64) (*ExchangeRate, error) {
	if rate <= 0 {
		return nil, errors.New("exchange rate must be positive")
	}
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)
	if from == to {
		return nil, errors.New("from and to currency must differ")
	}
	return &ExchangeRate{from: from, to: to, rate: rate}, nil
}

// From returns the source currency code.
func (r *ExchangeRate) From() string { return r.from }

// To returns the target currency code.
func (r *ExchangeRate) To() string { return r.to }

// Rate returns the numeric exchange rate.
func (r *ExchangeRate) Rate() float64 { return r.rate }

// Invert returns a new ExchangeRate that converts in the opposite direction
// using the reciprocal rate.
//
//	rate, _ := NewExchangeRate("USD", "EUR", 0.92)
//	inv := rate.Invert() // EUR → USD at ~1.0869...
func (r *ExchangeRate) Invert() *ExchangeRate {
	return &ExchangeRate{from: r.to, to: r.from, rate: 1 / r.rate}
}

// Convert converts m from the source currency to the target currency.
// The result is rounded to the target currency's smallest unit using RoundHalfUp.
// Returns an error when m's currency does not match the rate's source currency.
//
//	rate, _ := NewExchangeRate("USD", "EUR", 0.92)
//	euros, _ := rate.Convert(money.New(1000, "USD")) // €9.20 (920 cents)
func (r *ExchangeRate) Convert(m *Money) (*Money, error) {
	if !strings.EqualFold(m.Currency().Code, r.from) {
		return nil, fmt.Errorf(
			"exchange rate converts from %s, but Money is in %s",
			r.from, m.Currency().Code,
		)
	}
	converted := int64(math.Round(float64(m.Amount()) * r.rate))
	return New(converted, r.to), nil
}
