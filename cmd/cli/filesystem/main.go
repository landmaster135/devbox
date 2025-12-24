package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	usecases "github.com/landmaster135/devbox/internal/filesystem/usecases"
)

type cliConfig struct {
	operation      string
	path           string
	source         string
	destination    string
	content        string
	pattern        string
	excludePattern string
	allowedDirs    []string
}

func main() {
	flag.Usage = printUsage

	cfg := parseFlags()
	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *cliConfig {
	var cfg cliConfig
	var allowedDirsRaw string

	flag.StringVar(&cfg.operation, "operation", "", "実行する操作 (read_file, write_file, create_directory など)")
	flag.StringVar(&cfg.path, "path", "", "対象のファイルまたはディレクトリのパス")
	flag.StringVar(&cfg.source, "source", "", "移動元のパス (move_file操作用)")
	flag.StringVar(&cfg.destination, "destination", "", "移動先のパス (move_file操作用)")
	flag.StringVar(&cfg.content, "content", "", "ファイルに書き込む内容 (write_file操作用)")
	flag.StringVar(&cfg.pattern, "pattern", "", "検索に使用する文字列 (search_files操作用)")
	flag.StringVar(&cfg.excludePattern, "exclude-pattern", "", "除外するパターン (search_files操作用)")
	flag.StringVar(&allowedDirsRaw, "allowed-dirs", "", "カンマ区切りで追加の許可ディレクトリを指定")
	flag.Parse()

	cfg.allowedDirs = parseAllowedDirectories(allowedDirsRaw)

	if cfg.operation == "" {
		printUsage()
		os.Exit(1)
	}

	return &cfg
}

func run(cfg *cliConfig) error {
	switch cfg.operation {
	case "read_file":
		return runReadFile(cfg)
	case "write_file":
		return runWriteFile(cfg)
	case "create_directory":
		return runCreateDirectory(cfg)
	case "list_directory":
		return runListDirectory(cfg)
	case "directory_tree":
		return runDirectoryTree(cfg)
	case "move_file":
		return runMoveFile(cfg)
	case "search_files":
		return runSearchFiles(cfg)
	case "get_file_info":
		return runGetFileInfo(cfg)
	case "list_allowed_directories":
		return runListAllowedDirectories(cfg)
	default:
		printUsage()
		return fmt.Errorf("未対応のoperationです: %s", cfg.operation)
	}
}

func runReadFile(cfg *cliConfig) error {
	if cfg.path == "" {
		return fmt.Errorf("read_file操作では-pathが必要です")
	}

	service := newFileSystemService(cfg.allowedDirs, cfg.path)
	content, err := service.ReadFile(cfg.path)
	if err != nil {
		return err
	}

	fmt.Print(content)
	return nil
}

func runWriteFile(cfg *cliConfig) error {
	if cfg.path == "" {
		return fmt.Errorf("write_file操作では-pathが必要です")
	}

	service := newFileSystemService(cfg.allowedDirs, cfg.path)
	if err := service.WriteFile(cfg.path, cfg.content); err != nil {
		return err
	}

	fmt.Printf("ファイル %s への書き込みに成功しました\n", cfg.path)
	return nil
}

func runCreateDirectory(cfg *cliConfig) error {
	if cfg.path == "" {
		return fmt.Errorf("create_directory操作では-pathが必要です")
	}

	service := newFileSystemService(cfg.allowedDirs, cfg.path)
	if err := service.CreateDirectory(cfg.path); err != nil {
		return err
	}

	fmt.Printf("ディレクトリ %s の作成に成功しました\n", cfg.path)
	return nil
}

func runListDirectory(cfg *cliConfig) error {
	if cfg.path == "" {
		return fmt.Errorf("list_directory操作では-pathが必要です")
	}

	service := newFileSystemService(cfg.allowedDirs, cfg.path)
	entries, err := service.ListDirectory(cfg.path)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("(空のディレクトリ)")
		return nil
	}

	fmt.Println(strings.Join(entries, "\n"))
	return nil
}

func runDirectoryTree(cfg *cliConfig) error {
	if cfg.path == "" {
		return fmt.Errorf("directory_tree操作では-pathが必要です")
	}

	service := newFileSystemService(cfg.allowedDirs, cfg.path)
	result, err := service.GetDirectoryTreeAsJSON(cfg.path)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}

func runMoveFile(cfg *cliConfig) error {
	if cfg.source == "" || cfg.destination == "" {
		return fmt.Errorf("move_file操作では-sourceと-destinationが必要です")
	}

	service := newFileSystemService(cfg.allowedDirs, cfg.source, cfg.destination)
	if err := service.MoveFile(cfg.source, cfg.destination); err != nil {
		return err
	}

	fmt.Printf("%s から %s への移動に成功しました\n", cfg.source, cfg.destination)
	return nil
}

func runSearchFiles(cfg *cliConfig) error {
	if cfg.path == "" {
		return fmt.Errorf("search_files操作では-pathが必要です")
	}
	if cfg.pattern == "" {
		return fmt.Errorf("search_files操作では-patternが必要です")
	}

	service := newFileSystemService(cfg.allowedDirs, cfg.path)
	var exclude []string
	if cfg.excludePattern != "" {
		exclude = []string{cfg.excludePattern}
	}

	results, err := service.SearchFiles(cfg.path, cfg.pattern, exclude)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("一致するものが見つかりませんでした")
		return nil
	}

	fmt.Println(strings.Join(results, "\n"))
	return nil
}

func runGetFileInfo(cfg *cliConfig) error {
	if cfg.path == "" {
		return fmt.Errorf("get_file_info操作では-pathが必要です")
	}

	service := newFileSystemService(cfg.allowedDirs, cfg.path)
	info, err := service.GetFileInfoAsText(cfg.path)
	if err != nil {
		return err
	}

	fmt.Print(info)
	return nil
}

func runListAllowedDirectories(cfg *cliConfig) error {
	service := newFileSystemService(cfg.allowedDirs, cfg.path)
	fmt.Print(service.GetAllowedDirectoriesAsText())
	return nil
}

func parseAllowedDirectories(raw string) []string {
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	var dirs []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			dirs = append(dirs, trimmed)
		}
	}
	return dirs
}

func newFileSystemService(base []string, extras ...string) *usecases.FileSystemService {
	dirs := mergeAllowedDirectories(base, extras...)
	if len(dirs) == 0 {
		// フォールバックとしてカレントディレクトリを許可
		cwd, err := os.Getwd()
		if err == nil {
			dirs = []string{cwd}
		}
	}

	return usecases.NewFileSystemService(dirs)
}

func mergeAllowedDirectories(base []string, extras ...string) []string {
	combined := append([]string{}, base...)
	combined = append(combined, extras...)

	seen := make(map[string]struct{})
	result := make([]string, 0, len(combined))
	for _, dir := range combined {
		trimmed := strings.TrimSpace(dir)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}

	return result
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `filesystem CLI - ファイル読み書きや検索を安全に行うツール

使用例:
  go run ./cmd/cli/filesystem -operation=read_file -path=/path/to/file
  go run ./cmd/cli/filesystem -operation=write_file -path=./memo.txt -content="hello"
  go run ./cmd/cli/filesystem -operation=search_files -path=. -pattern=main

利用可能なoperation:
  read_file                指定したファイルを読み取ります
  write_file               ファイルへ内容を書き込みます
  create_directory         ディレクトリを作成します
  list_directory           ディレクトリ配下の一覧を表示します
  directory_tree           ディレクトリ構造をJSONで表示します
  move_file                ファイルまたはディレクトリを移動・リネームします
  search_files             パターンに一致するファイルやディレクトリを検索します
  get_file_info            ファイルやディレクトリのメタ情報を表示します
  list_allowed_directories 現在許可されているディレクトリ一覧を表示します

フラグ:
`)
	flag.PrintDefaults()
}
