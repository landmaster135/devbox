package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_iam/config"
	usecases "github.com/landmaster135/devbox/internal/gcloud_genset_iam/usecases"
)

func main() {
	config, err := cfg.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		cfg.PrintUsage()
		os.Exit(1)
	}

	if config.Help {
		cfg.PrintUsage()
		return
	}

	service := usecases.NewService()

	var (
		command  string
		buildErr error
	)

	switch config.Operation {
	case cfg.OperationAddIamPolicyBindingToProject:
		command, buildErr = service.BuildAddIamPolicyBindingToProjectCommand(usecases.AddIamPolicyBindingToProjectParams{
			ProjectID:        config.ProjectID,
			ServiceAccountID: config.ServiceAccountID,
			Role:             config.Role,
		})
	case cfg.OperationAddIamPolicyBindingToServiceAccount:
		command, buildErr = service.BuildAddIamPolicyBindingToServiceAccountCommand(usecases.AddIamPolicyBindingToServiceAccountParams{
			ServiceAccountEmail: config.ServiceAccountEmail,
			Member:              config.Member,
			Role:                config.Role,
			Condition:           config.Condition,
			ConditionFromFile:   config.ConditionFromFile,
		})
	case cfg.OperationAddWorkloadIdentityBindingToServiceAccount:
		command, buildErr = service.BuildAddWorkloadIdentityBindingToServiceAccountCommand(usecases.AddWorkloadIdentityBindingToServiceAccountParams{
			ServiceAccountEmail: config.ServiceAccountEmail,
			ProjectNumber:       config.ProjectNumber,
			PoolID:              config.PoolID,
			RepositoryOwner:     config.RepositoryOwner,
			RepositoryName:      config.RepositoryName,
			ProviderID:          config.ProviderID,
			Condition:           config.Condition,
			ConditionFromFile:   config.ConditionFromFile,
		})
	case cfg.OperationCreateServiceAccount:
		command, buildErr = service.BuildCreateServiceAccountCommand(usecases.CreateServiceAccountParams{
			ServiceAccountID: config.ServiceAccountID,
			ProjectID:        config.ProjectID,
			Role:             config.Role,
		})
	case cfg.OperationListServiceAccounts:
		command, buildErr = service.BuildListServiceAccountsCommand(usecases.ListServiceAccountsParams{
			Filter: config.Filter,
			SortBy: config.SortBy,
		})
	case cfg.OperationDisableServiceAccount:
		command, buildErr = service.BuildDisableServiceAccountCommand(usecases.ServiceAccountEmailParams{ServiceAccountEmail: config.ServiceAccountEmail})
	case cfg.OperationEnableServiceAccount:
		command, buildErr = service.BuildEnableServiceAccountCommand(usecases.ServiceAccountEmailParams{ServiceAccountEmail: config.ServiceAccountEmail})
	case cfg.OperationDeleteServiceAccount:
		command, buildErr = service.BuildDeleteServiceAccountCommand(usecases.ServiceAccountEmailParams{ServiceAccountEmail: config.ServiceAccountEmail})
	case cfg.OperationUndeleteServiceAccount:
		command, buildErr = service.BuildUndeleteServiceAccountCommand(usecases.ServiceAccountEmailParams{ServiceAccountEmail: config.ServiceAccountEmail})
	case cfg.OperationUpdateServiceAccount:
		command, buildErr = service.BuildUpdateServiceAccountCommand(usecases.UpdateServiceAccountParams{
			ServiceAccountEmail: config.ServiceAccountEmail,
			Description:         config.Description,
			DisplayName:         config.DisplayName,
		})
	case cfg.OperationDescribeServiceAccount:
		command, buildErr = service.BuildDescribeServiceAccountCommand(usecases.ServiceAccountEmailParams{ServiceAccountEmail: config.ServiceAccountEmail})
	case cfg.OperationCreateWorkloadIdentityPool:
		command, buildErr = service.BuildCreateWorkloadIdentityPoolCommand(usecases.CreateWorkloadIdentityPoolParams{
			WorkloadIdentityPoolBaseParams: usecases.WorkloadIdentityPoolBaseParams{
				ProjectID: config.ProjectID,
				PoolID:    config.PoolID,
				Location:  config.Location,
			},
			Description: config.PoolDescription,
		})
	case cfg.OperationListWorkloadIdentityPools:
		command, buildErr = service.BuildListWorkloadIdentityPoolsCommand(usecases.ListWorkloadIdentityPoolsParams{
			ProjectID:   config.ProjectID,
			Location:    config.Location,
			ShowDeleted: config.ShowDeleted,
			Filter:      config.Filter,
		})
	case cfg.OperationDescribeWorkloadIdentityPool:
		command, buildErr = service.BuildDescribeWorkloadIdentityPoolCommand(usecases.WorkloadIdentityPoolBaseParams{
			ProjectID: config.ProjectID,
			PoolID:    config.PoolID,
			Location:  config.Location,
		})
	case cfg.OperationDeleteWorkloadIdentityPool:
		command, buildErr = service.BuildDeleteWorkloadIdentityPoolCommand(usecases.WorkloadIdentityPoolBaseParams{
			ProjectID: config.ProjectID,
			PoolID:    config.PoolID,
			Location:  config.Location,
		})
	case cfg.OperationUndeleteWorkloadIdentityPool:
		command, buildErr = service.BuildUndeleteWorkloadIdentityPoolCommand(usecases.WorkloadIdentityPoolBaseParams{
			ProjectID: config.ProjectID,
			PoolID:    config.PoolID,
			Location:  config.Location,
		})
	case cfg.OperationUpdateWorkloadIdentityPool:
		command, buildErr = service.BuildUpdateWorkloadIdentityPoolCommand(usecases.UpdateWorkloadIdentityPoolParams{
			WorkloadIdentityPoolBaseParams: usecases.WorkloadIdentityPoolBaseParams{
				ProjectID: config.ProjectID,
				PoolID:    config.PoolID,
				Location:  config.Location,
			},
			Description: config.Description,
			Disabled:    config.Disabled,
			DisplayName: config.DisplayName,
		})
	case cfg.OperationCreateOidcWorkloadIdentityPoolProvider:
		command, buildErr = service.BuildCreateOidcWorkloadIdentityPoolProviderCommand(usecases.CreateOidcWorkloadIdentityPoolProviderParams{
			WorkloadIdentityPoolProviderBaseParams: usecases.WorkloadIdentityPoolProviderBaseParams{
				ProjectID:  config.ProjectID,
				PoolID:     config.PoolID,
				ProviderID: config.ProviderID,
				Location:   config.Location,
			},
			IssuerURI:          config.IssuerURI,
			AttributeMapping:   config.AttributeMapping,
			AttributeCondition: config.AttributeCondition,
		})
	case cfg.OperationCreateOidcWorkloadIdentityPoolProviderForGitHubActions:
		command, buildErr = service.BuildCreateOidcWorkloadIdentityPoolProviderForGitHubActionsCommand(usecases.CreateOidcWorkloadIdentityPoolProviderForGitHubActionsParams{
			WorkloadIdentityPoolProviderBaseParams: usecases.WorkloadIdentityPoolProviderBaseParams{
				ProjectID:  config.ProjectID,
				PoolID:     config.PoolID,
				ProviderID: config.ProviderID,
				Location:   config.Location,
			},
			RepositoryOwner: config.RepositoryOwner,
		})
	case cfg.OperationListWorkloadIdentityPoolProviders:
		command, buildErr = service.BuildListWorkloadIdentityPoolProvidersCommand(usecases.ListWorkloadIdentityPoolProvidersParams{
			ProjectID:   config.ProjectID,
			PoolID:      config.PoolID,
			Location:    config.Location,
			ShowDeleted: config.ShowDeleted,
			Filter:      config.Filter,
		})
	case cfg.OperationDescribeWorkloadIdentityPoolProvider:
		command, buildErr = service.BuildDescribeWorkloadIdentityPoolProviderCommand(usecases.WorkloadIdentityPoolProviderBaseParams{
			ProjectID:  config.ProjectID,
			PoolID:     config.PoolID,
			ProviderID: config.ProviderID,
			Location:   config.Location,
		})
	case cfg.OperationDeleteWorkloadIdentityPoolProvider:
		command, buildErr = service.BuildDeleteWorkloadIdentityPoolProviderCommand(usecases.WorkloadIdentityPoolProviderBaseParams{
			ProjectID:  config.ProjectID,
			PoolID:     config.PoolID,
			ProviderID: config.ProviderID,
			Location:   config.Location,
		})
	case cfg.OperationUndeleteWorkloadIdentityPoolProvider:
		command, buildErr = service.BuildUndeleteWorkloadIdentityPoolProviderCommand(usecases.WorkloadIdentityPoolProviderBaseParams{
			ProjectID:  config.ProjectID,
			PoolID:     config.PoolID,
			ProviderID: config.ProviderID,
			Location:   config.Location,
		})
	case cfg.OperationUpdateOidcWorkloadIdentityPoolProvider:
		command, buildErr = service.BuildUpdateOidcWorkloadIdentityPoolProviderCommand(usecases.UpdateOidcWorkloadIdentityPoolProviderParams{
			WorkloadIdentityPoolProviderBaseParams: usecases.WorkloadIdentityPoolProviderBaseParams{
				ProjectID:  config.ProjectID,
				PoolID:     config.PoolID,
				ProviderID: config.ProviderID,
				Location:   config.Location,
			},
			AllowedAudiences:   config.AllowedAudiences,
			AttributeCondition: config.AttributeCondition,
			AttributeMapping:   config.AttributeMapping,
			Description:        config.Description,
			Disabled:           config.Disabled,
			DisplayName:        config.DisplayName,
			IssuerURI:          config.IssuerURI,
			JWKJSONPath:        config.JWKJSONPath,
		})
	case cfg.OperationSetupWorkloadIdentityFederation:
		command, buildErr = service.BuildSetupWorkloadIdentityFederationScript(usecases.SetupWorkloadIdentityFederationParams{
			ProjectID:        config.ProjectID,
			PoolID:           config.PoolID,
			ProviderID:       config.ProviderID,
			ServiceAccountID: config.ServiceAccountID,
			RepositoryOwner:  config.RepositoryOwner,
			RepositoryName:   config.RepositoryName,
			Location:         config.Location,
			PoolDescription:  config.PoolDescription,
		})
	case cfg.OperationCleanupWorkloadIdentityFederation:
		command, buildErr = service.BuildCleanupWorkloadIdentityFederationScript(usecases.CleanupWorkloadIdentityFederationParams{
			ProjectID:        config.ProjectID,
			PoolID:           config.PoolID,
			ProviderID:       config.ProviderID,
			ServiceAccountID: config.ServiceAccountID,
			Location:         config.Location,
			SkipConfirmation: config.SkipConfirmation,
		})
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作です: %s\n", config.Operation)
		cfg.PrintUsage()
		os.Exit(1)
	}

	if buildErr != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", buildErr)
		os.Exit(1)
	}

	service.PrintHighlightedCommand(command)
}
