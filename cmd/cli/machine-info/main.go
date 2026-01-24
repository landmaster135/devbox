package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// MachineInfo システム情報を格納する構造体
type MachineInfo struct {
	CPUName                string  `json:"cpu_name"`
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

// NetworkStats ネットワーク統計情報
type NetworkStats struct {
	SentBytes     uint64
	ReceivedBytes uint64
}

// getCPUInfo CPU情報を取得
func getCPUInfo() (name string, cores int, logicalProcessors int, maxClockSpeed float64, err error) {
	// CPU名を取得 (インデントされている可能性があるので trim を使用)
	cmd := exec.Command("sh", "-c", "lscpu | grep 'Model name:' | sed 's/.*Model name:[[:space:]]*//'")
	output, err := cmd.Output()
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("CPU名の取得に失敗: %w", err)
	}
	name = strings.TrimSpace(string(output))

	// 物理コア数を取得 (Core(s) per socket)
	cmd = exec.Command("sh", "-c", "lscpu | grep 'Core(s) per socket:' | sed 's/.*Core(s) per socket:[[:space:]]*//'")
	output, err = cmd.Output()
	if err != nil {
		return name, 0, 0, 0, fmt.Errorf("コア数の取得に失敗: %w", err)
	}
	coresPerSocket, _ := strconv.Atoi(strings.TrimSpace(string(output)))

	// ソケット数を取得
	cmd = exec.Command("sh", "-c", "lscpu | grep 'Socket(s):' | sed 's/.*Socket(s):[[:space:]]*//'")
	output, err = cmd.Output()
	if err != nil {
		return name, 0, 0, 0, fmt.Errorf("ソケット数の取得に失敗: %w", err)
	}
	sockets, _ := strconv.Atoi(strings.TrimSpace(string(output)))
	cores = coresPerSocket * sockets

	// 論理プロセッサ数を取得
	cmd = exec.Command("sh", "-c", "lscpu | grep '^CPU(s):' | sed 's/.*CPU(s):[[:space:]]*//'")
	output, err = cmd.Output()
	if err != nil {
		return name, cores, 0, 0, fmt.Errorf("論理プロセッサ数の取得に失敗: %w", err)
	}
	logicalProcessors, _ = strconv.Atoi(strings.TrimSpace(string(output)))

	// 最大クロック速度を取得 (MHz)
	cmd = exec.Command("sh", "-c", "lscpu | grep 'CPU max MHz:' | sed 's/.*CPU max MHz:[[:space:]]*//'")
	output, err = cmd.Output()
	if err == nil && len(strings.TrimSpace(string(output))) > 0 {
		maxClockSpeed, _ = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	} else {
		// CPU max MHzがない場合、CPU MHz:を試す
		cmd = exec.Command("sh", "-c", "lscpu | grep 'CPU MHz:' | sed 's/.*CPU MHz:[[:space:]]*//'")
		output, err = cmd.Output()
		if err == nil && len(strings.TrimSpace(string(output))) > 0 {
			maxClockSpeed, _ = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
		} else {
			// それでもダメなら /sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq を試す
			cmd = exec.Command("sh", "-c", "cat /sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq")
			output, err = cmd.Output()
			if err == nil {
				// kHz単位なのでMHzに変換
				freqKHz, _ := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
				maxClockSpeed = freqKHz / 1000
			} else {
				// それでも取得できない場合は -1 を設定（後で現在のクロック速度で補完）
				maxClockSpeed = -1
			}
		}
	}

	return name, cores, logicalProcessors, maxClockSpeed, nil
}

// getCurrentCPUClockSpeed 現在のCPUクロック速度を取得
func getCurrentCPUClockSpeed() (float64, error) {
	// まず /proc/cpuinfo から取得を試みる
	cmd := exec.Command("sh", "-c", "cat /proc/cpuinfo | grep 'cpu MHz' | head -1 | sed 's/.*:[[:space:]]*//'")
	output, err := cmd.Output()
	if err == nil && len(strings.TrimSpace(string(output))) > 0 {
		speed, _ := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
		if speed > 0 {
			return speed, nil
		}
	}

	// /proc/cpuinfo で取得できない場合は /sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq を試す
	cmd = exec.Command("sh", "-c", "cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
	output, err = cmd.Output()
	if err == nil {
		// kHz単位なのでMHzに変換
		freqKHz, _ := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
		return freqKHz / 1000, nil
	}

	return 0, fmt.Errorf("現在のCPUクロック速度の取得に失敗")
}

// getCPUTemperature CPU温度を取得（複数回サンプリングして平均）
func getCPUTemperature(samplings int, interval time.Duration) (float64, error) {
	var totalTemp float64
	validSamples := 0

	for i := 0; i < samplings; i++ {
		temp := 0.0
		found := false

		// 方法1: sensors コマンドでTctl温度を取得（AMD Ryzen用）
		cmd := exec.Command("sh", "-c", "sensors | grep 'Tctl:' | awk '{print $2}' | tr -d '+°C'")
		output, err := cmd.Output()
		if err == nil && len(strings.TrimSpace(string(output))) > 0 {
			temp, err = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
			if err == nil && temp > 0 {
				found = true
			}
		}

		// 方法2: sensors コマンドでTdie温度を取得（AMD Ryzen用）
		if !found {
			cmd = exec.Command("sh", "-c", "sensors | grep 'Tdie:' | awk '{print $2}' | tr -d '+°C'")
			output, err = cmd.Output()
			if err == nil && len(strings.TrimSpace(string(output))) > 0 {
				temp, err = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
				if err == nil && temp > 0 {
					found = true
				}
			}
		}

		// 方法3: sensors コマンドでPackage id 0を取得（Intel用）
		if !found {
			cmd = exec.Command("sh", "-c", "sensors | grep 'Package id 0:' | awk '{print $4}' | tr -d '+°C'")
			output, err = cmd.Output()
			if err == nil && len(strings.TrimSpace(string(output))) > 0 {
				temp, err = strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
				if err == nil && temp > 0 {
					found = true
				}
			}
		}

		// 方法4: /sys/class/hwmon を使用
		if !found {
			// hwmonデバイスを探索
			cmd = exec.Command("sh", "-c", "find /sys/class/hwmon -name 'temp*_label' -exec grep -l 'Tctl\\|Tdie\\|Package' {} \\; 2>/dev/null | head -1")
			output, err = cmd.Output()
			if err == nil && len(strings.TrimSpace(string(output))) > 0 {
				labelPath := strings.TrimSpace(string(output))
				// temp1_label -> temp1_input に変換
				inputPath := strings.Replace(labelPath, "_label", "_input", 1)
				cmd = exec.Command("cat", inputPath)
				output, err = cmd.Output()
				if err == nil {
					tempMilliC, _ := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
					temp = tempMilliC / 1000 // ミリ度からセ氏に変換
					found = true
				}
			}
		}

		// 方法5: /sys/class/thermal/thermal_zone0/temp を試す
		if !found {
			cmd = exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
			output, err = cmd.Output()
			if err == nil {
				tempMilliC, _ := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
				temp = tempMilliC / 1000 // ミリ度からセ氏に変換
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

// getMemoryInfo メモリ情報を取得
func getMemoryInfo() (totalMB float64, usageMB float64, err error) {
	cmd := exec.Command("sh", "-c", "free -m | grep Mem: | awk '{print $2, $3}'")
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

// getHostname ホスト名を取得
func getHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("ホスト名の取得に失敗: %w", err)
	}
	return hostname, nil
}

// getNetworkStats ネットワークインターフェースの統計情報を取得
func getNetworkStats(interfaceName string) (*NetworkStats, error) {
	// 送信バイト数
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cat /sys/class/net/%s/statistics/tx_bytes", interfaceName))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("送信バイト数の取得に失敗: %w", err)
	}
	sentBytes, _ := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)

	// 受信バイト数
	cmd = exec.Command("sh", "-c", fmt.Sprintf("cat /sys/class/net/%s/statistics/rx_bytes", interfaceName))
	output, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("受信バイト数の取得に失敗: %w", err)
	}
	receivedBytes, _ := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)

	return &NetworkStats{
		SentBytes:     sentBytes,
		ReceivedBytes: receivedBytes,
	}, nil
}

// getEthernetSpeeds ネットワーク速度を計測（複数回サンプリングして平均）
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

		// 差分を計算してKbpsに変換
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

// getMachineInfo 全てのマシン情報を取得
func getMachineInfo(networkInterface string) (*MachineInfo, error) {
	fmt.Println("マシン情報を取得中...")

	// CPU情報
	cpuName, cpuCores, cpuLogicalProcessors, cpuMaxClockSpeed, err := getCPUInfo()
	if err != nil {
		fmt.Printf("警告: %v\n", err)
	}
	fmt.Printf("CPU名: %s\n", cpuName)
	fmt.Printf("CPUコア数: %d\n", cpuCores)
	fmt.Printf("論理プロセッサ数: %d\n", cpuLogicalProcessors)
	fmt.Printf("CPU最大クロック速度: %.2f MHz\n", cpuMaxClockSpeed)

	// 現在のCPUクロック速度
	cpuCurrentClockSpeed, err := getCurrentCPUClockSpeed()
	if err != nil {
		fmt.Printf("警告: %v\n", err)
	}
	fmt.Printf("CPU現在のクロック速度: %.2f MHz\n", cpuCurrentClockSpeed)

	// 最大クロック速度が取得できなかった場合は現在のクロック速度を使用
	if cpuMaxClockSpeed <= 0 {
		cpuMaxClockSpeed = cpuCurrentClockSpeed
		fmt.Printf("警告: CPU最大クロック速度を取得できなかったため、現在のクロック速度を使用します\n")
	}

	// CPU温度（10回サンプリング、1秒間隔）
	cpuTemp, err := getCPUTemperature(10, 1*time.Second)
	if err != nil {
		fmt.Printf("警告: %v\n", err)
	}
	fmt.Printf("CPU温度: %.2f °C\n", cpuTemp)

	// メモリ情報
	memoryTotal, memoryUsage, err := getMemoryInfo()
	if err != nil {
		fmt.Printf("警告: %v\n", err)
	}
	fmt.Printf("メモリ総容量: %.2f MB\n", memoryTotal)
	fmt.Printf("メモリ使用量: %.2f MB\n", memoryUsage)

	// ホスト名
	hostname, err := getHostname()
	if err != nil {
		fmt.Printf("警告: %v\n", err)
	}
	fmt.Printf("ホスト名: %s\n", hostname)

	// ネットワーク速度（10回サンプリング、1秒間隔）
	fmt.Println("ネットワーク速度を計測中...")
	sentKbps, receivedKbps, err := getEthernetSpeeds(networkInterface, 10, 1*time.Second)
	if err != nil {
		fmt.Printf("警告: %v\n", err)
	}
	fmt.Printf("平均送信速度: %.2f Kbps\n", sentKbps)
	fmt.Printf("平均受信速度: %.2f Kbps\n", receivedKbps)

	return &MachineInfo{
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
	}, nil
}

func main() {
	// ネットワークインターフェース名（環境に応じて変更）
	networkInterface := "eth0" // または "ens33", "enp0s3" など
	if len(os.Args) > 1 {
		networkInterface = os.Args[1]
	}

	info, err := getMachineInfo(networkInterface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// JSON形式で出力
	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON変換エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n取得したシステム情報JSON:")
	fmt.Println(string(jsonData))

	// ファイルに保存
	filename := fmt.Sprintf("log_%s.json", time.Now().Format("20060102-150405"))
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ファイル書き込みエラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nログファイルに保存しました: %s\n", filename)
}
