# Money

![alt text](http://i.imgur.com/c3XmCC6.jpg "Money")

[![Go Report Card](https://goreportcard.com/badge/github.com/rhymond/go-money)](https://goreportcard.com/report/github.com/rhymond/go-money)
[![Coverage Status](https://coveralls.io/repos/github/Rhymond/go-money/badge.svg?branch=master)](https://coveralls.io/github/Rhymond/go-money?branch=master)
[![Build Status](https://travis-ci.org/Rhymond/go-money.svg?branch=master)](https://travis-ci.org/Rhymond/go-money)
[![GoDoc](https://godoc.org/github.com/Rhymond/go-money?status.svg)](https://godoc.org/github.com/Rhymond/go-money)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**GoMoney** provides ability to work with [monetary value using a currency's smallest unit](https://martinfowler.com/eaaCatalog/money.html).
This package provides basic and precise Money operations such as rounding, splitting and allocating.  Monetary values should not be stored as floats due to small rounding differences.

```go
package main

import "github.com/Rhymond/go-money"

func main() {
    pound := money.New(100, "GBP")
    twoPounds, err := pound.Add(pound)

    if err != nil {
        log.Fatal(err)
    }

    parties, err := twoPounds.Split(3)

    if err != nil {
        log.Fatal(err)
    }

    parties[0].Display() // £0.67
    parties[1].Display() // £0.67
    parties[2].Display() // £0.66
}

```
Quick start
-
Get the package:

``` bash
$ go get github.com/Rhymond/go-money
```

## Features
* Provides a Money struct which stores information about an Money amount value and its currency.
* Provides a ```Money.Amount``` struct which encapsulates all information about a monetary unit.
* Represents monetary values as integers, in cents. This avoids floating point rounding errors.
* Represents currency as ```Money.Currency``` instances providing a high level of flexibility.

Usage
-
### Initialization
Initialize Money by using smallest unit value (e.g 100 represents 1 pound). Use ISO 4217 Currency Code to set money Currency
```go
pound := money.New(100, "GBP")
```
Comparison
-
**Go-money** provides base compare operations like:

* Equals
* GreaterThan
* GreaterThanOrEqual
* LessThan
* LessThanOrEqual

Comparisons must be made between the same currency units.

```go
pound := money.New(100, "GBP")
twoPounds := money.New(200, "GBP")
twoEuros := money.New(200, "EUR")

pound.GreaterThan(twoPounds) // false, nil
pound.LessThan(twoPounds) // true, nil
twoPounds.Equals(twoEuros) // false, error: Currencies don't match
```
Asserts
-
* IsZero
* IsNegative
* IsPositive

#### Zero value

To assert if Money value is equal to zero use `IsZero()`

```go
pound := money.New(100, "GBP")
result := pound.IsZero(pound) // false
```

#### Positive value

To assert if Money value is more than zero use `IsPositive()`

```go
pound := money.New(100, "GBP")
pound.IsPositive(pound) // true
```

#### Negative value

To assert if Money value is less than zero use `IsNegative()`

```go
pound := money.New(100, "GBP")
pound.IsNegative(pound) // false
```

Operations
-
* Add
* Subtract
* Divide
* Multiply
* Absolute
* Negative

Comparisons must be made between the same currency units.

#### Addition

Additions can be performed using `Add()`.

```go
pound := money.New(100, "GBP")
twoPounds := money.New(200, "GBP")

result, err := pound.Add(twoPounds) // £3.00, nil
```

#### Subtraction

Subtraction can be performed using `Subtract()`.

```go
pound := money.New(100, "GBP")
twoPounds := money.New(200, "GBP")

result, err := pound.Subtract(twoPounds) // -£1.00, nil
```

#### Multiplication

Multiplication can be performed using `Multiply()`.

```go
pound := money.New(100, "GBP")

result := pound.Multiply(2) // £2.00
```

#### Division

Division can be performed using `Divide()`.

```go
pound := money.New(100, "GBP")

result := pound.Divide(2) // £0.50
```

There is possibilities to lose pennies by using division operation e.g:
```go
money.New(100, "GBP").Divide(3) // £0.33
```
In order to split amount without losing use `Split()` operation.


#### Absolute

Return `absolute` value of Money structure

```go
pound := money.New(-100, "GBP")

result := pound.Absolute() // £1.00
```

#### Negative

Return `negative` value of Money structure

```go
pound := money.New(100, "GBP")

result := pound.Negative() // -£1.00
```

Allocation
-

* Split
* Allocate

#### Splitting

In order to split Money for parties without losing any pennies due to rounding differences, use `Split()`.

After division leftover pennies will be distributed round-robin amongst the parties. This means that parties listed first will likely receive more pennies than ones that are listed later.

```go
pound := money.New(100, "GBP")
parties, err := pound.Split(3)

if err != nil {
    log.Fatal(err)
}

parties[0].Display() // £0.34
parties[1].Display() // £0.33
parties[2].Display() // £0.33
```

#### Allocation

To perform allocation operation use `Allocate()`.

It splits money using the given ratios without losing pennies and as Split operations distributes leftover pennies amongst the parties with round-robin principle.

```go
pound := money.New(100, "GBP")
// Allocate is variadic function which can receive ratios as
// slice (int[]{33, 33, 33}...) or separated by a comma integers
parties, err := pound.Allocate(33, 33, 33)

if err != nil {
    log.Fatal(err)
}

parties[0].Display() // £0.34
parties[1].Display() // £0.33
parties[2].Display() // £0.33
```

Format
-

To format and return Money as a string use `Display()`.

```go
money.New(123456789, "EUR").Display() // €1,234,567.89
```
To format and return Money as a float64 representing the amount value in the currency's subunit use `AsSubunits()`.

```go
money.New(123456789, "EUR").AsSubunits() // 1234567.89
```

Contributing
-
Thank you for considering contributing!
Please use GitHub issues and Pull Requests for contributing.

License
-
The MIT License (MIT). Please see License File for more information.



[![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](https://github.com/Rhymond/go-money)
