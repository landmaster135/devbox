package steam_api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// buildURLWithParams はAPIキーとパラメータを含むURLを構築します
func buildURLWithParams(baseURL, apiKey string, params map[string]interface{}) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	query := u.Query()
	query.Set("key", apiKey)

	for key, value := range cleanParams(params) {
		query.Set(key, value)
	}

	u.RawQuery = query.Encode()
	return u.String()
}

// buildURLWithParamsForSearch は検索用のURLを構築します（APIキーなし）
func buildURLWithParamsForSearch(baseURL, searchTerm string, params map[string]interface{}) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	query := u.Query()
	query.Set("term", searchTerm)

	for key, value := range cleanParams(params) {
		query.Set(key, value)
	}

	u.RawQuery = query.Encode()
	return u.String()
}

// cleanParams はパラメータを文字列に変換し、nilや空の値を除外します
func cleanParams(params map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for key, value := range params {
		if value == nil {
			continue
		}

		switch v := value.(type) {
		case string:
			if v != "" {
				result[key] = v
			}
		case int:
			result[key] = strconv.Itoa(v)
		case int64:
			result[key] = strconv.FormatInt(v, 10)
		case bool:
			if v {
				result[key] = "true"
			} else {
				result[key] = "false"
			}
		case []string:
			if len(v) > 0 {
				result[key] = strings.Join(v, ",")
			}
		case []int:
			if len(v) > 0 {
				strSlice := make([]string, len(v))
				for i, num := range v {
					strSlice[i] = strconv.Itoa(num)
				}
				result[key] = strings.Join(strSlice, ",")
			}
		default:
			// その他の型は文字列として扱う
			result[key] = fmt.Sprintf("%v", v)
		}
	}

	return result
}

// mergeParams は2つのパラメータマップをマージします
func mergeParams(base, additional map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// ベースパラメータをコピー
	for key, value := range base {
		result[key] = value
	}

	// 追加パラメータをマージ（上書き）
	for key, value := range additional {
		result[key] = value
	}

	return result
}

// validateSteamID はSteam IDの形式を検証します
func validateSteamID(steamID string) bool {
	if steamID == "" {
		return false
	}

	// Steam ID は通常17桁の数字
	if len(steamID) != 17 {
		return false
	}

	// 数字のみかチェック
	for _, char := range steamID {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// validateAppID はApp IDの形式を検証します
func validateAppID(appID int) bool {
	return appID > 0
}

// parseCommaSeparatedInts はカンマ区切りの文字列を整数スライスに変換します
func parseCommaSeparatedInts(s string) ([]int, error) {
	if s == "" {
		return []int{}, nil
	}

	parts := strings.Split(s, ",")
	result := make([]int, len(parts))

	for i, part := range parts {
		num, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", part)
		}
		result[i] = num
	}

	return result, nil
}

// formatSteamIDList は複数のSteam IDを適切な形式に変換します
func formatSteamIDList(steamIDs []string) string {
	validIDs := make([]string, 0, len(steamIDs))

	for _, id := range steamIDs {
		if validateSteamID(id) {
			validIDs = append(validIDs, id)
		}
	}

	return strings.Join(validIDs, ",")
}

// chunkSteamIDs はSteam IDのリストを指定されたサイズのチャンクに分割します
// Steam APIは一度に処理できるIDの数に制限があるため
func chunkSteamIDs(steamIDs []string, chunkSize int) [][]string {
	if chunkSize <= 0 {
		chunkSize = 100 // デフォルトのチャンクサイズ
	}

	var chunks [][]string
	for i := 0; i < len(steamIDs); i += chunkSize {
		end := i + chunkSize
		if end > len(steamIDs) {
			end = len(steamIDs)
		}
		chunks = append(chunks, steamIDs[i:end])
	}

	return chunks
}
