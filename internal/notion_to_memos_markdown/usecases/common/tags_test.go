package common

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	domain "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/domain"
	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

func TestBuildFrontMatterPairs(t *testing.T) {
	t.Parallel()

	content := domain.Content{
		BoughtAt: " 2024-01-01T09:00:00+09:00 ",
		Score:    88,
		Price:    1200,
		ConID:    " CO0001 ",
		URL:      " https://example.com/item ",
	}

	got := BuildFrontMatterPairs(content)
	want := []string{
		"bought_at=2024-01-01T09:00:00+09:00",
		"score_of_100=88",
		"price_yen=1200",
		"con_id=CO0001",
		"url=https://example.com/item",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BuildFrontMatterPairs() = %#v, want %#v", got, want)
	}
}

func TestNormalizeFrontMatterURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "normal", input: "https://example.com", want: "https://example.com"},
		{name: "trim", input: "  https://example.com/path  ", want: "https://example.com/path"},
		{name: "empty", input: "", want: `""`},
		{name: "spaces", input: "   ", want: `""`},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NormalizeFrontMatterURL(tt.input)
			if got != tt.want {
				t.Fatalf("NormalizeFrontMatterURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBuildTagsForContent(t *testing.T) {
	t.Parallel()

	t.Run("normal with frequent tags and dedupe", func(t *testing.T) {
		t.Parallel()

		content := domain.Content{
			Category:     "software",
			OwningStatus: "already",
			Color:        " gray ",
			Tags: []domain.ContentTag{
				{PageTitle: " JavaScript "},
				{PageTitle: "#Go"},
				{PageTitle: "Custom"},
				{PageTitle: "javascript"},
				{PageTitle: " "},
			},
		}
		frequentTags := []string{
			"31-programming/language/javascript",
			"#31-programming/language/go",
			"custom",
		}

		got, err := BuildTagsForContent(content, frequentTags)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := []string{
			RequiredBackupTag,
			"0a-content/software",
			"01-p/own-status/2-already",
			"01-p/color/gray",
			"31-programming/language/javascript",
			"31-programming/language/go",
			"custom",
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("BuildTagsForContent() = %#v, want %#v", got, want)
		}
	})

	t.Run("color empty", func(t *testing.T) {
		t.Parallel()

		content := domain.Content{
			Category:     "software",
			OwningStatus: "already",
			Color:        " ",
		}
		got, err := BuildTagsForContent(content, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(got) != 3 {
			t.Fatalf("len(tags) = %d, want 3", len(got))
		}
		if strings.Join(got, ",") != strings.Join([]string{
			RequiredBackupTag,
			"0a-content/software",
			"01-p/own-status/2-already",
		}, ",") {
			t.Fatalf("unexpected tags: %#v", got)
		}
	})
}

func TestBuildTagsForContent_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content domain.Content
		errPart string
	}{
		{
			name: "unsupported category",
			content: domain.Content{
				Category:     "unknown",
				OwningStatus: "already",
			},
			errPart: "未対応のcategoryです",
		},
		{
			name: "unsupported owning status",
			content: domain.Content{
				Category:     "software",
				OwningStatus: "unknown",
			},
			errPart: "未対応のowning_statusです",
		},
		{
			name: "unsupported color",
			content: domain.Content{
				Category:     "software",
				OwningStatus: "already",
				Color:        "unknown",
			},
			errPart: "未対応のcolorです",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := BuildTagsForContent(tt.content, nil)
			if err == nil || !strings.Contains(err.Error(), tt.errPart) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestMapCategoryTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		want     string
		wantFail bool
	}{
		{name: "webclip", input: "web-clip", want: "0a-content/web-clip"},
		{name: "subscription month", input: "subscription month", want: "0a-content/subscription-month"},
		{name: "book", input: "book", want: "0a-content/book"},
		{name: "unsupported", input: "unknown", wantFail: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := MapCategoryTag(tt.input)
			if tt.wantFail {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMapOwningStatusTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		want     string
		wantFail bool
	}{
		{input: "yet", want: "01-p/own-status/1-yet"},
		{input: "already", want: "01-p/own-status/2-already"},
		{input: "gone", want: "01-p/own-status/3-gone"},
		{input: "unknown", wantFail: true},
	}

	for _, tt := range tests {
		got, err := MapOwningStatusTag(tt.input)
		if tt.wantFail {
			if err == nil {
				t.Fatalf("expected error for %q", tt.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tt.input, err)
		}
		if got != tt.want {
			t.Fatalf("got %q, want %q", got, tt.want)
		}
	}
}

func TestMapColorTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		want     string
		wantFail bool
	}{
		{input: "gray", want: "01-p/color/gray"},
		{input: "pink", want: "01-p/color/pink"},
		{input: "purple", want: "01-p/color/purple"},
		{input: "unknown", wantFail: true},
	}

	for _, tt := range tests {
		got, err := MapColorTag(tt.input)
		if tt.wantFail {
			if err == nil {
				t.Fatalf("expected error for %q", tt.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tt.input, err)
		}
		if got != tt.want {
			t.Fatalf("got %q, want %q", got, tt.want)
		}
	}
}

func TestNormalizeKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{input: "Web-Clip", want: "webclip"},
		{input: "subscription_year", want: "subscriptionyear"},
		{input: " subscription month ", want: "subscriptionmonth"},
	}
	for _, tt := range tests {
		got := NormalizeKey(tt.input)
		if got != tt.want {
			t.Fatalf("NormalizeKey(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestResolveFrequentTag(t *testing.T) {
	t.Parallel()

	frequentTags := []string{
		"chromeextension",
		"#31-programming/language/javascript",
		"31-programming/language/typescript",
	}

	tests := []struct {
		name string
		tag  string
		want string
	}{
		{name: "no slash exact", tag: "chromeextension", want: "chromeextension"},
		{name: "with slash suffix", tag: "#javascript", want: "31-programming/language/javascript"},
		{name: "trim and lowercase", tag: "  TypeScript  ", want: "31-programming/language/typescript"},
		{name: "fallback", tag: "go", want: "go"},
		{name: "blank", tag: " ", want: ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ResolveFrequentTag(tt.tag, frequentTags)
			if got != tt.want {
				t.Fatalf("ResolveFrequentTag(%q) = %q, want %q", tt.tag, got, tt.want)
			}
		})
	}
}

func TestParseConNumber(t *testing.T) {
	t.Parallel()

	number, err := ParseConNumber("CO000012345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if number != 12345 {
		t.Fatalf("number = %d, want 12345", number)
	}

	number, err = ParseConNumber("CO-12-item-34")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if number != 34 {
		t.Fatalf("number = %d, want 34", number)
	}

	if _, err := ParseConNumber("invalid"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestLoadFrequentTags(t *testing.T) {
	t.Parallel()

	t.Run("normal and dedupe", func(t *testing.T) {
		t.Parallel()

		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte("# tags\r\n\r\n## Frequent Tags\r\n#alpha #beta #alpha\r\n### ignored\r\n#gamma\r\n## Others\r\n#delta\r\n"), nil
		}

		got, err := LoadFrequentTags(mockRepo, "/tmp/tags.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := []string{"alpha", "beta", "gamma"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("LoadFrequentTags() = %#v, want %#v", got, want)
		}
	})

	t.Run("section not found", func(t *testing.T) {
		t.Parallel()

		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte("# tags\n## Other\n#alpha\n"), nil
		}

		_, err := LoadFrequentTags(mockRepo, "/tmp/tags.md")
		if err == nil || !strings.Contains(err.Error(), "## Frequent Tags") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("read file error", func(t *testing.T) {
		t.Parallel()

		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return nil, errors.New("read failed")
		}

		_, err := LoadFrequentTags(mockRepo, "/tmp/tags.md")
		if err == nil || !strings.Contains(err.Error(), "tags.md の読み込みに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestUniqueTags(t *testing.T) {
	t.Parallel()

	got := uniqueTags([]string{
		"#alpha",
		" beta ",
		"#alpha",
		" ",
		"",
		"#beta",
	})
	want := []string{"alpha", "beta"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("uniqueTags() = %#v, want %#v", got, want)
	}
}
