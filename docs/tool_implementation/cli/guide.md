# CLIツール実装ガイド

このドキュメントは、CLIツール実装に関する内容をまとめたものです。

## 実装パターン

### CLIツールのmain.go構造

```go
package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/{tool-name}/config"
	usecases "github.com/landmaster135/devbox/internal/{tool-name}/usecases"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	if cfg.Help {
		config.PrintUsage()
		return
	}

	switch cfg.Operation {
	case "operation1":
		handleOperation1(cfg)
	case "operation2":
		handleOperation2(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

func handleOperation1(cfg *config.Config) {
	service := usecases.NewService()
	result, err := service.HandleOperation1(cfg.Param1, cfg.Param2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(result)
}
```

### フラグ関連の責務分離

- `config` 層:
  - `Config` の生成と検証
  - `ParseFlags()` と `ParseFlagsWithParser(parser)` の公開
  - `PrintUsage()` で usage 出力処理を `infrastructures/flag_parser` に委譲
- `infrastructures/flag_parser` 層:
  - `FlagParser` interface
  - `StandardFlagParser` 実装
  - `PrintUsage(format)` 実装。`os.Stderr` と `os.Args[0]` を infrastructure 層で扱う
  - `MockFlagParser` 実装とその単体テスト。`infrastructures/flag_parser` に colocate して再利用可能にする。

想定ディレクトリ構成:

```text
internal/{tool-name}/
├─ config/
│  ├─ config.go
│  └─ config_test.go
└─ infrastructures/
   └─ flag_parser/
      ├─ flag_parser.go
      ├─ usage.go
      ├─ mock_flag_parser.go
      └─ *_test.go
```

### `PrintUsage(format)` の実装パターン

`PrintUsage(format)` は **出力先と実行ファイル名の注入のみ** を担当し、ヘルプ本文（usageTemplate）の組み立ては `config` 層で行います。

`infrastructures/flag_parser/usage.go`:

```go
package flag_parser

import (
	"fmt"
	"os"
)

func PrintUsage(format string) {
	fmt.Fprintf(os.Stderr, format, os.Args[0])
}
```

`config/config.go`:

```go
package config

import (
	"sort"
	"strconv"
	"strings"

	flagParser "github.com/landmaster135/devbox/internal/{tool-name}/infrastructures/flag_parser"
)

const usageTemplate = `使用方法: %[1]s [オプション]

共通オプション:
  -operation string 実行する操作 (%s)
  -timeout int     タイムアウト秒数 (デフォルト: %s)
`

func PrintUsage() {
	ops := make([]string, len(supportedOperations))
	copy(ops, supportedOperations)
	sort.Strings(ops)

	usage := strings.Replace(usageTemplate, "%s", strings.Join(ops, ", "), 1)
	usage = strings.Replace(usage, "%s", strconv.Itoa(defaultTimeoutSeconds), 1)
	flagParser.PrintUsage(usage)
}
```

## 実装アンチパターン

### CLIツールで結果を表示しない

```go
// ❌ 間違い: 結果を表示しない
service.HandleOperation(param)

// ✅ 正しい: 結果を標準出力に表示
result, err := service.HandleOperation(param)
if err != nil {
  fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
  os.Exit(1)
}
fmt.Print(result)
```

### フラグパーサーモックで事前設定値を反映しない（正しい実装は `internal/zip_compressor/infrastructures/flag_parser/mock_flag_parser.go` を参照）

```go
// ❌ 間違い: デフォルト値しか入らない
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
  *p = value
}

// ✅ 正しい: 事前設定値を優先して適用
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
  if presetValue, exists := m.stringValues[name]; exists {
    *p = presetValue
    return
  }
  if *p != "" {
    return
  }
  *p = value
}
```
