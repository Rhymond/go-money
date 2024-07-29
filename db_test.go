package money

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"
)

func TestMoney_Value(t *testing.T) {
	tests := []struct {
		have      *Money
		separator string
		want      string
		wantErr   bool
	}{
		{
			have:      New(10, CAD),
			separator: "|",
			want:      "10|CAD",
		},
		{
			have:      New(-10, USD),
			separator: "+-+",
			want:      "-10+-+USD",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.have), func(t *testing.T) {
			want := driver.Value(tt.want)
			DBMoneyValueSeparator = tt.separator
			got, err := tt.have.Value()
			if err != nil {
				t.Errorf("Value() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Value() got = %v, want %v", got, want)
			}
		})
	}
}

func TestMoney_Scan(t *testing.T) {
	tests := []struct {
		src       interface{}
		separator string
		want      *Money
		wantErr   bool
	}{
		{
			src:  "10|CAD",
			want: New(10, CAD),
		},
		{
			src:  "20|USD",
			want: New(20, USD),
		},
		{
			src:       "30000,IDR",
			separator: ",",
			want:      New(30000, IDR),
		},
		{
			src:     "10|",
			wantErr: true,
		},
		{
			src:     "|SAR",
			wantErr: true,
		},
		{
			src:     "10",
			wantErr: true,
		},
		{
			src:     "USD",
			wantErr: true,
		},
		{
			src:     "USD|10",
			wantErr: true,
		},
		{
			src:     "",
			wantErr: true,
		},
		{
			src:     "a|b|c",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.src), func(t *testing.T) {
			if tt.separator != "" {
				DBMoneyValueSeparator = tt.separator
			} else {
				DBMoneyValueSeparator = DefaultDBMoneyValueSeparator
			}
			got := &Money{}
			if err := got.Scan(tt.src); (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got == nil {
				t.Errorf("money.Scan() result was <nil>")
				return
			}
			eq, err := tt.want.Equals(got)
			if err != nil {
				t.Errorf(err.Error())
			}
			if !eq {
				t.Errorf("Value() got = %s %s, want %s %s", got.Display(), got.Currency().Code, tt.want.Display(), tt.want.Currency().Code)
			}
		})
	}
}

func TestCurrency_Value(t *testing.T) {
	for code, cc := range currencies {
		t.Run(code, func(t *testing.T) {
			want := driver.Value(code)

			got, err := cc.Value()
			if err != nil {
				t.Errorf("Value() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Value() got = %v, want %v", got, want)
			}
		})
	}
}

func TestCurrency_Scan(t *testing.T) {
	for code, want := range currencies {
		t.Run(code, func(t *testing.T) {
			src := interface{}(code)

			got := &Currency{}
			err := got.Scan(src)
			if err != nil {
				t.Errorf("Scan() error = %v", err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Scan() got %#v, want %#v", got, want)
			}
		})
	}
}
