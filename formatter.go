package money

import (
	"math"
	"strconv"
	"strings"
)

// Formatter stores Money formatting information.
type Formatter struct {
	Fraction int
	Decimal  string
	Thousand string
	Grapheme string
	Template string
	Code     string
}

// NewFormatter creates new Formatter instance.
func NewFormatter(fraction int, decimal, thousand, grapheme, template, code string) *Formatter {
	return &Formatter{
		Fraction: fraction,
		Decimal:  decimal,
		Thousand: thousand,
		Grapheme: grapheme,
		Template: template,
		Code:     code,
	}
}

// Format returns string of formatted integer using given currency template.
func (f *Formatter) Format(amount int64) string {
	// Work with absolute amount value
	sa := strconv.FormatInt(f.abs(amount), 10)

	if len(sa) <= f.Fraction {
		sa = strings.Repeat("0", f.Fraction-len(sa)+1) + sa
	}

	if f.Thousand != "" {
		for i := len(sa) - f.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + f.Thousand + sa[i:]
		}
	}

	if f.Fraction > 0 {
		sa = sa[:len(sa)-f.Fraction] + f.Decimal + sa[len(sa)-f.Fraction:]
	}
	sa = strings.Replace(f.Template, "1", sa, 1)
	sa = strings.Replace(sa, "$", f.Grapheme, 1)

	// Add minus sign for negative amount.
	if amount < 0 {
		sa = "-" + sa
	}

	return sa
}

// Format returns string of formatted integer using given currency template, with currency code instead of symbol.
func (f *Formatter) FormatWithCurrencyCode(amount int64) string {
	// Work with absolute amount value
	sa := strconv.FormatInt(f.abs(amount), 10)

	if len(sa) <= f.Fraction {
		sa = strings.Repeat("0", f.Fraction-len(sa)+1) + sa
	}

	if f.Thousand != "" {
		for i := len(sa) - f.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + f.Thousand + sa[i:]
		}
	}

	if f.Fraction > 0 {
		sa = sa[:len(sa)-f.Fraction] + f.Decimal + sa[len(sa)-f.Fraction:]
	}
	sa = strings.Replace(f.Template, "1", sa, 1)
	sa = strings.Replace(sa, "$", f.Code, 1)

	// Add minus sign for negative amount.
	if amount < 0 {
		sa = "-" + sa
	}

	return sa
}

// ToMajorUnits returns float64 representing the value in sub units using the currency data
func (f *Formatter) ToMajorUnits(amount int64) float64 {
	if f.Fraction == 0 {
		return float64(amount)
	}

	return float64(amount) / float64(math.Pow10(f.Fraction))
}

// abs return absolute value of given integer.
func (f Formatter) abs(amount int64) int64 {
	if amount < 0 {
		return -amount
	}

	return amount
}

func parseFormattedString(s string, currency *Currency) (int64, error) {
	// If the numeric string is empty, assume it's zero
	if len(s) == 0 {
		return 0, nil
	}

	// Remove currency code if in string
	s = strings.Replace(s, currency.Code, "", -1)

	// If the first character is a minus sign, we know this is a negative number
	negative := false
	if string(s[0]) == "-" {
		negative = true
		// Remove the minus symbol from the string
		s = s[1:]
	}

	// Remove all spaces
	s = strings.Replace(s, " ", "", -1)

	// Set the major amount to the full string
	majorAmountStr := s

	// Set the decimal amount to an empty int64
	var decimalAmount int64
	// If the string contains a decimal separator, parse
	if strings.Contains(s, currency.Decimal) {
		// Split the string by the decimal separator
		splitByDecimal := strings.Split(s, currency.Decimal)
		decimalAmountStr := splitByDecimal[1]

		// Ensure there's enough backing digits for the currency (e.g. 0.1 -> 0.10)
		for len(decimalAmountStr) < currency.Fraction {
			decimalAmountStr += "0"
		}

		// Parse the decimal amount to an integer
		decimalInt64, err := strconv.ParseInt(decimalAmountStr, 10, 64)
		if err != nil {
			return 0, err
		}

		// Set the decimal amount to the parsed amount, and set the major amount string to the
		// string without the decimal amount
		decimalAmount = decimalInt64
		majorAmountStr = splitByDecimal[0]
	}

	// Remove the thousands separator as we don't need this to determine the major amount
	majorAmountStr = strings.Replace(majorAmountStr, currency.Thousand, "", -1)

	// For the given currency, pad the major amount with the appropriate number of backing zeros
	for i := 0; i < currency.Fraction; i++ {
		majorAmountStr = majorAmountStr + "0"
	}

	// Convert major amount to int
	majorAmount, err := strconv.ParseInt(majorAmountStr, 10, 64)
	if err != nil {
		return 0, err
	}

	// Sum major and decimal amount, invert if negative
	sumAmount := majorAmount + decimalAmount
	if negative {
		sumAmount = sumAmount * -1
	}

	return sumAmount, nil
}
