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

type currencies map[string]*currency

func newCurrency(code string) *currency {
	return &currency{Code: strings.ToUpper(code)}
}

func (c *currency) get() (*currency, error) {
	cs, err := c.read("./resources/currencies.json")

	if err != nil {
		return nil, err
	}

	if _, ok := cs[c.Code]; !ok {
		return nil, errors.New("Currency not found")
	}

	cs[c.Code].Code = c.Code
	return cs[c.Code], nil
}

func (c *currency) equals(oc *currency) bool {
	return c.Code == oc.Code
}

func (c *currency) read(p string) (currencies, error) {
	file, err := ioutil.ReadFile(p)

	if err != nil {
		return nil, errors.New("Can't read currencies resource")
	}

	currs := make(currencies, 0)
	json.Unmarshal(file, &currs)

	return currs, nil
}
