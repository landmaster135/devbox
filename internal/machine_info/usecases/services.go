package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	filesystem "github.com/landmaster135/devbox/internal/machine_info/infrastructure/filesystem"
)

const (
	SHELL = "sh"
)

// MachineInfo はシステム情報を保持する
type MachineInfo struct {
	CPUName                 string  `json:"cpu_name"`
	CPUCurrentClockSpeedMHz float64 `json:"cpu_current_clock_speed_mhz"`
	CPUMaxClockSpeedMHz     float64 `json:"cpu_max_clock_speed_mhz"`
	CPUCores                int     `json:"cpu_cores"`
	CPULogicalProcessors    int     `json:"cpu_logical_processors"`
	CPUTemperature          float64 `json:"cpu_temperature"`
	MemoryTotalMB           float64 `json:"memory_total_mb"`
	MemoryUsageMB           float64 `json:"memory_usage_mb"`
	PCHostname              string  `json:"pc_hostname"`
	EthernetAvgSentKbps     float64 `json:"ethernet_avg_sent_kbps"`
	EthernetAvgReceivedKbps float64 `json:"ethernet_avg_received_kbps"`
	OutputAt                int64   `json:"output_at"`
}

// NetworkStats はネットワーク統計値を保持する
type NetworkStats struct {
	SentBytes     uint64
	ReceivedBytes uint64
}

// MachineInfoResult は収集結果と警告をまとめた構造体
type MachineInfoResult struct {
	Info     *MachineInfo
	Warnings []string
}

// MachineInfoService はマシン情報収集に関するユースケースを実行する
type MachineInfoService struct {
	fileSystem filesystem.Repository
}

// MachineInfoUsecase はmachine-infoユースケースの共通インターフェース
type MachineInfoUsecase interface {
	CollectUbuntuInfo(networkInterface string) (*MachineInfoResult, error)
	SaveMachineInfoLog(info *MachineInfo, outputDir string) (string, string, error)
	CollectAndSaveUbuntuInfo(networkInterface, outputDir string) (*MachineInfoResult, string, string, error)
}

// NewMachineInfoService はMachineInfoServiceを生成する
func NewMachineInfoService() *MachineInfoService {
	return &MachineInfoService{
		fileSystem: filesystem.NewRepository(),
	}
}

// NewMachineInfoServiceWithDependencies はDI用のコンストラクタ
func NewMachineInfoServiceWithDependencies(fs filesystem.Repository) *MachineInfoService {
	if fs == nil {
		fs = filesystem.NewRepository()
	}
	return &MachineInfoService{fileSystem: fs}
}

// CollectUbuntuInfo はUbuntu環境向けにマシン情報を収集する
func (s *MachineInfoService) CollectUbuntuInfo(networkInterface string) (*MachineInfoResult, error) {
	if strings.TrimSpace(networkInterface) == "" {
		return nil, fmt.Errorf("ネットワークインターフェース名を指定してください")
	}

	warnings := make([]string, 0)

	cpuName, cpuCores, cpuLogicalProcessors, cpuMaxClockSpeed, err := getCPUInfo()
	if err != nil {
		warnings = append(warnings, err.Error())
	}

	cpuCurrentClockSpeed, err := getCurrentCPUClockSpeed()
	if err != nil {
		warnings = append(warnings, err.Error())
	}

	if cpuMaxClockSpeed <= 0 {
		if cpuCurrentClockSpeed > 0 {
			cpuMaxClockSpeed = cpuCurrentClockSpeed
			warnings = append(warnings, "CPU最大クロック速度を取得できなかったため、現在のクロック速度を使用しました")
		} else {
			warnings = append(warnings, "CPU最大クロック速度および現在のクロック速度を取得できませんでした")
		}
	}

	cpuTemp, err := getCPUTemperature(10, 1*time.Second)
	if err != nil {
		warnings = append(warnings, err.Error())
	}

	memoryTotal, memoryUsage, err := getMemoryInfo()
	if err != nil {
		warnings = append(warnings, err.Error())
	}

	hostname, err := getHostname()
	if err != nil {
		warnings = append(warnings, err.Error())
	}

	sentKbps, receivedKbps, err := getEthernetSpeeds(networkInterface, 10, 1*time.Second)
	if err != nil {
		warnings = append(warnings, err.Error())
	}

	info := &MachineInfo{
		CPUName:                 cpuName,
		CPUCurrentClockSpeedMHz: cpuCurrentClockSpeed,
		CPUMaxClockSpeedMHz:     cpuMaxClockSpeed,
		CPUCores:                cpuCores,
		CPULogicalProcessors:    cpuLogicalProcessors,
		CPUTemperature:          cpuTemp,
		MemoryTotalMB:           memoryTotal,
		MemoryUsageMB:           memoryUsage,
		PCHostname:              hostname,
		EthernetAvgSentKbps:     sentKbps,
		EthernetAvgReceivedKbps: receivedKbps,
		OutputAt:                time.Now().Unix(),
	}

	return &MachineInfoResult{Info: info, Warnings: warnings}, nil
}

// SaveMachineInfoLog は収集した情報をJSONとして保存し、表示用JSON文字列と保存先パスを返す
func (s *MachineInfoService) SaveMachineInfoLog(info *MachineInfo, outputDir string) (string, string, error) {
	if info == nil {
		return "", "", fmt.Errorf("保存対象のマシン情報がありません")
	}

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("JSON変換エラー: %w", err)
	}

	trimmedOutputDir := strings.TrimSpace(outputDir)
	if trimmedOutputDir == "" {
		trimmedOutputDir = "."
	}

	if err := s.fileSystem.EnsureDir(trimmedOutputDir, 0o755); err != nil {
		return "", "", err
	}

	filename := fmt.Sprintf("log_%s.json", time.Now().Format("20060102-150405"))
	outputPath := s.fileSystem.Join(trimmedOutputDir, filename)
	if err := s.fileSystem.WriteFile(outputPath, jsonData, 0o644); err != nil {
		return "", "", err
	}

	return string(jsonData), outputPath, nil
}

// CollectAndSaveUbuntuInfo はマシン情報の収集とログ保存を一度に行う
func (s *MachineInfoService) CollectAndSaveUbuntuInfo(networkInterface, outputDir string) (*MachineInfoResult, string, string, error) {
	return collectAndSave(s, networkInterface, outputDir)
}

func collectAndSave(usecase MachineInfoUsecase, networkInterface, outputDir string) (*MachineInfoResult, string, string, error) {
	result, err := usecase.CollectUbuntuInfo(networkInterface)
	if err != nil {
		return nil, "", "", err
	}
	if result == nil || result.Info == nil {
		return nil, "", "", fmt.Errorf("マシン情報の取得に失敗しました")
	}

	jsonText, outputPath, err := usecase.SaveMachineInfoLog(result.Info, outputDir)
	if err != nil {
		return nil, "", "", err
	}

	return result, jsonText, outputPath, nil
}

func getCPUInfo() (name string, cores int, logicalProcessors int, maxClockSpeed float64, err error) {
	cmd := exec.Command(SHELL, "-c", "lscpu | grep 'Model name:' | sed 's/.*Model name:[[:space:]]*//'")
	output, err := cmd.Output()
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("CPU名の取得に失敗: %w", err)
	}
	name = strings.TrimSpace(string(output))

	cmd = exec.Command(SHELL, "-c", "lscpu | grep 'Core(s) per socket:' | sed 's/.*Core(s) per socket:[[:space:]]*//'")
	output, err = cmd.Output()
	if err != nil {
		return name, 0, 0, 0, fmt.Errorf("コア数の取得に失敗: %w", err)
	}
	coresPerSocket, _ := strconv.Atoi(strings.TrimSpace(string(output)))

	cmd = exec.Command(SHELL, "-c", "lscpu | grep 'Socket(s):' | sed 's/.*Socket(s):[[:space:]]*//'")
	output, err = cmd.Output()
	if err != nil {
		return name, 0, 0, 0, fmt.Errorf("ソケット数の取得に失敗: %w", err)
	}
	sockets, _ := strconv.Atoi(strings.TrimSpace(string(output)))
	cores = coresPerSocket * sockets

	cmd = exec.Command(SHELL, "-c", "lscpu | grep '^CPU(s):' | sed 's/.*CPU(s):[[:space:]]*//'")
	output, err = cmd.Output()
	if err != nil {
		return name, cores, 0, 0, fmt.Errorf("論理プロセッサ数の取得に失敗: %w", err)
	}
	logicalProcessors, _ = strconv.Atoi(strings.TrimSpace(string(output)))

	cmd = exec.Command(SHELL, "-c", "lscpu | grep 'CPU max MHz:' | sed 's/.*CPU max MHz:[[:space:]]*//'")
	output, err = cmd.Output()
	if err == nil && len(strings.TrimSpace(string(output))) > 0 {
		maxClockSpeed, _ = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	} else {
		cmd = exec.Command(SHELL, "-c", "lscpu | grep 'CPU MHz:' | sed 's/.*CPU MHz:[[:space:]]*//'")
		output, err = cmd.Output()
		if err == nil && len(strings.TrimSpace(string(output))) > 0 {
			maxClockSpeed, _ = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
		} else {
			cmd = exec.Command(SHELL, "-c", "cat /sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq")
			output, err = cmd.Output()
			if err == nil {
				freqKHz, _ := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
				maxClockSpeed = freqKHz / 1000
			} else {
				maxClockSpeed = -1
			}
		}
	}

	return name, cores, logicalProcessors, maxClockSpeed, nil
}

func getCurrentCPUClockSpeed() (float64, error) {
	cmd := exec.Command(SHELL, "-c", "cat /proc/cpuinfo | grep 'cpu MHz' | head -1 | sed 's/.*:[[:space:]]*//'")
	output, err := cmd.Output()
	if err == nil && len(strings.TrimSpace(string(output))) > 0 {
		speed, _ := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
		if speed > 0 {
			return speed, nil
		}
	}

	cmd = exec.Command(SHELL, "-c", "cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
	output, err = cmd.Output()
	if err == nil {
		freqKHz, _ := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
		return freqKHz / 1000, nil
	}

	return 0, fmt.Errorf("現在のCPUクロック速度の取得に失敗")
}

func getCPUTemperature(samplings int, interval time.Duration) (float64, error) {
	var totalTemp float64
	validSamples := 0

	for i := 0; i < samplings; i++ {
		temp := 0.0
		found := false

		cmd := exec.Command(SHELL, "-c", "sensors | grep 'Tctl:' | awk '{print $2}' | tr -d '+°C'")
		output, err := cmd.Output()
		if err == nil && len(strings.TrimSpace(string(output))) > 0 {
			temp, err = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
			if err == nil && temp > 0 {
				found = true
			}
		}

		if !found {
			cmd = exec.Command(SHELL, "-c", "sensors | grep 'Tdie:' | awk '{print $2}' | tr -d '+°C'")
			output, err = cmd.Output()
			if err == nil && len(strings.TrimSpace(string(output))) > 0 {
				temp, err = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
				if err == nil && temp > 0 {
					found = true
				}
			}
		}

		if !found {
			cmd = exec.Command(SHELL, "-c", "sensors | grep 'Package id 0:' | awk '{print $4}' | tr -d '+°C'")
			output, err = cmd.Output()
			if err == nil && len(strings.TrimSpace(string(output))) > 0 {
				temp, err = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
				if err == nil && temp > 0 {
					found = true
				}
			}
		}

		if !found {
			cmd = exec.Command(SHELL, "-c", "find /sys/class/hwmon -name 'temp*_label' -exec grep -l 'Tctl\\|Tdie\\|Package' {} \\; 2>/dev/null | head -1")
			output, err = cmd.Output()
			if err == nil && len(strings.TrimSpace(string(output))) > 0 {
				labelPath := strings.TrimSpace(string(output))
				inputPath := strings.Replace(labelPath, "_label", "_input", 1)
				cmd = exec.Command("cat", inputPath)
				output, err = cmd.Output()
				if err == nil {
					tempMilliC, _ := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
					temp = tempMilliC / 1000
					found = true
				}
			}
		}

		if !found {
			cmd = exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
			output, err = cmd.Output()
			if err == nil {
				tempMilliC, _ := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
				temp = tempMilliC / 1000
				found = true
			}
		}

		if found {
			totalTemp += temp
			validSamples++
		}

		if i < samplings-1 {
			time.Sleep(interval)
		}
	}

	if validSamples == 0 {
		return 0, fmt.Errorf("CPU温度の取得に失敗しました。sensors コマンドがインストールされているか確認してください (sudo apt update && sudo apt install lm-sensors && sudo sensors-detect)")
	}

	return totalTemp / float64(validSamples), nil
}

func getMemoryInfo() (totalMB float64, usageMB float64, err error) {
	cmd := exec.Command(SHELL, "-c", "free -m | grep Mem: | awk '{print $2, $3}'")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("メモリ情報の取得に失敗: %w", err)
	}

	parts := strings.Fields(strings.TrimSpace(string(output)))
	if len(parts) >= 2 {
		totalMB, _ = strconv.ParseFloat(parts[0], 64)
		usageMB, _ = strconv.ParseFloat(parts[1], 64)
	}

	return totalMB, usageMB, nil
}

func getHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("ホスト名の取得に失敗: %w", err)
	}
	return hostname, nil
}

func getNetworkStats(interfaceName string) (*NetworkStats, error) {
	cmd := exec.Command(SHELL, "-c", fmt.Sprintf("cat /sys/class/net/%s/statistics/tx_bytes", interfaceName))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("送信バイト数の取得に失敗: %w", err)
	}
	sentBytes, _ := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)

	cmd = exec.Command(SHELL, "-c", fmt.Sprintf("cat /sys/class/net/%s/statistics/rx_bytes", interfaceName))
	output, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("受信バイト数の取得に失敗: %w", err)
	}
	receivedBytes, _ := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)

	return &NetworkStats{SentBytes: sentBytes, ReceivedBytes: receivedBytes}, nil
}

func getEthernetSpeeds(interfaceName string, samplings int, interval time.Duration) (sentKbps float64, receivedKbps float64, err error) {
	var totalSentKbps, totalReceivedKbps float64
	nonZeroSentCount := 0
	nonZeroReceivedCount := 0

	for i := 0; i < samplings; i++ {
		startStats, err := getNetworkStats(interfaceName)
		if err != nil {
			return 0, 0, err
		}

		time.Sleep(interval)

		endStats, err := getNetworkStats(interfaceName)
		if err != nil {
			return 0, 0, err
		}

		sentDiff := endStats.SentBytes - startStats.SentBytes
		receivedDiff := endStats.ReceivedBytes - startStats.ReceivedBytes

		sentKbpsValue := float64(sentDiff*8) / 1000 / interval.Seconds()
		receivedKbpsValue := float64(receivedDiff*8) / 1000 / interval.Seconds()

		if sentKbpsValue > 0 {
			totalSentKbps += sentKbpsValue
			nonZeroSentCount++
		}
		if receivedKbpsValue > 0 {
			totalReceivedKbps += receivedKbpsValue
			nonZeroReceivedCount++
		}
	}

	if nonZeroSentCount > 0 {
		sentKbps = totalSentKbps / float64(nonZeroSentCount)
	}
	if nonZeroReceivedCount > 0 {
		receivedKbps = totalReceivedKbps / float64(nonZeroReceivedCount)
	}

	return sentKbps, receivedKbps, nil
}
