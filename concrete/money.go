package concrete

import (
	"errors"

	"github.com/Rhymond/go-money/currency"
	"github.com/Rhymond/go-money/interfaces"
)

// MoneyInt64 represents monetary value information, stores
// currency and amount value
type MoneyInt64 struct {
	amount   AmountInt64
	currency currency.Currency
}

// confirm MoneyInt64 implements Money
var _ interfaces.Money = (*MoneyInt64)(nil)

// NewMoneyInt64 creates a new int64-based money object
func NewMoneyInt64(amount interface{}, currencyCode string) (interfaces.Money, error) {
	amt, err := toAmount(amount)
	if err != nil {
		return nil, err
	}
	return MoneyInt64{
		amount:   amt,
		currency: currency.NewCurrency(currencyCode),
	}, nil
}

// Amount returns the internal amount object
func (m MoneyInt64) Amount() interfaces.Amount {
	return m.amount
}

// Currency returns the internal currency object
func (m MoneyInt64) Currency() currency.Currency {
	return m.currency
}

// SameCurrency check if given Money is equals by currency
func (m MoneyInt64) SameCurrency(om interfaces.Money) bool {
	return m.currency.Equals(om.Currency())
}

func (m MoneyInt64) assertSameCurrency(om interfaces.Money) error {
	if !m.SameCurrency(om) {
		return errors.New("Currencies don't match")
	}

	return nil
}

func (m MoneyInt64) compare(om interfaces.Money) int {
	switch {
	case m.Amount().Int64() > om.Amount().Int64():
		return 1
	case m.Amount().Int64() < om.Amount().Int64():
		return -1
	}

	return 0
}

// Equals checks equality between two Money types
func (m MoneyInt64) Equals(om interfaces.Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 0, nil
}

// GreaterThan checks whether the value of Money is greater than the other
func (m MoneyInt64) GreaterThan(om interfaces.Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 1, nil
}

// GreaterThanOrEqual checks whether the value of Money is greater or equal than the other
func (m MoneyInt64) GreaterThanOrEqual(om interfaces.Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) >= 0, nil
}

// LessThan checks whether the value of Money is less than the other
func (m MoneyInt64) LessThan(om interfaces.Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == -1, nil
}

// LessThanOrEqual checks whether the value of Money is less or equal than the other
func (m MoneyInt64) LessThanOrEqual(om interfaces.Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) <= 0, nil
}

// IsZero returns boolean of whether the value of Money is equals to zero
func (m MoneyInt64) IsZero() bool {
	return m.Amount().Int64() == 0
}

// IsPositive returns boolean of whether the value of Money is positive
func (m MoneyInt64) IsPositive() bool {
	return m.Amount().Int64() > 0
}

// IsNegative returns boolean of whether the value of Money is negative
func (m MoneyInt64) IsNegative() bool {
	return m.Amount().Int64() < 0
}

// Absolute returns new Money struct from given Money using absolute monetary value
func (m MoneyInt64) Absolute() interfaces.Money {
	return MoneyInt64{amount: m.amount.absolute(), currency: m.currency}
}

// Negative returns new Money struct from given Money using negative monetary value
func (m MoneyInt64) Negative() interfaces.Money {
	return MoneyInt64{amount: m.amount.negative(), currency: m.currency}
}

// Add returns new Money struct with value representing sum of Self and Other Money
//
// If the currencies don't match, the original Money object is returned
func (m MoneyInt64) Add(om interfaces.Money) interfaces.Money {
	if err := m.assertSameCurrency(om); err != nil {
		return m
	}
	return MoneyInt64{amount: m.amount.add(om.Amount().(AmountInt64)), currency: m.currency}
}

// Subtract returns new Money struct with value representing difference of Self and Other Money
//
// If the currencies don't match, the original Money object is returned
func (m MoneyInt64) Subtract(om interfaces.Money) interfaces.Money {
	if err := m.assertSameCurrency(om); err != nil {
		return m
	}
	return MoneyInt64{amount: m.amount.subtract(om.Amount().(AmountInt64)), currency: m.currency}
}

// Multiply returns new Money struct with value representing Self multiplied value by multiplier
func (m MoneyInt64) Multiply(mul interfaces.Money) interfaces.Money {
	return MoneyInt64{amount: m.amount.multiply(mul.Amount().(AmountInt64)), currency: m.currency}
}

// Divide returns new Money struct with value representing Self division value by given divider
func (m MoneyInt64) Divide(div interfaces.Money) interfaces.Money {
	return MoneyInt64{amount: m.amount.divide(div.Amount().(AmountInt64)), currency: m.currency}
}

// Round returns new Money struct with value rounded to nearest zero
func (m MoneyInt64) Round() interfaces.Money {
	return MoneyInt64{amount: m.amount.round(), currency: m.currency}
}

// Split returns slice of Money structs with split Self value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first will likely receive more pennies than ones that are listed later
func (m MoneyInt64) Split(n int) ([]interfaces.Money, error) {
	if n <= 0 {
		return nil, errors.New("Split must be higher than zero")
	}

	num, err := toAmount(n)
	if err != nil {
		return nil, err
	}
	a := m.amount.divide(num)
	ms := make([]MoneyInt64, n)

	for i := 0; i < n; i++ {
		ms[i] = MoneyInt64{a, m.currency}
	}

	l := m.amount.modulus(num).Int64()

	// Add leftovers to the first parties
	for p := 0; l != 0; p++ {
		ms[p].amount = ms[p].amount.add(AmountInt64{value: 1})
		l--
	}

	ret := make([]interfaces.Money, len(ms))
	for i, v := range ms {
		ret[i] = v
	}

	return ret, nil
}

// Allocate returns slice of Money structs with split Self value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
func (m MoneyInt64) Allocate(rs []int) ([]interfaces.Money, error) {
	if len(rs) == 0 {
		return nil, errors.New("No ratios specified")
	}

	// Calculate sum of ratios
	var sum int
	for _, r := range rs {
		sum += r
	}

	var total int64
	var ms []MoneyInt64
	for _, r := range rs {
		rA, _ := toAmount(r)
		sumA, _ := toAmount(sum)
		party := MoneyInt64{
			m.amount.allocate(rA, sumA),
			m.currency,
		}

		ms = append(ms, party)
		total += party.Amount().Int64()
	}

	// Calculate leftover value and divide to first parties
	lo := m.Amount().Int64() - total
	sub := int64(1)
	if lo < 0 {
		sub = -1
	}

	for p := 0; lo != 0; p++ {
		ms[p].amount = ms[p].amount.add(AmountInt64{value: sub})
		lo -= sub
	}

	ret := make([]interfaces.Money, len(ms))
	for i, v := range ms {
		ret[i] = v
	}

	return ret, nil
}

// Display lets represent Money struct as string in given Currency value
func (m MoneyInt64) Display() string {
	c := m.Currency().Get()
	return c.Formatter().Format(m.Amount().Int64())
}
