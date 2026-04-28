package money

import "math/big"

type calculator struct{}

func (c *calculator) add(a, b *Decimal) *Decimal {
	av, bv, exp := align(a, b)
	return &Decimal{val: av.Add(av, bv), exponent: exp}
}

func (c *calculator) subtract(a, b *Decimal) *Decimal {
	av, bv, exp := align(a, b)
	return &Decimal{val: av.Sub(av, bv), exponent: exp}
}

func (c *calculator) multiply(a, b *Decimal) *Decimal {
	return &Decimal{
		val:      new(big.Int).Mul(a.val, b.val),
		exponent: a.exponent + b.exponent,
	}
}

func (c *calculator) divide(a *Decimal, d int64) *Decimal {
	val := new(big.Int).Quo(a.val, big.NewInt(d))
	return &Decimal{val: val, exponent: a.exponent}
}

func (c *calculator) modulus(a *Decimal, d int64) *Decimal {
	val := new(big.Int).Rem(a.val, big.NewInt(d))
	return &Decimal{val: val, exponent: a.exponent}
}

// allocate computes a * r / s and returns the result at exponent a.exponent
// so the caller's leftover arithmetic (which compares against a.val directly)
// stays valid.
func (c *calculator) allocate(a, r, s *Decimal) *Decimal {
	if s.val.Sign() == 0 {
		return &Decimal{val: new(big.Int), exponent: a.exponent}
	}

	num := new(big.Int).Mul(a.val, r.val)

	// The natural exponent of num/s is a.exponent + r.exponent - s.exponent;
	// we want a.exponent. Adjust the numerator by 10^(s.exponent - r.exponent)
	// before dividing.
	diff := int64(s.exponent - r.exponent)
	if diff > 0 {
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(diff), nil)
		num.Mul(num, scale)
	} else if diff < 0 {
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(-diff), nil)
		num.Quo(num, scale)
	}

	val := new(big.Int).Quo(num, s.val)
	return &Decimal{val: val, exponent: a.exponent}
}

func (c *calculator) absolute(a *Decimal) *Decimal {
	return &Decimal{val: new(big.Int).Abs(a.val), exponent: a.exponent}
}

func (c *calculator) negative(a *Decimal) *Decimal {
	if a.val.Sign() > 0 {
		return &Decimal{val: new(big.Int).Neg(a.val), exponent: a.exponent}
	}
	return &Decimal{val: new(big.Int).Set(a.val), exponent: a.exponent}
}

func (c *calculator) round(a *Decimal, e int) *Decimal {
	if a.val.Sign() == 0 {
		return &Decimal{val: new(big.Int), exponent: a.exponent}
	}

	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(e)), nil)
	abs := new(big.Int).Abs(a.val)
	half := new(big.Int).Quo(exp, big.NewInt(2))

	rem := new(big.Int).Rem(abs, exp)
	if rem.Cmp(half) > 0 {
		abs.Add(abs, exp)
	}
	abs.Quo(abs, exp).Mul(abs, exp)

	if a.val.Sign() < 0 {
		abs.Neg(abs)
	}
	return &Decimal{val: abs, exponent: a.exponent}
}
