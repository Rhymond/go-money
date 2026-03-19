package money

import "testing"

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(12345, USD)
	}
}

func BenchmarkNewFromFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewFromFloat(123.45, USD)
	}
}

func BenchmarkMoney_Add(b *testing.B) {
	a := New(100, USD)
	c := New(200, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Add(c) //nolint:errcheck
	}
}

func BenchmarkMoney_Subtract(b *testing.B) {
	a := New(300, USD)
	c := New(100, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Subtract(c) //nolint:errcheck
	}
}

func BenchmarkMoney_Multiply(b *testing.B) {
	m := New(100, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Multiply(3)
	}
}

func BenchmarkMoney_MultiplyFloat(b *testing.B) {
	m := New(1000, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.MultiplyFloat(1.21)
	}
}

func BenchmarkMoney_Divide(b *testing.B) {
	m := New(1000, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Divide(3)
	}
}

func BenchmarkMoney_DivideWithRounding(b *testing.B) {
	m := New(1000, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.DivideWithRounding(3, RoundHalfEven)
	}
}

func BenchmarkMoney_Percentage(b *testing.B) {
	m := New(10000, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Percentage(8.5)
	}
}

func BenchmarkMoney_Split(b *testing.B) {
	m := New(10085, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Split(3) //nolint:errcheck
	}
}

func BenchmarkMoney_Allocate(b *testing.B) {
	m := New(10085, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Allocate(60, 30, 10) //nolint:errcheck
	}
}

func BenchmarkMoney_Display(b *testing.B) {
	m := New(123456789, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Display()
	}
}

func BenchmarkMoney_LocaleFormat(b *testing.B) {
	m := New(123456789, EUR)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.LocaleFormat("de-DE")
	}
}

func BenchmarkMoney_MarshalJSON(b *testing.B) {
	m := New(1234, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.MarshalJSON() //nolint:errcheck
	}
}

func BenchmarkMoney_MarshalBinary(b *testing.B) {
	m := New(1234, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.MarshalBinary() //nolint:errcheck
	}
}

func BenchmarkExchangeRate_Convert(b *testing.B) {
	rate, _ := NewExchangeRate("USD", "EUR", 0.92)
	m := New(1000, USD)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rate.Convert(m) //nolint:errcheck
	}
}

func BenchmarkSum(b *testing.B) {
	items := []*Money{
		New(100, USD), New(200, USD), New(300, USD),
		New(400, USD), New(500, USD),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sum(items...) //nolint:errcheck
	}
}
