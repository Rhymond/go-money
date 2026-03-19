package money

import "errors"

// Sum returns a new Money equal to the total of all given values.
// All values must share the same currency.
// Returns an error when no values are provided or currencies differ.
//
//	a := money.New(100, "USD")
//	b := money.New(200, "USD")
//	c := money.New(300, "USD")
//	total, _ := money.Sum(a, b, c) // $6.00
func Sum(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return nil, errors.New("sum requires at least one value")
	}
	result := ms[0]
	var err error
	for _, m := range ms[1:] {
		result, err = result.Add(m)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// Min returns the smallest Money value from the given set.
// All values must share the same currency.
// Returns an error when no values are provided or currencies differ.
//
//	a := money.New(100, "USD")
//	b := money.New(200, "USD")
//	c := money.New(50, "USD")
//	smallest, _ := money.Min(a, b, c) // $0.50
func Min(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return nil, errors.New("min requires at least one value")
	}
	min := ms[0]
	for _, m := range ms[1:] {
		lt, err := m.LessThan(min)
		if err != nil {
			return nil, err
		}
		if lt {
			min = m
		}
	}
	return min, nil
}

// Max returns the largest Money value from the given set.
// All values must share the same currency.
// Returns an error when no values are provided or currencies differ.
//
//	a := money.New(100, "USD")
//	b := money.New(200, "USD")
//	c := money.New(50, "USD")
//	largest, _ := money.Max(a, b, c) // $2.00
func Max(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return nil, errors.New("max requires at least one value")
	}
	max := ms[0]
	for _, m := range ms[1:] {
		gt, err := m.GreaterThan(max)
		if err != nil {
			return nil, err
		}
		if gt {
			max = m
		}
	}
	return max, nil
}

// Average returns the arithmetic mean of the given Money values,
// truncating any remainder to the smallest currency unit.
// All values must share the same currency.
// Returns an error when no values are provided or currencies differ.
//
//	a := money.New(100, "USD")
//	b := money.New(200, "USD")
//	c := money.New(300, "USD")
//	avg, _ := money.Average(a, b, c) // $2.00
func Average(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return nil, errors.New("average requires at least one value")
	}
	total, err := Sum(ms...)
	if err != nil {
		return nil, err
	}
	return total.Divide(int64(len(ms))), nil
}
