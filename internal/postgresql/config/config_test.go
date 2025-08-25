package config

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

func TestParseFlags_Normal(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--table-name=users",
		"--format=json",
		"--output-path=/tmp",
		"--limit=100",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() でエラーが発生しました: %v", err)
	}

	// 期待値の検証
	if cfg.Operation != "dump" {
		t.Errorf("Operation = %s, want dump", cfg.Operation)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost/testdb" {
		t.Errorf("DatabaseURL = %s, want postgres://user:pass@localhost/testdb", cfg.DatabaseURL)
	}
	if cfg.TableName != "users" {
		t.Errorf("TableName = %s, want users", cfg.TableName)
	}
	if cfg.Format != "json" {
		t.Errorf("Format = %s, want json", cfg.Format)
	}
	if cfg.OutputPath != "/tmp" {
		t.Errorf("OutputPath = %s, want /tmp", cfg.OutputPath)
	}
	if cfg.Limit == nil || *cfg.Limit != 100 {
		t.Errorf("Limit = %v, want 100", cfg.Limit)
	}
	if cfg.Help {
		t.Errorf("Help = %t, want false", cfg.Help)
	}
}

func TestParseFlags_MinimalArgs_Normal(t *testing.T) {
	// テスト用のコマンドライン引数を設定（最小限）
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--table-name=users",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() でエラーが発生しました: %v", err)
	}

	// 期待値の検証
	if cfg.Operation != "dump" {
		t.Errorf("Operation = %s, want dump", cfg.Operation)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost/testdb" {
		t.Errorf("DatabaseURL = %s, want postgres://user:pass@localhost/testdb", cfg.DatabaseURL)
	}
	if cfg.TableName != "users" {
		t.Errorf("TableName = %s, want users", cfg.TableName)
	}
	if cfg.Format != "json" {
		t.Errorf("Format = %s, want json (default)", cfg.Format)
	}
	if cfg.OutputPath != "" {
		t.Errorf("OutputPath = %s, want empty string (default)", cfg.OutputPath)
	}
	if cfg.Limit != nil {
		t.Errorf("Limit = %v, want nil (default)", cfg.Limit)
	}
}

func TestParseFlags_Help_Normal(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--help",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() でエラーが発生しました: %v", err)
	}

	// ヘルプフラグが設定されていることを確認
	if !cfg.Help {
		t.Errorf("Help = %t, want true", cfg.Help)
	}
}

func TestParseFlags_MissingOperation_Error(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--table-name=users",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() でエラーが期待されましたが、エラーが発生しませんでした")
	}

	expectedError := "--operation は必須です"
	if err.Error() != expectedError {
		t.Errorf("エラーメッセージ = %s, want %s", err.Error(), expectedError)
	}
}

func TestParseFlags_MissingDatabaseURL_Error(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump",
		"--table-name=users",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() でエラーが期待されましたが、エラーが発生しませんでした")
	}

	expectedError := "--database-url は必須です"
	if err.Error() != expectedError {
		t.Errorf("エラーメッセージ = %s, want %s", err.Error(), expectedError)
	}
}

func TestParseFlags_MissingTableName_Error(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump",
		"--database-url=postgres://user:pass@localhost/testdb",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() でエラーが期待されましたが、エラーが発生しませんでした")
	}

	expectedError := "--table-name は必須です (dump操作時)"
	if err.Error() != expectedError {
		t.Errorf("エラーメッセージ = %s, want %s", err.Error(), expectedError)
	}
}

func TestParseFlags_InvalidOperation_Error(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=invalid",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--table-name=users",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() でエラーが期待されましたが、エラーが発生しませんでした")
	}

	expectedError := "未対応の操作です: invalid (対応操作: dump, dump-all-tables, list-tables-minimum, list-tables)"
	if err.Error() != expectedError {
		t.Errorf("エラーメッセージ = %s, want %s", err.Error(), expectedError)
	}
}

func TestParseFlags_InvalidFormat_Error(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--table-name=users",
		"--format=invalid",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() でエラーが期待されましたが、エラーが発生しませんでした")
	}

	expectedError := "未対応のフォーマットです: invalid (対応フォーマット: json, csv, sql, text)"
	if err.Error() != expectedError {
		t.Errorf("エラーメッセージ = %s, want %s", err.Error(), expectedError)
	}
}

func TestParseFlags_InvalidLimit_Error(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--table-name=users",
		"--limit=-1",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() でエラーが期待されましたが、エラーが発生しませんでした")
	}

	expectedError := "--limit は正の数もしくは0である必要があります: -1"
	if err.Error() != expectedError {
		t.Errorf("エラーメッセージ = %s, want %s", err.Error(), expectedError)
	}
}

func TestParseFlags_ZeroLimit_Normal(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--table-name=users",
		"--limit=0",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() でエラーが発生しました: %v", err)
	}

	// limit=0の場合はnilになることを確認（制限なしとして扱われる）
	if cfg.Limit != nil {
		t.Errorf("Limit = %v, want nil (limit=0は制限なしとして扱われる)", cfg.Limit)
	}
}

func TestParseFlags_ListTablesMinimum_Normal(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=list-tables-minimum",
		"--database-url=postgres://user:pass@localhost/testdb",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() でエラーが発生しました: %v", err)
	}

	// 期待値の検証
	if cfg.Operation != "list-tables-minimum" {
		t.Errorf("Operation = %s, want list-tables-minimum", cfg.Operation)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost/testdb" {
		t.Errorf("DatabaseURL = %s, want postgres://user:pass@localhost/testdb", cfg.DatabaseURL)
	}
	if cfg.TableName != "" {
		t.Errorf("TableName = %s, want empty (list-tables-minimum操作では不要)", cfg.TableName)
	}
	if cfg.Format != "json" {
		t.Errorf("Format = %s, want json (default)", cfg.Format)
	}
}

func TestParseFlags_ListTables_Normal(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=list-tables",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--format=text",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() でエラーが発生しました: %v", err)
	}

	// 期待値の検証
	if cfg.Operation != "list-tables" {
		t.Errorf("Operation = %s, want list-tables", cfg.Operation)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost/testdb" {
		t.Errorf("DatabaseURL = %s, want postgres://user:pass@localhost/testdb", cfg.DatabaseURL)
	}
	if cfg.TableName != "" {
		t.Errorf("TableName = %s, want empty (list-tables操作では不要)", cfg.TableName)
	}
	if cfg.Format != "text" {
		t.Errorf("Format = %s, want text", cfg.Format)
	}
}

func TestParseFlags_ListTablesMinimumInvalidFormat_Error(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=list-tables-minimum",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--format=text",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() でエラーが期待されましたが、エラーが発生しませんでした")
	}

	expectedError := "list-tables-minimum操作ではjsonフォーマットのみ対応しています"
	if err.Error() != expectedError {
		t.Errorf("エラーメッセージ = %s, want %s", err.Error(), expectedError)
	}
}

func TestParseFlags_DumpAllTables_Normal(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump-all-tables",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--format=csv",
		"--output-path=/tmp/dumps",
		"--limit=500",
		"--concurrency=3",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() でエラーが発生しました: %v", err)
	}

	// 期待値の検証
	if cfg.Operation != "dump-all-tables" {
		t.Errorf("Operation = %s, want dump-all-tables", cfg.Operation)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost/testdb" {
		t.Errorf("DatabaseURL = %s, want postgres://user:pass@localhost/testdb", cfg.DatabaseURL)
	}
	if cfg.TableName != "" {
		t.Errorf("TableName = %s, want empty (dump-all-tables操作では不要)", cfg.TableName)
	}
	if cfg.Format != "csv" {
		t.Errorf("Format = %s, want csv", cfg.Format)
	}
	if cfg.OutputPath != "/tmp/dumps" {
		t.Errorf("OutputPath = %s, want /tmp/dumps", cfg.OutputPath)
	}
	if cfg.Limit == nil || *cfg.Limit != 500 {
		t.Errorf("Limit = %v, want 500", cfg.Limit)
	}
	if cfg.Concurrency == nil || *cfg.Concurrency != 3 {
		t.Errorf("Concurrency = %v, want 3", cfg.Concurrency)
	}
}

func TestParseFlags_Concurrency_Normal(t *testing.T) {
	// テスト用のコマンドライン引数を設定
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump-all-tables",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--concurrency=5",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() でエラーが発生しました: %v", err)
	}

	// concurrencyの検証
	if cfg.Concurrency == nil || *cfg.Concurrency != 5 {
		t.Errorf("Concurrency = %v, want 5", cfg.Concurrency)
	}
}

func TestParseFlags_ConcurrencyInvalid_Error(t *testing.T) {
	// テスト用のコマンドライン引数を設定（範囲外の値）
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump-all-tables",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--concurrency=15",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() でエラーが期待されましたが、エラーが発生しませんでした")
	}

	expectedError := "--concurrency は1以上かつ10以下である必要があります: 15"
	if err.Error() != expectedError {
		t.Errorf("エラーメッセージ = %s, want %s", err.Error(), expectedError)
	}
}

func TestParseFlags_ConcurrencyZero_Error(t *testing.T) {
	// テスト用のコマンドライン引数を設定（0の場合）
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump-all-tables",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--concurrency=0",
	}

	// flagパッケージをリセット
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() でエラーが期待されましたが、エラーが発生しませんでした")
	}

	expectedError := "--concurrency は1以上かつ10以下である必要があります: 0"
	if err.Error() != expectedError {
		t.Errorf("エラーメッセージ = %s, want %s", err.Error(), expectedError)
	}
}

// 境界値テスト: concurrencyの最大値
func TestParseFlags_ConcurrencyMax_Normal(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump-all-tables",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--concurrency=10",
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() でエラーが発生しました: %v", err)
	}

	if cfg.Concurrency == nil || *cfg.Concurrency != 10 {
		t.Errorf("Concurrency = %v, want 10", cfg.Concurrency)
	}
}

// 全フォーマットでの動作確認テスト
func TestParseFlags_AllFormats_Normal(t *testing.T) {
	formats := []string{"json", "csv", "sql", "text"}

	for _, format := range formats {
		t.Run("format_"+format, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = []string{
				"postgresql-cli",
				"--operation=dump",
				"--database-url=postgres://user:pass@localhost/testdb",
				"--table-name=users",
				"--format=" + format,
			}

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			cfg, err := ParseFlags()
			if err != nil {
				t.Fatalf("ParseFlags() でエラーが発生しました (format=%s): %v", format, err)
			}

			if cfg.Format != format {
				t.Errorf("Format = %s, want %s", cfg.Format, format)
			}
		})
	}
}

// デフォルト値確認テスト
func TestParseFlags_DefaultConcurrency_Normal(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump-all-tables",
		"--database-url=postgres://user:pass@localhost/testdb",
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() でエラーが発生しました: %v", err)
	}

	// concurrencyのデフォルト値が設定されていることを確認
	if cfg.Concurrency == nil {
		t.Error("Concurrency = nil, want non-nil (デフォルト値)")
	} else if *cfg.Concurrency < 1 || *cfg.Concurrency > 10 {
		t.Errorf("Concurrency = %d, want 1-10 range (デフォルト値)", *cfg.Concurrency)
	}
}

// 各操作でのテーブル名要否テスト
func TestParseFlags_TableNameRequirements_Normal(t *testing.T) {
	testCases := []struct {
		operation        string
		requireTableName bool
		description      string
	}{
		{"dump", true, "dump操作ではtable-nameが必須"},
		{"dump-all-tables", false, "dump-all-tables操作ではtable-nameは不要"},
		{"list-tables-minimum", false, "list-tables-minimum操作ではtable-nameは不要"},
		{"list-tables", false, "list-tables操作ではtable-nameは不要"},
	}

	for _, tc := range testCases {
		t.Run("operation_"+tc.operation, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			args := []string{
				"postgresql-cli",
				"--operation=" + tc.operation,
				"--database-url=postgres://user:pass@localhost/testdb",
			}

			// dump操作の場合のみtable-nameを追加
			if tc.requireTableName {
				args = append(args, "--table-name=users")
			}

			os.Args = args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			cfg, err := ParseFlags()
			if err != nil {
				t.Fatalf("ParseFlags() でエラーが発生しました (%s): %v", tc.description, err)
			}

			if cfg.Operation != tc.operation {
				t.Errorf("Operation = %s, want %s", cfg.Operation, tc.operation)
			}
		})
	}
}

// CSVフォーマットでの各操作テスト
func TestParseFlags_CSVFormat_Normal(t *testing.T) {
	operations := []string{"dump", "dump-all-tables", "list-tables"}

	for _, operation := range operations {
		t.Run("csv_"+operation, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			args := []string{
				"postgresql-cli",
				"--operation=" + operation,
				"--database-url=postgres://user:pass@localhost/testdb",
				"--format=csv",
			}

			// dump操作の場合のみtable-nameを追加
			if operation == "dump" {
				args = append(args, "--table-name=users")
			}

			os.Args = args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			cfg, err := ParseFlags()
			if err != nil {
				t.Fatalf("ParseFlags() でエラーが発生しました (operation=%s, format=csv): %v", operation, err)
			}

			if cfg.Format != "csv" {
				t.Errorf("Format = %s, want csv", cfg.Format)
			}
			if cfg.Operation != operation {
				t.Errorf("Operation = %s, want %s", cfg.Operation, operation)
			}
		})
	}
}

// 大きなlimit値でのテスト
func TestParseFlags_LargeLimit_Normal(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"postgresql-cli",
		"--operation=dump",
		"--database-url=postgres://user:pass@localhost/testdb",
		"--table-name=users",
		"--limit=1000000",
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() でエラーが発生しました: %v", err)
	}

	if cfg.Limit == nil || *cfg.Limit != 1000000 {
		t.Errorf("Limit = %v, want 1000000", cfg.Limit)
	}
}

// PrintUsage関数のテスト
func TestPrintUsage_Normal(t *testing.T) {
	// 標準エラー出力をキャプチャ
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// PrintUsage関数を実行
	PrintUsage()

	// パイプを閉じて出力を読み取り
	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// 期待される内容が含まれているかを確認
	expectedContents := []string{
		"PostgreSQL CLIツール",
		"使用方法:",
		"postgresql-cli [オプション]",
		"必須オプション:",
		"--operation string",
		"--database-url string",
		"--table-name string",
		"オプション:",
		"--output-path string",
		"--format string",
		"--limit int",
		"--concurrency int",
		"--help",
		"使用例:",
		"単一テーブルダンプ",
		"全テーブルダンプ",
		"テーブル一覧",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("PrintUsage()の出力に期待される内容が含まれていません: %s", expected)
		}
	}

	// 出力が空でないことを確認
	if len(output) == 0 {
		t.Error("PrintUsage()の出力が空です")
	}
}
