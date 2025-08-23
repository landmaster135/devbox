package domain

// MockAniListRepository はテスト用のAniListRepositoryモック
type MockAniListRepository struct {
	QueryAnimeListFunc func(req QueryAnimeRequest) (*AniListResponse, error)
}

// QueryAnimeList はアニメリストを取得する（モック）
func (m *MockAniListRepository) QueryAnimeList(req QueryAnimeRequest) (*AniListResponse, error) {
	if m.QueryAnimeListFunc != nil {
		return m.QueryAnimeListFunc(req)
	}
	return &AniListResponse{}, nil
}
