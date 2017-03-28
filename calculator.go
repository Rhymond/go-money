package money

import (
	"math"
)

type Calculator struct{}

func (c *Calculator) add(a, b *Number) *Number {
	return &Number{a.Amount + b.Amount}
}

func (c *Calculator) subtract(a, b *Number) *Number {
	return &Number{a.Amount - b.Amount}
}

func (c *Calculator) multiply(n *Number, m int) *Number {
	return &Number{n.Amount * m}
}

func (c *Calculator) divide(n *Number, d int) *Number {
	return &Number{n.Amount / d}
}

func (c *Calculator) modulus(n *Number, d int) *Number {
	return &Number{n.Amount % d}
}

func (c *Calculator) allocate(n *Number, r, s int) *Number {
	return &Number{n.Amount * r / s}
}

func (c *Calculator) round(n *Number) *Number {

	input := float64(n.Amount) / 100
	var o float64

	if input < 0 {
		o = math.Ceil(input - 0.5)
	} else {
		o = math.Floor(input + 0.5)
	}

	return &Number{int(o) * 100}
}