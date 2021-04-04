package money

import (
	"testing"
)

func TestFormatter_Format(t *testing.T) {
	tcs := []struct {
		fraction int
		decimal  string
		thousand string
		grapheme string
		template string
		amount   int64
		expected string
	}{
		{2, ".", ",", "$", "1 $", 0, "0.00 $"},
		{2, ".", ",", "$", "1 $", 1, "0.01 $"},
		{2, ".", ",", "$", "1 $", 12, "0.12 $"},
		{2, ".", ",", "$", "1 $", 123, "1.23 $"},
		{2, ".", ",", "$", "1 $", 1234, "12.34 $"},
		{2, ".", ",", "$", "1 $", 12345, "123.45 $"},
		{2, ".", ",", "$", "1 $", 123456, "1,234.56 $"},
		{2, ".", ",", "$", "1 $", 1234567, "12,345.67 $"},
		{2, ".", ",", "$", "1 $", 12345678, "123,456.78 $"},
		{2, ".", ",", "$", "1 $", 123456789, "1,234,567.89 $"},

		{2, ".", ",", "$", "1 $", -1, "-0.01 $"},
		{2, ".", ",", "$", "1 $", -12, "-0.12 $"},
		{2, ".", ",", "$", "1 $", -123, "-1.23 $"},
		{2, ".", ",", "$", "1 $", -1234, "-12.34 $"},
		{2, ".", ",", "$", "1 $", -12345, "-123.45 $"},
		{2, ".", ",", "$", "1 $", -123456, "-1,234.56 $"},
		{2, ".", ",", "$", "1 $", -1234567, "-12,345.67 $"},
		{2, ".", ",", "$", "1 $", -12345678, "-123,456.78 $"},
		{2, ".", ",", "$", "1 $", -123456789, "-1,234,567.89 $"},

		{3, ".", "", "$", "1 $", 1, "0.001 $"},
		{3, ".", "", "$", "1 $", 12, "0.012 $"},
		{3, ".", "", "$", "1 $", 123, "0.123 $"},
		{3, ".", "", "$", "1 $", 1234, "1.234 $"},
		{3, ".", "", "$", "1 $", 12345, "12.345 $"},
		{3, ".", "", "$", "1 $", 123456, "123.456 $"},
		{3, ".", "", "$", "1 $", 1234567, "1234.567 $"},
		{3, ".", "", "$", "1 $", 12345678, "12345.678 $"},
		{3, ".", "", "$", "1 $", 123456789, "123456.789 $"},

		{2, ".", ",", "£", "$1", 1, "£0.01"},
		{2, ".", ",", "£", "$1", 12, "£0.12"},
		{2, ".", ",", "£", "$1", 123, "£1.23"},
		{2, ".", ",", "£", "$1", 1234, "£12.34"},
		{2, ".", ",", "£", "$1", 12345, "£123.45"},
		{2, ".", ",", "£", "$1", 123456, "£1,234.56"},
		{2, ".", ",", "£", "$1", 1234567, "£12,345.67"},
		{2, ".", ",", "£", "$1", 12345678, "£123,456.78"},
		{2, ".", ",", "£", "$1", 123456789, "£1,234,567.89"},

		{0, ".", ",", "NT$", "$1", 1, "NT$1"},
		{0, ".", ",", "NT$", "$1", 12, "NT$12"},
		{0, ".", ",", "NT$", "$1", 123, "NT$123"},
		{0, ".", ",", "NT$", "$1", 1234, "NT$1,234"},
		{0, ".", ",", "NT$", "$1", 12345, "NT$12,345"},
		{0, ".", ",", "NT$", "$1", 123456, "NT$123,456"},
		{0, ".", ",", "NT$", "$1", 1234567, "NT$1,234,567"},
		{0, ".", ",", "NT$", "$1", 12345678, "NT$12,345,678"},
		{0, ".", ",", "NT$", "$1", 123456789, "NT$123,456,789"},

		{0, ".", ",", "NT$", "$1", -1, "-NT$1"},
		{0, ".", ",", "NT$", "$1", -12, "-NT$12"},
		{0, ".", ",", "NT$", "$1", -123, "-NT$123"},
		{0, ".", ",", "NT$", "$1", -1234, "-NT$1,234"},
		{0, ".", ",", "NT$", "$1", -12345, "-NT$12,345"},
		{0, ".", ",", "NT$", "$1", -123456, "-NT$123,456"},
		{0, ".", ",", "NT$", "$1", -1234567, "-NT$1,234,567"},
		{0, ".", ",", "NT$", "$1", -12345678, "-NT$12,345,678"},
		{0, ".", ",", "NT$", "$1", -123456789, "-NT$123,456,789"},
	}

	for _, tc := range tcs {
		formatter := NewFormatter(tc.fraction, tc.decimal, tc.thousand, tc.grapheme, tc.template, "")
		r := formatter.Format(tc.amount)

		if r != tc.expected {
			t.Errorf("Expected %d formatted to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}

func TestFormatter_FormatWithCurrencyCode(t *testing.T) {
	tcs := []struct {
		fraction int
		decimal  string
		thousand string
		grapheme string
		code     string
		template string
		amount   int64
		expected string
	}{
		{2, ".", ",", "$", "USD", "1 $", 0, "0.00 USD"},
		{2, ".", ",", "$", "USD", "1 $", 1, "0.01 USD"},
		{2, ".", ",", "$", "USD", "1 $", 12, "0.12 USD"},
		{2, ".", ",", "$", "USD", "1 $", 123, "1.23 USD"},
		{2, ".", ",", "$", "USD", "1 $", 1234, "12.34 USD"},
		{2, ".", ",", "$", "USD", "1 $", 12345, "123.45 USD"},
		{2, ".", ",", "$", "USD", "1 $", 123456, "1,234.56 USD"},
		{2, ".", ",", "$", "USD", "1 $", 1234567, "12,345.67 USD"},
		{2, ".", ",", "$", "USD", "1 $", 12345678, "123,456.78 USD"},
		{2, ".", ",", "$", "USD", "1 $", 123456789, "1,234,567.89 USD"},

		{2, ".", ",", "$", "USD", "1 $", -1, "-0.01 USD"},
		{2, ".", ",", "$", "USD", "1 $", -12, "-0.12 USD"},
		{2, ".", ",", "$", "USD", "1 $", -123, "-1.23 USD"},
		{2, ".", ",", "$", "USD", "1 $", -1234, "-12.34 USD"},
		{2, ".", ",", "$", "USD", "1 $", -12345, "-123.45 USD"},
		{2, ".", ",", "$", "USD", "1 $", -123456, "-1,234.56 USD"},
		{2, ".", ",", "$", "USD", "1 $", -1234567, "-12,345.67 USD"},
		{2, ".", ",", "$", "USD", "1 $", -12345678, "-123,456.78 USD"},
		{2, ".", ",", "$", "USD", "1 $", -123456789, "-1,234,567.89 USD"},

		{3, ".", "", "$", "USD", "1 $", 1, "0.001 USD"},
		{3, ".", "", "$", "USD", "1 $", 12, "0.012 USD"},
		{3, ".", "", "$", "USD", "1 $", 123, "0.123 USD"},
		{3, ".", "", "$", "USD", "1 $", 1234, "1.234 USD"},
		{3, ".", "", "$", "USD", "1 $", 12345, "12.345 USD"},
		{3, ".", "", "$", "USD", "1 $", 123456, "123.456 USD"},
		{3, ".", "", "$", "USD", "1 $", 1234567, "1234.567 USD"},
		{3, ".", "", "$", "USD", "1 $", 12345678, "12345.678 USD"},
		{3, ".", "", "$", "USD", "1 $", 123456789, "123456.789 USD"},

		{2, ".", ",", "£", "GBP", "$1", 1, "GBP0.01"},
		{2, ".", ",", "£", "GBP", "$1", 12, "GBP0.12"},
		{2, ".", ",", "£", "GBP", "$1", 123, "GBP1.23"},
		{2, ".", ",", "£", "GBP", "$1", 1234, "GBP12.34"},
		{2, ".", ",", "£", "GBP", "$1", 12345, "GBP123.45"},
		{2, ".", ",", "£", "GBP", "$1", 123456, "GBP1,234.56"},
		{2, ".", ",", "£", "GBP", "$1", 1234567, "GBP12,345.67"},
		{2, ".", ",", "£", "GBP", "$1", 12345678, "GBP123,456.78"},
		{2, ".", ",", "£", "GBP", "$1", 123456789, "GBP1,234,567.89"},

		{0, ".", ",", "NT$", "TWD", "$1", 1, "TWD1"},
		{0, ".", ",", "NT$", "TWD", "$1", 12, "TWD12"},
		{0, ".", ",", "NT$", "TWD", "$1", 123, "TWD123"},
		{0, ".", ",", "NT$", "TWD", "$1", 1234, "TWD1,234"},
		{0, ".", ",", "NT$", "TWD", "$1", 12345, "TWD12,345"},
		{0, ".", ",", "NT$", "TWD", "$1", 123456, "TWD123,456"},
		{0, ".", ",", "NT$", "TWD", "$1", 1234567, "TWD1,234,567"},
		{0, ".", ",", "NT$", "TWD", "$1", 12345678, "TWD12,345,678"},
		{0, ".", ",", "NT$", "TWD", "$1", 123456789, "TWD123,456,789"},

		{0, ".", ",", "NT$", "TWD", "$1", -1, "-TWD1"},
		{0, ".", ",", "NT$", "TWD", "$1", -12, "-TWD12"},
		{0, ".", ",", "NT$", "TWD", "$1", -123, "-TWD123"},
		{0, ".", ",", "NT$", "TWD", "$1", -1234, "-TWD1,234"},
		{0, ".", ",", "NT$", "TWD", "$1", -12345, "-TWD12,345"},
		{0, ".", ",", "NT$", "TWD", "$1", -123456, "-TWD123,456"},
		{0, ".", ",", "NT$", "TWD", "$1", -1234567, "-TWD1,234,567"},
		{0, ".", ",", "NT$", "TWD", "$1", -12345678, "-TWD12,345,678"},
		{0, ".", ",", "NT$", "TWD", "$1", -123456789, "-TWD123,456,789"},
	}

	for _, tc := range tcs {
		formatter := NewFormatter(tc.fraction, tc.decimal, tc.thousand, tc.grapheme, tc.template, tc.code)
		r := formatter.FormatWithCurrencyCode(tc.amount)

		if r != tc.expected {
			t.Errorf("Expected %d formatted to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}

func TestFormatter_ToMajorUnits(t *testing.T) {
	tcs := []struct {
		fraction int
		decimal  string
		thousand string
		grapheme string
		template string
		amount   int64
		expected float64
	}{
		{2, ".", ",", "$", "1 $", 0, 0.00},
		{2, ".", ",", "$", "1 $", 1, 0.01},
		{2, ".", ",", "$", "1 $", 12, 0.12},
		{2, ".", ",", "$", "1 $", 123, 1.23},
		{2, ".", ",", "$", "1 $", 1234, 12.34},
		{2, ".", ",", "$", "1 $", 12345, 123.45},
		{2, ".", ",", "$", "1 $", 123456, 1234.56},
		{2, ".", ",", "$", "1 $", 1234567, 12345.67},
		{2, ".", ",", "$", "1 $", 12345678, 123456.78},
		{2, ".", ",", "$", "1 $", 123456789, 1234567.89},

		{2, ".", ",", "$", "1 $", -1, -0.01},
		{2, ".", ",", "$", "1 $", -12, -0.12},
		{2, ".", ",", "$", "1 $", -123, -1.23},
		{2, ".", ",", "$", "1 $", -1234, -12.34},
		{2, ".", ",", "$", "1 $", -12345, -123.45},
		{2, ".", ",", "$", "1 $", -123456, -1234.56},
		{2, ".", ",", "$", "1 $", -1234567, -12345.67},
		{2, ".", ",", "$", "1 $", -12345678, -123456.78},
		{2, ".", ",", "$", "1 $", -123456789, -1234567.89},

		{3, ".", "", "$", "1 $", 1, 0.001},
		{3, ".", "", "$", "1 $", 12, 0.012},
		{3, ".", "", "$", "1 $", 123, 0.123},
		{3, ".", "", "$", "1 $", 1234, 1.234},
		{3, ".", "", "$", "1 $", 12345, 12.345},
		{3, ".", "", "$", "1 $", 123456, 123.456},
		{3, ".", "", "$", "1 $", 1234567, 1234.567},
		{3, ".", "", "$", "1 $", 12345678, 12345.678},
		{3, ".", "", "$", "1 $", 123456789, 123456.789},

		{2, ".", ",", "£", "$1", 1, 0.01},
		{2, ".", ",", "£", "$1", 12, 0.12},
		{2, ".", ",", "£", "$1", 123, 1.23},
		{2, ".", ",", "£", "$1", 1234, 12.34},
		{2, ".", ",", "£", "$1", 12345, 123.45},
		{2, ".", ",", "£", "$1", 123456, 1234.56},
		{2, ".", ",", "£", "$1", 1234567, 12345.67},
		{2, ".", ",", "£", "$1", 12345678, 123456.78},
		{2, ".", ",", "£", "$1", 123456789, 1234567.89},

		{0, ".", ",", "NT$", "$1", 1, 1},
		{0, ".", ",", "NT$", "$1", 12, 12},
		{0, ".", ",", "NT$", "$1", 123, 123},
		{0, ".", ",", "NT$", "$1", 1234, 1234},
		{0, ".", ",", "NT$", "$1", 12345, 12345},
		{0, ".", ",", "NT$", "$1", 123456, 123456},
		{0, ".", ",", "NT$", "$1", 1234567, 1234567},
		{0, ".", ",", "NT$", "$1", 12345678, 12345678},
		{0, ".", ",", "NT$", "$1", 123456789, 123456789},

		{0, ".", ",", "NT$", "$1", -1, -1},
		{0, ".", ",", "NT$", "$1", -12, -12},
		{0, ".", ",", "NT$", "$1", -123, -123},
		{0, ".", ",", "NT$", "$1", -1234, -1234},
		{0, ".", ",", "NT$", "$1", -12345, -12345},
		{0, ".", ",", "NT$", "$1", -123456, -123456},
		{0, ".", ",", "NT$", "$1", -1234567, -1234567},
		{0, ".", ",", "NT$", "$1", -12345678, -12345678},
		{0, ".", ",", "NT$", "$1", -123456789, -123456789},
	}

	for _, tc := range tcs {
		formatter := NewFormatter(tc.fraction, tc.decimal, tc.thousand, tc.grapheme, tc.template, "")
		r := formatter.ToMajorUnits(tc.amount)

		if r != tc.expected {
			t.Errorf("Expected %d formatted to major units to be %f got %f", tc.amount, tc.expected, r)
		}
	}
}

func Test_parseFormattedString(t *testing.T) {
	tests := []struct {
		name            string
		formattedString string
		currency        *Currency
		want            int64
		wantErr         bool
	}{
		{
			name:            "major amount and zero decimal amount",
			formattedString: "1.00",
			currency:        currencies["USD"],
			want:            100,
			wantErr:         false,
		},
		{
			name:            "major amount and decimal amount",
			formattedString: "1.23",
			currency:        currencies["USD"],
			want:            123,
			wantErr:         false,
		},
		{
			name:            "major amount with thousands and zero decimal amount",
			formattedString: "123,456.00",
			currency:        currencies["USD"],
			want:            12345600,
			wantErr:         false,
		},
		{
			name:            "major amount with thousands and decimal amount",
			formattedString: "123,456.78",
			currency:        currencies["USD"],
			want:            12345678,
			wantErr:         false,
		},
		{
			name:            "major amount with multiple thousands and zero decimal amount",
			formattedString: "123,456,789.00",
			currency:        currencies["USD"],
			want:            12345678900,
			wantErr:         false,
		},
		{
			name:            "major amount with no decimal amount",
			formattedString: "123",
			currency:        currencies["USD"],
			want:            12300,
			wantErr:         false,
		},
		{
			name:            "major amount with partial decimal amount",
			formattedString: "1.2",
			currency:        currencies["USD"],
			want:            120,
			wantErr:         false,
		},
		{
			name:            "major amount with three digit fraction",
			formattedString: "1.234",
			currency:        currencies["IQD"],
			want:            1234,
			wantErr:         false,
		},
		{
			name:            "major amount with zero digit fraction",
			formattedString: "123",
			currency:        currencies["GNF"],
			want:            123,
			wantErr:         false,
		},
		{
			name:            "comma decimal, period thousand",
			formattedString: "123.456,78",
			currency:        currencies["ANG"],
			want:            12345678,
			wantErr:         false,
		},
		{
			name:            "negative major amount",
			formattedString: "-123.45",
			currency:        currencies["USD"],
			want:            -12345,
			wantErr:         false,
		},
		{
			name:            "zero major amount with decimal amount",
			formattedString: "0.01",
			currency:        currencies["USD"],
			want:            1,
			wantErr:         false,
		},
		{
			name:            "amount with spaces",
			formattedString: "123,  456. 12",
			currency:        currencies["USD"],
			want:            12345612,
			wantErr:         false,
		},
		{
			name:            "invalid input",
			formattedString: "abf.def",
			currency:        currencies["USD"],
			want:            0,
			wantErr:         true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseFormattedString(tc.formattedString, tc.currency)
			if (err != nil) != tc.wantErr {
				t.Errorf("parseFormattedString() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("parseFormattedString() got = %v, want %v", got, tc.want)
			}
		})
	}
}
