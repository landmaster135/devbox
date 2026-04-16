package creategceinstanceandconfigure

import (
	"strings"

	addstartupscripttogceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/add_startup_script_to_gce_instance"
	creategceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_instance"
	setgceinstancemetadatafromyaml "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/set_gce_instance_metadata_from_yaml"
)

// Params はインスタンス作成 + metadata 設定 + スタートアップスクリプト登録コマンド生成に必要な値。
type Params struct {
	InstanceName      string
	Zone              string
	MachineType       string
	BootDiskSize      string
	BootDiskType      string
	MetadataYAMLPath  string
	StartupScriptPath string
}

type createGCEInstanceOperation interface {
	Build(params creategceinstance.Params) (string, error)
}

type setGCEInstanceMetadataFromYAMLOperation interface {
	Build(params setgceinstancemetadatafromyaml.Params) (string, error)
}

type addStartupScriptToGCEInstanceOperation interface {
	Build(params addstartupscripttogceinstance.Params) (string, error)
}

// Service は create-gce-instance-and-configure operation のコマンド生成を担当する。
type Service struct {
	createGCEInstanceOperation              createGCEInstanceOperation
	setGCEInstanceMetadataFromYAMLOperation setGCEInstanceMetadataFromYAMLOperation
	addStartupScriptToGCEInstanceOperation  addStartupScriptToGCEInstanceOperation
}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return newServiceWithOperations(
		creategceinstance.NewService(),
		setgceinstancemetadatafromyaml.NewService(),
		addstartupscripttogceinstance.NewService(),
	)
}

func newServiceWithOperations(
	createGCEInstanceOp createGCEInstanceOperation,
	setGCEInstanceMetadataFromYAMLOp setGCEInstanceMetadataFromYAMLOperation,
	addStartupScriptToGCEInstanceOp addStartupScriptToGCEInstanceOperation,
) *Service {
	return &Service{
		createGCEInstanceOperation:              createGCEInstanceOp,
		setGCEInstanceMetadataFromYAMLOperation: setGCEInstanceMetadataFromYAMLOp,
		addStartupScriptToGCEInstanceOperation:  addStartupScriptToGCEInstanceOp,
	}
}

// Build はインスタンス作成 + metadata 設定 + スタートアップスクリプト登録コマンドを生成する。
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

	setMetadataCommand, err := s.setGCEInstanceMetadataFromYAMLOperation.Build(setgceinstancemetadatafromyaml.Params{
		InstanceName:     params.InstanceName,
		Zone:             params.Zone,
		MetadataYAMLPath: params.MetadataYAMLPath,
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

	commands := []string{createCommand, setMetadataCommand, addStartupScriptCommand}
	return strings.Join(commands, " && \\\n"), nil
}
