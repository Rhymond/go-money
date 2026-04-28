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

func (c *calculator) allocate(a *Decimal, r, s int64) *Decimal {
	if s == 0 {
		return &Decimal{val: new(big.Int), exponent: a.exponent}
	}
	val := new(big.Int).Mul(a.val, big.NewInt(r))
	val.Quo(val, big.NewInt(s))
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
