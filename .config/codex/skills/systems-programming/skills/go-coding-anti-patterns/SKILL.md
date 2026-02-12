---
name: go-coding-anti-patterns
description: Use this skill when users ask to detect Go coding anti-patterns, perform Go code reviews, or refactor Go code for idiomatic and safe patterns.
---

# Go Coding Anti-patterns Detector

## Purpose
Go code review and refactor support focused on anti-pattern detection and concrete remediation.

## When to use
- User asks for Go anti-pattern detection.
- User asks for Go code review quality checks.
- User asks to make Go code more idiomatic or safer.

## Detection workflow
1. Identify review target files and scope.
2. Scan candidate patterns with fast text search (`rg`) and then verify context line by line.
3. Match findings against the anti-pattern catalog below.
4. Report findings with file path + line number, impact, and minimal fix example.
5. If editing is requested, apply fixes and run relevant tests.

## Reporting format
- `Severity` (`high`/`medium`/`low`)
- `Category`
- `File` and `line`
- `Current pattern`
- `Why it is a problem`
- `Recommended fix`

## Golang Coding Anti-patterns Collection

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
