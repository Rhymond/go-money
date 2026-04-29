package money

import (
	"errors"
	"math/big"
	"reflect"
)

// Integer matches every Go integer type — signed and unsigned, every width.
// The tilde lets named types (e.g. `type Cents int64`) satisfy the constraint.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Float matches both Go floating-point widths.
type Float interface {
	~float32 | ~float64
}

// Fixed-point fraction sizes for use with NewFromFixedPoint. Values are the
// number of decimal places represented by the fractional component.
const (
	FixedPointSizeDeci  uint = 1
	FixedPointSizeCenti uint = 2
	FixedPointSizeMilli uint = 3
	FixedPointSizeMicro uint = 6
	FixedPointSizeNano  uint = 9
)

// NewFromFixedPoint builds Money from a Google-style fixed-point amount
// (units in major currency units + a fractional component scaled to size
// decimal places). Total value is units + fraction * 10^-size.
//
// Per Google's Money convention, units and fraction should share the same
// sign (or one may be zero). Sub-currency-unit precision is truncated toward
// zero — e.g. a single nano-euro becomes €0.00.
//
// The smallest-unit value is held internally with arbitrary precision, but
// Money.Amount() truncates to int64 — see Amount's docs. Inputs whose
// smallest-unit value exceeds int64 (~±9.2e18) are stored exactly but
// cannot be read back through Amount, Display, etc. without truncation.
//
//	// €1.25 expressed as nanos
//	m := NewFromFixedPoint(1, 250_000_000, FixedPointSizeNano, EUR)
func NewFromFixedPoint[T Integer](units, fraction T, size uint, code string) *Money {
	c := newCurrency(code).get()

	raw := intToBigInt(units)
	raw.Mul(raw, pow10(int64(size)))
	raw.Add(raw, intToBigInt(fraction))

	return &Money{
		amount:   normalizeToSmallestUnit(&Decimal{val: raw, exponent: int(size)}, c.Fraction),
		currency: c,
	}
}

// NewFromMajorUnits builds Money from a whole-unit amount in the currency's
// major unit. NewFromMajorUnits(5, USD) yields $5.00; equivalent to
// New(5*100, USD) for currencies with Fraction == 2. Accepts any integer
// type (int, int32, uint64, ...).
//
// The smallest-unit value is amount * 10^Fraction held in big.Int internally,
// but Money.Amount() truncates to int64. Keep |amount| within ±9.2e18 /
// 10^Fraction (≈ 9.2e16 for USD/EUR, ≈ 9.2e18 for JPY) to avoid silent
// truncation when reading the value back.
func NewFromMajorUnits[T Integer](amount T, code string) *Money {
	c := newCurrency(code).get()
	val := intToBigInt(amount)
	val.Mul(val, pow10(int64(c.Fraction)))
	return &Money{amount: &Decimal{val: val}, currency: c}
}

// NewFromFloat builds Money from a major-unit float (e.g. 1.25 → $1.25 for
// USD). Accepts float32 or float64. Sub-currency-unit precision is truncated.
// Prefer NewFromString when the value originates as text — float64 cannot
// represent every decimal exactly (0.1 + 0.2 != 0.3, etc.). The same int64
// truncation note from NewFromMajorUnits applies on read-back.
func NewFromFloat[T Float](amount T, code string) (*Money, error) {
	d, err := NewDecimalFromFloat(amount)
	if err != nil {
		return nil, err
	}
	return NewFromDecimal(d, code)
}

// NewFromString builds Money from a major-unit decimal string (e.g. "1.25").
// Sub-currency-unit precision is truncated toward zero. The same int64
// truncation note from NewFromMajorUnits applies on read-back.
func NewFromString(amount, code string) (*Money, error) {
	d, err := NewDecimalFromString(amount)
	if err != nil {
		return nil, err
	}
	return NewFromDecimal(d, code)
}

// NewFromDecimal wraps a Decimal interpreted as major currency units into a
// Money. Sub-currency-unit precision is truncated toward zero. The same
// int64 truncation note from NewFromMajorUnits applies on read-back.
func NewFromDecimal(d *Decimal, code string) (*Money, error) {
	if d == nil {
		return nil, errors.New("nil *Decimal")
	}
	c := newCurrency(code).get()
	return &Money{amount: normalizeToSmallestUnit(d, c.Fraction), currency: c}, nil
}

// normalizeToSmallestUnit converts a Decimal expressed in major units into an
// integer count of the currency's smallest unit, packaged as a Decimal at
// exponent 0 — matching the convention used by New() so that the existing
// Display / Amount / AsMajorUnits paths stay correct.
func normalizeToSmallestUnit(d *Decimal, fraction int) *Decimal {
	val := new(big.Int).Set(d.val)
	diff := fraction - d.exponent
	switch {
	case diff > 0:
		val.Mul(val, pow10(int64(diff)))
	case diff < 0:
		val.Quo(val, pow10(int64(-diff)))
	}
	return &Decimal{val: val}
}

func pow10(n int64) *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(n), nil)
}

// intToBigInt converts any Integer-typed value to a *big.Int. We dispatch on
// reflect.Kind (not the boxed concrete type) so named types like
// `type Cents uint64` still take the SetUint64 path — a value above
// math.MaxInt64 would otherwise be reinterpreted as negative through int64.
// Signed kinds and small unsigned kinds (uint8/16/32) take the SetInt64 path,
// which is exact for both since their value range fits in int64.
func intToBigInt[T Integer](v T) *big.Int {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return new(big.Int).SetUint64(rv.Uint())
	}
	return new(big.Int).SetInt64(rv.Int())
}
