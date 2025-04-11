package repositories

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/domain/models"
	domainRepo "github.com/landmaster135/devbox/internal/domain/repositories"
)

// JSONRepositoryImpl はJSONRepositoryインターフェースの実装です
type JSONRepositoryImpl struct {
	fileRepo    domainRepo.FileRepository
	iso8601Repo domainRepo.ISO8601Repository
}

// NewJSONRepository は新しいJSONRepositoryImplインスタンスを作成します
func NewJSONRepository(fileRepo domainRepo.FileRepository, iso8601Repo domainRepo.ISO8601Repository) domainRepo.JSONRepository {
	return &JSONRepositoryImpl{
		fileRepo:    fileRepo,
		iso8601Repo: iso8601Repo,
	}
}

// FindJSONFiles はディレクトリ内のJSONファイルを検索します
func (r *JSONRepositoryImpl) FindJSONFiles(dirPath string, recursive bool) ([]string, error) {
	// 再帰的に検索する場合は、独自の実装が必要
	if recursive {
		var jsonFiles []string

		// ディレクトリ内のファイルとサブディレクトリを取得
		entries, err := r.fileRepo.ReadDir(dirPath)
		if err != nil {
			return nil, fmt.Errorf("ディレクトリの読み込みに失敗しました: %w", err)
		}

		// 各エントリを処理
		for _, entry := range entries {
			if entry.IsDir {
				// サブディレクトリの場合、再帰的に処理
				subFiles, err := r.FindJSONFiles(entry.Path, recursive)
				if err != nil {
					return nil, err
				}
				jsonFiles = append(jsonFiles, subFiles...)
			} else {
				// ファイルの場合、拡張子がJSONかどうかを確認
				if strings.HasSuffix(strings.ToLower(entry.Name), ".json") {
					jsonFiles = append(jsonFiles, entry.Path)
				}
			}
		}

		return jsonFiles, nil
	} else {
		// 非再帰的な場合は、FindFilesByExtを使用
		return r.fileRepo.FindFilesByExt(dirPath, ".json")
	}
}

// ConvertFile は単一のJSONファイルを処理します
func (r *JSONRepositoryImpl) ConvertFile(filePath, key string, dryRun bool) (bool, error) {
	// JSONファイルを読み込む
	jsonData, err := r.fileRepo.ReadJSONFile(filePath)
	if err != nil {
		return false, fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	// 型アサーション
	data, ok := jsonData.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("JSONデータをマップに変換できません")
	}

	// JSONデータを再帰的に処理
	newData, changed := r.ProcessJSONData(data, key)

	// 変換が行われた場合
	if changed {
		// ドライランでない場合は変更を保存
		if !dryRun {
			// JSONに変換
			jsonBytes, err := json.MarshalIndent(newData, "", "  ")
			if err != nil {
				return false, fmt.Errorf("JSONへの変換に失敗しました: %w", err)
			}

			// ファイルに書き込む
			lines := strings.Split(string(jsonBytes), "\n")
			content := models.NewFileContent(lines)
			if err := r.fileRepo.WriteFile(filePath, content); err != nil {
				return false, fmt.Errorf("ファイルの書き込みに失敗しました: %w", err)
			}
		}
	}

	return changed, nil
}

// ProcessJSONData はJSONデータを再帰的に処理します
func (r *JSONRepositoryImpl) ProcessJSONData(data interface{}, targetKey string) (interface{}, bool) {
	// 変更があったかどうかを示すフラグ
	changed := false

	switch v := data.(type) {
	case map[string]interface{}:
		// オブジェクトの場合、各キーを処理
		for k, val := range v {
			if k == targetKey {
				// ターゲットキーの場合、値を変換
				if strVal, ok := val.(string); ok {
					// ISO8601Repositoryを使用して変換
					timestamp, err := r.iso8601Repo.ParseISO8601(strVal)
					if err == nil {
						v[k] = timestamp
						changed = true
					}
				}
			} else {
				// 再帰的に処理
				newVal, childChanged := r.ProcessJSONData(val, targetKey)
				if childChanged {
					v[k] = newVal
					changed = true
				}
			}
		}
		return v, changed

	case []interface{}:
		// 配列の場合、各要素を処理
		for i, val := range v {
			newVal, childChanged := r.ProcessJSONData(val, targetKey)
			if childChanged {
				v[i] = newVal
				changed = true
			}
		}
		return v, changed

	default:
		// その他の型はそのまま返す
		return data, false
	}
}
