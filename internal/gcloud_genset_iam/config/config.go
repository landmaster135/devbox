package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	OperationAddIamPolicyBindingToProject                           = "add-iam-policy-binding-to-project"
	OperationAddIamPolicyBindingToServiceAccount                    = "add-iam-policy-binding-to-service-account"
	OperationAddWorkloadIdentityBindingToServiceAccount             = "add-workload-identity-binding-to-service-account"
	OperationCreateServiceAccount                                   = "create-service-account"
	OperationListServiceAccounts                                    = "list-service-accounts"
	OperationDisableServiceAccount                                  = "disable-service-account"
	OperationEnableServiceAccount                                   = "enable-service-account"
	OperationDeleteServiceAccount                                   = "delete-service-account"
	OperationUndeleteServiceAccount                                 = "undelete-service-account"
	OperationUpdateServiceAccount                                   = "update-service-account"
	OperationDescribeServiceAccount                                 = "describe-service-account"
	OperationCreateWorkloadIdentityPool                             = "create-workload-identity-pool"
	OperationListWorkloadIdentityPools                              = "list-workload-identity-pools"
	OperationDescribeWorkloadIdentityPool                           = "describe-workload-identity-pool"
	OperationDeleteWorkloadIdentityPool                             = "delete-workload-identity-pool"
	OperationUndeleteWorkloadIdentityPool                           = "undelete-workload-identity-pool"
	OperationUpdateWorkloadIdentityPool                             = "update-workload-identity-pool"
	OperationCreateOidcWorkloadIdentityPoolProvider                 = "create-oidc-workload-identity-pool-provider"
	OperationCreateOidcWorkloadIdentityPoolProviderForGitHubActions = "create-oidc-workload-identity-pool-provider-for-github-actions"
	OperationListWorkloadIdentityPoolProviders                      = "list-workload-identity-pool-providers"
	OperationDescribeWorkloadIdentityPoolProvider                   = "describe-workload-identity-pool-provider"
	OperationUpdateOidcWorkloadIdentityPoolProvider                 = "update-oidc-workload-identity-pool-provider"
	OperationDeleteWorkloadIdentityPoolProvider                     = "delete-workload-identity-pool-provider"
	OperationUndeleteWorkloadIdentityPoolProvider                   = "undelete-workload-identity-pool-provider"
	OperationSetupWorkloadIdentityFederation                        = "setup-workload-identity-federation"
	OperationCleanupWorkloadIdentityFederation                      = "cleanup-workload-identity-federation"

	defaultLocation = "global"
)

// Config は CLI から受け取るパラメータを保持する。
type Config struct {
	Operation string
	Help      bool

	ProjectID       string
	ProjectNumber   string
	PoolID          string
	PoolDescription string
	ProviderID      string
	RepositoryOwner string
	RepositoryName  string

	ServiceAccountID    string
	ServiceAccountEmail string
	Role                string
	Member              string

	Condition         string
	ConditionFromFile string

	Filter string
	SortBy string

	Description string
	DisplayName string
	Location    string

	AttributeMapping   string
	AttributeCondition string
	IssuerURI          string
	AllowedAudiences   string
	JWKJSONPath        string

	ShowDeleted      bool
	Disabled         bool
	SkipConfirmation bool
}

// ParseFlags は標準の flag パーサーで引数を解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたパーサーを使って引数を解析する。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{
		Location: defaultLocation,
	}

	parser.StringVar(&cfg.Operation, "operation", cfg.Operation, "実行する操作 (必須)")
	parser.BoolVar(&cfg.Help, "help", cfg.Help, "ヘルプを表示する")
	parser.StringVar(&cfg.ProjectID, "project-id", cfg.ProjectID, "対象となるプロジェクトID")
	parser.StringVar(&cfg.ProjectNumber, "project-number", cfg.ProjectNumber, "対象となるプロジェクト番号")
	parser.StringVar(&cfg.PoolID, "pool-id", cfg.PoolID, "ワークロードアイデンティティプールID")
	parser.StringVar(&cfg.PoolDescription, "pool-description", cfg.PoolDescription, "ワークロードアイデンティティプールの説明")
	parser.StringVar(&cfg.ProviderID, "provider-id", cfg.ProviderID, "ワークロードアイデンティティプロバイダーID")
	parser.StringVar(&cfg.RepositoryOwner, "repository-owner", cfg.RepositoryOwner, "GitHubリポジトリのオーナー名")
	parser.StringVar(&cfg.RepositoryName, "repository-name", cfg.RepositoryName, "GitHubリポジトリ名")

	parser.StringVar(&cfg.ServiceAccountID, "service-account-id", cfg.ServiceAccountID, "サービスアカウントID")
	parser.StringVar(&cfg.ServiceAccountEmail, "service-account-email", cfg.ServiceAccountEmail, "サービスアカウントメールアドレス")
	parser.StringVar(&cfg.Role, "role", cfg.Role, "IAMロール")
	parser.StringVar(&cfg.Member, "member", cfg.Member, "IAMメンバー")

	parser.StringVar(&cfg.Condition, "condition", cfg.Condition, "--condition オプションに指定する式")
	parser.StringVar(&cfg.ConditionFromFile, "condition-from-file", cfg.ConditionFromFile, "--condition-from-file オプションに指定するファイルパス")

	parser.StringVar(&cfg.Filter, "filter", cfg.Filter, "gcloud の --filter オプション")
	parser.StringVar(&cfg.SortBy, "sort-by", cfg.SortBy, "gcloud の --sort-by オプション")

	parser.StringVar(&cfg.Description, "description", cfg.Description, "説明 (description) フィールド")
	parser.StringVar(&cfg.DisplayName, "display-name", cfg.DisplayName, "表示名 (display-name)")
	parser.StringVar(&cfg.Location, "location", cfg.Location, "リソースのロケーション (デフォルト: global)")

	parser.StringVar(&cfg.AttributeMapping, "attribute-mapping", cfg.AttributeMapping, "属性マッピング (KEY=VALUE,...) ")
	parser.StringVar(&cfg.AttributeCondition, "attribute-condition", cfg.AttributeCondition, "属性条件")
	parser.StringVar(&cfg.IssuerURI, "issuer-uri", cfg.IssuerURI, "OIDC 発行者 URI")
	parser.StringVar(&cfg.AllowedAudiences, "allowed-audiences", cfg.AllowedAudiences, "許可するオーディエンス (カンマ区切り)")
	parser.StringVar(&cfg.JWKJSONPath, "jwk-json-path", cfg.JWKJSONPath, "JWK JSON ファイルパス")

	parser.BoolVar(&cfg.ShowDeleted, "show-deleted", cfg.ShowDeleted, "削除済みリソースも表示する")
	parser.BoolVar(&cfg.Disabled, "disabled", cfg.Disabled, "リソースを無効化する (--disabled)")
	parser.BoolVar(&cfg.SkipConfirmation, "skip-confirmation", cfg.SkipConfirmation, "確認プロンプトをスキップする")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if cfg.Help {
		return cfg, nil
	}

	normalizeConfig(cfg)

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func normalizeConfig(cfg *Config) {
	cfg.Operation = strings.TrimSpace(cfg.Operation)
	cfg.ProjectID = strings.TrimSpace(cfg.ProjectID)
	cfg.ProjectNumber = strings.TrimSpace(cfg.ProjectNumber)
	cfg.PoolID = strings.TrimSpace(cfg.PoolID)
	cfg.PoolDescription = strings.TrimSpace(cfg.PoolDescription)
	cfg.ProviderID = strings.TrimSpace(cfg.ProviderID)
	cfg.RepositoryOwner = strings.TrimSpace(cfg.RepositoryOwner)
	cfg.RepositoryName = strings.TrimSpace(cfg.RepositoryName)
	cfg.ServiceAccountID = strings.TrimSpace(cfg.ServiceAccountID)
	cfg.ServiceAccountEmail = strings.TrimSpace(cfg.ServiceAccountEmail)
	cfg.Role = strings.TrimSpace(cfg.Role)
	cfg.Member = strings.TrimSpace(cfg.Member)
	cfg.Condition = strings.TrimSpace(cfg.Condition)
	cfg.ConditionFromFile = strings.TrimSpace(cfg.ConditionFromFile)
	cfg.Filter = strings.TrimSpace(cfg.Filter)
	cfg.SortBy = strings.TrimSpace(cfg.SortBy)
	cfg.Description = strings.TrimSpace(cfg.Description)
	cfg.DisplayName = strings.TrimSpace(cfg.DisplayName)
	cfg.Location = strings.TrimSpace(cfg.Location)
	cfg.AttributeMapping = strings.TrimSpace(cfg.AttributeMapping)
	cfg.AttributeCondition = strings.TrimSpace(cfg.AttributeCondition)
	cfg.IssuerURI = strings.TrimSpace(cfg.IssuerURI)
	cfg.AllowedAudiences = strings.TrimSpace(cfg.AllowedAudiences)
	cfg.JWKJSONPath = strings.TrimSpace(cfg.JWKJSONPath)

	if cfg.Location == "" {
		cfg.Location = defaultLocation
	}
}

func validateConfig(cfg *Config) error {
	if cfg.Operation == "" {
		return fmt.Errorf("operation パラメータは必須です")
	}

	requireServiceAccountEmail := func() error {
		if cfg.ServiceAccountEmail == "" {
			return fmt.Errorf("service-account-email パラメータは必須です")
		}
		return nil
	}

	requireServiceAccountID := func() error {
		if cfg.ServiceAccountID == "" {
			return fmt.Errorf("service-account-id パラメータは必須です")
		}
		return nil
	}

	requireProjectID := func() error {
		if cfg.ProjectID == "" {
			return fmt.Errorf("project-id パラメータは必須です")
		}
		return nil
	}

	requirePoolID := func() error {
		if cfg.PoolID == "" {
			return fmt.Errorf("pool-id パラメータは必須です")
		}
		return nil
	}

	requireProviderID := func() error {
		if cfg.ProviderID == "" {
			return fmt.Errorf("provider-id パラメータは必須です")
		}
		return nil
	}

	requireRepository := func() error {
		if cfg.RepositoryOwner == "" {
			return fmt.Errorf("repository-owner パラメータは必須です")
		}
		if cfg.RepositoryName == "" {
			return fmt.Errorf("repository-name パラメータは必須です")
		}
		return nil
	}

	requireRepositoryOwnerOnly := func() error {
		if cfg.RepositoryOwner == "" {
			return fmt.Errorf("repository-owner パラメータは必須です")
		}
		return nil
	}

	requireProjectNumber := func() error {
		if cfg.ProjectNumber == "" {
			return fmt.Errorf("project-number パラメータは必須です")
		}
		return nil
	}

	requireRole := func() error {
		if cfg.Role == "" {
			return fmt.Errorf("role パラメータは必須です")
		}
		return nil
	}

	requireMember := func() error {
		if cfg.Member == "" {
			return fmt.Errorf("member パラメータは必須です")
		}
		return nil
	}

	requireConditionExclusive := func() error {
		if cfg.Condition != "" && cfg.ConditionFromFile != "" {
			return fmt.Errorf("condition と condition-from-file は同時に指定できません")
		}
		return nil
	}

	requireLocation := func() error {
		if cfg.Location == "" {
			cfg.Location = defaultLocation
		}
		return nil
	}

	requireDescriptionOrDisplay := func() error {
		if cfg.Description == "" && cfg.DisplayName == "" && !cfg.Disabled {
			return fmt.Errorf("description, display-name, disabled のいずれかを指定してください")
		}
		return nil
	}

	switch cfg.Operation {
	case OperationAddIamPolicyBindingToProject:
		if err := requireProjectID(); err != nil {
			return err
		}
		if err := requireServiceAccountID(); err != nil {
			return err
		}
		if err := requireRole(); err != nil {
			return err
		}
	case OperationAddIamPolicyBindingToServiceAccount:
		if err := requireServiceAccountEmail(); err != nil {
			return err
		}
		if err := requireMember(); err != nil {
			return err
		}
		if err := requireRole(); err != nil {
			return err
		}
		if err := requireConditionExclusive(); err != nil {
			return err
		}
	case OperationAddWorkloadIdentityBindingToServiceAccount:
		if err := requireServiceAccountEmail(); err != nil {
			return err
		}
		if err := requireProjectNumber(); err != nil {
			return err
		}
		if err := requirePoolID(); err != nil {
			return err
		}
		if err := requireRepository(); err != nil {
			return err
		}
		if err := requireConditionExclusive(); err != nil {
			return err
		}
	case OperationCreateServiceAccount:
		if err := requireServiceAccountID(); err != nil {
			return err
		}
		if err := requireProjectID(); err != nil {
			return err
		}
		if err := requireRole(); err != nil {
			return err
		}
	case OperationListServiceAccounts:
		// no mandatory fields
	case OperationDisableServiceAccount, OperationEnableServiceAccount, OperationDeleteServiceAccount,
		OperationDescribeServiceAccount:
		if err := requireServiceAccountEmail(); err != nil {
			return err
		}
	case OperationUndeleteServiceAccount:
		if err := requireServiceAccountEmail(); err != nil {
			return err
		}
	case OperationUpdateServiceAccount:
		if err := requireServiceAccountEmail(); err != nil {
			return err
		}
		if cfg.Description == "" && cfg.DisplayName == "" {
			return fmt.Errorf("description もしくは display-name のいずれかを指定してください")
		}
	case OperationCreateWorkloadIdentityPool:
		if err := requirePoolID(); err != nil {
			return err
		}
		if err := requireProjectID(); err != nil {
			return err
		}
		if err := requireLocation(); err != nil {
			return err
		}
	case OperationListWorkloadIdentityPools:
		if err := requireProjectID(); err != nil {
			return err
		}
		if err := requireLocation(); err != nil {
			return err
		}
	case OperationDescribeWorkloadIdentityPool, OperationDeleteWorkloadIdentityPool,
		OperationUndeleteWorkloadIdentityPool, OperationUpdateWorkloadIdentityPool:
		if err := requirePoolID(); err != nil {
			return err
		}
		if err := requireProjectID(); err != nil {
			return err
		}
		if err := requireLocation(); err != nil {
			return err
		}
		if cfg.Operation == OperationUpdateWorkloadIdentityPool {
			if err := requireDescriptionOrDisplay(); err != nil {
				return err
			}
		}
	case OperationCreateOidcWorkloadIdentityPoolProvider:
		if err := requireProviderID(); err != nil {
			return err
		}
		if err := requireProjectID(); err != nil {
			return err
		}
		if err := requirePoolID(); err != nil {
			return err
		}
		if cfg.IssuerURI == "" {
			return fmt.Errorf("issuer-uri パラメータは必須です")
		}
		if cfg.AttributeMapping == "" {
			return fmt.Errorf("attribute-mapping パラメータは必須です")
		}
		if cfg.AttributeCondition == "" {
			return fmt.Errorf("attribute-condition パラメータは必須です")
		}
	case OperationCreateOidcWorkloadIdentityPoolProviderForGitHubActions:
		if err := requireProviderID(); err != nil {
			return err
		}
		if err := requireProjectID(); err != nil {
			return err
		}
		if err := requirePoolID(); err != nil {
			return err
		}
		if err := requireRepositoryOwnerOnly(); err != nil {
			return err
		}
	case OperationListWorkloadIdentityPoolProviders:
		if err := requireProjectID(); err != nil {
			return err
		}
		if err := requirePoolID(); err != nil {
			return err
		}
		if err := requireLocation(); err != nil {
			return err
		}
	case OperationDescribeWorkloadIdentityPoolProvider, OperationDeleteWorkloadIdentityPoolProvider,
		OperationUndeleteWorkloadIdentityPoolProvider, OperationUpdateOidcWorkloadIdentityPoolProvider:
		if err := requireProviderID(); err != nil {
			return err
		}
		if err := requireProjectID(); err != nil {
			return err
		}
		if err := requirePoolID(); err != nil {
			return err
		}
		if err := requireLocation(); err != nil {
			return err
		}
		if cfg.Operation == OperationUpdateOidcWorkloadIdentityPoolProvider {
			if cfg.AllowedAudiences == "" && cfg.AttributeCondition == "" && cfg.AttributeMapping == "" &&
				cfg.Description == "" && !cfg.Disabled && cfg.DisplayName == "" && cfg.IssuerURI == "" && cfg.JWKJSONPath == "" {
				return fmt.Errorf("update 対象のパラメータを少なくとも1つ指定してください")
			}
		}
	case OperationSetupWorkloadIdentityFederation:
		if err := requireProjectID(); err != nil {
			return err
		}
		if err := requirePoolID(); err != nil {
			return err
		}
		if err := requireProviderID(); err != nil {
			return err
		}
		if err := requireServiceAccountID(); err != nil {
			return err
		}
		if err := requireRepository(); err != nil {
			return err
		}
		if err := requireLocation(); err != nil {
			return err
		}
	case OperationCleanupWorkloadIdentityFederation:
		if err := requireProjectID(); err != nil {
			return err
		}
		if err := requirePoolID(); err != nil {
			return err
		}
		if err := requireProviderID(); err != nil {
			return err
		}
		if err := requireServiceAccountID(); err != nil {
			return err
		}
		if err := requireLocation(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("未対応の operation です: %s", cfg.Operation)
	}

	return nil
}

// PrintUsage は CLI の使用方法を表示する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Google Cloud IAM / Workload Identity 向け gcloud コマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通オプション:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n")
	fmt.Fprintf(os.Stderr, "        実行する操作 (必須)\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプを表示\n")
	fmt.Fprintf(os.Stderr, "  -location string\n")
	fmt.Fprintf(os.Stderr, "        リソースのロケーション (デフォルト: %s)\n\n", defaultLocation)

	fmt.Fprintf(os.Stderr, "主要なオペレーション:\n")
	fmt.Fprintf(os.Stderr, "  %s\n", OperationAddIamPolicyBindingToProject)
	fmt.Fprintf(os.Stderr, "  %s\n", OperationAddIamPolicyBindingToServiceAccount)
	fmt.Fprintf(os.Stderr, "  %s\n", OperationAddWorkloadIdentityBindingToServiceAccount)
	fmt.Fprintf(os.Stderr, "  %s\n", OperationCreateServiceAccount)
	fmt.Fprintf(os.Stderr, "  %s\n", OperationListServiceAccounts)
	fmt.Fprintf(os.Stderr, "  %s\n", OperationCreateWorkloadIdentityPool)
	fmt.Fprintf(os.Stderr, "  %s\n", OperationCreateOidcWorkloadIdentityPoolProvider)
	fmt.Fprintf(os.Stderr, "  %s\n", OperationSetupWorkloadIdentityFederation)
	fmt.Fprintf(os.Stderr, "  %s\n", OperationCleanupWorkloadIdentityFederation)
	fmt.Fprintf(os.Stderr, "  (その他の操作については README を参照してください)\n\n")

	fmt.Fprintf(os.Stderr, "例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -project-id=my-project -service-account-id=my-sa -role=roles/storage.admin\n", os.Args[0], OperationAddIamPolicyBindingToProject)
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -service-account-email=sa@project.iam.gserviceaccount.com -member=user:example@example.com -role=roles/iam.serviceAccountUser\n", os.Args[0], OperationAddIamPolicyBindingToServiceAccount)
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -project-id=my-project -pool-id=my-pool -provider-id=my-provider -service-account-id=gha -repository-owner=my-org -repository-name=my-repo\n", os.Args[0], OperationSetupWorkloadIdentityFederation)
}
