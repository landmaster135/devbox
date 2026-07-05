package renameattachments

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/web_clipper/infrastructures/filesystem"
)

type Options struct {
	SrcDir     string
	Slug       string
	Start      int
	Digits     int
	SortByTime bool
	SortByName bool
	JSON       bool
	Verbose    bool
}

type Service struct {
	repository filesystem.Repository
}

type renameAttachmentResult struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type renameAttachmentsOutput struct {
	RenamedCount int                      `json:"renamed_count"`
	Files        []renameAttachmentResult `json:"files"`
}

func NewService(repository filesystem.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Execute(opts Options, now time.Time) (string, error) {
	normalizedOpts, err := validateOptions(opts)
	if err != nil {
		return "", err
	}

	entries, err := s.repository.ListDirectory(normalizedOpts.SrcDir)
	if err != nil {
		return "", fmt.Errorf("リネーム対象ディレクトリの読み込みに失敗しました: %w", err)
	}

	files := filterRegularFiles(entries)
	if len(files) == 0 {
		return "", fmt.Errorf("リネーム対象ファイルが見つかりませんでした")
	}

	sortAttachmentFiles(files, normalizedOpts.SortByTime)

	timestamp := now.Format("20060102-150405")
	results, err := buildRenameAttachmentResults(files, normalizedOpts, timestamp)
	if err != nil {
		return "", err
	}

	if err := s.ensureRenameTargetsAvailable(results); err != nil {
		return "", err
	}

	for _, result := range results {
		if result.From == result.To {
			continue
		}
		if err := s.repository.Rename(result.From, result.To); err != nil {
			return "", fmt.Errorf("ファイルのリネームに失敗しました (%s -> %s): %w", result.From, result.To, err)
		}
	}

	return formatOutput(results, normalizedOpts.JSON, normalizedOpts.Verbose)
}

func validateOptions(opts Options) (Options, error) {
	opts.SrcDir = strings.TrimSpace(opts.SrcDir)
	opts.Slug = strings.TrimSpace(opts.Slug)

	if opts.SrcDir == "" {
		return opts, fmt.Errorf("--src-dir は必須です")
	}
	if opts.Slug == "" {
		return opts, fmt.Errorf("--slug は必須です")
	}
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(opts.Slug) {
		return opts, fmt.Errorf("--slug は英小文字、数字、半角ハイフンのみ使用できます")
	}
	if opts.Start < 0 {
		return opts, fmt.Errorf("--start は 0 以上で指定してください")
	}
	if opts.Digits < 1 {
		return opts, fmt.Errorf("--digits は 1 以上で指定してください")
	}
	if opts.SortByTime == opts.SortByName {
		return opts, fmt.Errorf("-time または -name のいずれか一方を指定してください")
	}

	return opts, nil
}

func filterRegularFiles(entries []filesystem.FileInfo) []filesystem.FileInfo {
	files := make([]filesystem.FileInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir {
			continue
		}
		files = append(files, entry)
	}

	return files
}

func sortAttachmentFiles(files []filesystem.FileInfo, sortByTime bool) {
	sort.SliceStable(files, func(i, j int) bool {
		if sortByTime {
			if !files[i].ModTime.Equal(files[j].ModTime) {
				return files[i].ModTime.Before(files[j].ModTime)
			}
		}

		return files[i].Name < files[j].Name
	})
}

func buildRenameAttachmentResults(files []filesystem.FileInfo, opts Options, timestamp string) ([]renameAttachmentResult, error) {
	results := make([]renameAttachmentResult, 0, len(files))
	seenTargets := make(map[string]struct{}, len(files))

	for i, file := range files {
		serial := fmt.Sprintf("%0*d", opts.Digits, opts.Start+i)
		targetName := fmt.Sprintf("web-summary-%s-%s_%s%s", timestamp, opts.Slug, serial, filepath.Ext(file.Name))
		targetPath := filepath.Join(filepath.Dir(file.Path), targetName)

		if _, exists := seenTargets[targetPath]; exists {
			return nil, fmt.Errorf("リネーム先が重複しています: %s", targetPath)
		}
		seenTargets[targetPath] = struct{}{}

		results = append(results, renameAttachmentResult{
			From: file.Path,
			To:   targetPath,
		})
	}

	return results, nil
}

func (s *Service) ensureRenameTargetsAvailable(results []renameAttachmentResult) error {
	for _, result := range results {
		if result.From == result.To {
			continue
		}

		exists, err := s.repository.Exists(result.To)
		if err != nil {
			return fmt.Errorf("リネーム先の存在確認に失敗しました (%s): %w", result.To, err)
		}
		if exists {
			return fmt.Errorf("リネーム先ファイルが既に存在します: %s", result.To)
		}
	}

	return nil
}

func formatOutput(results []renameAttachmentResult, jsonOutput, verbose bool) (string, error) {
	output := renameAttachmentsOutput{
		RenamedCount: len(results),
		Files:        results,
	}

	if jsonOutput {
		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return "", fmt.Errorf("JSON出力の生成に失敗しました: %w", err)
		}
		return string(data) + "\n", nil
	}

	if !verbose {
		return fmt.Sprintf("リネームしました: %d件", len(results)), nil
	}

	lines := []string{fmt.Sprintf("リネームしました: %d件", len(results))}
	for _, result := range results {
		lines = append(lines, fmt.Sprintf("%s -> %s", result.From, result.To))
	}

	return strings.Join(lines, "\n"), nil
}
