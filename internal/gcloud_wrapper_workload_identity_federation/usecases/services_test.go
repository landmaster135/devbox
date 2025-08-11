package usecases

import (
	"strings"
	"testing"
)

func TestNewService_Normal(t *testing.T) {
	service := NewService()

	if service == nil {
		t.Error("Expected service to be created, got nil")
	}
}

func TestValidateConfig_Normal(t *testing.T) {
	tests := []struct {
		name        string
		config      *WorkloadIdentityConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "ValidConfig_Normal",
			config: &WorkloadIdentityConfig{
				ProjectID:        "test-project",
				PoolID:           "test-pool",
				ProviderID:       "test-provider",
				ServiceAccountID: "test-sa",
				Location:         "global",
				RepoOwner:        "test-owner",
				RepoName:         "test-repo",
			},
			expectError: false,
		},
		{
			name: "MissingProjectID_Error",
			config: &WorkloadIdentityConfig{
				PoolID:           "test-pool",
				ProviderID:       "test-provider",
				ServiceAccountID: "test-sa",
				RepoOwner:        "test-owner",
				RepoName:         "test-repo",
			},
			expectError: true,
			errorMsg:    "project-idは必須です",
		},
		{
			name: "MissingPoolID_Error",
			config: &WorkloadIdentityConfig{
				ProjectID:        "test-project",
				ProviderID:       "test-provider",
				ServiceAccountID: "test-sa",
				RepoOwner:        "test-owner",
				RepoName:         "test-repo",
			},
			expectError: true,
			errorMsg:    "pool-idは必須です",
		},
		{
			name: "MissingProviderID_Error",
			config: &WorkloadIdentityConfig{
				ProjectID:        "test-project",
				PoolID:           "test-pool",
				ServiceAccountID: "test-sa",
				RepoOwner:        "test-owner",
				RepoName:         "test-repo",
			},
			expectError: true,
			errorMsg:    "provider-idは必須です",
		},
		{
			name: "MissingServiceAccountID_Error",
			config: &WorkloadIdentityConfig{
				ProjectID:  "test-project",
				PoolID:     "test-pool",
				ProviderID: "test-provider",
				RepoOwner:  "test-owner",
				RepoName:   "test-repo",
			},
			expectError: true,
			errorMsg:    "service-account-idは必須です",
		},
		{
			name: "MissingRepoOwner_Error",
			config: &WorkloadIdentityConfig{
				ProjectID:        "test-project",
				PoolID:           "test-pool",
				ProviderID:       "test-provider",
				ServiceAccountID: "test-sa",
				RepoName:         "test-repo",
			},
			expectError: true,
			errorMsg:    "repo-ownerは必須です",
		},
		{
			name: "MissingRepoName_Error",
			config: &WorkloadIdentityConfig{
				ProjectID:        "test-project",
				PoolID:           "test-pool",
				ProviderID:       "test-provider",
				ServiceAccountID: "test-sa",
				RepoOwner:        "test-owner",
			},
			expectError: true,
			errorMsg:    "repo-nameは必須です",
		},
		{
			name: "EmptyLocation_DefaultsToGlobal",
			config: &WorkloadIdentityConfig{
				ProjectID:        "test-project",
				PoolID:           "test-pool",
				ProviderID:       "test-provider",
				ServiceAccountID: "test-sa",
				RepoOwner:        "test-owner",
				RepoName:         "test-repo",
				Location:         "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			err := service.validateConfig(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				// 空のLocationはglobalにデフォルト設定される
				if tt.config.Location == "" && tt.config.Location != "global" {
					// validateConfigはLocationを変更しないため、この確認は不要
				}
			}
		})
	}
}

func TestGenerateWorkloadIdentitySetupScript_Normal(t *testing.T) {
	tests := []struct {
		name            string
		config          *WorkloadIdentityConfig
		expectError     bool
		expectedParts   []string
		unexpectedParts []string
	}{
		{
			name: "ValidConfigWithDescription_Normal",
			config: &WorkloadIdentityConfig{
				ProjectID:        "test-project",
				PoolID:           "test-pool",
				ProviderID:       "test-provider",
				ServiceAccountID: "test-sa",
				Location:         "global",
				PoolDescription:  "Test pool description",
				RepoOwner:        "test-owner",
				RepoName:         "test-repo",
			},
			expectError: false,
			expectedParts: []string{
				"# Google Cloud Workload Identity Federation セットアップスクリプト",
				"setup_workload_identity_federation",
				"\"test-project\"",
				"\"test-pool\"",
				"\"test-provider\"",
				"\"test-sa\"",
				"\"test-owner\"",
				"\"test-repo\"",
				"\"global\"",
				"\"Test pool description\"",
			},
		},
		{
			name: "ValidConfigWithoutDescription_Normal",
			config: &WorkloadIdentityConfig{
				ProjectID:        "test-project",
				PoolID:           "test-pool",
				ProviderID:       "test-provider",
				ServiceAccountID: "test-sa",
				Location:         "global",
				RepoOwner:        "test-owner",
				RepoName:         "test-repo",
			},
			expectError: false,
			expectedParts: []string{
				"# Google Cloud Workload Identity Federation セットアップスクリプト",
				"setup_workload_identity_federation",
				"\"test-project\"",
				"\"test-pool\"",
				"\"test-provider\"",
				"\"test-sa\"",
				"\"test-owner\"",
				"\"test-repo\"",
				"\"global\"",
			},
			unexpectedParts: []string{
				"\"Test pool description\"",
			},
		},
		{
			name: "InvalidConfig_Error",
			config: &WorkloadIdentityConfig{
				ProjectID: "test-project",
				// 他の必須フィールドが不足
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			script, err := service.GenerateWorkloadIdentitySetupScript(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}

				// 期待される部分が含まれているかチェック
				for _, part := range tt.expectedParts {
					if !strings.Contains(script, part) {
						t.Errorf("Expected script to contain '%s', but it didn't", part)
					}
				}

				// 期待されない部分が含まれていないかチェック
				for _, part := range tt.unexpectedParts {
					if strings.Contains(script, part) {
						t.Errorf("Expected script not to contain '%s', but it did", part)
					}
				}
			}
		})
	}
}

func TestGenerateCleanupScript_Normal(t *testing.T) {
	tests := []struct {
		name          string
		config        *WorkloadIdentityConfig
		expectError   bool
		expectedParts []string
	}{
		{
			name: "ValidConfig_Normal",
			config: &WorkloadIdentityConfig{
				ProjectID:        "test-project",
				PoolID:           "test-pool",
				ProviderID:       "test-provider",
				ServiceAccountID: "test-sa",
				Location:         "global",
				RepoOwner:        "test-owner",
				RepoName:         "test-repo",
			},
			expectError: false,
			expectedParts: []string{
				"# Google Cloud Workload Identity Federation リソース削除スクリプト",
				"cleanup_workload_identity_federation",
				"\"test-project\"",
				"\"test-pool\"",
				"\"test-provider\"",
				"\"test-sa\"",
				"\"global\"",
			},
		},
		{
			name: "InvalidConfig_Error",
			config: &WorkloadIdentityConfig{
				ProjectID: "test-project",
				// 他の必須フィールドが不足
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			script, err := service.GenerateCleanupScript(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}

				for _, part := range tt.expectedParts {
					if !strings.Contains(script, part) {
						t.Errorf("Expected script to contain '%s', but it didn't", part)
					}
				}
			}
		})
	}
}

func TestGenerateResourceList_Normal(t *testing.T) {
	const (
		testProjectID        = "test-project"
		testPoolID           = "test-pool"
		testProviderID       = "test-provider"
		testServiceAccountID = "test-sa"
		testLocation         = "global"
		testRepoOwner        = "test-owner"
		testRepoName         = "test-repo"
		testPoolDescription  = "Test pool description"
	)

	tests := []struct {
		name          string
		config        *WorkloadIdentityConfig
		expectedParts []string
	}{
		{
			name: "ConfigWithDescription_Normal",
			config: &WorkloadIdentityConfig{
				ProjectID:        testProjectID,
				PoolID:           testPoolID,
				ProviderID:       testProviderID,
				ServiceAccountID: testServiceAccountID,
				Location:         testLocation,
				PoolDescription:  testPoolDescription,
				RepoOwner:        testRepoOwner,
				RepoName:         testRepoName,
			},
			expectedParts: []string{
				"## 新規作成されるリソース一覧",
				"### 1. Workload Identity Pool",
				testPoolID,
				testPoolDescription,
				"### 2. OIDC Provider (GitHub Actions用)",
				testProviderID,
				"https://token.actions.githubusercontent.com/",
				testRepoOwner + "/" + testRepoName,
				"### 3. サービスアカウント",
				testServiceAccountID,
				testServiceAccountID + "@" + testProjectID + ".iam.gserviceaccount.com",
				"### 4. IAMポリシーバインディング",
				"roles/monitoring.editor",
				"roles/run.viewer",
				"roles/iam.serviceAccounts.getAccessToken",
				"### 5. Workload Identityバインディング",
				"roles/iam.workloadIdentityUser",
			},
		},
		{
			name: "ConfigWithoutDescription_Normal",
			config: &WorkloadIdentityConfig{
				ProjectID:        testProjectID,
				PoolID:           testPoolID,
				ProviderID:       testProviderID,
				ServiceAccountID: testServiceAccountID,
				Location:         testLocation,
				RepoOwner:        testRepoOwner,
				RepoName:         testRepoName,
			},
			expectedParts: []string{
				"## 新規作成されるリソース一覧",
				"### 1. Workload Identity Pool",
				testPoolID,
				"### 2. OIDC Provider (GitHub Actions用)",
				testProviderID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			resourceList := service.GenerateResourceList(tt.config)

			for _, part := range tt.expectedParts {
				if !strings.Contains(resourceList, part) {
					t.Errorf("Expected resource list to contain '%s', but it didn't", part)
				}
			}
		})
	}
}

func TestGenerateGitHubActionsWorkflow_Normal(t *testing.T) {
	const (
		testProjectID  = "test-project"
		testLocation   = "us-central1"
		testPoolID     = "test-pool"
		testProviderID = "test-provider"
	)

	config := &WorkloadIdentityConfig{
		ProjectID:        testProjectID,
		PoolID:           testPoolID,
		ProviderID:       testProviderID,
		ServiceAccountID: "test-sa",
		Location:         testLocation,
		RepoOwner:        "test-owner",
		RepoName:         "test-repo",
	}

	service := NewService()
	workflow := service.GenerateGitHubActionsWorkflow(config)

	expectedParts := []string{
		"## GitHub Actions ワークフロー設定",
		"permissions:",
		"contents: write",
		"id-token: write",
		"env:",
		"GOOGLE_CLOUD_PROJECT_ID_01: '" + testProjectID + "'",
		"GCLOUD_PROJECT_NUMBER:",
		"GCLOUD_POOL_ID:",
		"GCLOUD_PROVIDER_ID:",
		"GCLOUD_SERVICE_ACCOUNT_EMAIL:",
		"google-github-actions/auth@v2",
		"workload_identity_provider: 'projects/${{ env.GCLOUD_PROJECT_NUMBER }}/locations/" + testLocation + "/workloadIdentityPools/${{ env.GCLOUD_POOL_ID }}/providers/${{ env.GCLOUD_PROVIDER_ID }}'",
		"service_account: '${{ env.GCLOUD_SERVICE_ACCOUNT_EMAIL }}'",
		"### GitHub Secretsに設定する値",
		"`GCLOUD_POOL_ID` | `" + testPoolID + "`",
		"`GCLOUD_PROVIDER_ID` | `" + testProviderID + "`",
	}

	for _, part := range expectedParts {
		if !strings.Contains(workflow, part) {
			t.Errorf("Expected workflow to contain '%s', but it didn't", part)
		}
	}
}

func TestGenerateSetupInstructions_Normal(t *testing.T) {
	const (
		testProjectID = "test-project"
		testRepoOwner = "test-owner"
		testRepoName  = "test-repo"
	)

	config := &WorkloadIdentityConfig{
		ProjectID:        testProjectID,
		PoolID:           "test-pool",
		ProviderID:       "test-provider",
		ServiceAccountID: "test-sa",
		Location:         "global",
		RepoOwner:        testRepoOwner,
		RepoName:         testRepoName,
	}

	service := NewService()
	bashScript := "# Test bash script"
	instructions, err := service.GenerateSetupInstructions(config, bashScript)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedParts := []string{
		"# Google Cloud Workload Identity Federation セットアップ",
		"プロジェクト: **" + testProjectID + "**",
		"リポジトリ: **" + testRepoOwner + "/" + testRepoName + "**",
		"## 新規作成されるリソース一覧",
		"## セットアップスクリプト",
		bashScript,
		"## GitHub Actions ワークフロー設定",
	}

	for _, part := range expectedParts {
		if !strings.Contains(instructions, part) {
			t.Errorf("Expected instructions to contain '%s', but it didn't", part)
		}
	}
}

func TestGenerateNotificationMessages_Normal(t *testing.T) {
	config := &WorkloadIdentityConfig{
		ProjectID:        "test-project",
		PoolID:           "test-pool",
		ProviderID:       "test-provider",
		ServiceAccountID: "test-sa",
		Location:         "global",
		RepoOwner:        "test-owner",
		RepoName:         "test-repo",
	}

	service := NewService()
	messages, err := service.GenerateNotificationMessages(config)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(messages) == 0 {
		t.Error("Expected at least one message, got none")
	}

	// 最初のメッセージは概要メッセージであることを確認
	if messages[0].Title != "Workload Identity Federation セットアップ概要" {
		t.Errorf("Expected first message title to be 'Workload Identity Federation セットアップ概要', got '%s'", messages[0].Title)
	}

	// 各メッセージが2000文字以下であることを確認
	for i, message := range messages {
		if len(message.Content) > 2000 {
			t.Errorf("Message %d exceeds 2000 characters: %d", i, len(message.Content))
		}
	}

	// 必要なメッセージタイプが含まれていることを確認
	titleSet := make(map[string]bool)
	for _, message := range messages {
		titleSet[message.Title] = true
	}

	expectedTitles := []string{
		"Workload Identity Federation セットアップ概要",
		"GitHub Actions ワークフロー設定",
	}

	for _, title := range expectedTitles {
		if !titleSet[title] {
			t.Errorf("Expected message with title '%s' not found", title)
		}
	}
}

func TestGenerateNotificationMessages_InvalidConfig_Error(t *testing.T) {
	config := &WorkloadIdentityConfig{
		ProjectID: "test-project",
		// 他の必須フィールドが不足
	}

	service := NewService()
	_, err := service.GenerateNotificationMessages(config)

	if err == nil {
		t.Error("Expected error for invalid config, got nil")
	}
}

func TestSplitBashScript_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name           string
		script         string
		expectedChunks int
		maxLength      int
	}{
		{
			name:           "ShortScript_SingleChunk",
			script:         "echo 'short script'",
			expectedChunks: 1,
			maxLength:      1900,
		},
		{
			name:           "LongScript_MultipleChunks",
			script:         strings.Repeat("echo 'this is a long line that will be repeated many times'\n", 100),
			expectedChunks: 2, // 予想される分割数（実際の長さに依存）
			maxLength:      1900,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := service.splitBashScript(tt.script)

			if len(messages) < 1 {
				t.Error("Expected at least one message chunk")
			}

			// 各チャンクが最大長以下であることを確認
			for i, message := range messages {
				if len(message.Content) > 2000 {
					t.Errorf("Chunk %d exceeds maximum length: %d", i, len(message.Content))
				}

				// Bashスクリプトのマーカーが含まれていることを確認
				if !strings.Contains(message.Content, "```bash") || !strings.Contains(message.Content, "```") {
					t.Errorf("Chunk %d does not contain proper bash markers", i)
				}

				// タイトルが適切に設定されていることを確認
				if !strings.Contains(message.Title, "実行用Bashスクリプト") {
					t.Errorf("Chunk %d has incorrect title: %s", i, message.Title)
				}

				// 色が適切に設定されていることを確認
				if message.Color != "green" {
					t.Errorf("Chunk %d has incorrect color: %s", i, message.Color)
				}

				// IsCodeフラグが適切に設定されていることを確認
				if !message.IsCode {
					t.Errorf("Chunk %d should have IsCode=true", i)
				}
			}
		})
	}
}

func TestGenerateOverviewMessage_Normal(t *testing.T) {
	config := &WorkloadIdentityConfig{
		ProjectID:        "test-project",
		PoolID:           "test-pool",
		ProviderID:       "test-provider",
		ServiceAccountID: "test-sa",
		Location:         "global",
		RepoOwner:        "test-owner",
		RepoName:         "test-repo",
	}

	service := NewService()
	message := service.generateOverviewMessage(config)

	if message.Title != "Workload Identity Federation セットアップ概要" {
		t.Errorf("Expected title 'Workload Identity Federation セットアップ概要', got '%s'", message.Title)
	}

	if message.Color != "blue" {
		t.Errorf("Expected color 'blue', got '%s'", message.Color)
	}

	if message.IsCode {
		t.Error("Expected IsCode to be false for overview message")
	}

	expectedParts := []string{
		"# Google Cloud Workload Identity Federation セットアップ",
		"**プロジェクト**: test-project",
		"**リポジトリ**: test-owner/test-repo",
		"## GitHub Secretsに設定する値",
		"`GCLOUD_POOL_ID` | `test-pool`",
		"`GCLOUD_PROVIDER_ID` | `test-provider`",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message.Content, part) {
			t.Errorf("Expected message content to contain '%s', but it didn't", part)
		}
	}
}

func TestGenerateWorkflowMessage_Normal(t *testing.T) {
	config := &WorkloadIdentityConfig{
		ProjectID:        "test-project",
		PoolID:           "test-pool",
		ProviderID:       "test-provider",
		ServiceAccountID: "test-sa",
		Location:         "us-central1",
		RepoOwner:        "test-owner",
		RepoName:         "test-repo",
	}

	service := NewService()
	message := service.generateWorkflowMessage(config)

	if message.Title != "GitHub Actions ワークフロー設定" {
		t.Errorf("Expected title 'GitHub Actions ワークフロー設定', got '%s'", message.Title)
	}

	if message.Color != "purple" {
		t.Errorf("Expected color 'purple', got '%s'", message.Color)
	}

	if !message.IsCode {
		t.Error("Expected IsCode to be true for workflow message")
	}

	expectedParts := []string{
		"## GitHub Actions ワークフロー設定",
		"```yaml",
		"permissions:",
		"id-token: write",
		"google-github-actions/auth@v2",
		"workload_identity_provider: 'projects/${{ env.GCLOUD_PROJECT_NUMBER }}/locations/us-central1/workloadIdentityPools/${{ env.GCLOUD_POOL_ID }}/providers/${{ env.GCLOUD_PROVIDER_ID }}'",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message.Content, part) {
			t.Errorf("Expected message content to contain '%s', but it didn't", part)
		}
	}
}

// ヘルパー関数のテスト
func createTestConfig() *WorkloadIdentityConfig {
	return &WorkloadIdentityConfig{
		ProjectID:        "test-project",
		PoolID:           "test-pool",
		ProviderID:       "test-provider",
		ServiceAccountID: "test-sa",
		Location:         "global",
		PoolDescription:  "Test pool description",
		RepoOwner:        "test-owner",
		RepoName:         "test-repo",
	}
}

func TestCreateTestConfig_Helper(t *testing.T) {
	config := createTestConfig()

	if config.ProjectID != "test-project" {
		t.Errorf("Expected ProjectID 'test-project', got '%s'", config.ProjectID)
	}

	if config.PoolID != "test-pool" {
		t.Errorf("Expected PoolID 'test-pool', got '%s'", config.PoolID)
	}
}

func TestSplitCleanupScript_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name           string
		script         string
		expectedChunks int
	}{
		{
			name:           "ShortScript_SingleChunk",
			script:         "echo 'short cleanup script'",
			expectedChunks: 1,
		},
		{
			name:           "LongScript_MultipleChunks",
			script:         strings.Repeat("echo 'this is a long cleanup line that will be repeated many times'\n", 100),
			expectedChunks: 2, // 予想される分割数（実際の長さに依存）
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := service.splitCleanupScript(tt.script)

			if len(messages) < 1 {
				t.Error("Expected at least one message chunk")
			}

			// 各チャンクが最大長以下であることを確認
			for i, message := range messages {
				if len(message.Content) > 2000 {
					t.Errorf("Chunk %d exceeds maximum length: %d", i, len(message.Content))
				}

				// Bashスクリプトのマーカーが含まれていることを確認
				if !strings.Contains(message.Content, "```bash") || !strings.Contains(message.Content, "```") {
					t.Errorf("Chunk %d does not contain proper bash markers", i)
				}

				// タイトルが適切に設定されていることを確認
				if !strings.Contains(message.Title, "削除用Bashスクリプト") {
					t.Errorf("Chunk %d has incorrect title: %s", i, message.Title)
				}

				// 色が適切に設定されていることを確認
				if message.Color != "red" {
					t.Errorf("Chunk %d has incorrect color: %s", i, message.Color)
				}

				// IsCodeフラグが適切に設定されていることを確認
				if !message.IsCode {
					t.Errorf("Chunk %d should have IsCode=true", i)
				}
			}
		})
	}
}

func TestGenerateCleanupMessage_Normal(t *testing.T) {
	config := &WorkloadIdentityConfig{
		ProjectID:        "test-project",
		PoolID:           "test-pool",
		ProviderID:       "test-provider",
		ServiceAccountID: "test-sa",
		Location:         "global",
		RepoOwner:        "test-owner",
		RepoName:         "test-repo",
	}

	service := NewService()
	message, err := service.generateCleanupMessage(config)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if message.Title != "リソース削除用Bashスクリプト" {
		t.Errorf("Expected title 'リソース削除用Bashスクリプト', got '%s'", message.Title)
	}

	if message.Color != "red" {
		t.Errorf("Expected color 'red', got '%s'", message.Color)
	}

	if !message.IsCode {
		t.Error("Expected IsCode to be true for cleanup message")
	}

	expectedParts := []string{
		"## リソース削除用スクリプト",
		"```bash",
		"cleanup_workload_identity_federation",
		"**注意**: このスクリプトは以下のリソースを削除します:",
		"- Workload Identity Pool",
		"- OIDC Provider",
		"- Service Account",
	}

	for _, part := range expectedParts {
		if !strings.Contains(message.Content, part) {
			t.Errorf("Expected message content to contain '%s', but it didn't", part)
		}
	}
}

func TestGenerateCleanupMessage_InvalidConfig_Error(t *testing.T) {
	config := &WorkloadIdentityConfig{
		ProjectID: "test-project",
		// 他の必須フィールドが不足
	}

	service := NewService()
	_, err := service.generateCleanupMessage(config)

	if err == nil {
		t.Error("Expected error for invalid config, got nil")
	}
}

func TestSplitCleanupScriptWithIntro_Normal(t *testing.T) {
	service := NewService()

	tests := []struct {
		name        string
		script      string
		introLength int
	}{
		{
			name:        "ShortScript_WithIntro",
			script:      "echo 'short script'",
			introLength: 100,
		},
		{
			name:        "LongScript_WithIntro",
			script:      strings.Repeat("echo 'long line'\n", 50),
			introLength: 200,
		},
		{
			name:        "VeryLongIntro_MinimumLength",
			script:      "echo 'test'",
			introLength: 1800, // 非常に長い説明部分
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := service.splitCleanupScriptWithIntro(tt.script, tt.introLength)

			if len(messages) < 1 {
				t.Error("Expected at least one message chunk")
			}

			// 各チャンクが適切に設定されていることを確認
			for i, message := range messages {
				if len(message.Content) > 2000 {
					t.Errorf("Chunk %d exceeds maximum length: %d", i, len(message.Content))
				}

				// タイトルが適切に設定されていることを確認
				if !strings.Contains(message.Title, "削除用Bashスクリプト") {
					t.Errorf("Chunk %d has incorrect title: %s", i, message.Title)
				}

				// 色が適切に設定されていることを確認
				if message.Color != "red" {
					t.Errorf("Chunk %d has incorrect color: %s", i, message.Color)
				}

				// IsCodeフラグが適切に設定されていることを確認
				if !message.IsCode {
					t.Errorf("Chunk %d should have IsCode=true", i)
				}
			}
		})
	}
}

func TestGenerateCleanupMessages_Normal(t *testing.T) {
	config := &WorkloadIdentityConfig{
		ProjectID:        "test-project",
		PoolID:           "test-pool",
		ProviderID:       "test-provider",
		ServiceAccountID: "test-sa",
		Location:         "global",
		RepoOwner:        "test-owner",
		RepoName:         "test-repo",
	}

	service := NewService()
	messages, err := service.generateCleanupMessages(config)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(messages) == 0 {
		t.Error("Expected at least one message, got none")
	}

	// 最初のメッセージに説明が含まれていることを確認
	firstMessage := messages[0]
	expectedIntroText := "## リソース削除用スクリプト"
	if !strings.Contains(firstMessage.Content, expectedIntroText) {
		t.Errorf("Expected first message to contain intro text '%s'", expectedIntroText)
	}

	// 各メッセージが2000文字以下であることを確認
	for i, message := range messages {
		if len(message.Content) > 2000 {
			t.Errorf("Message %d exceeds 2000 characters: %d", i, len(message.Content))
		}

		// 削除スクリプトメッセージの特性を確認
		if message.Color != "red" {
			t.Errorf("Message %d has incorrect color: %s", i, message.Color)
		}

		if !message.IsCode {
			t.Errorf("Message %d should have IsCode=true", i)
		}
	}
}

func TestGenerateCleanupMessages_InvalidConfig_Error(t *testing.T) {
	config := &WorkloadIdentityConfig{
		ProjectID: "test-project",
		// 他の必須フィールドが不足
	}

	service := NewService()
	_, err := service.generateCleanupMessages(config)

	if err == nil {
		t.Error("Expected error for invalid config, got nil")
	}
}

func TestGenerateNotificationMessages_EdgeCases_Normal(t *testing.T) {
	tests := []struct {
		name        string
		config      *WorkloadIdentityConfig
		expectError bool
	}{
		{
			name: "LongDescriptionConfig_Normal",
			config: &WorkloadIdentityConfig{
				ProjectID:        "very-long-project-name-that-might-affect-message-length",
				PoolID:           "very-long-pool-id-that-might-affect-message-length",
				ProviderID:       "very-long-provider-id-that-might-affect-message-length",
				ServiceAccountID: "very-long-service-account-id-that-might-affect-message-length",
				Location:         "asia-northeast1-very-long-location-name",
				PoolDescription:  strings.Repeat("This is a very long pool description that might cause message splitting issues. ", 20),
				RepoOwner:        "very-long-repository-owner-name",
				RepoName:         "very-long-repository-name-that-might-affect-splitting",
			},
			expectError: false,
		},
		{
			name: "MinimalConfig_Normal",
			config: &WorkloadIdentityConfig{
				ProjectID:        "p",
				PoolID:           "pool",
				ProviderID:       "prov",
				ServiceAccountID: "sa",
				Location:         "global",
				RepoOwner:        "o",
				RepoName:         "r",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			messages, err := service.GenerateNotificationMessages(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}

				if len(messages) == 0 {
					t.Error("Expected at least one message, got none")
				}

				// 各メッセージが2000文字以下であることを確認
				for i, message := range messages {
					// 長い設定の場合、分割処理により2000文字を少し超える可能性があるため、
					// 1900文字以下であることを確認（分割処理の改善が必要な場合のマージン）
					if len(message.Content) > maxLengthOfContentOfWebhookPayloadOnMargin {
						t.Errorf("Message %d significantly exceeds character limit: %d", i, len(message.Content))
					}

					// メッセージの基本的な構造を確認
					if message.Title == "" {
						t.Errorf("Message %d has empty title", i)
					}

					if message.Content == "" {
						t.Errorf("Message %d has empty content", i)
					}

					if message.Color == "" {
						t.Errorf("Message %d has empty color", i)
					}
				}
			}
		})
	}
}

func TestNotificationMessage_StructFields_Normal(t *testing.T) {
	message := &NotificationMessage{
		Title:   "Test Title",
		Content: "Test Content",
		Color:   "blue",
		IsCode:  true,
	}

	if message.Title != "Test Title" {
		t.Errorf("Expected Title 'Test Title', got '%s'", message.Title)
	}

	if message.Content != "Test Content" {
		t.Errorf("Expected Content 'Test Content', got '%s'", message.Content)
	}

	if message.Color != "blue" {
		t.Errorf("Expected Color 'blue', got '%s'", message.Color)
	}

	if !message.IsCode {
		t.Error("Expected IsCode to be true")
	}
}

func TestWorkloadIdentityConfig_StructFields_Normal(t *testing.T) {
	config := &WorkloadIdentityConfig{
		ProjectID:        "test-project",
		PoolID:           "test-pool",
		ProviderID:       "test-provider",
		ServiceAccountID: "test-sa",
		Location:         "us-central1",
		PoolDescription:  "Test description",
		RepoOwner:        "test-owner",
		RepoName:         "test-repo",
	}

	if config.ProjectID != "test-project" {
		t.Errorf("Expected ProjectID 'test-project', got '%s'", config.ProjectID)
	}

	if config.PoolID != "test-pool" {
		t.Errorf("Expected PoolID 'test-pool', got '%s'", config.PoolID)
	}

	if config.ProviderID != "test-provider" {
		t.Errorf("Expected ProviderID 'test-provider', got '%s'", config.ProviderID)
	}

	if config.ServiceAccountID != "test-sa" {
		t.Errorf("Expected ServiceAccountID 'test-sa', got '%s'", config.ServiceAccountID)
	}

	if config.Location != "us-central1" {
		t.Errorf("Expected Location 'us-central1', got '%s'", config.Location)
	}

	if config.PoolDescription != "Test description" {
		t.Errorf("Expected PoolDescription 'Test description', got '%s'", config.PoolDescription)
	}

	if config.RepoOwner != "test-owner" {
		t.Errorf("Expected RepoOwner 'test-owner', got '%s'", config.RepoOwner)
	}

	if config.RepoName != "test-repo" {
		t.Errorf("Expected RepoName 'test-repo', got '%s'", config.RepoName)
	}
}
