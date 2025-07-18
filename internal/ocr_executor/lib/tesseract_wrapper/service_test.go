package tesseract_wrapper

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

// MockCommandExecutor はテスト用のモック実装
type MockCommandExecutor struct {
	lookPathFunc func(file string) (string, error)
	commandFunc  func(name string, arg ...string) *exec.Cmd
}

func (m *MockCommandExecutor) LookPath(file string) (string, error) {
	if m.lookPathFunc != nil {
		return m.lookPathFunc(file)
	}
	return file, nil // デフォルトは成功
}

func (m *MockCommandExecutor) Command(name string, arg ...string) *exec.Cmd {
	if m.commandFunc != nil {
		return m.commandFunc(name, arg...)
	}
	// デフォルトは何もしないコマンド
	return exec.Command("echo", "mock output")
}

// MockCommand はテスト用のコマンドモック
type MockCommand struct {
	output []byte
	err    error
}

func (m *MockCommand) Output() ([]byte, error) {
	return m.output, m.err
}

// createMockCmd はモックコマンドを作成するヘルパー関数
func createMockCmd(output string, err error) *exec.Cmd {
	cmd := &exec.Cmd{}
	// 実際のコマンドの代わりにモックの出力を返すように設定
	if err == nil {
		cmd.Stdout = bytes.NewBufferString(output)
	}
	return cmd
}

// TestTesseractClient はTesseractClientのテストクラス
type TestTesseractClient struct {
	t *testing.T
}

func TestNewTesseractClient_Normal(t *testing.T) {
	mockExecutor := &MockCommandExecutor{
		lookPathFunc: func(file string) (string, error) {
			return "/usr/bin/tesseract", nil
		},
	}

	client, err := NewTesseractClient(WithExecutor(mockExecutor))

	if err != nil {
		t.Errorf("NewTesseractClient() エラーが発生しました: %v", err)
	}
	if client == nil {
		t.Error("NewTesseractClient() クライアントがnilです")
		return
	}
	if client.binPath != "tesseract" {
		t.Errorf("NewTesseractClient() binPath = %v, want %v", client.binPath, "tesseract")
	}
}

func TestNewTesseractClient_WithBinPath(t *testing.T) {
	customPath := "/custom/path/tesseract"
	mockExecutor := &MockCommandExecutor{
		lookPathFunc: func(file string) (string, error) {
			if file == customPath {
				return customPath, nil
			}
			return "", errors.New("not found")
		},
	}

	client, err := NewTesseractClient(
		WithBinPath(customPath),
		WithExecutor(mockExecutor),
	)

	if err != nil {
		t.Errorf("NewTesseractClient() エラーが発生しました: %v", err)
	}
	if client.binPath != customPath {
		t.Errorf("NewTesseractClient() binPath = %v, want %v", client.binPath, customPath)
	}
}

func TestNewTesseractClient_WithLanguages(t *testing.T) {
	languages := "jpn+eng"
	mockExecutor := &MockCommandExecutor{
		lookPathFunc: func(file string) (string, error) {
			return "/usr/bin/tesseract", nil
		},
	}

	client, err := NewTesseractClient(
		WithLanguages(languages),
		WithExecutor(mockExecutor),
	)

	if err != nil {
		t.Errorf("NewTesseractClient() エラーが発生しました: %v", err)
	}
	if client.languages != languages {
		t.Errorf("NewTesseractClient() languages = %v, want %v", client.languages, languages)
	}
}

func TestNewTesseractClient_WithMultipleOptions(t *testing.T) {
	customPath := "/custom/tesseract"
	languages := "jpn+eng+fra"
	mockExecutor := &MockCommandExecutor{
		lookPathFunc: func(file string) (string, error) {
			return customPath, nil
		},
	}

	client, err := NewTesseractClient(
		WithBinPath(customPath),
		WithLanguages(languages),
		WithExecutor(mockExecutor),
	)

	if err != nil {
		t.Errorf("NewTesseractClient() エラーが発生しました: %v", err)
	}
	if client.binPath != customPath {
		t.Errorf("NewTesseractClient() binPath = %v, want %v", client.binPath, customPath)
	}
	if client.languages != languages {
		t.Errorf("NewTesseractClient() languages = %v, want %v", client.languages, languages)
	}
}

func TestNewTesseractClient_BinaryNotFound(t *testing.T) {
	mockExecutor := &MockCommandExecutor{
		lookPathFunc: func(file string) (string, error) {
			return "", errors.New("tesseract not found")
		},
	}

	client, err := NewTesseractClient(WithExecutor(mockExecutor))

	if err == nil {
		t.Error("NewTesseractClient() エラーが期待されましたが、nilが返されました")
	}
	if client != nil {
		t.Error("NewTesseractClient() クライアントはnilであるべきです")
	}
	if !strings.Contains(err.Error(), "tesseractバイナリが見つかりません") {
		t.Errorf("NewTesseractClient() 予期しないエラーメッセージ: %v", err)
	}
}

func TestTextFromImageFile_Normal(t *testing.T) {
	expectedOutput := "Hello World"
	mockExecutor := &MockCommandExecutor{
		lookPathFunc: func(file string) (string, error) {
			return "/usr/bin/tesseract", nil
		},
		commandFunc: func(name string, arg ...string) *exec.Cmd {
			// 引数の検証
			expectedArgs := []string{"test.jpg", "stdout"}
			if len(arg) != len(expectedArgs) {
				t.Errorf("Command() 引数の数が違います: got %d, want %d", len(arg), len(expectedArgs))
			}
			for i, expected := range expectedArgs {
				if i < len(arg) && arg[i] != expected {
					t.Errorf("Command() 引数[%d] = %v, want %v", i, arg[i], expected)
				}
			}

			// モックコマンドを作成
			cmd := exec.Command("echo", expectedOutput)
			return cmd
		},
	}

	client, err := NewTesseractClient(WithExecutor(mockExecutor))
	if err != nil {
		t.Fatalf("NewTesseractClient() エラー: %v", err)
	}

	result, err := client.TextFromImageFile("test.jpg")

	if err != nil {
		t.Errorf("TextFromImageFile() エラーが発生しました: %v", err)
	}
	if result != expectedOutput {
		t.Errorf("TextFromImageFile() = %v, want %v", result, expectedOutput)
	}
}

func TestTextFromImageFile_WithLanguages(t *testing.T) {
	expectedOutput := "こんにちは世界"
	languages := "jpn+eng"
	mockExecutor := &MockCommandExecutor{
		lookPathFunc: func(file string) (string, error) {
			return "/usr/bin/tesseract", nil
		},
		commandFunc: func(name string, arg ...string) *exec.Cmd {
			// 引数の検証（言語設定を含む）
			expectedArgs := []string{"test.jpg", "stdout", "-l", languages}
			if len(arg) != len(expectedArgs) {
				t.Errorf("Command() 引数の数が違います: got %d, want %d", len(arg), len(expectedArgs))
			}
			for i, expected := range expectedArgs {
				if i < len(arg) && arg[i] != expected {
					t.Errorf("Command() 引数[%d] = %v, want %v", i, arg[i], expected)
				}
			}

			cmd := exec.Command("echo", expectedOutput)
			return cmd
		},
	}

	client, err := NewTesseractClient(
		WithLanguages(languages),
		WithExecutor(mockExecutor),
	)
	if err != nil {
		t.Fatalf("NewTesseractClient() エラー: %v", err)
	}

	result, err := client.TextFromImageFile("test.jpg")

	if err != nil {
		t.Errorf("TextFromImageFile() エラーが発生しました: %v", err)
	}
	if result != expectedOutput {
		t.Errorf("TextFromImageFile() = %v, want %v", result, expectedOutput)
	}
}

func TestTextFromImageFile_TrimWhitespace(t *testing.T) {
	outputWithWhitespace := "  Hello World  \n\t"
	expectedOutput := "Hello World"
	mockExecutor := &MockCommandExecutor{
		lookPathFunc: func(file string) (string, error) {
			return "/usr/bin/tesseract", nil
		},
		commandFunc: func(name string, arg ...string) *exec.Cmd {
			cmd := exec.Command("echo", outputWithWhitespace)
			return cmd
		},
	}

	client, err := NewTesseractClient(WithExecutor(mockExecutor))
	if err != nil {
		t.Fatalf("NewTesseractClient() エラー: %v", err)
	}

	result, err := client.TextFromImageFile("test.jpg")

	if err != nil {
		t.Errorf("TextFromImageFile() エラーが発生しました: %v", err)
	}
	if result != expectedOutput {
		t.Errorf("TextFromImageFile() = %q, want %q", result, expectedOutput)
	}
}

func TestTextFromImageFile_CommandError(t *testing.T) {
	mockExecutor := &MockCommandExecutor{
		lookPathFunc: func(file string) (string, error) {
			return "/usr/bin/tesseract", nil
		},
		commandFunc: func(name string, arg ...string) *exec.Cmd {
			// 失敗するコマンドを返す
			cmd := exec.Command("false") // falseコマンドは常に失敗する
			return cmd
		},
	}

	client, err := NewTesseractClient(WithExecutor(mockExecutor))
	if err != nil {
		t.Fatalf("NewTesseractClient() エラー: %v", err)
	}

	result, err := client.TextFromImageFile("test.jpg")

	if err == nil {
		t.Error("TextFromImageFile() エラーが期待されましたが、nilが返されました")
	}
	if result != "" {
		t.Errorf("TextFromImageFile() = %v, want empty string", result)
	}
	if !strings.Contains(err.Error(), "tesseract実行エラー") {
		t.Errorf("TextFromImageFile() 予期しないエラーメッセージ: %v", err)
	}
}
