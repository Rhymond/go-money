package money

import "math"

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

func (c *calculator) allocate(a Amount, r, s uint) Amount {
	if a == 0 || s == 0 {
		return 0
	}

	return a * int64(r) / int64(s)
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
