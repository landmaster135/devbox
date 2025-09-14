package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// ValidatePort はポート番号の妥当性を検証する
func ValidatePort(port int) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("ポート番号は1-65535の範囲で指定してください: %d", port)
	}
	return nil
}

// IsPortAvailable は指定されたポートが利用可能かどうかを確認する
func IsPortAvailable(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

// ParseAddress はアドレス文字列からホストとポートを分離する
func ParseAddress(address string) (host string, port int, err error) {
	parts := strings.Split(address, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("無効なアドレス形式: %s", address)
	}

	host = parts[0]
	if host == "" {
		host = "localhost"
	}

	port, err = strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("無効なポート番号: %s", parts[1])
	}

	if err = ValidatePort(port); err != nil {
		return "", 0, err
	}

	return host, port, nil
}

// LogServerInfo はサーバー情報をログに出力する
func LogServerInfo(serverType string, address string, additionalInfo map[string]string) {
	log.Printf("=== %s サーバー情報 ===", serverType)
	log.Printf("アドレス: %s", address)

	if additionalInfo != nil {
		for key, value := range additionalInfo {
			log.Printf("%s: %s", key, value)
		}
	}

	log.Printf("========================")
}

// GetLocalIP はローカルIPアドレスを取得する
func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// FormatServerAddress はサーバーアドレスを整形する
func FormatServerAddress(host string, port int) string {
	if host == "" || host == "0.0.0.0" {
		if localIP, err := GetLocalIP(); err == nil {
			return fmt.Sprintf("%s:%d", localIP, port)
		}
		return fmt.Sprintf("localhost:%d", port)
	}
	return fmt.Sprintf("%s:%d", host, port)
}

// LogStartupMessage は起動メッセージをログに出力する
func LogStartupMessage(serverType string, port int) {
	log.Printf("🚀 %s サーバーを開始しています...", serverType)
	log.Printf("📡 ポート: %d", port)

	if localIP, err := GetLocalIP(); err == nil {
		log.Printf("🌐 ローカルアクセス: http://localhost:%d", port)
		log.Printf("🌍 ネットワークアクセス: http://%s:%d", localIP, port)
	}
}

// LogShutdownMessage はシャットダウンメッセージをログに出力する
func LogShutdownMessage(serverType string) {
	log.Printf("🛑 %s サーバーをシャットダウンしています...", serverType)
}
