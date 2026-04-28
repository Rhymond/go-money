package money

import (
	"fmt"
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

// NewDecimal builds a Decimal from any supported numeric input.
// Supported types: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
// float32, float64, string, *big.Int, *Decimal, *Money.
func NewDecimal(value interface{}) (*Decimal, error) {
	switch v := value.(type) {
	case int:
		return &Decimal{val: big.NewInt(int64(v))}, nil
	case int8:
		return &Decimal{val: big.NewInt(int64(v))}, nil
	case int16:
		return &Decimal{val: big.NewInt(int64(v))}, nil
	case int32:
		return &Decimal{val: big.NewInt(int64(v))}, nil
	case int64:
		return &Decimal{val: big.NewInt(v)}, nil
	case uint:
		return &Decimal{val: new(big.Int).SetUint64(uint64(v))}, nil
	case uint8:
		return &Decimal{val: big.NewInt(int64(v))}, nil
	case uint16:
		return &Decimal{val: big.NewInt(int64(v))}, nil
	case uint32:
		return &Decimal{val: big.NewInt(int64(v))}, nil
	case uint64:
		return &Decimal{val: new(big.Int).SetUint64(v)}, nil
	case float32:
		return decimalFromString(strconv.FormatFloat(float64(v), 'f', -1, 32))
	case float64:
		return decimalFromString(strconv.FormatFloat(v, 'f', -1, 64))
	case string:
		return decimalFromString(v)
	case *big.Int:
		if v == nil {
			return nil, fmt.Errorf("nil *big.Int")
		}
		return &Decimal{val: new(big.Int).Set(v)}, nil
	case *Decimal:
		if v == nil {
			return nil, fmt.Errorf("nil *Decimal")
		}
		return &Decimal{val: new(big.Int).Set(v.val), exponent: v.exponent}, nil
	case *Money:
		if v == nil || v.amount == nil {
			return nil, fmt.Errorf("nil *Money")
		}
		return &Decimal{val: new(big.Int).Set(v.amount.val), exponent: v.amount.exponent}, nil
	default:
		return nil, fmt.Errorf("unsupported value type %T", value)
	}
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
// space). Callers expecting smallest-currency-unit values should ensure
// Decimal.exponent matches the currency's fraction.
func (d *Decimal) Int64() int64 {
	return d.val.Int64()
}

// Sign returns -1 if d < 0, 0 if d == 0, +1 if d > 0.
func (d *Decimal) Sign() int {
	return d.val.Sign()
}
