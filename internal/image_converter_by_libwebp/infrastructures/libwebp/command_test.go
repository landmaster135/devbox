package libwebp

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestCommandConverter_BuildArgs_Normal(t *testing.T) {
	converter := NewCommandConverter()

	got := converter.buildArgs("input.jpg", "output.webp", 99, 4, false)
	want := []string{
		"-preset", "photo",
		"-metadata", "icc",
		"-sharp_yuv",
		"-progress",
		"-short",
		"-m", "4",
		"-q", "99",
		"input.jpg",
		"-o", "output.webp",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildArgs() = %#v, want %#v", got, want)
	}
}

func TestCommandConverter_BuildArgs_Lossless_Normal(t *testing.T) {
	converter := NewCommandConverter()

	got := converter.buildArgs("input.png", "output.webp", 99, 6, true)
	want := []string{
		"-preset", "photo",
		"-metadata", "icc",
		"-sharp_yuv",
		"-progress",
		"-short",
		"-lossless",
		"-m", "6",
		"-q", "99",
		"input.png",
		"-o", "output.webp",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildArgs() = %#v, want %#v", got, want)
	}
}

func TestCommandConverter_CheckAvailable_Normal(t *testing.T) {
	converter := NewCommandConverter()
	converter.lookPath = func(file string) (string, error) {
		if file != "cwebp" {
			t.Fatalf("file = %q, want cwebp", file)
		}
		return "/usr/bin/cwebp", nil
	}

	if err := converter.CheckAvailable(); err != nil {
		t.Fatalf("CheckAvailable() error = %v", err)
	}
}

func TestCommandConverter_CheckAvailable_Error(t *testing.T) {
	converter := NewCommandConverter()
	converter.lookPath = func(file string) (string, error) {
		return "", errors.New("not found")
	}

	err := converter.CheckAvailable()
	if err == nil {
		t.Fatalf("CheckAvailable() error = nil")
	}
	if err.Error() != "libwebp パッケージが見つかりません: cwebp をインストールしてください" {
		t.Fatalf("CheckAvailable() error = %q", err.Error())
	}
}

func TestCommandConverter_ConvertToWebP_Normal(t *testing.T) {
	converter := NewCommandConverter()
	var gotName string
	var gotArgs []string
	converter.runCommand = func(ctx context.Context, name string, args ...string) (string, error) {
		gotName = name
		gotArgs = append([]string{}, args...)
		return "ok", nil
	}

	err := converter.ConvertToWebP(context.Background(), "input.jpg", "output.webp", 99, 4, false)
	if err != nil {
		t.Fatalf("ConvertToWebP() error = %v", err)
	}
	if gotName != "cwebp" {
		t.Fatalf("command name = %q, want cwebp", gotName)
	}
	wantArgs := converter.buildArgs("input.jpg", "output.webp", 99, 4, false)
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %#v, want %#v", gotArgs, wantArgs)
	}
}

func TestCommandConverter_ConvertToWebP_Error(t *testing.T) {
	converter := NewCommandConverter()
	converter.runCommand = func(ctx context.Context, name string, args ...string) (string, error) {
		return "failed output", errors.New("exit status 1")
	}

	err := converter.ConvertToWebP(context.Background(), "input.jpg", "output.webp", 99, 4, false)
	if err == nil {
		t.Fatalf("ConvertToWebP() error = nil")
	}
	if !strings.Contains(err.Error(), "cwebp 実行に失敗しました") {
		t.Fatalf("ConvertToWebP() error = %q", err.Error())
	}
	if !strings.Contains(err.Error(), "failed output") {
		t.Fatalf("ConvertToWebP() error = %q", err.Error())
	}
}
