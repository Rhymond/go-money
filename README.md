# go-money

[![CI](https://github.com/im-adarsh/go-money/actions/workflows/ci.yml/badge.svg)](https://github.com/im-adarsh/go-money/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/im-adarsh/go-money.svg)](https://pkg.go.dev/github.com/im-adarsh/go-money)
[![Go Report Card](https://goreportcard.com/badge/github.com/im-adarsh/go-money)](https://goreportcard.com/report/github.com/im-adarsh/go-money)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**go-money** is a zero-dependency Go library for precise monetary value
arithmetic. It implements [Martin Fowler's Money
pattern](https://martinfowler.com/eaaCatalog/money.html) by storing every
amount as an integer in the currency's **smallest unit** (cents, pence,
paise…), eliminating the floating-point rounding errors that plague `float64`
approaches.

```go
import "github.com/im-adarsh/go-money"

price  := money.NewFromFloat(10.99, money.USD)  // $10.99
tax    := price.Percentage(8.5)                  // $0.93
total, _ := price.Add(tax)                       // $11.92

// Split three ways — no penny lost
shares, _ := total.Split(3)
// $3.98, $3.97, $3.97

// Convert to euros
rate, _ := money.NewExchangeRate("USD", "EUR", 0.92)
euros, _ := rate.Convert(total)
fmt.Println(euros.Display())                     // "€10.97"
```

---

## Table of Contents

- [Why integer arithmetic?](#why-integer-arithmetic)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
  - [Creating Money](#creating-money)
  - [Accessing Values](#accessing-values)
  - [Comparisons](#comparisons)
  - [Assertions](#assertions)
  - [Arithmetic Operations](#arithmetic-operations)
  - [Percentage](#percentage)
  - [Rounding Modes](#rounding-modes)
  - [Allocation](#allocation)
  - [Exchange Rates & Conversion](#exchange-rates--conversion)
  - [Aggregates](#aggregates)
  - [Amount as Parts](#amount-as-parts)
  - [Formatting & Display](#formatting--display)
  - [JSON Serialization](#json-serialization)
  - [SQL / Database Integration](#sql--database-integration)
  - [Money in Words](#money-in-words)
  - [Custom Currencies](#custom-currencies)
  - [Currency Code Constants](#currency-code-constants)
- [Supported Currencies](#supported-currencies)
- [Comparison with Other Libraries](#comparison-with-other-libraries)
- [Contributing](#contributing)
- [License](#license)

---

## Why integer arithmetic?

```go
fmt.Println(0.1 + 0.2)          // 0.30000000000000004
fmt.Println(0.1 + 0.2 == 0.3)   // false
```

**go-money** stores `£1.00` as the integer `100` (pence). All operations are
performed on integers — the result is always exact.

---

## Installation

```bash
go get github.com/im-adarsh/go-money
```

Requires **Go 1.21** or later. **Zero external dependencies.**

---

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    money "github.com/im-adarsh/go-money"
)

func main() {
    // Create from float (human-friendly)
    price := money.NewFromFloat(99.99, money.USD) // $99.99

    // Or from integer smallest unit (precise)
    tax := money.New(850, money.USD) // $8.50

    total, err := price.Add(tax)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(total.Display())    // "$108.49"
    fmt.Println(total.AsFloat64())  // 108.49

    // Penny-perfect split
    shares, _ := total.Split(3)
    for _, s := range shares {
        fmt.Println(s.Display()) // $36.17, $36.16, $36.16
    }

    // Currency conversion
    rate, _ := money.NewExchangeRate(money.USD, money.EUR, 0.92)
    euros, _ := rate.Convert(total)
    fmt.Println(euros.Display()) // "€99.81"
}
```

---

## API Reference

### Creating Money

```go
// From integer (smallest unit): 100 = $1.00
m := money.New(100, "USD")
m := money.New(100, money.USD)  // using constants

// From float: 1.25 = $1.25 → stored as 125 cents
m := money.NewFromFloat(1.25, "USD")
m := money.NewFromFloat(100.0, "JPY") // JPY has no subunit

// Negative values are supported
m := money.New(-100, "EUR")   // -€1.00
```

### Accessing Values

```go
m := money.New(1234, "USD")

m.Amount()           // int64: 1234   (raw smallest-unit integer)
m.AsFloat64()        // float64: 12.34 (main unit — for display only)
m.Currency()         // *Currency{Code:"USD", Fraction:2, Grapheme:"$", ...}
m.Currency().Code    // "USD"
m.Display()          // "$12.34"
m.String()           // "$12.34" (same as Display; implements fmt.Stringer)
whole, frac := m.AsParts() // (12, 34)
```

### Comparisons

All comparison methods require identical currencies and return an error otherwise.

```go
a := money.New(100, "USD")
b := money.New(200, "USD")

a.Equals(b)              // false, nil
a.GreaterThan(b)         // false, nil
a.GreaterThanOrEqual(b)  // false, nil
a.LessThan(b)            // true, nil
a.LessThanOrEqual(b)     // true, nil
a.Compare(b)             // -1, nil  (-1 / 0 / +1 like strings.Compare)
a.SameCurrency(b)        // true (no error — pure boolean)

// Cross-currency returns an error:
c := money.New(100, "EUR")
a.Equals(c)              // false, error("currencies don't match")
```

### Assertions

```go
money.New(0,   "USD").IsZero()     // true
money.New(100, "USD").IsPositive() // true
money.New(-1,  "USD").IsNegative() // true
```

### Arithmetic Operations

All operations return a **new** Money; the receiver is unchanged.

```go
a := money.New(1000, "USD")  // $10.00
b := money.New(500, "USD")   // $5.00

// Add / Subtract (require same currency)
sum, _  := a.Add(b)        // $15.00
diff, _ := a.Subtract(b)   // $5.00

// Multiply / Divide (integer division)
a.Multiply(3)              // $30.00
a.Divide(3)                // $3.33  (remainder truncated — use Split for lossless)

// Sign operations
money.New(-150, "USD").Absolute() // $1.50
money.New(150, "USD").Negative()  // -$1.50
```

### Percentage

```go
price := money.New(10000, "USD")  // $100.00

price.Percentage(10)    // $10.00  (10%)
price.Percentage(8.5)   // $8.50   (8.5%)
price.Percentage(0.5)   // $0.50   (0.5%)
price.Percentage(100)   // $100.00 (100%)

// Practical use: adding VAT
net := money.NewFromFloat(49.99, "GBP")
vat := net.Percentage(20)               // 20% VAT
gross, _ := net.Add(vat)
```

### Rounding Modes

`Round()` uses the library's default behaviour (half-down). For explicit
control, use `RoundWithMode()`:

| Constant | Description | Example (1.5 → ?) |
|---|---|---|
| `RoundHalfUp` | Ties round away from zero | → 2 |
| `RoundHalfDown` | Ties round toward zero | → 1 |
| `RoundHalfEven` | Banker's rounding (ties to nearest even) | → 2 |
| `RoundUp` | Always round away from zero | → 2 |
| `RoundDown` | Always truncate toward zero | → 1 |

```go
m := money.New(150, "USD")  // $1.50 (stored as 150 cents)

m.RoundWithMode(money.RoundHalfUp)   // $2.00
m.RoundWithMode(money.RoundHalfDown) // $1.00
m.RoundWithMode(money.RoundHalfEven) // $2.00 (2 is even)

// Banker's rounding minimises cumulative error in batch processing:
money.New(250, "USD").RoundWithMode(money.RoundHalfEven) // $2.00 (2 is even)
money.New(350, "USD").RoundWithMode(money.RoundHalfEven) // $4.00 (4 is even)
```

### Allocation

#### Even split

```go
// No penny is ever lost — leftover distributed round-robin to first parties
parts, _ := money.New(100, "USD").Split(3)
// $0.34, $0.33, $0.33   (total: $1.00)
```

#### Ratio-based allocation

```go
profit := money.New(10000, "USD")  // $100.00

parts, _ := profit.Allocate(50, 30, 20)
// $50.00, $30.00, $20.00

// Uneven ratios with leftover pennies
parts, _ = money.New(100, "USD").Allocate(33, 33, 33)
// $0.34, $0.33, $0.33   (total: $1.00)

// 70/30 investor split
parts, _ = money.New(10000, "USD").Allocate(70, 30)
// $70.00, $30.00
```

### Exchange Rates & Conversion

```go
// Create a rate: 1 USD = 0.92 EUR
rate, err := money.NewExchangeRate("USD", "EUR", 0.92)

usd := money.New(1000, "USD")  // $10.00
eur, _ := rate.Convert(usd)   // €9.20
eur.Display()                  // "€9.20"

// Invert the rate (EUR → USD)
inv := rate.Invert()
usd2, _ := inv.Convert(eur)   // ~$9.99 (due to rounding)

// Access rate metadata
rate.From()  // "USD"
rate.To()    // "EUR"
rate.Rate()  // 0.92

// Cross-currency chain
gbpToUsd, _ := money.NewExchangeRate("GBP", "USD", 1.27)
gbp := money.New(500, "GBP")    // £5.00
usd3, _ := gbpToUsd.Convert(gbp) // $6.35
```

### Aggregates

```go
prices := []*money.Money{
    money.New(999, "USD"),   // $9.99
    money.New(1499, "USD"),  // $14.99
    money.New(299, "USD"),   // $2.99
}

total, _   := money.Sum(prices...)     // $27.97
cheapest, _ := money.Min(prices...)    // $2.99
priciest, _ := money.Max(prices...)    // $14.99
avg, _      := money.Average(prices...) // $9.32 (truncated)
```

### Amount as Parts

```go
whole, frac := money.New(1234, "USD").AsParts()  // 12, 34   → $12.34
whole, frac  = money.New(100,  "GBP").AsParts()  //  1,  0   → £1.00
whole, frac  = money.New(-550, "EUR").AsParts()  // -5, 50   → -€5.50
whole, frac  = money.New(500,  "JPY").AsParts()  // 500, 0   → ¥500
```

### Formatting & Display

```go
money.New(123456789, "USD").Display()  // "$1,234,567.89"
money.New(123456789, "EUR").Display()  // "€1,234,567.89"
money.New(123456789, "GBP").Display()  // "£1,234,567.89"
money.New(123456789, "JPY").Display()  // "¥123,456,789"
money.New(100,       "AED").Display()  // "1.00 .د.إ"
money.New(-150,      "USD").Display()  // "-$1.50"

// fmt.Stringer is implemented — works directly with fmt verbs
fmt.Printf("Total: %s\n", money.New(1000, "USD")) // "Total: $10.00"
fmt.Println(money.New(500, "GBP"))                // "£5.00"
```

### JSON Serialization

```go
m := money.New(1234, "USD")

// Marshal
b, _ := json.Marshal(m)
// {"amount":1234,"currency":"USD"}

// Unmarshal
var restored money.Money
json.Unmarshal(b, &restored)
restored.Display() // "$12.34"

// Works in structs
type Order struct {
    ID    int         `json:"id"`
    Price money.Money `json:"price"`
}
```

### SQL / Database Integration

`Money` implements `database/sql/driver.Valuer` and `sql.Scanner`, allowing
direct use with any `database/sql`-compatible driver (PostgreSQL, MySQL, SQLite…).

```go
// Stored as "<amount> <CURRENCY>", e.g. "1234 USD"

// Writing to database
price := money.New(1234, "USD")
db.Exec("INSERT INTO products (price) VALUES (?)", price)

// Reading from database
var price money.Money
db.QueryRow("SELECT price FROM products WHERE id = ?", 1).Scan(&price)
price.Display()  // "$12.34"

// Works in structs with database/sql
type Product struct {
    ID    int
    Price money.Money
}
```

### Money in Words

```go
money.New(100,   "USD").ToWords() // "one dollar only"
money.New(1,     "GBP").ToWords() // "one pound only"
money.New(1,     "EUR").ToWords() // "one euro only"
money.New(100,   "JPY").ToWords() // "one hundred yen only"
money.New(100,   "INR").ToWords() // "one hundred rupee only"
money.New(100,   "PHP").ToWords() // "one hundred pesos only"

// Direct call for sub-unit output (e.g. for cheques/invoices)
money.GetCurrencyAmountWords(1.25, "USD") // "one dollar and twenty-five cents only"
money.GetCurrencyAmountWords(0.75, "EUR") // "seventy-five cent only"

// Register a custom currency for words
money.AddCurrencyMeta("XYZ", "zorkmid", "zork")
```

80+ currencies supported — see [Supported Currencies](#supported-currencies).

### Custom Currencies

```go
// Register or override a currency
money.AddCurrency("BTC", "₿", "$1", ".", ",", 8)
money.New(100000000, "BTC").Display() // "₿1.00000000"

// Custom currency for words
money.AddCurrencyMeta("BTC", "bitcoin", "satoshi")
```

### Currency Code Constants

All major ISO 4217 currency codes are available as package-level constants:

```go
money.New(100, money.USD)  // $1.00
money.New(100, money.EUR)  // €1.00
money.New(100, money.GBP)  // £1.00
money.New(100, money.JPY)  // ¥100
money.New(100, money.INR)  // ₹1.00
```

See [`constants.go`](constants.go) for the full list.

---

## Supported Currencies

### Formatting (150+ currencies)

See [`currency.go`](currency.go) for the complete list.

Selected currencies:

| Code | Currency | Symbol | Fraction |
|------|----------|--------|----------|
| USD  | US Dollar | $ | 2 |
| EUR  | Euro | € | 2 |
| GBP  | Pound Sterling | £ | 2 |
| JPY  | Japanese Yen | ¥ | 0 |
| CNY  | Chinese Yuan | 元 | 2 |
| INR  | Indian Rupee | ₹ | 2 |
| KRW  | South Korean Won | ₩ | 0 |
| BHD  | Bahraini Dinar | .د.ب | 3 |
| KWD  | Kuwaiti Dinar | .د.ك | 3 |
| TND  | Tunisian Dinar | .د.ت | 3 |

### Words (ToWords — 80+ currencies)

Americas: USD, CAD, AUD, NZD, MXN, BRL, ARS, CLP, COP, PEN, BOB, UYU, TTD, JMD, DOP, BSD, BZD, GYD, SRD, PAB

Europe: EUR, GBP, CHF, SEK, NOK, DKK, PLN, CZK, HUF, RON, BGN, HRK, RUB, UAH, TRY, ISK, BAM, RSD, MKD, ALL

Asia-Pacific: JPY, CNY, INR, HKD, SGD, KRW, TWD, THB, MYR, IDR, PHP, VND, PKR, BDT, LKR, NPR, MMK, KHR, LAK, MNT, KZT, UZS, AZN, GEL

Middle East & Africa: AED, SAR, QAR, KWD, BHD, OMR, JOD, IQD, IRR, ILS, EGP, ZAR, NGN, KES, GHS, TZS, UGX, ETB, MUR, ZMW, MAD, TND, DZD, LYD

---

## Comparison with Other Libraries

| Feature | go-money | bojanz/currency | govalues/money |
|---------|----------|-----------------|----------------|
| Integer precision | ✅ | ✅ | ✅ |
| 150+ currencies | ✅ | ✅ | — |
| Multiple rounding modes | ✅ | ✅ | Banker's only |
| Exchange rate conversion | ✅ | Partial | ✅ |
| JSON marshaling | ✅ | ✅ | ✅ |
| SQL Scanner/Valuer | ✅ | ✅ | — |
| Money in words | ✅ (80+ currencies) | — | — |
| NewFromFloat | ✅ | — | — |
| Aggregate (Sum/Min/Max/Avg) | ✅ | — | — |
| Immutable values | ✅ | ✅ | ✅ |
| Zero dependencies | ✅ | ❌ (CLDR data) | ❌ |
| Currency constants | ✅ | — | — |

---

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for
guidelines. Report bugs and request features via [GitHub Issues](https://github.com/im-adarsh/go-money/issues).

---

## License

The MIT License (MIT). See [LICENSE](LICENSE) for details.
