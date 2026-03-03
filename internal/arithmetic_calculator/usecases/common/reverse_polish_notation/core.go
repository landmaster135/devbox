package reversePolishNotation

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// Evaluate は中置記法の数式を逆ポーランド記法に変換して計算します。
func Evaluate(expression string) (float64, error) {
	if expression == "" {
		return 0, fmt.Errorf("無効な数式です: %s", expression)
	}

	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}

	rpnTokens, err := toRPN(tokens)
	if err != nil {
		return 0, err
	}

	return evaluateRPN(rpnTokens)
}

type tokenType int

const (
	tokenNone tokenType = iota
	tokenNumber
	tokenOperator
	tokenFunction
	tokenLeftParen
	tokenRightParen
)

var operatorPrecedence = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"^": 3,
}

type unaryFunc func(float64) (float64, error)

var supportedFunctions = map[string]unaryFunc{
	"sqrt": func(value float64) (float64, error) {
		if value < 0 {
			return 0, fmt.Errorf("負数の平方根は計算できません")
		}
		return math.Sqrt(value), nil
	},
	"sin": func(value float64) (float64, error) {
		return math.Sin(value), nil
	},
	"cos": func(value float64) (float64, error) {
		return math.Cos(value), nil
	},
	"tan": func(value float64) (float64, error) {
		return math.Tan(value), nil
	},
}

func tokenize(expression string) ([]string, error) {
	var tokens []string
	prev := tokenNone

	for i := 0; i < len(expression); {
		ch := rune(expression[i])

		if unicode.IsSpace(ch) {
			i++
			continue
		}

		if unicode.IsDigit(ch) || ch == '.' {
			number, next, err := consumeNumber(expression, i)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, number)
			prev = tokenNumber
			i = next
			continue
		}

		if isOperatorRune(ch) {
			isUnary := (ch == '+' || ch == '-') && shouldTreatAsUnary(prev)
			if isUnary {
				if i+1 >= len(expression) {
					return nil, fmt.Errorf("符号の後に数値が必要です")
				}

				nextChar := rune(expression[i+1])
				if nextChar == '(' || unicode.IsLetter(nextChar) {
					factor := "1"
					if ch == '-' {
						factor = "-1"
					}
					tokens = append(tokens, factor, "*")
					prev = tokenOperator
					i++
					continue
				}

				if !unicode.IsDigit(nextChar) && nextChar != '.' {
					return nil, fmt.Errorf("符号の後に数値が必要です")
				}

				number, next, err := consumeNumber(expression, i+1)
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, string(ch)+number)
				prev = tokenNumber
				i = next
				continue
			}

			tokens = append(tokens, string(ch))
			prev = tokenOperator
			i++
			continue
		}

		if ch == '(' {
			tokens = append(tokens, "(")
			prev = tokenLeftParen
			i++
			continue
		}

		if ch == ')' {
			tokens = append(tokens, ")")
			prev = tokenRightParen
			i++
			continue
		}

		if unicode.IsLetter(ch) {
			start := i
			for i < len(expression) && unicode.IsLetter(rune(expression[i])) {
				i++
			}
			literal := strings.ToLower(expression[start:i])
			if _, ok := supportedFunctions[literal]; !ok {
				return nil, fmt.Errorf("未対応の識別子です: %s", literal)
			}
			tokens = append(tokens, literal)
			prev = tokenFunction
			continue
		}

		return nil, fmt.Errorf("無効な文字が含まれています: %c", ch)
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("無効な数式です: %s", expression)
	}

	return tokens, nil
}

func consumeNumber(expression string, start int) (string, int, error) {
	var builder []rune
	dotCount := 0
	hasDigit := false
	i := start

	for i < len(expression) {
		ch := rune(expression[i])
		if unicode.IsDigit(ch) {
			hasDigit = true
			builder = append(builder, ch)
			i++
			continue
		}
		if ch == '.' {
			dotCount++
			if dotCount > 1 {
				return "", start, fmt.Errorf("小数点の数が多すぎます")
			}
			builder = append(builder, ch)
			i++
			continue
		}
		break
	}

	if len(builder) == 0 || !hasDigit {
		return "", start, fmt.Errorf("数値が必要です")
	}

	return string(builder), i, nil
}

func toRPN(tokens []string) ([]string, error) {
	var output []string
	var stack []string

	for _, token := range tokens {
		if isFunctionToken(token) {
			stack = append(stack, token)
			continue
		}

		if isOperatorToken(token) {
			for len(stack) > 0 && isOperatorToken(stack[len(stack)-1]) &&
				((isLeftAssociative(token) && operatorPrecedence[token] <= operatorPrecedence[stack[len(stack)-1]]) ||
					(!isLeftAssociative(token) && operatorPrecedence[token] < operatorPrecedence[stack[len(stack)-1]])) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
			continue
		}

		if token == "(" {
			stack = append(stack, token)
			continue
		}

		if token == ")" {
			foundLeftParen := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top == "(" {
					foundLeftParen = true
					break
				}
				output = append(output, top)
			}
			if !foundLeftParen {
				return nil, fmt.Errorf("括弧の対応が正しくありません")
			}
			if len(stack) > 0 && isFunctionToken(stack[len(stack)-1]) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			continue
		}

		output = append(output, token)
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top == "(" || top == ")" {
			return nil, fmt.Errorf("括弧の対応が正しくありません")
		}
		output = append(output, top)
	}

	return output, nil
}

func evaluateRPN(tokens []string) (float64, error) {
	stack := make([]float64, 0, len(tokens))

	for _, token := range tokens {
		if fn, ok := supportedFunctions[token]; ok {
			if len(stack) < 1 {
				return 0, fmt.Errorf("無効な数式です")
			}
			value := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result, err := fn(value)
			if err != nil {
				return 0, err
			}
			stack = append(stack, result)
			continue
		}

		if isOperatorToken(token) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("無効な数式です")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("ゼロ除算は許可されていません")
				}
				result = a / b
			case "^":
				result = math.Pow(a, b)
			}

			stack = append(stack, result)
			continue
		}

		value, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return 0, fmt.Errorf("数値を解析できません: %s", token)
		}
		stack = append(stack, value)
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("無効な数式です")
	}

	return stack[0], nil
}

func isOperatorRune(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' || r == '^'
}

func isOperatorToken(token string) bool {
	_, ok := operatorPrecedence[token]
	return ok
}

func isFunctionToken(token string) bool {
	_, ok := supportedFunctions[token]
	return ok
}

func isLeftAssociative(operator string) bool {
	return operator != "^"
}

func shouldTreatAsUnary(prev tokenType) bool {
	return prev == tokenNone || prev == tokenOperator || prev == tokenLeftParen
}
