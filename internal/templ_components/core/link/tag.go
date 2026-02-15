package link

import (
	"context"
	"html"
	"io"
	"sort"
	"strings"

	"github.com/a-h/templ"
)

// Attributes represents the common attributes supported by <link> tags.
type Attributes struct {
	Rel             string
	Href            string
	Type            string
	Sizes           string
	Media           string
	Title           string
	CrossOrigin     string
	ReferrerPolicy  string
	As              string
	Integrity       string
	ID              string
	ExtraAttributes map[string]string
}

// Tag renders a <link> element with the provided attributes.
func Tag(attrs Attributes) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		var builder strings.Builder
		builder.WriteString("<link")
		writeAttr(&builder, "rel", attrs.Rel)
		writeAttr(&builder, "href", attrs.Href)
		writeAttr(&builder, "type", attrs.Type)
		writeAttr(&builder, "sizes", attrs.Sizes)
		writeAttr(&builder, "media", attrs.Media)
		writeAttr(&builder, "title", attrs.Title)
		writeAttr(&builder, "crossorigin", attrs.CrossOrigin)
		writeAttr(&builder, "referrerpolicy", attrs.ReferrerPolicy)
		writeAttr(&builder, "as", attrs.As)
		writeAttr(&builder, "integrity", attrs.Integrity)
		writeAttr(&builder, "id", attrs.ID)

		if len(attrs.ExtraAttributes) > 0 {
			keys := make([]string, 0, len(attrs.ExtraAttributes))
			for key := range attrs.ExtraAttributes {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				writeAttr(&builder, key, attrs.ExtraAttributes[key])
			}
		}

		builder.WriteString(">")
		_, err := io.WriteString(w, builder.String())
		return err
	})
}

func writeAttr(builder *strings.Builder, name, value string) {
	if name == "" || value == "" {
		return
	}
	builder.WriteByte(' ')
	builder.WriteString(name)
	builder.WriteString("=\"")
	builder.WriteString(html.EscapeString(value))
	builder.WriteString("\"")
}
