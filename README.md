# go-money

[![CI](https://github.com/im-adarsh/go-money/actions/workflows/ci.yml/badge.svg)](https://github.com/im-adarsh/go-money/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/im-adarsh/go-money.svg)](https://pkg.go.dev/github.com/im-adarsh/go-money)
[![Go Report Card](https://goreportcard.com/badge/github.com/im-adarsh/go-money)](https://goreportcard.com/report/github.com/im-adarsh/go-money)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**go-money** is a zero-dependency Go library for precise monetary value
arithmetic. It implements [Martin Fowler's Money
pattern](https://martinfowler.com/eaaCatalog/money.html) by storing every
amount as an integer in the currency's **smallest unit** (e.g. cents, pence,
paise), eliminating the floating-point rounding errors that plague naive
`float64` approaches.

```go
import "github.com/im-adarsh/go-money"

pound := money.New(100, "GBP")          // £1.00
twoPounds, _ := pound.Add(pound)        // £2.00

parties, _ := twoPounds.Split(3)
parties[0].Display() // £0.67
parties[1].Display() // £0.67
parties[2].Display() // £0.66
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
  - [Rounding](#rounding)
  - [Allocation](#allocation)
  - [Formatting & Display](#formatting--display)
  - [Amount as Parts](#amount-as-parts)
  - [JSON Serialization](#json-serialization)
  - [Money in Words](#money-in-words)
  - [Custom Currencies](#custom-currencies)
- [Supported Currencies](#supported-currencies)
- [Contributing](#contributing)
- [License](#license)

---

## Why integer arithmetic?

Floating-point arithmetic cannot represent many decimal fractions exactly:

```go
fmt.Println(0.1 + 0.2)          // 0.30000000000000004
fmt.Println(0.1 + 0.2 == 0.3)   // false
```

**go-money** stores `£1.00` as the integer `100` (pence).
All operations are performed on integers — the result is always exact.

---

## Installation

```bash
go get github.com/im-adarsh/go-money
```

Requires **Go 1.21** or later. The library has **no external dependencies**.

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
    // Amounts are in the currency's smallest unit.
    // 1000 cents  = $10.00
    price := money.New(1000, "USD")
    tax   := price.Percentage(8.5) // 8.5% tax

    total, err := price.Add(tax)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(price.Display()) // $10.00
    fmt.Println(tax.Display())   // $0.85
    fmt.Println(total.Display()) // $10.85

    // Split the bill three ways — no penny lost.
    shares, _ := total.Split(3)
    for _, s := range shares {
        fmt.Println(s.Display()) // $3.62, $3.62, $3.61 (sums to $10.85)
    }
}
```

---

## API Reference

### Creating Money

```go
// New creates a Money value.
// amount — integer in the currency's smallest unit (e.g. cents, pence).
// code   — ISO 4217 currency code (case-insensitive).
m := money.New(100, "USD")  // $1.00
m := money.New(100, "gbp")  // £1.00  (code is normalised to uppercase)
m := money.New(-50, "EUR")  // -€0.50
```

### Accessing Values

```go
m := money.New(1234, "USD")

m.Amount()          // int64: 1234  (raw smallest-unit value)
m.Currency()        // *Currency{Code:"USD", Fraction:2, Grapheme:"$", ...}
m.Currency().Code   // "USD"
m.Display()         // "$12.34"
```

### Comparisons

All comparison methods require **identical currencies**; they return an error
otherwise.

```go
a := money.New(100, "USD")
b := money.New(200, "USD")

a.Equals(b)             // false, nil
a.GreaterThan(b)        // false, nil
a.GreaterThanOrEqual(b) // false, nil
a.LessThan(b)           // true,  nil
a.LessThanOrEqual(b)    // true,  nil

// Compare returns -1, 0, or +1 (like strings.Compare).
a.Compare(b)            // -1, nil

// Cross-currency comparison returns an error:
c := money.New(100, "EUR")
a.Equals(c)             // false, error("currencies don't match")
```

### Assertions

These methods never return errors.

```go
money.New(0,   "USD").IsZero()     // true
money.New(100, "USD").IsPositive() // true
money.New(-1,  "USD").IsNegative() // true
```

### Arithmetic Operations

All arithmetic operations return a **new** Money instance; the receiver is
unchanged.

#### Add

```go
a, _ := money.New(100, "USD").Add(money.New(200, "USD"))
a.Display() // "$3.00"
```

#### Subtract

```go
a, _ := money.New(300, "USD").Subtract(money.New(100, "USD"))
a.Display() // "$2.00"
```

#### Multiply

```go
money.New(100, "USD").Multiply(3).Display() // "$3.00"
```

#### Divide

Integer division — fractional remainder is discarded.
Use [Split](#splitting) or [Allocate](#allocation-by-ratio) to avoid penny loss.

```go
money.New(100, "USD").Divide(3).Display() // "$0.33"  (1 cent lost)
```

#### Absolute

```go
money.New(-150, "USD").Absolute().Display() // "$1.50"
```

#### Negative

```go
money.New(150, "USD").Negative().Display() // "-$1.50"
```

### Percentage

```go
price := money.New(10000, "USD")    // $100.00

discount := price.Percentage(10)    // 10%  → $10.00
tax      := price.Percentage(8.5)   // 8.5% → $8.50
tip      := price.Percentage(18)    // 18%  → $18.00
```

The result is **rounded** to the nearest smallest unit.

### Rounding

`Round()` rounds to the nearest whole currency unit (e.g. nearest dollar).

```go
money.New(175, "USD").Round().Display() // "$2.00"
money.New(125, "USD").Round().Display() // "$1.00"
money.New(150, "USD").Round().Display() // "$2.00"
```

### Allocation

#### Splitting evenly

`Split(n)` divides `m` into `n` equal parts. Leftover pennies are distributed
to the **first** parties (round-robin), so the total is always preserved.

```go
parts, _ := money.New(100, "USD").Split(3)
// $0.34, $0.33, $0.33  (total: $1.00)
```

#### Allocation by ratio

`Allocate(ratios...)` distributes `m` according to arbitrary integer ratios.
Leftover pennies are again distributed round-robin from the first party.

```go
parts, _ := money.New(10000, "USD").Allocate(50, 30, 20)
// $5.00, $3.00, $2.00

parts, _ := money.New(100, "USD").Allocate(33, 33, 33)
// $0.34, $0.33, $0.33  (total: $1.00)

// Tip-splitting scenario: 70% / 30% split
parts, _ := money.New(10000, "USD").Allocate(70, 30)
// $7.00, $3.00
```

### Formatting & Display

`Display()` formats the amount using the currency's symbol, decimal separator,
and thousands separator as defined by the currency's metadata.

```go
money.New(123456789, "USD").Display() // "$1,234,567.89"
money.New(123456789, "EUR").Display() // "€1,234,567.89"
money.New(123456789, "GBP").Display() // "£1,234,567.89"
money.New(123456789, "JPY").Display() // "¥123,456,789"
money.New(100,       "AED").Display() // "1.00 .د.إ"
money.New(-150,      "USD").Display() // "-$1.50"
```

### Amount as Parts

`AsParts()` splits the raw integer into the whole-unit and fractional-unit
components — useful when you need to display or process each part separately.

```go
whole, frac := money.New(1234, "USD").AsParts() // 12, 34   → $12.34
whole, frac  = money.New(100,  "GBP").AsParts() //  1,  0   → £1.00
whole, frac  = money.New(-550, "GBP").AsParts() // -5, 50   → -£5.50
whole, frac  = money.New(500,  "JPY").AsParts() // 500, 0   → ¥500
```

### JSON Serialization

`Money` implements `json.Marshaler` and `json.Unmarshaler`.

```go
m := money.New(1234, "USD")

// Marshal
b, _ := json.Marshal(m)
// {"amount":1234,"currency":"USD"}

// Unmarshal
var restored money.Money
json.Unmarshal(b, &restored)
restored.Display() // "$12.34"
```

This makes it straightforward to store monetary values in databases or transmit
them via REST / gRPC APIs.

### Money in Words

`ToWords()` converts a monetary amount to English words.
The currency must be registered in `CountryCurrencyMeta`
(see [Supported Currencies](#supported-currencies)).

```go
money.New(100,   "USD").ToWords() // "one dollar only"
money.New(150,   "USD").ToWords() // "one dollar and fifty cents only"
money.New(50,    "USD").ToWords() // "fifty cents only"
money.New(100,   "GBP").ToWords() // "one pound only"
money.New(100,   "EUR").ToWords() // "one euro only"
money.New(10050, "INR").ToWords() // "one hundred rupee and fifty paisa only"
money.New(100,   "JPY").ToWords() // "one hundred yen only"
money.New(100,   "PHP").ToWords() // "one hundred pesos only"
```

#### Registering a custom currency for words

```go
money.AddCurrencyMeta("XYZ", "zorkmid", "zork")
money.New(150, "XYZ").ToWords() // "one zorkmid and fifty zork only"
```

### Custom Currencies

You can register currencies that are not built-in, or override existing ones.

```go
// AddCurrency(code, grapheme, template, decimal, thousand string, fraction int)
money.AddCurrency("BTC", "₿", "$1", ".", ",", 8)
money.New(100000000, "BTC").Display() // "₿1.00000000"
```

**Template** controls symbol placement:
- `"$1"` → symbol before amount (e.g. `$1.00`)
- `"1 $"` → amount then symbol with space (e.g. `1.00 kr`)
- `"1$"`  → amount directly followed by symbol (e.g. `1.00Gs`)

---

## Supported Currencies

### Formatting (Display)

The following ISO 4217 codes have built-in formatting metadata:

| Code | Currency | Symbol |
|------|----------|--------|
| AED  | UAE Dirham | .د.إ |
| AUD  | Australian Dollar | $ |
| BRL  | Brazilian Real | R$ |
| CAD  | Canadian Dollar | $ |
| CHF  | Swiss Franc | CHF |
| CNY  | Chinese Yuan | 元 |
| CZK  | Czech Koruna | Kč |
| DKK  | Danish Krone | kr |
| EUR  | Euro | € |
| GBP  | Pound Sterling | £ |
| HKD  | Hong Kong Dollar | $ |
| HUF  | Hungarian Forint | Ft |
| IDR  | Indonesian Rupiah | Rp |
| ILS  | Israeli Shekel | ₪ |
| INR  | Indian Rupee | ₹ |
| JPY  | Japanese Yen | ¥ |
| KRW  | South Korean Won | ₩ |
| KWD  | Kuwaiti Dinar | .د.ك |
| MXN  | Mexican Peso | $ |
| MYR  | Malaysian Ringgit | RM |
| NGN  | Nigerian Naira | ₦ |
| NOK  | Norwegian Krone | kr |
| NZD  | New Zealand Dollar | $ |
| PHP  | Philippine Peso | ₱ |
| PLN  | Polish Złoty | zł |
| RUB  | Russian Ruble | ₽ |
| SAR  | Saudi Riyal | ﷼ |
| SEK  | Swedish Krona | kr |
| SGD  | Singapore Dollar | $ |
| THB  | Thai Baht | ฿ |
| TRY  | Turkish Lira | ₺ |
| TWD  | New Taiwan Dollar | NT$ |
| UAH  | Ukrainian Hryvnia | ₴ |
| USD  | US Dollar | $ |
| VND  | Vietnamese Dong | ₫ |
| ZAR  | South African Rand | R |
| ZMW  | Zambian Kwacha | ZK |

And ~120 more — see [`currency.go`](currency.go) for the full list.

### Words (ToWords)

The following codes are supported by `ToWords()`:

Americas: USD, CAD, AUD, NZD, MXN, BRL, ARS, CLP, COP, PEN, BOB, UYU, TTD, JMD, DOP, BSD, BZD, GYD, SRD, PAB

Europe: EUR, GBP, CHF, SEK, NOK, DKK, PLN, CZK, HUF, RON, BGN, HRK, RUB, UAH, TRY, ISK, BAM, RSD, MKD, ALL

Asia-Pacific: JPY, CNY, INR, HKD, SGD, KRW, TWD, THB, MYR, IDR, PHP, VND, PKR, BDT, LKR, NPR, MMK, KHR, LAK, MNT, KZT, UZS, AZN, GEL

Middle East & Africa: AED, SAR, QAR, KWD, BHD, OMR, JOD, IQD, IRR, ILS, EGP, ZAR, NGN, KES, GHS, TZS, UGX, ETB, MUR, ZMW, MAD, TND, DZD, LYD

---

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for
guidelines on reporting bugs, requesting features, and submitting pull requests.

---

## License

The MIT License (MIT). See [LICENSE](LICENSE) for details.
