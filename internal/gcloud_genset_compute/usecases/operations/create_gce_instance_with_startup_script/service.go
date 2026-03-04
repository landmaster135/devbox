package creategceinstancewithstartupscript

import (
	"strings"

	addstartupscripttogceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/add_startup_script_to_gce_instance"
	creategceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_instance"
)

// Params はインスタンス作成 + スタートアップスクリプト登録コマンド生成に必要な値。
type Params struct {
	InstanceName      string
	Zone              string
	MachineType       string
	BootDiskSize      string
	BootDiskType      string
	StartupScriptPath string
}

type createGCEInstanceOperation interface {
	Build(params creategceinstance.Params) (string, error)
}

type addStartupScriptToGCEInstanceOperation interface {
	Build(params addstartupscripttogceinstance.Params) (string, error)
}

// Service は create-gce-instance-with-startup-script operation のコマンド生成を担当する。
type Service struct {
	createGCEInstanceOperation             createGCEInstanceOperation
	addStartupScriptToGCEInstanceOperation addStartupScriptToGCEInstanceOperation
}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return newServiceWithOperations(
		creategceinstance.NewService(),
		addstartupscripttogceinstance.NewService(),
	)
}

func newServiceWithOperations(
	createGCEInstanceOp createGCEInstanceOperation,
	addStartupScriptToGCEInstanceOp addStartupScriptToGCEInstanceOperation,
) *Service {
	return &Service{
		createGCEInstanceOperation:             createGCEInstanceOp,
		addStartupScriptToGCEInstanceOperation: addStartupScriptToGCEInstanceOp,
	}
}

// Build はインスタンス作成 + スタートアップスクリプト登録コマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	createCommand, err := s.createGCEInstanceOperation.Build(creategceinstance.Params{
		InstanceName: params.InstanceName,
		Zone:         params.Zone,
		MachineType:  params.MachineType,
		BootDiskSize: params.BootDiskSize,
		BootDiskType: params.BootDiskType,
	})
	if err != nil {
		return "", err
	}

	addStartupScriptCommand, err := s.addStartupScriptToGCEInstanceOperation.Build(addstartupscripttogceinstance.Params{
		InstanceName:      params.InstanceName,
		Zone:              params.Zone,
		StartupScriptPath: params.StartupScriptPath,
	})
	if err != nil {
		return "", err
	}

	commands := []string{createCommand, addStartupScriptCommand}
	return strings.Join(commands, " && \\\n"), nil
}
