package money

import (
	"math"
	"math/big"
)

type calculator struct{}

func (c *calculator) add(a, b Amount) Amount {
	return a + b
}

func (c *calculator) subtract(a, b Amount) Amount {
	return a - b
}

func (c *calculator) multiply(a Amount, m int64) Amount {
	return a * m
}

func (c *calculator) divide(a Amount, d int64) Amount {
	return a / d
}

func (c *calculator) modulus(a Amount, d int64) Amount {
	return a % d
}

func (c *calculator) allocate(a Amount, r, s int64) Amount {
	if a == 0 || s == 0 {
		return 0
	}

	// a*r can overflow int64 (e.g. Allocate on a large amount with a large
	// single ratio). Compute with big.Int; result always fits when r <= s.
	bigA := new(big.Int).SetInt64(a)
	bigA.Mul(bigA, big.NewInt(r))
	bigA.Quo(bigA, big.NewInt(s))
	return Amount(bigA.Int64())
}

func (c *calculator) absolute(a Amount) Amount {
	if a < 0 {
		return -a
	}

	return a
}

func (c *calculator) negative(a Amount) Amount {
	if a > 0 {
		return -a
	}

	return a
}

func (c *calculator) round(a Amount, e int) Amount {
	if a == 0 {
		return 0
	}

	absam := c.absolute(a)
	exp := int64(math.Pow(10, float64(e)))
	m := absam % exp

	if m > (exp / 2) {
		absam += exp
	}

	absam = (absam / exp) * exp

	if a < 0 {
		a = -absam
	} else {
		a = absam
	}

	return a
}
