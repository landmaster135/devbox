// Package valkey はValkeyサーバーとの連携機能を提供します
package valkey

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/valkey-io/valkey-go"
)

// Server はValkeyサーバーのプロセスを管理する構造体
type Server struct {
	cmd       *exec.Cmd
	port      int
	address   string
	isRunning bool
}

// NewServer は新しいServerインスタンスを作成します
// アドレスが空の場合はエラーを返します
func NewServer(address string) (*Server, error) {
	// アドレスが空の場合はエラー
	if address == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	// アドレスからポートを抽出
	port := 0
	parts := strings.Split(address, ":")
	if len(parts) > 1 {
		if p, err := strconv.Atoi(parts[1]); err == nil {
			port = p
		} else {
			return nil, fmt.Errorf("invalid port in address: %s", address)
		}
	} else {
		return nil, fmt.Errorf("invalid address format (expected host:port): %s", address)
	}

	return &Server{
		port:      port,
		address:   address,
		isRunning: false,
	}, nil
}

// checkPortAvailability はポートが利用可能かどうかを確認します
func (s *Server) checkPortAvailability() bool {
	conn, err := net.DialTimeout("tcp", s.address, time.Second)
	if err != nil {
		// エラーが発生した場合、ポートは利用可能
		return true
	}
	// 接続が成功した場合、ポートは既に使用されている
	conn.Close()
	return false
}

// Start はValkeyサーバーを起動します
func (s *Server) Start() error {
	if s.isRunning {
		return nil // 既に起動している場合は何もしない
	}

	// ポートが既に使用されているか確認
	if !s.checkPortAvailability() {
		// 既存のサーバーを使用する
		log.Printf("アドレス %s は既に使用されています。既存のサーバーを使用します。", s.address)
		s.isRunning = true
		return nil
	}

	// Valkeyサーバーを起動
	s.cmd = exec.Command("valkey-server", "--port", strconv.Itoa(s.port))

	// 標準出力と標準エラー出力をログに出力
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// 非同期でログを出力
	go s.pipeOutput(stdout, "VALKEY-STDOUT")
	go s.pipeOutput(stderr, "VALKEY-STDERR")

	// サーバー起動
	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start valkey server: %w", err)
	}

	s.isRunning = true

	// サーバーが起動するまで少し待つ
	time.Sleep(2 * time.Second)

	// 接続確認
	if err := s.checkConnection(); err != nil {
		s.Stop() // 接続できない場合は停止
		return fmt.Errorf("failed to connect to valkey server: %w", err)
	}

	log.Printf("Valkeyサーバーを起動しました (PID: %d, アドレス: %s)", s.cmd.Process.Pid, s.address)
	return nil
}

// Stop はValkeyサーバーを停止します
func (s *Server) Stop() error {
	if !s.isRunning || s.cmd == nil || s.cmd.Process == nil {
		return nil // 起動していない場合は何もしない
	}

	// プロセスを停止
	if err := s.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop valkey server: %w", err)
	}

	s.isRunning = false
	log.Printf("Valkeyサーバーを停止しました")
	return nil
}

// pipeOutput はコマンドの出力をログに出力します
func (s *Server) pipeOutput(reader io.Reader, prefix string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Printf("[%s] %s", prefix, scanner.Text())
	}
}

// checkConnection はValkeyサーバーへの接続を確認します
func (s *Server) checkConnection() error {
	// 接続確認のためにpingを送信
	opt, err := valkey.ParseURL(fmt.Sprintf("valkey://%s", s.address))
	if err != nil {
		return err
	}

	client, err := valkey.NewClient(opt)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// PINGコマンドを実行
	result := client.Do(ctx, client.B().Ping().Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// GetAddress はValkeyサーバーのアドレスを返します
func (s *Server) GetAddress() string {
	return s.address
}
