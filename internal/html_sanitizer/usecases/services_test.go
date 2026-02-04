package usecases

import (
	"errors"
	"io/fs"
	"strings"
	"testing"
)

func TestSanitizerService_SanitizeFile(t *testing.T) {
	t.Parallel()

	readCalled := false
	writeCalled := false
	sanitizeCalled := false

	reader := func(path string) ([]byte, error) {
		readCalled = true
		if path != "input.html" {
			t.Fatalf("unexpected input path: %s", path)
		}
		return []byte("<html><main>body</main></html>"), nil
	}

	sanitizer := func(html string, omit bool) (string, error) {
		sanitizeCalled = true
		if html != "<html><main>body</main></html>" {
			t.Fatalf("unexpected html: %s", html)
		}
		if !omit {
			t.Fatalf("expected omit flag to be true")
		}
		return "<body>sanitized</body>", nil
	}

	var writtenPath string
	var writtenData string
	writer := func(path string, data []byte, perm fs.FileMode) error {
		writeCalled = true
		writtenPath = path
		writtenData = string(data)
		if perm != 0o600 {
			t.Fatalf("expected perm 0600, got %o", perm)
		}
		return nil
	}

	svc := NewSanitizerService(
		WithFileReader(reader),
		WithFileWriter(writer),
		WithSanitizer(sanitizer),
		WithOutputPermission(0o600),
	)

	got, err := svc.SanitizeFile("input.html", "out.html", true)
	if err != nil {
		t.Fatalf("SanitizeFile returned error: %v", err)
	}

	if got != "<body>sanitized</body>" {
		t.Fatalf("unexpected sanitize result: %s", got)
	}
	if !readCalled || !writeCalled || !sanitizeCalled {
		t.Fatalf("expected all collaborators to be called")
	}
	if writtenPath != "out.html" {
		t.Fatalf("unexpected output path: %s", writtenPath)
	}
	if writtenData != "<body>sanitized</body>" {
		t.Fatalf("unexpected written data: %s", writtenData)
	}
}

func TestSanitizerService_SanitizeFile_Errors(t *testing.T) {
	t.Parallel()

	errRead := errors.New("read error")
	errSanitize := errors.New("sanitize error")
	errWrite := errors.New("write error")

	cases := []struct {
		name       string
		inputPath  string
		outputPath string
		reader     fileReadFunc
		sanitizer  sanitizeFunc
		writer     fileWriteFunc
		wantErr    string
	}{
		{
			name:       "missing input",
			inputPath:  "",
			outputPath: "out.html",
			wantErr:    "inputPath",
		},
		{
			name:       "missing output",
			inputPath:  "in.html",
			outputPath: "",
			wantErr:    "outputPath",
		},
		{
			name:       "read error",
			inputPath:  "in.html",
			outputPath: "out.html",
			reader: func(string) ([]byte, error) {
				return nil, errRead
			},
			wantErr: "読み込み",
		},
		{
			name:       "sanitize error",
			inputPath:  "in.html",
			outputPath: "out.html",
			reader: func(string) ([]byte, error) {
				return []byte("html"), nil
			},
			sanitizer: func(string, bool) (string, error) {
				return "", errSanitize
			},
			wantErr: "サニタイズ",
		},
		{
			name:       "write error",
			inputPath:  "in.html",
			outputPath: "out.html",
			reader: func(string) ([]byte, error) {
				return []byte("html"), nil
			},
			sanitizer: func(string, bool) (string, error) {
				return "clean", nil
			},
			writer: func(string, []byte, fs.FileMode) error {
				return errWrite
			},
			wantErr: "書き込み",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := []SanitizerServiceOption{}
			if tc.reader != nil {
				opts = append(opts, WithFileReader(tc.reader))
			}
			if tc.sanitizer != nil {
				opts = append(opts, WithSanitizer(tc.sanitizer))
			}
			if tc.writer != nil {
				opts = append(opts, WithFileWriter(tc.writer))
			}

			svc := NewSanitizerService(opts...)

			_, err := svc.SanitizeFile(tc.inputPath, tc.outputPath, true)
			if err == nil {
				t.Fatalf("expected error")
			}
			if tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
