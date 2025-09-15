package domain

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

// MockHTTPClient はテスト用のHTTPクライアントモック
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do はHTTPリクエストを実行する（モック）
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{}, nil
}

// #==============================================================#
// ##       normal Tests                                         ##
// #==============================================================#
// TestNewAniListClient_Normal はNewAniListClientメソッドの正常系テスト
func TestNewAniListClient_Normal(t *testing.T) {
	// Act
	result := NewAniListClient()

	// Assert
	if result == nil {
		t.Error("結果がnilです")
		return
	}
	if result.httpClient == nil {
		t.Error("httpClientがnilです")
	}
	if result.endpoint != AniListAPIEndpoint {
		t.Errorf("期待されるエンドポイント: %s, 実際: %s", AniListAPIEndpoint, result.endpoint)
	}
}

// TestNewAniListClientWithHTTPClient_Normal はNewAniListClientWithHTTPClientメソッドの正常系テスト
func TestNewAniListClientWithHTTPClient_Normal(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{}

	// Act
	result := NewAniListClientWithHTTPClient(mockClient)

	// Assert
	if result == nil {
		t.Error("結果がnilです")
		return
	}
	if result.httpClient != mockClient {
		t.Error("httpClientが設定されたモックと異なります")
	}
	if result.endpoint != AniListAPIEndpoint {
		t.Errorf("期待されるエンドポイント: %s, 実際: %s", AniListAPIEndpoint, result.endpoint)
	}
}

// TestAniListClient_QueryAnimeList_Normal はAniListClientのQueryAnimeListメソッドの正常系テスト
func TestAniListClient_QueryAnimeList_Normal(t *testing.T) {
	tests := []struct {
		name     string
		request  QueryAnimeRequest
		response string
	}{
		{
			name: "WithUsername_Normal",
			request: QueryAnimeRequest{
				Username: "testuser",
				UserID:   nil,
			},
			response: `{
				"data": {
					"MediaListCollection": {
						"lists": [
							{
								"entries": [
									{
										"media": {
											"id": 1,
											"title": {
												"native": "テストアニメ"
											},
											"coverImage": {
												"extraLarge": "https://example.com/image.jpg"
											},
											"siteUrl": "https://anilist.co/anime/1",
											"studios": {
												"nodes": [
													{
														"name": "テストスタジオ"
													}
												]
											}
										},
										"score": 85,
										"status": "COMPLETED",
										"progress": 12,
										"completedAt": {
											"year": 2023,
											"month": 6,
											"day": 15
										},
										"notes": "面白かった",
										"updatedAt": 1687123200
									}
								]
							}
						]
					}
				}
			}`,
		},
		{
			name: "WithUserID_Normal",
			request: QueryAnimeRequest{
				Username: "",
				UserID:   func() *int { id := 12345; return &id }(),
			},
			response: `{
				"data": {
					"MediaListCollection": {
						"lists": []
					}
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// リクエストの検証
					if req.Method != "POST" {
						t.Errorf("期待されるメソッド: POST, 実際: %s", req.Method)
					}
					if req.Header.Get("Content-Type") != "application/json" {
						t.Errorf("期待されるContent-Type: application/json, 実際: %s", req.Header.Get("Content-Type"))
					}
					if req.Header.Get("Accept") != "application/json" {
						t.Errorf("期待されるAccept: application/json, 実際: %s", req.Header.Get("Accept"))
					}

					// レスポンスを返す
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(tt.response)),
					}, nil
				},
			}

			client := NewAniListClientWithHTTPClient(mockClient)

			// Act
			result, err := client.QueryAnimeList(tt.request)

			// Assert
			if err != nil {
				t.Errorf("エラーが発生しました: %v", err)
				return
			}
			if result == nil {
				t.Error("結果がnilです")
				return
			}
			if result.Data == nil {
				t.Error("データがnilです")
				return
			}
			if result.Data.MediaListCollection == nil {
				t.Error("MediaListCollectionがnilです")
			}
		})
	}
}

// TestAniListClient_QueryAnimeList_HTTPError はAniListClientのQueryAnimeListメソッドのHTTPエラーテスト
func TestAniListClient_QueryAnimeList_HTTPError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
	}{
		{
			name:       "NotFound_Error",
			statusCode: 404,
			response:   "Not Found",
		},
		{
			name:       "InternalServerError_Error",
			statusCode: 500,
			response:   "Internal Server Error",
		},
		{
			name:       "BadRequest_Error",
			statusCode: 400,
			response:   "Bad Request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(strings.NewReader(tt.response)),
					}, nil
				},
			}

			client := NewAniListClientWithHTTPClient(mockClient)
			request := QueryAnimeRequest{Username: "testuser"}

			// Act
			result, err := client.QueryAnimeList(request)

			// Assert
			if err == nil {
				t.Error("エラーが期待されましたが、nilが返されました")
			}
			if result != nil {
				t.Error("エラー時は結果がnilである必要があります")
			}
			if !strings.Contains(err.Error(), "HTTPエラー") {
				t.Errorf("期待されるエラーメッセージにHTTPエラーが含まれていません: %s", err.Error())
			}
		})
	}
}

// TestAniListClient_QueryAnimeList_GraphQLError はAniListClientのQueryAnimeListメソッドのGraphQLエラーテスト
func TestAniListClient_QueryAnimeList_GraphQLError(t *testing.T) {
	// Arrange
	errorResponse := `{
		"errors": [
			{
				"message": "User not found",
				"locations": [{"line": 2, "column": 3}],
				"path": ["MediaListCollection"]
			}
		]
	}`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(errorResponse)),
			}, nil
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	request := QueryAnimeRequest{Username: "nonexistentuser"}

	// Act
	result, err := client.QueryAnimeList(request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	if !strings.Contains(err.Error(), "GraphQLエラー") {
		t.Errorf("期待されるエラーメッセージにGraphQLエラーが含まれていません: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "User not found") {
		t.Errorf("期待されるエラーメッセージにUser not foundが含まれていません: %s", err.Error())
	}
}

// networkError はネットワークエラーをシミュレートするためのカスタムエラー
type networkError struct {
	message string
}

func (e *networkError) Error() string {
	return e.message
}

// TestAniListClient_QueryAnimeList_NetworkError はAniListClientのQueryAnimeListメソッドのネットワークエラーテスト
func TestAniListClient_QueryAnimeList_NetworkError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, &networkError{message: "network error: no such host"}
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	request := QueryAnimeRequest{Username: "testuser"}

	// Act
	result, err := client.QueryAnimeList(request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	if !strings.Contains(err.Error(), "HTTPリクエストの実行に失敗しました") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %s", err.Error())
	}
}

// TestAniListClient_QueryAnimeList_InvalidJSON はAniListClientのQueryAnimeListメソッドの無効JSONテスト
func TestAniListClient_QueryAnimeList_InvalidJSON(t *testing.T) {
	// Arrange
	invalidJSON := `{"data": invalid json}`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(invalidJSON)),
			}, nil
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	request := QueryAnimeRequest{Username: "testuser"}

	// Act
	result, err := client.QueryAnimeList(request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	if !strings.Contains(err.Error(), "レスポンスのJSONデコードに失敗しました") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %s", err.Error())
	}
}

// TestAniListClient_QueryAnimeList_ReadBodyError はAniListClientのQueryAnimeListメソッドのボディ読み取りエラーテスト
func TestAniListClient_QueryAnimeList_ReadBodyError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       &errorReader{},
			}, nil
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	request := QueryAnimeRequest{Username: "testuser"}

	// Act
	result, err := client.QueryAnimeList(request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	if !strings.Contains(err.Error(), "レスポンスボディの読み取りに失敗しました") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %s", err.Error())
	}
}

// errorReader は読み取り時にエラーを返すReader
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (e *errorReader) Close() error {
	return nil
}

// TestAniListClient_QueryAnimeList_BothParameters はAniListClientのQueryAnimeListメソッドの両パラメータテスト
func TestAniListClient_QueryAnimeList_BothParameters(t *testing.T) {
	// Arrange
	response := `{
		"data": {
			"MediaListCollection": {
				"lists": []
			}
		}
	}`

	var capturedRequestBody []byte
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// リクエストボディをキャプチャ
			body, _ := io.ReadAll(req.Body)
			capturedRequestBody = body
			req.Body = io.NopCloser(bytes.NewReader(body))

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(response)),
			}, nil
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	userID := 12345
	request := QueryAnimeRequest{
		Username: "testuser",
		UserID:   &userID,
	}

	// Act
	result, err := client.QueryAnimeList(request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
		return
	}
	if result == nil {
		t.Error("結果がnilです")
		return
	}

	// リクエストボディに両方のパラメータが含まれていることを確認
	requestBodyStr := string(capturedRequestBody)
	if !strings.Contains(requestBodyStr, "testuser") {
		t.Error("リクエストボディにユーザー名が含まれていません")
	}
	if !strings.Contains(requestBodyStr, "12345") {
		t.Error("リクエストボディにユーザーIDが含まれていません")
	}
}

// TestAniListClient_QueryMangaList_Normal はAniListClientのQueryMangaListメソッドの正常系テスト
func TestAniListClient_QueryMangaList_Normal(t *testing.T) {
	tests := []struct {
		name     string
		request  QueryMangaRequest
		response string
	}{
		{
			name: "WithUsername_Normal",
			request: QueryMangaRequest{
				Username: "testuser",
				UserID:   nil,
			},
			response: `{
				"data": {
					"MediaListCollection": {
						"lists": [
							{
								"entries": [
									{
										"media": {
											"id": 1,
											"title": {
												"native": "テストマンガ"
											},
											"coverImage": {
												"extraLarge": "https://example.com/manga.jpg"
											},
											"siteUrl": "https://anilist.co/manga/1"
										},
										"score": 85,
										"status": "COMPLETED",
										"progress": 120,
										"progressVolumes": 12,
										"completedAt": {
											"year": 2023,
											"month": 6,
											"day": 15
										},
										"notes": "面白かった",
										"updatedAt": 1687123200
									}
								]
							}
						]
					}
				}
			}`,
		},
		{
			name: "WithUserID_Normal",
			request: QueryMangaRequest{
				Username: "",
				UserID:   func() *int { id := 12345; return &id }(),
			},
			response: `{
				"data": {
					"MediaListCollection": {
						"lists": []
					}
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// リクエストの検証
					if req.Method != "POST" {
						t.Errorf("期待されるメソッド: POST, 実際: %s", req.Method)
					}
					if req.Header.Get("Content-Type") != "application/json" {
						t.Errorf("期待されるContent-Type: application/json, 実際: %s", req.Header.Get("Content-Type"))
					}
					if req.Header.Get("Accept") != "application/json" {
						t.Errorf("期待されるAccept: application/json, 実際: %s", req.Header.Get("Accept"))
					}

					// レスポンスを返す
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(tt.response)),
					}, nil
				},
			}

			client := NewAniListClientWithHTTPClient(mockClient)

			// Act
			result, err := client.QueryMangaList(tt.request)

			// Assert
			if err != nil {
				t.Errorf("エラーが発生しました: %v", err)
				return
			}
			if result == nil {
				t.Error("結果がnilです")
				return
			}
			if result.Data == nil {
				t.Error("データがnilです")
				return
			}
			if result.Data.MediaListCollection == nil {
				t.Error("MediaListCollectionがnilです")
			}
		})
	}
}

// TestAniListClient_QueryMangaList_HTTPError はAniListClientのQueryMangaListメソッドのHTTPエラーテスト
func TestAniListClient_QueryMangaList_HTTPError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
	}{
		{
			name:       "NotFound_Error",
			statusCode: 404,
			response:   "Not Found",
		},
		{
			name:       "InternalServerError_Error",
			statusCode: 500,
			response:   "Internal Server Error",
		},
		{
			name:       "BadRequest_Error",
			statusCode: 400,
			response:   "Bad Request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(strings.NewReader(tt.response)),
					}, nil
				},
			}

			client := NewAniListClientWithHTTPClient(mockClient)
			request := QueryMangaRequest{Username: "testuser"}

			// Act
			result, err := client.QueryMangaList(request)

			// Assert
			if err == nil {
				t.Error("エラーが期待されましたが、nilが返されました")
			}
			if result != nil {
				t.Error("エラー時は結果がnilである必要があります")
			}
			if !strings.Contains(err.Error(), "HTTPエラー") {
				t.Errorf("期待されるエラーメッセージにHTTPエラーが含まれていません: %s", err.Error())
			}
		})
	}
}

// TestAniListClient_QueryMangaList_GraphQLError はAniListClientのQueryMangaListメソッドのGraphQLエラーテスト
func TestAniListClient_QueryMangaList_GraphQLError(t *testing.T) {
	// Arrange
	errorResponse := `{
		"errors": [
			{
				"message": "User not found",
				"locations": [{"line": 2, "column": 3}],
				"path": ["MediaListCollection"]
			}
		]
	}`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(errorResponse)),
			}, nil
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	request := QueryMangaRequest{Username: "nonexistentuser"}

	// Act
	result, err := client.QueryMangaList(request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	if !strings.Contains(err.Error(), "GraphQLエラー") {
		t.Errorf("期待されるエラーメッセージにGraphQLエラーが含まれていません: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "User not found") {
		t.Errorf("期待されるエラーメッセージにUser not foundが含まれていません: %s", err.Error())
	}
}

// TestMockAniListRepository_QueryAnimeList はMockAniListRepositoryのQueryAnimeListメソッドテスト
func TestMockAniListRepository_QueryAnimeList(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func(req QueryAnimeRequest) (*AniListResponse, error)
		request        QueryAnimeRequest
		expectedResult *AniListResponse
		expectedError  error
	}{
		{
			name: "WithMockFunc_Normal",
			mockFunc: func(req QueryAnimeRequest) (*AniListResponse, error) {
				return &AniListResponse{
					Data: &MediaListCollectionData{
						MediaListCollection: &MediaListCollection{
							Lists: []MediaList{},
						},
					},
				}, nil
			},
			request: QueryAnimeRequest{Username: "testuser"},
			expectedResult: &AniListResponse{
				Data: &MediaListCollectionData{
					MediaListCollection: &MediaListCollection{
						Lists: []MediaList{},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:           "WithoutMockFunc_Normal",
			mockFunc:       nil,
			request:        QueryAnimeRequest{Username: "testuser"},
			expectedResult: &AniListResponse{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mock := &MockAniListRepository{
				QueryAnimeListFunc: tt.mockFunc,
			}

			// Act
			result, err := mock.QueryAnimeList(tt.request)

			// Assert
			if (err == nil) != (tt.expectedError == nil) {
				t.Errorf("期待されるエラー: %v, 実際: %v", tt.expectedError, err)
			}
			if result == nil && tt.expectedResult != nil {
				t.Error("結果がnilですが、期待される結果があります")
			}
			if result != nil && tt.expectedResult == nil {
				t.Error("結果がnilではありませんが、期待される結果がnilです")
			}
		})
	}
}

// TestMockAniListRepository_QueryMangaList はMockAniListRepositoryのQueryMangaListメソッドテスト
func TestMockAniListRepository_QueryMangaList(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func(req QueryMangaRequest) (*AniListResponse, error)
		request        QueryMangaRequest
		expectedResult *AniListResponse
		expectedError  error
	}{
		{
			name: "WithMockFunc_Normal",
			mockFunc: func(req QueryMangaRequest) (*AniListResponse, error) {
				return &AniListResponse{
					Data: &MediaListCollectionData{
						MediaListCollection: &MediaListCollection{
							Lists: []MediaList{},
						},
					},
				}, nil
			},
			request: QueryMangaRequest{Username: "testuser"},
			expectedResult: &AniListResponse{
				Data: &MediaListCollectionData{
					MediaListCollection: &MediaListCollection{
						Lists: []MediaList{},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:           "WithoutMockFunc_Normal",
			mockFunc:       nil,
			request:        QueryMangaRequest{Username: "testuser"},
			expectedResult: &AniListResponse{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mock := &MockAniListRepository{
				QueryMangaListFunc: tt.mockFunc,
			}

			// Act
			result, err := mock.QueryMangaList(tt.request)

			// Assert
			if (err == nil) != (tt.expectedError == nil) {
				t.Errorf("期待されるエラー: %v, 実際: %v", tt.expectedError, err)
			}
			if result == nil && tt.expectedResult != nil {
				t.Error("結果がnilですが、期待される結果があります")
			}
			if result != nil && tt.expectedResult == nil {
				t.Error("結果がnilではありませんが、期待される結果がnilです")
			}
		})
	}
}

// TestAniListClient_QueryMangaList_NetworkError はAniListClientのQueryMangaListメソッドのネットワークエラーテスト
func TestAniListClient_QueryMangaList_NetworkError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, &networkError{message: "network error: no such host"}
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	request := QueryMangaRequest{Username: "testuser"}

	// Act
	result, err := client.QueryMangaList(request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	if !strings.Contains(err.Error(), "HTTPリクエストの実行に失敗しました") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %s", err.Error())
	}
}

// TestAniListClient_QueryMangaList_InvalidJSON はAniListClientのQueryMangaListメソッドの無効JSONテスト
func TestAniListClient_QueryMangaList_InvalidJSON(t *testing.T) {
	// Arrange
	invalidJSON := `{"data": invalid json}`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(invalidJSON)),
			}, nil
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	request := QueryMangaRequest{Username: "testuser"}

	// Act
	result, err := client.QueryMangaList(request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	if !strings.Contains(err.Error(), "レスポンスのJSONデコードに失敗しました") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %s", err.Error())
	}
}

// TestAniListClient_QueryMangaList_ReadBodyError はAniListClientのQueryMangaListメソッドのボディ読み取りエラーテスト
func TestAniListClient_QueryMangaList_ReadBodyError(t *testing.T) {
	// Arrange
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       &errorReader{},
			}, nil
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	request := QueryMangaRequest{Username: "testuser"}

	// Act
	result, err := client.QueryMangaList(request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	if !strings.Contains(err.Error(), "レスポンスボディの読み取りに失敗しました") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %s", err.Error())
	}
}

// TestAniListClient_QueryMangaList_BothParameters はAniListClientのQueryMangaListメソッドの両パラメータテスト
func TestAniListClient_QueryMangaList_BothParameters(t *testing.T) {
	// Arrange
	response := `{
		"data": {
			"MediaListCollection": {
				"lists": []
			}
		}
	}`

	var capturedRequestBody []byte
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// リクエストボディをキャプチャ
			body, _ := io.ReadAll(req.Body)
			capturedRequestBody = body
			req.Body = io.NopCloser(bytes.NewReader(body))

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(response)),
			}, nil
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	userID := 12345
	request := QueryMangaRequest{
		Username: "testuser",
		UserID:   &userID,
	}

	// Act
	result, err := client.QueryMangaList(request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
		return
	}
	if result == nil {
		t.Error("結果がnilです")
		return
	}

	// リクエストボディに両方のパラメータが含まれていることを確認
	requestBodyStr := string(capturedRequestBody)
	if !strings.Contains(requestBodyStr, "testuser") {
		t.Error("リクエストボディにユーザー名が含まれていません")
	}
	if !strings.Contains(requestBodyStr, "12345") {
		t.Error("リクエストボディにユーザーIDが含まれていません")
	}
}

// TestMockHTTPClient_Do はMockHTTPClientのDoメソッドテスト
func TestMockHTTPClient_Do(t *testing.T) {
	tests := []struct {
		name         string
		doFunc       func(req *http.Request) (*http.Response, error)
		expectedResp *http.Response
		expectedErr  error
	}{
		{
			name: "WithDoFunc_Normal",
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: 200}, nil
			},
			expectedResp: &http.Response{StatusCode: 200},
			expectedErr:  nil,
		},
		{
			name:         "WithoutDoFunc_Normal",
			doFunc:       nil,
			expectedResp: &http.Response{},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mock := &MockHTTPClient{
				DoFunc: tt.doFunc,
			}
			req := &http.Request{}

			// Act
			resp, err := mock.Do(req)

			// Assert
			if (err == nil) != (tt.expectedErr == nil) {
				t.Errorf("期待されるエラー: %v, 実際: %v", tt.expectedErr, err)
			}
			if resp == nil && tt.expectedResp != nil {
				t.Error("レスポンスがnilですが、期待されるレスポンスがあります")
			}
			if resp != nil && tt.expectedResp == nil {
				t.Error("レスポンスがnilではありませんが、期待されるレスポンスがnilです")
			}
			if resp != nil && tt.expectedResp != nil && resp.StatusCode != tt.expectedResp.StatusCode {
				t.Errorf("期待されるステータスコード: %d, 実際: %d", tt.expectedResp.StatusCode, resp.StatusCode)
			}
		})
	}
}

// #==============================================================#
// ##       Abnormal Tests                                       ##
// #==============================================================#
// TestAniListClient_QueryAnimeList_HTTPRequestCreationError はAniListClientのQueryAnimeListメソッドのHTTPリクエスト作成エラーテスト
func TestAniListClient_QueryAnimeList_HTTPRequestCreationError(t *testing.T) {
	// Arrange
	// 無効なエンドポイントを設定してHTTPリクエスト作成エラーを発生させる
	client := &AniListClient{
		httpClient: &MockHTTPClient{},
		endpoint:   "://invalid-url", // 無効なURL
	}

	request := QueryAnimeRequest{Username: "testuser"}

	// Act
	result, err := client.QueryAnimeList(request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	if !strings.Contains(err.Error(), "HTTPリクエストの作成に失敗しました") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %s", err.Error())
	}
}

// TestAniListClient_QueryMangaList_HTTPRequestCreationError はAniListClientのQueryMangaListメソッドのHTTPリクエスト作成エラーテスト
func TestAniListClient_QueryMangaList_HTTPRequestCreationError(t *testing.T) {
	// Arrange
	// 無効なエンドポイントを設定してHTTPリクエスト作成エラーを発生させる
	client := &AniListClient{
		httpClient: &MockHTTPClient{},
		endpoint:   "://invalid-url", // 無効なURL
	}

	request := QueryMangaRequest{Username: "testuser"}

	// Act
	result, err := client.QueryMangaList(request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
	if !strings.Contains(err.Error(), "HTTPリクエストの作成に失敗しました") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %s", err.Error())
	}
}

// TestAniListClient_QueryAnimeList_EmptyParameters はAniListClientのQueryAnimeListメソッドの空パラメータテスト
func TestAniListClient_QueryAnimeList_EmptyParameters(t *testing.T) {
	// Arrange
	response := `{
		"data": {
			"MediaListCollection": {
				"lists": []
			}
		}
	}`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(response)),
			}, nil
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	request := QueryAnimeRequest{} // 空のパラメータ

	// Act
	result, err := client.QueryAnimeList(request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
		return
	}
	if result == nil {
		t.Error("結果がnilです")
		return
	}
	if result.Data == nil {
		t.Error("データがnilです")
		return
	}
	if result.Data.MediaListCollection == nil {
		t.Error("MediaListCollectionがnilです")
	}
}

// TestAniListClient_QueryMangaList_EmptyParameters はAniListClientのQueryMangaListメソッドの空パラメータテスト
func TestAniListClient_QueryMangaList_EmptyParameters(t *testing.T) {
	// Arrange
	response := `{
		"data": {
			"MediaListCollection": {
				"lists": []
			}
		}
	}`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(response)),
			}, nil
		},
	}

	client := NewAniListClientWithHTTPClient(mockClient)
	request := QueryMangaRequest{} // 空のパラメータ

	// Act
	result, err := client.QueryMangaList(request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
		return
	}
	if result == nil {
		t.Error("結果がnilです")
		return
	}
	if result.Data == nil {
		t.Error("データがnilです")
		return
	}
	if result.Data.MediaListCollection == nil {
		t.Error("MediaListCollectionがnilです")
	}
}
