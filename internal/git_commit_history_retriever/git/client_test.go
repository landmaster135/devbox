package git

import (
	"fmt"
	"strings"
	"testing"
)

// #==============================================================#
// ##          Mocks                                             ##
// #==============================================================#
// MockGitExecutor はテスト用のモックGitExecutor
type MockGitExecutor struct {
	ExecuteFunc func(workingDir string, args ...string) ([]byte, error)
}

// Execute はGitコマンドを実行する（モック）
func (m *MockGitExecutor) Execute(workingDir string, args ...string) ([]byte, error) {
	return m.ExecuteFunc(workingDir, args...)
}

// #==============================================================#
// ##          Tests                                             ##
// #==============================================================#
// TestClient_ExtractCommitHashes_Normal はExtractCommitHashesの正常系テスト
func TestClient_ExtractCommitHashes_Normal(t *testing.T) {
	// Arrange
	client := NewClient("/test/repo")
	history := `* abc1234 - feat: test commit (1 hour ago) <test@example.com>
* def5678 - fix: bug fix (2 hours ago) <test@example.com>
* abc9012 - refactor: code cleanup (3 hours ago) <test@example.com>`

	// Act
	hashes := client.ExtractCommitHashes(history)

	// Debug output - 詳細な文字列解析
	t.Logf("History input:\n%s", history)
	t.Logf("History input (hex): %x", []byte(history))

	lines := strings.Split(history, "\n")
	for i, line := range lines {
		t.Logf("Line %d: %q (hex: %x)", i, line, []byte(line))
	}

	t.Logf("Found %d hashes: %v", len(hashes), hashes)

	// Assert
	expectedHashes := []string{"abc1234", "def5678", "abc9012"}
	if len(hashes) != len(expectedHashes) {
		t.Errorf("Expected %d hashes, got %d", len(expectedHashes), len(hashes))
		return // 早期リターンでpanicを防ぐ
	}

	for i, expected := range expectedHashes {
		if i >= len(hashes) || hashes[i] != expected {
			t.Errorf("Expected hash[%d] to be %s, got %s", i, expected, hashes[i])
		}
	}
}

// TestClient_ExtractCommitHashes_EmptyHistory は空の履歴の場合のテスト
func TestClient_ExtractCommitHashes_EmptyHistory(t *testing.T) {
	// Arrange
	client := NewClient("/test/repo")
	history := ""

	// Act
	hashes := client.ExtractCommitHashes(history)

	// Assert
	if len(hashes) != 0 {
		t.Errorf("Expected 0 hashes for empty history, got %d", len(hashes))
	}
}

// TestClient_ExtractCommitHashes_NoValidHashes は有効なハッシュがない場合のテスト
func TestClient_ExtractCommitHashes_NoValidHashes(t *testing.T) {
	// Arrange
	client := NewClient("/test/repo")
	history := "No commits found\nInvalid format line"

	// Act
	hashes := client.ExtractCommitHashes(history)

	// Assert
	if len(hashes) != 0 {
		t.Errorf("Expected 0 hashes for invalid format, got %d", len(hashes))
	}
}

// TestClient_ExtractCommitHashes_MixedFormat は混在フォーマットの場合のテスト
func TestClient_ExtractCommitHashes_MixedFormat(t *testing.T) {
	// Arrange
	client := NewClient("/test/repo")
	history := `* abc1234 - feat: test commit (1 hour ago) <test@example.com>
Invalid line without hash
* def5678 - fix: bug fix (2 hours ago) <test@example.com>
Another invalid line`

	// Act
	hashes := client.ExtractCommitHashes(history)

	// Assert
	expectedHashes := []string{"abc1234", "def5678"}
	if len(hashes) != len(expectedHashes) {
		t.Errorf("Expected %d hashes, got %d", len(expectedHashes), len(hashes))
	}

	for i, expected := range expectedHashes {
		if i >= len(hashes) || hashes[i] != expected {
			t.Errorf("Expected hash[%d] to be %s, got %s", i, expected, hashes[i])
		}
	}
}

// TestClient_ExtractCommitHashes_RealGitOutput は実際のgit logの出力形式のテスト
func TestClient_ExtractCommitHashes_RealGitOutput(t *testing.T) {
	// Arrange
	client := NewClient("/test/repo")
	history := `* d817ea9 - (origin/webp) chore: Updated coverage badge [skip ci] (3 hours ago) <GitHub Action>
*   a1b5281 - Merge remote-tracking branch 'refs/remotes/origin/webp' into webp (3 hours ago) <test_user>
|\
| * 362999f - chore: Updated coverage badge [skip ci] (7 hours ago) <GitHub Action>
* | 8d2ffe6 - test: add comprehensive test suite for MCP filesystem handlers (3 hours ago) <test_user>
* | 4bacc2c - refactor: implement clean architecture for filesystem MCP server (4 hours ago) <test_user>
* | 7b07288 - feat: initialize youtube_transcript module structure (5 hours ago) <test_user>
|/
* 7236624 - test: add comprehensive test suite for sequential thinking MCP server (7 hours ago) <test_user>
* f6f722f - test: add comprehensive test suite for sequential thinking service (7 hours ago) <test_user>`

	// Act
	hashes := client.ExtractCommitHashes(history)

	// Debug output
	t.Logf("Found %d hashes: %v", len(hashes), hashes)

	// Assert
	expectedHashes := []string{"d817ea9", "a1b5281", "362999f", "8d2ffe6", "4bacc2c", "7b07288", "7236624", "f6f722f"}
	if len(hashes) != len(expectedHashes) {
		t.Errorf("Expected %d hashes, got %d", len(expectedHashes), len(hashes))
		return
	}

	for i, expected := range expectedHashes {
		if i >= len(hashes) || hashes[i] != expected {
			t.Errorf("Expected hash[%d] to be %s, got %s", i, expected, hashes[i])
		}
	}
}

// TestClient_GetCommitDetails_Normal はGetCommitDetailsの正常系テスト
func TestClient_GetCommitDetails_Normal(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) >= 3 && args[0] == "show" && args[1] == "--stat" {
				return []byte("commit abc1234\nAuthor: test@example.com\nDate: test date\n\n    feat: test commit\n\n file1.go | 10 ++++++++++\n 1 file changed, 10 insertions(+)"), nil
			}
			return []byte(""), nil
		},
	}

	client := NewClientWithExecutor("/test/repo", mockExecutor)
	commitHashes := []string{"abc1234"}

	// Act
	result, err := client.GetCommitDetails(commitHashes)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, "commit abc1234") {
		t.Error("Expected result to contain commit hash")
	}

	if !strings.Contains(result, "feat: test commit") {
		t.Error("Expected result to contain commit message")
	}

	if !strings.Contains(result, "file1.go") {
		t.Error("Expected result to contain file changes")
	}
}

// TestClient_GetCommitDetails_EmptyHashes は空のハッシュリストの場合のテスト
func TestClient_GetCommitDetails_EmptyHashes(t *testing.T) {
	// Arrange
	client := NewClient("/test/repo")
	commitHashes := []string{}

	// Act
	result, err := client.GetCommitDetails(commitHashes)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "" {
		t.Errorf("Expected empty result for empty hashes, got %q", result)
	}
}

// TestClient_GetCommitDetails_GitShowError はgit showエラーの場合のテスト
func TestClient_GetCommitDetails_GitShowError(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) >= 3 && args[0] == "show" && args[1] == "--stat" {
				return nil, fmt.Errorf("git show failed")
			}
			return []byte(""), nil
		},
	}

	client := NewClientWithExecutor("/test/repo", mockExecutor)
	commitHashes := []string{"abc1234"}

	// Act
	_, err := client.GetCommitDetails(commitHashes)

	// Assert
	if err == nil {
		t.Error("Expected error for git show failure, got nil")
	}

	if !strings.Contains(err.Error(), "コミット詳細の取得に失敗しました") {
		t.Errorf("Expected error message about commit details failure, got %v", err)
	}
}

// TestClient_GetCommitDetails_MultipleHashes は複数ハッシュの場合のテスト
func TestClient_GetCommitDetails_MultipleHashes(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) >= 4 && args[0] == "show" && args[1] == "--stat" {
				// 複数のコミットハッシュが渡されることを確認
				if args[2] == "abc1234" && args[3] == "def5678" {
					return []byte("commit abc1234\nfeat: test commit\n\ncommit def5678\nfix: bug fix"), nil
				}
			}
			return []byte(""), nil
		},
	}

	client := NewClientWithExecutor("/test/repo", mockExecutor)
	commitHashes := []string{"abc1234", "def5678"}

	// Act
	result, err := client.GetCommitDetails(commitHashes)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, "abc1234") {
		t.Error("Expected result to contain first commit hash")
	}

	if !strings.Contains(result, "def5678") {
		t.Error("Expected result to contain second commit hash")
	}
}
