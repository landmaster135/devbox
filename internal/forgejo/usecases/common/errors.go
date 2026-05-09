package common

import (
	"fmt"
	"net/http"
)

// RequestError は API リクエスト失敗を表すエラーです。
type RequestError struct {
	Status int
	Body   string
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("request failed with status %d: %s", e.Status, e.Body)
}

// IsNotFoundError は 404 応答由来のエラーかを判定します。
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	requestErr, ok := err.(*RequestError)
	if !ok {
		return false
	}
	return requestErr.Status == http.StatusNotFound
}
