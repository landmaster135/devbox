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

func TestSanitizeHTMLBody_RemovesFormNavAside(t *testing.T) {
	t.Parallel()

	input := `
<html>
  <body>
    <main>
      <nav>global navigation</nav>
      <article>
        <p>core content</p>
      </article>
      <form action="/submit">
        <input type="text" value="secret">
      </form>
      <aside>related links</aside>
    </main>
  </body>
</html>`

	got, err := SanitizeHTMLBody(input, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, "core content") {
		t.Fatalf("expected main content to remain, got: %s", got)
	}

	forbiddenFragments := []string{"<form", "<nav", "<aside", "global navigation", "secret", "related links"}
	for _, fragment := range forbiddenFragments {
		if strings.Contains(got, fragment) {
			t.Fatalf("expected sanitized HTML to remove %q, got: %s", fragment, got)
		}
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

func TestSanitizeHTMLBody(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       string
		omitFull    bool
		expectErr   string
		expectEqual bool
		contains    []string
		notContains []string
	}{
		{
			name:        "メイン要素あり・不要要素除去",
			input:       `<html><body><!--comment--><header>brand</header><div data-testid="wrapper"><main class="content" style="color:red"><script>console.log('xss')</script><p>text</p><!--secret--><button>click</button><div></div><div><div></div></div><span data-allow-mismatch="true" data-allow-missmatch="true" data-testid="main-span">keep</span><img src="/img.png" alt="logo" onerror="alert(1)" data-nuxt-img="true" sizes="100vw" srcset="/img.png 1x, /img@2x.png 2x"></main></div><footer>copyright</footer></body></html>`,
			contains:    []string{"<main", "<p>text</p>", "keep", "<img", "alt=\"logo\"", "src=\"/img.png\""},
			notContains: []string{"<script", "console.log", "class=", "style=", "<button", "<!--", "secret", "data-allow-missmatch", "data-allow-mismatch", "data-testid", "<header", "brand", "<footer", "copyright", "onerror=", "data-nuxt-img", "sizes=", "srcset=", "<div></div>", "<div><div></div></div>"},
		},
		{
			name:        "メイン要素なしは未加工で返す",
			input:       `<html><body><div>no-main</div></body></html>`,
			expectErr:   "main要素が見つかりません",
			expectEqual: true,
		},
		{
			name:        "空文字は未加工",
			input:       "",
			expectErr:   "HTMLが空です",
			expectEqual: true,
		},
		{
			name:        "プレーンテキストは未加工",
			input:       "plain text only",
			expectErr:   "HTMLタグが存在しません",
			expectEqual: true,
		},
		{
			name:        "余分な空行を圧縮",
			input:       "<html><body><main><p>line1</p>\n\n\n<p>line2</p></main></body></html>",
			notContains: []string{"\n\n\n"},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			got, err := SanitizeHTMLBody(testCase.input, testCase.omitFull)
			if testCase.expectErr != "" {
				if err == nil {
					t.Fatalf("expected error %q but got nil", testCase.expectErr)
				}
				if err.Error() != testCase.expectErr {
					t.Fatalf("unexpected error message: want %q but got %q", testCase.expectErr, err.Error())
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if testCase.expectEqual && got != testCase.input {
				t.Fatalf("expected sanitized result to equal input, diff: %q", got)
			}
			for _, fragment := range testCase.contains {
				if !strings.Contains(got, fragment) {
					t.Fatalf("expected %q to contain %q", got, fragment)
				}
			}
			for _, fragment := range testCase.notContains {
				if strings.Contains(got, fragment) {
					t.Fatalf("expected result to remove %q but was %q", fragment, got)
				}
			}
		})
	}
}
