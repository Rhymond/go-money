package interfaces

import "github.com/Rhymond/go-money/currency"

// Amount defines an interface through which the internal representation of
// calculated values can be retrieved.
//
// Convenience methods allow for conversion to the types you require.
type Amount interface {
	Uint() uint
	Uint8() uint8
	Uint16() uint16
	Uint32() uint32
	Uint64() uint64
	Int() int
	Int8() int8
	Int16() int16
	Int32() int32
	Int64() int64
}

// Money is an interface defining methods through which a Money object
// can be manipulated. This allows for multiple implementations of the same
// patterns. For example, if an arbitrary precision implementation is needed,
// one could be created based on math/big
type Money interface {
	SameCurrency(om Money) bool
	Equals(om Money) (bool, error)
	GreaterThan(om Money) (bool, error)
	GreaterThanOrEqual(om Money) (bool, error)
	LessThan(om Money) (bool, error)
	LessThanOrEqual(om Money) (bool, error)
	IsZero() bool
	IsPositive() bool
	IsNegative() bool
	Absolute() Money
	Negative() Money
	Add(om Money) Money
	Subtract(om Money) Money
	Multiply(om Money) Money
	Divide(om Money) Money
	Round() Money
	Split(n int) ([]Money, error)
	Allocate(rs []int) ([]Money, error)
	Display() string
	Amount() Amount
	Currency() currency.Currency
}
