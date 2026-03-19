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

// Locale-aware formatting
total.LocaleFormat("en-US")  // "$11.92"
total.LocaleFormat("de-DE")  // "11,92 $"
total.AccountingFormat("en-US") // "$11.92"
money.New(-500, money.USD).AccountingFormat("en-US") // "($5.00)"

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
  - [Comparisons & Predicates](#comparisons--predicates)
  - [Arithmetic Operations](#arithmetic-operations)
  - [Float Multiplication & Division](#float-multiplication--division)
  - [Percentage & Rounding](#percentage--rounding)
  - [Range Clamping](#range-clamping)
  - [Allocation](#allocation)
  - [Exchange Rates & Conversion](#exchange-rates--conversion)
  - [Aggregates](#aggregates)
  - [Amount as Parts](#amount-as-parts)
  - [Locale-Aware Formatting](#locale-aware-formatting)
  - [Accounting Format](#accounting-format)
  - [JSON Serialization](#json-serialization)
  - [SQL / Database Integration](#sql--database-integration)
  - [Binary Marshaling (MongoDB / BSON)](#binary-marshaling-mongodb--bson)
  - [Text Marshaling](#text-marshaling)
  - [XML Marshaling](#xml-marshaling)
  - [Nullable Money](#nullable-money)
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
    price := money.NewFromFloat(99.99, money.USD)  // $99.99

    // Arithmetic
    tax := price.Percentage(7.5)
    total, _ := price.Add(tax)

    // Locale-aware display
    fmt.Println(total.LocaleFormat("en-US"))  // "$107.49"
    fmt.Println(total.LocaleFormat("de-DE"))  // "107,49 $"
    fmt.Println(total.LocaleFormat("fr-FR"))  // "107,49 $"

    // Accounting style (negatives in parentheses)
    refund := total.Negative()
    fmt.Println(refund.AccountingFormat("en-US")) // "($107.49)"

    // Penny-perfect allocation
    shares, _ := total.Allocate(50, 30, 20)
    // $53.75, $32.25, $21.49

    // Exchange rate
    rate, _ := money.NewExchangeRate(money.USD, money.EUR, 0.92)
    euros, _ := rate.Convert(total)
    fmt.Println(euros.Display())  // "€98.89"

    // Country-based creation
    jpy, _ := money.NewForCountry(5000, "JP") // ¥5000

    // Multi-serialization
    b, _ := total.MarshalBinary()  // compact binary (works with MongoDB BSON)
    t2, _ := total.MarshalText()   // "10749 USD"
    _ = b; _ = t2

    // Nullable for SQL NULL columns
    n := money.NewNullMoney(&total)
    fmt.Println(n.Valid) // true
    _ = jpy
}
```

---

## API Reference

### Creating Money

```go
money.New(1099, "USD")                   // $10.99 (from smallest unit)
money.New(100, money.GBP)                // £1.00  (using constant)
money.NewFromFloat(10.99, "USD")         // $10.99 (from float)
money.NewForCountry(1000, "JP")          // ¥1000  (from country code)
money.New(0, "EUR").WithAmount(250)      // €2.50  (clone with new amount)
```

### Accessing Values

| Method | Returns | Description |
|--------|---------|-------------|
| `Amount()` | `int64` | Raw amount in smallest unit (e.g. 1099 for $10.99) |
| `AsFloat64()` | `float64` | Main-unit float — display only |
| `Currency()` | `*Currency` | Associated currency |
| `AsParts()` | `(int64, int64)` | Whole and fractional parts |
| `Display()` | `string` | Formatted e.g. `"$10.99"` |
| `String()` | `string` | Same as `Display()` (implements `fmt.Stringer`) |
| `LocaleFormat(locale)` | `string` | Locale-aware format e.g. `"10,99 €"` |
| `AccountingFormat(locale)` | `string` | Accounting style; negatives in parens |
| `ToWords()` | `string` | English words e.g. `"ten dollar only"` |
| `Sign()` | `int` | -1 / 0 / +1 |
| `IsWhole()` | `bool` | `true` if no fractional sub-unit |

### Comparisons & Predicates

```go
m.Equals(om)             (bool, error)
m.GreaterThan(om)        (bool, error)
m.GreaterThanOrEqual(om) (bool, error)
m.LessThan(om)           (bool, error)
m.LessThanOrEqual(om)    (bool, error)
m.Compare(om)            (int, error)   // -1 / 0 / +1
m.SameCurrency(om)       bool
m.IsZero()               bool
m.IsPositive()           bool
m.IsNegative()           bool
m.IsWhole()              bool
m.Sign()                 int
```

### Arithmetic Operations

```go
m.Add(om)      (*Money, error)  // m + om
m.Subtract(om) (*Money, error)  // m - om
m.Multiply(n)  *Money           // m × n  (integer)
m.Divide(n)    *Money           // m ÷ n  (truncated)
m.Absolute()   *Money           // |m|
m.Negative()   *Money           // negate
```

All methods return a new `*Money`; the receiver is never modified.

### Float Multiplication & Division

```go
price.MultiplyFloat(1.21)                    // ×1.21 (VAT factor)
annual.MultiplyFloat(1.0 / 52)               // weekly rate

m.DivideWithRounding(3, money.RoundHalfUp)   // ÷3, round half-up
m.DivideWithRounding(3, money.RoundHalfEven) // ÷3, banker's rounding
```

### Percentage & Rounding

```go
price.Percentage(8.5)   //  8.5% of price
price.Percentage(-10)   // -10% discount

m.Round()                                    // RoundHalfUp default
m.RoundWithMode(money.RoundHalfUp)           // ties away from zero
m.RoundWithMode(money.RoundHalfDown)         // ties toward zero
m.RoundWithMode(money.RoundHalfEven)         // banker's rounding
m.RoundWithMode(money.RoundUp)               // always away from zero
m.RoundWithMode(money.RoundDown)             // truncate
```

### Range Clamping

```go
min := money.New(100, "USD")
max := money.New(500, "USD")

money.New(50,  "USD").Clamp(min, max) // → $1.00 (clamped to min)
money.New(300, "USD").Clamp(min, max) // → $3.00 (unchanged)
money.New(600, "USD").Clamp(min, max) // → $5.00 (clamped to max)
```

### Allocation

```go
// Equal split (round-robin remainder)
shares, _ := money.New(100, "GBP").Split(3)
// £0.34, £0.33, £0.33

// Weighted by ratio
parts, _ := money.New(100, "USD").Allocate(50, 30, 20)
// $0.50, $0.30, $0.20
```

### Exchange Rates & Conversion

```go
rate, _ := money.NewExchangeRate("USD", "EUR", 0.92)
euros, _ := rate.Convert(money.New(1000, "USD")) // €9.20
inv := rate.Invert()                              // EUR→USD
```

### Aggregates

```go
prices := []*money.Money{ money.New(100, "USD"), money.New(200, "USD"), money.New(300, "USD") }

money.Sum(prices...)      // $6.00
money.Min(prices...)      // $1.00
money.Max(prices...)      // $3.00
money.Average(prices...)  // $2.00
```

### Amount as Parts

```go
money.New(1234, "USD").AsParts() // (12, 34) → $12.34
money.New(-550, "GBP").AsParts() // (-5, 50) → -£5.50
```

### Locale-Aware Formatting

`LocaleFormat` formats using locale-specific decimal/thousand separators and symbol placement. Supports **60+ locales** with exact match and language-subtag fallback. Unrecognised locales fall back to `Display()`.

```go
m := money.New(123456, "EUR") // €1,234.56

m.LocaleFormat("en-US") // "€1,234.56"  — symbol before, dot decimal
m.LocaleFormat("de-DE") // "1.234,56 €" — symbol after, comma decimal
m.LocaleFormat("fr-FR") // "1 234,56 €" — symbol after, narrow-NBSP thousand
m.LocaleFormat("ja-JP") // "€1,234.56"  — symbol before, comma thousand
m.LocaleFormat("pt-BR") // "€ 1.234,56" — symbol before with space
m.LocaleFormat("vi-VN") // "1.234,56€"  — symbol after, no space
```

You can add custom locales at runtime:

```go
money.LocaleFormats["my-LC"] = money.LocaleConfig{
    Decimal: ".", Thousand: "'", SymbolPos: money.SymbolAfterSpace,
}
```

### Accounting Format

Negative values in parentheses — standard for financial reports:

```go
money.New(123456,  "USD").AccountingFormat("en-US") // "$1,234.56"
money.New(-123456, "USD").AccountingFormat("en-US") // "($1,234.56)"
money.New(-100,    "EUR").AccountingFormat("de-DE") // "(1,00 €)"
```

### JSON Serialization

```go
m := money.New(1099, "USD")
b, _ := json.Marshal(m)      // {"amount":1099,"currency":"USD"}
json.Unmarshal(b, &m)        // round-trip
```

### SQL / Database Integration

```go
// Stored as "1099 USD" in any TEXT/VARCHAR column
db.Exec("INSERT INTO orders (price) VALUES (?)", price)

var price money.Money
row.Scan(&price)
```

### Binary Marshaling (MongoDB / BSON)

`Money` implements `encoding.BinaryMarshaler` and `encoding.BinaryUnmarshaler`. MongoDB's `go.mongodb.org/mongo-driver` automatically calls `MarshalBinary` for types that implement this interface, so go-money works with BSON natively — **no extra dependency needed**.

```go
// 11 bytes: 8-byte little-endian int64 + ASCII currency code
b, _ := money.New(1234, "USD").MarshalBinary()

var m money.Money
m.UnmarshalBinary(b) // m == $12.34
```

### Text Marshaling

```go
b, _ := money.New(1234, "USD").MarshalText()  // []byte("1234 USD")

var m money.Money
m.UnmarshalText([]byte("1234 USD"))            // m == $12.34
```

### XML Marshaling

```go
// <Money><amount>1234</amount><currency>USD</currency></Money>
xml.Marshal(money.New(1234, "USD"))

var m money.Money
xml.Unmarshal(xmlBytes, &m)
```

### Nullable Money

For SQL NULL columns and absent JSON fields:

```go
// Create
n := money.NewNullMoney(money.New(100, "USD")) // Valid=true
var empty money.NullMoney                       // Valid=false

// SQL
db.Exec("INSERT INTO t (price) VALUES (?)", n) // NULL when not valid
row.Scan(&n)

// JSON
json.Marshal(n)           // {"amount":100,"currency":"USD"} or null
json.Unmarshal(b, &n)

// Text
n.MarshalText()           // "100 USD" or ""
n.UnmarshalText(data)     // empty data → Valid=false
```

### Money in Words

```go
money.New(1, "USD").ToWords()               // "one dollar only"
money.GetCurrencyAmountWords(1.50, "USD")   // "one dollar and fifty cents only"
money.GetCurrencyAmountWords(10.99, "EUR")  // "ten euro and ninety-nine cent only"
```

Supports **80+ currencies**. Register custom ones:
```go
money.AddCurrencyMeta("BTC", "bitcoin", "satoshi")
```

### Country to Currency Mapping

```go
money.CurrencyForCountry("US")           // "USD", true
money.CurrencyForCountry("DE")           // "EUR", true
money.NewForCountry(1000, "JP")          // ¥1000 JPY
```

Covers **180+ countries**.

### Custom Currencies

```go
money.AddCurrency("BTC", "₿", "1$", ".", ",", 8)
money.AddCurrencyMeta("BTC", "bitcoin", "satoshi")
money.New(100000000, "BTC").Display() // "₿1.00000000"
```

### Currency Code Constants

```go
money.New(100, money.USD)  // not "USD" — catches typos at compile time
money.New(100, money.EUR)
// 80+ constants: USD, EUR, GBP, JPY, CHF, AUD, CAD, NZD…
```

---

## Supported Currencies

**Formatting** — 150+ ISO 4217 currencies via `currency.go`

**Money in Words** — 80+ currencies via `currencyToWords.go`

**Locale Formatting** — 60+ locales via `locale.go`
`en-US/GB/AU/CA/NZ/SG/IN` · `de-DE/AT/CH` · `fr-FR/BE/CA/CH` · `es-ES/MX/AR/CO/CL` · `it-IT/CH` · `nl-NL/BE` · `pt-BR/PT` · `ru-RU` · `uk-UA` · `tr-TR` · `pl-PL` · `sv-SE` · `nb-NO` · `da-DK` · `fi-FI` · `ja-JP` · `zh-CN/TW/HK` · `ko-KR` · `ar-SA/AE/EG` · `hi-IN` · `th-TH` · `id-ID` · `ms-MY` · `vi-VN` · `he-IL` · `el-GR` · `ro-RO` · `hu-HU` · `cs-CZ` · `sk-SK` · `bg-BG` · `hr-HR` · `sr-RS` · `kk-KZ` · `bn-BD`…

---

## Comparison with Other Libraries

| Feature | **go-money** | [bojanz/currency](https://github.com/bojanz/currency) | [govalues/money](https://github.com/govalues/money) | [Rhymond/go-money](https://github.com/Rhymond/go-money) |
|---------|:---:|:---:|:---:|:---:|
| Integer precision | ✅ | ✅ | ✅ | ✅ |
| Exchange rates | ✅ | manual only | ✅ | ❌ |
| 5 rounding modes | ✅ | ✅ | 4 | 0 |
| Split / Allocate | ✅ | ❌ | split only | ✅ |
| Aggregates (Sum/Min/Max/Avg) variadic | ✅ | ❌ | pairwise only | ❌ |
| **Locale-aware formatting (60+ locales)** | ✅ | ✅ CLDR/370 | ❌ | ❌ |
| **Accounting format** | ✅ | ✅ | ❌ | ❌ |
| Money in words (80+ currencies) | ✅ | ❌ | ❌ | ❌ |
| Country → Currency (180+ countries) | ✅ | limited | ❌ | ❌ |
| Float multiply / divide-with-rounding | ✅ | ❌ | ❌ | ❌ |
| Range clamping | ✅ | ❌ | ✅ | ❌ |
| Sign() | ✅ | ❌ | ✅ | ❌ |
| IsWhole() | ✅ | ❌ | ✅ | ❌ |
| JSON | ✅ | ✅ | ✅ | ✅ |
| SQL Scan / Value | ✅ | ✅ | ✅ | ✅ |
| **Binary (encoding.BinaryMarshaler)** | ✅ | ✅ | ✅ | ❌ |
| **Text (encoding.TextMarshaler)** | ✅ | ❌ | ✅ | ❌ |
| **XML (encoding.xml.Marshaler)** | ✅ | ❌ | ❌ | ❌ |
| **BSON / MongoDB** | ✅ via BinaryMarshaler | ❌ | ✅ native | ❌ |
| **NullMoney type** | ✅ | ❌ | ✅ NullCurrency | ❌ |
| Currency code constants | ✅ | ❌ | ❌ | ✅ |
| Custom currencies | ✅ | ✅ | ❌ | ✅ |
| Zero external dependencies | ✅ | ❌ uses `cockroachdb/apd` | ✅ | ✅ |
| **100% test coverage (CI enforced)** | ✅ | ❌ | ❌ | ❌ |
| **Fuzz tests** | ✅ | ❌ | ✅ | ❌ |
| **Benchmarks** | ✅ | ✅ | ✅ | ❌ |
| Runnable godoc examples | ✅ 30+ | partial | partial | ❌ |
| Min. Go version | 1.21 | 1.18 | **1.22** | 1.11 |

### When to choose each library

| Use case | **Best choice** |
|----------|-----------------|
| General fintech, e-commerce, billing | ✅ **go-money** — all-in-one, 100% coverage |
| Locale-aware display (CLDR, 370+ locales) | bojanz/currency for exhaustive CLDR; **go-money** covers 60 common locales |
| HFT / order-book / zero-alloc at µs level | govalues/money — zero-alloc, 23 ns ops |
| MongoDB / BSON | ✅ **go-money** — `MarshalBinary` works natively, no extra dep |
| XML serialization | ✅ **go-money** — only library with XML support |
| Accounting reports | ✅ **go-money** — `AccountingFormat` with locale awareness |
| Money in words (cheques, legal) | ✅ **go-money** — only library with this |
| Multi-country apps | ✅ **go-money** — 180-country currency mapping |

---

## Contributing

```bash
git clone https://github.com/im-adarsh/go-money
cd go-money
go test ./...                            # run tests
go test -race ./...                      # race detector
go test -coverprofile=c.out ./... && go tool cover -func=c.out
# total must show 100.0%
go test -bench=. ./...                   # run benchmarks
go test -run FuzzMoney -fuzz=FuzzMoneyAdd -fuzztime=30s .  # fuzz
```

---

## License

MIT — see [LICENSE](LICENSE).
