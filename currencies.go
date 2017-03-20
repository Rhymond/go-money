package gocash

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Currencies map[string]*Currency

func (c *Currencies) read(path string) Currencies {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalf("Can't read currencies data, because: %v", err)
	}

	currencies := make(Currencies, 0)
	json.Unmarshal(file, &currencies)
	return currencies
}
