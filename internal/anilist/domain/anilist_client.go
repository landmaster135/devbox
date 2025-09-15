package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// AniListAPIEndpoint はAniList GraphQL APIのエンドポイント
	AniListAPIEndpoint = "https://graphql.anilist.co"

	// GraphQLクエリ
	MediaListCollectionQuery = `
		query ($username: String, $id: Int) {
			MediaListCollection(userName: $username, userId: $id, type: ANIME) {
				lists {
					entries {
						media {
							id
							title {
								native
							}
							coverImage {
								extraLarge
							}
							siteUrl
							studios (isMain: true, sort: FAVOURITES) {
								nodes {
									name
								}
							}
						}
						score (format: POINT_100)
						status
						progress
						completedAt {
							year
							month
							day
						}
						notes
						updatedAt
					}
				}
			}
		}
	`

	// MangaListCollectionQuery はマンガリスト取得用のGraphQLクエリ
	MangaListCollectionQuery = `
		query ($username: String, $id: Int) {
			MediaListCollection(userName: $username, userId: $id, type: MANGA) {
				lists {
					entries {
						media {
							id
							title {
								native
							}
							coverImage {
								extraLarge
							}
							siteUrl
						}
						score (format: POINT_100)
						status
						progress
						completedAt {
							year
							month
							day
						}
						notes
						updatedAt
					}
				}
			}
		}
	`
)

// #==============================================================#
// ##          Repositories                                      ##
// #==============================================================#
// AniListRepository はAniListデータ取得のインターフェース
type AniListRepository interface {
	QueryAnimeList(req QueryAnimeRequest) (*AniListResponse, error)
	QueryMangaList(req QueryMangaRequest) (*AniListResponse, error)
}

// HTTPClient はHTTPクライアントのインターフェース
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// #==============================================================#
// ##          Mocks                                             ##
// #==============================================================#
// MockAniListRepository はテスト用のAniListRepositoryモック
type MockAniListRepository struct {
	QueryAnimeListFunc func(req QueryAnimeRequest) (*AniListResponse, error)
	QueryMangaListFunc func(req QueryMangaRequest) (*AniListResponse, error)
}

// QueryAnimeList はアニメリストを取得する（モック）
func (m *MockAniListRepository) QueryAnimeList(req QueryAnimeRequest) (*AniListResponse, error) {
	if m.QueryAnimeListFunc != nil {
		return m.QueryAnimeListFunc(req)
	}
	return &AniListResponse{}, nil
}

// QueryMangaList はマンガリストを取得する（モック）
func (m *MockAniListRepository) QueryMangaList(req QueryMangaRequest) (*AniListResponse, error) {
	if m.QueryMangaListFunc != nil {
		return m.QueryMangaListFunc(req)
	}
	return &AniListResponse{}, nil
}

// #==============================================================#
// ##          Implements                                        ##
// #==============================================================#
// AniListClient はAniList APIクライアント
type AniListClient struct {
	httpClient HTTPClient
	endpoint   string
}

// NewAniListClient は新しいAniListClientを作成する
func NewAniListClient() *AniListClient {
	return &AniListClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		endpoint: AniListAPIEndpoint,
	}
}

// NewAniListClientWithHTTPClient はHTTPクライアントを指定してAniListClientを作成する
func NewAniListClientWithHTTPClient(httpClient HTTPClient) *AniListClient {
	return &AniListClient{
		httpClient: httpClient,
		endpoint:   AniListAPIEndpoint,
	}
}

// QueryAnimeList はアニメリストを取得する
func (c *AniListClient) QueryAnimeList(req QueryAnimeRequest) (*AniListResponse, error) {
	// GraphQLリクエストを構築
	variables := make(map[string]any)
	if req.Username != "" {
		variables["username"] = req.Username
	}
	if req.UserID != nil {
		variables["id"] = *req.UserID
	}

	graphQLReq := GraphQLRequest{
		Query:     MediaListCollectionQuery,
		Variables: variables,
	}

	// JSONにエンコード
	reqBody, err := json.Marshal(graphQLReq)
	if err != nil {
		return nil, fmt.Errorf("リクエストのJSONエンコードに失敗しました: %v", err)
	}

	// HTTPリクエストを作成
	httpReq, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの作成に失敗しました: %v", err)
	}

	// ヘッダーを設定
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// リクエストを実行
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの実行に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンスボディの読み取りに失敗しました: %v", err)
	}

	// HTTPステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTPエラー: %d - %s", resp.StatusCode, string(body))
	}

	// JSONをデコード
	var aniListResp AniListResponse
	if err := json.Unmarshal(body, &aniListResp); err != nil {
		return nil, fmt.Errorf("レスポンスのJSONデコードに失敗しました: %v", err)
	}

	// GraphQLエラーをチェック
	if len(aniListResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQLエラー: %s", aniListResp.Errors[0].Message)
	}

	return &aniListResp, nil
}

// QueryMangaList はマンガリストを取得する
func (c *AniListClient) QueryMangaList(req QueryMangaRequest) (*AniListResponse, error) {
	// GraphQLリクエストを構築
	variables := make(map[string]any)
	if req.Username != "" {
		variables["username"] = req.Username
	}
	if req.UserID != nil {
		variables["id"] = *req.UserID
	}

	graphQLReq := GraphQLRequest{
		Query:     MangaListCollectionQuery,
		Variables: variables,
	}

	// JSONにエンコード
	reqBody, err := json.Marshal(graphQLReq)
	if err != nil {
		return nil, fmt.Errorf("リクエストのJSONエンコードに失敗しました: %v", err)
	}

	// HTTPリクエストを作成
	httpReq, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの作成に失敗しました: %v", err)
	}

	// ヘッダーを設定
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// リクエストを実行
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの実行に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンスボディの読み取りに失敗しました: %v", err)
	}

	// HTTPステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTPエラー: %d - %s", resp.StatusCode, string(body))
	}

	// JSONをデコード
	var aniListResp AniListResponse
	if err := json.Unmarshal(body, &aniListResp); err != nil {
		return nil, fmt.Errorf("レスポンスのJSONデコードに失敗しました: %v", err)
	}

	// GraphQLエラーをチェック
	if len(aniListResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQLエラー: %s", aniListResp.Errors[0].Message)
	}

	return &aniListResp, nil
}
