package concrete

import (
	"fmt"
	"reflect"
)

// toAmount takes an interface representing an integer format (not float) and
// converts it to the internal amount type.
func toAmount(amount interface{}) AmountInt64 {
	switch amt := amount.(type) {
	case uint:
		return AmountInt64{value: int64(amt)}
	case uint8:
		return AmountInt64{value: int64(amt)}
	case uint16:
		return AmountInt64{value: int64(amt)}
	case uint32:
		return AmountInt64{value: int64(amt)}
	case uint64:
		return AmountInt64{value: int64(amt)}
	case int:
		return AmountInt64{value: int64(amt)}
	case int8:
		return AmountInt64{value: int64(amt)}
	case int16:
		return AmountInt64{value: int64(amt)}
	case int32:
		return AmountInt64{value: int64(amt)}
	case int64:
		return AmountInt64{value: amt}
	}
	panic(fmt.Sprintf("Unable to convert to Amount. Unsupported type: %s", reflect.TypeOf(amount).Name()))
}

// AmountInt64 is the type in which calculation results are stored.
//
// It defines methods for retrieving the internal value in different representations.
type AmountInt64 struct {
	value int64
}

// Uint returns a uint representation of the internal value of the Amount object
func (a AmountInt64) Uint() uint {
	return uint(a.value)
}

// Uint8 returns a uint8 representation of the internal value of the Amount object
func (a AmountInt64) Uint8() uint8 {
	return uint8(a.value)
}

// Uint16 returns a uint16 representation of the internal value of the Amount object
func (a AmountInt64) Uint16() uint16 {
	return uint16(a.value)
}

// Uint32 returns a uint32 representation of the internal value of the Amount object
func (a AmountInt64) Uint32() uint32 {
	return uint32(a.value)
}

// Uint64 returns a uint64 representation of the internal value of the Amount object
func (a AmountInt64) Uint64() uint64 {
	return uint64(a.value)
}

// Int returns a int representation of the internal value of the Amount object
func (a AmountInt64) Int() int {
	return int(a.value)
}

// Int8 returns a int8 representation of the internal value of the Amount object
func (a AmountInt64) Int8() int8 {
	return int8(a.value)
}

// Int16 returns a int16 representation of the internal value of the Amount object
func (a AmountInt64) Int16() int16 {
	return int16(a.value)
}

// Int32 returns a int32 representation of the internal value of the Amount object
func (a AmountInt64) Int32() int32 {
	return int32(a.value)
}

// Int64 returns a int64 representation of the internal value of the Amount object
func (a AmountInt64) Int64() int64 {
	return a.value
}

func (a AmountInt64) add(arg AmountInt64) AmountInt64 {
	return AmountInt64{a.value + arg.value}
}

func (a AmountInt64) subtract(arg AmountInt64) AmountInt64 {
	return AmountInt64{a.value - arg.value}
}

func (a AmountInt64) multiply(arg AmountInt64) AmountInt64 {
	return AmountInt64{a.value * arg.value}
}

func (a AmountInt64) divide(arg AmountInt64) AmountInt64 {
	return AmountInt64{a.value / arg.value}
}

func (a AmountInt64) modulus(arg AmountInt64) AmountInt64 {
	return AmountInt64{a.value % arg.value}
}

func (a AmountInt64) allocate(per, sum AmountInt64) AmountInt64 {
	return AmountInt64{(a.value * per.value) / sum.value}
}

func (a AmountInt64) absolute() AmountInt64 {
	if a.value < 0 {
		return AmountInt64{-a.value}
	}

	return AmountInt64{a.value}
}

func (a AmountInt64) negative() AmountInt64 {
	if a.value > 0 {
		return AmountInt64{-a.value}
	}

	return AmountInt64{a.value}
}

func (a AmountInt64) round() AmountInt64 {

	if a.value == 0 {
		return AmountInt64{0}
	}

	ab := a.absolute()
	m := ab.value % 100

	if m > 50 {
		ab.value += 100
	}

	ab.value = (ab.value / 100) * 100

	if a.value < 0 {
		a.value = -ab.value
	} else {
		a.value = ab.value
	}

	return AmountInt64{a.value}
}
