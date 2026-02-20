package renamebodiesbycategoryid

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

func TestService_Execute_Normal(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{}
	mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
		return []byte(`[
			{"con_id":"CO000000150","category_id":"SF0001"},
			{"con_id":"CO000000151","category_id":"SF0002"},
			{"con_id":"CO000000251","category_id":"SF9999"}
		]`), nil
	}
	mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
		return []string{
			filepath.Join(dirPath, "SF0001_image_01.webp"),
			filepath.Join(dirPath, "SF0002-body.md"),
			filepath.Join(dirPath, "SF9999_note.md"),
			filepath.Join(dirPath, "_tmp.txt"),
		}, nil
	}
	mockRepo.FileExistsFunc = func(path string) (bool, error) {
		return false, nil
	}

	service := NewService(mockRepo)
	result, err := service.Execute("content", 150, 200, "/tmp/contents.json", "/tmp/resource")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "対象Content件数=2") ||
		!strings.Contains(result, "対象ファイル総数=4") ||
		!strings.Contains(result, "リネーム成功=2") ||
		!strings.Contains(result, "スキップ(プレフィックスなし)=1") ||
		!strings.Contains(result, "スキップ(マップ未対応)=1") {
		t.Fatalf("unexpected result: %s", result)
	}

	if len(mockRepo.RenameFileCalls) != 2 {
		t.Fatalf("rename count = %d, want 2", len(mockRepo.RenameFileCalls))
	}

	if mockRepo.RenameFileCalls[0].SrcPath != "/tmp/resource/SF0001_image_01.webp" ||
		mockRepo.RenameFileCalls[0].DstPath != "/tmp/resource/CO000000150_image_01.webp" {
		t.Fatalf("unexpected first rename call: %+v", mockRepo.RenameFileCalls[0])
	}
	if mockRepo.RenameFileCalls[1].SrcPath != "/tmp/resource/SF0002-body.md" ||
		mockRepo.RenameFileCalls[1].DstPath != "/tmp/resource/CO000000151-body.md" {
		t.Fatalf("unexpected second rename call: %+v", mockRepo.RenameFileCalls[1])
	}
}

func TestService_Execute_Error(t *testing.T) {
	t.Parallel()

	t.Run("unsupported page type", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{})
		_, err := service.Execute("artifact", 1, 10, "/tmp/contents.json", "/tmp/resource")
		if err == nil || err.Error() != "未対応のpage-typeです: artifact" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("invalid con number range", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{})
		_, err := service.Execute("content", 10, 1, "/tmp/contents.json", "/tmp/resource")
		if err == nil || err.Error() != "con_number_start は con_number_end 以下である必要があります" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("missing src json file", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{})
		_, err := service.Execute("content", 1, 10, " ", "/tmp/resource")
		if err == nil || err.Error() != "src-json-file パラメータは必須です" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("missing src resource dir", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{})
		_, err := service.Execute("content", 1, 10, "/tmp/contents.json", " ")
		if err == nil || err.Error() != "src-resource-dir パラメータは必須です" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("read json error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return nil, errors.New("read failed")
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", 1, 10, "/tmp/contents.json", "/tmp/resource")
		if err == nil || !strings.Contains(err.Error(), "src-json-file の読み込みに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte("{invalid"), nil
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", 1, 10, "/tmp/contents.json", "/tmp/resource")
		if err == nil || !strings.Contains(err.Error(), "Content JSONの解析に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("invalid con_id", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"invalid","category_id":"SF0001"}]`), nil
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", 1, 10, "/tmp/contents.json", "/tmp/resource")
		if err == nil || !strings.Contains(err.Error(), "con_id の解析に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("duplicate category id", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[
				{"con_id":"CO000000101","category_id":"SF0001"},
				{"con_id":"CO000000102","category_id":"SF0001"}
			]`), nil
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", 100, 200, "/tmp/contents.json", "/tmp/resource")
		if err == nil || !strings.Contains(err.Error(), "category_id が重複しています") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("list files recursive error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"CO000000101","category_id":"SF0001"}]`), nil
		}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			return nil, errors.New("read dir failed")
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", 100, 200, "/tmp/contents.json", "/tmp/resource")
		if err == nil || !strings.Contains(err.Error(), "src-resource-dir の読み取りに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("destination exists", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"CO000000101","category_id":"SF0001"}]`), nil
		}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			return []string{filepath.Join(dirPath, "SF0001_body.md")}, nil
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return true, nil
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", 100, 200, "/tmp/contents.json", "/tmp/resource")
		if err == nil || !strings.Contains(err.Error(), "リネーム先ファイルが既に存在します") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("destination exists check error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"CO000000101","category_id":"SF0001"}]`), nil
		}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			return []string{filepath.Join(dirPath, "SF0001_body.md")}, nil
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return false, errors.New("stat failed")
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", 100, 200, "/tmp/contents.json", "/tmp/resource")
		if err == nil || !strings.Contains(err.Error(), "リネーム先ファイルの確認に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("rename error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"CO000000101","category_id":"SF0001"}]`), nil
		}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			return []string{filepath.Join(dirPath, "SF0001_body.md")}, nil
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return false, nil
		}
		mockRepo.RenameFileFunc = func(srcPath, dstPath string) error {
			return errors.New("rename failed")
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", 100, 200, "/tmp/contents.json", "/tmp/resource")
		if err == nil || !strings.Contains(err.Error(), "ファイル名変更に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestExtractPrefixToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fileName string
		want     string
		wantOK   bool
	}{
		{name: "alpha numeric", fileName: "SF0001_image.png", want: "SF0001", wantOK: true},
		{name: "starts with symbol", fileName: "_tmp.png", want: "", wantOK: false},
		{name: "empty", fileName: " ", want: "", wantOK: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := extractPrefixToken(tt.fileName)
			if got != tt.want || ok != tt.wantOK {
				t.Fatalf("extractPrefixToken(%q) = (%q, %v), want (%q, %v)", tt.fileName, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}
