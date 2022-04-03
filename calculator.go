package money

import "math"

type calculator[T Numeric] struct{}

func (c *calculator[T]) add(a, b *Amount[T]) *Amount[T] {
	return &Amount[T]{a.val + b.val}
}

func (c *calculator[T]) subtract(a, b *Amount[T]) *Amount[T] {
	return &Amount[T]{a.val - b.val}
}

func (c *calculator[T]) multiply(a *Amount[T], m T) *Amount[T] {
	return &Amount[T]{a.val * m}
}

func (c *calculator[T]) divide(a *Amount[T], d T) *Amount[T] {
	return &Amount[T]{a.val / d}
}

func (c *calculator[T]) modulus(a *Amount[T], d T) *Amount[T] {
	return &Amount[T]{a.val % d}
}

func (c *calculator[T]) allocate(a *Amount[T], r, s T) *Amount[T] {
	return &Amount[T]{a.val * r / s}
}

func (c *calculator[T]) absolute(a *Amount[T]) *Amount[T] {
	if a.val < 0 {
		return &Amount[T]{-a.val}
	}

	return &Amount[T]{a.val}
}

func (c *calculator[T]) negative(a *Amount[T]) *Amount[T] {
	if a.val > 0 {
		return &Amount[T]{-a.val}
	}

	return &Amount[T]{a.val}
}

func (c *calculator[T]) round(a *Amount[T], e int) *Amount[T] {
	if a.val == 0 {
		return &Amount[T]{0}
	}

	absam := c.absolute(a)
	exp := T(math.Pow(10, float64(e)))
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

	return &Amount[T]{a.val}
}
