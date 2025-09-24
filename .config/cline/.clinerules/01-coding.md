# What is important for coding

## Common things
- When encountering modules or packages with unclear specifications, use the MCP tool to search and investigate their specifications.
- When adding new functions or methods, write them after the final line of similar function groups within the file.
- Create functions and methods in small units whenever possible.
- Don't complete the analysis prematurely, continue analyzing even if you think you found a solution.
- Keep each file's line count under 700 lines. Before adding test code, count the number of lines in the file you're planning to add to using the shell command: "awk 'END { print NR }' ". Then, determine whether that number exceeds 700 using the MCP tool. If it's less than 700, append to the existing test file. If it's more than 700, create a new file in the same directory and append to that file.
- Avoid using hard coding as much as possible.
- Don’t forget to add logging.

## Test Driven Development (TDD)
- When implementing features requested by users, always implement test code for those features as well.
- After implementing features requested by users, confirm with the user before implementing the test code for those features.
- When implementing test code, keep the Red-Green-Refactor cycle in mind.
- Test methods are implemented based on the Arrange-Act-Assert pattern.
- Once the implemented test code passes normally, check the coverage and report it to the user. Then, improve the coverage.

## SOLID principles
Be mindful of SOLID principles. SOLID acronym stands for five design principles that help make software more maintainable and scalable:
- Single Responsibility Principle
- Open/Closed Principle
- Liskov Substitution Principle
- Interface Segregation Principle
- Dependency Inversion Principle

## Extra things

### Pull request
- Before creating a pull request, confirm Git commit history with the MCP tool: `get_git_commit_history_with_details`. I need to ask the user for the first date, last date, and keywords (mainly these tags: "feat:", "refactor:", "fix:", "test:", or "doc:") of the commits necessary to create a pull request.

## Python
- When creating functions, basically make them as instance methods of a class. When creating test functions, also create a test class and make them as instance methods.
- When testing modules, execute with `python -m pytest --cov=src --cov-branch --tb=short -vv`.

### Naming Conventions

Use PascalCase for the following:
- Classes
- Exceptions
- Test classes

Use lower_snake_case for the following:
- Directory names
- File names
- Modules
- Methods
- Functions
- Variables

Begin the words `test_` to use lower_snake_case for the following:
- Test methods
- Test functions

Use UPPER_SNAKE_CASE for the following:
- Environment variables
- Constants
- Global configurations


Components
Type definitions
Interfaces


Use kebab-case for the following:

Directory names (e.g., components/auth-wizard)
File names (e.g., user-profile.tsx)


Use camelCase for the following:

Variables
Functions
Methods
Hooks
Properties
Props


Use uppercase for the following:

Environment variables
Constants
Global configurations


For boolean variables, use a verb as prefix: isLoading, hasError, canSubmit

## Go

### Points to keep in mind during development
- When creating functions, basically make them as methods of a struct. When creating test functions, it's fine to make them as standalone functions.
- Always implement functions, structs, and objects starting from the base level. That is, the next thing you implement should be something that will be called by what you previously implemented.
- After finishing the CLI tool implementation, update the README within the CLI package before testing the functionality.

### Points to keep in mind during implementing tests and verification
- When creating test functions, include the name of the struct after the prefix 'Test'. And, add the suffix '_Normal' for test case names that test the normal path.
- Once the implemented test code passes normally, check the coverage and report it to the user. Then, to improve coverage, run the `go tool cover -html=coverage.out -o coverage.html` command. From the results, add test cases to cover the parts of the functionality you were implementing that aren't covered by tests yet.
- When testing modules, execute with `go test -coverprofile=coverage.out ./...`.
- Always run tests for the entire package when executing tests.
- Use the `go run` command for verification when testing CLI tool functionality.

### Naming Conventions
Use PascalCase for the following:
- Public structs
- Public interfaces
- Public methods
- Public functions
- Public variables
- Public constants
- Public fields
- Any identifier that needs to be accessible outside the package

Use camelCase for the following:
- Private methods
- Private functions
- Private variables
- Private constants
- Private fields
- Local variables

Use short names for the following:
- Loop counters (i, j, k, etc.)
- Temporary variables (with short scope)
- Method receivers (typically 1-2 characters)
- Error variables (err)

Use lowercase for the following:
- Package names (a single lowercase word, or kebab-case if it doesn't fit)
- Directory names (a single lowercase word, or kebab-case if it doesn't fit)
- File names (lowercase with underscores)

Special cases:
- Acronyms (URL, HTTP): all uppercase when public, all lowercase when private
- Interfaces: single method interfaces are named with the method name + "er"
- Don't forget that test file names in Go have the suffix '_test'.

### How to TDD with mock
1. Code processes with interface for the following process: HTTP request, OS file, runtime calling and etc...
e.g. `util.go` for HTTP request:
```go
package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPClient インターフェースを定義
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GitHubClient 構造体を修正（テスト用）
type GitHubClient struct {
	httpClient HTTPClient
	token      string
}

// NewGitHubClient は新しいGitHubクライアントを作成します
func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{},
		token:      token,
	}
}

// doRequest はHTTPリクエストを実行し、レスポンスを処理します（テスト用に簡略化）
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

2. Definite struct and method for mock.
e.g. `util_test.go`:
```go
package github

import (
  "net/http"
)

// 'MockHTTPClient' struct is a mock for HTTP client.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// 'Do' method executes HTTP request.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
```

### Golang Coding Anti-patterns Collection

| Category | Anti-pattern | Problem | Improvement |
|---------|---------------|---------|-------------|
| **nil Check** | `if slice != nil && len(slice) > 0` | len() for nil slices is defined as zero, making nil check unnecessary | `if len(slice) > 0` |
| **nil Check** | `if map != nil && len(map) > 0` | len() for nil maps is defined as zero, making nil check unnecessary | `if len(map) > 0` |
| **Variable Declaration** | `var s string = ""` | Explicit assignment is unnecessary when initializing with zero value | `var s string` |
| **Variable Declaration** | `var i int = 0` | Explicit assignment is unnecessary when initializing with zero value | `var i int` |
| **Loop** | `for i := 0; i < len(slice); i++ { v := slice[i] }` | Using range is more concise and safe | `for _, v := range slice` |
| **String Comparison** | `strings.ToLower(s1) == strings.ToLower(s2)` | Use dedicated function for case-insensitive comparison | `strings.EqualFold(s1, s2)` |
| **Error Handling** | Writing `if err != nil { return nil, err }` repeatedly | Not utilizing error wrapping or custom errors | `if err != nil { return nil, fmt.Errorf("operation failed: %w", err) }` |
| **Empty String Check** | `if len(s) == 0` | Direct comparison is recommended for string empty check | `if s == ""` |
| **Boolean Variable** | `if condition == true` | Explicit comparison with boolean values is unnecessary | `if condition` |
| **Boolean Variable** | `if condition == false` | Explicit comparison with boolean values is unnecessary | `if !condition` |
| **Loop Variable** | `for i, _ := range slice` | Omit unused variables | `for i := range slice` |
| **Slice Creation** | `slice := make([]int, 0, 0)` | Capacity of 0 can be omitted | `slice := make([]int, 0)` or `var slice []int` |
| **Type Conversion** | `int(float64(i))` | Unnecessary intermediate type conversion | Convert directly to required type |
| **Struct Initialization** | `MyStruct{Field1: value1, Field2: "", Field3: 0}` | Explicit setting of zero value fields is unnecessary | `MyStruct{Field1: value1}` |
| **defer Usage** | `f, err := os.Open(file); defer f.Close()` | Writing defer before error check is dangerous | `f, err := os.Open(file); if err != nil { return err }; defer f.Close()` |
| **interface{} Usage** | `func process(data interface{})` | Use generics after Go 1.18 | `func process[T any](data T)` |
| **goroutine** | `go func() { /* no error handling */ }()` | Improper error handling in goroutines | Use error channels or context |
| **String Concatenation** | Using `s += str` in loops | Inefficient | Use `strings.Builder` or `strings.Join` |
| **time Comparison** | `time1.Unix() == time2.Unix()` | Precision limited to seconds | `time1.Equal(time2)` |
| **Channel** | Calling `close(ch)` from multiple places | Causes panic | Use sync.Once or proper design patterns |
| **HTTP Response** | `resp, _ := http.Get(url)` | Ignoring errors is dangerous | Always perform error handling |
| **JSON Operations** | Reusing error variables with `json.Unmarshal(data, &v); if err != nil` | Variable scope issues | Declare error variables in appropriate scope |
| **File Operations** | `ioutil.ReadFile()` (after Go 1.16) | Using deprecated packages | `os.ReadFile()` |
| **mutex** | Passing `var mu sync.Mutex` by value in methods | mutex should be passed by reference | Use pointer or embed in struct |
| **context** | Always using `context.Background()` | Cannot propagate context properly | Use context received from upper layer |
| **context Arguments** | `func badFunc(k favContextKey, ctx context.Context)` | context.Context should be the first argument by convention | `func goodFunc(ctx context.Context, k favContextKey)` |
| **Concurrency** | Using mutex on channels | Channels are concurrency control mechanisms themselves | Learn proper channel usage patterns |
| **Export** | Exported functions returning unexported types | Difficult to use from outside the package | Export appropriate types or return interfaces |
| **goroutine Leak** | `go func() { /* no error handling */ }()` | goroutines may persist permanently | Proper termination conditions and error handling |
| **Channel Operations** | Using select for single channel operations | Unnecessary complexity | Use direct channel operations |
| **time.Timer** | Sharing time.Timer among multiple goroutines | Causes race conditions | Use individual Timer per goroutine |
| **Slice Joining** | Repeatedly using append in loops | Inefficient | Use `append(slice1, slice2...)` variadic version |
| **Return Statement** | `func foo() { ...; return }` | Unnecessary return statement in functions that don't return values | `func foo() { ... }` |
| **Switch Statement** | Assuming C-style fall-through | Go requires explicit fall-through specification | Understand that each case terminates automatically |
| **Error Ignoring** | `result, _ := someFunc()` | Ignoring errors is dangerous | Always perform error handling |
| **Constant Channel** | `ch := make(chan int, 0)` | Using magic numbers | Use named constants (except for debugging purposes) |
| **Type Assertion** | Type assertions that cause panic | Causes runtime errors | `value, ok := interface{}.(Type)` safe form |
| **Concurrent Access** | Shared variables with potential data races | Unexpected behavior or crashes | Use appropriate synchronization primitives |
| **dot import** | `import . "package"` | Reduces code readability | Use explicit package names |
| **init Function** | Flag initialization in init() | Makes testing difficult | Prefer initialization in main() |
| **context.Value** | Excessive use of context.Value | Lack of type safety, difficult testing | Prioritize explicit parameter passing |
| **Channel Size** | Deadlock with unbuffered channels | Synchronization issues between sender and receiver | Appropriate buffer size or asynchronous patterns |
| **sync Value Copy** | Passing mutex or WaitGroup by value | Synchronization doesn't work | Use with pointer or embedding |
| **Unnecessary Wrapper** | `func run(cmd string) error { return runRemote(cmd) }` | Unnecessary indirection layer | Call required function directly |
| **Generic Types** | Excessive use of `interface{}` (after Go 1.18) | Lack of type safety | Use appropriate generics |
| **Package Names** | Generic names like `util`, `tools`, `misc` | Package purpose is unclear | Use specific and descriptive names |
| **Pointer Abuse** | Using pointers to scalar values | Unnecessary complexity and increased memory usage | Use value passing as default |
| **Error Messages** | `assert.Equal(t, false, true)` | Meaningless test messages | Specific and understandable error messages |
| **Slice Search** | `for _, v := range slice { if v == target { return true } }` | Inefficient and verbose | `slices.Contains(slice, target)` (Go 1.21+) |
| **Slice Search** | Manual loop for searching in slices | Manual implementation of binary search | `slices.BinarySearch(sortedSlice, target)` (Go 1.21+) |
| **Slice Operations** | Manual slice concatenation and deletion | Error-prone and inefficient | Use `slices.Concat()`, `slices.Delete()` etc. (Go 1.21+) |
| **Slice Sorting** | Verbose use of `sort.Slice()` | Lack of type safety, performance issues | `slices.Sort()`, `slices.SortFunc()` (Go 1.21+) |
| **String Operations** | `if strings.HasPrefix(str, prefix) { str = str[len(prefix):] }` | Conditional prefix removal is verbose and error-prone | `str = strings.TrimPrefix(str, prefix)` |
| **Function Signature** | `func process() (error, string)` | Error should be the last return value by Go convention | `func process() (string, error)` |
| **Code Formatting** | `fmt.Printf("long format string with many args", arg1, arg2, arg3, arg4)` | Long function calls with multiple arguments on single line reduce readability | Break arguments into multiple lines with proper indentation |
| **Static Analysis Warning** | Continuing execution after nil check in tests | Checking for nil but not stopping execution can cause nil pointer dereference | Use `t.Fatal()` or `require.NotNil()` to stop execution after nil check |
| **Error Handling** | `fmt.Errorf("invalid input")` | fmt.Errorf is unnecessary when creating simple error messages without formatting | `errors.New("invalid input")` |
| **Testing Mock** | `mockObj.On("Method", mock.AnythingOfType("string")).Return(value, nil); mockObj.AssertExpectations(t)` | Using testify/mock framework creates complex setup, external dependency, and requires assertion calls with hard coding | Use function field-based mocks: `mockObj.MethodFunc = func(param string) (Type, error) { return value, nil }` |
