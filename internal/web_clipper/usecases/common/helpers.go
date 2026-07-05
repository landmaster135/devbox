package common

import "strings"

func NormalizeNewlines(content string) string {
	return strings.ReplaceAll(content, "\r\n", "\n")
}

func ContainsHeadingLevel4OrMore(markdownContent string) bool {
	lines := strings.Split(markdownContent, "\n")
	for _, line := range lines {
		if HeadingLevel(line) >= 4 {
			return true
		}
	}
	return false
}

func HeadingLevel(line string) int {
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" || trimmed[0] != '#' {
		return 0
	}

	level := 0
	for level < len(trimmed) && trimmed[level] == '#' {
		level++
	}
	if level == 0 {
		return 0
	}
	if level < len(trimmed) && trimmed[level] != ' ' && trimmed[level] != '\t' {
		return 0
	}

	return level
}
