package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// #==============================================================#
// ##          Interfaces for Dependency Injection              ##
// #==============================================================#

// FileOpener はファイルオープン操作のインターフェースです
type FileOpener interface {
	Open(name string) (*os.File, error)
}

// DefaultFileOpener は標準のos.Openを使用する実装です
type DefaultFileOpener struct{}

func (o *DefaultFileOpener) Open(name string) (*os.File, error) {
	return os.Open(name)
}

// FileWriter はファイル書き込み操作のインターフェースです
type FileWriter interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// DefaultFileWriter は標準のos.WriteFileを使用する実装です
type DefaultFileWriter struct{}

func (w *DefaultFileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// DirectoryReader はディレクトリ読み取り操作のインターフェースです
type DirectoryReader interface {
	ReadDir(name string) ([]os.DirEntry, error)
}

// DefaultDirectoryReader は標準のos.ReadDirを使用する実装です
type DefaultDirectoryReader struct{}

func (r *DefaultDirectoryReader) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

// FileStat はファイル情報取得操作のインターフェースです
type FileStat interface {
	Stat(name string) (os.FileInfo, error)
}

// DefaultFileStat は標準のos.Statを使用する実装です
type DefaultFileStat struct{}

func (s *DefaultFileStat) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// JSONMarshaler はJSON変換操作のインターフェースです
type JSONMarshaler interface {
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
}

// DefaultJSONMarshaler は標準のjson.MarshalIndentを使用する実装です
type DefaultJSONMarshaler struct{}

func (m *DefaultJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// YAMLMarshaler はYAML変換操作のインターフェースです
type YAMLMarshaler interface {
	Marshal(v interface{}) ([]byte, error)
}

// DefaultYAMLMarshaler は標準のyaml.Marshalを使用する実装です
type DefaultYAMLMarshaler struct{}

func (m *DefaultYAMLMarshaler) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// #==============================================================#
// ##          Data Structures                                   ##
// #==============================================================#

// FileTreeEntry はディレクトリツリーのエントリを表す構造体です
type FileTreeEntry struct {
	Name                string            `json:"name" yaml:"name"`
	Type                FileTreeEntryType `json:"type" yaml:"type"`
	Children            []FileTreeEntry   `json:"children,omitempty" yaml:"children,omitempty"`
	TruncatedChildCount int               `json:"truncatedChildren,omitempty" yaml:"truncated_children,omitempty"`
}

// FileTreeEntryType はディレクトリエントリの種類を表します
type FileTreeEntryType string

const (
	FileTreeEntryTypeDirectory FileTreeEntryType = "directory"
	FileTreeEntryTypeFile      FileTreeEntryType = "file"
	FileTreeEntryTypeSymlink   FileTreeEntryType = "symlink"
	FileTreeEntryTypeUnknown   FileTreeEntryType = "unknown"
)

// DirectoryTreeOptions はdirectory_tree操作のスライス条件を表します
type DirectoryTreeOptions struct {
	Offset int
	Limit  int
	Depth  int
}

// ListDirectoryOptions はlist_directory操作のスライス条件を表します
type ListDirectoryOptions struct {
	Offset int
	Limit  int
}

// FileInfo はファイル情報を表す構造体です
type FileInfo struct {
	Size        int64     `json:"size"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
	Accessed    time.Time `json:"accessed"`
	IsDirectory bool      `json:"isDirectory"`
	IsFile      bool      `json:"isFile"`
	Permissions string    `json:"permissions"`
}

// #==============================================================#
// ##          FileSystemService                                 ##
// #==============================================================#

// FileSystemService はファイルシステム関連の機能を提供する構造体です
type FileSystemService struct {
	allowedDirectories []string
	fileOpener         FileOpener
	fileWriter         FileWriter
	directoryReader    DirectoryReader
	fileStat           FileStat
	jsonMarshaler      JSONMarshaler
	yamlMarshaler      YAMLMarshaler
}

// NewFileSystemService は新しいFileSystemServiceを作成します
func NewFileSystemService(allowedDirs []string) *FileSystemService {
	// パスを正規化し、シンボリックリンクを解決
	normalizedDirs := make([]string, len(allowedDirs))
	for i, dir := range allowedDirs {
		expandedPath := expandHome(dir)
		absolutePath, err := filepath.Abs(expandedPath)
		if err != nil {
			normalizedDirs[i] = filepath.Clean(expandedPath)
			continue
		}

		// シンボリックリンクを解決
		realPath, err := filepath.EvalSymlinks(absolutePath)
		if err != nil {
			normalizedDirs[i] = filepath.Clean(absolutePath)
		} else {
			normalizedDirs[i] = filepath.Clean(realPath)
		}
	}
	return &FileSystemService{
		allowedDirectories: normalizedDirs,
		fileOpener:         &DefaultFileOpener{},
		fileWriter:         &DefaultFileWriter{},
		directoryReader:    &DefaultDirectoryReader{},
		fileStat:           &DefaultFileStat{},
		jsonMarshaler:      &DefaultJSONMarshaler{},
		yamlMarshaler:      &DefaultYAMLMarshaler{},
	}
}

// NewFileSystemServiceWithDependencies はテスト用に依存性を注入できるFileSystemServiceを作成します
func NewFileSystemServiceWithDependencies(
	allowedDirs []string,
	fileOpener FileOpener,
	fileWriter FileWriter,
	directoryReader DirectoryReader,
	fileStat FileStat,
	jsonMarshaler JSONMarshaler,
	yamlMarshaler YAMLMarshaler,
) *FileSystemService {
	normalizedDirs := make([]string, len(allowedDirs))
	for i, dir := range allowedDirs {
		expandedPath := expandHome(dir)
		absolutePath, err := filepath.Abs(expandedPath)
		if err != nil {
			normalizedDirs[i] = filepath.Clean(expandedPath)
			continue
		}

		realPath, err := filepath.EvalSymlinks(absolutePath)
		if err != nil {
			normalizedDirs[i] = filepath.Clean(absolutePath)
		} else {
			normalizedDirs[i] = filepath.Clean(realPath)
		}
	}
	return &FileSystemService{
		allowedDirectories: normalizedDirs,
		fileOpener:         fileOpener,
		fileWriter:         fileWriter,
		directoryReader:    directoryReader,
		fileStat:           fileStat,
		jsonMarshaler:      jsonMarshaler,
		yamlMarshaler:      yamlMarshaler,
	}
}

// expandHome はパス内の ~ をホームディレクトリに展開します
func expandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == "~" {
		return home
	}
	return filepath.Join(home, path[2:])
}

// isPathAllowed はパスが許可されたディレクトリ内にあるかチェックします
func (fs *FileSystemService) isPathAllowed(targetPath string) bool {
	// 両方のパスでシンボリックリンクを解決して比較
	targetReal, err := filepath.EvalSymlinks(targetPath)
	if err != nil {
		targetReal = targetPath
	}
	targetNormalized := filepath.Clean(targetReal)

	for _, dir := range fs.allowedDirectories {
		dirReal, err := filepath.EvalSymlinks(dir)
		if err != nil {
			dirReal = dir
		}
		dirNormalized := filepath.Clean(dirReal)

		// パスの比較（末尾のスラッシュを統一）
		if targetNormalized == dirNormalized || strings.HasPrefix(targetNormalized+string(filepath.Separator), dirNormalized+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

// ValidatePath はパスが許可されたディレクトリ内にあるか確認します
func (fs *FileSystemService) ValidatePath(requestedPath string) (string, error) {
	expandedPath := expandHome(requestedPath)
	absolutePath, err := filepath.Abs(expandedPath)
	if err != nil {
		return "", fmt.Errorf("パスの解決に失敗しました: %w", err)
	}

	// パスが許可されているかチェック
	if fs.isPathAllowed(absolutePath) {
		// シンボリックリンクを解決して返す
		realPath, err := filepath.EvalSymlinks(absolutePath)
		if err != nil {
			return absolutePath, nil
		}
		return realPath, nil
	}

	// 新しいファイルの場合、親ディレクトリを確認
	parentDir := filepath.Dir(absolutePath)
	if fs.isPathAllowed(parentDir) {
		// シンボリックリンクを解決して返す
		realPath, err := filepath.EvalSymlinks(absolutePath)
		if err != nil {
			return absolutePath, nil
		}
		return realPath, nil
	}

	return "", fmt.Errorf("アクセス拒否 - パスが許可されたディレクトリの外にあります: %s", absolutePath)
}

// ReadFile はファイルの内容を読み取ります
func (fs *FileSystemService) ReadFile(path string, offset, limit int) (string, error) {
	validPath, err := fs.ValidatePath(path)
	if err != nil {
		return "", err
	}
	if offset <= 0 {
		return "", fmt.Errorf("offsetは1以上で指定してください")
	}
	if limit <= 0 {
		return "", fmt.Errorf("limitは1以上で指定してください")
	}

	file, err := fs.fileOpener.Open(validPath)
	if err != nil {
		return "", fmt.Errorf("ファイルの読み取りに失敗しました: %w", err)
	}
	defer file.Close()

	content, err := os.ReadFile(validPath)
	if err != nil {
		return "", fmt.Errorf("ファイルの読み取りに失敗しました: %w", err)
	}

	return formatReadFileContent(string(content), offset, limit)
}

const readFileLineMaxLength = 500

func formatReadFileContent(content string, offset, limit int) (string, error) {
	lines := splitLinesForReadFile(content)
	if len(lines) == 0 {
		if offset <= 1 {
			return "", nil
		}
		return "", fmt.Errorf("offsetがファイルの総行数を超えています")
	}

	startIndex := offset - 1
	if startIndex < 0 || startIndex >= len(lines) {
		return "", fmt.Errorf("offsetがファイルの総行数を超えています")
	}

	endIndex := startIndex + limit
	if endIndex > len(lines) {
		endIndex = len(lines)
	}

	var builder strings.Builder
	for i := startIndex; i < endIndex; i++ {
		if builder.Len() > 0 {
			builder.WriteByte('\n')
		}
		builder.WriteByte('L')
		builder.WriteString(strconv.Itoa(i + 1))
		builder.WriteString(": ")
		builder.WriteString(truncateReadFileLine(lines[i], readFileLineMaxLength))
	}

	return builder.String(), nil
}

func splitLinesForReadFile(content string) []string {
	if content == "" {
		return nil
	}

	lines := strings.Split(content, "\n")
	if strings.HasSuffix(content, "\n") && len(lines) > 0 {
		lines = lines[:len(lines)-1]
	}

	for i, line := range lines {
		lines[i] = strings.TrimSuffix(line, "\r")
	}

	return lines
}

func truncateReadFileLine(line string, limit int) string {
	if limit <= 0 {
		return ""
	}

	runeCount := 0
	for i := range line {
		if runeCount == limit {
			return line[:i]
		}
		runeCount++
	}

	return line
}

// WriteFile はファイルに内容を書き込みます
func (fs *FileSystemService) WriteFile(path string, content string) error {
	validPath, err := fs.ValidatePath(path)
	if err != nil {
		return err
	}

	// 親ディレクトリが存在することを確認
	parentDir := filepath.Dir(validPath)
	if _, err := fs.fileStat.Stat(parentDir); os.IsNotExist(err) {
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("ディレクトリの作成に失敗しました: %w", err)
		}
	}

	err = fs.fileWriter.WriteFile(validPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}

// CreateDirectory はディレクトリを作成します
func (fs *FileSystemService) CreateDirectory(path string) error {
	validPath, err := fs.ValidatePath(path)
	if err != nil {
		return err
	}

	err = os.MkdirAll(validPath, 0755)
	if err != nil {
		return fmt.Errorf("ディレクトリの作成に失敗しました: %w", err)
	}

	return nil
}

// ListDirectory はディレクトリの内容を一覧表示します
func (fs *FileSystemService) ListDirectory(path string) ([]string, error) {
	return fs.ListDirectoryWithOptions(path, ListDirectoryOptions{
		Offset: 1,
		Limit:  0,
	})
}

// ListDirectoryWithOptions はoffset/limitを考慮してディレクトリを一覧表示します
func (fs *FileSystemService) ListDirectoryWithOptions(path string, opts ListDirectoryOptions) ([]string, error) {
	normalized := normalizeListDirectoryOptions(opts)

	validPath, err := fs.ValidatePath(path)
	if err != nil {
		return nil, err
	}

	entries, err := fs.directoryReader.ReadDir(validPath)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの読み取りに失敗しました: %w", err)
	}
	totalEntries := len(entries)

	if totalEntries == 0 {
		if normalized.Offset == 1 {
			return []string{}, nil
		}
		return nil, fmt.Errorf("offset exceeds directory entry count")
	}

	startIndex := normalized.Offset - 1
	if startIndex < 0 || startIndex >= totalEntries {
		return nil, fmt.Errorf("offset exceeds directory entry count")
	}

	endIndex := totalEntries
	if normalized.Limit > 0 {
		endIndex = startIndex + normalized.Limit
		if endIndex > totalEntries {
			endIndex = totalEntries
		}
	}

	slicedEntries := entries[startIndex:endIndex]
	result := make([]string, 0, len(slicedEntries)+1)
	for _, entry := range slicedEntries {
		prefix := "[FILE]"
		if entry.IsDir() {
			prefix = "[DIR]"
		}
		result = append(result, fmt.Sprintf("%s %s", prefix, entry.Name()))
	}

	if normalized.Limit > 0 && endIndex < totalEntries {
		nextOffset := normalized.Offset + len(slicedEntries)
		result = append(result, fmt.Sprintf("More than %d entries found (try -offset=%d)", normalized.Limit, nextOffset))
	}

	return result, nil
}

// GetDirectoryTree はディレクトリの階層構造を取得します
func (fs *FileSystemService) GetDirectoryTree(path string) ([]FileTreeEntry, error) {
	entries, _, _, err := fs.GetDirectoryTreeWithOptions(path, DirectoryTreeOptions{
		Offset: 1,
		Limit:  0,
		Depth:  0,
	})
	return entries, err
}

// GetDirectoryTreeWithOptions はoffset/limit/depthを考慮してツリーを取得します
func (fs *FileSystemService) GetDirectoryTreeWithOptions(path string, opts DirectoryTreeOptions) ([]FileTreeEntry, bool, int, error) {
	normalized := normalizeDirectoryTreeOptions(opts)

	validPath, err := fs.ValidatePath(path)
	if err != nil {
		return nil, false, 0, err
	}

	entries, err := fs.buildDirectoryTreeEntries(validPath, normalized.Depth)
	if err != nil {
		return nil, false, 0, err
	}

	return sliceFileTreeEntries(entries, normalized.Offset, normalized.Limit)
}

func (fs *FileSystemService) buildDirectoryTreeEntries(path string, remainingDepth int) ([]FileTreeEntry, error) {
	entries, err := fs.directoryReader.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの読み取りに失敗しました: %w", err)
	}

	var result []FileTreeEntry
	for _, entry := range entries {
		fileEntry := FileTreeEntry{
			Name: entry.Name(),
			Type: determineEntryType(entry),
		}

		if entry.IsDir() {
			subPath := filepath.Join(path, entry.Name())
			if shouldDescend(remainingDepth) {
				nextDepth := remainingDepth
				if remainingDepth > 1 {
					nextDepth--
				}
				children, err := fs.buildDirectoryTreeEntries(subPath, nextDepth)
				if err != nil {
					return nil, err
				}
				fileEntry.Children = children
			} else {
				count, err := fs.countDirectoryEntries(subPath)
				if err != nil {
					return nil, err
				}
				fileEntry.TruncatedChildCount = count
			}
		}

		result = append(result, fileEntry)
	}

	return result, nil
}

func shouldDescend(remainingDepth int) bool {
	return remainingDepth == 0 || remainingDepth > 1
}

func determineEntryType(entry os.DirEntry) FileTreeEntryType {
	if entry.IsDir() {
		return FileTreeEntryTypeDirectory
	}

	mode := entry.Type()
	if mode == os.ModeIrregular {
		info, err := entry.Info()
		if err != nil {
			return FileTreeEntryTypeUnknown
		}
		mode = info.Mode().Type()
	}

	switch {
	case mode&os.ModeSymlink != 0:
		return FileTreeEntryTypeSymlink
	case mode == 0:
		return FileTreeEntryTypeFile
	default:
		return FileTreeEntryTypeUnknown
	}
}

func (fs *FileSystemService) countDirectoryEntries(path string) (int, error) {
	entries, err := fs.directoryReader.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("ディレクトリの読み取りに失敗しました: %w", err)
	}
	return len(entries), nil
}

func normalizeDirectoryTreeOptions(opts DirectoryTreeOptions) DirectoryTreeOptions {
	normalized := DirectoryTreeOptions{
		Offset: opts.Offset,
		Limit:  opts.Limit,
		Depth:  opts.Depth,
	}

	if normalized.Offset < 1 {
		normalized.Offset = 1
	}

	if normalized.Limit < 0 {
		normalized.Limit = 0
	}

	if normalized.Depth < 0 {
		normalized.Depth = 0
	}

	return normalized
}

func sliceFileTreeEntries(entries []FileTreeEntry, offset, limit int) ([]FileTreeEntry, bool, int, error) {
	if len(entries) == 0 {
		return []FileTreeEntry{}, false, offset, nil
	}

	startIndex := offset - 1
	if startIndex >= len(entries) {
		return nil, false, 0, fmt.Errorf("offset exceeds directory entry count")
	}

	endIndex := len(entries)
	if limit > 0 {
		endIndex = startIndex + limit
		if endIndex > len(entries) {
			endIndex = len(entries)
		}
	}

	sliced := make([]FileTreeEntry, endIndex-startIndex)
	copy(sliced, entries[startIndex:endIndex])
	hasMore := limit > 0 && endIndex < len(entries)
	nextOffset := offset + (endIndex - startIndex)
	return sliced, hasMore, nextOffset, nil
}

func normalizeListDirectoryOptions(opts ListDirectoryOptions) ListDirectoryOptions {
	normalized := ListDirectoryOptions{
		Offset: opts.Offset,
		Limit:  opts.Limit,
	}

	if normalized.Offset < 1 {
		normalized.Offset = 1
	}

	if normalized.Limit < 0 {
		normalized.Limit = 0
	}

	return normalized
}

// GetDirectoryTreeAsJSON はディレクトリの階層構造をJSON文字列として取得します
func (fs *FileSystemService) GetDirectoryTreeAsJSON(path string) (string, error) {
	return fs.GetDirectoryTreeAsJSONWithOptions(path, DirectoryTreeOptions{
		Offset: 1,
		Limit:  0,
		Depth:  0,
	})
}

// GetDirectoryTreeAsJSONWithOptions はJSON出力に対してoffset/limit/depthを適用します
func (fs *FileSystemService) GetDirectoryTreeAsJSONWithOptions(path string, opts DirectoryTreeOptions) (string, error) {
	tree, _, _, err := fs.GetDirectoryTreeWithOptions(path, opts)
	if err != nil {
		return "", err
	}

	jsonBytes, err := fs.jsonMarshaler.MarshalIndent(tree, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSONの生成に失敗しました: %w", err)
	}

	return string(jsonBytes), nil
}

// GetDirectoryTreeAsYAML はディレクトリの階層構造をYAML文字列として取得します
func (fs *FileSystemService) GetDirectoryTreeAsYAML(path string) (string, error) {
	return fs.GetDirectoryTreeAsYAMLWithOptions(path, DirectoryTreeOptions{
		Offset: 1,
		Limit:  0,
		Depth:  0,
	})
}

// GetDirectoryTreeAsYAMLWithOptions はYAML出力に対してoffset/limit/depthを適用します
func (fs *FileSystemService) GetDirectoryTreeAsYAMLWithOptions(path string, opts DirectoryTreeOptions) (string, error) {
	tree, hasMore, nextOffset, err := fs.GetDirectoryTreeWithOptions(path, opts)
	if err != nil {
		return "", err
	}

	yamlBytes, err := fs.yamlMarshaler.Marshal(tree)
	if err != nil {
		return "", fmt.Errorf("YAMLの生成に失敗しました: %w", err)
	}

	result := string(yamlBytes)
	if opts.Limit > 0 && hasMore {
		if !strings.HasSuffix(result, "\n") {
			result += "\n"
		}
		result += fmt.Sprintf("More than %d entries found (try -offset=%d)", opts.Limit, nextOffset)
	}

	return result, nil
}

// MoveFile はファイルを移動します
func (fs *FileSystemService) MoveFile(source, destination string) error {
	validSource, err := fs.ValidatePath(source)
	if err != nil {
		return err
	}

	validDest, err := fs.ValidatePath(destination)
	if err != nil {
		return err
	}

	err = os.Rename(validSource, validDest)
	if err != nil {
		return fmt.Errorf("ファイルの移動に失敗しました: %w", err)
	}

	return nil
}

// SearchFiles はファイルを検索します
func (fs *FileSystemService) SearchFiles(rootPath, pattern string, excludePatterns []string) ([]string, error) {
	validPath, err := fs.ValidatePath(rootPath)
	if err != nil {
		return nil, err
	}

	var results []string
	err = filepath.Walk(validPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // エラーがあっても続行
		}

		// 除外パターンに一致するかチェック
		for _, excludePattern := range excludePatterns {
			matched, err := filepath.Match(excludePattern, filepath.Base(path))
			if err == nil && matched {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// パターンに一致するかチェック
		if strings.Contains(strings.ToLower(filepath.Base(path)), strings.ToLower(pattern)) {
			results = append(results, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ファイルの検索に失敗しました: %w", err)
	}

	return results, nil
}

// GetFileInfo はファイル情報を取得します
func (fs *FileSystemService) GetFileInfo(path string) (*FileInfo, error) {
	validPath, err := fs.ValidatePath(path)
	if err != nil {
		return nil, err
	}

	info, err := fs.fileStat.Stat(validPath)
	if err != nil {
		return nil, fmt.Errorf("ファイル情報の取得に失敗しました: %w", err)
	}

	// 一部のプラットフォームでは、これらの時間が利用できない場合があります
	var created, accessed time.Time
	// 時間取得
	created = info.ModTime()  // フォールバックとして変更時間を使用
	accessed = info.ModTime() // フォールバックとして変更時間を使用

	return &FileInfo{
		Size:        info.Size(),
		Created:     created,
		Modified:    info.ModTime(),
		Accessed:    accessed,
		IsDirectory: info.IsDir(),
		IsFile:      !info.IsDir(),
		Permissions: fmt.Sprintf("%o", info.Mode().Perm()),
	}, nil
}

// GetFileInfoAsText はファイル情報をテキスト形式で取得します
func (fs *FileSystemService) GetFileInfoAsText(path string) (string, error) {
	info, err := fs.GetFileInfo(path)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("サイズ: %d バイト\n", info.Size)
	result += fmt.Sprintf("作成日時: %s\n", info.Created.Format(time.RFC3339))
	result += fmt.Sprintf("最終変更: %s\n", info.Modified.Format(time.RFC3339))
	result += fmt.Sprintf("最終アクセス: %s\n", info.Accessed.Format(time.RFC3339))
	result += fmt.Sprintf("ディレクトリ: %t\n", info.IsDirectory)
	result += fmt.Sprintf("ファイル: %t\n", info.IsFile)
	result += fmt.Sprintf("権限: %s", info.Permissions)

	return result, nil
}

// GetAllowedDirectories は許可されたディレクトリのリストを取得します
func (fs *FileSystemService) GetAllowedDirectories() []string {
	return fs.allowedDirectories
}

// GetAllowedDirectoriesAsText は許可されたディレクトリのリストをテキスト形式で取得します
func (fs *FileSystemService) GetAllowedDirectoriesAsText() string {
	return "許可されたディレクトリ:\n" + strings.Join(fs.allowedDirectories, "\n")
}
