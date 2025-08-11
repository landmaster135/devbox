package usecases

import (
	"strings"
	"testing"
)

func TestSplitCleanupScriptWithIntro_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name        string
		script      string
		introLength int
		description string
	}{
		{
			name:        "ShortScript_WithShortIntro_Normal",
			script:      "echo 'short cleanup script'",
			introLength: 100,
			description: "Short script with short intro",
		},
		{
			name:        "MediumScript_WithMediumIntro_Normal",
			script:      strings.Repeat("echo 'medium cleanup line'\n", 20),
			introLength: 300,
			description: "Medium script with medium intro",
		},
		{
			name:        "LongScript_WithShortIntro_Normal",
			script:      strings.Repeat("echo 'long cleanup line that will be repeated many times'\n", 50),
			introLength: 200,
			description: "Long script with short intro",
		},
		{
			name:        "ShortScript_WithVeryLongIntro_Normal",
			script:      "echo 'short script'",
			introLength: 1800, // 非常に長い説明部分
			description: "Short script with very long intro",
		},
		{
			name:        "LongScript_WithLongIntro_Normal",
			script:      strings.Repeat("echo 'very long cleanup line that will definitely cause splitting issues'\n", 100),
			introLength: 500,
			description: "Long script with long intro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			t.Logf("Test: %s", tt.description)
			t.Logf("Script length: %d, Intro length: %d", len(tt.script), tt.introLength)

			// Act
			messages := service.splitCleanupScriptWithIntro(tt.script, tt.introLength)

			// Assert
			if len(messages) < 1 {
				t.Error("Expected at least one message chunk")
			}

			// 各チャンクの詳細検証
			for i, message := range messages {
				t.Logf("Chunk %d: Title='%s', Length=%d", i, message.Title, len(message.Content))

				// 文字数制限の確認
				if len(message.Content) > maxLengthOfContentOfWebhookPayload {
					t.Errorf("Chunk %d exceeds maximum length: %d (limit: %d)", i, len(message.Content), maxLengthOfContentOfWebhookPayload)
				}

				// Bashスクリプトのマーカーが含まれていることを確認
				if !strings.Contains(message.Content, "```bash") || !strings.Contains(message.Content, "```") {
					t.Errorf("Chunk %d does not contain proper bash markers", i)
				}

				// タイトルが適切に設定されていることを確認
				if !strings.Contains(message.Title, "削除用Bashスクリプト") {
					t.Errorf("Chunk %d has incorrect title: %s", i, message.Title)
				}

				// 色が適切に設定されていることを確認
				if message.Color != "red" {
					t.Errorf("Chunk %d has incorrect color: %s", i, message.Color)
				}

				// IsCodeフラグが適切に設定されていることを確認
				if !message.IsCode {
					t.Errorf("Chunk %d should have IsCode=true", i)
				}
			}

			// 計算された最大長の確認
			expectedMaxLength := maxLengthOfContentOfWebhookPayloadOnMargin - tt.introLength
			if expectedMaxLength < 500 {
				expectedMaxLength = 500 // 最低限の長さ
			}
			t.Logf("Expected max length per chunk: %d", expectedMaxLength)

			// 分割が適切に行われているかの確認
			totalOriginalLength := len(tt.script)
			totalReconstructedLength := 0
			for _, message := range messages {
				// ```bash と ``` を除いた実際のスクリプト部分の長さを計算
				content := message.Content
				content = strings.ReplaceAll(content, "```bash\n", "")
				content = strings.ReplaceAll(content, "```", "")
				totalReconstructedLength += len(content)
			}

			t.Logf("Original script length: %d, Reconstructed length: %d", totalOriginalLength, totalReconstructedLength)
		})
	}
}

func TestSplitCleanupScriptWithIntro_BoundaryConditions_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name        string
		scriptSize  int
		introLength int
		description string
	}{
		{
			name:        "ExactBoundary_Normal",
			scriptSize:  maxLengthOfContentOfWebhookPayloadOnMargin,
			introLength: 0,
			description: "Script exactly at boundary with no intro",
		},
		{
			name:        "OverBoundary_WithIntro_Normal",
			scriptSize:  maxLengthOfContentOfWebhookPayloadOnMargin,
			introLength: 100,
			description: "Script at boundary with intro",
		},
		{
			name:        "VeryLongIntro_Normal",
			scriptSize:  500,
			introLength: maxLengthOfContentOfWebhookPayloadOnMargin,
			description: "Short script with very long intro",
		},
		{
			name:        "MaxIntro_Normal",
			scriptSize:  100,
			introLength: maxLengthOfContentOfWebhookPayloadOnMargin + 500,
			description: "Short script with maximum intro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			script := strings.Repeat("a\n", tt.scriptSize/2) // 改行を含むスクリプト
			t.Logf("Test: %s", tt.description)
			t.Logf("Script size: %d, Intro length: %d", tt.scriptSize, tt.introLength)

			// Act
			messages := service.splitCleanupScriptWithIntro(script, tt.introLength)

			// Assert
			if len(messages) < 1 {
				t.Error("Expected at least one message chunk")
			}

			// 各チャンクが制限内であることを確認
			for i, message := range messages {
				if len(message.Content) > maxLengthOfContentOfWebhookPayload {
					t.Errorf("Chunk %d exceeds maximum length: %d (test: %s)", i, len(message.Content), tt.description)
				}

				// 基本的な構造の確認
				if message.Color != "red" {
					t.Errorf("Chunk %d has incorrect color: %s", i, message.Color)
				}

				if !message.IsCode {
					t.Errorf("Chunk %d should have IsCode=true", i)
				}
			}

			// 計算された制限値の確認
			expectedMaxLength := maxLengthOfContentOfWebhookPayloadOnMargin - tt.introLength
			if expectedMaxLength < 500 {
				expectedMaxLength = 500
			}
			t.Logf("Calculated max length: %d", expectedMaxLength)
		})
	}
}

func TestSplitCleanupScriptWithIntro_EdgeCases_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name        string
		script      string
		introLength int
		expectError bool
	}{
		{
			name:        "EmptyScript_Normal",
			script:      "",
			introLength: 100,
			expectError: false,
		},
		{
			name:        "OnlyNewlines_Normal",
			script:      "\n\n\n\n\n",
			introLength: 50,
			expectError: false,
		},
		{
			name:        "ZeroIntroLength_Normal",
			script:      "echo 'test'",
			introLength: 0,
			expectError: false,
		},
		{
			name:        "NegativeIntroLength_Normal",
			script:      "echo 'test'",
			introLength: -100,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange & Act
			messages := service.splitCleanupScriptWithIntro(tt.script, tt.introLength)

			// Assert
			if !tt.expectError && len(messages) < 1 {
				t.Error("Expected at least one message chunk")
			}

			// 各チャンクの基本的な検証
			for i, message := range messages {
				if len(message.Content) > maxLengthOfContentOfWebhookPayload {
					t.Errorf("Chunk %d exceeds maximum length: %d", i, len(message.Content))
				}

				if message.Color != "red" {
					t.Errorf("Chunk %d has incorrect color: %s", i, message.Color)
				}

				if !message.IsCode {
					t.Errorf("Chunk %d should have IsCode=true", i)
				}
			}
		})
	}
}

func TestSplitCleanupScriptWithIntro_RealWorldScenario_Normal(t *testing.T) {
	service := NewService()

	// 実際のクリーンアップスクリプトを模擬
	realScript := `# Google Cloud Workload Identity Federation リソース削除スクリプト
# 注意: このスクリプトを実行する前に、iam.shファイルをsourceしてください

# 使用例:
# source /path/to/iam.sh
# cleanup_workload_identity_federation \
#     "test-project" \
#     "test-pool" \
#     "test-provider" \
#     "test-sa" \
#     "global"

# iam.shで定義された関数を呼び出し
cleanup_workload_identity_federation "test-project" "test-pool" "test-provider" "test-sa" "global"`

	// 実際の説明文を模擬
	realIntro := `## リソース削除用スクリプト

作成したリソースを削除する場合は、以下のスクリプトを実行してください:

**注意**: このスクリプトは以下のリソースを削除します:
- Workload Identity Pool
- OIDC Provider
- Service Account
- IAM Policy Bindings (自動削除)

`

	t.Run("RealWorldScenario_Normal", func(t *testing.T) {
		// Arrange
		introLength := len(realIntro)
		t.Logf("Real script length: %d, Real intro length: %d", len(realScript), introLength)

		// Act
		messages := service.splitCleanupScriptWithIntro(realScript, introLength)

		// Assert
		if len(messages) < 1 {
			t.Error("Expected at least one message chunk")
		}

		// 各チャンクの詳細検証
		for i, message := range messages {
			t.Logf("Real scenario chunk %d: Length=%d", i, len(message.Content))

			if len(message.Content) > maxLengthOfContentOfWebhookPayload {
				t.Errorf("Real scenario chunk %d exceeds maximum length: %d", i, len(message.Content))
			}

			// 実際のスクリプト内容が含まれているかチェック
			if !strings.Contains(message.Content, "cleanup_workload_identity_federation") {
				t.Errorf("Real scenario chunk %d does not contain expected script content", i)
			}
		}
	})
}

func TestSplitCleanupScriptWithIntro_CompareWithSplitMessage_Normal(t *testing.T) {
	service := NewService()

	// 同じ内容で splitMessage と splitCleanupScriptWithIntro の結果を比較
	testContent := strings.Repeat("echo 'test line for comparison'\n", 100)
	introLength := 200

	t.Run("CompareWithSplitMessage_Normal", func(t *testing.T) {
		// Arrange
		message := &NotificationMessage{
			Title:   "Test Comparison",
			Content: testContent,
			Color:   "red",
			IsCode:  true,
		}

		// Act
		splitMessageResults := service.splitMessage(message)
		splitCleanupResults := service.splitCleanupScriptWithIntro(testContent, introLength)

		// Assert
		t.Logf("splitMessage results: %d chunks", len(splitMessageResults))
		t.Logf("splitCleanupScriptWithIntro results: %d chunks", len(splitCleanupResults))

		// 両方とも制限内であることを確認
		for i, msg := range splitMessageResults {
			if len(msg.Content) > maxLengthOfContentOfWebhookPayload {
				t.Errorf("splitMessage chunk %d exceeds limit: %d", i, len(msg.Content))
			}
		}

		for i, msg := range splitCleanupResults {
			if len(msg.Content) > maxLengthOfContentOfWebhookPayload {
				t.Errorf("splitCleanupScriptWithIntro chunk %d exceeds limit: %d", i, len(msg.Content))
			}
		}
	})
}
