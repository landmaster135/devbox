package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_OutFlagRequired(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-x1", "10", "-y1", "20", "-x2", "300", "-y2", "400"}, &stdout, &stderr)

	if code != exitCodeError {
		t.Fatalf("run() exitCode = %d, want %d", code, exitCodeError)
	}

	stderrText := stderr.String()
	if !strings.Contains(stderrText, "エラー: -out は必須です。出力先ディレクトリを指定してください。") {
		t.Fatalf("stderr = %q, want required -out error", stderrText)
	}
	if !strings.Contains(stderrText, "Usage of image-trimmer:") {
		t.Fatalf("stderr = %q, want usage output", stderrText)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}

func TestRun_OutFlagEmptyRequired(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-out", "", "-x1", "10", "-y1", "20", "-x2", "300", "-y2", "400"}, &stdout, &stderr)

	if code != exitCodeError {
		t.Fatalf("run() exitCode = %d, want %d", code, exitCodeError)
	}

	stderrText := stderr.String()
	if !strings.Contains(stderrText, "エラー: -out は必須です。出力先ディレクトリを指定してください。") {
		t.Fatalf("stderr = %q, want required -out error", stderrText)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}

func TestRun_InvalidCoordinatesWithOutFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-out", "./trimmed_images", "-x1", "300", "-y1", "20", "-x2", "10", "-y2", "400"}, &stdout, &stderr)

	if code != exitCodeError {
		t.Fatalf("run() exitCode = %d, want %d", code, exitCodeError)
	}

	stderrText := stderr.String()
	if !strings.Contains(stderrText, "エラー: 無効なトリミング座標です。x2 > x1, y2 > y1 である必要があります。") {
		t.Fatalf("stderr = %q, want invalid coordinates error", stderrText)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}
