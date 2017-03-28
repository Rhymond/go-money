package money

import (
	"log"
	"strings"
)

type Currency struct {
	Code   string
	Name   string
	Format string
}

func (c *Currency) Get(code string) *Currency {
	code = strings.ToUpper(code)
	cs := new(Currencies).read("./data/currencies.json")

	if _, ok := cs[code]; !ok {
		log.Fatalf("Currency %s not found", code)
	}

	return cs[code]
}

func (c *Currency) Equals(oc *Currency) bool {
	return c.Code == oc.Code
}
