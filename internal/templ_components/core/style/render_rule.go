package style

import (
	"strings"

	"github.com/a-h/templ"
)

func rule(selector string, declarations ...templ.SafeCSS) CSSRule {
	return CSSRule{
		Selector:     strings.TrimSpace(selector),
		Declarations: filterDeclarations(declarations),
	}
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
