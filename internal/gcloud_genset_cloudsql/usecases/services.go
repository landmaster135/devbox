package usecases

import (
	"fmt"
	"strings"
)

// Service は Cloud SQL 向け gcloud コマンドを生成する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

const (
	deletionModeEnable  = "enable"
	deletionModeDisable = "disable"

	activationPolicyAlways = "always"
	activationPolicyNever  = "never"
)

// DeleteInstanceParams はインスタンス削除に必要なパラメータ。
type DeleteInstanceParams struct {
	InstanceName string
}

// PatchDeletionProtectionParams は削除保護の更新に必要なパラメータ。
type PatchDeletionProtectionParams struct {
	InstanceName string
	Mode         string
}

// PatchActivationPolicyParams は起動ポリシー変更に必要なパラメータ。
type PatchActivationPolicyParams struct {
	InstanceName string
	Policy       string
}

// InstanceParams は単純にインスタンス名のみを必要とする操作用のパラメータ。
type InstanceParams struct {
	InstanceName string
}

// BuildDeleteInstanceCommand は安全に削除を行うための複合コマンドを生成する。
func (s *Service) BuildDeleteInstanceCommand(params DeleteInstanceParams) (string, error) {
	instance, err := normalizeInstanceName(params.InstanceName)
	if err != nil {
		return "", err
	}

	startCmd, err := s.BuildStartInstanceCommand(InstanceParams{InstanceName: instance})
	if err != nil {
		return "", err
	}

	disableProtectionCmd, err := s.BuildPatchDeletionProtectionCommand(PatchDeletionProtectionParams{
		InstanceName: instance,
		Mode:         deletionModeDisable,
	})
	if err != nil {
		return "", err
	}

	commands := []string{
		startCmd,
		disableProtectionCmd,
		fmt.Sprintf("gcloud sql instances delete %s", shellQuote(instance)),
	}

	return strings.Join(commands, " && \\\n"), nil
}

// BuildPatchDeletionProtectionCommand は削除保護の有効化・無効化コマンドを生成する。
func (s *Service) BuildPatchDeletionProtectionCommand(params PatchDeletionProtectionParams) (string, error) {
	instance, err := normalizeInstanceName(params.InstanceName)
	if err != nil {
		return "", err
	}

	mode := strings.ToLower(strings.TrimSpace(params.Mode))
	if mode == "" {
		return "", fmt.Errorf("deletionProtectionMode は必須です")
	}

	switch mode {
	case deletionModeEnable:
		return fmt.Sprintf("gcloud sql instances patch %s --deletion-protection", shellQuote(instance)), nil
	case deletionModeDisable:
		return fmt.Sprintf("gcloud sql instances patch %s --no-deletion-protection", shellQuote(instance)), nil
	default:
		return "", fmt.Errorf("deletionProtectionMode には enable または disable を指定してください")
	}
}

// BuildPatchActivationPolicyCommand は起動ポリシー変更コマンドを生成する。
func (s *Service) BuildPatchActivationPolicyCommand(params PatchActivationPolicyParams) (string, error) {
	instance, err := normalizeInstanceName(params.InstanceName)
	if err != nil {
		return "", err
	}

	policy := strings.ToLower(strings.TrimSpace(params.Policy))
	if policy == "" {
		return "", fmt.Errorf("activationPolicy は必須です")
	}

	switch policy {
	case activationPolicyAlways:
		return fmt.Sprintf("gcloud sql instances patch %s --activation-policy=ALWAYS", shellQuote(instance)), nil
	case activationPolicyNever:
		return fmt.Sprintf("gcloud sql instances patch %s --activation-policy=never", shellQuote(instance)), nil
	default:
		return "", fmt.Errorf("activationPolicy には always または never を指定してください")
	}
}

// BuildStartInstanceCommand は起動を保証するためのコマンドを生成する。
func (s *Service) BuildStartInstanceCommand(params InstanceParams) (string, error) {
	return s.BuildPatchActivationPolicyCommand(PatchActivationPolicyParams{
		InstanceName: params.InstanceName,
		Policy:       activationPolicyAlways,
	})
}

// BuildStopInstanceCommand は停止を保証するためのコマンドを生成する。
func (s *Service) BuildStopInstanceCommand(params InstanceParams) (string, error) {
	return s.BuildPatchActivationPolicyCommand(PatchActivationPolicyParams{
		InstanceName: params.InstanceName,
		Policy:       activationPolicyNever,
	})
}

// PrintHighlightedCommand は生成されたコマンドを装飾して出力する。
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}

func normalizeInstanceName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", fmt.Errorf("instanceName は必須です")
	}
	return trimmed, nil
}

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}
