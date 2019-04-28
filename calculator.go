package money

import "math"

type calculator struct{}

func (c *calculator) add(a, b *Amount) *Amount {
	return &Amount{a.val + b.val}
}

func (c *calculator) subtract(a, b *Amount) *Amount {
	return &Amount{a.val - b.val}
}

func (c *calculator) multiply(a *Amount, m int64) *Amount {
	return &Amount{a.val * m}
}

func (c *calculator) divide(a *Amount, d int64) *Amount {
	return &Amount{a.val / d}
}

func (c *calculator) modulus(a *Amount, d int64) *Amount {
	return &Amount{a.val % d}
}

func (c *calculator) allocate(a *Amount, r, s int) *Amount {
	return &Amount{a.val * int64(r) / int64(s)}
}

func (c *calculator) absolute(a *Amount) *Amount {
	if a.val < 0 {
		return &Amount{-a.val}
	}

	return &Amount{a.val}
}

func (c *calculator) negative(a *Amount) *Amount {
	if a.val > 0 {
		return &Amount{-a.val}
	}

	return &Amount{a.val}
}

func (c *calculator) round(a *Amount, e int) *Amount {
	if a.val == 0 {
		return &Amount{0}
	}

	absam := c.absolute(a)
	exp := int64(math.Pow(10, float64(e)))
	m := absam.val % exp

	if m > (exp / 2) {
		absam.val += exp
	}

	absam.val = (absam.val / exp) * exp

	if a.val < 0 {
		a.val = -absam.val
	} else {
		a.val = absam.val
	}

	return &Amount{a.val}
}
