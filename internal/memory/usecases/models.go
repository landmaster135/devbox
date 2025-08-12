package usecases

// Entity は知識グラフのエンティティを表す構造体
type Entity struct {
	Name         string   `json:"name"`
	EntityType   string   `json:"entityType"`
	Observations []string `json:"observations"`
}

// Relation は知識グラフのリレーションを表す構造体
type Relation struct {
	From         string `json:"from"`
	To           string `json:"to"`
	RelationType string `json:"relationType"`
}

// KnowledgeGraph は知識グラフ全体を表す構造体
type KnowledgeGraph struct {
	Entities  []Entity  `json:"entities"`
	Relations []Relation `json:"relations"`
}

// ObservationInput は観察事項追加用の入力構造体
type ObservationInput struct {
	EntityName string   `json:"entityName"`
	Contents   []string `json:"contents"`
}

// ObservationResult は観察事項追加の結果構造体
type ObservationResult struct {
	EntityName        string   `json:"entityName"`
	AddedObservations []string `json:"addedObservations"`
}

// DeletionInput は観察事項削除用の入力構造体
type DeletionInput struct {
	EntityName   string   `json:"entityName"`
	Observations []string `json:"observations"`
}
