package money

import (
	"reflect"
	"regexp"
	"testing"
)

func Test_patternFactory(t *testing.T) {
	twoDigitsPattern := `^\d{1,3}(,\d{3})*(\.\d{2})?$`
	regex := regexp.MustCompile(twoDigitsPattern)
	type args struct {
		fractionNumber int
	}
	tests := []struct {
		name    string
		args    args
		want    *regexp.Regexp
		wantErr bool
	}{
		{
			name:    "wrong fractional number",
			args:    args{fractionNumber: 5},
			wantErr: true,
		},
		{
			name: "success fractional number",
			args: args{fractionNumber: 2},
			want: regex,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := patternFactory(tt.args.fractionNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("patternFactory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("patternFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_amountIsValid(t *testing.T) {
	type args struct {
		amount   string
		fraction int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "invalid amount - First - 2 fractions",
			args: args{
				amount:   "1,0000.32",
				fraction: 2,
			},
			want: false,
		},
		{
			name: "invalid amount - Second - 2 fractions",
			args: args{
				amount:   "1,000,00.320",
				fraction: 2,
			},
			want: false,
		},
		{
			name: "invalid amount - #3 - 2 fractions",
			args: args{
				amount:   "1,000,0.32",
				fraction: 2,
			},
			want: false,
		},
		{
			name: "success amount - 2 fractions",
			args: args{
				amount:   "1,000,000.32",
				fraction: 2,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := amountIsValid(tt.args.amount, tt.args.fraction)
			if (err != nil) != tt.wantErr {
				t.Errorf("amountIsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("amountIsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanString(t *testing.T) {
	assertAmountClean := "1234423677.33"

	type args struct {
		validAmount string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid amount to clean",
			args: args{
				validAmount: "1,234,423,677.33",
			},
			want: assertAmountClean,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanString(tt.args.validAmount); got != tt.want {
				t.Errorf("cleanString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseStringToFloat64(t *testing.T) {
	type args struct {
		validAmountFormated string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "valid formated amount",
			args: args{
				validAmountFormated: "123456789.900",
			},
			want: float64(123456789.900),
		},
		{
			name: "invalid formated amount",
			args: args{
				validAmountFormated: "helloworld123456789.900",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStringToFloat64(tt.args.validAmountFormated)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseStringToFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseStringToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertStringAmount(t *testing.T) {
	usdCurrency := GetCurrency(USD)
	mockAmount := NewFromFloat(10000000.99, USD)
	type args struct {
		amount   string
		currency *Currency
	}
	tests := []struct {
		name    string
		args    args
		want    *Money
		wantErr bool
	}{
		{
			name: "success amount - USD",
			args: args{
				amount:   "10,000,000.99",
				currency: usdCurrency,
			},
			want: mockAmount,
		},
		{
			name: "invalid amount - USD",
			args: args{
				amount:   "10,000,000.9999",
				currency: usdCurrency,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertStringAmount(tt.args.amount, tt.args.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertStringAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertStringAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}
