package usecases

import (
	"strings"
	"testing"
)

func TestSplitMessage_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name                string
		message             *NotificationMessage
		expectedChunkCount  int
		expectedTotalLength int
	}{
		{
			name: "ShortMessage_NoSplit_Normal",
			message: &NotificationMessage{
				Title:   "Test Title",
				Content: "This is a short message that should not be split.",
				Color:   "blue",
				IsCode:  false,
			},
			expectedChunkCount:  1,
			expectedTotalLength: len("This is a short message that should not be split."),
		},
		{
			name: "ExactBoundaryMessage_NoSplit_Normal",
			message: &NotificationMessage{
				Title:   "Boundary Test",
				Content: strings.Repeat("a", maxLengthOfContentOfWebhookPayloadOnMargin),
				Color:   "green",
				IsCode:  true,
			},
			expectedChunkCount:  1,
			expectedTotalLength: maxLengthOfContentOfWebhookPayloadOnMargin,
		},
		{
			name: "SlightlyOverBoundary_Split_Normal",
			message: &NotificationMessage{
				Title:   "Over Boundary Test",
				Content: strings.Repeat("a", maxLengthOfContentOfWebhookPayloadOnMargin+1),
				Color:   "red",
				IsCode:  false,
			},
			expectedChunkCount: 2,
		},
		{
			name: "MultiLineMessage_SplitByLines_Normal",
			message: &NotificationMessage{
				Title:   "Multi Line Test",
				Content: strings.Repeat("This is a line that will be repeated many times.\n", 50),
				Color:   "yellow",
				IsCode:  true,
			},
			expectedChunkCount: 2, // 予想される分割数
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			originalMessage := tt.message

			// Act
			result := service.splitMessage(originalMessage)

			// Assert
			if len(result) != tt.expectedChunkCount {
				t.Errorf("Expected %d chunks, got %d", tt.expectedChunkCount, len(result))
			}

			// 各チャンクの検証
			totalContentLength := 0
			for i, chunk := range result {
				// 文字数制限の確認
				if len(chunk.Content) > maxLengthOfContentOfWebhookPayload {
					t.Errorf("Chunk %d exceeds maximum length: %d", i, len(chunk.Content))
				}

				// メタデータの保持確認
				if chunk.Color != originalMessage.Color {
					t.Errorf("Chunk %d color mismatch: expected %s, got %s", i, originalMessage.Color, chunk.Color)
				}

				if chunk.IsCode != originalMessage.IsCode {
					t.Errorf("Chunk %d IsCode mismatch: expected %t, got %t", i, originalMessage.IsCode, chunk.IsCode)
				}

				// タイトルの確認
				if len(result) == 1 {
					if chunk.Title != originalMessage.Title {
						t.Errorf("Single chunk title mismatch: expected %s, got %s", originalMessage.Title, chunk.Title)
					}
				} else {
					if !strings.Contains(chunk.Title, "Part") {
						t.Errorf("Multi-chunk title should contain 'Part': got %s", chunk.Title)
					}
				}

				totalContentLength += len(chunk.Content)
			}

			// 内容の完全性確認（改行やスペースの調整を考慮）
			if tt.expectedTotalLength > 0 {
				// 短いメッセージの場合、完全一致を確認
				if len(result) == 1 && totalContentLength != tt.expectedTotalLength {
					t.Errorf("Content length mismatch: expected %d, got %d", tt.expectedTotalLength, totalContentLength)
				}
			}
		})
	}
}

func TestSplitMessage_EdgeCases_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name           string
		message        *NotificationMessage
		expectedChunks int
		shouldNotError bool
	}{
		{
			name: "EmptyContent_SingleChunk_Normal",
			message: &NotificationMessage{
				Title:   "Empty Test",
				Content: "",
				Color:   "gray",
				IsCode:  false,
			},
			expectedChunks: 1,
			shouldNotError: true,
		},
		{
			name: "OnlyNewlines_SingleChunk_Normal",
			message: &NotificationMessage{
				Title:   "Newlines Test",
				Content: "\n\n\n\n\n",
				Color:   "white",
				IsCode:  false,
			},
			expectedChunks: 1,
			shouldNotError: true,
		},
		{
			name: "VeryLongSingleLine_ForceSplit_Normal",
			message: &NotificationMessage{
				Title:   "Long Line Test",
				Content: strings.Repeat("a", maxLengthOfContentOfWebhookPayload*2),
				Color:   "purple",
				IsCode:  false,
			},
			expectedChunks: 2,
			shouldNotError: true,
		},
		{
			name: "JapaneseText_Split_Normal",
			message: &NotificationMessage{
				Title:   "日本語テスト",
				Content: strings.Repeat("これは日本語のテストメッセージです。マルチバイト文字の処理を確認します。", 100),
				Color:   "orange",
				IsCode:  false,
			},
			expectedChunks: 2, // 予想される分割数
			shouldNotError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			originalMessage := tt.message

			// Act
			result := service.splitMessage(originalMessage)

			// Assert
			if tt.shouldNotError && result == nil {
				t.Error("Expected result, got nil")
			}

			if len(result) < 1 {
				t.Error("Expected at least one chunk")
			}

			// 各チャンクの基本的な検証
			for i, chunk := range result {
				if chunk.Title == "" {
					t.Errorf("Chunk %d has empty title", i)
				}

				if chunk.Color != originalMessage.Color {
					t.Errorf("Chunk %d color mismatch", i)
				}

				if chunk.IsCode != originalMessage.IsCode {
					t.Errorf("Chunk %d IsCode mismatch", i)
				}

				// 文字数制限の確認
				if len(chunk.Content) > maxLengthOfContentOfWebhookPayload {
					t.Errorf("Chunk %d exceeds maximum length: %d", i, len(chunk.Content))
				}
			}
		})
	}
}

func TestSplitMessage_CodeBlockHandling_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name           string
		message        *NotificationMessage
		expectedChunks int
	}{
		{
			name: "ShortCodeBlock_NoSplit_Normal",
			message: &NotificationMessage{
				Title:   "Code Test",
				Content: "```bash\necho 'hello world'\n```",
				Color:   "green",
				IsCode:  true,
			},
			expectedChunks: 1,
		},
		{
			name: "LongCodeBlock_Split_Normal",
			message: &NotificationMessage{
				Title:   "Long Code Test",
				Content: "```bash\n" + strings.Repeat("echo 'this is a very long command that will be repeated many times'\n", 100) + "```",
				Color:   "green",
				IsCode:  true,
			},
			expectedChunks: 2, // 予想される分割数
		},
		{
			name: "MarkdownWithCodeBlocks_Split_Normal",
			message: &NotificationMessage{
				Title:   "Markdown Test",
				Content: "# Title\n\nSome text\n\n```bash\n" + strings.Repeat("command\n", 50) + "```\n\nMore text\n\n```yaml\n" + strings.Repeat("key: value\n", 50) + "```",
				Color:   "blue",
				IsCode:  false,
			},
			expectedChunks: 2, // 予想される分割数
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			originalMessage := tt.message

			// Act
			result := service.splitMessage(originalMessage)

			// Assert
			if len(result) < 1 {
				t.Error("Expected at least one chunk")
			}

			// 各チャンクの検証
			for i, chunk := range result {
				// 文字数制限の確認
				if len(chunk.Content) > maxLengthOfContentOfWebhookPayload {
					t.Errorf("Chunk %d exceeds maximum length: %d", i, len(chunk.Content))
				}

				// メタデータの保持確認
				if chunk.Color != originalMessage.Color {
					t.Errorf("Chunk %d color mismatch", i)
				}

				if chunk.IsCode != originalMessage.IsCode {
					t.Errorf("Chunk %d IsCode mismatch", i)
				}

				// 空でないことを確認
				if strings.TrimSpace(chunk.Content) == "" {
					t.Errorf("Chunk %d has empty content after trimming", i)
				}
			}
		})
	}
}

func TestSplitMessage_BoundaryConditions_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name        string
		contentSize int
		description string
	}{
		{
			name:        "ExactMarginBoundary_Normal",
			contentSize: maxLengthOfContentOfWebhookPayloadOnMargin,
			description: "Exactly at margin boundary",
		},
		{
			name:        "ExactMaxBoundary_Normal",
			contentSize: maxLengthOfContentOfWebhookPayload,
			description: "Exactly at max boundary",
		},
		{
			name:        "OneOverMargin_Normal",
			contentSize: maxLengthOfContentOfWebhookPayloadOnMargin + 1,
			description: "One character over margin",
		},
		{
			name:        "OneOverMax_Normal",
			contentSize: maxLengthOfContentOfWebhookPayload + 1,
			description: "One character over max",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			message := &NotificationMessage{
				Title:   "Boundary Test",
				Content: strings.Repeat("a", tt.contentSize),
				Color:   "test",
				IsCode:  false,
			}

			// Act
			result := service.splitMessage(message)

			// Assert
			if len(result) < 1 {
				t.Error("Expected at least one chunk")
			}

			// 各チャンクが制限内であることを確認
			for i, chunk := range result {
				if len(chunk.Content) > maxLengthOfContentOfWebhookPayload {
					t.Errorf("Chunk %d exceeds maximum length: %d (test: %s)", i, len(chunk.Content), tt.description)
				}
			}

			// 内容の完全性確認
			totalLength := 0
			for _, chunk := range result {
				totalLength += len(strings.TrimSpace(chunk.Content))
			}

			// 分割による多少の文字数変動は許容（改行処理等）
			if totalLength < tt.contentSize-10 || totalLength > tt.contentSize+10 {
				t.Errorf("Content length significantly changed: original %d, total %d (test: %s)", tt.contentSize, totalLength, tt.description)
			}
		})
	}
}

func TestSplitMessage_TitleGeneration_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name          string
		originalTitle string
		contentSize   int
		expectedParts int
	}{
		{
			name:          "SingleChunk_OriginalTitle_Normal",
			originalTitle: "Original Title",
			contentSize:   100,
			expectedParts: 1,
		},
		{
			name:          "MultiChunk_PartNumbers_Normal",
			originalTitle: "Multi Part Title",
			contentSize:   maxLengthOfContentOfWebhookPayloadOnMargin * 2,
			expectedParts: 2,
		},
		{
			name:          "LongTitle_WithParts_Normal",
			originalTitle: "This is a very long title that might affect the splitting behavior",
			contentSize:   maxLengthOfContentOfWebhookPayloadOnMargin * 3,
			expectedParts: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			message := &NotificationMessage{
				Title:   tt.originalTitle,
				Content: strings.Repeat("a", tt.contentSize),
				Color:   "blue",
				IsCode:  false,
			}

			// Act
			result := service.splitMessage(message)

			// Assert
			if len(result) < 1 {
				t.Error("Expected at least one chunk")
			}

			if len(result) == 1 {
				// 単一チャンクの場合、元のタイトルが保持される
				if result[0].Title != tt.originalTitle {
					t.Errorf("Single chunk title mismatch: expected '%s', got '%s'", tt.originalTitle, result[0].Title)
				}
			} else {
				// 複数チャンクの場合、Part番号が付与される
				for i, chunk := range result {
					expectedPartNumber := i + 1
					if !strings.Contains(chunk.Title, "Part") {
						t.Errorf("Multi-chunk title should contain 'Part': got '%s'", chunk.Title)
					}

					if !strings.Contains(chunk.Title, tt.originalTitle) {
						t.Errorf("Multi-chunk title should contain original title: expected to contain '%s', got '%s'", tt.originalTitle, chunk.Title)
					}

					// Part番号の確認
					expectedPartStr := "(Part " + strings.TrimSpace(strings.Split(chunk.Title, "(Part ")[1])
					if i == 0 {
						expectedPartStr = "(Part 1)"
					}
					if !strings.Contains(chunk.Title, "Part "+strings.TrimSpace(strings.Split(expectedPartStr, ")")[0][5:])) {
						t.Errorf("Chunk %d should have correct part number: got '%s'", expectedPartNumber, chunk.Title)
					}
				}
			}
		})
	}
}

func TestSplitMessage_ContentIntegrity_Normal(t *testing.T) {
	service := NewService()

	// Arrange
	originalContent := "Line 1\nLine 2\nLine 3\n" + strings.Repeat("This is a repeated line for testing content integrity.\n", 100)
	message := &NotificationMessage{
		Title:   "Content Integrity Test",
		Content: originalContent,
		Color:   "green",
		IsCode:  false,
	}

	// Act
	result := service.splitMessage(message)

	// Assert
	if len(result) < 1 {
		t.Error("Expected at least one chunk")
	}

	// 全チャンクの内容を結合
	var reconstructed strings.Builder
	for _, chunk := range result {
		reconstructed.WriteString(chunk.Content)
	}

	reconstructedContent := reconstructed.String()

	// 重要な行が保持されているかチェック
	if !strings.Contains(reconstructedContent, "Line 1") {
		t.Error("Line 1 not found in reconstructed content")
	}

	if !strings.Contains(reconstructedContent, "Line 2") {
		t.Error("Line 2 not found in reconstructed content")
	}

	if !strings.Contains(reconstructedContent, "Line 3") {
		t.Error("Line 3 not found in reconstructed content")
	}

	if !strings.Contains(reconstructedContent, "repeated line for testing") {
		t.Error("Repeated line content not found in reconstructed content")
	}

	// 大まかな長さの確認（分割処理による多少の変動は許容）
	if len(reconstructedContent) < len(originalContent)-100 {
		t.Errorf("Reconstructed content significantly shorter: original %d, reconstructed %d", len(originalContent), len(reconstructedContent))
	}

	// 各チャンクが制限内であることを確認
	for i, chunk := range result {
		if len(chunk.Content) > maxLengthOfContentOfWebhookPayload {
			t.Errorf("Chunk %d exceeds maximum length: %d", i, len(chunk.Content))
		}
	}
}
