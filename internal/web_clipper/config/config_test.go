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
		srcDir             string
		slug               string
		start              int
		digits             int
		sortByTime         bool
		sortByName         bool
		jsonOutput         bool
		verbose            bool
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
			start:              defaultRenameAttachmentStart,
			digits:             defaultRenameAttachmentDigit,
		},
		{
			name:            "正常系_file指定",
			operation:       OperationPatchMarkdown,
			targetTitle:     "OpenAI Blog",
			targetURL:       "https://openai.com/blog",
			srcMarkdownFile: "./in.md",
			outFilePath:     "./out.md",
			topHeadingLevel: 2,
			start:           defaultRenameAttachmentStart,
			digits:          defaultRenameAttachmentDigit,
		},
		{
			name:       "正常系_renameAttachments_name指定",
			operation:  OperationRenameAttachments,
			srcDir:     "./attachments",
			slug:       "openai-blog",
			start:      1,
			digits:     3,
			sortByName: true,
		},
		{
			name:       "正常系_renameAttachments_time指定",
			operation:  OperationRenameAttachments,
			srcDir:     "./attachments",
			slug:       "openai-blog",
			start:      0,
			digits:     2,
			sortByTime: true,
			jsonOutput: true,
			verbose:    true,
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
			start:              defaultRenameAttachmentStart,
			digits:             defaultRenameAttachmentDigit,
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
			start:              defaultRenameAttachmentStart,
			digits:             defaultRenameAttachmentDigit,
			expectError:        true,
			errorMessage:       "--operation には patch-markdown, rename-attachments を指定してください",
		},
		{
			name:            "異常系_入力ソース未指定",
			operation:       OperationPatchMarkdown,
			targetTitle:     "OpenAI Blog",
			targetURL:       "https://openai.com/blog",
			outFilePath:     "./out.md",
			topHeadingLevel: 2,
			start:           defaultRenameAttachmentStart,
			digits:          defaultRenameAttachmentDigit,
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
			start:              defaultRenameAttachmentStart,
			digits:             defaultRenameAttachmentDigit,
			expectError:        true,
			errorMessage:       "--top-heading-level は 1 以上で指定してください",
		},
		{
			name:         "異常系_renameAttachments_slug未指定",
			operation:    OperationRenameAttachments,
			srcDir:       "./attachments",
			start:        1,
			digits:       3,
			sortByName:   true,
			expectError:  true,
			errorMessage: "--slug は必須です (--operation=rename-attachments)",
		},
		{
			name:         "異常系_renameAttachments_srcDir未指定",
			operation:    OperationRenameAttachments,
			slug:         "openai-blog",
			start:        1,
			digits:       3,
			sortByName:   true,
			expectError:  true,
			errorMessage: "--src-dir は必須です (--operation=rename-attachments)",
		},
		{
			name:         "異常系_renameAttachments_slugに英大文字を含む",
			operation:    OperationRenameAttachments,
			srcDir:       "./attachments",
			slug:         "OpenAI-blog",
			start:        1,
			digits:       3,
			sortByName:   true,
			expectError:  true,
			errorMessage: "--slug は英小文字、数字、半角ハイフンのみ使用できます",
		},
		{
			name:         "異常系_renameAttachments_slugにアンダースコアを含む",
			operation:    OperationRenameAttachments,
			srcDir:       "./attachments",
			slug:         "openai_blog",
			start:        1,
			digits:       3,
			sortByName:   true,
			expectError:  true,
			errorMessage: "--slug は英小文字、数字、半角ハイフンのみ使用できます",
		},
		{
			name:         "異常系_renameAttachments_start未指定",
			operation:    OperationRenameAttachments,
			srcDir:       "./attachments",
			slug:         "openai-blog",
			start:        defaultRenameAttachmentStart,
			digits:       3,
			sortByName:   true,
			expectError:  true,
			errorMessage: "--start は 0 以上で指定してください (--operation=rename-attachments)",
		},
		{
			name:         "異常系_renameAttachments_digits不正",
			operation:    OperationRenameAttachments,
			srcDir:       "./attachments",
			slug:         "openai-blog",
			start:        1,
			digits:       0,
			sortByName:   true,
			expectError:  true,
			errorMessage: "--digits は 1 以上で指定してください",
		},
		{
			name:         "異常系_renameAttachments_sort未指定",
			operation:    OperationRenameAttachments,
			srcDir:       "./attachments",
			slug:         "openai-blog",
			start:        1,
			digits:       3,
			expectError:  true,
			errorMessage: "-time または -name のいずれか一方を指定してください",
		},
		{
			name:         "異常系_renameAttachments_sort同時指定",
			operation:    OperationRenameAttachments,
			srcDir:       "./attachments",
			slug:         "openai-blog",
			start:        1,
			digits:       3,
			sortByTime:   true,
			sortByName:   true,
			expectError:  true,
			errorMessage: "-time または -name のいずれか一方を指定してください",
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
			start:              defaultRenameAttachmentStart,
			digits:             defaultRenameAttachmentDigit,
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
				tt.srcDir,
				tt.slug,
				tt.start,
				tt.digits,
				tt.sortByTime,
				tt.sortByName,
				tt.jsonOutput,
				tt.verbose,
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
