# Money

![alt text](http://i.imgur.com/c3XmCC6.jpg "Money")

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
	parties[0].format() // £0.67
	parties[1].format() // £0.67
	parties[2].format() // £0.66
}

```
## Quick start
Get the package via

``` bash
$ go get github.com/rhymond/go-money
```

## Features
* Provides a Money struct which stores information about an Money amount value and it's currency.
* Provides a ```Money.Amount``` struct which encapsulates all information about a monetary unit.
* Represents monetary values as integers, in cents. This avoids floating point rounding errors.
* Represents currency as ```Money.Currency``` instances providing a high level of flexibility.

## Usage
### Init
```go
// Initialise Money by using smallest unit value (e.g 100 represents 1 pound)
// and use ISO 4217 Currency Code to set money Currency
pound := money.New(100, "GBP")
```
### Comparison

**Go-money** lets you to use base compare operations like:

* Equals
* GreaterThan
* GreaterThanOrEqual
* LessThan
* LessThanOrEqual
In order to use them currencies must be equal

```go
// Initialise Money by using smallest unit value (e.g 100 represents 1 pound)
// and use ISO 4217 Currency Code to set money Currency
pound := money.New(100, "GBP")
twoPounds := money.New(200, "GBP")
twoEuros := money.New(200, "EUR")

pound.GreaterThan(twoPounds) // false
pound.LessThan(twoPounds) // true
twoPounds.Equals(twoEuros) // Error: Currencies don't match
```











