package usecase_style

import (
	"context"
	"io"

	"github.com/a-h/templ"
	style "github.com/landmaster135/devbox/internal/templ_components/core/style"
)

func Tag() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		css, err := buildUsecaseStyle()
		if err != nil {
			return err
		}
		return style.Tag(css).Render(ctx, w)
	})
}
