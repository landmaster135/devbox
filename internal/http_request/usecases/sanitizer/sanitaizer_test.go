package services

import (
	"strings"
	"testing"
)

func TestSanitizeHTMLBody(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       string
		expectFound bool
		expectEqual bool
		contains    []string
		notContains []string
	}{
		{
			name:        "メイン要素あり・不要要素除去",
			input:       `<html><body><!--comment--><header>brand</header><div data-testid="wrapper"><main class="content" style="color:red"><script>console.log('xss')</script><p>text</p><!--secret--><button>click</button><div></div><div><div></div></div><span data-allow-mismatch="true" data-allow-missmatch="true" data-testid="main-span">keep</span><img src="/img.png" alt="logo" onerror="alert(1)" data-nuxt-img="true" sizes="100vw" srcset="/img.png 1x, /img@2x.png 2x"></main></div><footer>copyright</footer></body></html>`,
			expectFound: true,
			notContains: []string{"<script", "console.log", "class=", "style=", "<button", "<!--", "secret", "data-allow-missmatch", "data-allow-mismatch", "data-testid", "<header", "brand", "<footer", "copyright", "onerror=", "data-nuxt-img", "sizes=", "srcset=", "<div></div>", "<div><div></div></div>"},
		},
		{
			name:        "メイン要素なしは未加工で返す",
			input:       `<html><body><div>no-main</div></body></html>`,
			expectFound: false,
			expectEqual: true,
		},
		{
			name:        "空文字は未加工",
			input:       "",
			expectFound: false,
			expectEqual: true,
		},
		{
			name:        "プレーンテキストは未加工",
			input:       "plain text only",
			expectFound: false,
			expectEqual: true,
		},
		{
			name:        "余分な空行を圧縮",
			input:       "<html><body><main><p>line1</p>\n\n\n<p>line2</p></main></body></html>",
			expectFound: true,
			notContains: []string{"\n\n\n"},
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			got, found := SanitizeHTMLBody(testCase.input, false)
			if found != testCase.expectFound {
				t.Fatalf("expectFound=%v but got %v", testCase.expectFound, found)
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
