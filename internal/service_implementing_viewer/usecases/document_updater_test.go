package usecases

import (
	"strings"
	"testing"
)

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

func containsLine(content, line string) bool {
	return strings.Contains(content, line)
}
