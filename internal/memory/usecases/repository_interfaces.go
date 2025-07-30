package usecases

// KnowledgeGraphRepository は知識グラフの永続化を抽象化するインターフェース
type KnowledgeGraphRepository interface {
	// LoadGraph は知識グラフを読み込む
	LoadGraph() (*KnowledgeGraph, error)

	// SaveGraph は知識グラフを保存する
	SaveGraph(graph *KnowledgeGraph) error
}
