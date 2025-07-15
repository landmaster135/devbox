package repositories

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	models "github.com/landmaster135/devbox/internal/api_client/domain/models"
)

// TestNewAPIRepository はNewAPIRepositoryメソッドのテストです
func TestNewAPIRepository(t *testing.T) {
	// Act
	repo := NewAPIRepository()

	// Assert
	if repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}
	if repo.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

// TestConvertEncoding_ShiftJIS はShift_JIS変換のテストです
func TestConvertEncoding_ShiftJIS(t *testing.T) {
	repo := NewAPIRepository()

	// Shift_JISエンコードされたHTMLサンプル
	htmlContent := `<html><head><title>阿部寛のホームページ</title></head><body>テスト</body></html>`
	shiftJISHTML, err := encodeToShiftJIS(htmlContent)
	if err != nil {
		t.Fatalf("Failed to encode HTML to Shift_JIS: %v", err)
	}

	testCases := []struct {
		name              string
		body              []byte
		contentType       string
		specifiedEncoding string
		expectConversion  bool
	}{
		{
			name:              "Shift_JIS指定でHTML変換",
			body:              shiftJISHTML,
			contentType:       "text/html",
			specifiedEncoding: "shift_jis",
			expectConversion:  true,
		},
		{
			name:              "Shift-JIS指定でHTML変換（ハイフン付き）",
			body:              shiftJISHTML,
			contentType:       "text/html",
			specifiedEncoding: "shift-jis",
			expectConversion:  true,
		},
		{
			name:              "UTF-8指定で変換なし",
			body:              []byte(htmlContent),
			contentType:       "text/html",
			specifiedEncoding: "utf-8",
			expectConversion:  false,
		},
		{
			name:              "JSON形式では変換しない",
			body:              shiftJISHTML,
			contentType:       "application/json",
			specifiedEncoding: "shift_jis",
			expectConversion:  false,
		},
		{
			name:              "auto指定でContent-Typeから検出",
			body:              shiftJISHTML,
			contentType:       "text/html; charset=shift_jis",
			specifiedEncoding: "auto",
			expectConversion:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := repo.convertEncoding(tc.body, tc.contentType, tc.specifiedEncoding)

			// Assert
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if tc.expectConversion {
				// 変換されたテキストに日本語が含まれていることを確認
				resultStr := string(result)
				if !strings.Contains(resultStr, "阿部寛") {
					t.Errorf("Expected converted text to contain Japanese characters, got %s", resultStr)
				}
			} else {
				// 変換されていないことを確認
				if !bytes.Equal(result, tc.body) {
					t.Errorf("Expected no conversion, but content was changed")
				}
			}
		})
	}
}

// TestConvertEncoding_AutoDetection は自動検出のテストです
func TestConvertEncoding_AutoDetection(t *testing.T) {
	repo := NewAPIRepository()

	testCases := []struct {
		name           string
		htmlContent    string
		contentType    string
		expectShiftJIS bool
	}{
		{
			name:           "Content-Typeからcharset検出",
			htmlContent:    `<html><head><title>テスト</title></head></html>`,
			contentType:    "text/html; charset=shift_jis",
			expectShiftJIS: true,
		},
		{
			name:           "HTMLメタタグからcharset検出",
			htmlContent:    `<html><head><meta charset="shift_jis"><title>テスト</title></head></html>`,
			contentType:    "text/html",
			expectShiftJIS: true,
		},
		{
			name:           "http-equiv形式のメタタグ検出",
			htmlContent:    `<html><head><meta http-equiv="Content-Type" content="text/html; charset=shift_jis"><title>テスト</title></head></html>`,
			contentType:    "text/html",
			expectShiftJIS: true,
		},
		{
			name:           "UTF-8の場合は変換しない",
			htmlContent:    `<html><head><meta charset="utf-8"><title>テスト</title></head></html>`,
			contentType:    "text/html",
			expectShiftJIS: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			var err error

			if tc.expectShiftJIS {
				// Shift_JISにエンコード
				body, err = encodeToShiftJIS(tc.htmlContent)
				if err != nil {
					t.Fatalf("Failed to encode to Shift_JIS: %v", err)
				}
			} else {
				// UTF-8のまま
				body = []byte(tc.htmlContent)
			}

			// Act
			result, err := repo.convertEncoding(body, tc.contentType, "auto")

			// Assert
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			resultStr := string(result)
			if tc.expectShiftJIS {
				// Shift_JISから変換されて日本語が正しく表示されることを確認
				if !strings.Contains(resultStr, "テスト") {
					t.Errorf("Expected converted text to contain Japanese characters, got %s", resultStr)
				}
			} else {
				// UTF-8の場合はそのまま
				if !strings.Contains(resultStr, "テスト") {
					t.Errorf("Expected UTF-8 text to remain unchanged, got %s", resultStr)
				}
			}
		})
	}
}

// TestSendRequest_WithEncoding はエンコーディング指定付きリクエストのテストです
func TestSendRequest_WithEncoding(t *testing.T) {
	// Shift_JISエンコードされたHTMLレスポンスを返すテストサーバー
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmlContent := `<html><head><title>阿部寛のホームページ</title></head><body>こんにちは</body></html>`
		shiftJISBytes, _ := encodeToShiftJIS(htmlContent)

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		w.Write(shiftJISBytes)
	}))
	defer server.Close()

	repo := NewAPIRepository()

	testCases := []struct {
		name             string
		encoding         string
		expectConversion bool
	}{
		{
			name:             "Shift_JIS指定で変換",
			encoding:         "shift_jis",
			expectConversion: true,
		},
		{
			name:             "UTF-8指定で変換なし",
			encoding:         "utf-8",
			expectConversion: false,
		},
		{
			name:             "auto指定で自動検出",
			encoding:         "auto",
			expectConversion: false, // Content-Typeにcharsetがないため変換されない
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			request := &models.APIRequest{
				URL:      server.URL,
				Method:   "GET",
				Headers:  map[string]string{"Accept": "text/html"},
				Body:     nil,
				Encoding: tc.encoding,
			}

			// Act
			response, err := repo.SendRequest(request)

			// Assert
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if response.StatusCode != 200 {
				t.Errorf("Expected status code 200, got %d", response.StatusCode)
			}

			bodyStr := string(response.Body)
			if tc.expectConversion {
				// 変換されて日本語が正しく表示されることを確認
				if !strings.Contains(bodyStr, "阿部寛") {
					t.Errorf("Expected converted text to contain Japanese characters, got %s", bodyStr)
				}
			}
		})
	}
}

// TestExtractCharsetFromContentType はContent-Typeからのcharset抽出テストです
func TestExtractCharsetFromContentType(t *testing.T) {
	repo := NewAPIRepository()

	testCases := []struct {
		name        string
		contentType string
		expected    string
	}{
		{
			name:        "標準的なcharset指定",
			contentType: "text/html; charset=shift_jis",
			expected:    "shift_jis",
		},
		{
			name:        "UTF-8指定",
			contentType: "text/html; charset=utf-8",
			expected:    "utf-8",
		},
		{
			name:        "複数パラメータ付き",
			contentType: "text/html; charset=shift_jis; boundary=something",
			expected:    "shift_jis",
		},
		{
			name:        "charset指定なし",
			contentType: "text/html",
			expected:    "",
		},
		{
			name:        "空文字列",
			contentType: "",
			expected:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := repo.extractCharsetFromContentType(tc.contentType)

			// Assert
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestExtractCharsetFromHTML はHTMLからのcharset抽出テストです
func TestExtractCharsetFromHTML(t *testing.T) {
	repo := NewAPIRepository()

	testCases := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "meta charset形式",
			html:     `<html><head><meta charset="shift_jis"><title>Test</title></head></html>`,
			expected: "shift_jis",
		},
		{
			name:     "http-equiv形式",
			html:     `<html><head><meta http-equiv="Content-Type" content="text/html; charset=shift_jis"><title>Test</title></head></html>`,
			expected: "shift_jis",
		},
		{
			name:     "content属性が先の形式",
			html:     `<html><head><meta content="text/html; charset=shift_jis" http-equiv="Content-Type"><title>Test</title></head></html>`,
			expected: "shift_jis",
		},
		{
			name:     "charset指定なし",
			html:     `<html><head><title>Test</title></head></html>`,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := repo.extractCharsetFromHTML(tc.html)

			// Assert
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestConvertEncoding_SpecificShiftJISBytes は特定のShift_JISバイト列のテストです
func TestConvertEncoding_SpecificShiftJISBytes(t *testing.T) {
	repo := NewAPIRepository()

	// 「阿部寛のホームページ」のShift_JISバイト列
	// [136 162 149 148 138 176 130 204 131 122 129 91 131 128 131 121 129 91 131 87]
	shiftJISBytes := []byte{136, 162, 149, 148, 138, 176, 130, 204, 131, 122, 129, 91, 131, 128, 131, 121, 129, 91, 131, 87}

	// HTMLコンテンツとして組み立て
	htmlPrefix := []byte(`<html><head><title>`)
	htmlSuffix := []byte(`</title></head><body>test</body></html>`)
	htmlContent := append(append(htmlPrefix, shiftJISBytes...), htmlSuffix...)

	testCases := []struct {
		name              string
		body              []byte
		contentType       string
		specifiedEncoding string
		expectedText      string
	}{
		{
			name:              "特定のShift_JISバイト列をshift_jis指定で変換",
			body:              htmlContent,
			contentType:       "text/html",
			specifiedEncoding: "shift_jis",
			expectedText:      "阿部寛のホームページ",
		},
		{
			name:              "特定のShift_JISバイト列をshift-jis指定で変換",
			body:              htmlContent,
			contentType:       "text/html",
			specifiedEncoding: "shift-jis",
			expectedText:      "阿部寛のホームページ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := repo.convertEncoding(tc.body, tc.contentType, tc.specifiedEncoding)

			// Assert
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			resultStr := string(result)
			if !strings.Contains(resultStr, tc.expectedText) {
				t.Errorf("Expected converted text to contain '%s', got %s", tc.expectedText, resultStr)
			}

			// より詳細な検証：タイトル部分を抽出して確認
			titleStart := strings.Index(resultStr, "<title>")
			titleEnd := strings.Index(resultStr, "</title>")
			if titleStart != -1 && titleEnd != -1 {
				titleContent := resultStr[titleStart+7 : titleEnd]
				if titleContent != tc.expectedText {
					t.Errorf("Expected title to be '%s', got '%s'", tc.expectedText, titleContent)
				}
			} else {
				t.Error("Could not find title tags in converted HTML")
			}
		})
	}
}

// encodeToShiftJIS はテキストをShift_JISにエンコードするヘルパー関数です
func encodeToShiftJIS(text string) ([]byte, error) {
	encoder := japanese.ShiftJIS.NewEncoder()
	reader := transform.NewReader(strings.NewReader(text), encoder)
	return io.ReadAll(reader)
}
