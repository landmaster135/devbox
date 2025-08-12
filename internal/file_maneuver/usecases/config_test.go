package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// #==============================================================#
// ##          Tests for Config                                  ##
// #==============================================================#
// TestConfigCreation はConfig作成のテストクラスです
type TestConfigCreation struct{}

// TestConfigCreation_Normal は正常なConfig作成をテストします
func (tc *TestConfigCreation) TestConfigCreation_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir1 := filepath.Join(tempDir, "src1")
	srcDir2 := filepath.Join(tempDir, "src2")
	destDir := filepath.Join(tempDir, "dest")

	// テスト用ディレクトリを作成
	os.MkdirAll(srcDir1, 0755)
	os.MkdirAll(srcDir2, 0755)
	os.MkdirAll(destDir, 0755)

	srcDirs := []string{srcDir1, srcDir2}
	extensions := []string{"jpg", "png"}

	// Act
	config, err := NewConfig(srcDirs, extensions, destDir, true, 4, false, false, false)

	// Assert
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	if len(config.SrcDirs) != 2 {
		t.Errorf("ソースディレクトリ数が期待値と異なります。期待値: 2, 実際: %d", len(config.SrcDirs))
	}

	if len(config.Extensions) != 2 {
		t.Errorf("拡張子数が期待値と異なります。期待値: 2, 実際: %d", len(config.Extensions))
	}

	// 拡張子の正規化確認
	if config.Extensions[0] != ".jpg" || config.Extensions[1] != ".png" {
		t.Errorf("拡張子の正規化が正しくありません。実際: %v", config.Extensions)
	}

	if config.DestDir != destDir {
		t.Errorf("宛先ディレクトリが期待値と異なります。期待値: %s, 実際: %s", destDir, config.DestDir)
	}

	if !config.Recursive {
		t.Error("再帰フラグが期待値と異なります")
	}

	if config.Workers != 4 {
		t.Errorf("ワーカー数が期待値と異なります。期待値: 4, 実際: %d", config.Workers)
	}
}

// 各テストクラスのインスタンスを作成してテストを実行
func TestConfigCreation_Normal(t *testing.T) {
	tc := &TestConfigCreation{}
	tc.TestConfigCreation_Normal(t)
}

// TestConfigValidation はConfig検証のテストクラスです
type TestConfigValidation struct{}

// TestConfigValidation_EmptySrcDirs は空のソースディレクトリでエラーになることをテストします
func (tc *TestConfigValidation) TestConfigValidation_EmptySrcDirs(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(destDir, 0755)

	// Act
	_, err := NewConfig([]string{}, []string{"jpg"}, destDir, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("空のソースディレクトリでエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "ソースディレクトリが指定されていません") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

func TestConfigValidation_EmptySrcDirs(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_EmptySrcDirs(t)
}

// TestConfigValidation_EmptyExtensions は空の拡張子でエラーになることをテストします
func (tc *TestConfigValidation) TestConfigValidation_EmptyExtensions(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// Act
	_, err := NewConfig([]string{srcDir}, []string{}, destDir, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("空の拡張子でエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "拡張子が指定されていません") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

func TestConfigValidation_EmptyExtensions(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_EmptyExtensions(t)
}

// TestConfigValidation_NonExistentSrcDir は存在しないソースディレクトリでエラーになることをテストします
func (tc *TestConfigValidation) TestConfigValidation_NonExistentSrcDir(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(destDir, 0755)

	nonExistentDir := filepath.Join(tempDir, "nonexistent")

	// Act
	_, err := NewConfig([]string{nonExistentDir}, []string{"jpg"}, destDir, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("存在しないソースディレクトリでエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "にアクセスできません") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

func TestConfigValidation_NonExistentSrcDir(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_NonExistentSrcDir(t)
}

// TestConfigValidation_NonExistentDestDir は存在しない宛先ディレクトリでエラーになることをテストします
func (tc *TestConfigValidation) TestConfigValidation_NonExistentDestDir(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	os.MkdirAll(srcDir, 0755)

	nonExistentDir := filepath.Join(tempDir, "nonexistent")

	// Act
	_, err := NewConfig([]string{srcDir}, []string{"jpg"}, nonExistentDir, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("存在しない宛先ディレクトリでエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "にアクセスできません") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

func TestConfigValidation_NonExistentDestDir(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_NonExistentDestDir(t)
}

// TestConfigValidation_EmptyStringInSrcDirs は空文字列のソースディレクトリでエラーになることをテストします
func (tc *TestConfigValidation) TestConfigValidation_EmptyStringInSrcDirs(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	validDir := filepath.Join(tempDir, "valid")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(validDir, 0755)
	os.MkdirAll(destDir, 0755)

	// Act
	_, err := NewConfig([]string{validDir, ""}, []string{"jpg"}, destDir, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("空文字列のソースディレクトリでエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "空のディレクトリパスが含まれています") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

func TestConfigValidation_EmptyStringInSrcDirs(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_EmptyStringInSrcDirs(t)
}

// TestConfigValidation_FileAsDirectory はファイルをディレクトリとして指定した場合のエラーをテストします
func (tc *TestConfigValidation) TestConfigValidation_FileAsDirectory(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(destDir, 0755)

	// ファイルを作成
	filePath := filepath.Join(tempDir, "notdir.txt")
	os.WriteFile(filePath, []byte("test"), 0644)

	// Act
	_, err := NewConfig([]string{filePath}, []string{"jpg"}, destDir, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("ファイルをディレクトリとして指定してもエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "はディレクトリではありません") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

func TestConfigValidation_FileAsDirectory(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_FileAsDirectory(t)
}

// TestConfigValidation_EmptyStringInExtensions は空文字列の拡張子でエラーになることをテストします
func (tc *TestConfigValidation) TestConfigValidation_EmptyStringInExtensions(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// Act
	_, err := NewConfig([]string{srcDir}, []string{"jpg", ""}, destDir, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("空文字列の拡張子でエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "空の拡張子が含まれています") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

func TestConfigValidation_EmptyStringInExtensions(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_EmptyStringInExtensions(t)
}

// TestConfigValidation_ExtensionNormalization は拡張子の正規化をテストします
func (tc *TestConfigValidation) TestConfigValidation_ExtensionNormalization(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// Act
	config, err := NewConfig([]string{srcDir}, []string{"JPG", ".png", "GIF"}, destDir, false, 1, false, false, false)

	// Assert
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	expected := []string{".jpg", ".png", ".gif"}
	if len(config.Extensions) != len(expected) {
		t.Errorf("拡張子数が期待値と異なります。期待値: %d, 実際: %d", len(expected), len(config.Extensions))
	}

	for i, ext := range expected {
		if config.Extensions[i] != ext {
			t.Errorf("拡張子[%d]が期待値と異なります。期待値: %s, 実際: %s", i, ext, config.Extensions[i])
		}
	}
}

func TestConfigValidation_ExtensionNormalization(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_ExtensionNormalization(t)
}

// TestConfigValidation_FileAsDestDir はファイルを宛先ディレクトリとして指定した場合のエラーをテストします
func (tc *TestConfigValidation) TestConfigValidation_FileAsDestDir(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	os.MkdirAll(srcDir, 0755)

	// ファイルを作成
	filePath := filepath.Join(tempDir, "notdir.txt")
	os.WriteFile(filePath, []byte("test"), 0644)

	// Act
	_, err := NewConfig([]string{srcDir}, []string{"jpg"}, filePath, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("ファイルを宛先ディレクトリとして指定してもエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "はディレクトリではありません") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

func TestConfigValidation_FileAsDestDir(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_FileAsDestDir(t)
}
