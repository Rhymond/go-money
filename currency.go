package money

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type Currency struct {
	Code     string
	Fraction int
	Grapheme string
	Template string
}

type Currencies map[string]*Currency

func (c *Currency) Get(code string) *Currency {
	code = strings.ToUpper(code)
	cs := c.read("./resources/currencies.json")

	if _, ok := cs[code]; !ok {
		log.Fatalf("Currency %s not found", code)
	}

	cs[code].Code = code

	return cs[code]
}

func (c *Currency) Equals(oc *Currency) bool {
	return c.Code == oc.Code
}

func (c *Currency) read(p string) Currencies {
	file, err := ioutil.ReadFile(p)

	if err != nil {
		log.Fatalf("Can't read currencies from file %s: %v", p, err)
	}

	currencies := make(Currencies, 0)
	json.Unmarshal(file, &currencies)

	return currencies
}
