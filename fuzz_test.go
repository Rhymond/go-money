package money

import "testing"

// FuzzMoneyAdd verifies that addition is commutative and consistent.
func FuzzMoneyAdd(f *testing.F) {
	f.Add(int64(100), int64(200))
	f.Add(int64(0), int64(0))
	f.Add(int64(-100), int64(200))
	f.Add(int64(9223372036854775806), int64(1))

	f.Fuzz(func(t *testing.T, a, b int64) {
		m1 := New(a, USD)
		m2 := New(b, USD)

		ab, err := m1.Add(m2)
		if err != nil {
			t.Skip()
		}
		ba, err := m2.Add(m1)
		if err != nil {
			t.Skip()
		}
		// Commutativity: a+b == b+a
		if ab.Amount() != ba.Amount() {
			t.Errorf("Add not commutative: %d+%d=%d but %d+%d=%d",
				a, b, ab.Amount(), b, a, ba.Amount())
		}
		// Identity: (a+b) - b == a
		result, err := ab.Subtract(m2)
		if err != nil {
			t.Skip()
		}
		if result.Amount() != a {
			t.Errorf("(a+b)-b != a: a=%d b=%d got %d", a, b, result.Amount())
		}
	})
}

// FuzzMoneyAllocate checks that allocations never lose money.
func FuzzMoneyAllocate(f *testing.F) {
	f.Add(int64(100), 33, 33, 34)
	f.Add(int64(1000), 50, 30, 20)
	f.Add(int64(1), 1, 1, 1)

	f.Fuzz(func(t *testing.T, amount int64, r1, r2, r3 int) {
		// Ensure ratios are positive to avoid allocation errors.
		if r1 <= 0 {
			r1 = 1
		}
		if r2 <= 0 {
			r2 = 1
		}
		if r3 <= 0 {
			r3 = 1
		}
		m := New(amount, USD)
		parts, err := m.Allocate(r1, r2, r3)
		if err != nil {
			t.Skip()
		}
		var total int64
		for _, p := range parts {
			total += p.Amount()
		}
		if total != amount {
			t.Errorf("Allocate lost money: input=%d total=%d", amount, total)
		}
	})
}

// FuzzMoneySplit checks that split never loses money.
func FuzzMoneySplit(f *testing.F) {
	f.Add(int64(100), 3)
	f.Add(int64(1), 1)
	f.Add(int64(1000), 7)

	f.Fuzz(func(t *testing.T, amount int64, n int) {
		if n <= 0 || n > 1000 {
			t.Skip()
		}
		m := New(amount, USD)
		parts, err := m.Split(n)
		if err != nil {
			t.Skip()
		}
		var total int64
		for _, p := range parts {
			total += p.Amount()
		}
		if total != amount {
			t.Errorf("Split lost money: input=%d total=%d", amount, total)
		}
	})
}

// FuzzMoneyScan verifies that scan never panics on arbitrary input.
func FuzzMoneyScan(f *testing.F) {
	f.Add("100 USD")
	f.Add("0 EUR")
	f.Add("-500 GBP")

	f.Fuzz(func(t *testing.T, s string) {
		var m Money
		// Must not panic regardless of input.
		m.Scan(s) //nolint:errcheck
	})
}
