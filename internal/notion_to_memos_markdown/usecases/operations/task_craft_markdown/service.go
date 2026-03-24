package taskcraftmarkdown

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	markdownconfig "github.com/landmaster135/devbox/internal/markdown_crafter/config"

	domain "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/domain"
	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
	common "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/common"
)

const (
	taskSplitHeadingLevel = 2

	doneStatusID    = "0198d2e4-9a5e-7127-82b2-3541e7eda226"
	wipStatusID     = "0198d2e4-9a5e-7165-9b5b-80fde571d270"
	plannedStatusID = "0198d2e4-9a5e-775c-9998-2074bf1a5ba9"
	planned2Status  = "0198d2e4-9a5e-74e9-9d85-94797dd1f663"
)

var (
	digitsOnlyRegexp   = regexp.MustCompile(`^\d+$`)
	priorityDigitRegex = regexp.MustCompile(`\d+`)
)

type Service struct {
	fileSystem filesystem.Repository
}

func NewService(fileSystem filesystem.Repository) *Service {
	return &Service{
		fileSystem: fileSystem,
	}
}

func (s *Service) Execute(pageType, category string, skipsNoSrcBody bool, conNumberStart, conNumberEnd int, srcJSONFile, srcBodyDir, outDir string) (string, error) {
	trimmedPageType := strings.TrimSpace(pageType)
	if trimmedPageType != common.SupportedPageTypeTask {
		return "", fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
	}
	_ = category

	if conNumberStart <= 0 || conNumberEnd <= 0 {
		return "", fmt.Errorf("con_number_start と con_number_end は1以上で指定してください")
	}
	if conNumberStart > conNumberEnd {
		return "", fmt.Errorf("con_number_start は con_number_end 以下である必要があります")
	}

	rawJSON, err := s.fileSystem.ReadFile(srcJSONFile)
	if err != nil {
		return "", fmt.Errorf("src-json-file の読み込みに失敗しました: %w", err)
	}

	var tasks []domain.Task
	if err := json.Unmarshal(rawJSON, &tasks); err != nil {
		return "", fmt.Errorf("Task JSONの解析に失敗しました: %w", err)
	}

	if err := s.fileSystem.MkdirAll(outDir, common.DefaultDirectoryPerm); err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗しました (%s): %w", outDir, err)
	}

	markdownService := common.NewMarkdownService(s.fileSystem)

	srcBodyFiles, err := s.fileSystem.ListMarkdownFiles(srcBodyDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("src-body-dir の読み込みに失敗しました (%s): %w", srcBodyDir, err)
		}
		srcBodyFiles = []string{}
	}
	splitIndex := common.BuildSplitMarkdownIndex(srcBodyFiles)

	totalTarget := 0
	processed := 0
	skipped := 0
	seenConID := map[string]struct{}{}

	for _, task := range tasks {
		conID := strings.TrimSpace(task.ConID)
		if conID == "" {
			continue
		}
		if _, exists := seenConID[conID]; exists {
			continue
		}
		seenConID[conID] = struct{}{}

		conNumber, err := common.ParseConNumber(conID)
		if err != nil {
			return "", fmt.Errorf("con_id の解析に失敗しました (%s): %w", conID, err)
		}
		if conNumber < conNumberStart || conNumber > conNumberEnd {
			continue
		}

		timestamp, err := resolveTimestampForTask(task)
		if err != nil {
			return "", fmt.Errorf("日時の解決に失敗しました (con_id=%s): %w", conID, err)
		}
		tags := buildTagsForTask(task)

		outputPaths, missingSource, err := s.prepareTaskOutputs(markdownService, conID, srcBodyDir, outDir, splitIndex)
		if err != nil {
			return "", err
		}

		if missingSource {
			totalTarget++
			if skipsNoSrcBody {
				skipped++
				continue
			}
			fallbackPath := filepath.Join(outDir, conID+"_01.md")
			if err := s.fileSystem.WriteFile(fallbackPath, []byte("")); err != nil {
				return "", fmt.Errorf("空Markdownの作成に失敗しました (%s): %w", fallbackPath, err)
			}
			outputPaths = []string{fallbackPath}
		} else {
			totalTarget += len(outputPaths)
		}

		for i, outputPath := range outputPaths {
			if err := applyTaskMarkdown(markdownService, outputPath, task.PageTitle, tags); err != nil {
				return "", fmt.Errorf("Task Markdownの加工に失敗しました (con_id=%s): %w", conID, err)
			}

			suffix := resolveTaskSuffix(conID, outputPath, i)
			renamedPath := filepath.Join(outDir, fmt.Sprintf("%s_%s.md", timestamp, suffix))
			if err := s.fileSystem.RenameFile(outputPath, renamedPath); err != nil {
				return "", fmt.Errorf("ファイル名の変更に失敗しました (%s -> %s): %w", outputPath, renamedPath, err)
			}
			processed++
		}
	}

	return fmt.Sprintf("処理完了\n対象件数=%d, 加工成功=%d, スキップ=%d", totalTarget, processed, skipped), nil
}

func (s *Service) prepareTaskOutputs(markdownService commonMarkdownService, conID, srcBodyDir, outDir string, splitIndex common.SplitMarkdownIndex) ([]string, bool, error) {
	exactSrcPath := filepath.Join(srcBodyDir, conID+".md")
	exactExists, err := s.fileSystem.FileExists(exactSrcPath)
	if err != nil {
		return nil, false, fmt.Errorf("コピー元ファイルの確認に失敗しました (%s): %w", exactSrcPath, err)
	}

	if exactExists {
		splitPaths, splitSuccess, err := s.splitSourceMarkdown(markdownService, exactSrcPath, outDir)
		if err != nil {
			return nil, false, err
		}
		if splitSuccess {
			return splitPaths, false, nil
		}

		fallbackPath := filepath.Join(outDir, conID+"_01.md")
		if err := s.fileSystem.CopyFile(exactSrcPath, fallbackPath); err != nil {
			return nil, false, fmt.Errorf("Markdownのコピーに失敗しました (%s -> %s): %w", exactSrcPath, fallbackPath, err)
		}
		return []string{fallbackPath}, false, nil
	}

	sourcePaths := splitIndex.Resolve(conID)
	if len(sourcePaths) == 0 {
		return nil, true, nil
	}

	outputPaths := make([]string, 0, len(sourcePaths))
	for _, srcPath := range sourcePaths {
		dstPath := filepath.Join(outDir, filepath.Base(srcPath))
		if err := s.fileSystem.CopyFile(srcPath, dstPath); err != nil {
			return nil, false, fmt.Errorf("Markdownのコピーに失敗しました (%s -> %s): %w", srcPath, dstPath, err)
		}
		outputPaths = append(outputPaths, dstPath)
	}

	return outputPaths, false, nil
}

type commonMarkdownService interface {
	SplitHeadings(filePath string, headingLevel int, outputDir string) (string, error)
	AddHeading1(filePath, headingText, headingPosition string) (string, error)
	AddTags(filePath string, tagsCSV string) (string, error)
}

func (s *Service) splitSourceMarkdown(markdownService commonMarkdownService, sourcePath, outputDir string) ([]string, bool, error) {
	beforeFiles, err := s.fileSystem.ListMarkdownFiles(outputDir)
	if err != nil {
		return nil, false, fmt.Errorf("split前ファイル一覧の取得に失敗しました (%s): %w", outputDir, err)
	}
	beforeSet := make(map[string]struct{}, len(beforeFiles))
	for _, path := range beforeFiles {
		beforeSet[path] = struct{}{}
	}

	_, err = markdownService.SplitHeadings(sourcePath, taskSplitHeadingLevel, outputDir)
	if err != nil {
		if strings.Contains(err.Error(), "見出しレベル") {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("split-headings の実行に失敗しました (%s): %w", sourcePath, err)
	}

	files, err := s.fileSystem.ListMarkdownFiles(outputDir)
	if err != nil {
		return nil, false, fmt.Errorf("split後ファイル一覧の取得に失敗しました (%s): %w", outputDir, err)
	}

	result := make([]string, 0, len(files))
	for _, path := range files {
		if _, exists := beforeSet[path]; exists {
			continue
		}
		baseName := filepath.Base(path)
		stem := strings.TrimSuffix(baseName, filepath.Ext(baseName))
		if digitsOnlyRegexp.MatchString(stem) {
			result = append(result, path)
		}
	}
	if len(result) == 0 {
		return nil, false, nil
	}
	return result, true, nil
}

func applyTaskMarkdown(markdownService commonMarkdownService, filePath, pageTitle string, tags []string) error {
	headingText := strings.TrimSpace(pageTitle)
	if headingText == "" {
		return fmt.Errorf("page_title が空です")
	}
	if _, err := markdownService.AddHeading1(filePath, headingText, markdownconfig.HeadingPositionHead); err != nil {
		return fmt.Errorf("見出し追加に失敗しました: %w", err)
	}

	if _, err := markdownService.AddTags(filePath, strings.Join(tags, ",")); err != nil {
		return fmt.Errorf("タグ追加に失敗しました: %w", err)
	}
	return nil
}

func buildTagsForTask(task domain.Task) []string {
	tags := make([]string, 0, 3)
	if statusTag := mapTaskStatusTag(task.StatusID); statusTag != "" {
		tags = append(tags, statusTag)
	}
	if priorityTag := mapTaskPriorityTag(task.Priority); priorityTag != "" {
		tags = append(tags, priorityTag)
	}
	tags = append(tags, common.RequiredBackupTagTask)

	return uniqueTaskTags(tags)
}

func mapTaskStatusTag(statusID *string) string {
	normalizedStatusID := normalizeStatusID(statusID)
	switch normalizedStatusID {
	case doneStatusID:
		return "01-p/todo-status/3-done"
	case wipStatusID:
		return "01-p/todo-status/2-wip"
	case plannedStatusID, planned2Status:
		return "01-p/todo-status/1-planned"
	default:
		return ""
	}
}

func mapTaskPriorityTag(priority domain.TaskPriority) string {
	level, ok := resolveTaskPriorityLevel(priority)
	if !ok {
		return ""
	}
	switch level {
	case 5:
		return "01-p/todo-prior/5-high"
	case 3:
		return "01-p/todo-prior/3-mid"
	case 1:
		return "01-p/todo-prior/1-low"
	default:
		return ""
	}
}

func resolveTaskPriorityLevel(priority domain.TaskPriority) (int, bool) {
	if priority.IntValue != nil {
		return *priority.IntValue, true
	}

	digits := priorityDigitRegex.FindString(strings.TrimSpace(priority.Text))
	if digits == "" {
		return 0, false
	}
	priorityValue, err := strconv.Atoi(digits)
	if err != nil {
		return 0, false
	}
	return priorityValue, true
}

func resolveTimestampForTask(task domain.Task) (string, error) {
	timeSource := strings.TrimSpace(task.UpdatedAt)
	if normalizeStatusID(task.StatusID) == doneStatusID {
		doneAtStart := strings.TrimSpace(task.DoneAtStart)
		if doneAtStart != "" {
			timeSource = doneAtStart
		}
	}
	if timeSource == "" {
		return "", fmt.Errorf("updated_at が空です")
	}

	parsed, err := time.Parse(time.RFC3339Nano, timeSource)
	if err != nil {
		return "", fmt.Errorf("日時の解析に失敗しました: %w", err)
	}
	return parsed.Format("20060102150405"), nil
}

func normalizeStatusID(statusID *string) string {
	if statusID == nil {
		return ""
	}
	return strings.TrimSpace(*statusID)
}

func resolveTaskSuffix(conID, filePath string, index int) string {
	baseName := filepath.Base(filePath)
	stem := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	candidate := ""
	if strings.HasPrefix(stem, conID+"_") {
		candidate = strings.TrimSpace(strings.TrimPrefix(stem, conID+"_"))
	} else if digitsOnlyRegexp.MatchString(stem) {
		candidate = stem
	} else if underscoreIndex := strings.LastIndex(stem, "_"); underscoreIndex >= 0 && underscoreIndex < len(stem)-1 {
		suffix := strings.TrimSpace(stem[underscoreIndex+1:])
		if digitsOnlyRegexp.MatchString(suffix) {
			candidate = suffix
		}
	}

	if digitsOnlyRegexp.MatchString(candidate) {
		number, err := strconv.Atoi(candidate)
		if err == nil {
			return fmt.Sprintf("%02d", number)
		}
	}
	return fmt.Sprintf("%02d", index+1)
}

func uniqueTaskTags(tags []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(tags))

	for _, tag := range tags {
		normalized := strings.TrimPrefix(strings.TrimSpace(tag), "#")
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result
}
