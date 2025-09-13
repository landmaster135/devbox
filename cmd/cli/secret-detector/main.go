package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/secret_detector/config"
	"github.com/landmaster135/devbox/internal/secret_detector/domain"
	"github.com/landmaster135/devbox/internal/secret_detector/usecases"
)

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// バージョンが要求された場合
	if cfg.Version {
		fmt.Println("Secret Detector v1.0.0")
		return
	}

	// ヘルプが要求された場合
	if cfg.Help {
		config.PrintUsage()
		return
	}

	fmt.Printf("%s🔍 Checking for secrets in JSON configuration files...%s\n", domain.Green, domain.Reset)

	// コマンド実行器を初期化
	commandExecutor := usecases.NewCommandExecutor()

	// シークレット検知サービスを初期化
	service := usecases.NewSecretDetectorService(cfg.Verbose, cfg.DryRun, commandExecutor)

	var allResults []domain.SecretResult
	var totalFiles int

	// 特定ファイルが指定された場合
	if cfg.ConfigFile != "" {
		fmt.Printf("%s📁 Checking specific configuration file: %s%s\n", domain.Green, cfg.ConfigFile, domain.Reset)
		totalFiles = 1

		results, err := service.CheckSpecificFile(cfg.ConfigFile)
		if err != nil {
			fmt.Printf("%s❌ Error checking file %s: %v%s\n", domain.Red, cfg.ConfigFile, err, domain.Reset)
			os.Exit(1)
		}

		allResults = append(allResults, results...)
	} else {
		// Gitのステージされたファイルを取得
		stagedFiles, err := service.GetStagedFiles()
		if err != nil {
			fmt.Printf("%s❌ Error getting staged files: %v%s\n", domain.Red, err, domain.Reset)
			os.Exit(1)
		}

		// 設定ファイルをフィルタリング
		configFiles := service.FilterConfigFiles(stagedFiles)

		if len(configFiles) == 0 {
			fmt.Printf("%sℹ️  No configuration files found in this commit.%s\n", domain.Green, domain.Reset)
			os.Exit(0)
		}

		totalFiles = len(configFiles)

		for _, file := range configFiles {
			fmt.Printf("%s📁 Checking configuration file: %s%s\n", domain.Green, file, domain.Reset)

			results, err := service.CheckFile(file)
			if err != nil {
				fmt.Printf("%s❌ Error checking file %s: %v%s\n", domain.Red, file, err, domain.Reset)
				continue
			}

			allResults = append(allResults, results...)
		}
	}

	// 結果の分析と表示
	summary := service.AnalyzeResults(allResults)

	fmt.Println()
	if summary.HasRealSecrets {
		fmt.Printf("%s❌ COMMIT BLOCKED: Potential secrets detected!%s\n", domain.Red, domain.Reset)
		fmt.Println()
		fmt.Printf("%s💡 Suggestions:%s\n", domain.Yellow, domain.Reset)
		fmt.Println("  1. Replace actual values with placeholders like 'YOUR_API_KEY'")
		fmt.Println("  2. Move secrets to environment variables or secure storage")
		fmt.Println("  3. Use a separate config file that's git-ignored")
		fmt.Println("  4. If this is intentional, use: git commit --no-verify")
		fmt.Println()
		os.Exit(1)
	} else {
		fmt.Printf("%s✅ Secret scan completed. No actual secrets detected in %d file(s).%s\n",
			domain.Green, totalFiles, domain.Reset)
		os.Exit(0)
	}
}
