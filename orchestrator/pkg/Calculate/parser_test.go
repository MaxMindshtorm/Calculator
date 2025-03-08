package Calculate

import (
	"reflect"
	"testing"
)

func TestPrefixToPostfix(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:    "Basic expression",
			input:   "2 + 2",
			want:    []string{"2", "2", "+"},
			wantErr: false,
		},
		{
			name:    "Expression with multiplication",
			input:   "2 + 2 * 2",
			want:    []string{"2", "2", "2", "*", "+"},
			wantErr: false,
		},
		{
			name:    "Expression with parentheses",
			input:   "(2 + 2) * 2",
			want:    []string{"2", "2", "+", "2", "*"},
			wantErr: false,
		},
		{
			name:    "Complex expression with all operations",
			input:   "3 + 4 * 2 / ( 1 - 5 )",
			want:    []string{"3", "4", "2", "*", "1", "5", "-", "/", "+"},
			wantErr: false,
		},
		{
			name:    "Expression with decimal numbers",
			input:   "3.14 + 2.71",
			want:    []string{"3.14", "2.71", "+"},
			wantErr: false,
		},
		{
			name:    "Expression with negative number at start",
			input:   "-5 + 3",
			want:    []string{"-1", "5", "*", "3", "+"},
			wantErr: false,
		},
		{
			name:    "Expression with negative number in parentheses",
			input:   "2 * (-5 + 3)",
			want:    []string{"2", "-1", "5", "*", "3", "+", "*"},
			wantErr: false,
		},
		{
			name:    "Expression with unbalanced opening bracket",
			input:   "2 + (3 * 4",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Expression with unbalanced closing bracket",
			input:   "2 + 3) * 4",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Expression with unknown character",
			input:   "2 + 3 $ 4",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Expression with spaces",
			input:   "  2  +  2  ",
			want:    []string{"2", "2", "+"},
			wantErr: false,
		},
		{
			name:    "Complex expression with negative numbers",
			input:   "10 / (-2) + (-3) * 4",
			want:    []string{"10", "-1", "2", "*", "/", "-1", "3", "*", "4", "*", "+"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrefixToPostfix(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrefixToPostfix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrefixToPostfix() = %v, want %v", got, tt.want)
			}
		})
	}
}
