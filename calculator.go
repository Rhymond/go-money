package money

type Calculator struct{}

func (c *Calculator) add(a, b *Amount) *Amount {
	return &Amount{a.val + b.val}
}

func (c *Calculator) subtract(a, b *Amount) *Amount {
	return &Amount{a.val - b.val}
}

func (c *Calculator) multiply(a *Amount, m int) *Amount {
	return &Amount{a.val * m}
}

func (c *Calculator) divide(a *Amount, d int) *Amount {
	return &Amount{a.val / d}
}

func (c *Calculator) modulus(a *Amount, d int) *Amount {
	return &Amount{a.val % d}
}

func (c *Calculator) allocate(a *Amount, r, s int) *Amount {
	return &Amount{a.val * r / s}
}

func (c *Calculator) absolute(a *Amount) *Amount {
	if a.val < 0 {
		return &Amount{-a.val}
	}

	return &Amount{a.val}
}

func (c *Calculator) negative(a *Amount) *Amount {
	if a.val > 0 {
		return &Amount{-a.val}
	}

	return &Amount{a.val}
}

func (c *Calculator) round(a *Amount) *Amount {

	if a.val == 0 {
		return &Amount{0}
	}

	absam := calc.absolute(a)
	m := absam.val % 100

	if m > 50 {
		absam.val += 100
	}

	absam.val = (absam.val / 100) * 100

	if a.val < 0 {
		a.val = -absam.val
	} else {
		a.val = absam.val
	}

	return &Amount{a.val}
}
