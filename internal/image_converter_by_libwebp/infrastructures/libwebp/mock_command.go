package libwebp

import (
	"context"
	"sync"
)

// ConvertCall は MockConverter へ渡された変換引数を保持します。
type ConvertCall struct {
	InputPath  string
	OutputPath string
	Quality    int
	Method     int
	Lossless   bool
}

// MockConverter は usecase テスト用の Converter 実装です。
type MockConverter struct {
	CheckAvailableFunc func() error
	ConvertToWebPFunc  func(ctx context.Context, inputPath string, outputPath string, quality int, method int, lossless bool) error
	Calls              []ConvertCall
	mu                 sync.Mutex
}

func (m *MockConverter) CheckAvailable() error {
	if m.CheckAvailableFunc != nil {
		return m.CheckAvailableFunc()
	}
	return nil
}

func (m *MockConverter) ConvertToWebP(ctx context.Context, inputPath string, outputPath string, quality int, method int, lossless bool) error {
	m.mu.Lock()
	m.Calls = append(m.Calls, ConvertCall{
		InputPath:  inputPath,
		OutputPath: outputPath,
		Quality:    quality,
		Method:     method,
		Lossless:   lossless,
	})
	m.mu.Unlock()
	if m.ConvertToWebPFunc != nil {
		return m.ConvertToWebPFunc(ctx, inputPath, outputPath, quality, method, lossless)
	}
	return nil
}

func (m *MockConverter) CallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Calls)
}
