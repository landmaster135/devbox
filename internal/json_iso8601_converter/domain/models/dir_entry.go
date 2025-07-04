package models

// DirEntry はディレクトリエントリを表す構造体です
type DirEntry struct {
	Name  string
	IsDir bool
	Path  string
}

// NewDirEntry は新しいDirEntryインスタンスを作成します
func NewDirEntry(name string, isDir bool, path string) *DirEntry {
	return &DirEntry{
		Name:  name,
		IsDir: isDir,
		Path:  path,
	}
}
