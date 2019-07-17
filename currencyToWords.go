package money

import (
	"fmt"
	"github.com/im-adarsh/go-money/num2words"
	"strconv"
	"strings"
)

type CurrencyMeta struct {
	MainWord string
	SubWord  string
}

var CountryCurrencyMeta = map[string]CurrencyMeta{
	"PHP": {MainWord: "pesos", SubWord: "sentimos"},
	"SGD": {MainWord: "dollar", SubWord: "cents"},
}

func GetCurrencyAmountWords(amount float64, currencyCode string) string {
	amount = ConvertTo2DecimalPlaces(amount)
	strAmount := fmt.Sprintf("%+v", amount)
	str := strings.Split(strAmount, ".")
	if len(str) > 2 {
		return fmt.Sprintf("%v", amount)
	}

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

func ConvertTo2DecimalPlaces(d float64) float64 {
	amountPaidStr := fmt.Sprintf("%.2f", d)
	amountPaid, err := strconv.ParseFloat(amountPaidStr, 64)
	if err != nil {
		return d
	}
	return amountPaid
}
