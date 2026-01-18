package usecases

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func normalizeMapInput(input string) (string, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return "", errors.New("値が入力されていません")
	}

	args := make([]string, 0, len(fields))
	for i := 0; i < len(fields); {
		flagToken := fields[i]
		if !strings.HasPrefix(flagToken, "-") {
			return "", fmt.Errorf("フラグ %q は '-' で始まっている必要があります", flagToken)
		}

		name, inlineValue, hasInline, err := splitFlagToken(flagToken)
		if err != nil {
			return "", err
		}

		if hasInline {
			args = append(args, formatKeyValue(name, inlineValue))
			i++
			continue
		}

		if i+1 >= len(fields) {
			return "", fmt.Errorf("フラグ %q の値が不足しています", name)
		}

		valueToken := fields[i+1]
		if looksLikeFlagToken(valueToken) {
			return "", fmt.Errorf("フラグ %q の値が不足しています", name)
		}

		args = append(args, formatKeyValue(name, valueToken))
		i += 2
	}

	return strings.Join(args, " "), nil
}

func splitFlagToken(token string) (name, value string, hasInline bool, err error) {
	trimmed := strings.TrimLeft(token, "-")
	if trimmed == "" {
		return "", "", false, fmt.Errorf("無効なフラグです: %s", token)
	}

	parts := strings.SplitN(trimmed, "=", 2)
	name = strings.TrimSpace(parts[0])
	if name == "" {
		return "", "", false, fmt.Errorf("無効なフラグです: %s", token)
	}

	if len(parts) == 2 {
		value = parts[1]
		if value == "" {
			return "", "", false, fmt.Errorf("フラグ %q の値が空です", name)
		}
		return name, value, true, nil
	}

	return name, "", false, nil
}

func looksLikeFlagToken(token string) bool {
	if len(token) < 2 || token[0] != '-' {
		return false
	}

	runes := []rune(token[1:])
	if len(runes) == 0 {
		return false
	}

	first := runes[0]
	return unicode.IsLetter(first) || first == '-'
}
