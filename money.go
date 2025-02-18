package money

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

// Injection points for backward compatibility.
// If you need to keep your JSON marshal/unmarshal way, overwrite them like below.
//
//	money.UnmarshalJSON = func (m *Money, b []byte) error { ... }
//	money.MarshalJSON = func (m Money) ([]byte, error) { ... }
var (
	// UnmarshalJSON is injection point of json.Unmarshaller for money.Money
	UnmarshalJSON = defaultUnmarshalJSON
	// MarshalJSON is injection point of json.Marshaller for money.Money
	MarshalJSON = defaultMarshalJSON

	// ErrCurrencyMismatch happens when two compared Money don't have the same currency.
	ErrCurrencyMismatch = errors.New("currencies don't match")

	// ErrInvalidJSONUnmarshal happens when the default money.UnmarshalJSON fails to unmarshal Money because of invalid data.
	ErrInvalidJSONUnmarshal = errors.New("invalid json unmarshal")

	ErrInvalidCurrency = errors.New("invalid currency passed")
)

func defaultUnmarshalJSON(m *Money, b []byte) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	var amount float64
	if amountRaw, ok := data["amount"]; ok {
		amount, ok = amountRaw.(float64)
		if !ok {
			return ErrInvalidJSONUnmarshal
		}
	}

	var currency string
	if currencyRaw, ok := data["currency"]; ok {
		currency, ok = currencyRaw.(string)
		if !ok {
			return ErrInvalidJSONUnmarshal
		}
		if c := GetCurrency(currency); c == nil {
			return ErrInvalidCurrency
		}
	}

	var ref *Money
	if amount == 0 && currency == "" {
		ref = &Money{}
	} else {
		ref = New(int64(amount), currency)
	}

	*m = *ref
	return nil
}

func defaultMarshalJSON(m Money) ([]byte, error) {
	if m == (Money{}) {
		m = *New(0, "")
	}

	buff := bytes.NewBufferString(fmt.Sprintf(`{"amount": %d, "currency": "%s"}`, m.Amount(), m.Currency().Code))
	return buff.Bytes(), nil
}

// Amount is a data structure that stores the amount being used for calculations.
type Amount = int64

// Money represents monetary value information, stores
// currency and amount value.
type Money struct {
	Amount_   Amount    `json:"amount" swaggertype:"primitive,integer"`
	Currency_ *Currency `json:"currency" swaggertype:"primitive,string"`
}

// New creates and returns new instance of Money.
func New(amount int64, code string) *Money {
	return &Money{
		Amount_:   amount,
		Currency_: newCurrency(code).get(),
	}
}

// NewFromFloat creates and returns new instance of Money from a float64.
// Always rounding trailing decimals down.
func NewFromFloat(amount float64, code string) *Money {
	currencyDecimals := math.Pow10(newCurrency(code).get().Fraction)
	return New(int64(amount*currencyDecimals), code)
}

// Currency returns the currency used by Money.
func (m *Money) Currency() *Currency {
	return m.Currency_
}

// Amount returns a copy of the internal monetary value as an int64.
func (m *Money) Amount() int64 {
	return m.Amount_
}

// SameCurrency check if given Money is equals by currency.
func (m *Money) SameCurrency(om *Money) bool {
	return m.Currency_.equals(om.Currency_)
}

func (m *Money) assertSameCurrency(om *Money) error {
	if !m.SameCurrency(om) {
		return ErrCurrencyMismatch
	}

	return nil
}

func (m *Money) compare(om *Money) int {
	switch {
	case m.Amount_ > om.Amount_:
		return 1
	case m.Amount_ < om.Amount_:
		return -1
	}

	return 0
}

// Equals checks equality between two Money types.
func (m *Money) Equals(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 0, nil
}

// GreaterThan checks whether the value of Money is greater than the other.
func (m *Money) GreaterThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 1, nil
}

// GreaterThanOrEqual checks whether the value of Money is greater or equal than the other.
func (m *Money) GreaterThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) >= 0, nil
}

// LessThan checks whether the value of Money is less than the other.
func (m *Money) LessThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == -1, nil
}

// LessThanOrEqual checks whether the value of Money is less or equal than the other.
func (m *Money) LessThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) <= 0, nil
}

// IsZero returns boolean of whether the value of Money is equals to zero.
func (m *Money) IsZero() bool {
	return m.Amount_ == 0
}

// IsPositive returns boolean of whether the value of Money is positive.
func (m *Money) IsPositive() bool {
	return m.Amount_ > 0
}

// IsNegative returns boolean of whether the value of Money is negative.
func (m *Money) IsNegative() bool {
	return m.Amount_ < 0
}

// Absolute returns new Money struct from given Money using absolute monetary value.
func (m *Money) Absolute() *Money {
	return &Money{Amount_: mutate.calc.absolute(m.Amount_), Currency_: m.Currency_}
}

// Negative returns new Money struct from given Money using negative monetary value.
func (m *Money) Negative() *Money {
	return &Money{Amount_: mutate.calc.negative(m.Amount_), Currency_: m.Currency_}
}

// Add returns new Money struct with value representing sum of Self and Other Money.
func (m *Money) Add(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return m, nil
	}

	k := New(0, m.Currency_.Code)

	for _, m2 := range ms {
		if err := m.assertSameCurrency(m2); err != nil {
			return nil, err
		}

		k.Amount_ = mutate.calc.add(k.Amount_, m2.Amount_)
	}

	return &Money{Amount_: mutate.calc.add(m.Amount_, k.Amount_), Currency_: m.Currency_}, nil
}

// Subtract returns new Money struct with value representing difference of Self and Other Money.
func (m *Money) Subtract(ms ...*Money) (*Money, error) {
	if len(ms) == 0 {
		return m, nil
	}

	k := New(0, m.Currency_.Code)

	for _, m2 := range ms {
		if err := m.assertSameCurrency(m2); err != nil {
			return nil, err
		}

		k.Amount_ = mutate.calc.add(k.Amount_, m2.Amount_)
	}

	return &Money{Amount_: mutate.calc.subtract(m.Amount_, k.Amount_), Currency_: m.Currency_}, nil
}

// Multiply returns new Money struct with value representing Self multiplied value by multiplier.
func (m *Money) Multiply(muls ...int64) *Money {
	if len(muls) == 0 {
		panic("At least one multiplier is required to multiply")
	}

	k := New(1, m.Currency_.Code)

	for _, m2 := range muls {
		k.Amount_ = mutate.calc.multiply(k.Amount_, m2)
	}

	return &Money{Amount_: mutate.calc.multiply(m.Amount_, k.Amount_), Currency_: m.Currency_}
}

// Round returns new Money struct with value rounded to nearest zero.
func (m *Money) Round() *Money {
	return &Money{Amount_: mutate.calc.round(m.Amount_, m.Currency_.Fraction), Currency_: m.Currency_}
}

// Split returns slice of Money structs with split Self value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first will likely receive more pennies than ones that are listed later.
func (m *Money) Split(n int) ([]*Money, error) {
	if n <= 0 {
		return nil, errors.New("split must be higher than zero")
	}

	a := mutate.calc.divide(m.Amount_, int64(n))
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{Amount_: a, Currency_: m.Currency_}
	}

	r := mutate.calc.modulus(m.Amount_, int64(n))
	l := mutate.calc.absolute(r)
	// Add leftovers to the first parties.

	v := int64(1)
	if m.Amount_ < 0 {
		v = -1
	}
	for p := 0; l != 0; p++ {
		ms[p].Amount_ = mutate.calc.add(ms[p].Amount_, v)
		l--
	}

	return ms, nil
}

// Allocate returns slice of Money structs with split Self value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
func (m *Money) Allocate(rs ...int) ([]*Money, error) {
	if len(rs) == 0 {
		return nil, errors.New("no ratios specified")
	}

	// Calculate sum of ratios.
	var sum int64
	for _, r := range rs {
		if r < 0 {
			return nil, errors.New("negative ratios not allowed")
		}
		if int64(r) > (math.MaxInt64 - sum) {
			return nil, errors.New("sum of given ratios exceeds max int")
		}
		sum += int64(r)
	}

	var total int64
	ms := make([]*Money, 0, len(rs))
	for _, r := range rs {
		party := &Money{
			Amount_:   mutate.calc.allocate(m.Amount_, int64(r), sum),
			Currency_: m.Currency_,
		}

		ms = append(ms, party)
		total += party.Amount_
	}

	// if the sum of all ratios is zero, then we just returns zeros and don't do anything
	// with the leftover
	if sum == 0 {
		return ms, nil
	}

	// Calculate leftover value and divide to first parties.
	lo := m.Amount_ - total
	sub := int64(1)
	if lo < 0 {
		sub = -sub
	}

	for p := 0; lo != 0; p++ {
		ms[p].Amount_ = mutate.calc.add(ms[p].Amount_, sub)
		lo -= sub
	}

	return ms, nil
}

// Display lets represent Money struct as string in given Currency value.
func (m *Money) Display() string {
	c := m.Currency_.get()
	return c.Formatter().Format(m.Amount_)
}

// AsMajorUnits lets represent Money struct as subunits (float64) in given Currency value
func (m *Money) AsMajorUnits() float64 {
	c := m.Currency_.get()
	return c.Formatter().ToMajorUnits(m.Amount_)
}

// UnmarshalJSON is implementation of json.Unmarshaller
func (m *Money) UnmarshalJSON(b []byte) error {
	return UnmarshalJSON(m, b)
}

// MarshalJSON is implementation of json.Marshaller
func (m Money) MarshalJSON() ([]byte, error) {
	return MarshalJSON(m)
}

// Compare function compares two money of the same type
//
//	if m.amount > om.amount returns (1, nil)
//	if m.amount == om.amount returns (0, nil
//	if m.amount < om.amount returns (-1, nil)
//
// If compare moneys from distinct currency, return (m.amount, ErrCurrencyMismatch)
func (m *Money) Compare(om *Money) (int, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return int(m.Amount_), err
	}

	return m.compare(om), nil
}
