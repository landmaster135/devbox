package html

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"

	test "github.com/landmaster135/devbox/internal/templ_components/core/test"
)

func TestDocumentWithLangRendersDoctypeAndEscapedAttribute(t *testing.T) {
	t.Parallel()

	lang := `en" onclick="alert('x')`
	child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<main id="app"></main>`)
		return err
	})

	markup := renderDocument(t, lang, child)

	trimmed := strings.TrimLeft(markup, "\ufeff\r\n\t ")
	if !strings.HasPrefix(strings.ToLower(trimmed), "<!doctype html>") {
		t.Fatalf("document should start with doctype, got %s", markup)
	}

	if !strings.Contains(markup, `<main id="app"></main>`) {
		t.Fatalf("expected child markup to be rendered, got %s", markup)
	}

	if strings.Contains(markup, `lang="en" onclick="alert('x')"`) {
		t.Fatalf("lang attribute should have been escaped, got %s", markup)
	}

	htmlNode := test.ParseSingleElement(t, markup, false)
	if got := test.GetAttr(htmlNode, "lang"); got != lang {
		t.Fatalf("expected lang %q but got %q", lang, got)
	}
}

func TestDocumentWithoutLangOmitsAttribute(t *testing.T) {
	t.Parallel()

	child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<section class="content"></section>`)
		return err
	})

	markup := renderDocument(t, "", child)
	htmlNode := test.ParseSingleElement(t, markup, false)

	if test.HasAttr(htmlNode, "lang") {
		t.Fatalf("expected no lang attribute on html tag, got %s", markup)
	}
}

func TestDocumentIgnoresNilChildren(t *testing.T) {
	t.Parallel()

	realChild := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<p data-role="content">child</p>`)
		return err
	})

	markup := renderDocument(t, "en", nil, realChild, nil)

	if count := strings.Count(markup, `<p data-role="content">`); count != 1 {
		t.Fatalf("expected exactly one paragraph child, got %d in %s", count, markup)
	}

	if strings.Contains(markup, `<nil>`) {
		t.Fatalf("nil children should not render anything, got %s", markup)
	}
}

func renderDocument(t *testing.T, lang string, children ...templ.Component) string {
	t.Helper()

	component := Document(lang, children...)
	if component == nil {
		t.Fatalf("Document returned nil component")
	}

	var buf strings.Builder
	return test.RenderComponent(t, component, &buf)
}
