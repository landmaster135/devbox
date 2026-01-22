package fetchers

import (
	"errors"
	"strings"
	"testing"
)

func TestSanitizeHTMLBody_RemovesUnwantedContent(t *testing.T) {
	t.Parallel()

	input := `
<html>
  <body>
    <header>skip me</header>
    <main class="main" data-testid="main">
      <div class="content" style="color:red" data-testid="content">
        <span class="inner" data-allow-mismatch="true">Hello</span>
        <script>alert('x')</script>
        <!-- comment -->
        <img class="img" onerror="alert(1)" sizes="100vw" src="image.png" srcset="large">
        <picture>
          <source sizes="50vw" srcset="small.jpg 1x, medium.jpg 2x">
        </picture>
      </div>
      <button>click</button>
    </main>
    <footer>bye</footer>
  </body>
</html>`

	got, err := SanitizeHTMLBody(input, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, "Hello") {
		t.Fatalf("expected sanitized HTML to contain text content, got: %s", got)
	}
	for _, forbidden := range []string{"<script", "class=", "style=", "data-testid", "button", "<!--", "sizes=", "srcset="} {
		if strings.Contains(got, forbidden) {
			t.Fatalf("sanitized HTML still contains %q: %s", forbidden, got)
		}
	}
}

func TestSanitizeHTMLBody_RemovesNestedStructures(t *testing.T) {
	t.Parallel()

	input := `
<html>
  <body>
    <main>
      <header>nested header</header>
      <section>
        <div class="wrapper">
          <div></div>
          <div>content</div>
        </div>
      </section>
      <footer>nested footer</footer>
    </main>
  </body>
</html>`

	got, err := SanitizeHTMLBody(input, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, forbidden := range []string{"<header", "<footer", "<div></div>"} {
		if strings.Contains(got, forbidden) {
			t.Fatalf("sanitized HTML still contains %q: %s", forbidden, got)
		}
	}
	if !strings.Contains(got, "content") {
		t.Fatalf("expected remaining text to survive sanitization: %s", got)
	}
}

func TestSanitizeHTMLBody_FallbacksToArticle(t *testing.T) {
	t.Parallel()

	input := `
<html>
  <body>
    <article>short</article>
    <article>
      <p>this is the longest article content available for selection</p>
    </article>
  </body>
</html>`

	got, err := SanitizeHTMLBody(input, true)
	if err != nil {
		t.Fatalf("unexpected error when falling back to article: %v", err)
	}

	if !strings.Contains(got, "longest article content") {
		t.Fatalf("expected fallback article content in sanitized HTML, got: %s", got)
	}
	if strings.Contains(got, "short") {
		t.Fatalf("expected only the longest article to be kept, got: %s", got)
	}
}

func TestSanitizeHTMLBody_MainAndArticleMissing(t *testing.T) {
	t.Parallel()

	_, err := SanitizeHTMLBody("<html><body><div>No main</div></body></html>", true)
	if err == nil {
		t.Fatalf("expected error when no main or article elements are present")
	}

	var notFound *MainNotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected mainNotFoundError, got %T", err)
	}
}

func TestCollapseBlankLines(t *testing.T) {
	t.Parallel()

	src := "line1\n\n\nline2\n\nline3"
	got := collapseBlankLines(src)
	if strings.Contains(got, "\n\n\n") {
		t.Fatalf("expected blank lines to collapse, got: %q", got)
	}
	expected := "line1\n\nline2\n\nline3"
	if got != expected {
		t.Fatalf("unexpected collapsed output. want %q, got %q", expected, got)
	}
}

func TestSanitizeHTMLBody_EmptyInput(t *testing.T) {
	t.Parallel()

	if _, err := SanitizeHTMLBody("   ", true); err == nil {
		t.Fatalf("expected error for empty html")
	}
}

func TestCollapseBlankLines_NoChange(t *testing.T) {
	t.Parallel()

	src := "single line"
	if got := collapseBlankLines(src); got != src {
		t.Fatalf("expected no change for single line, got %q", got)
	}
}
