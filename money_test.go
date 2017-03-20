package gocash

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New(500, "EUR")

	if m.Amount != 500 {
		t.Errorf("Expected %d got %d", 500, m.Amount)
	}

	if m.Currency.Code != "EUR" {
		t.Errorf("Expected currency %s got %s", "EUR", m.Currency.Code)
	}

	mm := New(-100, "EUR")

	if mm.Amount != -100 {
		t.Errorf("Expected %d got %d", -100, mm.Amount)
	}
}

func TestMoney_IsSameCurrency(t *testing.T) {
	me := New(0, "EUR")
	mu := New(0, "USD")

	if me.IsSameCurrency(mu) {
		t.Errorf("Expected %s not to be same as %s", me.Currency.Code, mu.Currency.Code)
	}

	if !me.IsSameCurrency(me) {
		t.Errorf("Expected %s to be same as %s", me.Currency.Code, mu.Currency.Code)
	}
}

func TestMoney_Equals(t *testing.T) {
	m := New(0, "EUR")
	om := New(0, "USD")

	if m.Equals(om) {
		t.Errorf("Expected %v not to be same as %v", m, om)
	}

	om = New(500, "EUR")

	if m.Equals(om) {
		t.Errorf("Expected %v not to be same as %v", m, om)
	}

	om = New(0, "EUR")

	if !m.Equals(om) {
		t.Errorf("Expected %v to be same as %v", m, om)
	}
}


