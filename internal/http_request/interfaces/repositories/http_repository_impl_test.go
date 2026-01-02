package repositories

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andybalholm/brotli"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	models "github.com/landmaster135/devbox/internal/http_request/domain/models"
)

// TestNewHTTPRepository はNewHTTPRepositoryメソッドのテストです
func TestNewHTTPRepository(t *testing.T) {
	// Act
	repo := NewHTTPRepository()

	// Assert
	if repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}
	if repo.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
	if repo.userAgent == "" {
		t.Error("Expected User-Agent to be initialized")
	}
	if len(repo.defaultHeaders) == 0 {
		t.Error("Expected default headers to be initialized")
	}
	if repo.retryPolicy.maxAttempts < 1 {
		t.Error("Expected retry policy to be configured")
	}
}

// TestBuildDefaultUserAgent はデフォルトUser-Agent構築のテストです
func TestBuildDefaultUserAgent(t *testing.T) {
	// Act
	userAgent := buildDefaultUserAgent()

	// Assert
	if userAgent == "" {
		t.Fatal("Expected User-Agent to be non-empty")
	}

	// User-Agentに必要な情報が含まれていることを確認
	if !strings.Contains(userAgent, "Mozilla/5.0") {
		t.Error("Expected User-Agent to contain browser prefix")
	}
	if !strings.Contains(userAgent, "Chrome/") {
		t.Error("Expected User-Agent to contain Chrome token")
	}
	if !strings.Contains(userAgent, "Safari/537.36") {
		t.Error("Expected User-Agent to contain Safari token")
	}
	if !strings.Contains(userAgent, getOSVersion()) {
		t.Error("Expected User-Agent to contain OS information")
	}
}

// TestGetOSVersion はOS情報取得のテストです
func TestGetOSVersion(t *testing.T) {
	// Act
	osVersion := getOSVersion()

	// Assert
	if osVersion == "" {
		t.Fatal("Expected OS version to be non-empty")
	}

	// 現在のOSに応じた適切な情報が返されることを確認
	switch runtime.GOOS {
	case "windows":
		if !strings.Contains(osVersion, "Windows NT") {
			t.Error("Expected Windows OS version to contain 'Windows NT'")
		}
	case "darwin":
		if !strings.Contains(osVersion, "Macintosh") {
			t.Error("Expected macOS version to contain 'Macintosh'")
		}
	case "linux":
		if !strings.Contains(osVersion, "Linux") {
			t.Error("Expected Linux version to contain 'Linux'")
		}
	default:
		if osVersion != runtime.GOOS {
			t.Errorf("Expected OS version to be %s for unknown OS, got %s", runtime.GOOS, osVersion)
		}
	}
}

// TestSendRequest_UserAgentSet はUser-Agentが正しく設定されるテストです
func TestSendRequest_UserAgentSet(t *testing.T) {
	// HTTPサーバーのモックを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")

		// User-Agentが設定されていることを確認
		if userAgent == "" {
			t.Error("Expected User-Agent header to be set")
		}

		// User-Agentに必要な情報が含まれていることを確認
		if !strings.Contains(userAgent, "Mozilla/5.0") {
			t.Errorf("Expected User-Agent to contain browser prefix, got: %s", userAgent)
		}
		if !strings.Contains(userAgent, "Chrome/") {
			t.Errorf("Expected User-Agent to contain Chrome token, got: %s", userAgent)
		}
		if !strings.Contains(userAgent, getOSVersion()) {
			t.Errorf("Expected User-Agent to contain OS info, got: %s", userAgent)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	repo := NewHTTPRepository()
	request := &models.HTTPRequest{
		URL:     server.URL,
		Method:  "GET",
		Headers: map[string]string{},
	}

	// Act
	_, err := repo.SendRequest(request)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// TestSendRequest_CustomUserAgent はカスタムUser-Agentが上書きされるテストです
func TestSendRequest_CustomUserAgent(t *testing.T) {
	customUserAgent := "Custom-Client/2.0"

	// HTTPサーバーのモックを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")

		// カスタムUser-Agentが設定されていることを確認
		if userAgent != customUserAgent {
			t.Errorf("Expected User-Agent to be %s, got: %s", customUserAgent, userAgent)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	repo := NewHTTPRepository()
	request := &models.HTTPRequest{
		URL:    server.URL,
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": customUserAgent,
		},
	}

	// Act
	_, err := repo.SendRequest(request)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSendRequest_DefaultBrowserHeaders(t *testing.T) {
	var acceptHeader, languageHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptHeader = r.Header.Get("Accept")
		languageHeader = r.Header.Get("Accept-Language")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	repo := NewHTTPRepository()
	request := &models.HTTPRequest{
		URL:     server.URL,
		Method:  "GET",
		Headers: map[string]string{},
	}

	if _, err := repo.SendRequest(request); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(acceptHeader, "text/html") {
		t.Errorf("expected Accept header to include text/html, got %s", acceptHeader)
	}
	if !strings.Contains(languageHeader, "ja") {
		t.Errorf("expected Accept-Language header to prioritize ja, got %s", languageHeader)
	}
}

func TestSendRequest_CustomAcceptOverridesDefault(t *testing.T) {
	const customAccept = "application/json"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if header := r.Header.Get("Accept"); header != customAccept {
			t.Fatalf("expected Accept=%s, got %s", customAccept, header)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	repo := NewHTTPRepository()
	request := &models.HTTPRequest{
		URL:    server.URL,
		Method: "GET",
		Headers: map[string]string{
			"Accept": customAccept,
		},
	}

	if _, err := repo.SendRequest(request); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendRequest_DecodesGzipResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gz.Write([]byte("hello gzip"))
	}))
	defer server.Close()

	repo := NewHTTPRepository()
	request := &models.HTTPRequest{URL: server.URL, Method: "GET", Headers: map[string]string{}}

	resp, err := repo.SendRequest(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(resp.Body) != "hello gzip" {
		t.Errorf("expected body to be decoded, got %s", string(resp.Body))
	}
}

func TestSendRequest_DecodesBrotliResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Encoding", "br")
		writer := brotli.NewWriter(w)
		defer writer.Close()
		writer.Write([]byte("hello brotli"))
	}))
	defer server.Close()

	repo := NewHTTPRepository()
	request := &models.HTTPRequest{URL: server.URL, Method: "GET", Headers: map[string]string{}}

	resp, err := repo.SendRequest(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(resp.Body) != "hello brotli" {
		t.Errorf("expected brotli body to be decoded, got %s", string(resp.Body))
	}
}

func TestSendRequest_AddsCloudflareWarning(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "cloudflare")
		w.Header().Set("Cf-Mitigated", "challenge")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Just a moment"))
	}))
	defer server.Close()

	repo := NewHTTPRepository()
	repo.retryPolicy.maxAttempts = 1
	request := &models.HTTPRequest{URL: server.URL, Method: "GET", Headers: map[string]string{}}

	resp, err := repo.SendRequest(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Warnings) == 0 {
		t.Fatalf("expected warning to be populated")
	}
	if !strings.Contains(resp.Warnings[0], "Cloudflare") {
		t.Errorf("expected warning to mention Cloudflare, got %s", resp.Warnings[0])
	}
}

func TestSendRequest_RetriesOnCloudflareForbidden(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt32(&callCount, 1)
		if current == 1 {
			w.Header().Set("Server", "cloudflare")
			w.Header().Set("Cf-Mitigated", "challenge")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Just a moment"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	repo := NewHTTPRepository()
	repo.retryPolicy.initialBackoff = time.Millisecond
	repo.retryPolicy.maxBackoff = 2 * time.Millisecond
	request := &models.HTTPRequest{URL: server.URL, Method: "GET", Headers: map[string]string{}}

	resp, err := repo.SendRequest(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected final status 200, got %d", resp.StatusCode)
	}
	attempts := atomic.LoadInt32(&callCount)
	if attempts < 2 {
		t.Fatalf("expected at least two attempts, got %d", attempts)
	}
}

// TestConvertEncoding_ShiftJIS はShift_JIS変換のテストです
func TestConvertEncoding_ShiftJIS(t *testing.T) {
	repo := NewHTTPRepository()

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
	repo := NewHTTPRepository()

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

	repo := NewHTTPRepository()

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
			request := &models.HTTPRequest{
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
	repo := NewHTTPRepository()

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
	repo := NewHTTPRepository()

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
	repo := NewHTTPRepository()

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
