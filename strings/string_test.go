package strings

import "testing"

func TestSnackToCamel(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should be ok", args{"dd_dd"}, "DdDd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnackToCamel(tt.args.str); got != tt.want {
				t.Errorf("SnackToCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}
