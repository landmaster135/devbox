package usecases

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type mockDocumentUpdaterRepository struct {
	readFileFunc        func(path string) ([]byte, error)
	writeFileFunc       func(path string, data []byte, perm os.FileMode) error
	lastWritePath       string
	lastWriteContent    []byte
	lastWritePermission os.FileMode
}

func (m *mockDocumentUpdaterRepository) ReadFile(path string) ([]byte, error) {
	if m.readFileFunc != nil {
		return m.readFileFunc(path)
	}
	return nil, nil
}

func (m *mockDocumentUpdaterRepository) WriteFile(path string, data []byte, perm os.FileMode) error {
	m.lastWritePath = path
	m.lastWriteContent = append([]byte(nil), data...)
	m.lastWritePermission = perm
	if m.writeFileFunc != nil {
		return m.writeFileFunc(path, data, perm)
	}
	return nil
}

func (m *mockDocumentUpdaterRepository) ListDirectories(path string) ([]string, error) {
	return []string{}, nil
}

func (m *mockDocumentUpdaterRepository) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func TestBuildUpdatedDocument(t *testing.T) {
	original := "# Title\n\n" + implementationHeading + "\n\n" + "| old | table |\n\n" + statisticsHeading + "\n\n- old stats\n\n## 次章\n内容\n"

	table := "| service | cli |\n| test | ✅ |"
	stats := &ServiceStatistics{TotalServices: 1, CLICount: 1}

	updated, err := buildUpdatedDocument(original, table, stats)
	if err != nil {
		t.Fatalf("予期せぬエラー: %v", err)
	}

	if !containsLine(updated, "| test | ✅ |") {
		t.Errorf("テーブルが更新されていません\n%s", updated)
	}

	if containsLine(updated, "| old | table |") {
		t.Errorf("旧テーブルが残っています")
	}

	if !containsLine(updated, "- **総サービス数**: 1") {
		t.Errorf("統計情報が更新されていません")
	}
	if containsLine(updated, "- old stats") {
		t.Errorf("旧統計情報が残っています")
	}
}

func TestBuildUpdatedDocumentMissingHeading(t *testing.T) {
	original := "# Title\n\n## No sections"
	table := "| service | cli |"
	stats := &ServiceStatistics{}

	if _, err := buildUpdatedDocument(original, table, stats); err == nil {
		t.Fatalf("エラーが発生しませんでした")
	}
}

func TestUpdateDocumentationFileWithRepository_Normal(t *testing.T) {
	original := "# Title\n\n" + implementationHeading + "\n\n" + "| old | table |\n\n" + statisticsHeading + "\n\n- old stats\n"
	stats := &ServiceStatistics{TotalServices: 2, CLICount: 1}
	table := "| service | cli |\n| alpha | ✅ |"
	mockRepo := &mockDocumentUpdaterRepository{
		readFileFunc: func(path string) ([]byte, error) {
			return []byte(original), nil
		},
	}

	err := updateDocumentationFileWithRepository(mockRepo, "docs/status.md", table, stats)
	if err != nil {
		t.Fatalf("予期せぬエラー: %v", err)
	}

	if mockRepo.lastWritePath != "docs/status.md" {
		t.Fatalf("書き込みパスが期待値と異なります。期待値: %s, 実際: %s", "docs/status.md", mockRepo.lastWritePath)
	}
	if mockRepo.lastWritePermission != 0o644 {
		t.Fatalf("書き込みパーミッションが期待値と異なります。期待値: %o, 実際: %o", 0o644, mockRepo.lastWritePermission)
	}

	updated := string(mockRepo.lastWriteContent)
	if !containsLine(updated, "| alpha | ✅ |") {
		t.Fatalf("テーブルが更新されていません: %s", updated)
	}
	if !containsLine(updated, "- **総サービス数**: 2") {
		t.Fatalf("統計情報が更新されていません: %s", updated)
	}
}

func TestUpdateDocumentationFileWithRepository_ReadError(t *testing.T) {
	mockRepo := &mockDocumentUpdaterRepository{
		readFileFunc: func(path string) ([]byte, error) {
			return nil, errors.New("read failed")
		},
	}

	err := updateDocumentationFileWithRepository(mockRepo, "docs/status.md", "", &ServiceStatistics{})
	if err == nil {
		t.Fatal("エラーが発生しませんでした")
	}
}

func TestUpdateDocumentationFileWithRepository_WriteError(t *testing.T) {
	original := "# Title\n\n" + implementationHeading + "\n\n" + statisticsHeading + "\n"
	mockRepo := &mockDocumentUpdaterRepository{
		readFileFunc: func(path string) ([]byte, error) {
			return []byte(original), nil
		},
		writeFileFunc: func(path string, data []byte, perm os.FileMode) error {
			return errors.New("write failed")
		},
	}

	err := updateDocumentationFileWithRepository(mockRepo, "docs/status.md", "", &ServiceStatistics{})
	if err == nil {
		t.Fatal("エラーが発生しませんでした")
	}
}

func TestUpdateDocumentationFile_Normal(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "status.md")
	original := "# Title\n\n" + implementationHeading + "\n\n" + statisticsHeading + "\n\n## Next\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatalf("テストファイル作成に失敗: %v", err)
	}

	err := UpdateDocumentationFile(path, "| service | cli |\n| app | ✅ |", &ServiceStatistics{TotalServices: 1})
	if err != nil {
		t.Fatalf("予期せぬエラー: %v", err)
	}

	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("更新後ファイルの読み込みに失敗: %v", err)
	}
	text := string(updated)
	if !containsLine(text, "| app | ✅ |") {
		t.Fatalf("テーブルが更新されていません: %s", text)
	}
	if !containsLine(text, "- **総サービス数**: 1") {
		t.Fatalf("統計が更新されていません: %s", text)
	}
}

func TestBuildStatisticsSection_NilStats(t *testing.T) {
	lines := buildStatisticsSection(nil)
	if len(lines) == 0 {
		t.Fatal("統計行が空です")
	}
	if !containsLine(strings.Join(lines, "\n"), "- **総サービス数**: 0") {
		t.Fatal("nil stats の既定値が反映されていません")
	}
}

func containsLine(content, line string) bool {
	return strings.Contains(content, line)
}
