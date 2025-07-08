package git

import (
	"errors"
	"strings"
	"testing"
)

// TestClient_NewClient_Normal はClient作成の正常系テスト
func TestClient_NewClient_Normal(t *testing.T) {
	// Arrange
	workingDir := "/tmp/test"

	// Act
	client := NewClient(workingDir)

	// Assert
	if client == nil {
		t.Error("Clientの作成に失敗しました")
		return
	}
	if client.workingDir != workingDir {
		t.Errorf("workingDirが期待値と異なります。期待値: %s, 実際: %s", workingDir, client.workingDir)
	}
	if client.executor == nil {
		t.Error("executorが設定されていません")
	}
}

// TestClient_NewClientWithExecutor_Normal は依存性注入Client作成の正常系テスト
func TestClient_NewClientWithExecutor_Normal(t *testing.T) {
	// Arrange
	workingDir := "/tmp/test"
	mockExecutor := NewMockGitExecutor()

	// Act
	client := NewClientWithExecutor(workingDir, mockExecutor)

	// Assert
	if client == nil {
		t.Error("Clientの作成に失敗しました")
		return
	}
	if client.workingDir != workingDir {
		t.Errorf("workingDirが期待値と異なります。期待値: %s, 実際: %s", workingDir, client.workingDir)
	}
	if client.executor != mockExecutor {
		t.Error("executorが正しく設定されていません")
	}
}

// TestClient_GetRepositoryName_Normal はGetRepositoryName正常系テスト
func TestClient_GetRepositoryName_Normal(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "/home/user/projects/test-repo\n"
	mockExecutor.SetResponse("rev-parse --show-toplevel", []byte(mockOutput))

	// Act
	repoName, err := client.GetRepositoryName()

	// Assert
	if err != nil {
		t.Errorf("GetRepositoryNameでエラーが発生しました: %v", err)
	}
	expected := "test-repo"
	if repoName != expected {
		t.Errorf("リポジトリ名が期待値と異なります。期待値: %s, 実際: %s", expected, repoName)
	}
}

// TestClient_GetRepositoryName_Error はGetRepositoryNameエラー系テスト
func TestClient_GetRepositoryName_Error(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	expectedError := errors.New("git command failed")
	mockExecutor.SetError("rev-parse --show-toplevel", expectedError)

	// Act
	_, err := client.GetRepositoryName()

	// Assert
	if err == nil {
		t.Error("エラーが発生するべきです")
	}
	if !strings.Contains(err.Error(), "リポジトリのルートディレクトリを取得できませんでした") {
		t.Errorf("期待されるエラーメッセージと異なります: %v", err)
	}
}

// TestClient_GetCurrentBranch_Normal はGetCurrentBranch正常系テスト
func TestClient_GetCurrentBranch_Normal(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "main\n"
	mockExecutor.SetResponse("rev-parse --abbrev-ref HEAD", []byte(mockOutput))

	// Act
	branch, err := client.GetCurrentBranch()

	// Assert
	if err != nil {
		t.Errorf("GetCurrentBranchでエラーが発生しました: %v", err)
	}
	expected := "main"
	if branch != expected {
		t.Errorf("ブランチ名が期待値と異なります。期待値: %s, 実際: %s", expected, branch)
	}
}

// TestClient_GetCurrentBranch_Error はGetCurrentBranchエラー系テスト
func TestClient_GetCurrentBranch_Error(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	expectedError := errors.New("git command failed")
	mockExecutor.SetError("rev-parse --abbrev-ref HEAD", expectedError)

	// Act
	_, err := client.GetCurrentBranch()

	// Assert
	if err == nil {
		t.Error("エラーが発生するべきです")
	}
	if !strings.Contains(err.Error(), "現在のブランチを取得できませんでした") {
		t.Errorf("期待されるエラーメッセージと異なります: %v", err)
	}
}

// TestClient_GetLatestCommitHash_Normal はGetLatestCommitHash正常系テスト
func TestClient_GetLatestCommitHash_Normal(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "abcdef1234567890abcdef1234567890abcdef12\n"
	mockExecutor.SetResponse("rev-parse HEAD", []byte(mockOutput))

	// Act
	hash, err := client.GetLatestCommitHash()

	// Assert
	if err != nil {
		t.Errorf("GetLatestCommitHashでエラーが発生しました: %v", err)
	}
	expected := "abcdef12" // 最初の8文字
	if hash != expected {
		t.Errorf("コミットハッシュが期待値と異なります。期待値: %s, 実際: %s", expected, hash)
	}
}

// TestClient_GetLatestCommitHash_ShortHash は短いハッシュのテスト
func TestClient_GetLatestCommitHash_ShortHash(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "abc123\n"
	mockExecutor.SetResponse("rev-parse HEAD", []byte(mockOutput))

	// Act
	hash, err := client.GetLatestCommitHash()

	// Assert
	if err != nil {
		t.Errorf("GetLatestCommitHashでエラーが発生しました: %v", err)
	}
	expected := "abc123" // 8文字未満の場合はそのまま
	if hash != expected {
		t.Errorf("コミットハッシュが期待値と異なります。期待値: %s, 実際: %s", expected, hash)
	}
}

// TestClient_GetDiff_StagedOnly はステージング済み差分取得のテスト
func TestClient_GetDiff_StagedOnly(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "diff --git a/test.txt b/test.txt\n+added line\n"
	mockExecutor.SetResponse("diff --cached", []byte(mockOutput))

	// Act
	diff, err := client.GetDiff(true)

	// Assert
	if err != nil {
		t.Errorf("GetDiffでエラーが発生しました: %v", err)
	}
	if !strings.Contains(diff, "diff --git a/test.txt b/test.txt") {
		t.Error("期待される差分内容が含まれていません")
	}
}

// TestClient_GetDiff_AllChanges は全変更差分取得のテスト
func TestClient_GetDiff_AllChanges(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "diff --git a/test.txt b/test.txt\n+added line\n-removed line\n"
	mockExecutor.SetResponse("diff HEAD", []byte(mockOutput))

	// Act
	diff, err := client.GetDiff(false)

	// Assert
	if err != nil {
		t.Errorf("GetDiffでエラーが発生しました: %v", err)
	}
	if !strings.Contains(diff, "diff --git a/test.txt b/test.txt") {
		t.Error("期待される差分内容が含まれていません")
	}
}

// TestClient_GetStatus_Normal はGetStatus正常系テスト
func TestClient_GetStatus_Normal(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "M  modified.txt\nA  added.txt\nD  deleted.txt\n?? untracked.txt\n"
	mockExecutor.SetResponse("status --porcelain", []byte(mockOutput))

	// Act
	status, err := client.GetStatus()

	// Assert
	if err != nil {
		t.Errorf("GetStatusでエラーが発生しました: %v", err)
	}
	if !strings.Contains(status, "M  modified.txt") {
		t.Error("変更ファイル情報が含まれていません")
	}
	if !strings.Contains(status, "A  added.txt") {
		t.Error("追加ファイル情報が含まれていません")
	}
}

// TestClient_GetNewFiles_StagedOnly はステージング済み新規ファイル取得のテスト
func TestClient_GetNewFiles_StagedOnly(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "A  staged_new.txt\n?? untracked.txt\nM  modified.txt\n"
	mockExecutor.SetResponse("status --porcelain", []byte(mockOutput))

	// Act
	newFiles, err := client.GetNewFiles(true)

	// Assert
	if err != nil {
		t.Errorf("GetNewFilesでエラーが発生しました: %v", err)
	}
	if len(newFiles) != 1 {
		t.Errorf("新規ファイル数が期待値と異なります。期待値: 1, 実際: %d", len(newFiles))
	}
	if newFiles[0] != "staged_new.txt" {
		t.Errorf("新規ファイル名が期待値と異なります。期待値: staged_new.txt, 実際: %s", newFiles[0])
	}
}

// TestClient_GetNewFiles_AllFiles は全新規ファイル取得のテスト
func TestClient_GetNewFiles_AllFiles(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "A  staged_new.txt\n?? untracked.txt\nM  modified.txt\n"
	mockExecutor.SetResponse("status --porcelain", []byte(mockOutput))

	// Act
	newFiles, err := client.GetNewFiles(false)

	// Assert
	if err != nil {
		t.Errorf("GetNewFilesでエラーが発生しました: %v", err)
	}
	if len(newFiles) != 2 {
		t.Errorf("新規ファイル数が期待値と異なります。期待値: 2, 実際: %d", len(newFiles))
	}

	expectedFiles := []string{"staged_new.txt", "untracked.txt"}
	for _, expected := range expectedFiles {
		found := false
		for _, actual := range newFiles {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("期待されるファイルが見つかりません: %s", expected)
		}
	}
}

// TestClient_GetDeletedFiles_Normal は削除ファイル取得のテスト
func TestClient_GetDeletedFiles_Normal(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "D  staged_deleted.txt\n D unstaged_deleted.txt\nM  modified.txt\n"
	mockExecutor.SetResponse("status --porcelain", []byte(mockOutput))

	// Act
	deletedFiles, err := client.GetDeletedFiles(false)

	// Assert
	if err != nil {
		t.Errorf("GetDeletedFilesでエラーが発生しました: %v", err)
	}
	if len(deletedFiles) != 2 {
		t.Errorf("削除ファイル数が期待値と異なります。期待値: 2, 実際: %d", len(deletedFiles))
	}
}

// TestClient_GetModifiedFilesCount_Normal は変更ファイル数取得のテスト
func TestClient_GetModifiedFilesCount_Normal(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "M  staged_modified.txt\n M unstaged_modified.txt\nMM both_modified.txt\nA  added.txt\n"
	mockExecutor.SetResponse("status --porcelain", []byte(mockOutput))

	// Act
	count, err := client.GetModifiedFilesCount(false)

	// Assert
	if err != nil {
		t.Errorf("GetModifiedFilesCountでエラーが発生しました: %v", err)
	}
	expected := 3 // M , M, MM
	if count != expected {
		t.Errorf("変更ファイル数が期待値と異なります。期待値: %d, 実際: %d", expected, count)
	}
}

// TestClient_GetModifiedFilesCount_StagedOnly はステージング済み変更ファイル数取得のテスト
func TestClient_GetModifiedFilesCount_StagedOnly(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	client := NewClientWithExecutor("/tmp/test", mockExecutor)

	mockOutput := "M  staged_modified.txt\n M unstaged_modified.txt\nMM both_modified.txt\nA  added.txt\n"
	mockExecutor.SetResponse("status --porcelain", []byte(mockOutput))

	// Act
	count, err := client.GetModifiedFilesCount(true)

	// Assert
	if err != nil {
		t.Errorf("GetModifiedFilesCountでエラーが発生しました: %v", err)
	}
	expected := 1 // M のみ
	if count != expected {
		t.Errorf("変更ファイル数が期待値と異なります。期待値: %d, 実際: %d", expected, count)
	}
}

// TestMockGitExecutor_SetResponse はモックのレスポンス設定テスト
func TestMockGitExecutor_SetResponse(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	command := "test command"
	response := []byte("test response")

	// Act
	mockExecutor.SetResponse(command, response)
	result, err := mockExecutor.Execute("/tmp", "test", "command")

	// Assert
	if err != nil {
		t.Errorf("Executeでエラーが発生しました: %v", err)
	}
	if string(result) != string(response) {
		t.Errorf("レスポンスが期待値と異なります。期待値: %s, 実際: %s", string(response), string(result))
	}
}

// TestMockGitExecutor_SetError はモックのエラー設定テスト
func TestMockGitExecutor_SetError(t *testing.T) {
	// Arrange
	mockExecutor := NewMockGitExecutor()
	command := "test command"
	expectedError := errors.New("test error")

	// Act
	mockExecutor.SetError(command, expectedError)
	_, err := mockExecutor.Execute("/tmp", "test", "command")

	// Assert
	if err != expectedError {
		t.Errorf("エラーが期待値と異なります。期待値: %v, 実際: %v", expectedError, err)
	}
}
