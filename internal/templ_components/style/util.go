package style

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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
	return CSSRule{
		Selector:     strings.TrimSpace(selector),
		Declarations: filterDeclarations(declarations),
	}
}

// Media builds a CSSMediaRule while dropping empty child rules.
func Media(query string, rules ...CSSRule) CSSMediaRule {
	filtered := make([]CSSRule, 0, len(rules))
	for _, rule := range rules {
		if strings.TrimSpace(rule.Selector) == "" {
			continue
		}
		if len(rule.Declarations) == 0 {
			continue
		}
		filtered = append(filtered, rule)
	}
	return CSSMediaRule{Query: strings.TrimSpace(query), Rules: filtered}
}

// Property sanitizes a CSS property/value pair.
func Property(name string, value string) templ.SafeCSS {
	return templ.SanitizeCSS(name, value)
}

// MustSafeProperty allows callers to provide a value already marked safe.
func MustSafeProperty(name string, value templ.SafeCSSProperty) templ.SafeCSS {
	return templ.SanitizeCSS(name, value)
}

func filterDeclarations(values []templ.SafeCSS) []templ.SafeCSS {
	filtered := make([]templ.SafeCSS, 0, len(values))
	for _, decl := range values {
		if strings.TrimSpace(string(decl)) == "" {
			continue
		}
		filtered = append(filtered, decl)
	}
	return filtered
}

func safeLinearGradient(property, angle string, stops ...string) (templ.SafeCSS, error) {
	angle = strings.TrimSpace(angle)
	if !gradientAnglePattern.MatchString(angle) {
		return "", fmt.Errorf("invalid gradient angle %q", angle)
	}
	if len(stops) < 2 {
		return "", errors.New("linear-gradient requires at least two color stops")
	}
	sanitizedStops := make([]string, len(stops))
	for i, stop := range stops {
		sanitized, err := sanitizeColorStop(stop)
		if err != nil {
			return "", fmt.Errorf("invalid color stop %q: %w", stop, err)
		}
		sanitizedStops[i] = sanitized
	}
	value := fmt.Sprintf("linear-gradient(%s, %s)", angle, strings.Join(sanitizedStops, ", "))
	return MustSafeProperty(property, templ.SafeCSSProperty(value)), nil
}

func safeBoxShadow(property string, offsets []string, colors []string) (templ.SafeCSS, error) {
	if len(offsets) == 0 {
		return "", errors.New("box-shadow requires at least one offset value")
	}
	if len(colors) == 0 {
		return "", errors.New("box-shadow requires at least one color")
	}
	sanitizedOffsets := make([]string, len(offsets))
	for i, off := range offsets {
		length, err := sanitizeLength(off, true)
		if err != nil {
			return "", fmt.Errorf("invalid box-shadow offset %q: %w", off, err)
		}
		sanitizedOffsets[i] = length
	}
	sanitizedColors := make([]string, len(colors))
	for i, color := range colors {
		sanitized, err := sanitizeColor(color)
		if err != nil {
			return "", fmt.Errorf("invalid box-shadow color %q: %w", color, err)
		}
		sanitizedColors[i] = sanitized
	}
	value := strings.Join(sanitizedOffsets, " ") + " " + strings.Join(sanitizedColors, ", ")
	return MustSafeProperty(property, templ.SafeCSSProperty(value)), nil
}

func safeBorder(property, width, style string, colors []string) (templ.SafeCSS, error) {
	width = strings.TrimSpace(width)
	if width == "" {
		return "", errors.New("border width is required")
	}
	sanitizedWidth, err := sanitizeLength(width, false)
	if err != nil {
		return "", fmt.Errorf("invalid border width %q: %w", width, err)
	}
	styleLower := strings.ToLower(strings.TrimSpace(style))
	if _, ok := allowedBorderStyles[styleLower]; !ok {
		return "", fmt.Errorf("invalid border style %q", style)
	}
	if len(colors) == 0 {
		return "", errors.New("border requires at least one color")
	}
	sanitizedColors := make([]string, len(colors))
	for i, color := range colors {
		sanitized, err := sanitizeColor(color)
		if err != nil {
			return "", fmt.Errorf("invalid border color %q: %w", color, err)
		}
		sanitizedColors[i] = sanitized
	}
	value := fmt.Sprintf("%s %s %s", sanitizedWidth, styleLower, strings.Join(sanitizedColors, ", "))
	return MustSafeProperty(property, templ.SafeCSSProperty(value)), nil
}

func MustSafeLinearGradient(property, angle string, stops ...string) templ.SafeCSS {
	css, err := safeLinearGradient(property, angle, stops...)
	if err != nil {
		panic(err)
	}
	return css
}

func safeColorProperty(property, color string) (templ.SafeCSS, error) {
	sanitized, err := sanitizeColor(color)
	if err != nil {
		return "", err
	}
	return MustSafeProperty(property, templ.SafeCSSProperty(sanitized)), nil
}

func MustSafeColorProperty(property, color string) templ.SafeCSS {
	css, err := safeColorProperty(property, color)
	if err != nil {
		panic(err)
	}
	return css
}

func MustSafeBoxShadow(property string, offsets []string, colors []string) templ.SafeCSS {
	css, err := safeBoxShadow(property, offsets, colors)
	if err != nil {
		panic(err)
	}
	return css
}

func MustSafeBorder(property, width, style string, colors []string) templ.SafeCSS {
	css, err := safeBorder(property, width, style, colors)
	if err != nil {
		panic(err)
	}
	return css
}

var (
	gradientAnglePattern = regexp.MustCompile(`^(?:-?\d+(?:\.\d+)?)(?:deg|rad|turn)$`)
	lengthPattern        = regexp.MustCompile(`^-?\d+(?:\.\d+)?(?:px|rem|em|vw|vh|%)$`)
	unitlessZeroPattern  = regexp.MustCompile(`^-?0(?:\.0+)?$`)
	hexColorPattern      = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	rgbaColorPattern     = regexp.MustCompile(`^rgba?\(\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*\d{1,3}(?:\s*,\s*(?:0|1|0?\.\d+))?\s*\)$`)
	percentPattern       = regexp.MustCompile(`^(?:100|\d{1,2})(?:\.\d+)?%$`)
)

var allowedBorderStyles = map[string]struct{}{
	"solid":  {},
	"dashed": {},
	"dotted": {},
	"double": {},
	"groove": {},
	"ridge":  {},
	"inset":  {},
	"outset": {},
	"none":   {},
	"hidden": {},
}

var allowedColorKeywords = map[string]struct{}{
	"transparent":  {},
	"currentcolor": {},
}

func sanitizeLength(value string, allowUnitlessZero bool) (string, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return "", errors.New("length cannot be empty")
	}
	if allowUnitlessZero && unitlessZeroPattern.MatchString(v) {
		return v, nil
	}
	if lengthPattern.MatchString(v) {
		return v, nil
	}
	return "", errors.New("invalid length")
}

func sanitizeColor(value string) (string, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return "", errors.New("color cannot be empty")
	}
	if hexColorPattern.MatchString(v) || rgbaColorPattern.MatchString(v) {
		return v, nil
	}
	if _, ok := allowedColorKeywords[strings.ToLower(v)]; ok {
		return strings.ToLower(v), nil
	}
	return "", errors.New("invalid color")
}

func sanitizeColorStop(value string) (string, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return "", errors.New("color stop cannot be empty")
	}
	parts := strings.Fields(v)
	if len(parts) == 0 || len(parts) > 2 {
		return "", errors.New("invalid color stop format")
	}
	color, err := sanitizeColor(parts[0])
	if err != nil {
		return "", err
	}
	if len(parts) == 1 {
		return color, nil
	}
	position := parts[1]
	if !percentPattern.MatchString(position) {
		return "", errors.New("invalid gradient position")
	}
	return color + " " + position, nil
}
