package money

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

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
