# 共通実装ガイド

## 実装戦略

- ユーザーから特別に指示が無い限り、後方互換性は考慮しない。
- `domain/`、`usecases/`、`infrastructures/` を実装。
- **infrastructure層は、他のいかなる層からも独立していなければならない。**

## 実装パターン

### 1. 共通サービス層の実装

```go
package usecases

type Service struct {
	// 依存関係の注入
}

func NewService() *Service {
	return &Service{}
}

// CLIツールとMCPツールで共通のビジネスロジック
func (s *Service) HandleOperation1(param1, param2 string) (string, error) {
	// ビジネスロジックの実装
	// 戻り値は文字列で統一（CLIツールで標準出力、MCPツールでクライアントに返却）
	
	result := fmt.Sprintf("処理結果: %s, %s", param1, param2)
	return result, nil
}
```

### 2. コマンド実行ラッパーによる非シェル化

外部入力を伴う処理で OS コマンドを実行する場合、`exec.Command` を直接 `name` と `args` に分割して呼び出すことで、シェル展開を回避しコマンドインジェクション対策を実装しなくても安全な構造を保てます。

```go
// usecases/command_executor.go
type CommandExecutor struct{}

func (e *CommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

func (e *CommandExecutor) ExecuteInDir(dir, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// 呼び出し側（例）
args := []string{"run", executionFile}
args = append(args, strings.Fields(parameters)...)
output, err := executor.ExecuteInDir(absRootDir, "go", args...)

// ポイント
// - シェル (`sh -c` / `bash -c`) を介さないので、`;` や `&&` を含む入力でも追加コマンド化されない
// - 入力は引数スライスとして渡され、Go ランタイムがエスケープを扱う
// - 出力は呼び出し元で検査・加工し、副作用が必要な場合のみ実施
```

このパターンを守れば、追加のエスケープ処理を実装せずにコマンドインジェクションを防げます。必要に応じて、許可するコマンド名や引数のホワイトリスト検証を組み合わせるとより堅牢になります。

### 3. 特定の処理を infrastructure 層へ集約し、usecase から呼び出す

infrastructure層で実装すべき主なリソース（シリアライズ境界は除く）:
- HTTP通信: `http.Client`、認証ヘッダ、タイムアウト、リトライ
- 外部APIアダプタ: Memos など外部サービス呼び出しの実装
- 永続化: DBアクセス、トランザクション、Repository実装
- メッセージ基盤: Queue/Kafka/PubSub への publish/consume
- ファイル/OS操作: ファイルI/O、パス処理、MIME判定、ディレクトリ走査
- 時刻/乱数/ID生成: `time.Now()`、UUID、乱数生成器
- 設定/シークレット取得: 環境変数、設定ファイル、Secret Manager
- 外部SDK連携: S3/GCS/GitHub/Notion などのSDKラッパ
- プロセス実行: `exec.Command` による外部コマンド実行

例えば、ファイル操作を含む操作を扱うときは、CLI で `os.ReadFile` せずに以下の責務分離を行います。

- CLI層: フラグ解析と `service` 呼び出しのみ（パス文字列をそのまま渡す）
- usecase層: ファイル読み込みを含む業務フローやビジネスロジック（例: content 解決、添付作成、既存添付マージ）
- infrastructure層: 実ファイルシステム実装（`os.ReadFile`、MIME 判定など）

実装例:
```go
// infrastructure
type FileSystem interface {
  ReadFile(filePath string) ([]byte, error)
  ReadAttachmentFile(filePath string) (*AttachmentFile, error)
}

type OSFileSystem struct{}

// usecase
type ServiceOptions struct {
  FileSystem infrastructures.FileSystem
}

func NewService(opts ServiceOptions) *Service {
  fs := opts.FileSystem
  if fs == nil {
    fs = infrastructures.NewOSFileSystem()
  }
  return &Service{fileSystem: fs}
}

func (s *Service) PatchFiles(ctx context.Context, memo string, filePaths []string, replaces bool) error {
  // fileSystem 経由で読み込み
  // CreateAttachment -> (必要なら ListMemoAttachments) -> SetMemoAttachments
  return nil
}
```

テスト方針:
- usecase テストでは `mockFileSystem` を注入し、ファイルI/Oに依存せず分岐を検証する
  - `infrastructures` 層の `Repository` に対するモックは、対応する `infrastructures/{resource}/` 直下で定義する
  - usecases 側のテストから上記モックを import して利用し、usecases 配下で同等モックを再定義しない
- CLI テストではファイル内容ではなく、`service` への引数委譲を検証する

### 4. operation別サブディレクトリ + common分離

operation が増えて `usecases/services.go` が肥大化する場合は、`usecases` 直下に処理をフラットに並べず、operationごとにサブディレクトリへ分割します。共通の業務ロジックは `common` サブディレクトリへ集約し、HTTP クライアントなどの外部I/Oは `infrastructures` 層へ分離します。

ディレクトリ構成の例:
```text
internal/{tool}/usecases/
├── services.go                      // Facadeと依存注入
├── service_operations.go            // 公開メソッドの委譲
├── service_contract.go              // 公開インターフェース + テスト用公開Mock
├── common/
│   ├── helpers.go                   // 正規化、共通バリデーション
│   ├── models.go                    // 共通DTO
│   └── requests.go                  // 共通request payload
├── operations/
│   ├── create_memo/
│   │   ├── service.go
│   │   └── service_test.go
│   ├── get_memo/
│   │   ├── service.go
│   │   └── service_test.go
│   └── ...
└── testutil/
    └── helpers.go                   // 共通テストヘルパ

internal/{tool}/infrastructures/
├── http/
│   ├── client.go                    // JSON HTTPクライアント実装
│   └── mock_client.go               // HTTPクライアント向けmock（usecaseテストから利用）
└── filesystem/
    ├── repository.go
    ├── impl.go
    └── mock_repository.go           // Repository向けmock（usecaseテストから利用）
```

実装ルール:
- `services.go` は Facade として薄く保つ（公開APIと依存注入に限定）
- 各 operation のロジックは `operations/{operation}/` に閉じ込める
- operation 間で共有する業務ロジック・DTO・request payload は `common/` に移す
- HTTPクライアントやファイル操作などの外部I/Oは `infrastructures/` に実装する
- `infrastructures` 層の `Repository` モックは `infrastructures/{resource}/` 直下に置く
- operation 専用テストは同じ operation ディレクトリに置く
- operation 実装ファイル名は `service.go`、テストファイル名は `service_test.go` とする

適用目安:
- operation が 4 種類以上あり、単一 `services.go` の見通しが悪くなっている場合
- operation ごとに依存（FileSystem/API client/補助関数）が異なり、責務分離が必要な場合

## 実装時の重要な注意点

### ツール種別ごとの注意点
- CLIツール固有の注意点: `docs/tool_implementation/cli/guide.md`
- MCPツール固有の注意点: `docs/tool_implementation/mcp/guide.md`

### ディレクトリ操作の注意点
- `cmd.Dir = dir` の設定は必須
- 相対パスは `filepath.Abs()` で絶対パスに変換
- ディレクトリの存在確認と種別確認を実装

実装例:
```go
// ✅ 正しい実装: ExecuteInDirメソッドでcmd.Dirを設定
func (e *DefaultCommandExecutor) ExecuteInDir(dir, name string, args ...string) ([]byte, error) {
  cmd := exec.Command(name, args...)
  cmd.Dir = dir  // 重要: 実行ディレクトリを設定
  return cmd.CombinedOutput()
}

// ✅ 正しい実装: 絶対パス変換とディレクトリ確認
func (s *Service) ExecuteOperation(directory string) (string, error) {
  // ディレクトリの存在確認
  if !s.directoryChecker.Exists(directory) {
    return "", fmt.Errorf("指定されたディレクトリが存在しません: %s", directory)
  }
  if !s.directoryChecker.IsDirectory(directory) {
    return "", fmt.Errorf("指定されたパスはディレクトリではありません: %s", directory)
  }

  // 絶対パスに変換（重要）
  absDir, err := filepath.Abs(directory)
  if err != nil {
    return "", fmt.Errorf("ディレクトリパスの変換に失敗しました: %v", err)
  }

  // 指定されたディレクトリでコマンド実行
  output, err := s.commandExecutor.ExecuteInDir(absDir, "go", "test", "./...")
  // ...
}
```

### exec.ExitErrorのExitCode()によるエラーの区別
- `err.(*exec.ExitError)` で型アサーション
- `exitError.ExitCode()` で終了コードを取得し、終了コードに応じて分岐
- **exit status 1**: テスト失敗 → 正常な動作として扱い、出力を返却
- **exit status 2以上**: システムエラー → エラーとして処理
- **その他のエラー**: コマンド未発見等 → エラーとして処理

例えば、Goのテストコマンドでは、テスト失敗とシステムエラーを以下のように区別します。

実装例:
```go
// ✅ 正しい実装: テスト失敗を正常な動作として扱う
func (s *GolangOpsService) ExecuteTestCoverage(directory, grepPattern string) (string, error) {
  // ... ディレクトリ確認等の処理 ...

  // go test -cover ./... を実行
  output, err := s.commandExecutor.ExecuteInDir(absDir, "go", "test", "-cover", "./...")
  if err != nil {
    if exitError, ok := err.(*exec.ExitError); ok {
      if exitError.ExitCode() == 1 {
        // テスト失敗は正常な動作として扱う
        // エラーを返さず、出力をそのまま処理する
      } else {
        // exit code が1以外の場合は実際のエラー
        return "", fmt.Errorf("テストカバレッジの実行でエラーが発生しました: %v\n出力: %s", err, string(output))
      }
    } else {
      // exec.ExitError以外のエラー
      return "", fmt.Errorf("コマンドの実行に失敗しました: %v", err)
    }
  }

  // テスト失敗の場合でも出力を処理して返却
  result.Write(output)
  return result.String(), nil
}
```

テスト実装での注意点:
```go
// テスト用の実際のExitErrorを作成
cmd := exec.Command("sh", "-c", "exit 1")
err := cmd.Run()
exitError, _ := err.(*exec.ExitError)

// モックでExitErrorを返却
mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-cover", "./..."}).Return(expectedOutput, exitError)
```

## 付録: 具体的な実装例の比較

### arithmetic-calculator の実装例
- **CLIツール** (`cmd/cli/arithmetic-calculator/main.go`)
- **MCPツール** (`cmd/mcp/arithmetic_calculator/mcp.go`)

### ops-for-golang の実装例
- **CLIツール** (`cmd/cli/ops-for-golang/main.go`)
- **MCPツール** (`cmd/mcp/ops_for_golang/mcp.go`)
