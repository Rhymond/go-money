package money

import (
	"math/big"
)

// zeroDecimal is the value used by Multiply when it encounters a zero-valued
// Multiplier{}. We allocate a fresh *big.Int per call so callers can't mutate
// the shared sentinel.
func zeroDecimal() *Decimal { return &Decimal{val: new(big.Int)} }

// Multiplier is a pre-validated value usable with Money.Multiply. Because
// validation happens at construction time, Money.Multiply itself does not
// return an error.
//
// Construct via NewMultiplier (ints / floats / *Decimal — cannot fail for
// well-formed input) or NewMultiplierFromString (the only constructor that
// can fail, since string parsing is the only validation that can reject
// well-typed input).
type Multiplier struct {
	d *Decimal
}

// NewMultiplier constructs a Multiplier from any native numeric type (any
// int / uint width, float32, float64) or from an existing *Decimal. The
// underlying Decimal is defensively copied so later mutations to the source
// don't bleed in.
//
// It panics on NaN or ±Inf floats and on a nil *Decimal — those represent
// programming bugs at the call site that should be caught loudly. The check
// itself lives in NewDecimal; integers always succeed.
func NewMultiplier[T Integer | Float | *Decimal](v T) Multiplier {
	d, err := NewDecimal(v)
	if err != nil {
		panic("money: NewMultiplier: " + err.Error())
	}
	return Multiplier{d: d}
}

// NewMultiplierFromString parses a decimal string ("1.25", "-0.075") into a
// Multiplier. This is the only constructor that can return an error — strings
// are the only input that can be malformed.
func NewMultiplierFromString(s string) (Multiplier, error) {
	d, err := NewDecimal(s)
	if err != nil {
		return Multiplier{}, err
	}
	return Multiplier{d: d}, nil
}

// Multiply returns a new Money struct whose value is Self multiplied by every
// supplied Multiplier in order. Multipliers are pre-validated at construction
// time, so this operation cannot fail. With no multipliers it returns Self
// unchanged.
//
// A zero-valued Multiplier{} (e.g. one obtained by ignoring an error from
// NewMultiplierFromString) is treated as zero — multiplying by it collapses
// the result to zero Money, matching ordinary arithmetic semantics.
//
// The result preserves full decimal precision — e.g. $0.01 * 1.5 yields 15
// at exponent 3 ($0.015). Call .Round() if you need to collapse to the
// currency's smallest unit.
func (m *Money) Multiply(muls ...Multiplier) *Money {
	result := m.amount
	for _, mul := range muls {
		d := mul.d
		if d == nil {
			d = zeroDecimal()
		}
		result = mutate.calc.multiply(result, d)
	}
	return &Money{amount: result, currency: m.currency}
}
