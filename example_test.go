package money_test

import (
	"fmt"

	money "github.com/im-adarsh/go-money"
)

func ExampleNew() {
	m := money.New(100, "USD")
	fmt.Println(m.Display())
	// Output: $1.00
}

func ExampleNewFromFloat() {
	m := money.NewFromFloat(10.99, "USD")
	fmt.Println(m.Display())
	// Output: $10.99
}

func ExampleMoney_Add() {
	a := money.New(100, "USD")
	b := money.New(200, "USD")
	result, _ := a.Add(b)
	fmt.Println(result.Display())
	// Output: $3.00
}

func ExampleMoney_Subtract() {
	a := money.New(300, "USD")
	b := money.New(100, "USD")
	result, _ := a.Subtract(b)
	fmt.Println(result.Display())
	// Output: $2.00
}

func ExampleMoney_Multiply() {
	m := money.New(100, "USD")
	fmt.Println(m.Multiply(3).Display())
	// Output: $3.00
}

func ExampleMoney_MultiplyFloat() {
	price := money.New(1000, "USD") // $10.00
	withVAT := price.MultiplyFloat(1.21)
	fmt.Println(withVAT.Display())
	// Output: $12.10
}

func ExampleMoney_Divide() {
	m := money.New(300, "USD")
	fmt.Println(m.Divide(3).Display())
	// Output: $1.00
}

func ExampleMoney_DivideWithRounding() {
	m := money.New(11, "USD") // 11 cents

	up := m.DivideWithRounding(2, money.RoundHalfUp)
	down := m.DivideWithRounding(2, money.RoundHalfDown)

	fmt.Println(up.Amount(), down.Amount())
	// Output: 6 5
}

func ExampleMoney_Percentage() {
	price := money.New(10000, "USD") // $100.00
	tax := price.Percentage(8.5)
	fmt.Println(tax.Display())
	// Output: $8.50
}

func ExampleMoney_Split() {
	m := money.New(100, "GBP") // £1.00
	parts, _ := m.Split(3)
	for _, p := range parts {
		fmt.Println(p.Display())
	}
	// Output:
	// £0.34
	// £0.33
	// £0.33
}

func ExampleMoney_Allocate() {
	m := money.New(100, "USD") // $1.00
	parts, _ := m.Allocate(50, 30, 20)
	for _, p := range parts {
		fmt.Println(p.Display())
	}
	// Output:
	// $0.50
	// $0.30
	// $0.20
}

func ExampleMoney_RoundWithMode() {
	m := money.New(150, "USD") // $1.50 — tie

	fmt.Println(m.RoundWithMode(money.RoundHalfUp).Display())
	fmt.Println(m.RoundWithMode(money.RoundHalfDown).Display())
	fmt.Println(m.RoundWithMode(money.RoundHalfEven).Display())
	// Output:
	// $2.00
	// $1.00
	// $2.00
}

func ExampleMoney_IsWhole() {
	fmt.Println(money.New(100, "USD").IsWhole()) // $1.00
	fmt.Println(money.New(150, "USD").IsWhole()) // $1.50
	// Output:
	// true
	// false
}

func ExampleMoney_Clamp() {
	min := money.New(100, "USD")
	max := money.New(500, "USD")

	below, _ := money.New(50, "USD").Clamp(min, max)
	within, _ := money.New(300, "USD").Clamp(min, max)
	above, _ := money.New(600, "USD").Clamp(min, max)

	fmt.Println(below.Display(), within.Display(), above.Display())
	// Output: $1.00 $3.00 $5.00
}

func ExampleMoney_AsParts() {
	m := money.New(1234, "USD") // $12.34
	whole, frac := m.AsParts()
	fmt.Printf("%d dollars and %02d cents\n", whole, frac)
	// Output: 12 dollars and 34 cents
}

func ExampleMoney_AsFloat64() {
	m := money.New(1099, "USD")
	fmt.Printf("%.2f\n", m.AsFloat64())
	// Output: 10.99
}

func ExampleMoney_WithAmount() {
	base := money.New(0, "EUR")
	fee := base.WithAmount(250) // €2.50
	fmt.Println(fee.Display())
	// Output: €2.50
}

func ExampleNewExchangeRate() {
	rate, _ := money.NewExchangeRate("USD", "EUR", 0.92)
	euros, _ := rate.Convert(money.New(1000, "USD"))
	fmt.Println(euros.Display())
	// Output: €9.20
}

func ExampleExchangeRate_Invert() {
	rate, _ := money.NewExchangeRate("USD", "EUR", 1.0)
	inv := rate.Invert()
	fmt.Println(inv.From(), inv.To())
	// Output: EUR USD
}

func ExampleSum() {
	result, _ := money.Sum(
		money.New(100, "USD"),
		money.New(200, "USD"),
		money.New(300, "USD"),
	)
	fmt.Println(result.Display())
	// Output: $6.00
}

func ExampleMin() {
	result, _ := money.Min(
		money.New(300, "USD"),
		money.New(100, "USD"),
		money.New(200, "USD"),
	)
	fmt.Println(result.Display())
	// Output: $1.00
}

func ExampleMax() {
	result, _ := money.Max(
		money.New(100, "USD"),
		money.New(300, "USD"),
		money.New(200, "USD"),
	)
	fmt.Println(result.Display())
	// Output: $3.00
}

func ExampleAverage() {
	result, _ := money.Average(
		money.New(100, "USD"),
		money.New(200, "USD"),
		money.New(300, "USD"),
	)
	fmt.Println(result.Display())
	// Output: $2.00
}

func ExampleCurrencyForCountry() {
	code, ok := money.CurrencyForCountry("US")
	fmt.Println(code, ok)
	// Output: USD true
}

func ExampleNewForCountry() {
	m, _ := money.NewForCountry(1000, "US")
	fmt.Println(m.Display())
	// Output: $10.00
}

func ExampleMoney_ToWords() {
	m := money.New(1, "USD")
	fmt.Println(m.ToWords())
	// Output: one dollar only
}

func ExampleMoney_Display() {
	fmt.Println(money.New(123456789, "EUR").Display())
	// Output: €1,234,567.89
}

func ExampleMoney_LocaleFormat() {
	m := money.New(123456, "EUR") // €1,234.56

	fmt.Println(m.LocaleFormat("en-US")) // symbol before, dot decimal
	fmt.Println(m.LocaleFormat("de-DE")) // symbol after, comma decimal, dot thousand
	// Output:
	// €1,234.56
	// 1.234,56 €
}

func ExampleMoney_AccountingFormat() {
	profit := money.New(123456, "USD")
	loss := money.New(-123456, "USD")

	fmt.Println(profit.AccountingFormat("en-US"))
	fmt.Println(loss.AccountingFormat("en-US"))
	// Output:
	// $1,234.56
	// ($1,234.56)
}

func ExampleMoney_Sign() {
	fmt.Println(money.New(100, "USD").Sign())
	fmt.Println(money.New(0, "USD").Sign())
	fmt.Println(money.New(-100, "USD").Sign())
	// Output:
	// 1
	// 0
	// -1
}

func ExampleMoney_MarshalBinary() {
	m := money.New(1234, "USD")
	b, _ := m.MarshalBinary()
	fmt.Println(len(b)) // 8 bytes int64 + 3 bytes "USD"
	// Output: 11
}

func ExampleMoney_MarshalText() {
	m := money.New(1234, "USD")
	b, _ := m.MarshalText()
	fmt.Println(string(b))
	// Output: 1234 USD
}

func ExampleNullMoney() {
	// Valid NullMoney
	n := money.NewNullMoney(money.New(100, "USD"))
	fmt.Println(n.Valid, n.Money.Display())

	// Null
	var empty money.NullMoney
	fmt.Println(empty.Valid)
	// Output:
	// true $1.00
	// false
}
