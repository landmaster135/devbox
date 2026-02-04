package style

import (
	"slices"
	"strings"
)

// Media builds a CSSMediaRule while dropping empty child rules.
func media(query string, rules ...CSSRule) CSSMediaRule {
	filtered := make([]CSSRule, 0, len(rules))
	for _, rule := range rules {
		if !shouldRenderRule(rule) {
			continue
		}
		filtered = append(filtered, rule)
	}
	return CSSMediaRule{Query: strings.TrimSpace(query), Rules: filtered}
}

func trimmedSelector(rule CSSRule) string {
	return strings.TrimSpace(rule.Selector)
}

func renderableDeclarations(rule CSSRule) []string {
	filtered := make([]string, 0, len(rule.Declarations))
	for _, decl := range rule.Declarations {
		trimmed := strings.TrimSpace(string(decl))
		if trimmed == "" {
			continue
		}
		filtered = append(filtered, trimmed)
	}
	return filtered
}

func mediaQuery(media CSSMediaRule) string {
	return strings.TrimSpace(media.Query)
}

func shouldRenderRule(rule CSSRule) bool {
	if trimmedSelector(rule) == "" {
		return false
	}
	return len(renderableDeclarations(rule)) > 0
}

func mediaHasRenderableRules(media CSSMediaRule) bool {
	if mediaQuery(media) == "" {
		return false
	}
	return slices.ContainsFunc(media.Rules, shouldRenderRule)
}
