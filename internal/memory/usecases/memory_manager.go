package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	config "github.com/landmaster135/devbox/internal/memory/config"
)

// KnowledgeGraphManager は知識グラフの管理を行う構造体
type KnowledgeGraphManager struct {
	fileReader config.FileReader
	fileWriter config.FileWriter
	memoryFile string
}

// NewKnowledgeGraphManager は新しいKnowledgeGraphManagerを作成する
func NewKnowledgeGraphManager(memoryFile string) *KnowledgeGraphManager {
	return &KnowledgeGraphManager{
		fileReader: &config.StandardFileReader{},
		fileWriter: &config.StandardFileWriter{},
		memoryFile: memoryFile,
	}
}

// NewKnowledgeGraphManagerWithDependencies は依存性注入版のKnowledgeGraphManagerを作成する
func NewKnowledgeGraphManagerWithDependencies(fileReader config.FileReader, fileWriter config.FileWriter, memoryFile string) *KnowledgeGraphManager {
	return &KnowledgeGraphManager{
		fileReader: fileReader,
		fileWriter: fileWriter,
		memoryFile: memoryFile,
	}
}

// loadGraph は知識グラフをファイルから読み込む
func (m *KnowledgeGraphManager) loadGraph() (*KnowledgeGraph, error) {
	// ファイルが存在しない場合は空のグラフを返す
	if _, err := os.Stat(m.memoryFile); os.IsNotExist(err) {
		return &KnowledgeGraph{
			Entities:  []Entity{},
			Relations: []Relation{},
		}, nil
	}

	data, err := m.fileReader.ReadFile(m.memoryFile)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込みエラー: %v", err)
	}

	var graph KnowledgeGraph
	if err := json.Unmarshal(data, &graph); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %v", err)
	}

	return &graph, nil
}

// saveGraph は知識グラフをファイルに保存する
func (m *KnowledgeGraphManager) saveGraph(graph *KnowledgeGraph) error {
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(m.memoryFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("ディレクトリ作成エラー: %v", err)
	}

	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON変換エラー: %v", err)
	}

	if err := m.fileWriter.WriteFile(m.memoryFile, data, 0644); err != nil {
		return fmt.Errorf("ファイル書き込みエラー: %v", err)
	}

	return nil
}

// CreateEntities は複数のエンティティを作成する
func (m *KnowledgeGraphManager) CreateEntities(entities []Entity) ([]Entity, error) {
	graph, err := m.loadGraph()
	if err != nil {
		return nil, err
	}

	var newEntities []Entity
	for _, entity := range entities {
		// 既存のエンティティ名をチェック
		exists := false
		for _, existingEntity := range graph.Entities {
			if existingEntity.Name == entity.Name {
				exists = true
				break
			}
		}

		if !exists {
			graph.Entities = append(graph.Entities, entity)
			newEntities = append(newEntities, entity)
		}
	}

	if err := m.saveGraph(graph); err != nil {
		return nil, err
	}

	return newEntities, nil
}

// CreateRelations は複数のリレーションを作成する
func (m *KnowledgeGraphManager) CreateRelations(relations []Relation) ([]Relation, error) {
	graph, err := m.loadGraph()
	if err != nil {
		return nil, err
	}

	var newRelations []Relation
	for _, relation := range relations {
		// 重複チェック
		exists := false
		for _, existingRelation := range graph.Relations {
			if existingRelation.From == relation.From &&
				existingRelation.To == relation.To &&
				existingRelation.RelationType == relation.RelationType {
				exists = true
				break
			}
		}

		if !exists {
			graph.Relations = append(graph.Relations, relation)
			newRelations = append(newRelations, relation)
		}
	}

	if err := m.saveGraph(graph); err != nil {
		return nil, err
	}

	return newRelations, nil
}

// AddObservations は既存のエンティティに観察事項を追加する
func (m *KnowledgeGraphManager) AddObservations(observations []ObservationInput) ([]ObservationResult, error) {
	graph, err := m.loadGraph()
	if err != nil {
		return nil, err
	}

	var results []ObservationResult
	for _, obs := range observations {
		// エンティティを検索
		var targetEntity *Entity
		for i := range graph.Entities {
			if graph.Entities[i].Name == obs.EntityName {
				targetEntity = &graph.Entities[i]
				break
			}
		}

		if targetEntity == nil {
			return nil, fmt.Errorf("エンティティが見つかりません: %s", obs.EntityName)
		}

		// 新しい観察事項を追加（重複チェック）
		var addedObservations []string
		for _, content := range obs.Contents {
			exists := false
			for _, existingObs := range targetEntity.Observations {
				if existingObs == content {
					exists = true
					break
				}
			}

			if !exists {
				targetEntity.Observations = append(targetEntity.Observations, content)
				addedObservations = append(addedObservations, content)
			}
		}

		results = append(results, ObservationResult{
			EntityName:        obs.EntityName,
			AddedObservations: addedObservations,
		})
	}

	if err := m.saveGraph(graph); err != nil {
		return nil, err
	}

	return results, nil
}

// DeleteEntities はエンティティとその関連リレーションを削除する
func (m *KnowledgeGraphManager) DeleteEntities(entityNames []string) error {
	graph, err := m.loadGraph()
	if err != nil {
		return err
	}

	// エンティティを削除
	var filteredEntities []Entity
	for _, entity := range graph.Entities {
		shouldDelete := false
		for _, name := range entityNames {
			if entity.Name == name {
				shouldDelete = true
				break
			}
		}
		if !shouldDelete {
			filteredEntities = append(filteredEntities, entity)
		}
	}
	graph.Entities = filteredEntities

	// 関連するリレーションを削除
	var filteredRelations []Relation
	for _, relation := range graph.Relations {
		shouldDelete := false
		for _, name := range entityNames {
			if relation.From == name || relation.To == name {
				shouldDelete = true
				break
			}
		}
		if !shouldDelete {
			filteredRelations = append(filteredRelations, relation)
		}
	}
	graph.Relations = filteredRelations

	return m.saveGraph(graph)
}

// DeleteObservations は特定の観察事項を削除する
func (m *KnowledgeGraphManager) DeleteObservations(deletions []DeletionInput) error {
	graph, err := m.loadGraph()
	if err != nil {
		return err
	}

	for _, deletion := range deletions {
		// エンティティを検索
		for i := range graph.Entities {
			if graph.Entities[i].Name == deletion.EntityName {
				// 指定された観察事項を削除
				var filteredObservations []string
				for _, obs := range graph.Entities[i].Observations {
					shouldDelete := false
					for _, delObs := range deletion.Observations {
						if obs == delObs {
							shouldDelete = true
							break
						}
					}
					if !shouldDelete {
						filteredObservations = append(filteredObservations, obs)
					}
				}
				graph.Entities[i].Observations = filteredObservations
				break
			}
		}
	}

	return m.saveGraph(graph)
}

// DeleteRelations は特定のリレーションを削除する
func (m *KnowledgeGraphManager) DeleteRelations(relations []Relation) error {
	graph, err := m.loadGraph()
	if err != nil {
		return err
	}

	// 指定されたリレーションを削除
	var filteredRelations []Relation
	for _, relation := range graph.Relations {
		shouldDelete := false
		for _, delRelation := range relations {
			if relation.From == delRelation.From &&
				relation.To == delRelation.To &&
				relation.RelationType == delRelation.RelationType {
				shouldDelete = true
				break
			}
		}
		if !shouldDelete {
			filteredRelations = append(filteredRelations, relation)
		}
	}
	graph.Relations = filteredRelations

	return m.saveGraph(graph)
}

// ReadGraph は知識グラフ全体を読み取る
func (m *KnowledgeGraphManager) ReadGraph() (*KnowledgeGraph, error) {
	return m.loadGraph()
}

// SearchNodes はクエリに基づいてノードを検索する
func (m *KnowledgeGraphManager) SearchNodes(query string) (*KnowledgeGraph, error) {
	graph, err := m.loadGraph()
	if err != nil {
		return nil, err
	}

	queryLower := strings.ToLower(query)

	// エンティティをフィルタリング
	var filteredEntities []Entity
	for _, entity := range graph.Entities {
		match := false

		// 名前での検索
		if strings.Contains(strings.ToLower(entity.Name), queryLower) {
			match = true
		}

		// エンティティタイプでの検索
		if strings.Contains(strings.ToLower(entity.EntityType), queryLower) {
			match = true
		}

		// 観察事項での検索
		for _, obs := range entity.Observations {
			if strings.Contains(strings.ToLower(obs), queryLower) {
				match = true
				break
			}
		}

		if match {
			filteredEntities = append(filteredEntities, entity)
		}
	}

	// フィルタリングされたエンティティ名のセットを作成
	entityNames := make(map[string]bool)
	for _, entity := range filteredEntities {
		entityNames[entity.Name] = true
	}

	// フィルタリングされたエンティティ間のリレーションのみを含める
	var filteredRelations []Relation
	for _, relation := range graph.Relations {
		if entityNames[relation.From] && entityNames[relation.To] {
			filteredRelations = append(filteredRelations, relation)
		}
	}

	return &KnowledgeGraph{
		Entities:  filteredEntities,
		Relations: filteredRelations,
	}, nil
}

// OpenNodes は指定された名前のノードを取得する
func (m *KnowledgeGraphManager) OpenNodes(names []string) (*KnowledgeGraph, error) {
	graph, err := m.loadGraph()
	if err != nil {
		return nil, err
	}

	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}

	// 指定された名前のエンティティをフィルタリング
	var filteredEntities []Entity
	for _, entity := range graph.Entities {
		if nameSet[entity.Name] {
			filteredEntities = append(filteredEntities, entity)
		}
	}

	// フィルタリングされたエンティティ間のリレーションのみを含める
	var filteredRelations []Relation
	for _, relation := range graph.Relations {
		if nameSet[relation.From] && nameSet[relation.To] {
			filteredRelations = append(filteredRelations, relation)
		}
	}

	return &KnowledgeGraph{
		Entities:  filteredEntities,
		Relations: filteredRelations,
	}, nil
}
