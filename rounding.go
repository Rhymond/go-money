package money

import "math"

// RoundingMode controls how intermediate values are rounded when the result
// does not fit exactly in the currency's smallest unit.
type RoundingMode uint8

const (
	// RoundHalfUp rounds to the nearest unit; ties go away from zero.
	// This is the most common financial rounding mode.
	//   1.5 → 2,  -1.5 → -2
	RoundHalfUp RoundingMode = iota

	// RoundHalfDown rounds to the nearest unit; ties go toward zero.
	//   1.5 → 1,  -1.5 → 1
	RoundHalfDown

	// RoundHalfEven rounds to the nearest unit; ties go to the nearest
	// even number (banker's rounding). This minimises cumulative rounding
	// errors in large data sets.
	//   0.5 → 0,  1.5 → 2,  2.5 → 2,  3.5 → 4
	RoundHalfEven

	// RoundUp always rounds away from zero, regardless of the remainder.
	//   1.1 → 2,  -1.1 → -2
	RoundUp

	// RoundDown truncates toward zero, discarding any remainder.
	//   1.9 → 1,  -1.9 → -1
	RoundDown
)

// RoundWithMode returns a new Money whose amount is rounded to the nearest
// whole currency unit using the specified rounding strategy.
//
//	money.New(150, "USD").RoundWithMode(money.RoundHalfUp)  // $2.00
//	money.New(150, "USD").RoundWithMode(money.RoundHalfDown) // $1.00
//	money.New(250, "USD").RoundWithMode(money.RoundHalfEven) // $2.00 (ties to even)
func (m *Money) RoundWithMode(mode RoundingMode) *Money {
	return &Money{
		amount:   roundWithMode(m.amount, m.currency.Fraction, mode),
		currency: m.currency,
	}
}

// roundWithMode performs the rounding computation for a given Amount, fraction
// (number of decimal places), and RoundingMode.
func roundWithMode(a *Amount, fraction int, mode RoundingMode) *Amount {
	if a.val == 0 || fraction == 0 {
		return &Amount{a.val}
	}

	exp := int64(math.Pow10(fraction))

	sign := int64(1)
	absVal := a.val
	if a.val < 0 {
		sign = -1
		absVal = -a.val
	}

	quotient := absVal / exp
	remainder := absVal % exp

	var rounded int64

	switch mode {
	case RoundHalfUp:
		if remainder*2 >= exp {
			rounded = (quotient + 1) * exp
		} else {
			rounded = quotient * exp
		}

	case RoundHalfDown:
		if remainder*2 > exp {
			rounded = (quotient + 1) * exp
		} else {
			rounded = quotient * exp
		}

	case RoundHalfEven:
		doubled := remainder * 2
		if doubled > exp {
			rounded = (quotient + 1) * exp
		} else if doubled == exp {
			// Tie: round to nearest even quotient
			if quotient%2 == 0 {
				rounded = quotient * exp
			} else {
				rounded = (quotient + 1) * exp
			}
		} else {
			rounded = quotient * exp
		}

	case RoundUp:
		if remainder > 0 {
			rounded = (quotient + 1) * exp
		} else {
			rounded = quotient * exp
		}

	case RoundDown:
		rounded = quotient * exp

	default:
		// Default to RoundHalfUp
		if remainder*2 >= exp {
			rounded = (quotient + 1) * exp
		} else {
			rounded = quotient * exp
		}
	}

	return &Amount{sign * rounded}
}
