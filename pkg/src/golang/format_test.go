package golang

import "testing"

func TestMakeIdentifier(t *testing.T) {
	tests := []struct {
		v    string
		want string
	}{
		{"", "_"},
		{" ", "_"},
		{"123", "_"},
		{" hello world", "HelloWorld"},
		{"-1hello2world", "Hello2world"},
		{"öüä%&§@'\"$!=?<+", "_"},
	}
	for _, tt := range tests {
		t.Run(tt.v, func(t *testing.T) {
			if got := MakeIdentifier(tt.v); got != tt.want {
				t.Errorf("MakeIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
