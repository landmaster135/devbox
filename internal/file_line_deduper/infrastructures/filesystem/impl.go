package filesystem

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/landmaster135/devbox/internal/file_line_deduper/domain/models"
)

// osRepository は Repository のOS実装です。
type osRepository struct{}

// ReadFile はファイルを読み込み、FileContentを返します。
func (r *osRepository) ReadFile(path string) (*models.FileContent, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ファイルを開けませんでした: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ファイルの読み込み中にエラーが発生しました: %w", err)
	}

	return models.NewFileContent(lines), nil
}

// WriteFile はFileContentをファイルに書き込みます。
func (r *osRepository) WriteFile(path string, content *models.FileContent) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("出力ファイルを作成できませんでした: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range content.Lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("ファイルへの書き込み中にエラーが発生しました: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("ファイルへの書き込み中にエラーが発生しました: %w", err)
	}

	return nil
}

// FileExists はファイルが存在するかどうかを確認します。
func (r *osRepository) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FindFilesByExt はディレクトリ内の指定された拡張子のファイルのパスリストを返します。
func (r *osRepository) FindFilesByExt(dirPath, ext string) ([]string, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("指定されたディレクトリが存在しません: %s", dirPath)
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの読み込みに失敗しました: %w", err)
	}

	matchedFiles := []string{}
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ext {
			matchedFiles = append(matchedFiles, filepath.Join(dirPath, file.Name()))
		}
	}

	return matchedFiles, nil
}

// HasFilesWithExt はディレクトリ内に指定された拡張子のファイルが存在するかどうかを確認します。
func (r *osRepository) HasFilesWithExt(dirPath, ext string) (bool, error) {
	files, err := r.FindFilesByExt(dirPath, ext)
	if err != nil {
		return false, err
	}
	return len(files) > 0, nil
}

// ReadJSONFile はJSONファイルを読み込み、インターフェース型として返します。
func (r *osRepository) ReadJSONFile(path string) (interface{}, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	if len(fileData) >= 3 && fileData[0] == 0xEF && fileData[1] == 0xBB && fileData[2] == 0xBF {
		fileData = fileData[3:]
	}

	var jsonData interface{}
	if err := json.Unmarshal(fileData, &jsonData); err != nil {
		return nil, fmt.Errorf("JSONデータのパースに失敗しました: %w", err)
	}

	return jsonData, nil
}

// GetDirectoryPath はファイルパスからディレクトリパスを取得します。
func (r *osRepository) GetDirectoryPath(path string) string {
	return filepath.Dir(path)
}

// CreateDirectory はディレクトリを作成します。
func (r *osRepository) CreateDirectory(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("ディレクトリの作成に失敗しました: %w", err)
		}
	}
	return nil
}

// ReadDir はディレクトリ内のファイルとサブディレクトリのエントリを返します。
func (r *osRepository) ReadDir(dirPath string) ([]*models.DirEntry, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("指定されたディレクトリが存在しません: %s", dirPath)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの読み込みに失敗しました: %w", err)
	}

	result := make([]*models.DirEntry, 0, len(entries))
	for _, entry := range entries {
		path := filepath.Join(dirPath, entry.Name())
		result = append(result, models.NewDirEntry(entry.Name(), entry.IsDir(), path))
	}

	return result, nil
}
