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
