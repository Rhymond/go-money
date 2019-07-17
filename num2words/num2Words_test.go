package num2words

import "testing"

func TestConvert(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{number: 0},
			want: "zero",
		},
		{
			name: "",
			args: args{number: 1},
			want: "one",
		},
		{
			name: "",
			args: args{number: 5},
			want: "five",
		},
		{
			name: "",
			args: args{number: 10},
			want: "ten",
		},
		{
			name: "",
			args: args{number: 11},
			want: "eleven",
		},
		{
			name: "",
			args: args{number: 12},
			want: "twelve",
		},
		{
			name: "",
			args: args{number: 17},
			want: "seventeen",
		},

		{
			name: "",
			args: args{number: 20},
			want: "twenty",
		},

		{
			name: "",
			args: args{number: 30},
			want: "thirty",
		},
		{
			name: "",
			args: args{number: 21},
			want: "twenty-one",
		},
		{
			name: "",
			args: args{number: 99},
			want: "ninety-nine",
		},
		{
			name: "",
			args: args{number: 100},
			want: "one hundred",
		},
		{
			name: "",
			args: args{number: 123},
			want: "one hundred twenty-three",
		},
		{
			name: "",
			args: args{number: 1024},
			want: "one thousand twenty-four",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Convert(tt.args.number); got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertAnd(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{number: 0},
			want: "zero",
		},
		{
			name: "",
			args: args{number: 1},
			want: "one",
		},
		{
			name: "",
			args: args{number: 5},
			want: "five",
		},
		{
			name: "",
			args: args{number: 10},
			want: "ten",
		},
		{
			name: "",
			args: args{number: 11},
			want: "eleven",
		},
		{
			name: "",
			args: args{number: 12},
			want: "twelve",
		},
		{
			name: "",
			args: args{number: 17},
			want: "seventeen",
		},

		{
			name: "",
			args: args{number: 20},
			want: "twenty",
		},

		{
			name: "",
			args: args{number: 30},
			want: "thirty",
		},
		{
			name: "",
			args: args{number: 21},
			want: "twenty-one",
		},
		{
			name: "",
			args: args{number: 99},
			want: "ninety-nine",
		},
		{
			name: "",
			args: args{number: 100},
			want: "one hundred",
		},
		{
			name: "",
			args: args{number: 123},
			want: "one hundred and twenty-three",
		},
		{
			name: "",
			args: args{number: 514},
			want: "five hundred and fourteen",
		},
		{
			name: "",
			args: args{number: 1111},
			want: "one thousand one hundred and eleven",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertAnd(tt.args.number); got != tt.want {
				t.Errorf("ConvertAnd() = %v, want %v", got, tt.want)
			}
		})
	}
}
