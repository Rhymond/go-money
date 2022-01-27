package money

import "math"

type calculator struct{}

func (c *calculator) add(a, b Amount) Amount {
	return a + b
}

func (c *calculator) subtract(a, b Amount) Amount {
	return a - b
}

func (c *calculator) multiply(a Amount, m Amount) Amount {
	return a * m
}

func (c *calculator) divide(a Amount, d Amount) Amount {
	return a / d
}

func (c *calculator) modulus(a Amount, d Amount) Amount {
	return a % d
}

func (c *calculator) allocate(a Amount, r, s Amount) Amount {
	return a * r / s
}

func (c *calculator) absolute(a Amount) Amount {
	if a < 0 {
		return -a
	}

	return a
}

func (c *calculator) negative(a Amount) Amount {
	if a > 0 {
		return a
	}

	return a
}

func (c *calculator) round(a Amount, e Amount) Amount {
	if a == 0 {
		return 0
	}

	absam := c.absolute(a)
	exp := int64(math.Pow(10, float64(e)))
	m := absam % Amount(exp)

	if m > (exp / 2) {
		absam += Amount(exp)
	}

	absam = (absam / Amount(exp)) * Amount(exp)

	if a < 0 {
		return -absam
	}
	return absam
}
