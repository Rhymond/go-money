package money

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

type currency struct {
	Code     string
	Fraction int
	Grapheme string
	Template string
}

var currencies = map[string]*currency{
	"AED": &currency{Fraction: 2, Grapheme: ".\u062f.\u0625", Template: "1 $"},
	"AFN": &currency{Fraction: 2, Grapheme: "\u060b", Template: "1 $"},
	"ALL": &currency{Fraction: 2, Grapheme: "L", Template: "$1"},
	"AMD": &currency{Fraction: 2, Grapheme: "\u0564\u0580.", Template: "1 $"},
	"ANG": &currency{Fraction: 2, Grapheme: "\u0192", Template: "$1"},
	"ARS": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"AUD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"AWG": &currency{Fraction: 2, Grapheme: "\u0192", Template: "$1"},
	"AZN": &currency{Fraction: 2, Grapheme: "\u20bc", Template: "$1"},
	"BAM": &currency{Fraction: 2, Grapheme: "KM", Template: "$1"},
	"BBD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"BGN": &currency{Fraction: 2, Grapheme: "\u043b\u0432", Template: "$1"},
	"BHD": &currency{Fraction: 3, Grapheme: ".\u062f.\u0628", Template: "1 $"},
	"BMD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"BND": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"BOB": &currency{Fraction: 2, Grapheme: "Bs.", Template: "$1"},
	"BRL": &currency{Fraction: 2, Grapheme: "R$", Template: "$1"},
	"BSD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"BWP": &currency{Fraction: 2, Grapheme: "P", Template: "$1"},
	"BYN": &currency{Fraction: 2, Grapheme: "p.", Template: "1 $"},
	"BYR": &currency{Fraction: 0, Grapheme: "p.", Template: "1 $"},
	"BZD": &currency{Fraction: 2, Grapheme: "BZ$", Template: "$1"},
	"CAD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"CLP": &currency{Fraction: 0, Grapheme: "$", Template: "$1"},
	"CNY": &currency{Fraction: 2, Grapheme: "\u5143", Template: "1 $"},
	"COP": &currency{Fraction: 0, Grapheme: "$", Template: "$1"},
	"CRC": &currency{Fraction: 2, Grapheme: "\u20a1", Template: "$1"},
	"CUP": &currency{Fraction: 2, Grapheme: "$MN", Template: "$1"},
	"CZK": &currency{Fraction: 2, Grapheme: "K\u010d", Template: "1 $"},
	"DKK": &currency{Fraction: 2, Grapheme: "kr", Template: "1 $"},
	"DOP": &currency{Fraction: 2, Grapheme: "RD$", Template: "$1"},
	"DZD": &currency{Fraction: 2, Grapheme: ".\u062f.\u062c", Template: "1 $"},
	"EEK": &currency{Fraction: 2, Grapheme: "kr", Template: "$1"},
	"EGP": &currency{Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"EUR": &currency{Fraction: 2, Grapheme: "\u20ac", Template: "$1"},
	"FJD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"FKP": &currency{Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"GBP": &currency{Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"GGP": &currency{Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"GHC": &currency{Fraction: 2, Grapheme: "\u00a2", Template: "$1"},
	"GIP": &currency{Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"GTQ": &currency{Fraction: 2, Grapheme: "Q", Template: "$1"},
	"GYD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"HKD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"HNL": &currency{Fraction: 2, Grapheme: "L", Template: "$1"},
	"HRK": &currency{Fraction: 2, Grapheme: "kn", Template: "$1"},
	"HUF": &currency{Fraction: 0, Grapheme: "Ft", Template: "$1"},
	"IDR": &currency{Fraction: 2, Grapheme: "Rp", Template: "$1"},
	"ILS": &currency{Fraction: 2, Grapheme: "\u20aa", Template: "$1"},
	"IMP": &currency{Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"INR": &currency{Fraction: 2, Grapheme: "\u20b9", Template: "$1"},
	"IQD": &currency{Fraction: 3, Grapheme: ".\u062f.\u0639", Template: "1 $"},
	"IRR": &currency{Fraction: 2, Grapheme: "\ufdfc", Template: "1 $"},
	"ISK": &currency{Fraction: 2, Grapheme: "kr", Template: "$1"},
	"JEP": &currency{Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"JMD": &currency{Fraction: 2, Grapheme: "J$", Template: "$1"},
	"JOD": &currency{Fraction: 3, Grapheme: ".\u062f.\u0625", Template: "1 $"},
	"JPY": &currency{Fraction: 0, Grapheme: "\u00a5", Template: "$1"},
	"KES": &currency{Fraction: 2, Grapheme: "KSh", Template: "$1"},
	"KGS": &currency{Fraction: 2, Grapheme: "\u0441\u043e\u043c", Template: "$1"},
	"KHR": &currency{Fraction: 2, Grapheme: "\u17db", Template: "$1"},
	"KPW": &currency{Fraction: 0, Grapheme: "\u20a9", Template: "$1"},
	"KRW": &currency{Fraction: 0, Grapheme: "\u20a9", Template: "$1"},
	"KWD": &currency{Fraction: 3, Grapheme: ".\u062f.\u0643", Template: "1 $"},
	"KYD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"KZT": &currency{Fraction: 2, Grapheme: "\u20b8", Template: "$1"},
	"LAK": &currency{Fraction: 2, Grapheme: "\u20ad", Template: "$1"},
	"LBP": &currency{Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"LKR": &currency{Fraction: 2, Grapheme: "\u20a8", Template: "$1"},
	"LRD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"LTL": &currency{Fraction: 2, Grapheme: "Lt", Template: "$1"},
	"LVL": &currency{Fraction: 2, Grapheme: "Ls", Template: "1 $"},
	"LYD": &currency{Fraction: 3, Grapheme: ".\u062f.\u0644", Template: "1 $"},
	"MAD": &currency{Fraction: 2, Grapheme: ".\u062f.\u0645", Template: "1 $"},
	"MKD": &currency{Fraction: 2, Grapheme: "\u0434\u0435\u043d", Template: "$1"},
	"MNT": &currency{Fraction: 2, Grapheme: "\u20ae", Template: "$1"},
	"MUR": &currency{Fraction: 2, Grapheme: "\u20a8", Template: "$1"},
	"MXN": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"MYR": &currency{Fraction: 2, Grapheme: "RM", Template: "$1"},
	"MZN": &currency{Fraction: 2, Grapheme: "MT", Template: "$1"},
	"NAD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"NGN": &currency{Fraction: 2, Grapheme: "\u20a6", Template: "$1"},
	"NIO": &currency{Fraction: 2, Grapheme: "C$", Template: "$1"},
	"NOK": &currency{Fraction: 2, Grapheme: "kr", Template: "1 $"},
	"NPR": &currency{Fraction: 2, Grapheme: "\u20a8", Template: "$1"},
	"NZD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"OMR": &currency{Fraction: 3, Grapheme: "\ufdfc", Template: "1 $"},
	"PAB": &currency{Fraction: 2, Grapheme: "B\/.", Template: "$1"},
	"PEN": &currency{Fraction: 2, Grapheme: "S\/", Template: "$1"},
	"PHP": &currency{Fraction: 2, Grapheme: "\u20b1", Template: "$1"},
	"PKR": &currency{Fraction: 2, Grapheme: "\u20a8", Template: "$1"},
	"PLN": &currency{Fraction: 2, Grapheme: "z\u0142", Template: "1 $"},
	"PYG": &currency{Fraction: 0, Grapheme: "Gs", Template: "1$"},
	"QAR": &currency{Fraction: 2, Grapheme: "\ufdfc", Template: "1 $"},
	"RON": &currency{Fraction: 2, Grapheme: "lei", Template: "$1"},
	"RSD": &currency{Fraction: 2, Grapheme: "\u0414\u0438\u043d.", Template: "$1"},
	"RUB": &currency{Fraction: 2, Grapheme: "\u20bd", Template: "1 $"},
	"RUR": &currency{Fraction: 2, Grapheme: "\u20bd", Template: "1 $"},
	"SAR": &currency{Fraction: 2, Grapheme: "\ufdfc", Template: "1 $"},
	"SBD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"SCR": &currency{Fraction: 2, Grapheme: "\u20a8", Template: "$1"},
	"SEK": &currency{Fraction: 2, Grapheme: "kr", Template: "1 $"},
	"SGD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"SHP": &currency{Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"SOS": &currency{Fraction: 2, Grapheme: "S", Template: "$1"},
	"SRD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"SVC": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"SYP": &currency{Fraction: 2, Grapheme: "\u00a3", Template: "$1"},
	"THB": &currency{Fraction: 2, Grapheme: "\u0e3f", Template: "$1"},
	"TND": &currency{Fraction: 3, Grapheme: ".\u062f.\u062a", Template: "1 $"},
	"TRL": &currency{Fraction: 2, Grapheme: "\u20a4", Template: "$1"},
	"TRY": &currency{Fraction: 2, Grapheme: "\u20ba", Template: "$1"},
	"TTD": &currency{Fraction: 2, Grapheme: "TT$", Template: "$1"},
	"TWD": &currency{Fraction: 0, Grapheme: "NT$", Template: "$1"},
	"TZS": &currency{Fraction: 0, Grapheme: "TSh", Template: "$1"},
	"UAH": &currency{Fraction: 2, Grapheme: "\u20b4", Template: "$1"},
	"UGX": &currency{Fraction: 0, Grapheme: "USh", Template: "$1"},
	"USD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"UYU": &currency{Fraction: 0, Grapheme: "$U", Template: "$1"},
	"UZS": &currency{Fraction: 2, Grapheme: "so\u2019m", Template: "$1"},
	"VEF": &currency{Fraction: 2, Grapheme: "Bs", Template: "$1"},
	"VND": &currency{Fraction: 0, Grapheme: "\u20ab", Template: "1 $"},
	"XCD": &currency{Fraction: 2, Grapheme: "$", Template: "$1"},
	"YER": &currency{Fraction: 2, Grapheme: "\ufdfc", Template: "1 $"},
	"ZAR": &currency{Fraction: 2, Grapheme: "R", Template: "$1"},
	"ZWD": &currency{Fraction: 2, Grapheme: "Z$", Template: "$1"},
}

func newCurrency(code string) *currency {
	return &currency{Code: strings.ToUpper(code)}
}

func (c *currency) get() (*currency, error) {
	if v, ok := currencies[c.Code]; ok {
		return v, nil
	}
	
	return nil, errors.New("Currency not found")
}

func (c *currency) equals(oc *currency) bool {
	return c.Code == oc.Code
}
