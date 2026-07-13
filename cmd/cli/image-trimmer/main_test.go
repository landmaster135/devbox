package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_OutputDirFlagRequired(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-x1", "10", "-y1", "20", "-x2", "300", "-y2", "400"}, &stdout, &stderr)

	if code != exitCodeError {
		t.Fatalf("run() exitCode = %d, want %d", code, exitCodeError)
	}

	stderrText := stderr.String()
	if !strings.Contains(stderrText, "エラー: -output-dir は必須です。出力先ディレクトリを指定してください。") {
		t.Fatalf("stderr = %q, want required -output-dir error", stderrText)
	}
	if !strings.Contains(stderrText, "Usage of image-trimmer:") {
		t.Fatalf("stderr = %q, want usage output", stderrText)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}

func TestRun_OutputDirFlagEmptyRequired(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-output-dir", "", "-x1", "10", "-y1", "20", "-x2", "300", "-y2", "400"}, &stdout, &stderr)

	if code != exitCodeError {
		t.Fatalf("run() exitCode = %d, want %d", code, exitCodeError)
	}

	stderrText := stderr.String()
	if !strings.Contains(stderrText, "エラー: -output-dir は必須です。出力先ディレクトリを指定してください。") {
		t.Fatalf("stderr = %q, want required -output-dir error", stderrText)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}

func TestRun_InvalidCoordinatesWithOutFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-output-dir", "./trimmed_images", "-x1", "300", "-y1", "20", "-x2", "10", "-y2", "400"}, &stdout, &stderr)

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

func TestRun_LegacyFlagsRejected(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedStderr string
	}{
		{
			name:           "OutFlagRejected",
			args:           []string{"-out", "./trimmed_images", "-x1", "10", "-y1", "20", "-x2", "300", "-y2", "400"},
			expectedStderr: "flag provided but not defined: -out",
		},
		{
			name:           "SrcFlagRejected",
			args:           []string{"-src", "./images", "-output-dir", "./trimmed_images", "-x1", "10", "-y1", "20", "-x2", "300", "-y2", "400"},
			expectedStderr: "flag provided but not defined: -src",
		},
		{
			name:           "ArcFlagRejected",
			args:           []string{"-archive-dir", "./archive", "-arc", "./old_archive", "-output-dir", "./trimmed_images", "-x1", "10", "-y1", "20", "-x2", "300", "-y2", "400"},
			expectedStderr: "flag provided but not defined: -arc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := run(tt.args, &stdout, &stderr)

			if code != exitCodeError {
				t.Fatalf("run() exitCode = %d, want %d", code, exitCodeError)
			}

			stderrText := stderr.String()
			if !strings.Contains(stderrText, tt.expectedStderr) {
				t.Fatalf("stderr = %q, want %q", stderrText, tt.expectedStderr)
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
		})
	}
}
