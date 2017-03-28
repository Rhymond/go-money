// Copyright © 2017 Raimondas Kazlauskas
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package money

import (
	"log"
)

type Money struct {
	Number   *Number
	Currency *Currency
}

var calc *Calculator

func New(am int, curr string) *Money {

	calc = new(Calculator)

	return &Money{
		&Number{am},
		new(Currency).Get(curr),
	}
}

func (m *Money) SameCurrency(om *Money) bool {
	return m.Currency.Equals(om.Currency)
}

func (m *Money) assertSameCurrency(om *Money) {
	if !m.SameCurrency(om) {
		log.Fatalf("Currencies %s and %s don't match", m.Currency.Code, om.Currency.Code)
	}
}

func (m *Money) compare(om *Money) int {

	m.assertSameCurrency(om)

	switch {
	case m.Number.Amount > om.Number.Amount:
		return 1
	case m.Number.Amount < om.Number.Amount:
		return -1
	}

	return 0
}

func (m *Money) Equals(om *Money) bool {
	return m.compare(om) == 0
}

func (m *Money) GreaterThan(om *Money) bool {
	return m.compare(om) == 1
}

func (m *Money) GreaterThanOrEqual(om *Money) bool {
	return m.compare(om) >= 0
}

func (m *Money) LessThan(om *Money) bool {
	return m.compare(om) == -1
}

func (m *Money) LessThanOrEqual(om *Money) bool {
	return m.compare(om) <= 0
}

func (m *Money) IsZero() bool {
	return m.Number.Amount == 0
}

func (m *Money) IsPositive() bool {
	return m.Number.Amount > 0
}

func (m *Money) IsNegative() bool {
	return m.Number.Amount < 0
}

func (m *Money) Absolute() *Money {
	m.Number.Absolute()
	return m
}

func (m *Money) Negative() *Money {
	m.Number.Negative()
	return m
}

func (m *Money) Add(om *Money) *Money {
	m.assertSameCurrency(om)
	return &Money{calc.add(m.Number, om.Number), m.Currency}
}

func (m *Money) Subtract(om *Money) *Money {
	m.assertSameCurrency(om)
	return &Money{calc.subtract(m.Number, om.Number), m.Currency}
}

func (m *Money) Multiply(mul int) *Money {
	return &Money{calc.multiply(m.Number, mul), m.Currency}
}

func (m *Money) Divide(div int) *Money {
	return &Money{calc.divide(m.Number, div), m.Currency}
}

func (m *Money) Round() *Money {
	return &Money{calc.round(m.Number), m.Currency}
}

func (m *Money) Split(n int) []*Money {
	if n <= 0 {
		log.Fatalf("Split must be higher than zero")
	}

	a := calc.divide(m.Number, n)
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{a, m.Currency}
	}

	l := calc.modulus(m.Number, n).Amount

	// Add leftovers to the first parties
	for p := 0; l != 0; p++ {
		ms[p].Number = calc.add(ms[p].Number, &Number{1})
		l -= 1
	}

	return ms
}

func (m *Money) Allocate(rs []int) []*Money {
	if len(rs) == 0 {
		log.Fatalf("No ratios specified")
	}

	// Calculate sum of ratios
	var sum int
	for _, r := range rs {
		sum += r
	}

	var total int
	var ms []*Money
	for _, r := range rs {
		party := &Money{
			calc.allocate(m.Number, r, sum),
			m.Currency,
		}

		ms = append(ms, party)
		total += party.Number.Amount
	}

	// Calculate leftover value and divide to first parties
	lo := m.Number.Amount - total
	sub := 1
	if lo < 0 {
		sub = -1
	}

	for p := 0; lo != 0; p++ {
		ms[p].Number = calc.add(ms[p].Number, &Number{sub})
		lo -= sub
	}

	return ms
}
