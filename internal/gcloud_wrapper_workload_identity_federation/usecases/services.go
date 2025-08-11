package usecases

import (
	"fmt"
	"strings"
)

// WorkloadIdentityConfig はWorkload Identity Federationの設定を保持する構造体
type WorkloadIdentityConfig struct {
	ProjectID        string
	PoolID           string
	ProviderID       string
	ServiceAccountID string
	Location         string
	PoolDescription  string
	RepoOwner        string
	RepoName         string
}

// Service はWorkload Identity Federation設定のサービス
type Service struct{}

// NewService は新しいServiceインスタンスを作成する
func NewService() *Service {
	return &Service{}
}

// GenerateWorkloadIdentitySetupScript はWorkload Identity Federationセットアップ用のBash関数を生成する
func (s *Service) GenerateWorkloadIdentitySetupScript(config *WorkloadIdentityConfig) (string, error) {
	if err := s.validateConfig(config); err != nil {
		return "", fmt.Errorf("設定の検証に失敗しました: %v", err)
	}

	var script strings.Builder

	script.WriteString("# Google Cloud Workload Identity Federation セットアップスクリプト\n")
	script.WriteString("# 注意: このスクリプトを実行する前に、iam.shファイルをsourceしてください\n\n")

	// 使用例
	script.WriteString("# 使用例:\n")
	script.WriteString("# source /path/to/iam.sh\n")
	script.WriteString("# setup_workload_identity_federation \\\n")
	script.WriteString(fmt.Sprintf("#     \"%s\" \\\n", config.ProjectID))
	script.WriteString(fmt.Sprintf("#     \"%s\" \\\n", config.PoolID))
	script.WriteString(fmt.Sprintf("#     \"%s\" \\\n", config.ProviderID))
	script.WriteString(fmt.Sprintf("#     \"%s\" \\\n", config.ServiceAccountID))
	script.WriteString(fmt.Sprintf("#     \"%s\" \\\n", config.RepoOwner))
	script.WriteString(fmt.Sprintf("#     \"%s\" \\\n", config.RepoName))
	script.WriteString(fmt.Sprintf("#     \"%s\"", config.Location))
	if config.PoolDescription != "" {
		script.WriteString(fmt.Sprintf(" \\\n#     \"%s\"", config.PoolDescription))
	}
	script.WriteString("\n\n")

	// 実際の関数呼び出し
	script.WriteString("# iam.shで定義された関数を呼び出し\n")
	script.WriteString("setup_workload_identity_federation")
	script.WriteString(fmt.Sprintf(" \"%s\"", config.ProjectID))
	script.WriteString(fmt.Sprintf(" \"%s\"", config.PoolID))
	script.WriteString(fmt.Sprintf(" \"%s\"", config.ProviderID))
	script.WriteString(fmt.Sprintf(" \"%s\"", config.ServiceAccountID))
	script.WriteString(fmt.Sprintf(" \"%s\"", config.RepoOwner))
	script.WriteString(fmt.Sprintf(" \"%s\"", config.RepoName))
	script.WriteString(fmt.Sprintf(" \"%s\"", config.Location))
	if config.PoolDescription != "" {
		script.WriteString(fmt.Sprintf(" \"%s\"", config.PoolDescription))
	}
	script.WriteString("\n")

	return script.String(), nil
}

// GenerateCleanupScript は既存のiam.sh関数を使用したリソース削除用のBashスクリプトを生成する
func (s *Service) GenerateCleanupScript(config *WorkloadIdentityConfig) (string, error) {
	if err := s.validateConfig(config); err != nil {
		return "", fmt.Errorf("設定の検証に失敗しました: %v", err)
	}

	var script strings.Builder

	script.WriteString("# Google Cloud Workload Identity Federation リソース削除スクリプト\n")
	script.WriteString("# 注意: このスクリプトを実行する前に、iam.shファイルをsourceしてください\n\n")

	// 使用例
	script.WriteString("# 使用例:\n")
	script.WriteString("# source /path/to/iam.sh\n")
	script.WriteString("# cleanup_workload_identity_federation \\\n")
	script.WriteString(fmt.Sprintf("#     \"%s\" \\\n", config.ProjectID))
	script.WriteString(fmt.Sprintf("#     \"%s\" \\\n", config.PoolID))
	script.WriteString(fmt.Sprintf("#     \"%s\" \\\n", config.ProviderID))
	script.WriteString(fmt.Sprintf("#     \"%s\" \\\n", config.ServiceAccountID))
	script.WriteString(fmt.Sprintf("#     \"%s\"\n", config.Location))
	script.WriteString("\n")

	// 実際の関数呼び出し
	script.WriteString("# iam.shで定義された関数を呼び出し\n")
	script.WriteString("cleanup_workload_identity_federation")
	script.WriteString(fmt.Sprintf(" \"%s\"", config.ProjectID))
	script.WriteString(fmt.Sprintf(" \"%s\"", config.PoolID))
	script.WriteString(fmt.Sprintf(" \"%s\"", config.ProviderID))
	script.WriteString(fmt.Sprintf(" \"%s\"", config.ServiceAccountID))
	script.WriteString(fmt.Sprintf(" \"%s\"", config.Location))
	script.WriteString("\n")

	return script.String(), nil
}

// GenerateResourceList は作成されるリソース一覧を生成する
func (s *Service) GenerateResourceList(config *WorkloadIdentityConfig) string {
	var resources strings.Builder

	resources.WriteString("## 新規作成されるリソース一覧\n\n")

	resources.WriteString("### 1. Workload Identity Pool\n")
	resources.WriteString(fmt.Sprintf("- **名前**: %s\n", config.PoolID))
	resources.WriteString(fmt.Sprintf("- **フルパス**: projects/%s/locations/%s/workloadIdentityPools/%s\n", config.ProjectID, config.Location, config.PoolID))
	if config.PoolDescription != "" {
		resources.WriteString(fmt.Sprintf("- **説明**: %s\n", config.PoolDescription))
	}
	resources.WriteString("\n")

	resources.WriteString("### 2. OIDC Provider (GitHub Actions用)\n")
	resources.WriteString(fmt.Sprintf("- **名前**: %s\n", config.ProviderID))
	resources.WriteString(fmt.Sprintf("- **フルパス**: projects/%s/locations/%s/workloadIdentityPools/%s/providers/%s\n", config.ProjectID, config.Location, config.PoolID, config.ProviderID))
	resources.WriteString("- **発行者URI**: https://token.actions.githubusercontent.com/\n")
	resources.WriteString(fmt.Sprintf("- **対象リポジトリ**: %s/%s\n", config.RepoOwner, config.RepoName))
	resources.WriteString("\n")

	resources.WriteString("### 3. サービスアカウント\n")
	resources.WriteString(fmt.Sprintf("- **ID**: %s\n", config.ServiceAccountID))
	resources.WriteString(fmt.Sprintf("- **メールアドレス**: %s@%s.iam.gserviceaccount.com\n", config.ServiceAccountID, config.ProjectID))
	resources.WriteString("\n")

	resources.WriteString("### 4. IAMポリシーバインディング\n")
	resources.WriteString("以下のロールがサービスアカウントに付与されます:\n")
	resources.WriteString("- **roles/monitoring.editor** - Cloud Monitoringの編集権限\n")
	resources.WriteString("- **roles/run.viewer** - Cloud Runの閲覧権限\n")
	resources.WriteString("- **roles/iam.serviceAccounts.getAccessToken** - サービスアカウントトークン取得権限\n")
	resources.WriteString("\n")

	resources.WriteString("### 5. Workload Identityバインディング\n")
	resources.WriteString(fmt.Sprintf("- **プリンシパル**: principalSet://iam.googleapis.com/projects/{PROJECT_NUMBER}/locations/%s/workloadIdentityPools/%s/attribute.repository/%s/%s\n", config.Location, config.PoolID, config.RepoOwner, config.RepoName))
	resources.WriteString("- **ロール**: roles/iam.workloadIdentityUser\n")

	return resources.String()
}

// GenerateGitHubActionsWorkflow はGitHub Actions用のワークフロー設定を生成する
func (s *Service) GenerateGitHubActionsWorkflow(config *WorkloadIdentityConfig) string {
	var workflow strings.Builder

	workflow.WriteString("## GitHub Actions ワークフロー設定\n\n")
	workflow.WriteString("以下の設定を `.github/workflows/` ディレクトリ内のYAMLファイルに追加してください:\n\n")
	workflow.WriteString("```yaml\n")
	workflow.WriteString("on:\n")
	workflow.WriteString("  push:\n\n")
	workflow.WriteString("permissions:\n")
	workflow.WriteString("  contents: write # リポジトリへの書き込み権限 (バッジ更新のため)\n")
	workflow.WriteString("  id-token: write # OIDCトークンをリクエストする権限 (Google Cloud認証のため)\n\n")
	workflow.WriteString("env:\n")
	workflow.WriteString(fmt.Sprintf("  GOOGLE_CLOUD_PROJECT_ID_01: '%s'\n", config.ProjectID))
	workflow.WriteString("  GCLOUD_PROJECT_NUMBER: ${{ secrets.GCLOUD_PROJECT_NUMBER }}\n")
	workflow.WriteString("  GCLOUD_POOL_ID: ${{ secrets.GCLOUD_POOL_ID }}\n")
	workflow.WriteString("  GCLOUD_PROVIDER_ID: ${{ secrets.GCLOUD_PROVIDER_ID }}\n")
	workflow.WriteString("  GCLOUD_SERVICE_ACCOUNT_EMAIL: ${{ secrets.GCLOUD_SERVICE_ACCOUNT_EMAIL }}\n\n")
	workflow.WriteString("jobs:\n")
	workflow.WriteString("  test:\n")
	workflow.WriteString("    runs-on: ubuntu-latest\n")
	workflow.WriteString("    if: contains(github.event.head_commit.message, '[skip ci]') == false\n")
	workflow.WriteString("    steps:\n")
	workflow.WriteString("      - name: Checkout\n")
	workflow.WriteString("        uses: actions/checkout@v4\n")
	workflow.WriteString("      - name: Set up Golang\n")
	workflow.WriteString("        uses: actions/setup-go@v5\n")
	workflow.WriteString("        with:\n")
	workflow.WriteString("          go-version-file: 'go.mod'\n")
	workflow.WriteString("      - id: 'gcloud_auth'\n")
	workflow.WriteString("        name: 'Authenticate to Google Cloud'\n")
	workflow.WriteString("        uses: 'google-github-actions/auth@v2'\n")
	workflow.WriteString("        with:\n")
	workflow.WriteString("          create_credentials_file: true\n")
	workflow.WriteString(fmt.Sprintf("          workload_identity_provider: 'projects/${{ env.GCLOUD_PROJECT_NUMBER }}/locations/%s/workloadIdentityPools/${{ env.GCLOUD_POOL_ID }}/providers/${{ env.GCLOUD_PROVIDER_ID }}'\n", config.Location))
	workflow.WriteString("          service_account: '${{ env.GCLOUD_SERVICE_ACCOUNT_EMAIL }}'\n")
	workflow.WriteString("      - name: 'Set up Cloud SDK'\n")
	workflow.WriteString("        uses: 'google-github-actions/setup-gcloud@v2'\n")
	workflow.WriteString("      - name: 'Use gcloud CLI'\n")
	workflow.WriteString("        run: 'gcloud info'\n")
	workflow.WriteString("```\n\n")

	workflow.WriteString("### GitHub Secretsに設定する値\n\n")
	workflow.WriteString("リポジトリの Settings > Secrets and variables > Actions で以下のSecretsを設定してください:\n\n")
	workflow.WriteString("| Secret名 | 値 | 説明 |\n")
	workflow.WriteString("|---------|----|---------|\n")
	workflow.WriteString("| `GCLOUD_PROJECT_NUMBER` | スクリプト実行後に表示される値 | Google Cloudプロジェクト番号 |\n")
	workflow.WriteString(fmt.Sprintf("| `GCLOUD_POOL_ID` | `%s` | Workload Identity Pool ID |\n", config.PoolID))
	workflow.WriteString(fmt.Sprintf("| `GCLOUD_PROVIDER_ID` | `%s` | OIDC Provider ID |\n", config.ProviderID))
	workflow.WriteString(fmt.Sprintf("| `GCLOUD_SERVICE_ACCOUNT_EMAIL` | `%s@%s.iam.gserviceaccount.com` | サービスアカウントのメールアドレス |\n", config.ServiceAccountID, config.ProjectID))

	return workflow.String()
}

// GenerateSetupInstructions はセットアップ手順を生成する
func (s *Service) GenerateSetupInstructions(config *WorkloadIdentityConfig, bashScript string) (string, error) {
	var instructions strings.Builder

	instructions.WriteString("# Google Cloud Workload Identity Federation セットアップ\n\n")
	instructions.WriteString(fmt.Sprintf("プロジェクト: **%s**\n", config.ProjectID))
	instructions.WriteString(fmt.Sprintf("リポジトリ: **%s/%s**\n\n", config.RepoOwner, config.RepoName))

	// リソース一覧
	resourceList := s.GenerateResourceList(config)
	instructions.WriteString(resourceList)
	instructions.WriteString("\n")

	// Bashスクリプト
	instructions.WriteString("## セットアップスクリプト\n\n")
	instructions.WriteString("以下のスクリプトを実行してWorkload Identity Federationをセットアップしてください:\n\n")
	instructions.WriteString("```bash\n")
	instructions.WriteString(bashScript)
	instructions.WriteString("```\n\n")

	// GitHub Actions設定
	workflowConfig := s.GenerateGitHubActionsWorkflow(config)
	instructions.WriteString(workflowConfig)

	return instructions.String(), nil
}

// NotificationMessage はDiscord通知用のメッセージ構造体
type NotificationMessage struct {
	Title   string
	Content string
	Color   string
	IsCode  bool
}

// GenerateNotificationMessages は全ての通知メッセージを生成し、2000文字制限に基づいて自動分割する
func (s *Service) GenerateNotificationMessages(config *WorkloadIdentityConfig) ([]*NotificationMessage, error) {
	if err := s.validateConfig(config); err != nil {
		return nil, fmt.Errorf("設定の検証に失敗しました: %v", err)
	}

	var messages []*NotificationMessage

	// 1. 概要とリソース一覧
	overviewMessage := s.generateOverviewMessage(config)
	messages = append(messages, overviewMessage)

	// 2. Bashスクリプト（分割）
	bashScript, err := s.GenerateWorkloadIdentitySetupScript(config)
	if err != nil {
		return nil, fmt.Errorf("bashスクリプトの生成に失敗しました: %v", err)
	}

	bashMessages := s.splitBashScript(bashScript)
	messages = append(messages, bashMessages...)

	// 3. 削除スクリプト（分割）
	cleanupMessages, err := s.generateCleanupMessages(config)
	if err != nil {
		return nil, fmt.Errorf("削除スクリプトの生成に失敗しました: %v", err)
	}
	messages = append(messages, cleanupMessages...)

	// 4. GitHub Actions設定
	workflowMessage := s.generateWorkflowMessage(config)
	messages = append(messages, workflowMessage)

	return messages, nil
}

// generateOverviewMessage は概要とリソース一覧のメッセージを生成する
func (s *Service) generateOverviewMessage(config *WorkloadIdentityConfig) *NotificationMessage {
	var content strings.Builder

	content.WriteString("# Google Cloud Workload Identity Federation セットアップ\n\n")
	content.WriteString(fmt.Sprintf("**プロジェクト**: %s\n", config.ProjectID))
	content.WriteString(fmt.Sprintf("**リポジトリ**: %s/%s\n\n", config.RepoOwner, config.RepoName))

	// リソース一覧
	resourceList := s.GenerateResourceList(config)
	content.WriteString(resourceList)

	content.WriteString("\n## GitHub Secretsに設定する値\n\n")
	content.WriteString("| Secret名 | 値 |\n")
	content.WriteString("|---------|----|\n")
	content.WriteString("| `GCLOUD_PROJECT_NUMBER` | スクリプト実行後に表示される値 |\n")
	content.WriteString(fmt.Sprintf("| `GCLOUD_POOL_ID` | `%s` |\n", config.PoolID))
	content.WriteString(fmt.Sprintf("| `GCLOUD_PROVIDER_ID` | `%s` |\n", config.ProviderID))
	content.WriteString(fmt.Sprintf("| `GCLOUD_SERVICE_ACCOUNT_EMAIL` | `%s@%s.iam.gserviceaccount.com` |\n", config.ServiceAccountID, config.ProjectID))

	return &NotificationMessage{
		Title:   "Workload Identity Federation セットアップ概要",
		Content: content.String(),
		Color:   "blue",
		IsCode:  false,
	}
}

// splitBashScript はBashスクリプトを2000文字制限に基づいて分割する
func (s *Service) splitBashScript(bashScript string) []*NotificationMessage {
	const maxLength = 1900 // 安全マージン100文字
	var messages []*NotificationMessage

	lines := strings.Split(bashScript, "\n")
	var currentChunk strings.Builder
	chunkNumber := 1

	currentChunk.WriteString("```bash\n")

	for _, line := range lines {
		// 次の行を追加した場合の長さをチェック
		testContent := currentChunk.String() + line + "\n```"
		if len(testContent) > maxLength && currentChunk.Len() > 10 { // 最低限の内容がある場合のみ分割
			// 現在のチunkを完了
			currentChunk.WriteString("```")

			title := fmt.Sprintf("実行用Bashスクリプト (Part %d)", chunkNumber)
			messages = append(messages, &NotificationMessage{
				Title:   title,
				Content: currentChunk.String(),
				Color:   "green",
				IsCode:  true,
			})

			// 新しいchunkを開始
			currentChunk.Reset()
			currentChunk.WriteString("```bash\n")
			chunkNumber++
		}

		currentChunk.WriteString(line + "\n")
	}

	// 最後のchunkを追加
	if currentChunk.Len() > 10 {
		currentChunk.WriteString("```")

		title := fmt.Sprintf("実行用Bashスクリプト (Part %d)", chunkNumber)
		if chunkNumber == 1 {
			title = "実行用Bashスクリプト"
		}

		messages = append(messages, &NotificationMessage{
			Title:   title,
			Content: currentChunk.String(),
			Color:   "green",
			IsCode:  true,
		})
	}

	return messages
}

// splitCleanupScript は削除スクリプトを2000文字制限に基づいて分割する
func (s *Service) splitCleanupScript(cleanupScript string) []*NotificationMessage {
	const maxLength = 1900 // 安全マージン100文字
	var messages []*NotificationMessage

	lines := strings.Split(cleanupScript, "\n")
	var currentChunk strings.Builder
	chunkNumber := 1

	currentChunk.WriteString("```bash\n")

	for _, line := range lines {
		// 次の行を追加した場合の長さをチェック
		testContent := currentChunk.String() + line + "\n```"
		if len(testContent) > maxLength && currentChunk.Len() > 10 { // 最低限の内容がある場合のみ分割
			// 現在のchunkを完了
			currentChunk.WriteString("```")

			title := fmt.Sprintf("削除用Bashスクリプト (Part %d)", chunkNumber)
			messages = append(messages, &NotificationMessage{
				Title:   title,
				Content: currentChunk.String(),
				Color:   "red",
				IsCode:  true,
			})

			// 新しいchunkを開始
			currentChunk.Reset()
			currentChunk.WriteString("```bash\n")
			chunkNumber++
		}

		currentChunk.WriteString(line + "\n")
	}

	// 最後のchunkを追加
	if currentChunk.Len() > 10 {
		currentChunk.WriteString("```")

		title := fmt.Sprintf("削除用Bashスクリプト (Part %d)", chunkNumber)
		if chunkNumber == 1 {
			title = "削除用Bashスクリプト"
		}

		messages = append(messages, &NotificationMessage{
			Title:   title,
			Content: currentChunk.String(),
			Color:   "red",
			IsCode:  true,
		})
	}

	return messages
}

// splitCleanupScriptWithIntro は削除スクリプトを説明部分の長さを考慮して分割する
func (s *Service) splitCleanupScriptWithIntro(cleanupScript string, introLength int) []*NotificationMessage {
	maxLength := 2000 - introLength - 100 // 説明部分の長さを差し引いて安全マージン100文字
	if maxLength < 500 {
		maxLength = 500 // 最低限の長さを確保
	}

	var messages []*NotificationMessage
	lines := strings.Split(cleanupScript, "\n")
	var currentChunk strings.Builder
	chunkNumber := 1

	currentChunk.WriteString("```bash\n")

	for _, line := range lines {
		// 次の行を追加した場合の長さをチェック
		testContent := currentChunk.String() + line + "\n```"
		if len(testContent) > maxLength && currentChunk.Len() > 10 { // 最低限の内容がある場合のみ分割
			// 現在のchunkを完了
			currentChunk.WriteString("```")

			title := fmt.Sprintf("削除用Bashスクリプト (Part %d)", chunkNumber)
			messages = append(messages, &NotificationMessage{
				Title:   title,
				Content: currentChunk.String(),
				Color:   "red",
				IsCode:  true,
			})

			// 新しいchunkを開始
			currentChunk.Reset()
			currentChunk.WriteString("```bash\n")
			chunkNumber++
		}

		currentChunk.WriteString(line + "\n")
	}

	// 最後のchunkを追加
	if currentChunk.Len() > 10 {
		currentChunk.WriteString("```")

		title := fmt.Sprintf("削除用Bashスクリプト (Part %d)", chunkNumber)
		if chunkNumber == 1 {
			title = "削除用Bashスクリプト"
		}

		messages = append(messages, &NotificationMessage{
			Title:   title,
			Content: currentChunk.String(),
			Color:   "red",
			IsCode:  true,
		})
	}

	return messages
}

// generateCleanupMessages は削除スクリプトのメッセージを生成し、分割する
func (s *Service) generateCleanupMessages(config *WorkloadIdentityConfig) ([]*NotificationMessage, error) {
	cleanupScript, err := s.GenerateCleanupScript(config)
	if err != nil {
		return nil, err
	}

	// 説明部分を準備
	var introContent strings.Builder
	introContent.WriteString("## リソース削除用スクリプト\n\n")
	introContent.WriteString("作成したリソースを削除する場合は、以下のスクリプトを実行してください:\n\n")
	introContent.WriteString("**注意**: このスクリプトは以下のリソースを削除します:\n")
	introContent.WriteString("- Workload Identity Pool\n")
	introContent.WriteString("- OIDC Provider\n")
	introContent.WriteString("- Service Account\n")
	introContent.WriteString("- IAM Policy Bindings (自動削除)\n\n")

	introText := introContent.String()
	introLength := len(introText)

	// 説明部分を考慮してより厳しい制限でスクリプトを分割
	scriptMessages := s.splitCleanupScriptWithIntro(cleanupScript, introLength)

	// 最初のメッセージに説明を追加
	if len(scriptMessages) > 0 {
		scriptMessages[0].Content = introText + scriptMessages[0].Content
	}

	return scriptMessages, nil
}

// generateCleanupMessage は削除スクリプトのメッセージを生成する（非推奨：分割版を使用）
func (s *Service) generateCleanupMessage(config *WorkloadIdentityConfig) (*NotificationMessage, error) {
	cleanupScript, err := s.GenerateCleanupScript(config)
	if err != nil {
		return nil, err
	}

	var content strings.Builder
	content.WriteString("## リソース削除用スクリプト\n\n")
	content.WriteString("作成したリソースを削除する場合は、以下のスクリプトを実行してください:\n\n")
	content.WriteString("```bash\n")
	content.WriteString(cleanupScript)
	content.WriteString("```\n\n")
	content.WriteString("**注意**: このスクリプトは以下のリソースを削除します:\n")
	content.WriteString("- Workload Identity Pool\n")
	content.WriteString("- OIDC Provider\n")
	content.WriteString("- Service Account\n")
	content.WriteString("- IAM Policy Bindings (自動削除)\n")

	return &NotificationMessage{
		Title:   "リソース削除用Bashスクリプト",
		Content: content.String(),
		Color:   "red",
		IsCode:  true,
	}, nil
}

// generateWorkflowMessage はGitHub Actions設定のメッセージを生成する
func (s *Service) generateWorkflowMessage(config *WorkloadIdentityConfig) *NotificationMessage {
	var content strings.Builder

	content.WriteString("## GitHub Actions ワークフロー設定\n\n")
	content.WriteString("以下の設定を `.github/workflows/` ディレクトリ内のYAMLファイルに追加してください:\n\n")
	content.WriteString("```yaml\n")
	content.WriteString("on:\n")
	content.WriteString("  push:\n\n")
	content.WriteString("permissions:\n")
	content.WriteString("  contents: write # リポジトリへの書き込み権限 (バッジ更新のため)\n")
	content.WriteString("  id-token: write # OIDCトークンをリクエストする権限 (Google Cloud認証のため)\n\n")
	content.WriteString("env:\n")
	content.WriteString("  GOOGLE_CLOUD_PROJECT_ID_01: 'any-project'\n")
	content.WriteString("  GCLOUD_PROJECT_NUMBER: ${{ secrets.GCLOUD_PROJECT_NUMBER }}\n")
	content.WriteString("  GCLOUD_POOL_ID: ${{ secrets.GCLOUD_POOL_ID }}\n")
	content.WriteString("  GCLOUD_PROVIDER_ID: ${{ secrets.GCLOUD_PROVIDER_ID }}\n")
	content.WriteString("  GCLOUD_SERVICE_ACCOUNT_EMAIL: ${{ secrets.GCLOUD_SERVICE_ACCOUNT_EMAIL }}\n\n")
	content.WriteString("jobs:\n")
	content.WriteString("  test:\n")
	content.WriteString("    runs-on: ubuntu-latest\n")
	content.WriteString("    if: contains(github.event.head_commit.message, '[skip ci]') == false\n")
	content.WriteString("    steps:\n")
	content.WriteString("      - name: Checkout\n")
	content.WriteString("        uses: actions/checkout@v4\n")
	content.WriteString("      - name: Set up Golang\n")
	content.WriteString("        uses: actions/setup-go@v5\n")
	content.WriteString("        with:\n")
	content.WriteString("          go-version-file: 'go.mod'\n")
	content.WriteString("      - id: 'gcloud_auth'\n")
	content.WriteString("        name: 'Authenticate to Google Cloud'\n")
	content.WriteString("        uses: 'google-github-actions/auth@v2'\n")
	content.WriteString("        with:\n")
	content.WriteString("          create_credentials_file: true\n")
	content.WriteString(fmt.Sprintf("          workload_identity_provider: 'projects/${{ env.GCLOUD_PROJECT_NUMBER }}/locations/%s/workloadIdentityPools/${{ env.GCLOUD_POOL_ID }}/providers/${{ env.GCLOUD_PROVIDER_ID }}'\n", config.Location))
	content.WriteString("          service_account: '${{ env.GCLOUD_SERVICE_ACCOUNT_EMAIL }}'\n")
	content.WriteString("      - name: 'Set up Cloud SDK'\n")
	content.WriteString("        uses: 'google-github-actions/setup-gcloud@v2'\n")
	content.WriteString("      - name: 'Use gcloud CLI'\n")
	content.WriteString("        run: 'gcloud info'\n")
	content.WriteString("```")

	return &NotificationMessage{
		Title:   "GitHub Actions ワークフロー設定",
		Content: content.String(),
		Color:   "purple",
		IsCode:  true,
	}
}

// validateConfig は設定の妥当性を検証する
func (s *Service) validateConfig(config *WorkloadIdentityConfig) error {
	if config.ProjectID == "" {
		return fmt.Errorf("project-idは必須です")
	}
	if config.PoolID == "" {
		return fmt.Errorf("pool-idは必須です")
	}
	if config.ProviderID == "" {
		return fmt.Errorf("provider-idは必須です")
	}
	if config.ServiceAccountID == "" {
		return fmt.Errorf("service-account-idは必須です")
	}
	if config.RepoOwner == "" {
		return fmt.Errorf("repo-ownerは必須です")
	}
	if config.RepoName == "" {
		return fmt.Errorf("repo-nameは必須です")
	}
	if config.Location == "" {
		config.Location = "global"
	}
	return nil
}
