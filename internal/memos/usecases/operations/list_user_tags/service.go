package listusertags

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"unicode"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

// Service は list_user_tags operation を扱う。
type Service struct {
	client     *common.JSONClient
	fileSystem infrastructures.FileSystem
}

func New(client *common.JSONClient, fileSystem infrastructures.FileSystem) *Service {
	return &Service{client: client, fileSystem: fileSystem}
}

func (s *Service) Execute(ctx context.Context, userID string, outputDir string) (*common.ListUserTagsOutput, error) {
	cleanUserID := strings.TrimSpace(userID)
	if cleanUserID == "" {
		return nil, fmt.Errorf("userID が空です")
	}
	cleanOutputDir := strings.TrimSpace(outputDir)
	if cleanOutputDir == "" {
		return nil, fmt.Errorf("outputDir が空です")
	}

	var stats common.UserStats
	requestPath := "/users/" + url.PathEscape(cleanUserID) + ":getStats"
	if err := s.client.DoJSON(ctx, http.MethodGet, requestPath, nil, nil, &stats); err != nil {
		return nil, err
	}

	tagCount := stats.TagCount
	if tagCount == nil {
		tagCount = map[string]int{}
	}

	outputRoot, err := s.fileSystem.EnsureDirectory(cleanOutputDir)
	if err != nil {
		return nil, err
	}

	data, err := json.MarshalIndent(tagCount, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("tagCount の JSON 整形に失敗しました: %w", err)
	}
	data = append(data, '\n')

	outputPath := filepath.Join(outputRoot, "user-tags_"+sanitizeFileComponent(cleanUserID)+".json")
	if err := s.fileSystem.WriteFile(outputPath, data); err != nil {
		return nil, err
	}

	return &common.ListUserTagsOutput{
		UserID:     cleanUserID,
		OutputPath: outputPath,
		TagCount:   tagCount,
	}, nil
}

func sanitizeFileComponent(value string) string {
	cleanValue := strings.TrimSpace(value)
	if cleanValue == "" {
		return "user"
	}

	var builder strings.Builder
	lastUnderscore := false
	for _, r := range cleanValue {
		if isSafeFileRune(r) {
			builder.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if lastUnderscore {
			continue
		}
		builder.WriteByte('_')
		lastUnderscore = true
	}

	result := strings.Trim(builder.String(), "._-")
	if result == "" {
		return "user"
	}
	return result
}

func isSafeFileRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
}
