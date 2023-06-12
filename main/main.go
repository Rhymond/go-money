package main

import (
	"fmt"

	"github.com/Rhymond/go-money"
)

func main() {
	jeffBezosFortune := money.NewFromString("114,000,000,000.99", money.USD)
	fmt.Println(jeffBezosFortune.Display())
}
