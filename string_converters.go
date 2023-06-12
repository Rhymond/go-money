package money

import (
	"errors"
	"regexp"
	"strconv"
)

func patternFactory(fractionNumber int) (*regexp.Regexp, error) {
	switch fractionNumber {
	case 0:
		zeroDigitsPattern := `^\d{1,3}(,\d{3})*$`
		return regexp.MustCompile(zeroDigitsPattern), nil
	case 2:
		twoDigitsPattern := `^\d{1,3}(,\d{3})*(\.\d{2})?$`
		return regexp.MustCompile(twoDigitsPattern), nil
	case 3:
		threeDigitsPattern := `^\d{1,3}(,\d{3})*(\.\d{3})?$`
		return regexp.MustCompile(threeDigitsPattern), nil
	}
	return nil, errors.New("invalid currency fraction")
}

func ConvertStringAmount(amount string, currency *Currency) (*Money, error) {

	ok, err := amountIsValid(amount, currency.Fraction)

	if err != nil || !ok {
		return nil, errors.New("invalid amount pattern")
	}

	amountFormated := cleanString(amount)

	amountAsFloat, err := parseStringToFloat64(amountFormated)

	if err != nil {
		return nil, errors.New("error parsing amount")
	}

	return NewFromFloat(amountAsFloat, currency.Code), nil
}

func amountIsValid(amount string, fraction int) (bool, error) {
	regex, err := patternFactory(fraction)

	if err != nil {
		return false, err
	}

	return regex.MatchString(amount), nil
}

func cleanString(validAmount string) string {
	cleanPattern := "[,]"

	regex := regexp.MustCompile(cleanPattern)

	amountFormated := regex.ReplaceAllString(validAmount, "")

	return amountFormated
}

func parseStringToFloat64(validAmountFormated string) (float64, error) {
	amountAsFloat, err := strconv.ParseFloat(validAmountFormated, 64)
	return amountAsFloat, err
}
