package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewExifMirrorService(t *testing.T) {
	service := NewExifMirrorService()
	if service == nil {
		t.Fatal("NewExifMirrorService() returned nil")
	}
}

func TestValidateDirectory(t *testing.T) {
	tests := []struct {
		name    string
		dirPath string
		wantErr bool
	}{
		{
			name:    "empty path",
			dirPath: "",
			wantErr: true,
		},
		{
			name:    "non-existent directory",
			dirPath: "/non/existent/path",
			wantErr: true,
		},
		{
			name:    "current directory",
			dirPath: ".",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDirectory(tt.dirPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateExtension(t *testing.T) {
	tests := []struct {
		name    string
		ext     string
		wantErr bool
	}{
		{
			name:    "empty extension",
			ext:     "",
			wantErr: true,
		},
		{
			name:    "valid jpg extension",
			ext:     "jpg",
			wantErr: false,
		},
		{
			name:    "valid jpg extension with dot",
			ext:     ".jpg",
			wantErr: false,
		},
		{
			name:    "valid png extension",
			ext:     "png",
			wantErr: false,
		},
		{
			name:    "valid webp extension",
			ext:     "webp",
			wantErr: false,
		},
		{
			name:    "invalid extension",
			ext:     "xyz",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExtension(tt.ext)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateExtension() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemoveExtension(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{
			name:     "file with jpg extension",
			fileName: "image.jpg",
			want:     "image",
		},
		{
			name:     "file with png extension",
			fileName: "photo.png",
			want:     "photo",
		},
		{
			name:     "file without extension",
			fileName: "noextension",
			want:     "noextension",
		},
		{
			name:     "file with multiple dots",
			fileName: "image.backup.jpg",
			want:     "image.backup",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeExtension(tt.fileName)
			if got != tt.want {
				t.Errorf("removeExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExifMirrorService_isTargetFile(t *testing.T) {
	service := NewExifMirrorService()

	tests := []struct {
		name            string
		filePath        string
		targetExtension string
		want            bool
	}{
		{
			name:            "jpg file with jpg target",
			filePath:        "/path/to/image.jpg",
			targetExtension: "jpg",
			want:            true,
		},
		{
			name:            "jpg file with png target",
			filePath:        "/path/to/image.jpg",
			targetExtension: "png",
			want:            false,
		},
		{
			name:            "png file with png target",
			filePath:        "/path/to/image.png",
			targetExtension: "png",
			want:            true,
		},
		{
			name:            "jpg file with dot extension",
			filePath:        "/path/to/image.jpg",
			targetExtension: ".jpg",
			want:            true,
		},
		{
			name:            "unsupported extension",
			filePath:        "/path/to/file.xyz",
			targetExtension: "",
			want:            false,
		},
		{
			name:            "webp file with empty target",
			filePath:        "/path/to/image.webp",
			targetExtension: "",
			want:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.isTargetFile(tt.filePath, tt.targetExtension)
			if got != tt.want {
				t.Errorf("ExifMirrorService.isTargetFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExifMirrorService_findCorrespondingSourceFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_mirror_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のファイル構造を作成
	sourceDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")

	err = os.MkdirAll(sourceDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// テスト用のソースファイルを作成
	sourceFile := filepath.Join(sourceDir, "test.jpg")
	err = os.WriteFile(sourceFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	service := NewExifMirrorService()
	config := &Config{
		SourceFolderPath: sourceDir,
		TargetFolderPath: targetDir,
		SourceExtension:  "jpg",
		TargetExtension:  "webp",
	}

	targetFilePath := filepath.Join(targetDir, "test.webp")

	result := service.findCorrespondingSourceFile(targetFilePath, config)
	expected := sourceFile

	if result != expected {
		t.Errorf("findCorrespondingSourceFile() = %v, want %v", result, expected)
	}
}

func TestExifMirrorService_findCorrespondingSourceFile_NotFound(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_mirror_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	targetDir := filepath.Join(tempDir, "target")

	err = os.MkdirAll(sourceDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	service := NewExifMirrorService()
	config := &Config{
		SourceFolderPath: sourceDir,
		TargetFolderPath: targetDir,
		SourceExtension:  "jpg",
		TargetExtension:  "webp",
	}

	// 存在しないターゲットファイル
	targetFilePath := filepath.Join(targetDir, "nonexistent.webp")

	result := service.findCorrespondingSourceFile(targetFilePath, config)

	if result != "" {
		t.Errorf("findCorrespondingSourceFile() = %v, want empty string", result)
	}
}

func TestExifMirrorService_hasExifTool(t *testing.T) {
	service := NewExifMirrorService()
	
	// このテストは環境に依存するため、エラーが出ないことだけ確認
	result := service.hasExifTool()
	
	// exiftoolがあるかないかは環境次第なので、実行されることだけ確認
	_ = result
}

func TestExifMirrorService_CopyFileExifSimple(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_mirror_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のファイルを作成
	sourceFile := filepath.Join(tempDir, "source.jpg")
	targetFile := filepath.Join(tempDir, "target.jpg")

	err = os.WriteFile(sourceFile, []byte("source content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(targetFile, []byte("target content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	service := NewExifMirrorService()
	
	// 基本的なファイル時刻コピーのテスト
	err = service.CopyFileExifSimple(sourceFile, targetFile)
	if err != nil {
		t.Errorf("CopyFileExifSimple() error = %v", err)
	}
}

func TestExifMirrorService_BackupFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_mirror_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のファイルを作成
	originalFile := filepath.Join(tempDir, "original.jpg")
	testContent := "test content for backup"
	
	err = os.WriteFile(originalFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	service := NewExifMirrorService()
	
	// バックアップ作成
	backupPath, err := service.BackupFile(originalFile)
	if err != nil {
		t.Fatalf("BackupFile() error = %v", err)
	}

	// バックアップファイルが作成されたか確認
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Errorf("Backup file was not created: %s", backupPath)
	}

	// バックアップファイルの内容が正しいか確認
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("Failed to read backup file: %v", err)
	}

	if string(backupContent) != testContent {
		t.Errorf("Backup content = %s, want %s", string(backupContent), testContent)
	}
}

