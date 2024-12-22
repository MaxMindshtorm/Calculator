package calcul_test

import (
	"testing"

	"github.com/MaxMindshtorm/calculator/pkg/calcul"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		expectErr  bool
	}{
		{"3 + 5", 8, false},
		{"10 - 2 * 3", 4, false},
		{"(2 + 3) * 4", 20, false},
		{"10 / 2", 5, false},
		{"10 / 0", 0, true},
		{"3 +", 0, true},
		{"a + 5", 0, true},
	}

	for _, test := range tests {
		result, err := calcul.Calc(test.expression)
		if (err != nil) != test.expectErr {
			t.Errorf("Calc(%q): unexpected error status: got %v, want error: %v", test.expression, err != nil, test.expectErr)
		}
		if !test.expectErr && result != test.expected {
			t.Errorf("Calc(%q): got %v, want %v", test.expression, result, test.expected)
		}
	}
}
