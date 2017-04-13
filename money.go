package money

import (
	"errors"
	"fmt"
	"reflect"
)

// Money represents monetary value information, stores
// currency and amount value
type Money struct {
	Amount   Amount
	currency Currency
}

func toAmount(amount interface{}) Amount {
	switch amt := amount.(type) {
	case uint:
		return Amount{value: int64(amt)}
	case uint8:
		return Amount{value: int64(amt)}
	case uint16:
		return Amount{value: int64(amt)}
	case uint32:
		return Amount{value: int64(amt)}
	case uint64:
		return Amount{value: int64(amt)}
	case int:
		return Amount{value: int64(amt)}
	case int8:
		return Amount{value: int64(amt)}
	case int16:
		return Amount{value: int64(amt)}
	case int32:
		return Amount{value: int64(amt)}
	case int64:
		return Amount{value: amt}
	}
	panic(fmt.Sprintf("Unable to convert to Amount. Unsupported type: %s", reflect.TypeOf(amount).Name()))
}

// New creates and returns new instance of Money
func New(amount interface{}, code string) Money {
	return Money{
		Amount:   toAmount(amount),
		currency: newCurrency(code),
	}
}

// SameCurrency check if given Money is equals by currency
func (m Money) SameCurrency(om Money) bool {
	return m.currency.equals(om.currency)
}

func (m Money) assertSameCurrency(om Money) error {
	if !m.SameCurrency(om) {
		return errors.New("Currencies don't match")
	}

	return nil
}

func (m Money) compare(om Money) int {
	switch {
	case m.Amount.value > om.Amount.value:
		return 1
	case m.Amount.value < om.Amount.value:
		return -1
	}

	return 0
}

// Equals checks equality between two Money types
func (m Money) Equals(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 0, nil
}

// GreaterThan checks whether the value of Money is greater than the other
func (m Money) GreaterThan(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 1, nil
}

// GreaterThanOrEqual checks whether the value of Money is greater or equal than the other
func (m Money) GreaterThanOrEqual(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) >= 0, nil
}

// LessThan checks whether the value of Money is less than the other
func (m Money) LessThan(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == -1, nil
}

// LessThanOrEqual checks whether the value of Money is less or equal than the other
func (m Money) LessThanOrEqual(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) <= 0, nil
}

// IsZero returns boolean of whether the value of Money is equals to zero
func (m Money) IsZero() bool {
	return m.Amount.value == 0
}

// IsPositive returns boolean of whether the value of Money is positive
func (m Money) IsPositive() bool {
	return m.Amount.value > 0
}

// IsNegative returns boolean of whether the value of Money is negative
func (m Money) IsNegative() bool {
	return m.Amount.value < 0
}

// Absolute returns new Money struct from given Money using absolute monetary value
func (m Money) Absolute() Money {
	return Money{Amount: m.Amount.absolute(), currency: m.currency}
}

// Negative returns new Money struct from given Money using negative monetary value
func (m Money) Negative() Money {
	return Money{Amount: m.Amount.negative(), currency: m.currency}
}

// Add returns new Money struct with value representing sum of Self and Other Money
//
// If the currencies don't match, the original Money object is returned
func (m Money) Add(om Money) Money {
	if err := m.assertSameCurrency(om); err != nil {
		return m
	}
	return Money{Amount: m.Amount.add(om.Amount), currency: m.currency}
}

// Subtract returns new Money struct with value representing difference of Self and Other Money
//
// If the currencies don't match, the original Money object is returned
func (m Money) Subtract(om Money) Money {
	if err := m.assertSameCurrency(om); err != nil {
		return m
	}
	return Money{Amount: m.Amount.subtract(om.Amount), currency: m.currency}
}

// Multiply returns new Money struct with value representing Self multiplied value by multiplier
func (m Money) Multiply(mul interface{}) Money {
	return Money{Amount: m.Amount.multiply(toAmount(mul)), currency: m.currency}
}

// Divide returns new Money struct with value representing Self division value by given divider
func (m Money) Divide(div interface{}) Money {
	return Money{Amount: m.Amount.divide(toAmount(div)), currency: m.currency}
}

// Round returns new Money struct with value rounded to nearest zero
func (m Money) Round() Money {
	return Money{Amount: m.Amount.round(), currency: m.currency}
}

// Split returns slice of Money structs with split Self value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first will likely receive more pennies than ones that are listed later
func (m Money) Split(n int) ([]Money, error) {
	if n <= 0 {
		return nil, errors.New("Split must be higher than zero")
	}

	num := toAmount(n)
	a := m.Amount.divide(num)
	ms := make([]Money, n)

	for i := 0; i < n; i++ {
		ms[i] = Money{a, m.currency}
	}

	l := m.Amount.modulus(num).value

	// Add leftovers to the first parties
	for p := 0; l != 0; p++ {
		ms[p].Amount = ms[p].Amount.add(Amount{value: 1})
		l--
	}

	return ms, nil
}

// Allocate returns slice of Money structs with split Self value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
func (m Money) Allocate(rs []int) ([]Money, error) {
	if len(rs) == 0 {
		return nil, errors.New("No ratios specified")
	}

	// Calculate sum of ratios
	var sum int
	for _, r := range rs {
		sum += r
	}

	var total int64
	var ms []Money
	for _, r := range rs {
		party := Money{
			m.Amount.allocate(toAmount(r), toAmount(sum)),
			m.currency,
		}

		ms = append(ms, party)
		total += party.Amount.value
	}

	// Calculate leftover value and divide to first parties
	lo := m.Amount.value - total
	sub := int64(1)
	if lo < 0 {
		sub = -1
	}

	for p := 0; lo != 0; p++ {
		ms[p].Amount = ms[p].Amount.add(Amount{value: sub})
		lo -= sub
	}

	return ms, nil
}

// Display lets represent Money struct as string in given Currency value
func (m Money) Display() string {
	c := m.currency.get()
	return c.Formatter().Format(m.Amount.value)
}
