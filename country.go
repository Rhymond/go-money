package money

import (
	"fmt"
	"strings"
)

// CountryToCurrency maps ISO 3166-1 alpha-2 country codes to their primary
// ISO 4217 currency code. The map covers ~180 countries and territories.
//
// Use CurrencyForCountry or NewForCountry for the idiomatic API.
var CountryToCurrency = map[string]string{
	// Americas
	"AG": "XCD", "AI": "XCD", "AN": "ANG", "AR": "ARS", "AW": "AWG",
	"BB": "BBD", "BL": "EUR", "BM": "BMD", "BO": "BOB", "BR": "BRL",
	"BS": "BSD", "BZ": "BZD", "CA": "CAD", "CL": "CLP", "CO": "COP",
	"CR": "CRC", "CU": "CUP", "DM": "XCD", "DO": "DOP", "EC": "USD",
	"GD": "XCD", "GP": "EUR", "GT": "GTQ", "GY": "GYD", "HN": "HNL",
	"HT": "HTG", "JM": "JMD", "KN": "XCD", "KY": "KYD", "LC": "XCD",
	"MF": "EUR", "MQ": "EUR", "MS": "XCD", "MX": "MXN", "NI": "NIO",
	"PA": "PAB", "PE": "PEN", "PM": "EUR", "PR": "USD", "PY": "PYG",
	"SR": "SRD", "SV": "USD", "TC": "USD", "TT": "TTD", "US": "USD",
	"UY": "UYU", "VC": "XCD", "VE": "VES", "VG": "USD", "VI": "USD",

	// Europe
	"AD": "EUR", "AL": "ALL", "AT": "EUR", "AX": "EUR", "AZ": "AZN",
	"BA": "BAM", "BE": "EUR", "BG": "BGN", "BY": "BYN", "CH": "CHF",
	"CY": "EUR", "CZ": "CZK", "DE": "EUR", "DK": "DKK", "EE": "EUR",
	"ES": "EUR", "FI": "EUR", "FO": "DKK", "FR": "EUR", "GB": "GBP",
	"GE": "GEL", "GG": "GBP", "GI": "GIP", "GL": "DKK", "GR": "EUR",
	"HR": "EUR", "HU": "HUF", "IE": "EUR", "IM": "GBP", "IS": "ISK",
	"IT": "EUR", "JE": "GBP", "KZ": "KZT", "LI": "CHF", "LT": "EUR",
	"LU": "EUR", "LV": "EUR", "MC": "EUR", "MD": "MDL", "ME": "EUR",
	"MK": "MKD", "MT": "EUR", "NL": "EUR", "NO": "NOK", "PL": "PLN",
	"PT": "EUR", "RO": "RON", "RS": "RSD", "RU": "RUB", "SE": "SEK",
	"SI": "EUR", "SJ": "NOK", "SK": "EUR", "SM": "EUR", "TR": "TRY",
	"UA": "UAH", "UZ": "UZS", "VA": "EUR", "XK": "EUR",

	// Asia-Pacific
	"AF": "AFN", "AM": "AMD", "AU": "AUD", "BD": "BDT", "BN": "BND",
	"BT": "BTN", "CC": "AUD", "CN": "CNY", "CX": "AUD", "FJ": "FJD",
	"FM": "USD", "GU": "USD", "HK": "HKD", "ID": "IDR", "IN": "INR",
	"JP": "JPY", "KG": "KGS", "KH": "KHR", "KI": "AUD", "KP": "KPW",
	"KR": "KRW", "LA": "LAK", "LK": "LKR", "MH": "USD", "MM": "MMK",
	"MN": "MNT", "MO": "MOP", "MP": "USD", "MV": "MVR", "MY": "MYR",
	"NC": "XPF", "NF": "AUD", "NP": "NPR", "NR": "AUD", "NU": "NZD",
	"NZ": "NZD", "PF": "XPF", "PG": "PGK", "PH": "PHP", "PK": "PKR",
	"PW": "USD", "SB": "SBD", "SG": "SGD", "TH": "THB", "TJ": "TJS",
	"TL": "USD", "TM": "TMT", "TO": "TOP", "TV": "AUD", "TW": "TWD",
	"VN": "VND", "VU": "VUV", "WF": "XPF", "WS": "WST",

	// Middle East
	"AE": "AED", "BH": "BHD", "IL": "ILS", "IQ": "IQD", "IR": "IRR",
	"JO": "JOD", "KW": "KWD", "LB": "LBP", "OM": "OMR", "PS": "ILS",
	"QA": "QAR", "SA": "SAR", "SY": "SYP", "YE": "YER",

	// Africa
	"AO": "AOA", "BF": "XOF", "BI": "BIF", "BJ": "XOF", "BW": "BWP",
	"CD": "CDF", "CF": "XAF", "CG": "XAF", "CI": "XOF", "CM": "XAF",
	"CV": "CVE", "DJ": "DJF", "DZ": "DZD", "EG": "EGP", "ER": "ERN",
	"ET": "ETB", "GA": "XAF", "GH": "GHS", "GM": "GMD", "GN": "GNF",
	"GQ": "XAF", "GW": "XOF", "KE": "KES", "KM": "KMF", "LR": "LRD",
	"LS": "LSL", "LY": "LYD", "MA": "MAD", "MG": "MGA", "ML": "XOF",
	"MR": "MRU", "MU": "MUR", "MW": "MWK", "MZ": "MZN", "NA": "NAD",
	"NE": "XOF", "NG": "NGN", "RE": "EUR", "RW": "RWF", "SC": "SCR",
	"SD": "SDG", "SL": "SLL", "SN": "XOF", "SO": "SOS", "SS": "SSP",
	"ST": "STN", "SZ": "SZL", "TD": "XAF", "TG": "XOF", "TN": "TND",
	"TZ": "TZS", "UG": "UGX", "YT": "EUR", "ZA": "ZAR", "ZM": "ZMW",
	"ZW": "ZWL",
}

// CurrencyForCountry returns the primary ISO 4217 currency code for the given
// ISO 3166-1 alpha-2 country code (case-insensitive).
// The second return value reports whether the country was found.
//
//	code, ok := money.CurrencyForCountry("US") // "USD", true
//	code, ok := money.CurrencyForCountry("GB") // "GBP", true
//	code, ok := money.CurrencyForCountry("XX") // "",    false
func CurrencyForCountry(countryCode string) (string, bool) {
	c, ok := CountryToCurrency[strings.ToUpper(countryCode)]
	return c, ok
}

// NewForCountry creates a Money value for the primary currency of the given
// ISO 3166-1 alpha-2 country code (case-insensitive).
// Returns an error if the country code is not in CountryToCurrency.
//
//	m, _ := money.NewForCountry(1000, "US") // $10.00 USD
//	m, _ := money.NewForCountry(500,  "JP") // ¥500  JPY
func NewForCountry(amount int64, countryCode string) (*Money, error) {
	code, ok := CurrencyForCountry(countryCode)
	if !ok {
		return nil, fmt.Errorf("no currency mapping for country code %q", strings.ToUpper(countryCode))
	}
	return New(amount, code), nil
}
