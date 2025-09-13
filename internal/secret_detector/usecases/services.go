package usecases

import (
	"fmt"
	"math"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/landmaster135/devbox/internal/secret_detector/config"
	"github.com/landmaster135/devbox/internal/secret_detector/domain"
)

// SecretDetectorService はシークレット検知サービス
type SecretDetectorService struct{}

// NewSecretDetectorService は新しいSecretDetectorServiceを作成
func NewSecretDetectorService() *SecretDetectorService {
	return &SecretDetectorService{}
}

// GetStagedFiles はGitのステージされたファイルを取得
func (s *SecretDetectorService) GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	var result []string
	for _, file := range files {
		if file != "" {
			result = append(result, file)
		}
	}
	return result, nil
}

// FilterConfigFiles は設定ファイルをフィルタリング
func (s *SecretDetectorService) FilterConfigFiles(files []string) []string {
	patterns := domain.GetConfigFilePatterns()

	var configFiles []string
	for _, file := range files {
		for _, pattern := range patterns {
			matched, _ := filepath.Match(pattern, filepath.Base(file))
			if matched {
				configFiles = append(configFiles, file)
				break
			}
		}
	}
	return configFiles
}

// CheckFile はファイルをチェックしてシークレットを検知
func (s *SecretDetectorService) CheckFile(filename string) ([]domain.SecretResult, error) {
	config, err := config.LoadConfig(filename)
	if err != nil {
		return nil, err
	}

	var results []domain.SecretResult

	fmt.Printf("  Found %d server(s) in configuration\n", len(config.MCPServers))

	totalEnvVars := 0
	for serverName, serverConfig := range config.MCPServers {
		if serverConfig.Env != nil {
			totalEnvVars += len(serverConfig.Env)
			for key, value := range serverConfig.Env {
				result := domain.SecretResult{
					File:          filename,
					Server:        serverName,
					Key:           key,
					Value:         value,
					IsPlaceholder: s.IsPlaceholder(value),
				}

				// 疑わしいキー名かチェック
				if s.IsSuspiciousKey(key) {
					result.MatchedPattern = "suspicious_key"
					results = append(results, result)
				}

				// 実際のシークレットパターンかチェック
				if pattern := s.MatchesSecretPattern(value); pattern != "" {
					result.MatchedPattern = pattern
					results = append(results, result)
				}
			}
		}
	}

	fmt.Printf("  Found %d environment variable(s) total\n", totalEnvVars)

	return results, nil
}

// IsPlaceholder はプレースホルダー値かどうかを判定
func (s *SecretDetectorService) IsPlaceholder(value string) bool {
	// 空文字は許可
	if value == "" {
		return true
	}

	// 許可されたプレースホルダー値
	allowedPlaceholders := domain.GetAllowedPlaceholders()
	for _, placeholder := range allowedPlaceholders {
		if value == placeholder {
			return true
		}
	}

	// YOUR_で始まる値も許可
	if strings.HasPrefix(value, "YOUR_") {
		return true
	}

	// 短すぎる値
	if len(value) < 8 {
		return true
	}

	// テスト用の値
	testPatterns := domain.GetTestPatterns()
	for _, pattern := range testPatterns {
		if matched, _ := regexp.MatchString(pattern, value); matched {
			return true
		}
	}

	return false
}

// IsSuspiciousKey は疑わしいキー名かどうかを判定
func (s *SecretDetectorService) IsSuspiciousKey(key string) bool {
	suspiciousPatterns := domain.GetSuspiciousPatterns()
	for _, pattern := range suspiciousPatterns {
		if matched, _ := regexp.MatchString(pattern, key); matched {
			return true
		}
	}
	return false
}

// MatchesSecretPattern は実際のシークレットパターンにマッチするかを判定
func (s *SecretDetectorService) MatchesSecretPattern(value string) string {
	realSecretPatterns := domain.GetRealSecretPatterns()
	for _, pattern := range realSecretPatterns {
		if matched, _ := regexp.MatchString(pattern, value); matched {
			return pattern
		}
	}
	return ""
}

// AnalyzeResults は結果を分析して表示
func (s *SecretDetectorService) AnalyzeResults(results []domain.SecretResult) domain.ScanSummary {
	summary := domain.ScanSummary{}

	for _, result := range results {
		if result.MatchedPattern == "suspicious_key" {
			if result.IsPlaceholder {
				fmt.Printf("    %s✅ %s.%s: safe placeholder value%s\n",
					domain.Green, result.Server, result.Key, domain.Reset)
				summary.PlaceholderCount++
			} else {
				fmt.Printf("    %s❌ %s.%s: potential real secret detected%s\n",
					domain.Red, result.Server, result.Key, domain.Reset)
				fmt.Printf("      %sValue:%s %s\n", domain.Yellow, domain.Reset, result.Value)
				summary.HasRealSecrets = true
				summary.SecretCount++
			}
		} else if result.MatchedPattern != "" {
			if result.IsPlaceholder {
				fmt.Printf("    %s✅ %s.%s: pattern matched but safe%s\n",
					domain.Green, result.Server, result.Key, domain.Reset)
				summary.PlaceholderCount++
			} else {
				fmt.Printf("    %s❌ %s.%s: suspicious pattern detected%s\n",
					domain.Red, result.Server, result.Key, domain.Reset)
				fmt.Printf("      %sPattern:%s %s\n", domain.Yellow, domain.Reset, result.MatchedPattern)
				fmt.Printf("      %sValue:%s %s\n", domain.Yellow, domain.Reset, result.Value)
				summary.HasRealSecrets = true
				summary.SecretCount++
			}
		}
	}

	// 高エントロピー値の検証
	for _, result := range results {
		if result.MatchedPattern == "" && !result.IsPlaceholder {
			entropy := s.CalculateEntropy(result.Value)
			if entropy > 4.0 && len(result.Value) > 20 {
				fmt.Printf("    %s⚠️  %s.%s: high entropy value (%.2f)%s\n",
					domain.Yellow, result.Server, result.Key, entropy, domain.Reset)
				fmt.Printf("      %sValue:%s %s\n", domain.Yellow, domain.Reset, result.Value)
			}
		}
	}

	fmt.Printf("\n  %sSummary:%s Found %d placeholder(s), %d potential secret(s)\n",
		domain.Blue, domain.Reset, summary.PlaceholderCount, summary.SecretCount)

	return summary
}

// CalculateEntropy はエントロピーを計算
func (s *SecretDetectorService) CalculateEntropy(str string) float64 {
	if len(str) == 0 {
		return 0
	}

	freq := make(map[rune]int)
	for _, char := range str {
		freq[char]++
	}

	entropy := 0.0
	length := float64(len(str))
	for _, count := range freq {
		p := float64(count) / length
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}

	return entropy
}
