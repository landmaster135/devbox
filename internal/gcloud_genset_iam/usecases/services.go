package usecases

import (
	"fmt"
	"strings"
)

// Service は IAM 関連の gcloud コマンドを生成する。
type Service struct{}

// NewService は Service を生成する。
func NewService() *Service {
	return &Service{}
}

// AddIamPolicyBindingToProjectParams はプロジェクトへのポリシーバインドに必要な情報。
type AddIamPolicyBindingToProjectParams struct {
	ProjectID        string
	ServiceAccountID string
	Role             string
}

// AddIamPolicyBindingToServiceAccountParams はサービスアカウントへのポリシーバインドに必要な情報。
type AddIamPolicyBindingToServiceAccountParams struct {
	ServiceAccountEmail string
	Member              string
	Role                string
	Condition           string
	ConditionFromFile   string
}

// AddWorkloadIdentityBindingToServiceAccountParams は Workload Identity バインドに必要な情報。
type AddWorkloadIdentityBindingToServiceAccountParams struct {
	ServiceAccountEmail string
	ProjectNumber       string
	PoolID              string
	RepositoryOwner     string
	RepositoryName      string
	ProviderID          string
	Condition           string
	ConditionFromFile   string
}

// CreateServiceAccountParams はサービスアカウント作成に必要な情報。
type CreateServiceAccountParams struct {
	ServiceAccountID string
	ProjectID        string
	Role             string
}

// ListServiceAccountsParams はサービスアカウント一覧取得のオプション。
type ListServiceAccountsParams struct {
	Filter string
	SortBy string
}

// ServiceAccountEmailParams は単一のサービスアカウントを扱う操作で使用する。
type ServiceAccountEmailParams struct {
	ServiceAccountEmail string
}

// UpdateServiceAccountParams はサービスアカウント更新のパラメータ。
type UpdateServiceAccountParams struct {
	ServiceAccountEmail string
	Description         string
	DisplayName         string
}

// WorkloadIdentityPoolBaseParams は Workload Identity Pool 操作用の共通パラメータ。
type WorkloadIdentityPoolBaseParams struct {
	ProjectID string
	PoolID    string
	Location  string
}

// CreateWorkloadIdentityPoolParams はプール作成のパラメータ。
type CreateWorkloadIdentityPoolParams struct {
	WorkloadIdentityPoolBaseParams
	Description string
}

// ListWorkloadIdentityPoolsParams はプール一覧取得のパラメータ。
type ListWorkloadIdentityPoolsParams struct {
	ProjectID   string
	Location    string
	ShowDeleted bool
	Filter      string
}

// UpdateWorkloadIdentityPoolParams はプール更新のパラメータ。
type UpdateWorkloadIdentityPoolParams struct {
	WorkloadIdentityPoolBaseParams
	Description string
	Disabled    bool
	DisplayName string
}

// WorkloadIdentityPoolProviderBaseParams はプロバイダー操作の共通パラメータ。
type WorkloadIdentityPoolProviderBaseParams struct {
	ProjectID  string
	PoolID     string
	ProviderID string
	Location   string
}

// CreateOidcWorkloadIdentityPoolProviderParams は OIDC プロバイダー作成のパラメータ。
type CreateOidcWorkloadIdentityPoolProviderParams struct {
	WorkloadIdentityPoolProviderBaseParams
	IssuerURI          string
	AttributeMapping   string
	AttributeCondition string
}

// CreateOidcWorkloadIdentityPoolProviderForGitHubActionsParams は GitHub Actions 向けプロバイダー作成のパラメータ。
type CreateOidcWorkloadIdentityPoolProviderForGitHubActionsParams struct {
	WorkloadIdentityPoolProviderBaseParams
	RepositoryOwner string
}

// ListWorkloadIdentityPoolProvidersParams はプロバイダー一覧のパラメータ。
type ListWorkloadIdentityPoolProvidersParams struct {
	ProjectID   string
	PoolID      string
	Location    string
	ShowDeleted bool
	Filter      string
}

// UpdateOidcWorkloadIdentityPoolProviderParams は OIDC プロバイダー更新のパラメータ。
type UpdateOidcWorkloadIdentityPoolProviderParams struct {
	WorkloadIdentityPoolProviderBaseParams
	AllowedAudiences   string
	AttributeCondition string
	AttributeMapping   string
	Description        string
	Disabled           bool
	DisplayName        string
	IssuerURI          string
	JWKJSONPath        string
}

// SetupWorkloadIdentityFederationParams はセットアップスクリプト生成のパラメータ。
type SetupWorkloadIdentityFederationParams struct {
	ProjectID        string
	PoolID           string
	ProviderID       string
	ServiceAccountID string
	RepositoryOwner  string
	RepositoryName   string
	Location         string
	PoolDescription  string
}

// CleanupWorkloadIdentityFederationParams はクリーンアップスクリプト生成のパラメータ。
type CleanupWorkloadIdentityFederationParams struct {
	ProjectID        string
	PoolID           string
	ProviderID       string
	ServiceAccountID string
	Location         string
	SkipConfirmation bool
}

// BuildAddIamPolicyBindingToProjectCommand は projects add-iam-policy-binding コマンドを生成する。
func (s *Service) BuildAddIamPolicyBindingToProjectCommand(params AddIamPolicyBindingToProjectParams) (string, error) {
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.ServiceAccountID == "" {
		return "", fmt.Errorf("service-account-id は必須です")
	}
	if params.Role == "" {
		return "", fmt.Errorf("role は必須です")
	}

	member := fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", params.ServiceAccountID, params.ProjectID)
	command := joinCommand(
		"gcloud",
		"projects",
		"add-iam-policy-binding",
		shellQuote(params.ProjectID),
		fmt.Sprintf("--member=%s", shellQuote(member)),
		fmt.Sprintf("--role=%s", shellQuote(params.Role)),
	)

	return command, nil
}

// BuildAddIamPolicyBindingToServiceAccountCommand は service-accounts add-iam-policy-binding を生成する。
func (s *Service) BuildAddIamPolicyBindingToServiceAccountCommand(params AddIamPolicyBindingToServiceAccountParams) (string, error) {
	if params.ServiceAccountEmail == "" {
		return "", fmt.Errorf("service-account-email は必須です")
	}
	if params.Member == "" {
		return "", fmt.Errorf("member は必須です")
	}
	if params.Role == "" {
		return "", fmt.Errorf("role は必須です")
	}
	if params.Condition != "" && params.ConditionFromFile != "" {
		return "", fmt.Errorf("condition と condition-from-file は同時に指定できません")
	}

	parts := []string{
		"gcloud",
		"iam",
		"service-accounts",
		"add-iam-policy-binding",
		shellQuote(params.ServiceAccountEmail),
		fmt.Sprintf("--member=%s", shellQuote(params.Member)),
		fmt.Sprintf("--role=%s", shellQuote(params.Role)),
	}

	if params.Condition != "" {
		parts = append(parts, fmt.Sprintf("--condition=%s", shellQuote(params.Condition)))
	}
	if params.ConditionFromFile != "" {
		parts = append(parts, fmt.Sprintf("--condition-from-file=%s", shellQuote(params.ConditionFromFile)))
	}

	return joinCommand(parts...), nil
}

// BuildAddWorkloadIdentityBindingToServiceAccountCommand は Workload Identity バインドコマンドを生成する。
func (s *Service) BuildAddWorkloadIdentityBindingToServiceAccountCommand(params AddWorkloadIdentityBindingToServiceAccountParams) (string, error) {
	if params.ServiceAccountEmail == "" {
		return "", fmt.Errorf("service-account-email は必須です")
	}
	if params.ProjectNumber == "" {
		return "", fmt.Errorf("project-number は必須です")
	}
	if params.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if params.RepositoryOwner == "" {
		return "", fmt.Errorf("repository-owner は必須です")
	}
	if params.RepositoryName == "" {
		return "", fmt.Errorf("repository-name は必須です")
	}
	if params.Condition != "" && params.ConditionFromFile != "" {
		return "", fmt.Errorf("condition と condition-from-file は同時に指定できません")
	}

	principalSet := fmt.Sprintf(
		"principalSet://iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/attribute.repository/%s/%s",
		params.ProjectNumber,
		params.PoolID,
		params.RepositoryOwner,
		params.RepositoryName,
	)

	parts := []string{
		"gcloud",
		"iam",
		"service-accounts",
		"add-iam-policy-binding",
		shellQuote(params.ServiceAccountEmail),
		fmt.Sprintf("--member=%s", shellQuote(principalSet)),
		fmt.Sprintf("--role=%s", shellQuote("roles/iam.workloadIdentityUser")),
	}

	if params.Condition != "" {
		parts = append(parts, fmt.Sprintf("--condition=%s", shellQuote(params.Condition)))
	}
	if params.ConditionFromFile != "" {
		parts = append(parts, fmt.Sprintf("--condition-from-file=%s", shellQuote(params.ConditionFromFile)))
	}

	return joinCommand(parts...), nil
}

// BuildCreateServiceAccountCommand はサービスアカウントを作成しロールを付与するコマンドを生成する。
func (s *Service) BuildCreateServiceAccountCommand(params CreateServiceAccountParams) (string, error) {
	if params.ServiceAccountID == "" {
		return "", fmt.Errorf("service-account-id は必須です")
	}
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.Role == "" {
		return "", fmt.Errorf("role は必須です")
	}

	createCmd := joinCommand(
		"gcloud",
		"iam",
		"service-accounts",
		"create",
		shellQuote(params.ServiceAccountID),
	)

	bindingCmd, err := s.BuildAddIamPolicyBindingToProjectCommand(AddIamPolicyBindingToProjectParams{
		ProjectID:        params.ProjectID,
		ServiceAccountID: params.ServiceAccountID,
		Role:             params.Role,
	})
	if err != nil {
		return "", err
	}

	return strings.Join([]string{createCmd, bindingCmd}, " &&\n"), nil
}

// BuildListServiceAccountsCommand はサービスアカウント一覧コマンドを生成する。
func (s *Service) BuildListServiceAccountsCommand(params ListServiceAccountsParams) (string, error) {
	parts := []string{"gcloud", "iam", "service-accounts", "list"}

	if params.Filter != "" {
		parts = append(parts, fmt.Sprintf("--filter=%s", shellQuote(params.Filter)))
	}
	if params.SortBy != "" {
		parts = append(parts, fmt.Sprintf("--sort-by=%s", shellQuote(params.SortBy)))
	}

	return joinCommand(parts...), nil
}

// BuildDisableServiceAccountCommand はサービスアカウントの無効化コマンドを生成する。
func (s *Service) BuildDisableServiceAccountCommand(params ServiceAccountEmailParams) (string, error) {
	if params.ServiceAccountEmail == "" {
		return "", fmt.Errorf("service-account-email は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"service-accounts",
		"disable",
		shellQuote(params.ServiceAccountEmail),
	), nil
}

// BuildEnableServiceAccountCommand はサービスアカウントの有効化コマンドを生成する。
func (s *Service) BuildEnableServiceAccountCommand(params ServiceAccountEmailParams) (string, error) {
	if params.ServiceAccountEmail == "" {
		return "", fmt.Errorf("service-account-email は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"service-accounts",
		"enable",
		shellQuote(params.ServiceAccountEmail),
	), nil
}

// BuildDeleteServiceAccountCommand はサービスアカウント削除コマンドを生成する。
func (s *Service) BuildDeleteServiceAccountCommand(params ServiceAccountEmailParams) (string, error) {
	if params.ServiceAccountEmail == "" {
		return "", fmt.Errorf("service-account-email は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"service-accounts",
		"delete",
		shellQuote(params.ServiceAccountEmail),
		"--quiet",
	), nil
}

// BuildUndeleteServiceAccountCommand はサービスアカウント復元コマンドを生成する。
func (s *Service) BuildUndeleteServiceAccountCommand(params ServiceAccountEmailParams) (string, error) {
	if params.ServiceAccountEmail == "" {
		return "", fmt.Errorf("service-account-email は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"service-accounts",
		"undelete",
		shellQuote(params.ServiceAccountEmail),
	), nil
}

// BuildUpdateServiceAccountCommand はサービスアカウント更新コマンドを生成する。
func (s *Service) BuildUpdateServiceAccountCommand(params UpdateServiceAccountParams) (string, error) {
	if params.ServiceAccountEmail == "" {
		return "", fmt.Errorf("service-account-email は必須です")
	}
	if params.Description == "" && params.DisplayName == "" {
		return "", fmt.Errorf("description もしくは display-name のいずれかを指定してください")
	}

	parts := []string{
		"gcloud",
		"iam",
		"service-accounts",
		"update",
		shellQuote(params.ServiceAccountEmail),
	}
	if params.Description != "" {
		parts = append(parts, fmt.Sprintf("--description=%s", shellQuote(params.Description)))
	}
	if params.DisplayName != "" {
		parts = append(parts, fmt.Sprintf("--display-name=%s", shellQuote(params.DisplayName)))
	}

	return joinCommand(parts...), nil
}

// BuildDescribeServiceAccountCommand はサービスアカウント詳細表示コマンドを生成する。
func (s *Service) BuildDescribeServiceAccountCommand(params ServiceAccountEmailParams) (string, error) {
	if params.ServiceAccountEmail == "" {
		return "", fmt.Errorf("service-account-email は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"service-accounts",
		"describe",
		shellQuote(params.ServiceAccountEmail),
	), nil
}

// BuildCreateWorkloadIdentityPoolCommand はプール作成コマンドを生成する。
func (s *Service) BuildCreateWorkloadIdentityPoolCommand(params CreateWorkloadIdentityPoolParams) (string, error) {
	base := params.WorkloadIdentityPoolBaseParams
	if base.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if base.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if base.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	parts := []string{
		"gcloud",
		"iam",
		"workload-identity-pools",
		"create",
		shellQuote(base.PoolID),
		fmt.Sprintf("--project=%s", shellQuote(base.ProjectID)),
		fmt.Sprintf("--location=%s", shellQuote(base.Location)),
	}
	if params.Description != "" {
		parts = append(parts, fmt.Sprintf("--description=%s", shellQuote(params.Description)))
	}

	return joinCommand(parts...), nil
}

// BuildListWorkloadIdentityPoolsCommand はプール一覧コマンドを生成する。
func (s *Service) BuildListWorkloadIdentityPoolsCommand(params ListWorkloadIdentityPoolsParams) (string, error) {
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	parts := []string{
		"gcloud",
		"iam",
		"workload-identity-pools",
		"list",
		fmt.Sprintf("--project=%s", shellQuote(params.ProjectID)),
		fmt.Sprintf("--location=%s", shellQuote(params.Location)),
	}
	if params.ShowDeleted {
		parts = append(parts, "--show-deleted")
	}
	if params.Filter != "" {
		parts = append(parts, fmt.Sprintf("--filter=%s", shellQuote(params.Filter)))
	}

	return joinCommand(parts...), nil
}

// BuildDescribeWorkloadIdentityPoolCommand はプール詳細表示コマンドを生成する。
func (s *Service) BuildDescribeWorkloadIdentityPoolCommand(params WorkloadIdentityPoolBaseParams) (string, error) {
	if params.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"workload-identity-pools",
		"describe",
		shellQuote(params.PoolID),
		fmt.Sprintf("--project=%s", shellQuote(params.ProjectID)),
		fmt.Sprintf("--location=%s", shellQuote(params.Location)),
	), nil
}

// BuildDeleteWorkloadIdentityPoolCommand はプール削除コマンドを生成する。
func (s *Service) BuildDeleteWorkloadIdentityPoolCommand(params WorkloadIdentityPoolBaseParams) (string, error) {
	if params.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"workload-identity-pools",
		"delete",
		shellQuote(params.PoolID),
		fmt.Sprintf("--project=%s", shellQuote(params.ProjectID)),
		fmt.Sprintf("--location=%s", shellQuote(params.Location)),
		"--quiet",
	), nil
}

// BuildUndeleteWorkloadIdentityPoolCommand はプール復元コマンドを生成する。
func (s *Service) BuildUndeleteWorkloadIdentityPoolCommand(params WorkloadIdentityPoolBaseParams) (string, error) {
	if params.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"workload-identity-pools",
		"undelete",
		shellQuote(params.PoolID),
		fmt.Sprintf("--project=%s", shellQuote(params.ProjectID)),
		fmt.Sprintf("--location=%s", shellQuote(params.Location)),
	), nil
}

// BuildUpdateWorkloadIdentityPoolCommand はプール更新コマンドを生成する。
func (s *Service) BuildUpdateWorkloadIdentityPoolCommand(params UpdateWorkloadIdentityPoolParams) (string, error) {
	base := params.WorkloadIdentityPoolBaseParams
	if base.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if base.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if base.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}
	if params.Description == "" && params.DisplayName == "" && !params.Disabled {
		return "", fmt.Errorf("description, display-name, disabled のいずれかを指定してください")
	}

	parts := []string{
		"gcloud",
		"iam",
		"workload-identity-pools",
		"update",
		shellQuote(base.PoolID),
		fmt.Sprintf("--project=%s", shellQuote(base.ProjectID)),
		fmt.Sprintf("--location=%s", shellQuote(base.Location)),
	}

	if params.Description != "" {
		parts = append(parts, fmt.Sprintf("--description=%s", shellQuote(params.Description)))
	}
	if params.Disabled {
		parts = append(parts, "--disabled")
	}
	if params.DisplayName != "" {
		parts = append(parts, fmt.Sprintf("--display-name=%s", shellQuote(params.DisplayName)))
	}

	return joinCommand(parts...), nil
}

// BuildCreateOidcWorkloadIdentityPoolProviderCommand は OIDC プロバイダー作成コマンドを生成する。
func (s *Service) BuildCreateOidcWorkloadIdentityPoolProviderCommand(params CreateOidcWorkloadIdentityPoolProviderParams) (string, error) {
	base := params.WorkloadIdentityPoolProviderBaseParams
	if base.ProviderID == "" {
		return "", fmt.Errorf("provider-id は必須です")
	}
	if base.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if base.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if base.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}
	if params.IssuerURI == "" {
		return "", fmt.Errorf("issuer-uri は必須です")
	}
	if params.AttributeMapping == "" {
		return "", fmt.Errorf("attribute-mapping は必須です")
	}
	if params.AttributeCondition == "" {
		return "", fmt.Errorf("attribute-condition は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"workload-identity-pools",
		"providers",
		"create-oidc",
		shellQuote(base.ProviderID),
		fmt.Sprintf("--project=%s", shellQuote(base.ProjectID)),
		fmt.Sprintf("--location=%s", shellQuote(base.Location)),
		fmt.Sprintf("--workload-identity-pool=%s", shellQuote(base.PoolID)),
		fmt.Sprintf("--issuer-uri=%s", shellQuote(params.IssuerURI)),
		fmt.Sprintf("--attribute-mapping=%s", shellQuote(params.AttributeMapping)),
		fmt.Sprintf("--attribute-condition=%s", shellQuote(params.AttributeCondition)),
	), nil
}

// BuildCreateOidcWorkloadIdentityPoolProviderForGitHubActionsCommand は GitHub Actions 向け OIDC プロバイダー作成コマンドを生成する。
func (s *Service) BuildCreateOidcWorkloadIdentityPoolProviderForGitHubActionsCommand(params CreateOidcWorkloadIdentityPoolProviderForGitHubActionsParams) (string, error) {
	base := params.WorkloadIdentityPoolProviderBaseParams
	if base.ProviderID == "" {
		return "", fmt.Errorf("provider-id は必須です")
	}
	if base.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if base.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if base.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}
	if params.RepositoryOwner == "" {
		return "", fmt.Errorf("repository-owner は必須です")
	}

	attributeMapping := "google.subject=assertion.sub,attribute.repository=assertion.repository"
	attributeCondition := fmt.Sprintf("assertion.repository_owner=='%s'", params.RepositoryOwner)

	return joinCommand(
		"gcloud",
		"iam",
		"workload-identity-pools",
		"providers",
		"create-oidc",
		shellQuote(base.ProviderID),
		fmt.Sprintf("--project=%s", shellQuote(base.ProjectID)),
		fmt.Sprintf("--location=%s", shellQuote(base.Location)),
		fmt.Sprintf("--workload-identity-pool=%s", shellQuote(base.PoolID)),
		fmt.Sprintf("--issuer-uri=%s", shellQuote("https://token.actions.githubusercontent.com/")),
		fmt.Sprintf("--attribute-mapping=%s", shellQuote(attributeMapping)),
		fmt.Sprintf("--attribute-condition=%s", shellQuote(attributeCondition)),
	), nil
}

// BuildListWorkloadIdentityPoolProvidersCommand はプロバイダー一覧コマンドを生成する。
func (s *Service) BuildListWorkloadIdentityPoolProvidersCommand(params ListWorkloadIdentityPoolProvidersParams) (string, error) {
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if params.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	parts := []string{
		"gcloud",
		"iam",
		"workload-identity-pools",
		"providers",
		"list",
		fmt.Sprintf("--project=%s", shellQuote(params.ProjectID)),
		fmt.Sprintf("--workload-identity-pool=%s", shellQuote(params.PoolID)),
		fmt.Sprintf("--location=%s", shellQuote(params.Location)),
	}
	if params.ShowDeleted {
		parts = append(parts, "--show-deleted")
	}
	if params.Filter != "" {
		parts = append(parts, fmt.Sprintf("--filter=%s", shellQuote(params.Filter)))
	}

	return joinCommand(parts...), nil
}

// BuildDescribeWorkloadIdentityPoolProviderCommand はプロバイダー詳細コマンドを生成する。
func (s *Service) BuildDescribeWorkloadIdentityPoolProviderCommand(params WorkloadIdentityPoolProviderBaseParams) (string, error) {
	if params.ProviderID == "" {
		return "", fmt.Errorf("provider-id は必須です")
	}
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if params.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"workload-identity-pools",
		"providers",
		"describe",
		shellQuote(params.ProviderID),
		fmt.Sprintf("--project=%s", shellQuote(params.ProjectID)),
		fmt.Sprintf("--workload-identity-pool=%s", shellQuote(params.PoolID)),
		fmt.Sprintf("--location=%s", shellQuote(params.Location)),
	), nil
}

// BuildDeleteWorkloadIdentityPoolProviderCommand はプロバイダー削除コマンドを生成する。
func (s *Service) BuildDeleteWorkloadIdentityPoolProviderCommand(params WorkloadIdentityPoolProviderBaseParams) (string, error) {
	if params.ProviderID == "" {
		return "", fmt.Errorf("provider-id は必須です")
	}
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if params.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"workload-identity-pools",
		"providers",
		"delete",
		shellQuote(params.ProviderID),
		fmt.Sprintf("--project=%s", shellQuote(params.ProjectID)),
		fmt.Sprintf("--workload-identity-pool=%s", shellQuote(params.PoolID)),
		fmt.Sprintf("--location=%s", shellQuote(params.Location)),
		"--quiet",
	), nil
}

// BuildUndeleteWorkloadIdentityPoolProviderCommand はプロバイダー復元コマンドを生成する。
func (s *Service) BuildUndeleteWorkloadIdentityPoolProviderCommand(params WorkloadIdentityPoolProviderBaseParams) (string, error) {
	if params.ProviderID == "" {
		return "", fmt.Errorf("provider-id は必須です")
	}
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if params.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	return joinCommand(
		"gcloud",
		"iam",
		"workload-identity-pools",
		"providers",
		"undelete",
		shellQuote(params.ProviderID),
		fmt.Sprintf("--project=%s", shellQuote(params.ProjectID)),
		fmt.Sprintf("--workload-identity-pool=%s", shellQuote(params.PoolID)),
		fmt.Sprintf("--location=%s", shellQuote(params.Location)),
	), nil
}

// BuildUpdateOidcWorkloadIdentityPoolProviderCommand は OIDC プロバイダー更新コマンドを生成する。
func (s *Service) BuildUpdateOidcWorkloadIdentityPoolProviderCommand(params UpdateOidcWorkloadIdentityPoolProviderParams) (string, error) {
	base := params.WorkloadIdentityPoolProviderBaseParams
	if base.ProviderID == "" {
		return "", fmt.Errorf("provider-id は必須です")
	}
	if base.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if base.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if base.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}
	if params.AllowedAudiences == "" && params.AttributeCondition == "" && params.AttributeMapping == "" &&
		params.Description == "" && !params.Disabled && params.DisplayName == "" && params.IssuerURI == "" && params.JWKJSONPath == "" {
		return "", fmt.Errorf("少なくとも1つの更新項目を指定してください")
	}

	parts := []string{
		"gcloud",
		"iam",
		"workload-identity-pools",
		"providers",
		"update-oidc",
		shellQuote(base.ProviderID),
		fmt.Sprintf("--project=%s", shellQuote(base.ProjectID)),
		fmt.Sprintf("--workload-identity-pool=%s", shellQuote(base.PoolID)),
		fmt.Sprintf("--location=%s", shellQuote(base.Location)),
	}

	if params.AllowedAudiences != "" {
		parts = append(parts, fmt.Sprintf("--allowed-audiences=%s", shellQuote(params.AllowedAudiences)))
	}
	if params.AttributeCondition != "" {
		parts = append(parts, fmt.Sprintf("--attribute-condition=%s", shellQuote(params.AttributeCondition)))
	}
	if params.AttributeMapping != "" {
		parts = append(parts, fmt.Sprintf("--attribute-mapping=%s", shellQuote(params.AttributeMapping)))
	}
	if params.Description != "" {
		parts = append(parts, fmt.Sprintf("--description=%s", shellQuote(params.Description)))
	}
	if params.Disabled {
		parts = append(parts, "--disabled")
	}
	if params.DisplayName != "" {
		parts = append(parts, fmt.Sprintf("--display-name=%s", shellQuote(params.DisplayName)))
	}
	if params.IssuerURI != "" {
		parts = append(parts, fmt.Sprintf("--issuer-uri=%s", shellQuote(params.IssuerURI)))
	}
	if params.JWKJSONPath != "" {
		parts = append(parts, fmt.Sprintf("--jwk-json-path=%s", shellQuote(params.JWKJSONPath)))
	}

	return joinCommand(parts...), nil
}

// BuildSetupWorkloadIdentityFederationScript は Workload Identity Federation を構築するスクリプトを生成する。
func (s *Service) BuildSetupWorkloadIdentityFederationScript(params SetupWorkloadIdentityFederationParams) (string, error) {
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if params.ProviderID == "" {
		return "", fmt.Errorf("provider-id は必須です")
	}
	if params.ServiceAccountID == "" {
		return "", fmt.Errorf("service-account-id は必須です")
	}
	if params.RepositoryOwner == "" {
		return "", fmt.Errorf("repository-owner は必須です")
	}
	if params.RepositoryName == "" {
		return "", fmt.Errorf("repository-name は必須です")
	}
	if params.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	serviceAccountEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", params.ServiceAccountID, params.ProjectID)
	attributeMapping := "google.subject=assertion.sub,attribute.repository=assertion.repository"
	attributeCondition := fmt.Sprintf("assertion.repository_owner=='%s'", params.RepositoryOwner)

	commands := []string{
		"set -e",
		fmt.Sprintf("gcloud iam workload-identity-pools create %s --project=%s --location=%s%s",
			shellQuote(params.PoolID),
			shellQuote(params.ProjectID),
			shellQuote(params.Location),
			optionalDescriptionFlag(params.PoolDescription),
		),
		fmt.Sprintf("gcloud iam workload-identity-pools providers create-oidc %s --project=%s --location=%s --workload-identity-pool=%s --issuer-uri=%s --attribute-mapping=%s --attribute-condition=%s",
			shellQuote(params.ProviderID),
			shellQuote(params.ProjectID),
			shellQuote(params.Location),
			shellQuote(params.PoolID),
			shellQuote("https://token.actions.githubusercontent.com/"),
			shellQuote(attributeMapping),
			shellQuote(attributeCondition),
		),
		fmt.Sprintf("gcloud iam service-accounts create %s",
			shellQuote(params.ServiceAccountID),
		),
		fmt.Sprintf("gcloud projects add-iam-policy-binding %s --member=%s --role=%s",
			shellQuote(params.ProjectID),
			shellQuote(fmt.Sprintf("serviceAccount:%s", serviceAccountEmail)),
			shellQuote("roles/monitoring.editor"),
		),
		fmt.Sprintf("gcloud projects add-iam-policy-binding %s --member=%s --role=%s",
			shellQuote(params.ProjectID),
			shellQuote(fmt.Sprintf("serviceAccount:%s", serviceAccountEmail)),
			shellQuote("roles/run.viewer"),
		),
		fmt.Sprintf("PROJECT_NUMBER=$(gcloud projects describe %s --format=value(projectNumber))",
			shellQuote(params.ProjectID),
		),
		fmt.Sprintf(
			"gcloud iam service-accounts add-iam-policy-binding %s --member=\"principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/%s/attribute.repository/%s/%s\" --role=%s",
			shellQuote(serviceAccountEmail),
			params.PoolID,
			params.RepositoryOwner,
			params.RepositoryName,
			shellQuote("roles/iam.workloadIdentityUser"),
		),
		"echo '--- GitHub Actions 設定例 ---'",
		"echo 'env:'",
		"echo \"  GCLOUD_PROJECT_NUMBER: ${PROJECT_NUMBER}\"",
		fmt.Sprintf("echo %s", shellQuote("  GCLOUD_POOL_ID: "+params.PoolID)),
		fmt.Sprintf("echo %s", shellQuote("  GCLOUD_PROVIDER_ID: "+params.ProviderID)),
		fmt.Sprintf("echo %s", shellQuote("  GCLOUD_SERVICE_ACCOUNT_EMAIL: "+serviceAccountEmail)),
	}

	return strings.Join(commands, "\n"), nil
}

// BuildCleanupWorkloadIdentityFederationScript は Workload Identity Federation のリソース削除スクリプトを生成する。
func (s *Service) BuildCleanupWorkloadIdentityFederationScript(params CleanupWorkloadIdentityFederationParams) (string, error) {
	if params.ProjectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}
	if params.PoolID == "" {
		return "", fmt.Errorf("pool-id は必須です")
	}
	if params.ProviderID == "" {
		return "", fmt.Errorf("provider-id は必須です")
	}
	if params.ServiceAccountID == "" {
		return "", fmt.Errorf("service-account-id は必須です")
	}
	if params.Location == "" {
		return "", fmt.Errorf("location は必須です")
	}

	serviceAccountEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", params.ServiceAccountID, params.ProjectID)

	commands := []string{}
	if !params.SkipConfirmation {
		commands = append(commands,
			"read -p 'この操作は Workload Identity Federation 関連リソースを削除します。続行しますか？ (y/N): ' -r reply",
			"if [[ ! $reply =~ ^[Yy]$ ]]; then",
			"  echo 'キャンセルしました。'",
			"  exit 0",
			"fi",
		)
	}

	commands = append(commands,
		"set -e",
		fmt.Sprintf("gcloud iam workload-identity-pools providers delete %s --project=%s --workload-identity-pool=%s --location=%s --quiet",
			shellQuote(params.ProviderID),
			shellQuote(params.ProjectID),
			shellQuote(params.PoolID),
			shellQuote(params.Location),
		),
		fmt.Sprintf("gcloud iam workload-identity-pools delete %s --project=%s --location=%s --quiet",
			shellQuote(params.PoolID),
			shellQuote(params.ProjectID),
			shellQuote(params.Location),
		),
		fmt.Sprintf("gcloud iam service-accounts delete %s --quiet",
			shellQuote(serviceAccountEmail),
		),
	)

	return strings.Join(commands, "\n"), nil
}

// PrintHighlightedCommand は生成したコマンドを装飾して表示する。
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println("\n==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}

func optionalDescriptionFlag(description string) string {
	if strings.TrimSpace(description) == "" {
		return ""
	}
	return fmt.Sprintf(" --description=%s", shellQuote(description))
}

func joinCommand(parts ...string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}
		filtered = append(filtered, part)
	}
	return strings.Join(filtered, " ")
}

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}
