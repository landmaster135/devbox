package detectors

import (
	"fmt"
	"strings"
)

const bodyScanLimit = 8 << 10 // 8 KiB

var (
	bodyIndicators = []string{
		"just a moment",
		"checking your browser",
		"enable javascript",
		"cloudflare",
		"cf-browser-verification",
		"redirecting",
		"attention required",
		"security check",
	}
	headerIndicators = []string{
		"cf-ray",
		"cf-mitigated",
		"cf-chl-bypass",
		"cf-cache-status",
		"cf-challenge-status",
	}
)

// IsCloudflareChallenge はレスポンスからCloudflareのボットチャレンジが発生しているかを推測します。

func IsCloudflareChallenge(statusCode int, headers map[string]string, body []byte) bool {
	normalized := normalizeHeaders(headers)
	if !isSuspiciousStatus(statusCode) && !strings.Contains(strings.ToLower(normalized["server"]), "cloudflare") {
		return false
	}

	if hasHeaderIndicators(normalized) {
		return true
	}

	snippet := strings.ToLower(string(bodySample(body)))
	for _, indicator := range bodyIndicators {
		if strings.Contains(snippet, indicator) {
			if strings.Contains(strings.ToLower(normalized["server"]), "cloudflare") || hasHeaderIndicators(normalized) {
				return true
			}
		}
	}

	return false
}

// BuildCloudflareWarning はユーザー向けの警告文を作成します。
func BuildCloudflareWarning(statusCode int, headers map[string]string) string {
	normalized := normalizeHeaders(headers)

	message := fmt.Sprintf("CloudflareのBot検出によりステータスコード%dが返却されました", statusCode)

	details := make([]string, 0, 2)
	if rayID := normalized["cf-ray"]; rayID != "" {
		details = append(details, "Ray ID: "+rayID)
	}
	if mitigation := normalized["cf-mitigated"]; mitigation != "" {
		details = append(details, "Mitigation: "+mitigation)
	}
	if len(details) > 0 {
		message += " (" + strings.Join(details, ", ") + ")"
	}

	return message + "。ブラウザに近いヘッダーや十分な待機時間を設けても改善しない場合は手動でページへアクセスしてください。"
}

func bodySample(body []byte) []byte {
	if len(body) <= bodyScanLimit {
		return body
	}
	return body[:bodyScanLimit]
}

func normalizeHeaders(headers map[string]string) map[string]string {
	normalized := make(map[string]string, len(headers))
	for key, value := range headers {
		normalized[strings.ToLower(key)] = value
	}
	return normalized
}

func hasHeaderIndicators(headers map[string]string) bool {
	for _, key := range headerIndicators {
		if headers[key] != "" {
			return true
		}
	}
	server := strings.ToLower(headers["server"])
	return strings.Contains(server, "cloudflare") && (headers["cf-ray"] != "" || headers["cf-mitigated"] != "")
}

func isSuspiciousStatus(statusCode int) bool {
	if statusCode == 403 || statusCode == 503 || statusCode == 429 {
		return true
	}
	return statusCode >= 520 && statusCode <= 524
}
