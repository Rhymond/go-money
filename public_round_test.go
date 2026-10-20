package money

import (
	"math"
	"testing"
)

func TestPublicSignedRound(t *testing.T) {
	for _, c := range []struct{ in, want int64 }{{math.MinInt64, -9223372036854775800}, {-151, -200}, {-150, -100}, {-149, -100}, {149, 100}, {150, 100}, {151, 200}, {0, 0}} {
		m := New(c.in, USD)
		got := m.Round()
		if got.Amount() != c.want {
			t.Errorf("Round(%d)=%d want %d", c.in, got.Amount(), c.want)
		}
		if m.Amount() != c.in {
			t.Error("input mutated")
		}
		if got.Currency().Code != USD {
			t.Error("currency changed")
		}
	}
}
