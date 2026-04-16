package config

import (
	"errors"
	"flag"
	"strings"
	"testing"

	flagParser "github.com/landmaster135/devbox/internal/memos_utility/infrastructures/flag_parser"
)

func TestConfig_ParseFlagsWithParser_Normal(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetString("operation", " CREATE-WEB-CLIP ")
	parser.SetString("base-url", " https://memos.example.com ")
	parser.SetString("api-token", " test-token ")
	parser.SetString("content-file", " /tmp/web-summary-20240719-231059-palworld.md ")
	parser.SetString("attachments", " ./a.png,./b.txt ")
	parser.SetInt("timeout", 45)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}

	if cfg.Operation != OperationCreateWebClip {
		t.Fatalf("operation = %q, want %q", cfg.Operation, OperationCreateWebClip)
	}
	if cfg.BaseURL != "https://memos.example.com" {
		t.Fatalf("baseURL = %q, want trimmed value", cfg.BaseURL)
	}
	if cfg.APIToken != "test-token" {
		t.Fatalf("apiToken = %q, want trimmed value", cfg.APIToken)
	}
	if cfg.ContentFile != "/tmp/web-summary-20240719-231059-palworld.md" {
		t.Fatalf("contentFile = %q, want trimmed value", cfg.ContentFile)
	}
	if cfg.Attachments != "./a.png,./b.txt" {
		t.Fatalf("attachments = %q, want trimmed value", cfg.Attachments)
	}
	if cfg.TimeoutSeconds != 45 {
		t.Fatalf("timeout = %d, want 45", cfg.TimeoutSeconds)
	}
}

func TestConfig_ParseFlagsWithParser_HelpByParseError_Normal(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetParseError(flag.ErrHelp)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if !cfg.Help {
		t.Fatal("help = false, want true")
	}
}

func TestConfig_ParseFlagsWithParser_ParseError_Error(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetParseError(errors.New("parse failed"))

	_, err := ParseFlagsWithParser(parser)
	if err == nil {
		t.Fatal("ParseFlagsWithParser() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "フラグの解析に失敗しました") {
		t.Fatalf("error = %v, want parse error prefix", err)
	}
}

func TestConfig_ParseFlagsWithParser_PositionalArgs_Error(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetArgs([]string{"extra"})

	_, err := ParseFlagsWithParser(parser)
	if err == nil {
		t.Fatal("ParseFlagsWithParser() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "未処理の位置引数があります") {
		t.Fatalf("error = %v, want positional args error", err)
	}
}
