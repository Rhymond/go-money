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

func (c *calculator) allocate(a Amount, r, s int64) Amount {
	if a == 0 || s == 0 {
		return 0
	}

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
		return -a
	}

	return a
}

func (c *calculator) round(a Amount, e int) Amount {
	if a == 0 {
		return 0
	}

	// Round using signed quotient and remainder so MinInt64 does not overflow
	// when converted to an absolute value. Preserve the existing tie behavior:
	// exact halves round toward zero.
	exp := int64(math.Pow(10, float64(e)))
	quotient := a / exp
	remainder := a % exp
	if remainder > exp/2 {
		quotient++
	} else if remainder < -(exp / 2) {
		quotient--
	}

	return quotient * exp
}
