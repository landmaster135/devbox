package usecases

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

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

func TestServiceRenameAttachments_Normal(t *testing.T) {
	t.Parallel()

	fixedNow := time.Date(2026, 7, 5, 17, 12, 34, 0, time.UTC)
	olderTime := time.Date(2026, 7, 5, 9, 0, 0, 0, time.UTC)
	newerTime := time.Date(2026, 7, 5, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name                string
		opts                RenameAttachmentsOptions
		entries             []filesystem.FileInfo
		expectedResult      string
		expectedRenameCalls []filesystem.RenameCall
	}{
		{
			name: "名前昇順で採番_Normal",
			opts: RenameAttachmentsOptions{
				SrcDir:     "/tmp/attachments",
				Slug:       "openai-blog",
				Start:      1,
				Digits:     3,
				SortByName: true,
			},
			entries: []filesystem.FileInfo{
				{Name: "b.jpg", Path: "/tmp/attachments/b.jpg", ModTime: newerTime},
				{Name: "nested", Path: "/tmp/attachments/nested", ModTime: olderTime, IsDir: true},
				{Name: "a.png", Path: "/tmp/attachments/a.png", ModTime: olderTime},
			},
			expectedResult: "リネームしました: 2件",
			expectedRenameCalls: []filesystem.RenameCall{
				{OldPath: "/tmp/attachments/a.png", NewPath: "/tmp/attachments/web-summary-20260705-171234-openai-blog_001.png"},
				{OldPath: "/tmp/attachments/b.jpg", NewPath: "/tmp/attachments/web-summary-20260705-171234-openai-blog_002.jpg"},
			},
		},
		{
			name: "更新時刻昇順で採番_Normal",
			opts: RenameAttachmentsOptions{
				SrcDir:     "/tmp/attachments",
				Slug:       "openai-blog",
				Start:      5,
				Digits:     2,
				SortByTime: true,
			},
			entries: []filesystem.FileInfo{
				{Name: "newer", Path: "/tmp/attachments/newer", ModTime: newerTime},
				{Name: "same-b.gif", Path: "/tmp/attachments/same-b.gif", ModTime: olderTime},
				{Name: "same-a.gif", Path: "/tmp/attachments/same-a.gif", ModTime: olderTime},
			},
			expectedResult: "リネームしました: 3件",
			expectedRenameCalls: []filesystem.RenameCall{
				{OldPath: "/tmp/attachments/same-a.gif", NewPath: "/tmp/attachments/web-summary-20260705-171234-openai-blog_05.gif"},
				{OldPath: "/tmp/attachments/same-b.gif", NewPath: "/tmp/attachments/web-summary-20260705-171234-openai-blog_06.gif"},
				{OldPath: "/tmp/attachments/newer", NewPath: "/tmp/attachments/web-summary-20260705-171234-openai-blog_07"},
			},
		},
		{
			name: "詳細出力_Normal",
			opts: RenameAttachmentsOptions{
				SrcDir:     "/tmp/attachments",
				Slug:       "openai-blog",
				Start:      1,
				Digits:     1,
				SortByName: true,
				Verbose:    true,
			},
			entries: []filesystem.FileInfo{
				{Name: "a.png", Path: "/tmp/attachments/a.png", ModTime: olderTime},
			},
			expectedResult: "リネームしました: 1件\n/tmp/attachments/a.png -> /tmp/attachments/web-summary-20260705-171234-openai-blog_1.png",
			expectedRenameCalls: []filesystem.RenameCall{
				{OldPath: "/tmp/attachments/a.png", NewPath: "/tmp/attachments/web-summary-20260705-171234-openai-blog_1.png"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &filesystem.MockRepository{
				ListDirectoryFunc: func(path string) ([]filesystem.FileInfo, error) {
					return tt.entries, nil
				},
			}
			service := NewService(mockRepo)
			service.now = func() time.Time {
				return fixedNow
			}

			result, err := service.RenameAttachments(tt.opts)
			if err != nil {
				t.Fatalf("RenameAttachmentsがエラーを返しました: %v", err)
			}
			if result != tt.expectedResult {
				t.Fatalf("戻り値が期待と異なります: got=%q want=%q", result, tt.expectedResult)
			}
			if len(mockRepo.RenameCalls) != len(tt.expectedRenameCalls) {
				t.Fatalf("Rename呼び出し回数が期待と異なります: got=%d want=%d", len(mockRepo.RenameCalls), len(tt.expectedRenameCalls))
			}
			for i, expectedCall := range tt.expectedRenameCalls {
				if mockRepo.RenameCalls[i] != expectedCall {
					t.Fatalf("Rename呼び出しが期待と異なります: got=%+v want=%+v", mockRepo.RenameCalls[i], expectedCall)
				}
			}
		})
	}
}

func TestServiceRenameAttachments_JSONOutput_Normal(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{
		ListDirectoryFunc: func(path string) ([]filesystem.FileInfo, error) {
			return []filesystem.FileInfo{
				{Name: "a.png", Path: "/tmp/attachments/a.png", ModTime: time.Date(2026, 7, 5, 9, 0, 0, 0, time.UTC)},
			}, nil
		},
	}
	service := NewService(mockRepo)
	service.now = func() time.Time {
		return time.Date(2026, 7, 5, 17, 12, 34, 0, time.UTC)
	}

	result, err := service.RenameAttachments(RenameAttachmentsOptions{
		SrcDir:     "/tmp/attachments",
		Slug:       "openai-blog",
		Start:      1,
		Digits:     3,
		SortByName: true,
		JSON:       true,
		Verbose:    true,
	})
	if err != nil {
		t.Fatalf("RenameAttachmentsがエラーを返しました: %v", err)
	}

	var decoded renameAttachmentsOutput
	if err := json.Unmarshal([]byte(result), &decoded); err != nil {
		t.Fatalf("JSON出力をdecodeできません: %v", err)
	}
	if decoded.RenamedCount != 1 {
		t.Fatalf("renamed_countが期待と異なります: got=%d want=1", decoded.RenamedCount)
	}
	if len(decoded.Files) != 1 {
		t.Fatalf("files件数が期待と異なります: got=%d want=1", len(decoded.Files))
	}
	if decoded.Files[0].To != "/tmp/attachments/web-summary-20260705-171234-openai-blog_001.png" {
		t.Fatalf("リネーム先が期待と異なります: got=%q", decoded.Files[0].To)
	}
}

func TestServiceRenameAttachments_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		opts          RenameAttachmentsOptions
		setupMock     func(*filesystem.MockRepository)
		errorMessage  string
		expectedCalls int
	}{
		{
			name: "srcDir未指定",
			opts: RenameAttachmentsOptions{
				Slug:       "openai-blog",
				Start:      1,
				Digits:     3,
				SortByName: true,
			},
			errorMessage: "--src-dir は必須です",
		},
		{
			name: "slug不正",
			opts: RenameAttachmentsOptions{
				SrcDir:     "/tmp/attachments",
				Slug:       "OpenAI",
				Start:      1,
				Digits:     3,
				SortByName: true,
			},
			errorMessage: "--slug は英小文字、数字、半角ハイフンのみ使用できます",
		},
		{
			name: "対象ファイルなし",
			opts: RenameAttachmentsOptions{
				SrcDir:     "/tmp/attachments",
				Slug:       "openai-blog",
				Start:      1,
				Digits:     3,
				SortByName: true,
			},
			setupMock: func(mockRepo *filesystem.MockRepository) {
				mockRepo.ListDirectoryFunc = func(path string) ([]filesystem.FileInfo, error) {
					return []filesystem.FileInfo{
						{Name: "nested", Path: "/tmp/attachments/nested", IsDir: true},
					}, nil
				}
			},
			errorMessage: "リネーム対象ファイルが見つかりませんでした",
		},
		{
			name: "リネーム先が既に存在",
			opts: RenameAttachmentsOptions{
				SrcDir:     "/tmp/attachments",
				Slug:       "openai-blog",
				Start:      1,
				Digits:     3,
				SortByName: true,
			},
			setupMock: func(mockRepo *filesystem.MockRepository) {
				mockRepo.ListDirectoryFunc = func(path string) ([]filesystem.FileInfo, error) {
					return []filesystem.FileInfo{
						{Name: "a.png", Path: "/tmp/attachments/a.png"},
					}, nil
				}
				mockRepo.ExistsFunc = func(path string) (bool, error) {
					return true, nil
				}
			},
			errorMessage: "リネーム先ファイルが既に存在します",
		},
		{
			name: "リネーム失敗",
			opts: RenameAttachmentsOptions{
				SrcDir:     "/tmp/attachments",
				Slug:       "openai-blog",
				Start:      1,
				Digits:     3,
				SortByName: true,
			},
			setupMock: func(mockRepo *filesystem.MockRepository) {
				mockRepo.ListDirectoryFunc = func(path string) ([]filesystem.FileInfo, error) {
					return []filesystem.FileInfo{
						{Name: "a.png", Path: "/tmp/attachments/a.png"},
					}, nil
				}
				mockRepo.RenameFunc = func(oldPath, newPath string) error {
					return errors.New("permission denied")
				}
			},
			errorMessage:  "ファイルのリネームに失敗しました",
			expectedCalls: 1,
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
			service.now = func() time.Time {
				return time.Date(2026, 7, 5, 17, 12, 34, 0, time.UTC)
			}

			_, err := service.RenameAttachments(tt.opts)
			if err == nil {
				t.Fatal("エラーが期待されますがnilでした")
			}
			if !strings.Contains(err.Error(), tt.errorMessage) {
				t.Fatalf("エラーメッセージが期待と異なります: got=%q want contains %q", err.Error(), tt.errorMessage)
			}
			if tt.expectedCalls > 0 && len(mockRepo.RenameCalls) != tt.expectedCalls {
				t.Fatalf("Rename呼び出し回数が期待と異なります: got=%d want=%d", len(mockRepo.RenameCalls), tt.expectedCalls)
			}
		})
	}
}
