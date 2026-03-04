package listdisktypes

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/infrastructures"
	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

var validDiskSizePattern = regexp.MustCompile(`^(\d+)GB-(\d+)GB$`)

// Params はディスクタイプ一覧処理に必要な値。
type Params struct {
	Zones      []string
	MinSizeGiB int
	MaxSizeGiB int
}

type diskType struct {
	Name          string `json:"name"`
	Zone          string `json:"zone"`
	ValidDiskSize string `json:"validDiskSize"`
}

// Service は list-disk-types operation の処理を担当する。
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

// Execute は gcloud を実行してディスクタイプ一覧を表形式で返す。
func (s *Service) Execute(params Params) (string, error) {
	if err := validateSizeRange(params.MinSizeGiB, params.MaxSizeGiB); err != nil {
		return "", err
	}

	args := []string{"compute", "disk-types", "list"}
	zones := normalizeZones(params.Zones)
	if len(zones) > 0 {
		args = append(args, "--zones="+strings.Join(zones, ","))
	}
	args = append(args, "--format=json(name,zone,validDiskSize)")

	output, err := s.commandExecutor.Execute("gcloud", args...)
	if err != nil {
		return "", fmt.Errorf("gcloud コマンドの実行に失敗しました: %v\n出力: %s", err, strings.TrimSpace(string(output)))
	}

	var diskTypes []diskType
	if err := json.Unmarshal(output, &diskTypes); err != nil {
		return "", fmt.Errorf("gcloud の出力(JSON)解析に失敗しました: %w\n出力: %s", err, strings.TrimSpace(string(output)))
	}

	rows := make([][]string, 0, len(diskTypes))
	for _, item := range diskTypes {
		minSize, maxSize, err := parseValidDiskSize(item.ValidDiskSize)
		if err != nil {
			return "", err
		}
		if !shouldIncludeDiskType(minSize, maxSize, params.MinSizeGiB, params.MaxSizeGiB) {
			continue
		}

		rows = append(rows, []string{
			item.Name,
			common.ZoneBasename(item.Zone),
			item.ValidDiskSize,
		})
	}

	return common.FormatTable([]string{"NAME", "ZONE", "VALID_DISK_SIZES"}, rows), nil
}

func validateSizeRange(minSizeGiB, maxSizeGiB int) error {
	if minSizeGiB < 0 {
		return fmt.Errorf("min-size-gib は0以上で指定してください")
	}
	if maxSizeGiB < 0 {
		return fmt.Errorf("max-size-gib は0以上で指定してください")
	}
	if minSizeGiB > 0 && maxSizeGiB > 0 && minSizeGiB > maxSizeGiB {
		return fmt.Errorf("min-size-gib は max-size-gib 以下で指定してください")
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

func parseValidDiskSize(value string) (int, int, error) {
	matches := validDiskSizePattern.FindStringSubmatch(strings.TrimSpace(value))
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("validDiskSize の形式が不正です: %s", value)
	}

	minSize, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("validDiskSize(min) の解析に失敗しました: %w", err)
	}
	maxSize, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf("validDiskSize(max) の解析に失敗しました: %w", err)
	}
	return minSize, maxSize, nil
}

func shouldIncludeDiskType(minSize, maxSize, minFilter, maxFilter int) bool {
	switch {
	case minFilter > 0 && maxFilter > 0:
		return minSize >= minFilter && maxSize <= maxFilter
	case minFilter > 0:
		return minSize >= minFilter
	case maxFilter > 0:
		return maxSize <= maxFilter
	default:
		return true
	}
}
