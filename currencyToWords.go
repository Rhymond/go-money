package money

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/im-adarsh/go-money/num2words"
)

// CurrencyMeta holds the English word names for a currency's main unit
// and its fractional sub-unit.
type CurrencyMeta struct {
	// MainWord is the singular English name for the whole unit (e.g. "dollar").
	MainWord string
	// SubWord is the singular English name for the fractional sub-unit (e.g. "cent").
	SubWord string
}

// CountryCurrencyMeta maps ISO 4217 currency codes to their English word
// representations used by ToWords / GetCurrencyAmountWords.
//
// New entries can be added at runtime via AddCurrencyMeta.
var CountryCurrencyMeta = map[string]CurrencyMeta{
	// Americas
	"USD": {MainWord: "dollar", SubWord: "cents"},
	"CAD": {MainWord: "dollar", SubWord: "cents"},
	"AUD": {MainWord: "dollar", SubWord: "cents"},
	"NZD": {MainWord: "dollar", SubWord: "cents"},
	"MXN": {MainWord: "peso", SubWord: "centavo"},
	"BRL": {MainWord: "real", SubWord: "centavo"},
	"ARS": {MainWord: "peso", SubWord: "centavo"},
	"CLP": {MainWord: "peso", SubWord: ""},
	"COP": {MainWord: "peso", SubWord: ""},
	"PEN": {MainWord: "sol", SubWord: "céntimo"},
	"BOB": {MainWord: "boliviano", SubWord: "centavo"},
	"UYU": {MainWord: "peso", SubWord: ""},
	"TTD": {MainWord: "dollar", SubWord: "cents"},
	"JMD": {MainWord: "dollar", SubWord: "cents"},
	"DOP": {MainWord: "peso", SubWord: "centavo"},
	"BSD": {MainWord: "dollar", SubWord: "cents"},
	"BZD": {MainWord: "dollar", SubWord: "cents"},
	"GYD": {MainWord: "dollar", SubWord: "cents"},
	"SRD": {MainWord: "dollar", SubWord: "cents"},
	"HTG": {MainWord: "gourde", SubWord: "centime"},
	"PAB": {MainWord: "balboa", SubWord: "centésimo"},

	// Europe
	"EUR": {MainWord: "euro", SubWord: "cent"},
	"GBP": {MainWord: "pound", SubWord: "penny"},
	"CHF": {MainWord: "franc", SubWord: "centime"},
	"SEK": {MainWord: "krona", SubWord: "öre"},
	"NOK": {MainWord: "krone", SubWord: "øre"},
	"DKK": {MainWord: "krone", SubWord: "øre"},
	"PLN": {MainWord: "złoty", SubWord: "grosz"},
	"CZK": {MainWord: "koruna", SubWord: "haléř"},
	"HUF": {MainWord: "forint", SubWord: ""},
	"RON": {MainWord: "leu", SubWord: "ban"},
	"BGN": {MainWord: "lev", SubWord: "stotinka"},
	"HRK": {MainWord: "kuna", SubWord: "lipa"},
	"RUB": {MainWord: "ruble", SubWord: "kopek"},
	"UAH": {MainWord: "hryvnia", SubWord: "kopiyka"},
	"TRY": {MainWord: "lira", SubWord: "kuruş"},
	"ISK": {MainWord: "króna", SubWord: "eyrir"},
	"BAM": {MainWord: "mark", SubWord: "fenig"},
	"RSD": {MainWord: "dinar", SubWord: "para"},
	"MKD": {MainWord: "denar", SubWord: "deni"},
	"ALL": {MainWord: "lek", SubWord: "qindarkë"},

	// Asia-Pacific
	"JPY": {MainWord: "yen", SubWord: ""},
	"CNY": {MainWord: "yuan", SubWord: "fen"},
	"INR": {MainWord: "rupee", SubWord: "paisa"},
	"HKD": {MainWord: "dollar", SubWord: "cents"},
	"SGD": {MainWord: "dollar", SubWord: "cents"},
	"KRW": {MainWord: "won", SubWord: ""},
	"TWD": {MainWord: "dollar", SubWord: ""},
	"THB": {MainWord: "baht", SubWord: "satang"},
	"MYR": {MainWord: "ringgit", SubWord: "sen"},
	"IDR": {MainWord: "rupiah", SubWord: "sen"},
	"PHP": {MainWord: "pesos", SubWord: "sentimos"},
	"VND": {MainWord: "dong", SubWord: ""},
	"PKR": {MainWord: "rupee", SubWord: "paisa"},
	"BDT": {MainWord: "taka", SubWord: "paisa"},
	"LKR": {MainWord: "rupee", SubWord: "cent"},
	"NPR": {MainWord: "rupee", SubWord: "paisa"},
	"MMK": {MainWord: "kyat", SubWord: "pya"},
	"KHR": {MainWord: "riel", SubWord: "sen"},
	"LAK": {MainWord: "kip", SubWord: "att"},
	"MNT": {MainWord: "tögrög", SubWord: "möngö"},
	"KZT": {MainWord: "tenge", SubWord: "tiyn"},
	"UZS": {MainWord: "soum", SubWord: "tiyin"},
	"AZN": {MainWord: "manat", SubWord: "qəpik"},
	"GEL": {MainWord: "lari", SubWord: "tetri"},

	// Middle East & Africa
	"AED": {MainWord: "dirham", SubWord: "fils"},
	"SAR": {MainWord: "riyal", SubWord: "halala"},
	"QAR": {MainWord: "riyal", SubWord: "dirham"},
	"KWD": {MainWord: "dinar", SubWord: "fils"},
	"BHD": {MainWord: "dinar", SubWord: "fils"},
	"OMR": {MainWord: "rial", SubWord: "baisa"},
	"JOD": {MainWord: "dinar", SubWord: "fils"},
	"IQD": {MainWord: "dinar", SubWord: "fils"},
	"IRR": {MainWord: "rial", SubWord: "dinar"},
	"ILS": {MainWord: "shekel", SubWord: "agora"},
	"EGP": {MainWord: "pound", SubWord: "piastre"},
	"ZAR": {MainWord: "rand", SubWord: "cent"},
	"NGN": {MainWord: "naira", SubWord: "kobo"},
	"KES": {MainWord: "shilling", SubWord: "cent"},
	"GHS": {MainWord: "cedi", SubWord: "pesewa"},
	"TZS": {MainWord: "shilling", SubWord: ""},
	"UGX": {MainWord: "shilling", SubWord: ""},
	"ETB": {MainWord: "birr", SubWord: "santim"},
	"MUR": {MainWord: "rupee", SubWord: "cent"},
	"ZMW": {MainWord: "kwacha", SubWord: "ngwee"},
	"MAD": {MainWord: "dirham", SubWord: "centime"},
	"TND": {MainWord: "dinar", SubWord: "millime"},
	"DZD": {MainWord: "dinar", SubWord: "centime"},
	"LYD": {MainWord: "dinar", SubWord: "dirham"},
}

// AddCurrencyMeta registers the English word names for a currency so that
// ToWords / GetCurrencyAmountWords can produce human-readable output for it.
//
//	money.AddCurrencyMeta("XYZ", "zorkmid", "zork")
func AddCurrencyMeta(code, mainWord, subWord string) {
	CountryCurrencyMeta[strings.ToUpper(code)] = CurrencyMeta{
		MainWord: mainWord,
		SubWord:  subWord,
	}
}

// GetCurrencyAmountWords converts a monetary amount to English words using
// the currency names registered in CountryCurrencyMeta.
//
// amount is expected to be in the currency's smallest unit (e.g. cents).
// currencyCode must be a key in CountryCurrencyMeta.
//
// If the currency code is not registered, the raw numeric string is returned.
//
//	GetCurrencyAmountWords(100, "USD")  // "one dollar only"
//	GetCurrencyAmountWords(150, "USD")  // "one dollar and fifty cents only"
//	GetCurrencyAmountWords(50, "USD")   // "fifty cents only"
func GetCurrencyAmountWords(amount float64, currencyCode string) string {
	amount = ConvertTo2DecimalPlaces(amount)
	strAmount := fmt.Sprintf("%+v", amount)
	str := strings.Split(strAmount, ".")

	if len(str) == 1 {
		str = append(str, "00")
	}
	mainAmt, err := strconv.Atoi(str[0])
	if err != nil {
		return fmt.Sprintf("%v", amount)
	}
	subAmt, err := strconv.Atoi(str[1])
	if err != nil {
		return fmt.Sprintf("%v", amount)
	}

	mainAmtStr := num2words.Convert(mainAmt)
	subAmtStr := num2words.Convert(subAmt)

	c, ok := CountryCurrencyMeta[currencyCode]
	if !ok {
		return fmt.Sprintf("%v", amount)
	}

	// Currency has no sub-unit (e.g. JPY, KRW)
	if c.SubWord == "" {
		if mainAmt != 0 {
			return fmt.Sprintf("%v %v only", mainAmtStr, c.MainWord)
		}
		return fmt.Sprintf("zero %v only", c.MainWord)
	}

	// Currency has no main unit name at all (should not happen with current data)
	if c.MainWord == "" {
		if mainAmt != 0 && subAmt != 0 {
			return fmt.Sprintf("%v and %v only", mainAmtStr, subAmtStr)
		} else if mainAmt != 0 {
			return fmt.Sprintf("%v only", mainAmtStr)
		} else if subAmt != 0 {
			return fmt.Sprintf("point %v only", subAmtStr)
		}
	}

	if mainAmt != 0 && subAmt != 0 {
		return fmt.Sprintf("%v %v and %v %v only", mainAmtStr, c.MainWord, subAmtStr, c.SubWord)
	} else if mainAmt != 0 {
		return fmt.Sprintf("%v %v only", mainAmtStr, c.MainWord)
	} else if subAmt != 0 {
		return fmt.Sprintf("%v %v only", subAmtStr, c.SubWord)
	}
	return strAmount
}

// ConvertTo2DecimalPlaces rounds a float64 value to exactly two decimal places.
func ConvertTo2DecimalPlaces(d float64) float64 {
	// strconv.ParseFloat never errors on output from fmt.Sprintf("%.2f", …).
	f, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", d), 64)
	return f
}
