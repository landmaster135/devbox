package usecases

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

type authStubHTTPClient struct {
	handler func(req *http.Request) (*http.Response, error)
}

func (c *authStubHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.handler(req)
}

func TestBuildAuthorizationURL(t *testing.T) {
	svc := NewAuthService(5 * time.Second)
	urlStr, err := svc.BuildAuthorizationURL(AuthorizationURLParams{
		ClientID:    "client",
		RedirectURI: "https://example.com/callback",
		Scope:       " user.metrics , user.activity",
		State:       "state123",
		Mode:        "demo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parsed, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("failed to parse generated url: %v", err)
	}
	q := parsed.Query()
	if q.Get("scope") != "user.metrics,user.activity" {
		t.Fatalf("scope normalization failed: %s", q.Get("scope"))
	}
	if q.Get("state") != "state123" || q.Get("mode") != "demo" {
		t.Fatalf("unexpected query params: %v", q)
	}
}

func TestExchangeAuthorizationCode(t *testing.T) {
	stub := &authStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(req.Body)
		form, _ := url.ParseQuery(string(body))
		if form.Get("grant_type") != "authorization_code" {
			t.Fatalf("unexpected grant_type: %s", form.Get("grant_type"))
		}
		respBody := `{"status":0,"body":{"userid":123,"access_token":"token","refresh_token":"refresh","expires_in":3600,"token_type":"Bearer","scope":"user.metrics"}}`
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(respBody)), Header: make(http.Header)}, nil
	}}

	svc := NewAuthServiceWithClient(stub, 10*time.Second)
	resp, err := svc.ExchangeAuthorizationCode(context.Background(), TokenRequest{
		ClientID:     "client",
		ClientSecret: "secret",
		Code:         "code",
		RedirectURI:  "https://example.com/callback",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Body.UserID != "123" || resp.Body.AccessToken != "token" {
		t.Fatalf("unexpected token response: %+v", resp.Body)
	}
}

func TestRefreshAccessToken(t *testing.T) {
	stub := &authStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(req.Body)
		form, _ := url.ParseQuery(string(body))
		if form.Get("grant_type") != "refresh_token" {
			t.Fatalf("unexpected grant_type: %s", form.Get("grant_type"))
		}
		respBody := `{"status":0,"body":{"userid":"321","access_token":"token2","refresh_token":"refresh2","expires_in":1800,"token_type":"Bearer","scope":"user.metrics,user.activity"}}`
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(respBody)), Header: make(http.Header)}, nil
	}}

	svc := NewAuthServiceWithClient(stub, 0)
	resp, err := svc.RefreshAccessToken(context.Background(), RefreshRequest{
		ClientID:     "client",
		ClientSecret: "secret",
		RefreshToken: "refresh",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Body.UserID != "321" || resp.Body.RefreshToken != "refresh2" {
		t.Fatalf("unexpected response: %+v", resp.Body)
	}
	if svc.requestTimeout != defaultOAuthTimeout {
		t.Fatalf("expected default timeout to be applied, got %v", svc.requestTimeout)
	}
}

func TestSetEndpoints(t *testing.T) {
	svc := NewAuthService(5 * time.Second)
	svc.SetEndpoints("https://auth.example", "https://token.example")
	if svc.authorizeURL != "https://auth.example" || svc.tokenEndpoint != "https://token.example" {
		t.Fatalf("SetEndpoints failed: authorize=%s token=%s", svc.authorizeURL, svc.tokenEndpoint)
	}
}

func TestNewAuthServiceWithEndpoints(t *testing.T) {
	svc := NewAuthServiceWithEndpoints(3*time.Second, " https://auth.example/ ", " https://token.example ")
	if svc.authorizeURL != "https://auth.example/" || svc.tokenEndpoint != "https://token.example" {
		t.Fatalf("endpoints not normalized: %s %s", svc.authorizeURL, svc.tokenEndpoint)
	}
	if svc.requestTimeout != 3*time.Second {
		t.Fatalf("timeout not applied: %v", svc.requestTimeout)
	}
}

func TestExchangeAuthorizationCodeValidation(t *testing.T) {
	svc := NewAuthService(5 * time.Second)
	_, err := svc.ExchangeAuthorizationCode(context.Background(), TokenRequest{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestRefreshAccessTokenValidation(t *testing.T) {
	svc := NewAuthService(5 * time.Second)
	_, err := svc.RefreshAccessToken(context.Background(), RefreshRequest{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestParseFlexibleString(t *testing.T) {
	value, err := parseFlexibleString([]byte("123"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "123" {
		t.Fatalf("unexpected value: %s", value)
	}
	value, err = parseFlexibleString([]byte("\"456\""))
	if err != nil || value != "456" {
		t.Fatalf("unexpected string value: %s, err: %v", value, err)
	}
}

func TestDoTokenRequestHTTPErrors(t *testing.T) {
	svc := NewAuthService(5 * time.Second)
	svc.httpClient = &authStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(strings.NewReader("oops")), Header: make(http.Header)}, nil
	}}

	_, err := svc.doTokenRequest(context.Background(), url.Values{})
	if err == nil {
		t.Fatal("expected HTTP error")
	}

	svc.httpClient = &authStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		respBody := `{"status":247,"error":"invalid"}`
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(respBody)), Header: make(http.Header)}, nil
	}}

	resp, err := svc.doTokenRequest(context.Background(), url.Values{"action": {"requesttoken"}})
	if err == nil || resp.Status != 247 {
		t.Fatalf("expected status error, got %#v err=%v", resp, err)
	}
}
