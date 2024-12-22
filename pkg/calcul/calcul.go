package calcul

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Calc(exp string) (float64, error) {
	exp = strings.ReplaceAll(exp, " ", "")
	postfix, err := PrefixToPostfix(exp)
	if err != nil {
		return 0, err
	}

	return EvaluatePostfix(postfix)
}

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
				return nil, errors.New("неверное количество скобок")
			}
			oper = oper[:len(oper)-1]
		default:
			return nil, fmt.Errorf("неизвестный символ: %c", char)
		}
		i++
	}

	for len(oper) > 0 {
		if oper[len(oper)-1] == '(' {
			return nil, errors.New("неверное количество скобок")
		}
		output = append(output, string(oper[len(oper)-1]))
		oper = oper[:len(oper)-1]
	}

	return output, nil
}
func EvaluatePostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, elem := range postfix {
		if val, err := strconv.ParseFloat(elem, 64); err == nil {
			stack = append(stack, val)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("неверное выражение")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch elem {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, errors.New("деление на ноль")
				}
				result = a / b
			default:
				return 0, fmt.Errorf("неизвестный оператор: %s", elem)
			}

			stack = append(stack, result)
		}
	}
	if len(stack) != 1 {
		return 0, errors.New("неверное выражение")
	}

	return stack[0], nil
}
