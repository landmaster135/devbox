package usecases

import (
	"os"
)

// FileSystem はファイルシステム操作のインターフェースです
type FileSystem interface {
	ReadDir(dirname string) ([]os.DirEntry, error)
	Stat(name string) (os.FileInfo, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	ReadFile(name string) ([]byte, error)
	MkdirAll(path string, perm os.FileMode) error
}

// InputReader は入力読み取りのインターフェースです
type InputReader interface {
	ReadString(delim byte) (string, error)
}

// READMEParser はREADMEファイル解析のインターフェースです
type READMEParser interface {
	ParseUsageExamples(content []byte) ([]string, error)
}

// ScriptGenerator はスクリプト生成のインターフェースです
type ScriptGenerator interface {
	GenerateContent(packageName, packagePath string, usageExamples []string) string
}
