package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	usecases "github.com/landmaster135/devbox/internal/filesystem/usecases"
)

type cliConfig struct {
	operation   string
	path        string
	source      string
	destination string
	content     string
	pattern     string
	offset      int
	limit       int
	depth       int
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

	flag.StringVar(&cfg.operation, "operation", "", "実行する操作 (read_file, write_file, create_directory など)")
	flag.StringVar(&cfg.path, "path", "", "対象のファイルまたはディレクトリのパス")
	flag.StringVar(&cfg.source, "source", "", "移動元のパス (move_file操作用)")
	flag.StringVar(&cfg.destination, "destination", "", "移動先のパス (move_file操作用)")
	flag.StringVar(&cfg.content, "content", "", "ファイルに書き込む内容 (write_file操作用)")
	flag.StringVar(&cfg.pattern, "pattern", "", "検索に使用する正規表現 (grep_files操作用)")
	flag.IntVar(&cfg.offset, "offset", 1, "read_file/list_directory/directory_treeでの開始位置 (1始まり)")
	flag.IntVar(&cfg.limit, "limit", 2000, "read_file/list_directory/directory_treeで返す最大件数")
	flag.IntVar(&cfg.depth, "depth", 0, "directory_treeで辿る最大深さ (0は無制限)")
	flag.Parse()

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
	case "grep_files":
		return runGrepFiles(cfg)
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
	if cfg.offset <= 0 {
		return fmt.Errorf("offsetは1以上で指定してください")
	}
	if cfg.limit <= 0 {
		return fmt.Errorf("limitは1以上で指定してください")
	}

	service := newFileSystemService(cfg.path)
	content, err := service.ReadFile(cfg.path, cfg.offset, cfg.limit)
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

	service := newFileSystemService(cfg.path)
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

	service := newFileSystemService(cfg.path)
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
	if cfg.offset <= 0 {
		return fmt.Errorf("list_directory操作では-offsetを1以上にしてください")
	}
	if cfg.limit <= 0 {
		return fmt.Errorf("list_directory操作では-limitを1以上にしてください")
	}

	service := newFileSystemService(cfg.path)
	options := usecases.ListDirectoryOptions{
		Offset: cfg.offset,
		Limit:  cfg.limit,
	}
	entries, err := service.ListDirectoryWithOptions(cfg.path, options)
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
	if cfg.offset <= 0 {
		return fmt.Errorf("directory_tree操作では-offsetを1以上にしてください")
	}
	if cfg.limit <= 0 {
		return fmt.Errorf("directory_tree操作では-limitを1以上にしてください")
	}
	if cfg.depth < 0 {
		return fmt.Errorf("directory_tree操作では-depthを0以上にしてください")
	}

	service := newFileSystemService(cfg.path)
	options := usecases.DirectoryTreeOptions{
		Offset: cfg.offset,
		Limit:  cfg.limit,
		Depth:  cfg.depth,
	}
	result, err := service.GetDirectoryTreeAsYAMLWithOptions(cfg.path, options)
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

	service := newFileSystemService(cfg.source, cfg.destination)
	if err := service.MoveFile(cfg.source, cfg.destination); err != nil {
		return err
	}

	fmt.Printf("%s から %s への移動に成功しました\n", cfg.source, cfg.destination)
	return nil
}

func runGrepFiles(cfg *cliConfig) error {
	if cfg.path == "" {
		return fmt.Errorf("grep_files操作では-pathが必要です")
	}
	if cfg.pattern == "" {
		return fmt.Errorf("grep_files操作では-patternが必要です")
	}

	service := newFileSystemService(cfg.path)
	results, err := service.SearchFiles(cfg.path, cfg.pattern)
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

	service := newFileSystemService(cfg.path)
	info, err := service.GetFileInfoAsText(cfg.path)
	if err != nil {
		return err
	}

	fmt.Print(info)
	return nil
}

func runListAllowedDirectories(cfg *cliConfig) error {
	service := newFileSystemService(cfg.path)
	fmt.Print(service.GetAllowedDirectoriesAsText())
	return nil
}

func newFileSystemService(paths ...string) *usecases.FileSystemService {
	dirs := uniquePaths(paths...)
	if len(dirs) == 0 {
		// フォールバックとしてカレントディレクトリを許可
		cwd, err := os.Getwd()
		if err == nil {
			dirs = []string{cwd}
		}
	}

	return usecases.NewFileSystemService(dirs)
}

func uniquePaths(paths ...string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(paths))
	for _, dir := range paths {
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
  go run ./cmd/cli/filesystem -operation=grep_files -path=. -pattern=main

利用可能なoperation:
  read_file                指定したファイルを読み取ります
  write_file               ファイルへ内容を書き込みます
  create_directory         ディレクトリを作成します
  list_directory           ディレクトリ配下の一覧を表示します
  directory_tree           ディレクトリ構造をYAMLで表示します
  move_file                ファイルまたはディレクトリを移動・リネームします
  grep_files               正規表現に一致するファイルやディレクトリを検索します
  get_file_info            ファイルやディレクトリのメタ情報を表示します
  list_allowed_directories 現在許可されているディレクトリ一覧を表示します

フラグ:
`)
	flag.PrintDefaults()
}
