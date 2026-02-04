package usecases

import (
	"fmt"
	"io/fs"
	"os"

	fetchers "github.com/landmaster135/devbox/internal/html_sanitizer/usecases/sanitizer"
)

// fileReadFunc は入力ファイルを読み込む関数の型です。
type fileReadFunc func(string) ([]byte, error)

// fileWriteFunc は結果をファイルへ書き込む関数の型です。
type fileWriteFunc func(string, []byte, fs.FileMode) error

// sanitizeFunc はHTML文字列をサニタイズする関数の型です。
type sanitizeFunc func(string, bool) (string, error)

// SanitizerService はファイル入出力とHTMLサニタイズを仲介します。
type SanitizerService struct {
	readFile   fileReadFunc
	writeFile  fileWriteFunc
	sanitize   sanitizeFunc
	outputPerm fs.FileMode
}

// SanitizerServiceOption はSanitizerServiceの依存を差し替えるためのオプションです。
type SanitizerServiceOption func(*SanitizerService)

// WithFileReader はファイル読み込み関数を上書きします。
func WithFileReader(reader fileReadFunc) SanitizerServiceOption {
	return func(s *SanitizerService) {
		if reader != nil {
			s.readFile = reader
		}
	}
}

// WithFileWriter はファイル書き込み関数を上書きします。
func WithFileWriter(writer fileWriteFunc) SanitizerServiceOption {
	return func(s *SanitizerService) {
		if writer != nil {
			s.writeFile = writer
		}
	}
}

// WithSanitizer はHTMLサニタイズ関数を上書きします。
func WithSanitizer(fn sanitizeFunc) SanitizerServiceOption {
	return func(s *SanitizerService) {
		if fn != nil {
			s.sanitize = fn
		}
	}
}

// WithOutputPermission は出力ファイルのパーミッションを設定します。
func WithOutputPermission(perm fs.FileMode) SanitizerServiceOption {
	return func(s *SanitizerService) {
		if perm != 0 {
			s.outputPerm = perm
		}
	}
}

// NewSanitizerService はSanitizerServiceのデフォルト実装を返します。
func NewSanitizerService(opts ...SanitizerServiceOption) *SanitizerService {
	svc := &SanitizerService{
		readFile:   os.ReadFile,
		writeFile:  os.WriteFile,
		sanitize:   fetchers.SanitizeHTMLBody,
		outputPerm: 0o644,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

// SanitizeFile は入力ファイルを読み込みサニタイズ結果を書き出します。
func (s *SanitizerService) SanitizeFile(inputPath, outputPath string, omitFullBody bool) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("inputPathが指定されていません")
	}
	if outputPath == "" {
		return "", fmt.Errorf("outputPathが指定されていません")
	}

	contents, err := s.readFile(inputPath)
	if err != nil {
		return "", fmt.Errorf("入力ファイル(%s)の読み込みに失敗しました: %w", inputPath, err)
	}

	sanitized, err := s.sanitize(string(contents), omitFullBody)
	if err != nil {
		return "", fmt.Errorf("HTMLのサニタイズに失敗しました: %w", err)
	}

	if err := s.writeFile(outputPath, []byte(sanitized), s.outputPerm); err != nil {
		return "", fmt.Errorf("出力ファイル(%s)の書き込みに失敗しました: %w", outputPath, err)
	}

	return sanitized, nil
}
