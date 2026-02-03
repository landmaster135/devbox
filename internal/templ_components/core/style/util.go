package style

import (
	"github.com/a-h/templ"
)

// CSSStyle describes a collection of CSS rules and optional media queries.
type CSSStyle struct {
	Rules      []CSSRule
	MediaRules []CSSMediaRule
}

// CSSRule pairs a selector with a set of sanitized declarations.
type CSSRule struct {
	Selector     string
	Declarations []templ.SafeCSS
}

// CSSMediaRule wraps a collection of CSS rules under an at-rule query (e.g. @media).
type CSSMediaRule struct {
	Query string
	Rules []CSSRule
}

// Rule builds a CSSRule with non-empty declarations.
func Rule(selector string, declarations ...templ.SafeCSS) CSSRule {
	return rule(selector, declarations...)
}

func Media(query string, rules ...CSSRule) CSSMediaRule {
	return media(query, rules...)
}

// Property sanitizes a CSS property/value pair.
func Property(name string, value string) templ.SafeCSS {
	return templ.SanitizeCSS(name, value)
}

func MustSafeBoxShadow(property string, offsets []string, colors []string) (templ.SafeCSS, error) {
	return safeBoxShadow(property, offsets, colors)
}

func MustSafeBorder(property, width, style string, colors []string) (templ.SafeCSS, error) {
	return safeBorder(property, width, style, colors)
}

func MustSafeColorProperty(property, color string) (templ.SafeCSS, error) {
	return safeColorProperty(property, color)
}

func MustSafeFontFamily(property string, families []string) (templ.SafeCSS, error) {
	return safeFontFamily(property, families)
}

func MustSafeLinearGradient(property, angle string, stops ...string) (templ.SafeCSS, error) {
	return safeLinearGradient(property, angle, stops...)
}
