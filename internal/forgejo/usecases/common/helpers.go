package common

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	forgejo "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

const TimeFormatDate = time.RFC3339

// HTTPClient は HTTP アクセスの抽象です。
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DecodeProjects は複数の API 応答形式を吸収して projects を返します。
func DecodeProjects(data []byte) ([]ProjectResponse, error) {
	var projects []ProjectResponse
	if err := json.Unmarshal(data, &projects); err == nil {
		return projects, nil
	}

	var wrapper struct {
		Data     []ProjectResponse `json:"data"`
		Projects []ProjectResponse `json:"projects"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	if len(wrapper.Data) > 0 {
		return wrapper.Data, nil
	}
	return wrapper.Projects, nil
}

// SetAuthHeader は API 呼び出し用ヘッダを設定します。
func SetAuthHeader(req *http.Request, token string) {
	if req == nil {
		return
	}
	if strings.TrimSpace(token) != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	req.Header.Set("Accept", "application/json")
}

// RepositoryName はリポジトリから owner/name を抽出します。
func RepositoryName(repo *forgejo.Repository) (owner string, name string) {
	name = strings.TrimSpace(repo.Name)
	if repo.Owner != nil {
		owner = strings.TrimSpace(repo.Owner.UserName)
	}
	if owner == "" {
		parts := strings.SplitN(repo.FullName, "/", 2)
		if len(parts) == 2 {
			owner = strings.TrimSpace(parts[0])
		}
	}
	return owner, name
}

// FormatDate は RFC3339 文字列へ変換します。
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(TimeFormatDate)
}

// PrimaryLanguage は最大値の言語名を返します。
func PrimaryLanguage(langs map[string]float64) string {
	if len(langs) == 0 {
		return ""
	}
	type item struct {
		name  string
		value float64
	}
	items := make([]item, 0, len(langs))
	for name, value := range langs {
		items = append(items, item{name: name, value: value})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].value == items[j].value {
			return items[i].name < items[j].name
		}
		return items[i].value > items[j].value
	})
	return items[0].name
}

// FirstNonEmptyString は最初の非空文字列を返します。
func FirstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
