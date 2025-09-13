package usecases

import (
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/landmaster135/devbox/internal/secret_detector/config"
	"github.com/landmaster135/devbox/internal/secret_detector/domain"
)

// SecretDetectorService はシークレット検知サービス
type SecretDetectorService struct {
	verbose         bool
	dryRun          bool
	commandExecutor CommandExecutorRepository
}

// NewSecretDetectorService は新しいSecretDetectorServiceを作成
func NewSecretDetectorService(verbose, dryRun bool, commandExecutor CommandExecutorRepository) *SecretDetectorService {
	return &SecretDetectorService{
		verbose:         verbose,
		dryRun:          dryRun,
		commandExecutor: commandExecutor,
	}
}

// GetStagedFiles はGitのステージされたファイルを取得
func (s *SecretDetectorService) GetStagedFiles() ([]string, error) {
	if s.dryRun {
		if s.verbose {
			fmt.Printf("%s[DRY-RUN] Skipping Git operations%s\n", domain.Yellow, domain.Reset)
		}
		return []string{}, nil
	}

	if s.verbose {
		fmt.Printf("%s[VERBOSE] Executing: git diff --cached --name-only --diff-filter=ACM%s\n", domain.Blue, domain.Reset)
	}

	output, err := s.commandExecutor.Execute("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
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

	if s.verbose {
		fmt.Printf("%s[VERBOSE] Found %d staged file(s)%s\n", domain.Blue, len(result), domain.Reset)
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

	if s.verbose {
		fmt.Printf("%s[VERBOSE] Filtered %d config file(s) from %d total file(s)%s\n",
			domain.Blue, len(configFiles), len(files), domain.Reset)
	}

	return configFiles
}

// CheckSpecificFile は特定のファイルをチェック
func (s *SecretDetectorService) CheckSpecificFile(filename string) ([]domain.SecretResult, error) {
	if s.verbose {
		fmt.Printf("%s[VERBOSE] Checking specific file: %s%s\n", domain.Blue, filename, domain.Reset)
	}

	return s.CheckFile(filename)
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

// stripProtocolPrefix はプロトコル識別子を除去する
func (s *SecretDetectorService) stripProtocolPrefix(value string) string {
	protocols := domain.GetProtocolPrefixes()
	for _, protocol := range protocols {
		if strings.HasPrefix(value, protocol) {
			if s.verbose {
				fmt.Printf("%s[VERBOSE] Stripped protocol prefix '%s' from value%s\n",
					domain.Blue, protocol, domain.Reset)
			}
			return value[len(protocol):]
		}
	}
	return value
}

// IsPlaceholder はプレースホルダー値かどうかを判定
func (s *SecretDetectorService) IsPlaceholder(value string) bool {
	// 空文字は許可
	if value == "" {
		return true
	}

	// プロトコル識別子を除去してからチェック
	strippedValue := s.stripProtocolPrefix(value)

	// 除去後の値が空の場合は許可
	if strippedValue == "" {
		return true
	}

	// 許可されたプレースホルダー値（元の値とstripped値の両方をチェック）
	allowedPlaceholders := domain.GetAllowedPlaceholders()
	for _, placeholder := range allowedPlaceholders {
		if value == placeholder || strippedValue == placeholder {
			return true
		}
	}

	// YOUR_で始まる値も許可（stripped値でチェック）
	if strings.HasPrefix(strippedValue, "YOUR_") {
		return true
	}

	// 短すぎる値（stripped値でチェック）
	if len(strippedValue) < 8 {
		return true
	}

	// テスト用の値（stripped値でチェック）
	testPatterns := domain.GetTestPatterns()
	for _, pattern := range testPatterns {
		if matched, _ := regexp.MatchString(pattern, strippedValue); matched {
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
