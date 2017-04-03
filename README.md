# Money

![alt text](http://i.imgur.com/c3XmCC6.jpg "Money")

[![Go Report Card](https://goreportcard.com/badge/github.com/rhymond/go-money)](https://goreportcard.com/report/github.com/rhymond/go-money)
[![Coverage Status](https://coveralls.io/repos/github/Rhymond/go-money/badge.svg?branch=wip)](https://coveralls.io/github/Rhymond/go-money?branch=wip)
[![Build Status](https://travis-ci.org/Rhymond/go-money.svg?branch=master)](https://travis-ci.org/Rhymond/go-money)

**GoMoney** provides ability to work with [monetary value using a currency's smallest unit](https://martinfowler.com/eaaCatalog/money.html).
Package allows you to use basic Money operations like rounding, splitting or allocating without losing a penny.
You shouldn't use float for monetary values, since they always carry small rounding differences.

```go
package main

import "github.com/rhymond/go-money"

func main() {
	pound := money.New(100, "GBP")
	twoPounds := pound.Add(pound)

	parties := twoPounds.Split(3)
	parties[0].Display() // £0.67
	parties[1].Display() // £0.67
	parties[2].Display() // £0.66
}

```
Quick start
-
Get the package via

``` bash
$ go get github.com/rhymond/go-money
```

## Features
* Provides a Money struct which stores information about an Money amount value and it's currency.
* Provides a ```Money.Amount``` struct which encapsulates all information about a monetary unit.
* Represents monetary values as integers, in cents. This avoids floating point rounding errors.
* Represents currency as ```Money.Currency``` instances providing a high level of flexibility.

Usage
-
### Instantiation
Initialise Money by using smallest unit value (e.g 100 represents 1 pound). Use ISO 4217 Currency Code to set money Currency
```go
pound := money.New(100, "GBP")
```
Comparison
-
**Go-money** lets you to use base compare operations like:

* Equals
* GreaterThan
* GreaterThanOrEqual
* LessThan
* LessThanOrEqual

In order to use them currencies must be equal

```go
pound := money.New(100, "GBP")
twoPounds := money.New(200, "GBP")
twoEuros := money.New(200, "EUR")

pound.GreaterThan(twoPounds) // false
pound.LessThan(twoPounds) // true
twoPounds.Equals(twoEuros) // Error: Currencies don't match
```
Asserts
-
* IsZero
* isNegative
* IsPositive

#### Zero value

To assert if Money value is equal zero use `IsZero()`

```go
pound := money.New(100, "GBP")
result := pound.IsZero(pound) // false
```

#### Positive value

To assert if Money value is more than zero `IsPositive()`

```go
pound := money.New(100, "GBP")
pound.IsPositive(pound) // true
```

#### Negative value

To assert if Money value is less than zero `IsNegative()`

```go
pound := money.New(100, "GBP")
pound.IsNegative(pound) // false
```

Operations
-
* Add
* Subtract
* Divide
* Modulus
* Multiply

In Order to use operations between Money structures Currencies must be equal

#### Addition

Additions can be performed using `Add()`.

```go
pound := money.New(100, "GBP")
twoPounds := money.New(200, "GBP")
 
result := pound.Add(twoPounds).Display() // £3.00
```

#### Subtraction

Subtraction can be performed using `Subtract()`.

```go
pound := money.New(100, "GBP")
twoPounds := money.New(200, "GBP")
 
result := pound.Subtract(twoPounds).Display() // -£3.00
```

#### Multiplication 

Multiplication can be performed using `Multiply()`.

```go
pound := money.New(100, "GBP")
 
result := pound.Multiply(2).Display() // £4.00
```

#### Division 

Division can be performed using `Divide()`.

```go
pound := money.New(100, "GBP")
 
result := pound.Divide(2).Display() // £4.00
```

There is possibilities to lose pennies by using division operation e.g:
```go
money.New(100, "GBP").Divide(3).Display() // £0.33
```
In order to split amount without losing pennies use `Split()` operation.


#### Absolute

Return `absolute` value of Money structure

```go
pound := money.New(-100, "GBP")
 
result := pound.Absolute().Display() // £1.00
```

#### Negative

Return `negative` value of Money structure

```go
pound := money.New(100, "GBP")
 
result := pound.Negative().Display() // -£1.00
```

Allocation
-

* Split
* Allocate

#### Splitting

In order to split Money for parties without losing any penny use `Split()`. 

After division leftover pennies will be distributed round-robin amongst the parties. This means that parties listed first will likely receive more pennies than ones that are listed later

```go
pound := money.New(100, "GBP")
parties := pound.Split(3)

parties[0].Display() // £0.34
parties[1].Display() // £0.33
parties[2].Display() // £0.33
```

#### Allocation

To perform allocation operation use `Allocate()`.

It lets split money by given ratios without losing pennies and as Split operations distributes leftover pennies amongst the parties with round-robin principle.

```go
pound := money.New(100, "GBP")
parties := pound.Split([]int{33, 33, 33})

parties[0].Display() // £0.34
parties[1].Display() // £0.33
parties[2].Display() // £0.33
```

Format
-

To format and return Money as string use `Display()`. 

```go
money.New(123456789, "EUR").Display() // €1,234,567.89
```

Contributing
-
Thank you for considering contributing! 
Please use GitHub issues and Pull Requests for Contributing.

License
-
The MIT License (MIT). Please see License File for more information.



[![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](https://github.com/Rhymond/go-money)