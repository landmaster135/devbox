package style

import (
	"context"
	"io"
	"strings"

	"github.com/a-h/templ"
)

func renderStyle(s CSSStyle) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "<style type=\"text/css\">"); err != nil {
			return err
		}
		if err := writeRules(w, s.Rules, ""); err != nil {
			return err
		}
		for _, media := range s.MediaRules {
			if len(media.Rules) == 0 || strings.TrimSpace(media.Query) == "" {
				continue
			}
			if _, err := io.WriteString(w, "\n@media "+media.Query+" {"); err != nil {
				return err
			}
			if err := writeRules(w, media.Rules, "  "); err != nil {
				return err
			}
			if _, err := io.WriteString(w, "\n}"); err != nil {
				return err
			}
		}
		_, err := io.WriteString(w, "\n</style>")
		return err
	})
}

func writeRules(w io.Writer, rules []CSSRule, indent string) error {
	for _, rule := range rules {
		selector := strings.TrimSpace(rule.Selector)
		if selector == "" || len(rule.Declarations) == 0 {
			continue
		}
		if _, err := io.WriteString(w, "\n"+indent+selector+" {"); err != nil {
			return err
		}
		for _, declaration := range rule.Declarations {
			value := strings.TrimSpace(string(declaration))
			if value == "" {
				continue
			}
			if _, err := io.WriteString(w, "\n"+indent+"  "+value); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, "\n"+indent+"}"); err != nil {
			return err
		}
	}
	return nil
}
