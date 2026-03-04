package listmachinetypes

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/infrastructures"
	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params はマシンタイプ一覧処理に必要な値。
type Params struct {
	Zones          []string
	MinDiskSizeGiB int
	MaxDiskSizeGiB int
}

type machineType struct {
	Name                         string          `json:"name"`
	Zone                         string          `json:"zone"`
	GuestCPUs                    int             `json:"guestCpus"`
	MemoryMB                     int             `json:"memoryMb"`
	MaximumPersistentDisksSizeGB json.RawMessage `json:"maximumPersistentDisksSizeGb"`
}

// Service は list-machine-types operation の処理を担当する。
type Service struct {
	commandExecutor infrastructures.CommandExecutor
}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return newServiceWithCommandExecutor(infrastructures.NewOSCommandExecutor())
}

func newServiceWithCommandExecutor(commandExecutor infrastructures.CommandExecutor) *Service {
	if commandExecutor == nil {
		commandExecutor = infrastructures.NewOSCommandExecutor()
	}
	return &Service{commandExecutor: commandExecutor}
}

// Execute は gcloud を実行してマシンタイプ一覧を表形式で返す。
func (s *Service) Execute(params Params) (string, error) {
	if err := validateSizeRange(params.MinDiskSizeGiB, params.MaxDiskSizeGiB); err != nil {
		return "", err
	}

	args := []string{"compute", "machine-types", "list"}
	zones := normalizeZones(params.Zones)
	if len(zones) > 0 {
		args = append(args, "--zones="+strings.Join(zones, ","))
	}
	args = append(args, "--format=json(name,zone,guestCpus,memoryMb,maximumPersistentDisksSizeGb)")

	output, err := s.commandExecutor.Execute("gcloud", args...)
	if err != nil {
		return "", fmt.Errorf("gcloud コマンドの実行に失敗しました: %v\n出力: %s", err, strings.TrimSpace(string(output)))
	}

	var machineTypes []machineType
	if err := json.Unmarshal(output, &machineTypes); err != nil {
		return "", fmt.Errorf("gcloud の出力(JSON)解析に失敗しました: %w\n出力: %s", err, strings.TrimSpace(string(output)))
	}

	rows := make([][]string, 0, len(machineTypes))
	for _, item := range machineTypes {
		maxPersistentDiskSizeGB, err := parseMaximumPersistentDisksSizeGB(item.MaximumPersistentDisksSizeGB)
		if err != nil {
			return "", err
		}
		if !shouldIncludeMachineType(maxPersistentDiskSizeGB, params.MinDiskSizeGiB, params.MaxDiskSizeGiB) {
			continue
		}

		rows = append(rows, []string{
			item.Name,
			common.ZoneBasename(item.Zone),
			strconv.Itoa(item.GuestCPUs),
			strconv.Itoa(item.MemoryMB),
			strconv.FormatInt(maxPersistentDiskSizeGB, 10),
		})
	}

	return common.FormatTable(
		[]string{"NAME", "ZONE", "GUEST_CPUS", "MEMORY_MB", "MAX_PERSISTENT_DISKS_SIZE_GB"},
		rows,
	), nil
}

func validateSizeRange(minSizeGiB, maxSizeGiB int) error {
	if minSizeGiB < 0 {
		return fmt.Errorf("min-disk-size-gib は0以上で指定してください")
	}
	if maxSizeGiB < 0 {
		return fmt.Errorf("max-disk-size-gib は0以上で指定してください")
	}
	if minSizeGiB > 0 && maxSizeGiB > 0 && minSizeGiB > maxSizeGiB {
		return fmt.Errorf("min-disk-size-gib は max-disk-size-gib 以下で指定してください")
	}
	return nil
}

func normalizeZones(zones []string) []string {
	normalized := make([]string, 0, len(zones))
	for _, zone := range zones {
		trimmed := strings.TrimSpace(zone)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}
	return normalized
}

func shouldIncludeMachineType(maxPersistentDiskSizeGB int64, minFilter, maxFilter int) bool {
	switch {
	case minFilter > 0 && maxFilter > 0:
		return maxPersistentDiskSizeGB >= int64(minFilter) && maxPersistentDiskSizeGB <= int64(maxFilter)
	case minFilter > 0:
		return maxPersistentDiskSizeGB >= int64(minFilter)
	case maxFilter > 0:
		return maxPersistentDiskSizeGB <= int64(maxFilter)
	default:
		return true
	}
}

func parseMaximumPersistentDisksSizeGB(raw json.RawMessage) (int64, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return 0, nil
	}

	var asInt int64
	if err := json.Unmarshal(raw, &asInt); err == nil {
		return asInt, nil
	}

	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		value, parseErr := strconv.ParseInt(strings.TrimSpace(asString), 10, 64)
		if parseErr != nil {
			return 0, fmt.Errorf("maximumPersistentDisksSizeGb の解析に失敗しました: %w", parseErr)
		}
		return value, nil
	}

	return 0, fmt.Errorf("maximumPersistentDisksSizeGb の形式が不正です: %s", string(raw))
}
