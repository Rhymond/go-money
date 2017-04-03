package money

type calculator struct{}

func (c *calculator) add(a, b *amount) *amount {
	return &amount{a.val + b.val}
}

func (c *calculator) subtract(a, b *amount) *amount {
	return &amount{a.val - b.val}
}

func (c *calculator) multiply(a *amount, m int) *amount {
	return &amount{a.val * m}
}

func (c *calculator) divide(a *amount, d int) *amount {
	return &amount{a.val / d}
}

func (c *calculator) modulus(a *amount, d int) *amount {
	return &amount{a.val % d}
}

func (c *calculator) allocate(a *amount, r, s int) *amount {
	return &amount{a.val * r / s}
}

func (c *calculator) absolute(a *amount) *amount {
	if a.val < 0 {
		return &amount{-a.val}
	}

	return &amount{a.val}
}

func (c *calculator) negative(a *amount) *amount {
	if a.val > 0 {
		return &amount{-a.val}
	}

	return &amount{a.val}
}

func (c *calculator) round(a *amount) *amount {

	if a.val == 0 {
		return &amount{0}
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

	return &amount{a.val}
}
