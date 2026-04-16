package usecases

import (
	"errors"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/web_clipper/infrastructures/filesystem"
)

func TestPatchMarkdown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		targetTitle        string
		targetURL          string
		srcMarkdownContent string
		srcMarkdownFile    string
		outFilePath        string
		topHeadingLevel    int
		setupMock          func(*filesystem.MockRepository)
		expectError        bool
		errorMessage       string
		expectedResult     string
		expectedReadPath   string
		expectedReadCalls  int
		expectedWritePath  string
		expectedWriteData  string
		expectedWriteCalls int
	}{
		{
			name:               "正常系_最初の対象見出し直後へリンクを挿入",
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			srcMarkdownContent: "## 記事タイトル 要約\n\n### 見出し1\n本文\n\n## 後続セクション\n",
			outFilePath:        "/tmp/out.md",
			topHeadingLevel:    2,
			expectedResult:     "出力しました: /tmp/out.md",
			expectedReadCalls:  0,
			expectedWritePath:  "/tmp/out.md",
			expectedWriteData:  "## 記事タイトル 要約\n- [OpenAI Blog](https://openai.com/blog)\n\n### 見出し1\n本文\n\n## 後続セクション\n",
			expectedWriteCalls: 1,
		},
		{
			name:            "正常系_file入力を利用",
			targetTitle:     "OpenAI Blog",
			targetURL:       "https://openai.com/blog",
			srcMarkdownFile: "/tmp/in.md",
			outFilePath:     "/tmp/out.md",
			topHeadingLevel: 2,
			setupMock: func(mockRepo *filesystem.MockRepository) {
				mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
					return []byte("## 記事タイトル 要約\n\n### 見出し1\n本文\n"), nil
				}
			},
			expectedResult:     "出力しました: /tmp/out.md",
			expectedReadPath:   "/tmp/in.md",
			expectedReadCalls:  1,
			expectedWritePath:  "/tmp/out.md",
			expectedWriteData:  "## 記事タイトル 要約\n- [OpenAI Blog](https://openai.com/blog)\n\n### 見出し1\n本文\n",
			expectedWriteCalls: 1,
		},
		{
			name:               "異常系_contentとfileの同時指定",
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			srcMarkdownContent: "## 記事タイトル 要約\n\n### 見出し1\n本文\n",
			srcMarkdownFile:    "/tmp/in.md",
			outFilePath:        "/tmp/out.md",
			topHeadingLevel:    2,
			expectError:        true,
			errorMessage:       "--src-markdown-content と --src-markdown-file は同時に指定できません",
			expectedReadCalls:  0,
			expectedWriteCalls: 0,
		},
		{
			name:               "異常系_見出しレベル4以上を含む",
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			srcMarkdownContent: "## 記事タイトル 要約\n\n#### NG\n本文\n",
			outFilePath:        "/tmp/out.md",
			topHeadingLevel:    2,
			expectError:        true,
			errorMessage:       "見出しレベル4以上",
			expectedReadCalls:  0,
			expectedWriteCalls: 0,
		},
		{
			name:               "異常系_指定見出しが見つからない",
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			srcMarkdownContent: "### 見出し1\n本文\n",
			outFilePath:        "/tmp/out.md",
			topHeadingLevel:    2,
			expectError:        true,
			errorMessage:       "見出しレベル2 が見つかりませんでした",
			expectedReadCalls:  0,
			expectedWriteCalls: 0,
		},
		{
			name:               "異常系_入力ソース未指定",
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			outFilePath:        "/tmp/out.md",
			topHeadingLevel:    2,
			expectError:        true,
			errorMessage:       "--src-markdown-content または --src-markdown-file のいずれかは必須です",
			expectedReadCalls:  0,
			expectedWriteCalls: 0,
		},
		{
			name:               "異常系_outFilePathにカンマを含む",
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			srcMarkdownContent: "## 見出し\n本文\n",
			outFilePath:        "/tmp/out,invalid.md",
			topHeadingLevel:    2,
			expectError:        true,
			errorMessage:       "--out-file-path にカンマは使用できません",
			expectedReadCalls:  0,
			expectedWriteCalls: 0,
		},
		{
			name:               "異常系_topHeadingLevelが0以下",
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			srcMarkdownContent: "## 見出し\n本文\n",
			outFilePath:        "/tmp/out.md",
			topHeadingLevel:    0,
			expectError:        true,
			errorMessage:       "--top-heading-level は 1 以上で指定してください",
			expectedReadCalls:  0,
			expectedWriteCalls: 0,
		},
		{
			name:               "異常系_出力書き込みに失敗",
			targetTitle:        "OpenAI Blog",
			targetURL:          "https://openai.com/blog",
			srcMarkdownContent: "## 見出し\n本文\n",
			outFilePath:        "/tmp/out.md",
			topHeadingLevel:    2,
			setupMock: func(mockRepo *filesystem.MockRepository) {
				mockRepo.WriteFileFunc = func(path string, data []byte) error {
					return errors.New("disk full")
				}
			},
			expectError:        true,
			errorMessage:       "出力ファイルへの書き込みに失敗しました",
			expectedReadCalls:  0,
			expectedWriteCalls: 1,
		},
		{
			name:            "異常系_file読み込みに失敗",
			targetTitle:     "OpenAI Blog",
			targetURL:       "https://openai.com/blog",
			srcMarkdownFile: "/tmp/in.md",
			outFilePath:     "/tmp/out.md",
			topHeadingLevel: 2,
			setupMock: func(mockRepo *filesystem.MockRepository) {
				mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
					return nil, errors.New("permission denied")
				}
			},
			expectError:        true,
			errorMessage:       "入力ファイルの読み込みに失敗しました",
			expectedReadPath:   "/tmp/in.md",
			expectedReadCalls:  1,
			expectedWriteCalls: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &filesystem.MockRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			service := NewService(mockRepo)
			result, err := service.PatchMarkdown(
				tt.targetTitle,
				tt.targetURL,
				tt.srcMarkdownContent,
				tt.srcMarkdownFile,
				tt.outFilePath,
				tt.topHeadingLevel,
			)

			if tt.expectError {
				if err == nil {
					t.Fatal("エラーが期待されますがnilでした")
				}
				if !strings.Contains(err.Error(), tt.errorMessage) {
					t.Fatalf("エラーメッセージが期待と異なります: got=%q want contains %q", err.Error(), tt.errorMessage)
				}
				if len(mockRepo.ReadFileCalls) != tt.expectedReadCalls {
					t.Fatalf("ReadFile呼び出し回数が期待と異なります: got=%d want=%d", len(mockRepo.ReadFileCalls), tt.expectedReadCalls)
				}
				if tt.expectedReadPath != "" && len(mockRepo.ReadFileCalls) > 0 && mockRepo.ReadFileCalls[0] != tt.expectedReadPath {
					t.Fatalf("ReadFile pathが期待と異なります: got=%q want=%q", mockRepo.ReadFileCalls[0], tt.expectedReadPath)
				}
				if len(mockRepo.WriteFileCalls) != tt.expectedWriteCalls {
					t.Fatalf("WriteFile呼び出し回数が期待と異なります: got=%d want=%d", len(mockRepo.WriteFileCalls), tt.expectedWriteCalls)
				}
				return
			}

			if err != nil {
				t.Fatalf("PatchMarkdownがエラーを返しました: %v", err)
			}

			if result != tt.expectedResult {
				t.Fatalf("戻り値が期待と異なります: got=%q want=%q", result, tt.expectedResult)
			}

			if len(mockRepo.ReadFileCalls) != tt.expectedReadCalls {
				t.Fatalf("ReadFile呼び出し回数が期待と異なります: got=%d want=%d", len(mockRepo.ReadFileCalls), tt.expectedReadCalls)
			}
			if tt.expectedReadPath != "" {
				if len(mockRepo.ReadFileCalls) == 0 {
					t.Fatal("ReadFileが呼び出されていません")
				}
				if mockRepo.ReadFileCalls[0] != tt.expectedReadPath {
					t.Fatalf("ReadFile pathが期待と異なります: got=%q want=%q", mockRepo.ReadFileCalls[0], tt.expectedReadPath)
				}
			}

			if len(mockRepo.WriteFileCalls) != tt.expectedWriteCalls {
				t.Fatalf("WriteFile呼び出し回数が期待と異なります: got=%d want=%d", len(mockRepo.WriteFileCalls), tt.expectedWriteCalls)
			}
			if len(mockRepo.WriteFileCalls) == 0 {
				t.Fatal("WriteFileが呼び出されていません")
			}

			gotCall := mockRepo.WriteFileCalls[0]
			if gotCall.Path != tt.expectedWritePath {
				t.Fatalf("WriteFile pathが期待と異なります: got=%q want=%q", gotCall.Path, tt.expectedWritePath)
			}
			if string(gotCall.Data) != tt.expectedWriteData {
				t.Fatalf("WriteFile dataが期待と異なります:\n--- got ---\n%s\n--- want ---\n%s", string(gotCall.Data), tt.expectedWriteData)
			}
		})
	}
}
