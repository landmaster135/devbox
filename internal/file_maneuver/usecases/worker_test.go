package usecases

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// #==============================================================#
// ##          Tests for Worker                                  ##
// #==============================================================#
// TestWorkerNormalization はワーカー数正規化のテストクラスです
type TestWorkerNormalization struct{}

// TestWorkerNormalization_ZeroWorkers はワーカー数0の場合の正規化をテストします
func (tc *TestWorkerNormalization) TestWorkerNormalization_ZeroWorkers(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// Act
	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, destDir, false, 0, false, false, false)

	// Assert
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	expectedWorkers := runtime.NumCPU()
	if config.Workers != expectedWorkers {
		t.Errorf("ワーカー数が期待値と異なります。期待値: %d, 実際: %d", expectedWorkers, config.Workers)
	}
}

func TestWorkerNormalization_ZeroWorkers(t *testing.T) {
	tc := &TestWorkerNormalization{}
	tc.TestWorkerNormalization_ZeroWorkers(t)
}

// TestWorkerNormalization_ExcessiveWorkers は過剰なワーカー数の場合の正規化をテストします
func (tc *TestWorkerNormalization) TestWorkerNormalization_ExcessiveWorkers(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	excessiveWorkers := runtime.NumCPU()*2 + 10

	// Act
	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, destDir, false, excessiveWorkers, false, false, false)

	// Assert
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	expectedMaxWorkers := runtime.NumCPU() * 2
	if config.Workers != expectedMaxWorkers {
		t.Errorf("ワーカー数が期待値と異なります。期待値: %d, 実際: %d", expectedMaxWorkers, config.Workers)
	}
}

func TestWorkerNormalization_ExcessiveWorkers(t *testing.T) {
	tc := &TestWorkerNormalization{}
	tc.TestWorkerNormalization_ExcessiveWorkers(t)
}
