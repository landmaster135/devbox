package creategceinstance

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params はインスタンス作成コマンド生成に必要な値。
type Params struct {
	InstanceName string
	Zone         string
	MachineType  string
	BootDiskSize string
	BootDiskType string
}

// Service は create-gce-instance operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build はインスタンス作成コマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.InstanceName) {
		return "", fmt.Errorf("instance-name は必須です")
	}
	if common.IsBlank(params.Zone) {
		return "", fmt.Errorf("zone は必須です")
	}
	if common.IsBlank(params.MachineType) {
		return "", fmt.Errorf("machine-type は必須です")
	}
	if common.IsBlank(params.BootDiskSize) {
		return "", fmt.Errorf("boot-disk-size は必須です")
	}
	if common.IsBlank(params.BootDiskType) {
		return "", fmt.Errorf("boot-disk-type は必須です")
	}

	return fmt.Sprintf(
		"gcloud compute instances create %s --zone=%s --machine-type=%s --no-address --boot-disk-size=%s --boot-disk-type=%s",
		common.ShellQuote(params.InstanceName),
		common.ShellQuote(params.Zone),
		common.ShellQuote(params.MachineType),
		common.ShellQuote(params.BootDiskSize),
		common.ShellQuote(params.BootDiskType),
	), nil
}
