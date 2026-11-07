package domain

import "testing"

func TestValidateProfileText(t *testing.T) {
	tests := []struct {
		name  string
		value string
		max   int
		want  bool
	}{
		{"valid", " Mona ", 6, true},
		{"blank", " \t", 10, false},
		{"too long", "Mona", 3, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ValidateProfileText(test.value, test.max); got != test.want {
				t.Fatalf("ValidateProfileText() = %v, want %v", got, test.want)
			}
		})
	}
}
