package money

import "testing"

func TestRoundWithMode(t *testing.T) {
	tcs := []struct {
		name     string
		amount   int64
		mode     RoundingMode
		expected int64
	}{
		// RoundHalfUp — ties round away from zero
		{"HalfUp_below_half", 124, RoundHalfUp, 100},
		{"HalfUp_exactly_half", 150, RoundHalfUp, 200},
		{"HalfUp_above_half", 175, RoundHalfUp, 200},
		{"HalfUp_negative_half", -150, RoundHalfUp, -200},
		{"HalfUp_negative_below", -124, RoundHalfUp, -100},

		// RoundHalfDown — ties round toward zero
		{"HalfDown_below_half", 124, RoundHalfDown, 100},
		{"HalfDown_exactly_half", 150, RoundHalfDown, 100},
		{"HalfDown_above_half", 175, RoundHalfDown, 200},
		{"HalfDown_negative_half", -150, RoundHalfDown, -100},

		// RoundHalfEven (banker's) — ties round to nearest even
		{"HalfEven_tie_0.5_to_0", 50, RoundHalfEven, 0},   // 0 is even
		{"HalfEven_tie_1.5_to_2", 150, RoundHalfEven, 200}, // 2 is even
		{"HalfEven_tie_2.5_to_2", 250, RoundHalfEven, 200}, // 2 is even
		{"HalfEven_tie_3.5_to_4", 350, RoundHalfEven, 400}, // 4 is even
		{"HalfEven_above_half", 175, RoundHalfEven, 200},
		{"HalfEven_below_half", 124, RoundHalfEven, 100},

		// RoundUp — always away from zero
		{"RoundUp_small", 101, RoundUp, 200},
		{"RoundUp_half", 150, RoundUp, 200},
		{"RoundUp_exact", 100, RoundUp, 100},
		{"RoundUp_negative", -101, RoundUp, -200},

		// RoundDown — always truncate toward zero
		{"RoundDown_large", 199, RoundDown, 100},
		{"RoundDown_half", 150, RoundDown, 100},
		{"RoundDown_exact", 100, RoundDown, 100},
		{"RoundDown_negative", -199, RoundDown, -100},

		// Zero
		{"Zero_HalfUp", 0, RoundHalfUp, 0},
		{"Zero_HalfEven", 0, RoundHalfEven, 0},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			m := New(tc.amount, "EUR") // EUR has Fraction=2
			r := m.RoundWithMode(tc.mode)
			if r.Amount() != tc.expected {
				t.Errorf("RoundWithMode(%d, %v): expected %d got %d",
					tc.amount, tc.mode, tc.expected, r.Amount())
			}
		})
	}
}

func TestRoundWithMode_NoFraction(t *testing.T) {
	// JPY has Fraction=0 — rounding should be a no-op
	m := New(100, "JPY")
	for _, mode := range []RoundingMode{RoundHalfUp, RoundHalfDown, RoundHalfEven, RoundUp, RoundDown} {
		r := m.RoundWithMode(mode)
		if r.Amount() != 100 {
			t.Errorf("RoundWithMode on zero-fraction currency: expected 100 got %d", r.Amount())
		}
	}
}
