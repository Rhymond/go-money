package money

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"
)

func TestMoney_Value(t *testing.T) {
	tests := []struct {
		have    *Money
		want    []byte
		wantErr bool
	}{
		{
			have: New(10, CAD),
			want: []byte(`{"amount":10,"currency":"CAD"}`),
		},
		{
			have: New(-10, USD),
			want: []byte(`{"amount":-10,"currency":"USD"}`),
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.have), func(t *testing.T) {
			want := driver.Value(tt.want)
			got, err := tt.have.Value()
			if err != nil {
				t.Errorf("Value() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Value() got = %v, want %v", string((got).([]byte)), string((want).([]byte)))
			}
		})
	}
}

func TestMoney_Scan(t *testing.T) {
	tests := []struct {
		name    string
		src     interface{}
		want    *Money
		wantErr bool
	}{
		{
			name: "10 cad",
			src:  []byte(`{"amount":10,"currency":"CAD"}`),
			want: New(10, CAD),
		},
		{
			name: "20 usd",
			src:  []byte(`{"amount":20,"currency":"USD"}`),
			want: New(20, USD),
		},
		{
			name: "30.00 IDR",
			src:  []byte(`{"amount":30000,"currency":"IDR"}`),
			want: New(30000, IDR),
		},
		{
			name:    "empty currency",
			src:     []byte(`{"amount":10,"currency":""}`),
			wantErr: true,
		},
		{
			name:    "missing amount, parse error",
			src:     []byte(`{"amount":,"currency":"SAR"}`),
			wantErr: true,
		},
		{
			name:    "missing currency",
			src:     []byte(`{"amount":10}`),
			wantErr: true,
		},
		{
			name: "missing amount",
			src:  []byte(`{"currency":"USD"}`),
			// unfortunately, this is not an error
			want: New(0, USD),
		},
		{
			name:    "empty",
			src:     "{}",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Money{}
			if err := got.Scan(tt.src); (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got == nil {
				t.Errorf("money.Scan() result was <nil> %v", tt.want)
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
