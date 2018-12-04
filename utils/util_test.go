package utils

import "testing"

func TestMinInt(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"shold be 10", args{10, 12}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinInt(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MinInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
