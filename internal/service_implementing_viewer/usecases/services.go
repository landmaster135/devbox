package usecases

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// ServiceStatus はサービスの実装状況を表す構造体
type ServiceStatus struct {
	ServiceName string
	Directories map[string]bool // ディレクトリ名 -> 存在フラグ
}

// ServiceImplementingViewerService はサービス実装状況を確認するサービス
type ServiceImplementingViewerService struct {
	rootDir    string
	targetDirs []string
}

// NewServiceImplementingViewerService は新しいサービスを作成する
func NewServiceImplementingViewerService(rootDir string, targetDirs []string) *ServiceImplementingViewerService {
	return &ServiceImplementingViewerService{
		rootDir:    rootDir,
		targetDirs: targetDirs,
	}
}

// GetServiceImplementingStatus はサービス実装状況を取得する
func (s *ServiceImplementingViewerService) GetServiceImplementingStatus() (string, error) {
	// 各対象ディレクトリからサービス名を収集
	allServices := make(map[string]bool)
	servicesByDir := make(map[string][]string)

	for _, targetDir := range s.targetDirs {
		dirPath := filepath.Join(s.rootDir, targetDir)
		services, err := s.getServicesInDirectory(dirPath)
		if err != nil {
			return "", fmt.Errorf("ディレクトリ %s の読み取りに失敗しました: %v", dirPath, err)
		}

		servicesByDir[targetDir] = services
		for _, service := range services {
			normalizedName := s.normalizeServiceName(service)
			allServices[normalizedName] = true
		}
	}

	// サービス名をソート
	sortedServices := make([]string, 0, len(allServices))
	for serviceName := range allServices {
		sortedServices = append(sortedServices, serviceName)
	}
	sort.Strings(sortedServices)

	// 各サービスの実装状況を確認
	serviceStatuses := make([]ServiceStatus, 0, len(sortedServices))
	for _, serviceName := range sortedServices {
		status := ServiceStatus{
			ServiceName: serviceName,
			Directories: make(map[string]bool),
		}

		for _, targetDir := range s.targetDirs {
			status.Directories[targetDir] = s.isServiceImplementedInDirectory(serviceName, servicesByDir[targetDir])
		}

		serviceStatuses = append(serviceStatuses, status)
	}

	// 表形式で出力
	return s.formatAsTable(serviceStatuses), nil
}

// getServicesInDirectory は指定されたディレクトリ内のサービス名を取得する
func (s *ServiceImplementingViewerService) getServicesInDirectory(dirPath string) ([]string, error) {
	var services []string

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// ディレクトリが存在しない場合は空のスライスを返す
			if strings.Contains(err.Error(), "no such file or directory") {
				return filepath.SkipDir
			}
			return err
		}

		// ルートディレクトリ自体はスキップ
		if path == dirPath {
			return nil
		}

		// ディレクトリのみを対象とし、1階層のみ
		if d.IsDir() && filepath.Dir(path) == dirPath {
			services = append(services, d.Name())
		}

		return nil
	})

	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		return nil, err
	}

	return services, nil
}

// normalizeServiceName はサービス名を正規化する（「_」と「-」を統一）
func (s *ServiceImplementingViewerService) normalizeServiceName(serviceName string) string {
	return strings.ReplaceAll(serviceName, "_", "-")
}

// isServiceImplementedInDirectory は指定されたサービスがディレクトリに実装されているかチェック
func (s *ServiceImplementingViewerService) isServiceImplementedInDirectory(normalizedServiceName string, servicesInDir []string) bool {
	for _, service := range servicesInDir {
		if s.normalizeServiceName(service) == normalizedServiceName {
			return true
		}
	}
	return false
}

// formatAsTable は結果を表形式でフォーマットする
func (s *ServiceImplementingViewerService) formatAsTable(serviceStatuses []ServiceStatus) string {
	var result strings.Builder

	// ヘッダー行を作成
	result.WriteString("| service")
	defaultSerNameLen := 7 // "service"の長さ
	maxServiceNameLen := defaultSerNameLen
	for _, status := range serviceStatuses {
		if len(status.ServiceName) > maxServiceNameLen {
			maxServiceNameLen = len(status.ServiceName)
		}
	}

	// サービス名の列幅を調整
	for i := defaultSerNameLen; i < maxServiceNameLen; i++ {
		result.WriteString(" ")
	}

	for _, targetDir := range s.targetDirs {
		result.WriteString(" | ")
		result.WriteString(targetDir)
	}
	result.WriteString(" |\n")

	// セパレーター行を作成
	result.WriteString("| :")
	for i := 1; i < maxServiceNameLen - 1; i++ {
		result.WriteString("-")
	}
	result.WriteString(": ")

	for range s.targetDirs {
		result.WriteString("| :-: ")
	}
	result.WriteString("|\n")

	// データ行を作成
	for _, status := range serviceStatuses {
		result.WriteString("| ")
		result.WriteString(status.ServiceName)

		// サービス名の列幅を調整
		for i := len(status.ServiceName); i < maxServiceNameLen; i++ {
			result.WriteString(" ")
		}

		for _, targetDir := range s.targetDirs {
			result.WriteString(" | ")
			if status.Directories[targetDir] {
				result.WriteString("✅")
			} else {
				result.WriteString("❌️")
			}
			result.WriteString(" ")
		}
		result.WriteString("|\n")
	}

	return result.String()
}
