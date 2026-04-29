package money

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Decimal represents an arbitrary-precision decimal number with the
// value val * 10^(-exponent).
type Decimal struct {
	val      *big.Int
	exponent int
}

// NewDecimalFromInt builds a Decimal from any signed or unsigned integer width.
// Always succeeds — every integer value fits exactly.
func NewDecimalFromInt[T Integer](v T) *Decimal {
	return &Decimal{val: intToBigInt(v)}
}

// NewDecimalFromFloat builds a Decimal from a float32 or float64 using its
// shortest round-trip decimal representation. Returns an error for NaN or ±Inf.
func NewDecimalFromFloat[T Float](v T) (*Decimal, error) {
	bitSize := 64
	if _, ok := any(v).(float32); ok {
		bitSize = 32
	}
	return decimalFromFloat(float64(v), bitSize)
}

// NewDecimalFromString parses a decimal string like "1.25" or "-3" into a Decimal.
func NewDecimalFromString(s string) (*Decimal, error) {
	return decimalFromString(s)
}

// NewDecimalFromBigInt wraps a *big.Int into a Decimal. The input is defensively
// copied so later mutations to the source don't bleed in.
func NewDecimalFromBigInt(b *big.Int) (*Decimal, error) {
	if b == nil {
		return nil, errors.New("nil *big.Int")
	}
	return &Decimal{val: new(big.Int).Set(b)}, nil
}

// NewDecimalFromMoney returns the Money's amount as a Decimal. The underlying
// value is defensively copied.
func NewDecimalFromMoney(m *Money) (*Decimal, error) {
	if m == nil || m.amount == nil {
		return nil, errors.New("nil *Money")
	}
	return &Decimal{val: new(big.Int).Set(m.amount.val), exponent: m.amount.exponent}, nil
}

// decimalFromFloat is the single entry point for float → Decimal conversion.
// It rejects NaN and ±Inf with a clear error before delegating to the string
// parser; otherwise FormatFloat would emit "NaN"/"+Inf" and the failure would
// surface as a generic "invalid number" message.
func decimalFromFloat(v float64, bitSize int) (*Decimal, error) {
	if math.IsNaN(v) {
		return nil, fmt.Errorf("NaN is not a valid decimal")
	}
	if math.IsInf(v, 0) {
		return nil, fmt.Errorf("infinity is not a valid decimal")
	}
	return decimalFromString(strconv.FormatFloat(v, 'f', -1, bitSize))
}

func decimalFromString(s string) (*Decimal, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty string")
	}

	neg := false
	switch s[0] {
	case '+':
		s = s[1:]
	case '-':
		neg = true
		s = s[1:]
	}
	if s == "" {
		return nil, fmt.Errorf("invalid number: missing digits")
	}

	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}
	if intPart == "" {
		intPart = "0"
	}

	digits := intPart + fracPart
	if digits == "" {
		return nil, fmt.Errorf("invalid number: no digits")
	}
	for _, r := range digits {
		if r < '0' || r > '9' {
			return nil, fmt.Errorf("invalid number: %q", s)
		}
	}
	val, ok := new(big.Int).SetString(digits, 10)
	if !ok {
		return nil, fmt.Errorf("invalid number: %q", s)
	}
	if neg {
		val.Neg(val)
	}
	return &Decimal{val: val, exponent: len(fracPart)}, nil
}

// align returns the underlying integer values of a and b scaled to a common
// exponent (the larger of the two). Original Decimals are not mutated.
func align(a, b *Decimal) (av, bv *big.Int, exp int) {
	if a.exponent == b.exponent {
		return new(big.Int).Set(a.val), new(big.Int).Set(b.val), a.exponent
	}
	if a.exponent > b.exponent {
		diff := int64(a.exponent - b.exponent)
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(diff), nil)
		return new(big.Int).Set(a.val), new(big.Int).Mul(b.val, scale), a.exponent
	}
	diff := int64(b.exponent - a.exponent)
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(diff), nil)
	return new(big.Int).Mul(a.val, scale), new(big.Int).Set(b.val), b.exponent
}

// Int64 returns the underlying integer value as int64 (in the current exponent
// space). Callers expecting smallest-currency-unit values should use
// Money.Amount, which is exponent-aware. The result silently truncates if val
// exceeds the int64 range — see (*big.Int).Int64.
func (d *Decimal) Int64() int64 {
	return d.val.Int64()
}

// Sign returns -1 if d < 0, 0 if d == 0, +1 if d > 0.
func (d *Decimal) Sign() int {
	return d.val.Sign()
}

// add returns a new Decimal equal to d + b, scaled to the larger of the two
// exponents.
func (d *Decimal) add(b *Decimal) *Decimal {
	av, bv, exp := align(d, b)
	return &Decimal{val: av.Add(av, bv), exponent: exp}
}

// subtract returns a new Decimal equal to d - b, scaled to the larger of the
// two exponents.
func (d *Decimal) subtract(b *Decimal) *Decimal {
	av, bv, exp := align(d, b)
	return &Decimal{val: av.Sub(av, bv), exponent: exp}
}

// multiply returns a new Decimal equal to d * b. The result's exponent is the
// sum of the two operand exponents (full precision is preserved).
func (d *Decimal) multiply(b *Decimal) *Decimal {
	return &Decimal{
		val:      new(big.Int).Mul(d.val, b.val),
		exponent: d.exponent + b.exponent,
	}
}

// divide returns a new Decimal equal to d / n, integer-divided in val space
// (truncated toward zero). The exponent is preserved.
func (d *Decimal) divide(n int64) *Decimal {
	val := new(big.Int).Quo(d.val, big.NewInt(n))
	return &Decimal{val: val, exponent: d.exponent}
}

// modulus returns a new Decimal equal to d mod n in val space. The exponent
// is preserved.
func (d *Decimal) modulus(n int64) *Decimal {
	val := new(big.Int).Rem(d.val, big.NewInt(n))
	return &Decimal{val: val, exponent: d.exponent}
}

// allocate returns a new Decimal equal to d * r / s, used by Money.Allocate
// to compute each party's share. Returns zero when s == 0 to avoid division
// by zero. The exponent is preserved.
func (d *Decimal) allocate(r, s int64) *Decimal {
	if s == 0 {
		return &Decimal{val: new(big.Int), exponent: d.exponent}
	}
	val := new(big.Int).Mul(d.val, big.NewInt(r))
	val.Quo(val, big.NewInt(s))
	return &Decimal{val: val, exponent: d.exponent}
}

// absolute returns a new Decimal equal to |d|.
func (d *Decimal) absolute() *Decimal {
	return &Decimal{val: new(big.Int).Abs(d.val), exponent: d.exponent}
}

// negative returns a new Decimal equal to -|d| (zero stays zero).
func (d *Decimal) negative() *Decimal {
	if d.val.Sign() > 0 {
		return &Decimal{val: new(big.Int).Neg(d.val), exponent: d.exponent}
	}
	return &Decimal{val: new(big.Int).Set(d.val), exponent: d.exponent}
}

// round returns a new Decimal with val rounded to the nearest 10^e in its
// own exponent space using half-away-from-zero. The exponent is preserved.
// e.g. round(2) on val=150 yields 200; round(2) on val=-150 yields -200.
func (d *Decimal) round(e int) *Decimal {
	if d.val.Sign() == 0 {
		return &Decimal{val: new(big.Int), exponent: d.exponent}
	}

	exp := pow10(int64(e))
	abs := new(big.Int).Abs(d.val)
	half := new(big.Int).Quo(exp, big.NewInt(2))

	rem := new(big.Int).Rem(abs, exp)
	if rem.Cmp(half) >= 0 {
		abs.Add(abs, exp)
	}
	abs.Quo(abs, exp).Mul(abs, exp)

	if d.val.Sign() < 0 {
		abs.Neg(abs)
	}
	return &Decimal{val: abs, exponent: d.exponent}
}

// roundToSmallestUnit returns a new Decimal with sub-smallest-unit precision
// removed using half-away-from-zero. The result has exponent 0, so callers
// like Money.Amount() see the rounded value directly. For exponent ≤ 0 the
// value already lives at smallest-unit precision and is returned with
// exponent 0 unchanged.
func (d *Decimal) roundToSmallestUnit() *Decimal {
	if d.exponent <= 0 {
		return &Decimal{val: new(big.Int).Set(d.val), exponent: 0}
	}

	scale := pow10(int64(d.exponent))
	half := new(big.Int).Quo(scale, big.NewInt(2))

	abs := new(big.Int).Abs(d.val)
	quo, rem := new(big.Int), new(big.Int)
	quo.QuoRem(abs, scale, rem)
	if rem.Cmp(half) >= 0 {
		quo.Add(quo, big.NewInt(1))
	}
	if d.val.Sign() < 0 {
		quo.Neg(quo)
	}
	return &Decimal{val: quo, exponent: 0}
}
