package models

import (
	"testing"
)

func TestFileContent_RemoveDuplicateLines(t *testing.T) {
	tests := []struct {
		name       string
		lines      []string
		startPos   int
		endPos     int
		wantLines  []string
		wantRemoved int
	}{
		{
			name: "通常のケース",
			lines: []string{
				"# Goマイクロサービス関連リンク集",
				"",
				"798. [Golang Logging: A Comprehensive Guide for Developers - Last9](https://last9.io/blog/golang-logging-guide-for-developers/)",
				"799. [End-to-End DevOps for a Golang Web App: Docker, EKS, AWS CI/CD | by Amit Maurya](https://towardsaws.com/end-to-end-devops-for-a-golang-web-app-docker-eks-aws-ci-cd-f7793e230ae9)",
				"800. [The 5 Best Logging Libraries for Golang - Highlight.io](https://www.highlight.io/blog/5-best-logging-libraries-for-go)",
				"801. [End-to-End DevOps for a Golang Web App: Docker, EKS, AWS CI/CD | by Amit Maurya](https://towardsaws.com/end-to-end-devops-for-a-golang-web-app-docker-eks-aws-ci-cd-f7793e230ae9)",
				"802. [The 5 Best Logging Library for C++ - Highlight.io](https://www.highlight.io/blog/5-best-logging-libraries-for-go)",
				"803. [The 5 Best Logging Libraries for Golang - Highlight.io](https://www.highlight.io/blog/5-best-logging-libraries-for-go)",
			},
			startPos: 5,
			endPos:   200,
			wantLines: []string{
				"# Goマイクロサービス関連リンク集",
				"",
				"798. [Golang Logging: A Comprehensive Guide for Developers - Last9](https://last9.io/blog/golang-logging-guide-for-developers/)",
				"799. [End-to-End DevOps for a Golang Web App: Docker, EKS, AWS CI/CD | by Amit Maurya](https://towardsaws.com/end-to-end-devops-for-a-golang-web-app-docker-eks-aws-ci-cd-f7793e230ae9)",
				"800. [The 5 Best Logging Libraries for Golang - Highlight.io](https://www.highlight.io/blog/5-best-logging-libraries-for-go)",
				"802. [The 5 Best Logging Library for C++ - Highlight.io](https://www.highlight.io/blog/5-best-logging-libraries-for-go)",
			},
			wantRemoved: 2,
		},
		{
			name: "空のファイル",
			lines: []string{},
			startPos: 5,
			endPos:   200,
			wantLines: []string{},
			wantRemoved: 0,
		},
		{
			name: "開始位置が行の長さを超える",
			lines: []string{
				"Line 1",
				"Line 2",
				"Line 3",
			},
			startPos: 10,
			endPos:   20,
			wantLines: []string{
				"Line 1",
				"Line 2",
				"Line 3",
			},
			wantRemoved: 0,
		},
		{
			name: "開始位置が負の場合",
			lines: []string{
				"Line 1",
				"Line 2",
				"Line 3",
			},
			startPos: -1,
			endPos:   5,
			wantLines: []string{
				"Line 1",
				"Line 2",
				"Line 3",
			},
			wantRemoved: 0,
		},
		{
			name: "終了位置が開始位置より小さい",
			lines: []string{
				"Line 1",
				"Line 2",
				"Line 3",
			},
			startPos: 5,
			endPos:   2,
			wantLines: []string{
				"Line 1",
				"Line 2",
				"Line 3",
			},
			wantRemoved: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := NewFileContent(tt.lines)
			gotRemoved, err := fc.RemoveDuplicateLines(tt.startPos, tt.endPos)

			// 特殊なケースではエラーが返されることを期待
			if tt.startPos < 0 || tt.endPos <= tt.startPos {
				if err == nil {
					t.Errorf("RemoveDuplicateLines() エラーが期待されましたが、nilが返されました")
				}
				// エラーが返された場合は、他のチェックをスキップ
				return
			} else if err != nil {
				t.Errorf("RemoveDuplicateLines() 予期しないエラー: %v", err)
				return
			}

			if gotRemoved != tt.wantRemoved {
				t.Errorf("RemoveDuplicateLines() = %v, want %v", gotRemoved, tt.wantRemoved)
			}

			if len(fc.Lines) != len(tt.wantLines) {
				t.Errorf("Lines length = %v, want %v", len(fc.Lines), len(tt.wantLines))
				return
			}

			for i := range fc.Lines {
				if fc.Lines[i] != tt.wantLines[i] {
					t.Errorf("Lines[%d] = %v, want %v", i, fc.Lines[i], tt.wantLines[i])
				}
			}
		})
	}
}
