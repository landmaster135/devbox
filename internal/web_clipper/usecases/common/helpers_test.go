package common

import "testing"

func TestNormalizeNewlines_Normal(t *testing.T) {
	t.Parallel()

	got := NormalizeNewlines("a\r\nb\nc")
	want := "a\nb\nc"
	if got != want {
		t.Fatalf("NormalizeNewlines() = %q, want %q", got, want)
	}
}

func TestHeadingLevel_Normal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		line string
		want int
	}{
		{name: "見出し1", line: "# title", want: 1},
		{name: "タブ区切り", line: "###\ttitle", want: 3},
		{name: "前方空白", line: "  ## title", want: 2},
		{name: "空白なしは見出しではない", line: "##title", want: 0},
		{name: "本文", line: "text", want: 0},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := HeadingLevel(tt.line)
			if got != tt.want {
				t.Fatalf("HeadingLevel() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestContainsHeadingLevel4OrMore_Normal(t *testing.T) {
	t.Parallel()

	if !ContainsHeadingLevel4OrMore("## ok\n\n#### ng\n") {
		t.Fatal("ContainsHeadingLevel4OrMore() = false, want true")
	}
	if ContainsHeadingLevel4OrMore("## ok\n\n### ok\n") {
		t.Fatal("ContainsHeadingLevel4OrMore() = true, want false")
	}
}
