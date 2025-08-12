package utils

import (
	"fmt"
	"net"
	"strconv"
)

// FindAvailablePort は利用可能なポートを検索する
func FindAvailablePort(preferredPort int) (int, error) {
	// 優先ポートが指定されている場合は、まずそれを試す
	if preferredPort > 0 {
		if isPortAvailable(preferredPort) {
			return preferredPort, nil
		}
	}

	// 優先ポートが使用できない場合は、利用可能なポートを自動選択
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("利用可能なポートの検索に失敗しました: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// isPortAvailable は指定されたポートが利用可能かどうかを確認する
func isPortAvailable(port int) bool {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// CalculateDefaultPort はサーバーURLハッシュから一貫したデフォルトポートを計算する
func CalculateDefaultPort(serverURLHash string) int {
	if len(serverURLHash) < 4 {
		return 3334 // フォールバック
	}

	// ハッシュの最初の4文字を16進数として解釈
	offset := 0
	for i, char := range serverURLHash[:4] {
		if i >= 4 {
			break
		}
		var value int
		switch {
		case char >= '0' && char <= '9':
			value = int(char - '0')
		case char >= 'a' && char <= 'f':
			value = int(char - 'a' + 10)
		case char >= 'A' && char <= 'F':
			value = int(char - 'A' + 10)
		default:
			value = 0
		}
		offset = offset*16 + value
	}

	// 3335から49151の範囲でポートを計算
	return 3335 + (offset % 45816)
}
