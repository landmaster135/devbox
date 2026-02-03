package script

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func Tag(content string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, "<script>"); err != nil {
			return err
		}
		if _, err := io.WriteString(w, content); err != nil {
			return err
		}
		_, err := io.WriteString(w, "</script>")
		return err
	})
}
