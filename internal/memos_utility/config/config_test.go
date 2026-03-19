package config

import (
	"strings"
	"testing"
)

func TestConfig_ParseFlagsFromArgs_CreateWebClip_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=create-web-clip",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content-file=/tmp/web-summary-20240719-231059-palworld-steam.md",
		"-attachments=./a.png, ./b.txt",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}

	if cfg.Operation != OperationCreateWebClip {
		t.Fatalf("operation = %s, want %s", cfg.Operation, OperationCreateWebClip)
	}
	if cfg.TimeoutSeconds != defaultTimeoutSeconds {
		t.Fatalf("timeout = %d, want %d", cfg.TimeoutSeconds, defaultTimeoutSeconds)
	}
	if cfg.ContentFile != "/tmp/web-summary-20240719-231059-palworld-steam.md" {
		t.Fatalf("contentFile = %s, want expected path", cfg.ContentFile)
	}
	if cfg.Attachments != "./a.png, ./b.txt" {
		t.Fatalf("attachments = %s, want ./a.png, ./b.txt", cfg.Attachments)
	}
}

func TestConfig_ParseFlagsFromArgs_CreateMovieClip_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{
		"-operation=create-movie-clip",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content-file=/tmp/movie-summary-20260319-055716-trump-masako.md",
		"-timeout=45",
	})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}

	if cfg.Operation != OperationCreateMovieClip {
		t.Fatalf("operation = %s, want %s", cfg.Operation, OperationCreateMovieClip)
	}
	if cfg.TimeoutSeconds != 45 {
		t.Fatalf("timeout = %d, want 45", cfg.TimeoutSeconds)
	}
}

func TestConfig_ParseFlagsFromArgs_Help_Normal(t *testing.T) {
	cfg, err := ParseFlagsFromArgs([]string{"-help"})
	if err != nil {
		t.Fatalf("ParseFlagsFromArgs() error = %v", err)
	}
	if !cfg.Help {
		t.Fatal("help = false, want true")
	}
}

func TestConfig_ParseFlagsFromArgs_ValidationError_Error(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name: "OperationMissing",
			args: []string{
				"-base-url=https://memos.example.com",
				"-api-token=test-token",
				"-content-file=/tmp/web-summary-20240719-231059-palworld.md",
			},
			wantErr: "operation パラメータは必須です",
		},
		{
			name: "OperationUnsupported",
			args: []string{
				"-operation=create",
				"-base-url=https://memos.example.com",
				"-api-token=test-token",
				"-content-file=/tmp/web-summary-20240719-231059-palworld.md",
			},
			wantErr: "未対応の operation です",
		},
		{
			name: "BaseURLMissing",
			args: []string{
				"-operation=create-web-clip",
				"-api-token=test-token",
				"-content-file=/tmp/web-summary-20240719-231059-palworld.md",
			},
			wantErr: "base-url パラメータは必須です",
		},
		{
			name: "APITokenMissing",
			args: []string{
				"-operation=create-web-clip",
				"-base-url=https://memos.example.com",
				"-content-file=/tmp/web-summary-20240719-231059-palworld.md",
			},
			wantErr: "api-token パラメータは必須です",
		},
		{
			name: "ContentFileMissing",
			args: []string{
				"-operation=create-web-clip",
				"-base-url=https://memos.example.com",
				"-api-token=test-token",
			},
			wantErr: "content-file パラメータは必須です",
		},
		{
			name: "TimeoutInvalid",
			args: []string{
				"-operation=create-web-clip",
				"-base-url=https://memos.example.com",
				"-api-token=test-token",
				"-content-file=/tmp/web-summary-20240719-231059-palworld.md",
				"-timeout=0",
			},
			wantErr: "timeout パラメータは 1 以上",
		},
		{
			name: "WebClipPatternInvalid",
			args: []string{
				"-operation=create-web-clip",
				"-base-url=https://memos.example.com",
				"-api-token=test-token",
				"-content-file=/tmp/movie-summary-20260319-055716-trump.md",
			},
			wantErr: "web-summary-YYYYMMDD-hhmmss-<slug>.md",
		},
		{
			name: "MovieClipPatternInvalid",
			args: []string{
				"-operation=create-movie-clip",
				"-base-url=https://memos.example.com",
				"-api-token=test-token",
				"-content-file=/tmp/web-summary-20240719-231059-palworld.md",
			},
			wantErr: "movie-summary-YYYYMMDD-hhmmss-<slug>.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFlagsFromArgs(tt.args)
			if err == nil {
				t.Fatal("ParseFlagsFromArgs() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want %s", err, tt.wantErr)
			}
		})
	}
}
