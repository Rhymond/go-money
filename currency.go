package money

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type currency struct {
	Code     string
	Fraction int
	Grapheme string
	Template string
}

type currencies map[string]*currency

func (c *currency) Get(code string) *currency {
	code = strings.ToUpper(code)
	cs := c.read("./resources/currencies.json")

	if _, ok := cs[code]; !ok {
		log.Fatalf("currency %s not found", code)
	}

	cs[code].Code = code

	return cs[code]
}

func (c *currency) Equals(oc *currency) bool {
	return c.Code == oc.Code
}

func (c *currency) read(p string) currencies {
	file, err := ioutil.ReadFile(p)

	if err != nil {
		log.Fatalf("Can't read currencies from file %s: %v", p, err)
	}

	currs := make(currencies, 0)
	json.Unmarshal(file, &currs)

	return currs
}
