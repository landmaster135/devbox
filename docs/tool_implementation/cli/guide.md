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
	fmt.Print(result) // 重要: 結果を標準出力に表示
}
```

## 実装時のチェックリスト

- [ ] `fmt.Print(result)` で結果を標準出力に表示している
- [ ] エラー時は `fmt.Fprintf(os.Stderr, ...)` でエラーメッセージを出力
- [ ] エラー時は `os.Exit(1)` でプロセスを終了
- [ ] ヘルプ機能（`-help` フラグ）を実装
- [ ] 操作タイプの switch 文で未対応操作をハンドリング

## よくある実装ミス

### 1. CLIツールで結果を表示しない

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

### 2. フラグパーサーのモック実装の間違い（正しい実装は `internal/zip_compressor/config/config_test.go` を参照）

```go
// ❌ 間違い: フラグ定義後に値を設定しても反映されない
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
  *p = value // デフォルト値のみ
}

// ✅ 正しい: 事前設定値をチェックして適用
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
  if presetValue, exists := m.stringValues[name]; exists {
    *p = presetValue // 事前設定値を優先
  } else {
    *p = value // デフォルト値
  }
  m.stringVars[name] = p
}
```
