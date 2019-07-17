package money

import "testing"

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
			args: args{amount: 123, countryCode: "PH"},
			want: "one hundred twenty-three pesos only",
		},
		{
			name: "0 after decimal",
			args: args{amount: 123.00, countryCode: "PH"},
			want: "one hundred twenty-three pesos only",
		},
		{
			name: "0 before decimal",
			args: args{amount: 00.11, countryCode: "PH"},
			want: "eleven sentimos only",
		},
		{
			name: "decimal number",
			args: args{amount: 123.11, countryCode: "PH"},
			want: "one hundred twenty-three pesos and eleven sentimos only",
		},
		{
			name: "decimal number",
			args: args{amount: 123.01, countryCode: "PH"},
			want: "one hundred twenty-three pesos and one sentimos only",
		},
		{
			name: "decimal number",
			args: args{amount: 123.10, countryCode: "PH"},
			want: "one hundred twenty-three pesos and one sentimos only",
		},

		{
			name: "No decimal",
			args: args{amount: 123, countryCode: "SG"},
			want: "one hundred twenty-three dollar only",
		},
		{
			name: "0 after decimal",
			args: args{amount: 123.00, countryCode: "SG"},
			want: "one hundred twenty-three dollar only",
		},
		{
			name: "0 before decimal",
			args: args{amount: 00.11, countryCode: "SG"},
			want: "eleven cents only",
		},
		{
			name: "decimal number",
			args: args{amount: 123.11, countryCode: "SG"},
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