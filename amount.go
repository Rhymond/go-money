package money

// Amount is the type in which calculation results are stored.
//
// It defines methods for retrieving the internal value in different representations.
type Amount struct {
	value int64
}

// Uint returns a uint representation of the internal value of the Amount object
func (a Amount) Uint() uint {
	return uint(a.value)
}

// Uint8 returns a uint8 representation of the internal value of the Amount object
func (a Amount) Uint8() uint8 {
	return uint8(a.value)
}

// Uint16 returns a uint16 representation of the internal value of the Amount object
func (a Amount) Uint16() uint16 {
	return uint16(a.value)
}

// Uint32 returns a uint32 representation of the internal value of the Amount object
func (a Amount) Uint32() uint32 {
	return uint32(a.value)
}

// Uint64 returns a uint64 representation of the internal value of the Amount object
func (a Amount) Uint64() uint64 {
	return uint64(a.value)
}

// Int returns a int representation of the internal value of the Amount object
func (a Amount) Int() int {
	return int(a.value)
}

// Int8 returns a int8 representation of the internal value of the Amount object
func (a Amount) Int8() int8 {
	return int8(a.value)
}

// Int16 returns a int16 representation of the internal value of the Amount object
func (a Amount) Int16() int16 {
	return int16(a.value)
}

// Int32 returns a int32 representation of the internal value of the Amount object
func (a Amount) Int32() int32 {
	return int32(a.value)
}

// Int64 returns a int64 representation of the internal value of the Amount object
func (a Amount) Int64() int64 {
	return a.value
}

func (a Amount) add(arg Amount) Amount {
	return Amount{a.value + arg.value}
}

func (a Amount) subtract(arg Amount) Amount {
	return Amount{a.value - arg.value}
}

func (a Amount) multiply(arg Amount) Amount {
	return Amount{a.value * arg.value}
}

func (a Amount) divide(arg Amount) Amount {
	return Amount{a.value / arg.value}
}

func (a Amount) modulus(arg Amount) Amount {
	return Amount{a.value % arg.value}
}

func (a Amount) allocate(per, sum Amount) Amount {
	return Amount{(a.value * per.value) / sum.value}
}

func (a Amount) absolute() Amount {
	if a.value < 0 {
		return Amount{-a.value}
	}

	return Amount{a.value}
}

func (a Amount) negative() Amount {
	if a.value > 0 {
		return Amount{-a.value}
	}

	return Amount{a.value}
}

func (a Amount) round() Amount {

	if a.value == 0 {
		return Amount{0}
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

	return Amount{a.value}
}
