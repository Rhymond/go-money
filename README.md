# go-money

[![CI](https://github.com/im-adarsh/go-money/actions/workflows/ci.yml/badge.svg)](https://github.com/im-adarsh/go-money/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/im-adarsh/go-money.svg)](https://pkg.go.dev/github.com/im-adarsh/go-money)
[![Go Report Card](https://goreportcard.com/badge/github.com/im-adarsh/go-money)](https://goreportcard.com/report/github.com/im-adarsh/go-money)
[![Coverage: 100%](https://img.shields.io/badge/coverage-100%25-brightgreen)](https://github.com/im-adarsh/go-money/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**go-money** is a zero-dependency Go library for precise monetary value arithmetic. It implements [Martin Fowler's Money pattern](https://martinfowler.com/eaaCatalog/money.html) by storing every amount as an integer in the currency's **smallest unit** (cents, pence, paise…), eliminating the floating-point rounding errors that plague `float64` approaches.

```go
import money "github.com/im-adarsh/go-money"

price  := money.NewFromFloat(10.99, money.USD) // $10.99
tax    := price.Percentage(8.5)                // $0.93
total, _ := price.Add(tax)                    // $11.92

// Split three ways — no penny lost
shares, _ := total.Split(3)
// $3.98, $3.97, $3.97

// Convert to euros
rate, _ := money.NewExchangeRate("USD", "EUR", 0.92)
euros, _ := rate.Convert(total)
fmt.Println(euros.Display()) // "€10.97"
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
  - [Float Multiplication & Division](#float-multiplication--division)
  - [Percentage](#percentage)
  - [Rounding Modes](#rounding-modes)
  - [Range Clamping](#range-clamping)
  - [Allocation](#allocation)
  - [Exchange Rates & Conversion](#exchange-rates--conversion)
  - [Aggregates](#aggregates)
  - [Amount as Parts](#amount-as-parts)
  - [Formatting & Display](#formatting--display)
  - [JSON Serialization](#json-serialization)
  - [SQL / Database Integration](#sql--database-integration)
  - [Money in Words](#money-in-words)
  - [Country to Currency Mapping](#country-to-currency-mapping)
  - [Custom Currencies](#custom-currencies)
  - [Currency Code Constants](#currency-code-constants)
- [Supported Currencies](#supported-currencies)
- [Comparison with Other Libraries](#comparison-with-other-libraries)
- [Contributing](#contributing)
- [License](#license)

---

## Why integer arithmetic?

```go
fmt.Println(0.1 + 0.2)        // 0.30000000000000004
fmt.Println(0.1 + 0.2 == 0.3) // false
```

**go-money** stores £1.00 as the integer `100` (pence). All operations are performed on integers — the result is always exact.

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
    money "github.com/im-adarsh/go-money"
)

func main() {
    // Create from float (human-friendly)
    price := money.NewFromFloat(99.99, money.USD) // $99.99

    // Arithmetic
    tax := price.Percentage(7.5)
    total, _ := price.Add(tax)
    fmt.Println(total.Display()) // "$107.49"

    // Penny-perfect allocation
    shares, _ := total.Allocate(50, 30, 20)
    for _, s := range shares {
        fmt.Println(s.Display())
    }
    // $53.75
    // $32.25
    // $21.49

    // Exchange rate
    rate, _ := money.NewExchangeRate(money.USD, money.EUR, 0.92)
    euros, _ := rate.Convert(total)
    fmt.Println(euros.Display()) // "€98.89"

    // Aggregates
    prices := []*money.Money{
        money.NewFromFloat(10.00, money.USD),
        money.NewFromFloat(25.50, money.USD),
        money.NewFromFloat(7.99, money.USD),
    }
    sum, _ := money.Sum(prices...)
    avg, _ := money.Average(prices...)
    fmt.Println(sum.Display(), avg.Display()) // "$43.49" "$14.49"
}
```

---

## API Reference

### Creating Money

```go
// From smallest unit (cents, pence, paise…)
m := money.New(1099, "USD")           // $10.99
m := money.New(100, money.GBP)        // £1.00 (using constant)

// From main unit float — multiplied by 10^Fraction, rounded
m := money.NewFromFloat(10.99, "USD") // $10.99 → 1099 cents
m := money.NewFromFloat(1.005, "USD") // rounds  → $1.01

// For a country code (ISO 3166-1 alpha-2)
m, err := money.NewForCountry(1000, "JP") // ¥1000
m, err := money.NewForCountry(500, "DE")  // €5.00

// Clone with a new raw amount (keeps currency)
fee := money.New(0, "EUR").WithAmount(250) // €2.50
```

### Accessing Values

| Method | Returns | Description |
|--------|---------|-------------|
| `Amount()` | `int64` | Raw amount in smallest unit (e.g. 1099 for $10.99) |
| `AsFloat64()` | `float64` | Main-unit float — display only, never use in calculations |
| `Currency()` | `*Currency` | Associated currency |
| `AsParts()` | `(int64, int64)` | Whole and fractional parts — (10, 99) for $10.99 |
| `Display()` | `string` | Formatted string e.g. `"$10.99"` |
| `String()` | `string` | Same as `Display()`, implements `fmt.Stringer` |
| `ToWords()` | `string` | English words e.g. `"ten dollar only"` |
| `IsWhole()` | `bool` | `true` if no fractional sub-unit remainder |

### Comparisons

```go
m.Equals(om)             (bool, error)
m.GreaterThan(om)        (bool, error)
m.GreaterThanOrEqual(om) (bool, error)
m.LessThan(om)           (bool, error)
m.LessThanOrEqual(om)    (bool, error)
m.Compare(om)            (int, error)  // -1 / 0 / +1
m.SameCurrency(om)       bool
```

All comparison methods return an error when the currencies differ.

### Assertions

```go
m.IsZero()     bool   // amount == 0
m.IsPositive() bool   // amount > 0
m.IsNegative() bool   // amount < 0
m.IsWhole()    bool   // no fractional sub-unit
```

### Arithmetic Operations

```go
m.Add(om)      (*Money, error)  // m + om  (same currency required)
m.Subtract(om) (*Money, error)  // m - om
m.Multiply(n)  *Money           // m × n  (integer multiplier)
m.Divide(n)    *Money           // m ÷ n  (truncated toward zero)
m.Absolute()   *Money           // |m|
m.Negative()   *Money           // negate
```

All operations return a new `*Money`; the receiver is never modified (immutable style).

### Float Multiplication & Division

```go
// Multiply by a float (VAT, interest rate, exchange factor)
withVAT := price.MultiplyFloat(1.21)        // price × 1.21
weekRate := annual.MultiplyFloat(1.0 / 52)

// Divide with configurable rounding (instead of silent truncation)
half := m.DivideWithRounding(2, money.RoundHalfUp)
half := m.DivideWithRounding(2, money.RoundHalfEven) // banker's rounding
third := m.DivideWithRounding(3, money.RoundDown)
```

### Percentage

```go
tax  := price.Percentage(8.5)   //  8.5% of price (rounded to nearest cent)
disc := price.Percentage(-10)   // -10% discount
```

### Rounding Modes

go-money provides five configurable rounding strategies:

| Mode | Description | Tie example |
|------|-------------|-------------|
| `RoundHalfUp` | Nearest; ties away from zero | 1.5 → 2, −1.5 → −2 |
| `RoundHalfDown` | Nearest; ties toward zero | 1.5 → 1 |
| `RoundHalfEven` | Nearest; ties to even (banker's rounding) | 0.5 → 0, 1.5 → 2, 2.5 → 2 |
| `RoundUp` | Always away from zero | 1.1 → 2, −1.1 → −2 |
| `RoundDown` | Truncate toward zero | 1.9 → 1, −1.9 → −1 |

```go
m := money.New(150, "USD") // $1.50

m.Round()                                    // $2.00 (default: RoundHalfUp)
m.RoundWithMode(money.RoundHalfUp)           // $2.00
m.RoundWithMode(money.RoundHalfDown)         // $1.00
m.RoundWithMode(money.RoundHalfEven)         // $2.00 (2 is even)
m.RoundWithMode(money.RoundUp)               // $2.00
m.RoundWithMode(money.RoundDown)             // $1.00
```

### Range Clamping

```go
min := money.New(100, "USD") // $1.00
max := money.New(500, "USD") // $5.00

money.New(50,  "USD").Clamp(min, max) // → $1.00 (clamped to min)
money.New(300, "USD").Clamp(min, max) // → $3.00 (unchanged)
money.New(600, "USD").Clamp(min, max) // → $5.00 (clamped to max)
```

Returns an error if currencies differ or if `min > max`.

### Allocation

Distribute money without losing a single penny — remainders are spread round-robin:

```go
// Equal split
shares, _ := money.New(100, "GBP").Split(3)
// £0.34, £0.33, £0.33

// Weighted by ratio
parts, _ := money.New(100, "USD").Allocate(50, 30, 20)
// $0.50, $0.30, $0.20

// Unequal ratios — leftover cent goes to first party
parts, _ = money.New(100, "USD").Allocate(33, 33, 33)
// $0.34, $0.33, $0.33
```

### Exchange Rates & Conversion

```go
rate, err := money.NewExchangeRate("USD", "EUR", 0.92)

// Convert — result rounded to target currency's smallest unit
euros, err := rate.Convert(money.New(1000, "USD")) // 920 cents → €9.20

// Invert
inv := rate.Invert() // EUR → USD at ~1.0869

// Accessors
rate.From() // "USD"
rate.To()   // "EUR"
rate.Rate() // 0.92
```

`NewExchangeRate` returns an error when rate ≤ 0 or `from == to`.

### Aggregates

```go
prices := []*money.Money{
    money.New(100, "USD"),
    money.New(200, "USD"),
    money.New(300, "USD"),
}

total, _ := money.Sum(prices...)      // $6.00
min, _   := money.Min(prices...)      // $1.00
max, _   := money.Max(prices...)      // $3.00
avg, _   := money.Average(prices...)  // $2.00
```

All aggregate functions require the same currency; they return an error for mixed currencies or empty input.

### Amount as Parts

```go
whole, frac := money.New(1234, "USD").AsParts()  // 12, 34   ($12.34)
whole, frac  = money.New(-550, "GBP").AsParts()  // -5, 50   (-£5.50)
whole, frac  = money.New(100,  "JPY").AsParts()  // 100, 0   (¥100)
```

### Formatting & Display

```go
money.New(123456789, "USD").Display() // "$1,234,567.89"
money.New(123456789, "EUR").Display() // "€1,234,567.89"
money.New(100, "JPY").Display()       // "¥100"
money.New(100, "GBP").Display()       // "£1.00"

// Custom formatter
f := money.NewFormatter(2, ".", ",", "£", "1$")
f.Format(123456) // "£1,234.56"
```

### JSON Serialization

`Money` implements `json.Marshaler` and `json.Unmarshaler`:

```go
m := money.New(1099, "USD")

b, _ := json.Marshal(m)
// {"amount":1099,"currency":"USD"}

var m2 money.Money
json.Unmarshal(b, &m2) // m2.Amount() == 1099, m2.Currency().Code == "USD"
```

### SQL / Database Integration

`Money` implements `driver.Valuer` and `sql.Scanner`:

```go
// Stored as the string "1099 USD"
_, err = db.Exec("INSERT INTO orders (price) VALUES (?)", price)

// Read back transparently
var price money.Money
row.Scan(&price) // $10.99 USD
```

Works with any `database/sql`-compatible driver (MySQL, PostgreSQL, SQLite, etc.).

### Money in Words

```go
// Via Money.ToWords() — converts the raw integer amount
money.New(1, "USD").ToWords()   // "one dollar only"
money.New(50, "GBP").ToWords()  // "fifty pound only"

// Via GetCurrencyAmountWords — pass a human-scale float for sub-unit words
money.GetCurrencyAmountWords(1.50, "USD")  // "one dollar and fifty cents only"
money.GetCurrencyAmountWords(10.99, "EUR") // "ten euro and ninety-nine cent only"
money.GetCurrencyAmountWords(0.75, "GBP")  // "seventy-five penny only"
money.GetCurrencyAmountWords(1000, "JPY")  // "one thousand yen only"
```

Supports **80+ currencies**. Add custom ones:

```go
money.AddCurrencyMeta("BTC", "bitcoin", "satoshi")
money.GetCurrencyAmountWords(1.0, "BTC") // "one bitcoin only"
```

### Country to Currency Mapping

Map ISO 3166-1 alpha-2 country codes to their primary ISO 4217 currency — covers **180+ countries**:

```go
code, ok := money.CurrencyForCountry("US") // "USD", true
code, ok  = money.CurrencyForCountry("JP") // "JPY", true
code, ok  = money.CurrencyForCountry("DE") // "EUR", true
code, ok  = money.CurrencyForCountry("gb") // "GBP", true  (case-insensitive)
code, ok  = money.CurrencyForCountry("XX") // "",    false

// Directly create Money for a country
m, err := money.NewForCountry(1000, "US")  // $10.00 USD
m, err  = money.NewForCountry(500,  "JP")  // ¥500  JPY
m, err  = money.NewForCountry(100,  "XX")  // error
```

### Custom Currencies

```go
// Register a new currency (persists for the process lifetime)
money.AddCurrency("BTC", "₿", "1$", ".", ",", 8)
m := money.New(100000000, "BTC") // 1.00000000 BTC
m.Display()                       // "₿1.00000000"

// Register English words
money.AddCurrencyMeta("BTC", "bitcoin", "satoshi")
money.New(1, "BTC").ToWords() // "one bitcoin only"
```

### Currency Code Constants

Use compile-time constants instead of raw strings to catch typos early:

```go
money.New(100, money.USD)  // not "USD"
money.New(100, money.EUR)
money.New(100, money.GBP)
// 80+ constants: USD, EUR, GBP, JPY, CHF, AUD, CAD, NZD, SGD, HKD, CNY,
//                INR, KRW, TWD, THB, MYR, IDR, PHP, VND, PKR, BDT, AED,
//                SAR, ILS, EGP, ZAR, NGN, KES, BRL, MXN, RUB, TRY, PLN…
```

---

## Supported Currencies

### Formatting (150+ via `currency.go`)

All ISO 4217 currencies with symbol, separators, and fraction digits:

`USD` `EUR` `GBP` `JPY` `CHF` `AUD` `CAD` `NZD` `SGD` `HKD` `CNY` `INR` `KRW` `BRL` `MXN` `SEK` `NOK` `DKK` `PLN` `CZK` `HUF` `RON` `RUB` `TRY` `ZAR` `AED` `SAR` `THB` `MYR` `IDR` `PHP` `VND` `PKR` `BDT` `EGP` `NGN` `KES` `GHS` `ETB` `MAD` `TND` `DZD` `LYD` `KWD` `BHD` `OMR` `JOD` `IQD` `QAR` `ILS`…

### Money in Words (80+ via `currencyToWords.go`)

`USD` `EUR` `GBP` `JPY` `CNY` `INR` `AUD` `CAD` `CHF` `SGD` `HKD` `KRW` `TWD` `THB` `MYR` `IDR` `PHP` `VND` `PKR` `BDT` `LKR` `NPR` `MMK` `AED` `SAR` `QAR` `KWD` `BHD` `OMR` `JOD` `ILS` `EGP` `ZAR` `NGN` `KES` `GHS` `ETB` `BRL` `MXN` `ARS` `COP` `PEN` `CLP` `RUB` `UAH` `TRY` `PLN` `SEK` `NOK` `DKK`…

---

## Comparison with Other Libraries

| Feature | **go-money** | [bojanz/currency](https://github.com/bojanz/currency) | [govalues/money](https://github.com/govalues/money) | [Rhymond/go-money](https://github.com/Rhymond/go-money) |
|---------|:---:|:---:|:---:|:---:|
| Integer precision (no float errors) | ✅ | ✅ | ✅ | ✅ |
| Exchange rates & conversion | ✅ | ❌ | ✅ | ❌ |
| 5 rounding modes | ✅ | basic | banker's only | basic |
| Split / Allocate (penny-perfect) | ✅ | ❌ | ❌ | ✅ |
| Aggregates (Sum / Min / Max / Avg) | ✅ | ❌ | ❌ | ❌ |
| Money in words (80+ currencies) | ✅ | ❌ | ❌ | ❌ |
| Country → Currency mapping (180+ countries) | ✅ | ❌ | ❌ | ❌ |
| Float multiply & divide-with-rounding | ✅ | ❌ | ❌ | ❌ |
| Range clamping (Clamp) | ✅ | ❌ | ❌ | ❌ |
| IsWhole / WithAmount helpers | ✅ | ❌ | ❌ | ❌ |
| JSON marshal / unmarshal | ✅ | ✅ | ✅ | ✅ |
| SQL Valuer / Scanner | ✅ | ✅ | ✅ | ✅ |
| `fmt.Stringer` | ✅ | ✅ | ✅ | ✅ |
| Currency code constants | ✅ | ❌ | ❌ | ✅ |
| Locale-aware CLDR formatting | ❌ | ✅ | ❌ | ❌ |
| Zero external dependencies | ✅ | ✅ | ✅ | ✅ |
| **100% test coverage (enforced in CI)** | ✅ | ❌ | ❌ | ❌ |
| Runnable godoc examples | ✅ | partial | partial | ❌ |
| Minimum Go version | 1.21 | 1.18 | **1.22** | 1.11 |

### When to choose each library

- **go-money** — best all-round choice: most features, 100% coverage, supports Go 1.21+. Pick this when you need allocation, aggregates, exchange, words, or country mapping in the same library.
- **bojanz/currency** — best for apps that need locale-aware number formatting (e.g. "1.234,56 €" vs "$1,234.56") using CLDR data.
- **govalues/money** — best for high-frequency trading systems where zero heap allocation matters and you can require Go 1.22.
- **shopspring/decimal** — best when you need arbitrary-precision decimal arithmetic without currency semantics.

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for setup, coding standards, and PR guidelines.

```bash
git clone https://github.com/im-adarsh/go-money
cd go-money
go test ./...                          # run all tests
go test -race ./...                    # race detector
go test -coverprofile=c.out ./... && go tool cover -func=c.out
# total must show 100.0%
```

---

## License

MIT — see [LICENSE](LICENSE).
