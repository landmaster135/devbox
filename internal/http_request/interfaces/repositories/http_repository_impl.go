package repositories

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/landmaster135/devbox/internal/http_request/domain/detectors"
	models "github.com/landmaster135/devbox/internal/http_request/domain/models"
)

// HTTPRepositoryImpl はHTTPRepositoryインターフェースの実装です
type HTTPRepositoryImpl struct {
	client         *http.Client
	userAgent      string
	defaultHeaders map[string]string
	retryPolicy    retryPolicy
}

type retryPolicy struct {
	maxAttempts    int
	initialBackoff time.Duration
	maxBackoff     time.Duration
}

const (
	maxResponseBodySize = 10 << 20 // 10 MiB
	sniffBufferSize     = 4 << 10  // 4 KiB
)

var (
	charsetContentTypeRegex = regexp.MustCompile(`charset=([^;]+)`)
	metaCharsetRegex        = regexp.MustCompile(`<meta\s+charset\s*=\s*["']?([^"'>\s]+)["']?`)
	metaHTTPEquivRegex      = regexp.MustCompile(`<meta\s+http-equiv\s*=\s*["']?content-type["']?\s+content\s*=\s*["'][^"']*charset=([^"';\s]+)`)
	metaContentFirstRegex   = regexp.MustCompile(`<meta\s+content\s*=\s*["'][^"']*charset=([^"';\s]+)[^"']*["']\s+http-equiv\s*=\s*["']?content-type["']?`)
)

// NewHTTPRepository は新しいHTTPRepositoryインスタンスを作成します
func NewHTTPRepository() *HTTPRepositoryImpl {
	return &HTTPRepositoryImpl{
		client:         &http.Client{},
		userAgent:      buildDefaultUserAgent(),
		defaultHeaders: buildDefaultHeaders(),
		retryPolicy: retryPolicy{
			maxAttempts:    3,
			initialBackoff: 500 * time.Millisecond,
			maxBackoff:     5 * time.Second,
		},
	}
}

func buildDefaultHeaders() map[string]string {
	return map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language":           "ja,en-US;q=0.9,en;q=0.8",
		"Accept-Encoding":           "gzip, deflate, br",
		"Cache-Control":             "no-cache",
		"Pragma":                    "no-cache",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
		"Upgrade-Insecure-Requests": "1",
		"Connection":                "keep-alive",
	}
}

// buildDefaultUserAgent はブラウザ風のUser-Agentを返します
func buildDefaultUserAgent() string {
	return fmt.Sprintf(
		"Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		getOSVersion(),
	)
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
	retryable := isRetryableMethod(request.Method)
	attempts := 1
	if retryable {
		attempts = r.retryPolicy.maxAttempts
	}

	var (
		response   *models.HTTPResponse
		cfDetected bool
		err        error
	)

	for attempt := 0; attempt < attempts; attempt++ {
		response, cfDetected, err = r.executeRequest(request)
		if err != nil {
			if !retryable || attempt == attempts-1 {
				return nil, err
			}
		} else {
			if !retryable || !r.shouldRetryResponse(response, cfDetected) || attempt == attempts-1 {
				return response, nil
			}
		}

		if retryable {
			time.Sleep(r.calculateBackoffDuration(response, attempt))
		}
	}

	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *HTTPRepositoryImpl) executeRequest(request *models.HTTPRequest) (*models.HTTPResponse, bool, error) {
	req, err := http.NewRequest(request.Method, request.URL, bytes.NewReader(request.Body))
	if err != nil {
		return nil, false, fmt.Errorf("HTTPリクエストの作成に失敗しました: %w", err)
	}

	r.applyHeaders(req, request.Headers)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("HTTPリクエストの送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	bodyReader, decompressor, err := r.buildBodyReader(resp)
	if err != nil {
		return nil, false, err
	}
	if decompressor != nil {
		defer decompressor.Close()
	}

	contentType := resp.Header.Get("Content-Type")
	limited := &io.LimitedReader{R: bodyReader, N: maxResponseBodySize + 1}
	reader := bufio.NewReader(limited)

	body, convErr := r.readResponseBody(reader, contentType, request.Encoding)
	if limited.N <= 0 {
		return nil, false, fmt.Errorf("レスポンスボディが最大サイズ(%dバイト)を超えました", maxResponseBodySize)
	}
	if convErr != nil {
		// 文字コード変換に失敗した場合はエラーを無視して未変換のボディを返す
	}

	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	cfDetected := detectors.IsCloudflareChallenge(resp.StatusCode, headers, body)
	warnings := make([]string, 0)
	if cfDetected {
		warnings = append(warnings, detectors.BuildCloudflareWarning(resp.StatusCode, headers))
	}

	return &models.HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       body,
		Warnings:   warnings,
	}, cfDetected, nil
}

func (r *HTTPRepositoryImpl) applyHeaders(req *http.Request, custom map[string]string) {
	req.Header.Set("User-Agent", r.userAgent)
	for key, value := range r.defaultHeaders {
		if req.Header.Get(key) == "" {
			req.Header.Set(key, value)
		}
	}
	for key, value := range custom {
		req.Header.Set(key, value)
	}
}

func (r *HTTPRepositoryImpl) buildBodyReader(resp *http.Response) (io.Reader, io.Closer, error) {
	encoding := normalizeEncoding(resp.Header.Get("Content-Encoding"))
	switch encoding {
	case "gzip":
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("gzip圧縮の展開に失敗しました: %w", err)
		}
		return gz, gz, nil
	case "deflate":
		flateReader := flate.NewReader(resp.Body)
		return flateReader, flateReader, nil
	case "br":
		return brotli.NewReader(resp.Body), nil, nil
	default:
		return resp.Body, nil, nil
	}
}

func normalizeEncoding(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	if len(parts) == 0 {
		return strings.ToLower(strings.TrimSpace(value))
	}
	return strings.ToLower(strings.TrimSpace(parts[0]))
}

func (r *HTTPRepositoryImpl) shouldRetryResponse(response *models.HTTPResponse, cfDetected bool) bool {
	if response == nil {
		return false
	}

	if cfDetected && response.StatusCode == http.StatusForbidden {
		return true
	}

	switch response.StatusCode {
	case http.StatusTooManyRequests, http.StatusServiceUnavailable, http.StatusGatewayTimeout, http.StatusBadGateway:
		return true
	}

	if response.StatusCode >= 520 && response.StatusCode <= 524 {
		return true
	}

	return false
}

func (r *HTTPRepositoryImpl) calculateBackoffDuration(response *models.HTTPResponse, attempt int) time.Duration {
	if response != nil {
		if retryAfter, ok := parseRetryAfter(response.Headers); ok {
			if retryAfter <= 0 {
				return r.retryPolicy.initialBackoff
			}
			if retryAfter > r.retryPolicy.maxBackoff {
				return r.retryPolicy.maxBackoff
			}
			return retryAfter
		}
	}

	backoff := r.retryPolicy.initialBackoff * time.Duration(1<<attempt)
	if backoff > r.retryPolicy.maxBackoff {
		return r.retryPolicy.maxBackoff
	}
	if backoff <= 0 {
		return r.retryPolicy.initialBackoff
	}
	return backoff
}

func parseRetryAfter(headers map[string]string) (time.Duration, bool) {
	for key, value := range headers {
		if strings.EqualFold(key, "Retry-After") {
			if seconds, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
				return time.Duration(seconds) * time.Second, true
			}
			if when, err := http.ParseTime(value); err == nil {
				return time.Until(when), true
			}
		}
	}
	return 0, false
}

func isRetryableMethod(method string) bool {
	switch strings.ToUpper(method) {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

// LoadJSONFile はJSONファイルを読み込み、バイト配列として返します
func (r *HTTPRepositoryImpl) LoadJSONFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルを開けませんでした: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルの読み込みに失敗しました: %w", err)
	}

	if len(content) >= 3 && content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		content = content[3:]
	}

	var js interface{}
	if err := json.Unmarshal(content, &js); err != nil {
		return nil, fmt.Errorf("無効なJSON形式です: %w", err)
	}

	return content, nil
}

// readResponseBody はレスポンスボディを読み込み、必要に応じてエンコーディングを変換します
func (r *HTTPRepositoryImpl) readResponseBody(reader *bufio.Reader, contentType, specifiedEncoding string) ([]byte, error) {
	isHTML := strings.Contains(strings.ToLower(contentType), "text/html")
	charset := r.detectCharset(reader, contentType, specifiedEncoding, isHTML)

	if isHTML && isShiftJIS(charset) {
		var rawBuf bytes.Buffer
		tee := io.TeeReader(reader, &rawBuf)
		decoded := transform.NewReader(tee, japanese.ShiftJIS.NewDecoder())
		var convertedBuf bytes.Buffer
		if _, err := convertedBuf.ReadFrom(decoded); err != nil {
			return rawBuf.Bytes(), fmt.Errorf("Shift_JISからUTF-8への変換に失敗しました: %w", err)
		}
		return convertedBuf.Bytes(), nil
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, fmt.Errorf("レスポンスボディの読み込みに失敗しました: %w", err)
	}
	return buf.Bytes(), nil
}

// convertEncoding はテスト互換性のために残しているヘルパー
func (r *HTTPRepositoryImpl) convertEncoding(body []byte, contentType string, specifiedEncoding string) ([]byte, error) {
	limited := &io.LimitedReader{R: bytes.NewReader(body), N: maxResponseBodySize + 1}
	reader := bufio.NewReader(limited)
	converted, err := r.readResponseBody(reader, contentType, specifiedEncoding)
	if limited.N <= 0 {
		return nil, fmt.Errorf("レスポンスボディが最大サイズ(%dバイト)を超えました", maxResponseBodySize)
	}
	return converted, err
}

// detectCharset はレスポンスのcharsetを推測します
func (r *HTTPRepositoryImpl) detectCharset(reader *bufio.Reader, contentType, specifiedEncoding string, isHTML bool) string {
	if specifiedEncoding != "" && specifiedEncoding != "auto" {
		return specifiedEncoding
	}

	if charset := r.extractCharsetFromContentType(contentType); charset != "" {
		return charset
	}

	if !isHTML {
		return ""
	}

	peek, _ := reader.Peek(sniffBufferSize)
	if len(peek) > 0 {
		if charset := r.extractCharsetFromHTML(string(peek)); charset != "" {
			return charset
		}
	}

	return ""
}

func isShiftJIS(charset string) bool {
	switch strings.ToLower(strings.TrimSpace(charset)) {
	case "shift_jis", "shift-jis":
		return true
	default:
		return false
	}
}

// extractCharsetFromContentType はContent-Typeヘッダーからcharsetを抽出します
func (r *HTTPRepositoryImpl) extractCharsetFromContentType(contentType string) string {
	matches := charsetContentTypeRegex.FindStringSubmatch(contentType)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// extractCharsetFromHTML はHTMLのmetaタグからcharsetを抽出します
func (r *HTTPRepositoryImpl) extractCharsetFromHTML(html string) string {
	lowerHTML := strings.ToLower(html)

	matches := metaCharsetRegex.FindStringSubmatch(lowerHTML)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	matches = metaHTTPEquivRegex.FindStringSubmatch(lowerHTML)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	matches = metaContentFirstRegex.FindStringSubmatch(lowerHTML)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	if strings.Contains(html, "�") {
		return "shift_jis"
	}

	return ""
}
