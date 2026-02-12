---
name: golang-senior
description: Skill for Go implementation, refactoring, and test design. Use when applying the `## Go` rules in `instructions.md`, including struct-method-first design, naming conventions, CLI verification, and coverage improvement workflows.
---

# Golang Senior

Follow these rules for Go coding.

## Points to keep in mind during development

- When creating functions, basically make them as methods of a struct. When creating test functions, it's fine to make them as standalone functions.
- Always implement functions, structs, and objects starting from the base level. That is, the next thing you implement should be something that will be called by what you previously implemented.
- After finishing the CLI tool implementation, update the README within the CLI package before testing the functionality.

## Points to keep in mind during implementing tests and verification

- When creating test functions, include the name of the struct after the prefix `Test`. And, add the suffix `_Normal` for test case names that test the normal path.
- Once the implemented test code passes normally, check the coverage and report it to the user. Then, to improve coverage, run the `go tool cover -html=coverage.out -o coverage.html` command. From the results, add test cases to cover the parts of the functionality you were implementing that aren't covered by tests yet.
- When testing modules, execute with `go test -coverprofile=coverage.out ./...` in the local environment.
- Always run tests for the entire package when executing tests.
- Use the `go run` command for verification when testing CLI tool functionality.
- If, when you run `go test`, you encounter an error like:

```text
failed in sandbox LinuxSeccomp with execution error: sandbox denied exec error, exit code: 1, stdout: , stderr: open /home/user/.cache/go-build/dc/dcaa44daf65c89ad04ba373e820b4ec33f282d5855df0cfc997a4e3d8a29bdf9-d: permission denied
```

stop the work and explain the reason for the interruption.

## Naming Conventions

Use PascalCase for:
- Public structs
- Public interfaces
- Public methods
- Public functions
- Public variables
- Public constants
- Public fields
- Any identifier that needs to be accessible outside the package

Use camelCase for:
- Private methods
- Private functions
- Private variables
- Private constants
- Private fields
- Local variables

Use short names for:
- Loop counters (`i`, `j`, `k`, etc.)
- Temporary variables (with short scope)
- Method receivers (typically 1-2 characters)
- Error variables (`err`)

Use lowercase for:
- Package names (a single lowercase word, or kebab-case if it doesn't fit)
- Directory names (a single lowercase word, or kebab-case if it doesn't fit)
- File names (lowercase with underscores)

Special cases:
- Acronyms (`URL`, `HTTP`): all uppercase when public, all lowercase when private
- Interfaces: single method interfaces are named with the method name + `er`
- Test file names in Go have the suffix `_test`

## How to TDD with mock

1. Code processes with interfaces for external interactions such as HTTP request, OS file, and runtime calls.
2. Define mock structs and methods for the interfaces.

Example (`util.go`):

```go
package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GitHubClient struct {
	httpClient HTTPClient
	token      string
}

func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{},
		token:      token,
	}
}

func (c *GitHubClient) doRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	if method == "POST" || method == "PATCH" || method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var ghError GitHubError
		if err := json.Unmarshal(respBody, &ghError); err != nil {
			return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(respBody))
		}
		ghError.StatusCode = resp.StatusCode
		return nil, &ghError
	}

	return respBody, nil
}
```

Example (`util_test.go`):

```go
package github

import "net/http"

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
```
