package money

import (
	"math"
	"testing"
)

func TestGetCurrencyAmountWords(t *testing.T) {
	type args struct {
		amount      float64
		countryCode string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No decimal",
			args: args{amount: 123, countryCode: "PHP"},
			want: "one hundred twenty-three pesos only",
		},
		{
			name: "0 after decimal",
			args: args{amount: 123.00, countryCode: "PHP"},
			want: "one hundred twenty-three pesos only",
		},
		{
			name: "0 before decimal",
			args: args{amount: 00.11, countryCode: "PHP"},
			want: "eleven sentimos only",
		},
		{
			name: "decimal number",
			args: args{amount: 123.11, countryCode: "PHP"},
			want: "one hundred twenty-three pesos and eleven sentimos only",
		},
		{
			name: "decimal number",
			args: args{amount: 123.01, countryCode: "PHP"},
			want: "one hundred twenty-three pesos and one sentimos only",
		},
		{
			name: "decimal number",
			args: args{amount: 123.10, countryCode: "PHP"},
			want: "one hundred twenty-three pesos and one sentimos only",
		},

		{
			name: "No decimal",
			args: args{amount: 123, countryCode: "SGD"},
			want: "one hundred twenty-three dollar only",
		},
		{
			name: "0 after decimal",
			args: args{amount: 123.00, countryCode: "SGD"},
			want: "one hundred twenty-three dollar only",
		},
		{
			name: "0 before decimal",
			args: args{amount: 00.11, countryCode: "SGD"},
			want: "eleven cents only",
		},
		{
			name: "decimal number",
			args: args{amount: 123.11, countryCode: "SGD"},
			want: "one hundred twenty-three dollar and eleven cents only",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCurrencyAmountWords(tt.args.amount, tt.args.countryCode); got != tt.want {
				t.Errorf("GetCurrencyAmountWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrencyAmountWords_EdgeCases(t *testing.T) {
	// NaN triggers Atoi error on str[0] — returns fmt.Sprintf("%v", NaN) = "NaN"
	if got := GetCurrencyAmountWords(math.NaN(), "USD"); got != "NaN" {
		t.Errorf("NaN: got %q, want %q", got, "NaN")
	}

	// 1.5e10 triggers Atoi error on str[1] — returns fmt.Sprintf("%v", 1.5e10) = "1.5e+10"
	if got := GetCurrencyAmountWords(1.5e10, "USD"); got != "1.5e+10" {
		t.Errorf("1.5e10: got %q, want %q", got, "1.5e+10")
	}

	// Unknown currency code — returns fmt.Sprintf("%v", amount) = "1"
	if got := GetCurrencyAmountWords(1.0, "UNKNOWN_CCY"); got != "1" {
		t.Errorf("unknown currency: got %q, want %q", got, "1")
	}

	// Zero amount + JPY (no SubWord) — hits "zero %v only" branch
	if got := GetCurrencyAmountWords(0.0, "JPY"); got != "zero yen only" {
		t.Errorf("zero JPY: got %q, want %q", got, "zero yen only")
	}

	// Zero amount + USD — both parts zero, returns strAmount ("0")
	if got := GetCurrencyAmountWords(0.0, "USD"); got != "0" {
		t.Errorf("zero USD: got %q, want %q", got, "0")
	}
}

func TestGetCurrencyAmountWords_NoMainWord(t *testing.T) {
	// Register a currency with no main word
	AddCurrencyMeta("NOUNIT", "", "pieces")

	tcs := []struct {
		amount float64
		want   string
	}{
		// mainAmt != 0 && subAmt != 0 (1.51 → str[1]="51" → 51 pieces)
		{1.51, "one and fifty-one only"},
		// mainAmt != 0 && subAmt == 0 (1.00 → no decimal point in %+v → str[1]="00" → 0)
		{1.00, "one only"},
		// mainAmt == 0 && subAmt != 0 (0.05 → str[1]="05" → Atoi=5)
		{0.05, "point five only"},
	}
	for _, tc := range tcs {
		if got := GetCurrencyAmountWords(tc.amount, "NOUNIT"); got != tc.want {
			t.Errorf("NOUNIT(%.2f): got %q, want %q", tc.amount, got, tc.want)
		}
	}
}
