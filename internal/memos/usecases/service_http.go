package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (s *Service) doJSON(ctx context.Context, method, requestPath string, query url.Values, payload any, out any) error {
	req, err := s.newRequest(ctx, method, requestPath, query, payload)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("Memos API %s %s の呼び出しに失敗しました: %w", method, requestPath, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Memos API %s %s が失敗しました: status=%d body=%s", method, requestPath, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if out == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("Memos API %s %s のレスポンスデコードに失敗しました: %w", method, requestPath, err)
	}
	return nil
}

func (s *Service) newRequest(ctx context.Context, method, requestPath string, query url.Values, payload any) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return nil, fmt.Errorf("リクエストボディのエンコードに失敗しました: %w", err)
		}
		body = buf
	}

	requestURL := s.baseURL + requestPath
	if len(query) > 0 {
		requestURL = requestURL + "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("HTTP リクエストの作成に失敗しました: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if s.apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiToken)
	}
	return req, nil
}
