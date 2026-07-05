package usecases

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

type Service struct {
	repository filesystem.Repository
	now        func() time.Time
}

func NewService(repository filesystem.Repository) *Service {
	repo := repository
	if repo == nil {
		repo = filesystem.NewRepository()
	}

	return &Service{
		repository: repo,
		now:        time.Now,
	}
}

func (s *Service) PatchMarkdown(targetTitle, targetURL, srcMarkdownContent, srcMarkdownFile, outFilePath string, topHeadingLevel int) (string, error) {
	trimmedTargetTitle := strings.TrimSpace(targetTitle)
	trimmedTargetURL := strings.TrimSpace(targetURL)
	trimmedSrcMarkdownFile := strings.TrimSpace(srcMarkdownFile)
	trimmedOutFilePath := strings.TrimSpace(outFilePath)

	if trimmedTargetTitle == "" {
		return "", fmt.Errorf("--target-title は必須です")
	}
	if trimmedTargetURL == "" {
		return "", fmt.Errorf("--target-url は必須です")
	}
	if strings.TrimSpace(srcMarkdownContent) == "" && trimmedSrcMarkdownFile == "" {
		return "", fmt.Errorf("--src-markdown-content または --src-markdown-file のいずれかは必須です")
	}
	if strings.TrimSpace(srcMarkdownContent) != "" && trimmedSrcMarkdownFile != "" {
		return "", fmt.Errorf("--src-markdown-content と --src-markdown-file は同時に指定できません")
	}
	if trimmedOutFilePath == "" {
		return "", fmt.Errorf("--out-file-path は必須です")
	}
	if strings.Contains(trimmedOutFilePath, ",") {
		return "", fmt.Errorf("--out-file-path にカンマは使用できません")
	}
	if topHeadingLevel < 1 {
		return "", fmt.Errorf("--top-heading-level は 1 以上で指定してください")
	}

	markdownContent, err := s.resolveMarkdownContent(srcMarkdownContent, trimmedSrcMarkdownFile)
	if err != nil {
		return "", err
	}

	normalizedContent := normalizeNewlines(markdownContent)
	if containsHeadingLevel4OrMore(normalizedContent) {
		return "", fmt.Errorf("見出しレベル4以上（#### 以降）は使用できません")
	}

	patchedContent, err := addWebArticleInfo(normalizedContent, trimmedTargetTitle, trimmedTargetURL, topHeadingLevel)
	if err != nil {
		return "", err
	}

	if err := s.repository.WriteFile(trimmedOutFilePath, []byte(patchedContent)); err != nil {
		return "", fmt.Errorf("出力ファイルへの書き込みに失敗しました: %w", err)
	}

	return fmt.Sprintf("出力しました: %s", trimmedOutFilePath), nil
}

type RenameAttachmentsOptions struct {
	SrcDir     string
	Slug       string
	Start      int
	Digits     int
	SortByTime bool
	SortByName bool
	JSON       bool
	Verbose    bool
}

type renameAttachmentResult struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type renameAttachmentsOutput struct {
	RenamedCount int                      `json:"renamed_count"`
	Files        []renameAttachmentResult `json:"files"`
}

func (s *Service) RenameAttachments(opts RenameAttachmentsOptions) (string, error) {
	normalizedOpts, err := s.validateRenameAttachmentsOptions(opts)
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

	timestamp := s.now().Format("20060102-150405")
	results, err := s.buildRenameAttachmentResults(files, normalizedOpts, timestamp)
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

	return formatRenameAttachmentsOutput(results, normalizedOpts.JSON, normalizedOpts.Verbose)
}

func (s *Service) validateRenameAttachmentsOptions(opts RenameAttachmentsOptions) (RenameAttachmentsOptions, error) {
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

func (s *Service) buildRenameAttachmentResults(files []filesystem.FileInfo, opts RenameAttachmentsOptions, timestamp string) ([]renameAttachmentResult, error) {
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

func formatRenameAttachmentsOutput(results []renameAttachmentResult, jsonOutput, verbose bool) (string, error) {
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

func (s *Service) resolveMarkdownContent(srcMarkdownContent, srcMarkdownFile string) (string, error) {
	if srcMarkdownFile != "" {
		data, err := s.repository.ReadFile(srcMarkdownFile)
		if err != nil {
			return "", fmt.Errorf("入力ファイルの読み込みに失敗しました: %w", err)
		}
		return string(data), nil
	}

	return srcMarkdownContent, nil
}

func normalizeNewlines(content string) string {
	return strings.ReplaceAll(content, "\r\n", "\n")
}

func containsHeadingLevel4OrMore(markdownContent string) bool {
	lines := strings.Split(markdownContent, "\n")
	for _, line := range lines {
		if headingLevel(line) >= 4 {
			return true
		}
	}
	return false
}

func addWebArticleInfo(markdownContent, targetTitle, targetURL string, topHeadingLevel int) (string, error) {
	lines := strings.Split(markdownContent, "\n")
	targetHeadingIndex := -1

	for idx, line := range lines {
		if headingLevel(line) == topHeadingLevel {
			targetHeadingIndex = idx
			break
		}
	}

	if targetHeadingIndex == -1 {
		return "", fmt.Errorf("見出しレベル%d が見つかりませんでした", topHeadingLevel)
	}

	linkLine := fmt.Sprintf("- [%s](%s)", targetTitle, targetURL)

	outputLines := make([]string, 0, len(lines)+1)
	outputLines = append(outputLines, lines[:targetHeadingIndex+1]...)
	outputLines = append(outputLines, linkLine)
	outputLines = append(outputLines, lines[targetHeadingIndex+1:]...)

	return strings.Join(outputLines, "\n"), nil
}

func headingLevel(line string) int {
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" || trimmed[0] != '#' {
		return 0
	}

	level := 0
	for level < len(trimmed) && trimmed[level] == '#' {
		level++
	}
	if level == 0 {
		return 0
	}
	if level < len(trimmed) && trimmed[level] != ' ' && trimmed[level] != '\t' {
		return 0
	}

	return level
}
