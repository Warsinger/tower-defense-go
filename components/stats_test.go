package components

import "testing"

func Test_makeDisplayName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"HighTowerLevel", args{"HighTowerLevel"}, "High Tower Level"},
		{"TowersBuilt", args{"TowersBuilt"}, "Towers Built"},
		{"Gold", args{"Gold"}, "Gold"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeDisplayName(tt.args.name); got != tt.want {
				t.Errorf("makeDisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}
