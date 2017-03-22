package gocash

type Calculator struct{}

func (c *Calculator) add(a, b int) int {
	return a + b
}

func (c *Calculator) subtract(a, b int) int {
	return a - b
}

func (c *Calculator) multiply(a, m int) int {
	return a * m
}

func (c *Calculator) divide(a, d int) int {
	return a / d
}


