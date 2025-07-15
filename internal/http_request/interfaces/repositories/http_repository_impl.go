package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	models "github.com/landmaster135/devbox/internal/http_request/domain/models"
)

// HTTPRepositoryImpl はHTTPRepositoryインターフェースの実装です
type HTTPRepositoryImpl struct {
	client    *http.Client
	userAgent string
}

// NewHTTPRepository は新しいHTTPRepositoryインスタンスを作成します
func NewHTTPRepository() *HTTPRepositoryImpl {
	return &HTTPRepositoryImpl{
		client:    &http.Client{},
		userAgent: buildDefaultUserAgent(),
	}
}

// buildDefaultUserAgent はデフォルトのUser-Agentを構築します
func buildDefaultUserAgent() string {
	return fmt.Sprintf(
		"HTTP-Request-CLI/1.0 (%s %s; %s) Go/%s",
		runtime.GOOS,
		runtime.GOARCH,
		getOSVersion(),
		runtime.Version())
}

// getOSVersion はOS情報を取得します
func getOSVersion() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows NT 10.0; Win64; x64"
	case "darwin":
		return "Macintosh; Intel Mac OS X 10_15_7"
	case "linux":
		return "X11; Linux x86_64"
	default:
		return runtime.GOOS
	}
}

// SendRequest はHTTPリクエストを送信し、レスポンスを返します
func (r *HTTPRepositoryImpl) SendRequest(request *models.HTTPRequest) (*models.HTTPResponse, error) {
	// HTTPリクエストを作成
	req, err := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(request.Body))
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの作成に失敗しました: %w", err)
	}

	// デフォルトのUser-Agentを設定（Windows Defender誤検知対策）
	req.Header.Set("User-Agent", r.userAgent)

	// ヘッダーを設定（User-Agentが明示的に指定されている場合は上書き）
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	// リクエストを送信
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンスボディの読み込みに失敗しました: %w", err)
	}

	// レスポンスヘッダーをマップに変換
	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// 文字エンコーディングを検出・変換
	convertedBody, err := r.convertEncoding(body, headers["Content-Type"], request.Encoding)
	if err != nil {
		// 変換に失敗した場合は元のボディをそのまま使用
		convertedBody = body
	}

	// HTTPResponseを作成して返す
	return &models.HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       convertedBody,
	}, nil
}

// LoadJSONFile はJSONファイルを読み込み、バイト配列として返します
func (r *HTTPRepositoryImpl) LoadJSONFile(filePath string) ([]byte, error) {
	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルを開けませんでした: %w", err)
	}
	defer file.Close()

	// ファイルの内容を読み込む
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	// BOM（Byte Order Mark）を削除
	if len(content) >= 3 && content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		content = content[3:]
	}

	// JSONとして有効かチェック
	var js interface{}
	if err := json.Unmarshal(content, &js); err != nil {
		return nil, fmt.Errorf("無効なJSON形式です: %w", err)
	}

	return content, nil
}

// convertEncoding は文字エンコーディングを検出・変換します
func (r *HTTPRepositoryImpl) convertEncoding(body []byte, contentType string, specifiedEncoding string) ([]byte, error) {
	// HTMLの場合のみ文字エンコーディング変換を実行
	if !strings.Contains(strings.ToLower(contentType), "text/html") {
		return body, nil
	}

	var charset string

	// 指定されたエンコーディングがある場合は優先使用
	if specifiedEncoding != "" && specifiedEncoding != "auto" {
		charset = specifiedEncoding
	} else {
		// Content-Typeヘッダーからcharsetを検出
		charset = r.extractCharsetFromContentType(contentType)

		// Content-Typeにcharsetがない場合、HTMLのmetaタグから検出
		if charset == "" {
			charset = r.extractCharsetFromHTML(string(body))
		}
	}

	// Shift_JISの場合、UTF-8に変換
	if strings.ToLower(charset) == "shift_jis" || strings.ToLower(charset) == "shift-jis" {
		reader := transform.NewReader(bytes.NewReader(body), japanese.ShiftJIS.NewDecoder())
		converted, err := io.ReadAll(reader)
		if err != nil {
			return body, fmt.Errorf("Shift_JISからUTF-8への変換に失敗しました: %w", err)
		}
		return converted, nil
	}

	return body, nil
}

// extractCharsetFromContentType はContent-Typeヘッダーからcharsetを抽出します
func (r *HTTPRepositoryImpl) extractCharsetFromContentType(contentType string) string {
	re := regexp.MustCompile(`charset=([^;]+)`)
	matches := re.FindStringSubmatch(contentType)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// extractCharsetFromHTML はHTMLのmetaタグからcharsetを抽出します
func (r *HTTPRepositoryImpl) extractCharsetFromHTML(html string) string {
	// 大文字小文字を区別しない検索のため、小文字に変換
	lowerHTML := strings.ToLower(html)

	// <meta charset="..."> パターン
	re1 := regexp.MustCompile(`<meta\s+charset\s*=\s*["']?([^"'>\s]+)["']?`)
	matches := re1.FindStringSubmatch(lowerHTML)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// <meta http-equiv="Content-Type" content="text/html; charset=..."> パターン
	re2 := regexp.MustCompile(`<meta\s+http-equiv\s*=\s*["']?content-type["']?\s+content\s*=\s*["'][^"']*charset=([^"';\s]+)`)
	matches = re2.FindStringSubmatch(lowerHTML)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// content属性が先に来るパターンも対応
	re3 := regexp.MustCompile(`<meta\s+content\s*=\s*["'][^"']*charset=([^"';\s]+)[^"']*["']\s+http-equiv\s*=\s*["']?content-type["']?`)
	matches = re3.FindStringSubmatch(lowerHTML)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// 日本語サイトでよく使われるShift_JISを推測
	// 文字化けパターンが見つかった場合、Shift_JISと推測
	if strings.Contains(html, "�") {
		return "shift_jis"
	}

	return ""
}
