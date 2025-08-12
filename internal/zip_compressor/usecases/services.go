package usecases

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileSystemOperator インターフェースを定義（テスト用）
type FileSystemOperator interface {
	Stat(name string) (os.FileInfo, error)
	Create(name string) (*os.File, error)
	Open(name string) (*os.File, error)
	MkdirAll(path string, perm os.FileMode) error
	Walk(root string, walkFn filepath.WalkFunc) error
}

// DefaultFileSystemOperator は標準のosパッケージを使用する実装
type DefaultFileSystemOperator struct{}

func (f *DefaultFileSystemOperator) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (f *DefaultFileSystemOperator) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (f *DefaultFileSystemOperator) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (f *DefaultFileSystemOperator) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f *DefaultFileSystemOperator) Walk(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(root, walkFn)
}

// ZipCompressorService はZip圧縮・展開を行うサービスです
type ZipCompressorService struct {
	fileSystem FileSystemOperator
}

// NewZipCompressorService は新しいZipCompressorServiceを作成します
func NewZipCompressorService() *ZipCompressorService {
	return &ZipCompressorService{
		fileSystem: &DefaultFileSystemOperator{},
	}
}

// NewZipCompressorServiceWithFileSystem はFileSystemOperatorを注入した新しいZipCompressorServiceを作成します
func NewZipCompressorServiceWithFileSystem(fileSystem FileSystemOperator) *ZipCompressorService {
	return &ZipCompressorService{
		fileSystem: fileSystem,
	}
}

// HandleCompress はファイルまたはディレクトリを圧縮します
func (s *ZipCompressorService) HandleCompress(path string) (string, error) {
	// パスの存在確認
	info, err := s.fileSystem.Stat(path)
	if err != nil {
		return "", fmt.Errorf("指定されたパスが存在しません: %s", path)
	}

	// 出力ファイル名を決定
	var outputPath string
	if info.IsDir() {
		// ディレクトリの場合: ディレクトリ名.zip
		outputPath = path + ".zip"
	} else {
		// ファイルの場合: ファイル名.zip
		outputPath = path + ".zip"
	}

	// Zipファイルを作成
	zipFile, err := s.fileSystem.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("zipファイルの作成に失敗しました: %v", err)
	}
	defer zipFile.Close()

	// Zipライターを作成
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if info.IsDir() {
		// ディレクトリの場合
		err = s.compressDirectory(path, zipWriter)
	} else {
		// ファイルの場合
		err = s.compressFile(path, zipWriter)
	}

	if err != nil {
		return "", fmt.Errorf("圧縮処理に失敗しました: %v", err)
	}

	return fmt.Sprintf("圧縮が完了しました: %s", outputPath), nil
}

// compressFile は単一ファイルを圧縮します
func (s *ZipCompressorService) compressFile(filePath string, zipWriter *zip.Writer) error {
	// ファイルを開く
	file, err := s.fileSystem.Open(filePath)
	if err != nil {
		return fmt.Errorf("ファイルを開けませんでした: %v", err)
	}
	defer file.Close()

	// ファイル情報を取得
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("ファイル情報の取得に失敗しました: %v", err)
	}

	// Zipエントリのヘッダーを作成
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("zipヘッダーの作成に失敗しました: %v", err)
	}

	// ファイル名を設定（パスではなくファイル名のみ）
	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	// Zipエントリを作成
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("zipエントリの作成に失敗しました: %v", err)
	}

	// ファイル内容をコピー
	_, err = io.Copy(writer, file)
	if err != nil {
		return fmt.Errorf("ファイル内容のコピーに失敗しました: %v", err)
	}

	return nil
}

// compressDirectory はディレクトリを再帰的に圧縮します
func (s *ZipCompressorService) compressDirectory(dirPath string, zipWriter *zip.Writer) error {
	// ディレクトリを再帰的に走査
	return s.fileSystem.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 相対パスを計算
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return fmt.Errorf("相対パスの計算に失敗しました: %v", err)
		}

		// ルートディレクトリ自体はスキップ
		if relPath == "." {
			return nil
		}

		// Zipエントリのヘッダーを作成
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("zipヘッダーの作成に失敗しました: %v", err)
		}

		// パス区切り文字を統一（Windowsでも/を使用）
		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			// ディレクトリの場合は末尾に/を追加
			header.Name += "/"
		} else {
			// ファイルの場合は圧縮方法を設定
			header.Method = zip.Deflate
		}

		// Zipエントリを作成
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("zipエントリの作成に失敗しました: %v", err)
		}

		// ファイルの場合は内容をコピー
		if !info.IsDir() {
			file, err := s.fileSystem.Open(path)
			if err != nil {
				return fmt.Errorf("ファイルを開けませんでした: %v", err)
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return fmt.Errorf("ファイル内容のコピーに失敗しました: %v", err)
			}
		}

		return nil
	})
}

// HandleDecompress はZipファイルを展開します
func (s *ZipCompressorService) HandleDecompress(zipPath string) (string, error) {
	// Zipファイルの存在確認
	info, err := s.fileSystem.Stat(zipPath)
	if err != nil {
		return "", fmt.Errorf("指定されたZipファイルが存在しません: %s", zipPath)
	}

	if info.IsDir() {
		return "", fmt.Errorf("指定されたパスはディレクトリです。Zipファイルを指定してください: %s", zipPath)
	}

	// .zip拡張子の確認
	if !strings.HasSuffix(strings.ToLower(zipPath), ".zip") {
		return "", fmt.Errorf("指定されたファイルはZipファイルではありません: %s", zipPath)
	}

	// 出力ディレクトリ名を決定
	baseName := strings.TrimSuffix(filepath.Base(zipPath), ".zip")
	outputDir := filepath.Join(filepath.Dir(zipPath), baseName+"_decompressed")

	// 出力ディレクトリを作成
	err = s.fileSystem.MkdirAll(outputDir, 0755)
	if err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗しました: %v", err)
	}

	// Zipファイルを開く
	zipFile, err := s.fileSystem.Open(zipPath)
	if err != nil {
		return "", fmt.Errorf("zipファイルを開けませんでした: %v", err)
	}
	defer zipFile.Close()

	// Zipリーダーを作成
	zipReader, err := zip.NewReader(zipFile, info.Size())
	if err != nil {
		return "", fmt.Errorf("zipリーダーの作成に失敗しました: %v", err)
	}

	// 各エントリを展開
	for _, file := range zipReader.File {
		err = s.extractFile(file, outputDir)
		if err != nil {
			return "", fmt.Errorf("ファイルの展開に失敗しました: %v", err)
		}
	}

	return fmt.Sprintf("展開が完了しました: %s", outputDir), nil
}

// extractFile は単一のZipエントリを展開します
func (s *ZipCompressorService) extractFile(file *zip.File, outputDir string) error {
	// 出力パスを構築
	outputPath := filepath.Join(outputDir, file.Name)

	// パストラバーサル攻撃を防ぐ
	if !strings.HasPrefix(outputPath, filepath.Clean(outputDir)+string(os.PathSeparator)) {
		return fmt.Errorf("不正なパスが検出されました: %s", file.Name)
	}

	if file.FileInfo().IsDir() {
		// ディレクトリの場合
		return s.fileSystem.MkdirAll(outputPath, file.FileInfo().Mode())
	}

	// ファイルの場合
	// 親ディレクトリを作成
	err := s.fileSystem.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return fmt.Errorf("親ディレクトリの作成に失敗しました: %v", err)
	}

	// Zipエントリを開く
	zipFileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("zipエントリを開けませんでした: %v", err)
	}
	defer zipFileReader.Close()

	// 出力ファイルを作成
	outputFile, err := s.fileSystem.Create(outputPath)
	if err != nil {
		return fmt.Errorf("出力ファイルの作成に失敗しました: %v", err)
	}
	defer outputFile.Close()

	// ファイル内容をコピー
	_, err = io.Copy(outputFile, zipFileReader)
	if err != nil {
		return fmt.Errorf("ファイル内容のコピーに失敗しました: %v", err)
	}

	return nil
}
