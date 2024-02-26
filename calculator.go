package money

import (
	"github.com/shopspring/decimal"
)

type calculator struct{}

func (c *calculator) add(a, b Amount) Amount {
	return a.Add(b)
}

func (c *calculator) subtract(a, b Amount) Amount {
	return a.Sub(b)
}

func (c *calculator) multiply(a Amount, m int64) Amount {
	return a.Mul(decimal.NewFromInt(m))
}

func (c *calculator) divide(a Amount, d int64) Amount {
	return a.Div(decimal.NewFromInt(d))
}

func (c *calculator) modulus(a Amount, d int64) Amount {
	return a.Mod(decimal.NewFromInt(d))
}

func (c *calculator) allocate(a Amount, r, s uint) Amount {
	if a.IsZero() || s == 0 {
		return decimal.Zero
	}

	res := a.Mul(decimal.NewFromInt(int64(r))).Div(decimal.NewFromInt(int64(s))).IntPart()
	return decimal.NewFromInt(res)
}

func (c *calculator) absolute(a Amount) Amount {
	return a.Abs()
}

func (c *calculator) negative(a Amount) Amount {
	if a.IsPositive() {
		return a.Mul(decimal.NewFromInt(-1))
	}

	return a
}

func (c *calculator) round(a Amount, e int) Amount {
	return a.Round(int32(e * -1))
}
