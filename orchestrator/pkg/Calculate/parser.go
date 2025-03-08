package Calculate

import (
	"fmt"
	"strings"
	"unicode"
)

func PrefixToPostfix(exp string) ([]string, error) {
	var output []string
	var oper []rune

	priority := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
		'(': 0,
	}

	exp = strings.ReplaceAll(exp, " ", "")

	for i := 0; i < len(exp); {
		char := rune(exp[i])

		if unicode.IsDigit(char) || char == '.' {
			j := i
			for j < len(exp) && (unicode.IsDigit(rune(exp[j])) || exp[j] == '.') {
				j++
			}
			output = append(output, exp[i:j])
			i = j
			continue
		}

		switch char {
		case '+', '-', '*', '/':
			if char == '-' && (i == 0 || rune(exp[i-1]) == '(') {
				char = '*'
				output = append(output, "-1")
			}
			for len(oper) > 0 && priority[oper[len(oper)-1]] >= priority[char] {
				output = append(output, string(oper[len(oper)-1]))
				oper = oper[:len(oper)-1]
			}
			oper = append(oper, char)
		case '(':
			oper = append(oper, char)
		case ')':
			for len(oper) > 0 && oper[len(oper)-1] != '(' {
				output = append(output, string(oper[len(oper)-1]))
				oper = oper[:len(oper)-1]
			}
			if len(oper) == 0 {
				return nil, fmt.Errorf("wrong amount of brackets")
			}
			oper = oper[:len(oper)-1]
		default:
			return nil, fmt.Errorf("unknown char %c", char)
		}
		i++
	}

	for len(oper) > 0 {
		if oper[len(oper)-1] == '(' {
			return nil, fmt.Errorf("wrong amount of brackets")
		}
		output = append(output, string(oper[len(oper)-1]))
		oper = oper[:len(oper)-1]
	}

	return output, nil
}
