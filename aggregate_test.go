package money

import "testing"

func TestSum(t *testing.T) {
	tcs := []struct {
		amounts  []int64
		expected int64
	}{
		{[]int64{100, 200, 300}, 600},
		{[]int64{100}, 100},
		{[]int64{-100, 200}, 100},
		{[]int64{0, 0, 0}, 0},
	}

	for _, tc := range tcs {
		var ms []*Money
		for _, a := range tc.amounts {
			ms = append(ms, New(a, "USD"))
		}
		result, err := Sum(ms...)
		if err != nil {
			t.Fatalf("Sum: unexpected error: %v", err)
		}
		if result.Amount() != tc.expected {
			t.Errorf("Sum(%v): expected %d got %d", tc.amounts, tc.expected, result.Amount())
		}
	}
}

func TestSum_Empty(t *testing.T) {
	_, err := Sum()
	if err == nil {
		t.Error("Sum(): expected error for empty input")
	}
}

func TestSum_CurrencyMismatch(t *testing.T) {
	_, err := Sum(New(100, "USD"), New(100, "EUR"))
	if err == nil {
		t.Error("Sum(): expected error for currency mismatch")
	}
}

func TestMin(t *testing.T) {
	tcs := []struct {
		amounts  []int64
		expected int64
	}{
		{[]int64{100, 200, 50}, 50},
		{[]int64{100}, 100},
		{[]int64{-10, -100, -1}, -100},
		{[]int64{300, 300}, 300},
	}

	for _, tc := range tcs {
		var ms []*Money
		for _, a := range tc.amounts {
			ms = append(ms, New(a, "USD"))
		}
		result, err := Min(ms...)
		if err != nil {
			t.Fatalf("Min: unexpected error: %v", err)
		}
		if result.Amount() != tc.expected {
			t.Errorf("Min(%v): expected %d got %d", tc.amounts, tc.expected, result.Amount())
		}
	}
}

func TestMin_Empty(t *testing.T) {
	_, err := Min()
	if err == nil {
		t.Error("Min(): expected error for empty input")
	}
}

func TestMax(t *testing.T) {
	tcs := []struct {
		amounts  []int64
		expected int64
	}{
		{[]int64{100, 200, 50}, 200},
		{[]int64{100}, 100},
		{[]int64{-10, -100, -1}, -1},
		{[]int64{300, 300}, 300},
	}

	for _, tc := range tcs {
		var ms []*Money
		for _, a := range tc.amounts {
			ms = append(ms, New(a, "USD"))
		}
		result, err := Max(ms...)
		if err != nil {
			t.Fatalf("Max: unexpected error: %v", err)
		}
		if result.Amount() != tc.expected {
			t.Errorf("Max(%v): expected %d got %d", tc.amounts, tc.expected, result.Amount())
		}
	}
}

func TestMax_Empty(t *testing.T) {
	_, err := Max()
	if err == nil {
		t.Error("Max(): expected error for empty input")
	}
}

func TestAverage(t *testing.T) {
	tcs := []struct {
		amounts  []int64
		expected int64
	}{
		{[]int64{100, 200, 300}, 200},
		{[]int64{100, 101}, 100}, // truncated
		{[]int64{100}, 100},
	}

	for _, tc := range tcs {
		var ms []*Money
		for _, a := range tc.amounts {
			ms = append(ms, New(a, "USD"))
		}
		result, err := Average(ms...)
		if err != nil {
			t.Fatalf("Average: unexpected error: %v", err)
		}
		if result.Amount() != tc.expected {
			t.Errorf("Average(%v): expected %d got %d", tc.amounts, tc.expected, result.Amount())
		}
	}
}

func TestAverage_Empty(t *testing.T) {
	_, err := Average()
	if err == nil {
		t.Error("Average(): expected error for empty input")
	}
}
