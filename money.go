package money

import (
	"github.com/Rhymond/go-money/concrete"
	"github.com/Rhymond/go-money/interfaces"
)

// NewInt64 creates and returns new instance of Money
func NewInt64(amount interface{}, currencyCode string) (interfaces.Money, error) {
	return concrete.NewMoneyInt64(amount, currencyCode)
}
