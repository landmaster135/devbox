package usecases

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/landmaster135/devbox/internal/image_converter_by_libwebp/config"
	"github.com/landmaster135/devbox/internal/image_converter_by_libwebp/infrastructures/libwebp"
)

var supportedInputExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".tif":  {},
	".tiff": {},
	".webp": {},
}

// Result は変換結果の集計値です。
type Result struct {
	SuccessCount int
	ErrorCount   int
	OutputDir    string
}

// Service は libwebp による画像変換ユースケースを提供します。
type Service struct {
	converter libwebp.Converter
}

// NewService はデフォルト依存の Service を作成します。
func NewService() *Service {
	return NewServiceWithConverter(libwebp.NewCommandConverter())
}

// NewServiceWithConverter は依存を注入して Service を作成します。
func NewServiceWithConverter(converter libwebp.Converter) *Service {
	if converter == nil {
		converter = libwebp.NewCommandConverter()
	}
	return &Service{converter: converter}
}

// Convert は設定に従って複数画像を WebP へ変換します。
func (s *Service) Convert(ctx context.Context, cfg *config.Config) (*Result, error) {
	if cfg == nil {
		return nil, fmt.Errorf("設定が指定されていません")
	}
	if err := s.converter.CheckAvailable(); err != nil {
		return nil, err
	}

	srcDir, outDir, archiveDir, err := s.prepareDirectories(cfg)
	if err != nil {
		return nil, err
	}

	paths, err := s.findImageFiles(srcDir, cfg.Recursive)
	if err != nil {
		return nil, err
	}

	result := &Result{OutputDir: outDir}
	if len(paths) == 0 {
		return result, nil
	}

	jobs := make(chan string, len(paths))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				if err := s.convertOne(ctx, path, srcDir, outDir, archiveDir, cfg); err != nil {
					mu.Lock()
					result.ErrorCount++
					mu.Unlock()
					continue
				}
				mu.Lock()
				result.SuccessCount++
				mu.Unlock()
			}
		}()
	}

	for _, path := range paths {
		jobs <- path
	}
	close(jobs)
	wg.Wait()

	if result.ErrorCount > 0 {
		return result, fmt.Errorf("%d ファイルの変換に失敗しました", result.ErrorCount)
	}
	return result, nil
}

func (s *Service) prepareDirectories(cfg *config.Config) (string, string, string, error) {
	srcDir, err := filepath.Abs(cfg.SrcDir)
	if err != nil {
		return "", "", "", fmt.Errorf("入力ディレクトリパスの変換に失敗しました: %w", err)
	}
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return "", "", "", fmt.Errorf("入力ディレクトリを確認できません: %w", err)
	}
	if !srcInfo.IsDir() {
		return "", "", "", fmt.Errorf("入力パスはディレクトリではありません: %s", srcDir)
	}

	outDir, err := filepath.Abs(cfg.OutDir)
	if err != nil {
		return "", "", "", fmt.Errorf("出力ディレクトリパスの変換に失敗しました: %w", err)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", "", "", fmt.Errorf("出力ディレクトリの作成に失敗しました: %w", err)
	}

	archiveDir := ""
	if cfg.ArchiveDir != "" {
		archiveDir, err = filepath.Abs(cfg.ArchiveDir)
		if err != nil {
			return "", "", "", fmt.Errorf("退避ディレクトリパスの変換に失敗しました: %w", err)
		}
		if err := os.MkdirAll(archiveDir, 0o755); err != nil {
			return "", "", "", fmt.Errorf("退避ディレクトリの作成に失敗しました: %w", err)
		}
	}

	return srcDir, outDir, archiveDir, nil
}

func (s *Service) findImageFiles(srcDir string, recursive bool) ([]string, error) {
	paths := make([]string, 0)
	if recursive {
		err := filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !isSupportedInput(path) {
				return nil
			}
			paths = append(paths, path)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("入力ディレクトリの走査に失敗しました: %w", err)
		}
		return paths, nil
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, fmt.Errorf("入力ディレクトリの読み込みに失敗しました: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(srcDir, entry.Name())
		if isSupportedInput(path) {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

func (s *Service) convertOne(ctx context.Context, path string, srcDir string, outDir string, archiveDir string, cfg *config.Config) error {
	outPath, err := s.buildOutputPath(path, srcDir, outDir, cfg.OutExt)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("出力先ディレクトリの作成に失敗しました: %w", err)
	}
	if err := s.converter.ConvertToWebP(ctx, path, outPath, cfg.Quality, cfg.Method, cfg.Lossless); err != nil {
		return err
	}
	if archiveDir == "" {
		return nil
	}
	return s.archiveOriginal(path, srcDir, archiveDir, cfg.Move)
}

func (s *Service) buildOutputPath(path string, srcDir string, outDir string, outExt string) (string, error) {
	rel, err := filepath.Rel(srcDir, path)
	if err != nil {
		return "", fmt.Errorf("相対パスの算出に失敗しました: %w", err)
	}
	base := strings.TrimSuffix(rel, filepath.Ext(rel))
	return filepath.Join(outDir, base+"."+outExt), nil
}

func (s *Service) archiveOriginal(srcPath string, srcDir string, archiveDir string, move bool) error {
	rel, err := filepath.Rel(srcDir, srcPath)
	if err != nil {
		return fmt.Errorf("退避先相対パスの算出に失敗しました: %w", err)
	}
	destPath := filepath.Join(archiveDir, rel)
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("退避先ディレクトリの作成に失敗しました: %w", err)
	}
	if !move {
		return copyFile(srcPath, destPath)
	}
	if err := os.Rename(srcPath, destPath); err == nil {
		return nil
	}
	if err := copyFile(srcPath, destPath); err != nil {
		return err
	}
	return os.Remove(srcPath)
}

func isSupportedInput(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := supportedInputExtensions[ext]
	return ok
}

func copyFile(srcPath string, destPath string) error {
	in, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("退避元ファイルの読み込みに失敗しました: %w", err)
	}
	defer in.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("退避先ファイルの作成に失敗しました: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("退避ファイルのコピーに失敗しました: %w", err)
	}
	if err := out.Sync(); err != nil {
		return fmt.Errorf("退避ファイルの同期に失敗しました: %w", err)
	}
	return nil
}
