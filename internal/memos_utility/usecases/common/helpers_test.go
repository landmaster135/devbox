package common

import (
	"strings"
	"testing"

	memos "github.com/landmaster135/devbox/internal/memos/usecases"
)

func TestNormalizeOperation_Normal(t *testing.T) {
	got := NormalizeOperation("  CREATE-WEB-CLIP ")
	if got != OperationCreateWebClip {
		t.Fatalf("NormalizeOperation() = %s, want %s", got, OperationCreateWebClip)
	}
}

func TestNormalizeAttachments_Normal(t *testing.T) {
	got := NormalizeAttachments([]string{" ./a.png ", "", " ./b.txt "})
	if len(got) != 2 || got[0] != "./a.png" || got[1] != "./b.txt" {
		t.Fatalf("NormalizeAttachments() = %#v, want [./a.png ./b.txt]", got)
	}
}

func TestNormalizeAttachments_Empty_Normal(t *testing.T) {
	got := NormalizeAttachments([]string{"", "   "})
	if got != nil {
		t.Fatalf("NormalizeAttachments() = %#v, want nil", got)
	}
}

func TestBuildDisplayTime_Normal(t *testing.T) {
	got, err := BuildDisplayTime(OperationCreateMovieClip, "/tmp/movie-summary-20260319-055716-sample.md")
	if err != nil {
		t.Fatalf("BuildDisplayTime() error = %v", err)
	}
	if got != "2026-03-19T05:57:16+09:00" {
		t.Fatalf("BuildDisplayTime() = %s, want 2026-03-19T05:57:16+09:00", got)
	}
}

func TestBuildDisplayTime_Error(t *testing.T) {
	tests := []struct {
		name       string
		operation  string
		content    string
		wantErrSub string
	}{
		{
			name:       "UnsupportedOperation",
			operation:  "unknown",
			content:    "/tmp/web-summary-20240719-231059-palworld.md",
			wantErrSub: "未対応の operation",
		},
		{
			name:       "InvalidFormat",
			operation:  OperationCreateWebClip,
			content:    "/tmp/invalid.md",
			wantErrSub: "形式が不正",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildDisplayTime(tt.operation, tt.content)
			if err == nil {
				t.Fatal("BuildDisplayTime() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErrSub) {
				t.Fatalf("error = %v, want %s", err, tt.wantErrSub)
			}
		})
	}
}

func TestResolveMemoIdentifier_Normal(t *testing.T) {
	tests := []struct {
		name string
		memo *memos.Memo
		want string
	}{
		{name: "Name", memo: &memos.Memo{Name: "memos/1"}, want: "memos/1"},
		{name: "UID", memo: &memos.Memo{UID: "uid-1"}, want: "uid-1"},
		{name: "ID", memo: &memos.Memo{ID: 42}, want: "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveMemoIdentifier(tt.memo)
			if err != nil {
				t.Fatalf("ResolveMemoIdentifier() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("ResolveMemoIdentifier() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestResolveMemoIdentifier_Error(t *testing.T) {
	tests := []struct {
		name       string
		memo       *memos.Memo
		wantErrSub string
	}{
		{name: "NilMemo", memo: nil, wantErrSub: "メモ情報が空です"},
		{name: "NoIdentifier", memo: &memos.Memo{}, wantErrSub: "いずれも取得できません"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ResolveMemoIdentifier(tt.memo)
			if err == nil {
				t.Fatal("ResolveMemoIdentifier() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErrSub) {
				t.Fatalf("error = %v, want %s", err, tt.wantErrSub)
			}
		})
	}
}

func TestResolveContentBaseNameFromAttachment_Normal(t *testing.T) {
	got, err := ResolveContentBaseNameFromAttachment("web-summary-20241225-233435-daikokuyu-event-info_11.webp")
	if err != nil {
		t.Fatalf("ResolveContentBaseNameFromAttachment() error = %v", err)
	}
	want := "web-summary-20241225-233435-daikokuyu-event-info"
	if got != want {
		t.Fatalf("ResolveContentBaseNameFromAttachment() = %s, want %s", got, want)
	}
}

func TestResolveContentBaseNameFromAttachment_Error(t *testing.T) {
	_, err := ResolveContentBaseNameFromAttachment("invalid.webp")
	if err == nil {
		t.Fatal("ResolveContentBaseNameFromAttachment() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "attachment-dir 内のファイル名が不正です") {
		t.Fatalf("error = %v, want attachment-dir validation message", err)
	}
}

func TestPatterns_Normal(t *testing.T) {
	if !MatchWebClipFile("web-summary-20241225-233435-daikokuyu-event-info.md") {
		t.Fatal("MatchWebClipFile() = false, want true")
	}
	if !MatchMovieClipFile("movie-summary-20260319-055716-trump-masako-diplomacy.md") {
		t.Fatal("MatchMovieClipFile() = false, want true")
	}
	if ts, ok := ParseWebClipDisplayTime("web-summary-20241225-233435-daikokuyu-event-info.md"); !ok || ts != "20241225233435" {
		t.Fatalf("ParseWebClipDisplayTime() = (%s, %v), want (20241225233435, true)", ts, ok)
	}
	if ts, ok := ParseMovieClipDisplayTime("movie-summary-20260319-055716-trump-masako-diplomacy.md"); !ok || ts != "20260319055716" {
		t.Fatalf("ParseMovieClipDisplayTime() = (%s, %v), want (20260319055716, true)", ts, ok)
	}
	if base, ok := ParseWebAttachmentContentBaseName("web-summary-20241225-233435-daikokuyu-event-info_02.webp"); !ok || base != "web-summary-20241225-233435-daikokuyu-event-info" {
		t.Fatalf("ParseWebAttachmentContentBaseName() = (%s, %v), want expected", base, ok)
	}
	if base, ok := ParseMovieAttachmentContentBaseName("movie-summary-20260319-055716-trump-masako-diplomacy_01.webp"); !ok || base != "movie-summary-20260319-055716-trump-masako-diplomacy" {
		t.Fatalf("ParseMovieAttachmentContentBaseName() = (%s, %v), want expected", base, ok)
	}
}
