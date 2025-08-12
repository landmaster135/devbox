package repositories

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	models "github.com/landmaster135/devbox/internal/http_request/domain/models"
)

func TestHTTPRepositoryImpl_SendRequest_InvalidURL(t *testing.T) {
	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// 無効なURLを含むリクエストを作成
	request := &models.HTTPRequest{
		URL:     "://invalid-url", // 無効なURL
		Method:  "GET",
		Headers: map[string]string{"Accept": "application/json"},
		Body:    nil,
	}

	// リクエストを送信
	response, err := repo.SendRequest(request)

	// 検証
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "HTTPリクエストの作成に失敗しました")
}

func TestHTTPRepositoryImpl_SendRequest_ConnectionError(t *testing.T) {
	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// 存在しないサーバーへのリクエストを作成
	request := &models.HTTPRequest{
		URL:     "http://localhost:12345", // 存在しないポート
		Method:  "GET",
		Headers: map[string]string{"Accept": "application/json"},
		Body:    nil,
	}

	// リクエストを送信
	response, err := repo.SendRequest(request)

	// 検証
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "HTTPリクエストの送信に失敗しました")
}

func TestHTTPRepositoryImpl_SendRequest_Normal(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストメソッドの検証
		assert.Equal(t, "GET", r.Method)

		// リクエストヘッダーの検証
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		// レスポンスを設定
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// リクエストを作成
	request := &models.HTTPRequest{
		URL:     server.URL,
		Method:  "GET",
		Headers: map[string]string{"Accept": "application/json"},
		Body:    nil,
	}

	// リクエストを送信
	response, err := repo.SendRequest(request)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Headers["Content-Type"])

	// レスポンスボディを検証
	var responseBody map[string]interface{}
	err = json.Unmarshal(response.Body, &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "success", responseBody["message"])
}

func TestHTTPRepositoryImpl_LoadJSONFile_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "http-repo-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// テスト用のJSONファイルを作成
	jsonContent := `{"key": "value", "number": 42}`
	jsonFilePath := filepath.Join(tempDir, "test.json")
	err = os.WriteFile(jsonFilePath, []byte(jsonContent), 0644)
	assert.NoError(t, err)

	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// JSONファイルを読み込む
	content, err := repo.LoadJSONFile(jsonFilePath)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, content)

	// JSONとして解析できることを確認
	var jsonData map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	assert.NoError(t, err)
	assert.Equal(t, "value", jsonData["key"])
	assert.Equal(t, float64(42), jsonData["number"])
}

func TestHTTPRepositoryImpl_LoadJSONFile_FileNotFound(t *testing.T) {
	// 存在しないファイルパス
	nonExistentFilePath := "/path/to/non/existent/file.json"

	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// JSONファイルを読み込む
	content, err := repo.LoadJSONFile(nonExistentFilePath)

	// 検証
	assert.Error(t, err)
	assert.Nil(t, content)
	assert.Contains(t, err.Error(), "JSONファイルを開けませんでした")
}

func TestHTTPRepositoryImpl_LoadJSONFile_InvalidJSON(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "http-repo-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 無効なJSONファイルを作成
	invalidJSONContent := `{"key": "value", "number": 42,}` // 末尾のカンマが無効
	invalidJSONFilePath := filepath.Join(tempDir, "invalid.json")
	err = os.WriteFile(invalidJSONFilePath, []byte(invalidJSONContent), 0644)
	assert.NoError(t, err)

	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// JSONファイルを読み込む
	content, err := repo.LoadJSONFile(invalidJSONFilePath)

	// 検証
	assert.Error(t, err)
	assert.Nil(t, content)
	assert.Contains(t, err.Error(), "無効なJSON形式です")
}

func TestHTTPRepositoryImpl_LoadJSONFile_ReadError(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "http-repo-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// ディレクトリを作成（ファイルではなくディレクトリを読み込もうとする）
	dirPath := filepath.Join(tempDir, "test-dir")
	err = os.Mkdir(dirPath, 0755)
	assert.NoError(t, err)

	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// ディレクトリを読み込もうとする（失敗するはず）
	content, err := repo.LoadJSONFile(dirPath)

	// 検証
	assert.Error(t, err)
	assert.Nil(t, content)
	assert.Contains(t, err.Error(), "JSONファイルの読み込みに失敗しました")
}

func TestHTTPRepositoryImpl_LoadJSONFile_WithBOM(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "http-repo-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// BOMを含むJSONファイルを作成
	jsonContent := []byte{0xEF, 0xBB, 0xBF, '{', '"', 'k', 'e', 'y', '"', ':', ' ', '"', 'v', 'a', 'l', 'u', 'e', '"', '}'}
	jsonFilePath := filepath.Join(tempDir, "bom.json")
	err = os.WriteFile(jsonFilePath, jsonContent, 0644)
	assert.NoError(t, err)

	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// JSONファイルを読み込む
	content, err := repo.LoadJSONFile(jsonFilePath)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, content)

	// JSONとして解析できることを確認
	var jsonData map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	assert.NoError(t, err)
	assert.Equal(t, "value", jsonData["key"])
}

func TestHTTPRepositoryImpl_LoadJSONFile_WithCRLF(t *testing.T) {
	// プロジェクトルートからの相対パスでテストデータのパスを取得
	testDataPath := "./test_data/org/sample_request_with_crlf.json"

	// ファイルが存在することを確認
	_, err := os.Stat(testDataPath)
	assert.NoError(t, err, "テストデータファイルが見つかりません: %s", testDataPath)

	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// JSONファイルを読み込む
	content, err := repo.LoadJSONFile(testDataPath)

	// エラーがないことを確認
	assert.NoError(t, err, "JSONファイルの読み込みに失敗しました")
	assert.NotNil(t, content, "読み込まれたデータがnilです")
	assert.NotEmpty(t, content, "読み込まれたデータが空です")

	// JSONとして解析できることを確認
	var jsonObj map[string]interface{}
	err = json.Unmarshal(content, &jsonObj)
	assert.NoError(t, err, "JSONとして解析できません")

	// JSONオブジェクトの内容を確認
	assert.Equal(t, "テストユーザー", jsonObj["name"], "nameフィールドの値が一致しません")
	assert.Equal(t, "test@example.com", jsonObj["email"], "emailフィールドの値が一致しません")
	assert.Equal(t, float64(30), jsonObj["age"], "ageフィールドの値が一致しません")

	// interestsフィールドの確認
	interests, ok := jsonObj["interests"].([]interface{})
	assert.True(t, ok, "interestsフィールドが配列ではありません")
	assert.Len(t, interests, 3, "interestsフィールドの要素数が一致しません")
	assert.Equal(t, "プログラミング", interests[0], "interestsフィールドの最初の要素が一致しません")

	// addressフィールドの確認
	address, ok := jsonObj["address"].(map[string]interface{})
	assert.True(t, ok, "addressフィールドがオブジェクトではありません")
	assert.Equal(t, "日本", address["country"], "address.countryフィールドの値が一致しません")
	assert.Equal(t, "東京", address["city"], "address.cityフィールドの値が一致しません")
}

func TestNewHTTPRepository_Normal(t *testing.T) {
	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// 検証
	assert.NotNil(t, repo, "リポジトリがnilです")
	assert.NotNil(t, repo.client, "HTTPクライアントがnilです")
}

func TestHTTPRepositoryImpl_SendRequest_MultipleHeaderValues(t *testing.T) {
	// テスト用のHTTPサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストヘッダーの検証
		assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		// レスポンスヘッダーに複数の値を設定
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Custom-Header", "custom-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// リクエストを作成
	request := &models.HTTPRequest{
		URL:    server.URL,
		Method: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer token123",
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		},
		Body: []byte(`{"test": "data"}`),
	}

	// リクエストを送信
	response, err := repo.SendRequest(request)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Headers["Content-Type"])
	assert.Equal(t, "no-cache", response.Headers["Cache-Control"])
	assert.Equal(t, "custom-value", response.Headers["X-Custom-Header"])

	// レスポンスボディを検証
	var responseBody map[string]interface{}
	err = json.Unmarshal(response.Body, &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "success", responseBody["result"])
}

func TestHTTPRepositoryImpl_SendRequest_DifferentStatusCodes(t *testing.T) {
	testCases := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedStatus int
	}{
		{
			name:           "201 Created",
			statusCode:     http.StatusCreated,
			responseBody:   `{"id": 123, "created": true}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "400 Bad Request",
			statusCode:     http.StatusBadRequest,
			responseBody:   `{"error": "Invalid request"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "404 Not Found",
			statusCode:     http.StatusNotFound,
			responseBody:   `{"error": "Resource not found"}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   `{"error": "Internal server error"}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用のHTTPサーバーを作成
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.responseBody))
			}))
			defer server.Close()

			// テスト対象のリポジトリを作成
			repo := NewHTTPRepository()

			// リクエストを作成
			request := &models.HTTPRequest{
				URL:     server.URL,
				Method:  "GET",
				Headers: map[string]string{"Accept": "application/json"},
				Body:    nil,
			}

			// リクエストを送信
			response, err := repo.SendRequest(request)

			// 検証
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, tc.expectedStatus, response.StatusCode)
			assert.Equal(t, "application/json", response.Headers["Content-Type"])
			assert.Equal(t, tc.responseBody, string(response.Body))
		})
	}
}

func TestHTTPRepositoryImpl_SendRequest_DifferentMethods(t *testing.T) {
	testCases := []struct {
		name   string
		method string
		body   []byte
	}{
		{
			name:   "GET request",
			method: "GET",
			body:   nil,
		},
		{
			name:   "POST request",
			method: "POST",
			body:   []byte(`{"data": "test"}`),
		},
		{
			name:   "PUT request",
			method: "PUT",
			body:   []byte(`{"update": "data"}`),
		},
		{
			name:   "DELETE request",
			method: "DELETE",
			body:   nil,
		},
		{
			name:   "PATCH request",
			method: "PATCH",
			body:   []byte(`{"patch": "data"}`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用のHTTPサーバーを作成
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// リクエストメソッドの検証
				assert.Equal(t, tc.method, r.Method)

				// ボディがある場合は検証
				if tc.body != nil {
					body, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.Equal(t, tc.body, body)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"method": "` + tc.method + `"}`))
			}))
			defer server.Close()

			// テスト対象のリポジトリを作成
			repo := NewHTTPRepository()

			// リクエストを作成
			request := &models.HTTPRequest{
				URL:     server.URL,
				Method:  tc.method,
				Headers: map[string]string{"Accept": "application/json"},
				Body:    tc.body,
			}

			// リクエストを送信
			response, err := repo.SendRequest(request)

			// 検証
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, http.StatusOK, response.StatusCode)
			assert.Contains(t, string(response.Body), tc.method)
		})
	}
}

func TestHTTPRepositoryImpl_SendRequest_EmptyResponse(t *testing.T) {
	// テスト用のHTTPサーバーを作成（空のレスポンス）
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		// ボディは書き込まない
	}))
	defer server.Close()

	// テスト対象のリポジトリを作成
	repo := NewHTTPRepository()

	// リクエストを作成
	request := &models.HTTPRequest{
		URL:     server.URL,
		Method:  "DELETE",
		Headers: map[string]string{"Accept": "application/json"},
		Body:    nil,
	}

	// リクエストを送信
	response, err := repo.SendRequest(request)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
	assert.Empty(t, response.Body)
}
