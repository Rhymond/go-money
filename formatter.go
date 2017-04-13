package money

import (
	"strconv"
	"strings"
)

// formatter stores Money formatting information
type formatter struct {
	Fraction int
	Decimal  string
	Thousand string
	Grapheme string
	Template string
}

// newFormatter creates new Formatter instance
func newFormatter(fraction int, decimal, thousand, grapheme, template string) *formatter {
	return &formatter{
		Fraction: fraction,
		Decimal:  decimal,
		Thousand: thousand,
		Grapheme: grapheme,
		Template: template,
	}
}

// format returns string of formatted integer using given currency template
func (f formatter) format(amount Amount) string {
	// Work with absolute amount value
	sa := strconv.FormatInt(amount.absolute().value, 10)
	if len(sa) <= f.Fraction {
		sa = strings.Repeat("0", f.Fraction-len(sa)+1) + sa
	}

	if f.Thousand != "" {
		for i := len(sa) - f.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + f.Thousand + sa[i:]
		}
	}

	sa = sa[:len(sa)-f.Fraction] + f.Decimal + sa[len(sa)-f.Fraction:]
	sa = strings.Replace(f.Template, "1", sa, 1)
	sa = strings.Replace(sa, "$", f.Grapheme, 1)

	// Add minus sign for negative amount
	if amount.value < 0 {
		sa = "-" + sa
	}

	return sa
}
