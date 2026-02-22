package config

import "testing"

func TestNewConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		operation          string
		targetTitle        string
		targetURL          string
		srcMarkdownContent string
		srcMarkdownFile    string
		outFilePath        string
		topHeadingLevel    int
		help               bool
		expectError        bool
		errorMessage       string
	}{
		{
			name:               "正常系_content指定",
			operation:          OperationPatchMarkdown,
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			srcMarkdownContent: "## 見出し\n本文\n",
			outFilePath:        "./out.md",
			topHeadingLevel:    2,
		},
		{
			name:            "正常系_file指定",
			operation:       OperationPatchMarkdown,
			targetTitle:     "OpenAI Blog",
			targetURL:       "https://openai.com/blog",
			srcMarkdownFile: "./in.md",
			outFilePath:     "./out.md",
			topHeadingLevel: 2,
		},
		{
			name:               "異常系_contentとfileの同時指定",
			operation:          OperationPatchMarkdown,
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			srcMarkdownContent: "## 見出し\n本文\n",
			srcMarkdownFile:    "./in.md",
			outFilePath:        "./out.md",
			topHeadingLevel:    2,
			expectError:        true,
			errorMessage:       "--src-markdown-content と --src-markdown-file は同時に指定できません",
		},
		{
			name:               "異常系_未対応operation",
			operation:          "unsupported",
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			srcMarkdownContent: "## 見出し\n本文\n",
			outFilePath:        "./out.md",
			topHeadingLevel:    2,
			expectError:        true,
			errorMessage:       "--operation には patch-markdown を指定してください",
		},
		{
			name:            "異常系_入力ソース未指定",
			operation:       OperationPatchMarkdown,
			targetTitle:     "OpenAI Blog",
			targetURL:       "https://openai.com/blog",
			outFilePath:     "./out.md",
			topHeadingLevel: 2,
			expectError:     true,
			errorMessage:    "--src-markdown-content または --src-markdown-file のいずれかは必須です (--operation=patch-markdown)",
		},
		{
			name:               "異常系_topHeadingLevelが0以下",
			operation:          OperationPatchMarkdown,
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			srcMarkdownContent: "## 見出し\n本文\n",
			outFilePath:        "./out.md",
			topHeadingLevel:    0,
			expectError:        true,
			errorMessage:       "--top-heading-level は 1 以上で指定してください",
		},
		{
			name:               "正常系_help",
			operation:          "",
			targetTitle:        "",
			targetURL:          "",
			srcMarkdownContent: "",
			srcMarkdownFile:    "",
			outFilePath:        "",
			topHeadingLevel:    0,
			help:               true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewConfig(
				tt.operation,
				tt.targetTitle,
				tt.targetURL,
				tt.srcMarkdownContent,
				tt.srcMarkdownFile,
				tt.outFilePath,
				tt.topHeadingLevel,
				tt.help,
			)

			if tt.expectError {
				if err == nil {
					t.Fatal("エラーが期待されますがnilでした")
				}
				if err.Error() != tt.errorMessage {
					t.Fatalf("エラーメッセージが期待と異なります: got=%q want=%q", err.Error(), tt.errorMessage)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewConfigがエラーを返しました: %v", err)
			}
			if got == nil {
				t.Fatal("Configがnilです")
			}
		})
	}
}
