package money

import (
	"strconv"
	"strings"
	"fmt"
)

type Formatter struct {
	Fraction int
	Decimal  string
	Thousand string
	Grapheme string
	Template string
}

func NewFormatter(fraction int, decimal, thousand, grapheme, template string) *Formatter {
	return &Formatter{
		Fraction: fraction,
		Decimal:  decimal,
		Thousand: thousand,
		Grapheme: grapheme,
		Template: template,
	}
}

func (f *Formatter) Format(amount int) string {
	// Work with absolute amount value
	sa := strconv.Itoa(f.abs(amount))

	if len(sa) <= f.Fraction {
		sa = strings.Repeat("0", f.Fraction-len(sa)+1) + sa
	}

	if f.Thousand != "" {
		for i := len(sa) - f.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + f.Thousand + sa[i:]
		}
	}

	sa = sa[:len(sa) - f.Fraction] + f.Decimal + sa[len(sa) - f.Fraction:]
	sa = fmt.Sprintf(f.Template, sa)
	sa = strings.Replace(sa, "$", f.Grapheme, -1)

	// Add minus sign for negative amount
	if amount < 0 {
		sa = "-" + sa
	}

	return sa
}

func (f Formatter) abs(amount int) int {
	if amount < 0 {
		return -amount
	}

	return amount
}
