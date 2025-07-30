package usecases

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	config "github.com/landmaster135/devbox/internal/memory/config"
)

// MemoryService はメモリ操作を行うサービス
type MemoryService struct {
	manager *KnowledgeGraphManager
}

// NewMemoryService は新しいMemoryServiceを作成する（ファイルベース）
func NewMemoryService(memoryFile string) *MemoryService {
	fileRepo := NewFileRepository(memoryFile)
	return &MemoryService{
		manager: NewKnowledgeGraphManager(fileRepo),
	}
}

// NewMemoryServiceWithFile はファイルベースのMemoryServiceを作成する
func NewMemoryServiceWithFile(memoryFile string) *MemoryService {
	fileRepo := NewFileRepository(memoryFile)
	return &MemoryService{
		manager: NewKnowledgeGraphManager(fileRepo),
	}
}

// NewMemoryServiceWithValkey はキーを指定してValkeyベースのMemoryServiceを作成する
func NewMemoryServiceWithValkey(valkeyURL string, key string) (*MemoryService, error) {
	valkeyRepo, err := NewValkeyRepositoryWithKey(valkeyURL, key)
	if err != nil {
		return nil, fmt.Errorf("valkeyリポジトリの作成に失敗しました: %v", err)
	}

	return &MemoryService{
		manager: NewKnowledgeGraphManager(valkeyRepo),
	}, nil
}

// NewMemoryServiceWithDependencies は依存性注入版のMemoryServiceを作成する
func NewMemoryServiceWithDependencies(fileReader config.FileReader, fileWriter config.FileWriter, memoryFile string) *MemoryService {
	fileRepo := NewFileRepositoryWithDependencies(fileReader, fileWriter, memoryFile)
	return &MemoryService{
		manager: NewKnowledgeGraphManager(fileRepo),
	}
}

// HandleCreateEntities はエンティティ作成を処理する
func (s *MemoryService) HandleCreateEntities(entitiesJSON string) (string, error) {
	if entitiesJSON == "" {
		return "", fmt.Errorf("エンティティのJSONが指定されていません")
	}

	var entities []Entity
	if err := json.Unmarshal([]byte(entitiesJSON), &entities); err != nil {
		return "", fmt.Errorf("エンティティのJSON解析エラー: %v", err)
	}

	newEntities, err := s.manager.CreateEntities(entities)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(newEntities, "", "  ")
	if err != nil {
		return "", fmt.Errorf("結果のJSON変換エラー: %v", err)
	}

	return string(result), nil
}

// HandleCreateRelations はリレーション作成を処理する
func (s *MemoryService) HandleCreateRelations(relationsJSON string) (string, error) {
	if relationsJSON == "" {
		return "", fmt.Errorf("リレーションのJSONが指定されていません")
	}

	var relations []Relation
	if err := json.Unmarshal([]byte(relationsJSON), &relations); err != nil {
		return "", fmt.Errorf("リレーションのJSON解析エラー: %v", err)
	}

	newRelations, err := s.manager.CreateRelations(relations)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(newRelations, "", "  ")
	if err != nil {
		return "", fmt.Errorf("結果のJSON変換エラー: %v", err)
	}

	return string(result), nil
}

// HandleAddObservations は観察事項追加を処理する
func (s *MemoryService) HandleAddObservations(observationsJSON string) (string, error) {
	if observationsJSON == "" {
		return "", fmt.Errorf("観察事項のJSONが指定されていません")
	}

	var observations []ObservationInput
	if err := json.Unmarshal([]byte(observationsJSON), &observations); err != nil {
		return "", fmt.Errorf("観察事項のJSON解析エラー: %v", err)
	}

	results, err := s.manager.AddObservations(observations)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", fmt.Errorf("結果のJSON変換エラー: %v", err)
	}

	return string(result), nil
}

// HandleDeleteEntities はエンティティ削除を処理する
func (s *MemoryService) HandleDeleteEntities(entityNamesStr string) error {
	if entityNamesStr == "" {
		return fmt.Errorf("削除対象のエンティティ名が指定されていません")
	}

	entityNames := strings.Split(entityNamesStr, ",")
	for i := range entityNames {
		entityNames[i] = strings.TrimSpace(entityNames[i])
	}

	return s.manager.DeleteEntities(entityNames)
}

// HandleDeleteObservations は観察事項削除を処理する
func (s *MemoryService) HandleDeleteObservations(deletionsJSON string) error {
	if deletionsJSON == "" {
		return fmt.Errorf("削除対象の観察事項のJSONが指定されていません")
	}

	var deletions []DeletionInput
	if err := json.Unmarshal([]byte(deletionsJSON), &deletions); err != nil {
		return fmt.Errorf("削除対象のJSON解析エラー: %v", err)
	}

	return s.manager.DeleteObservations(deletions)
}

// HandleDeleteRelations はリレーション削除を処理する
func (s *MemoryService) HandleDeleteRelations(relationsJSON string) error {
	if relationsJSON == "" {
		return fmt.Errorf("削除対象のリレーションのJSONが指定されていません")
	}

	var relations []Relation
	if err := json.Unmarshal([]byte(relationsJSON), &relations); err != nil {
		return fmt.Errorf("削除対象のJSON解析エラー: %v", err)
	}

	return s.manager.DeleteRelations(relations)
}

// HandleReadGraph は知識グラフ全体の読み取りを処理する
func (s *MemoryService) HandleReadGraph() (string, error) {
	graph, err := s.manager.ReadGraph()
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return "", fmt.Errorf("結果のJSON変換エラー: %v", err)
	}

	return string(result), nil
}

// HandleSearchNodes はノード検索を処理する
func (s *MemoryService) HandleSearchNodes(query string) (string, error) {
	if query == "" {
		return "", fmt.Errorf("検索クエリが指定されていません")
	}

	graph, err := s.manager.SearchNodes(query)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return "", fmt.Errorf("結果のJSON変換エラー: %v", err)
	}

	return string(result), nil
}

// HandleOpenNodes は特定ノード取得を処理する
func (s *MemoryService) HandleOpenNodes(namesStr string) (string, error) {
	if namesStr == "" {
		return "", fmt.Errorf("ノード名が指定されていません")
	}

	names := strings.Split(namesStr, ",")
	for i := range names {
		names[i] = strings.TrimSpace(names[i])
	}

	graph, err := s.manager.OpenNodes(names)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return "", fmt.Errorf("結果のJSON変換エラー: %v", err)
	}

	return string(result), nil
}

// GetDefaultMemoryFile はデフォルトのメモリファイルパスを返す
func GetDefaultMemoryFile() string {
	return filepath.Join(".", "memory.json")
}
