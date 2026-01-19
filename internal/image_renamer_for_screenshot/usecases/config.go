package usecases

import (
	"fmt"
	"io"
	"os"
)

// Operation はリネーム対象となるスクリーンショットの種別を表します
type Operation string

const (
	OperationUnknown Operation = ""
	OperationVLC     Operation = "vlc"
	OperationWin     Operation = "win"
	OperationPixel   Operation = "pixel"
	OperationXiaomi  Operation = "xiaomi"
	operationAll     Operation = "all"
)

func (o Operation) isValidSelection() bool {
	switch o {
	case OperationVLC, OperationWin, OperationPixel, OperationXiaomi:
		return true
	default:
		return false
	}
}

// Config はプログラムの設定を保持する構造体です
type Config struct {
	SrcDir     string
	Recursive  bool
	Workers    int
	Operation  Operation
	ToDateTime bool
}

// validateConfig は設定の妥当性を検証します
func validateConfig(config Config, stderr io.Writer) error {
	// --to-datetimeが指定されている場合は、他のパターンは不要
	if config.ToDateTime {
		// ディレクトリの存在確認
		if _, err := os.Stat(config.SrcDir); err != nil {
			fmt.Fprintf(stderr, "エラー: ディレクトリ %s へのアクセスエラー: %v\n", config.SrcDir, err)
			return err
		}
		return nil
	}

	if !config.Operation.isValidSelection() {
		fmt.Fprintln(stderr, "エラー: -operation には 'vlc'、'win'、'pixel'、または 'xiaomi' のいずれかを指定する必要があります。")
		fmt.Fprintln(stderr, "例: ./image-renamer-for-screenshot -operation=vlc")
		fmt.Fprintln(stderr, "例: ./image-renamer-for-screenshot -to-datetime")
		return fmt.Errorf("無効なoperationが指定されています")
	}

	// ディレクトリの存在確認
	if _, err := os.Stat(config.SrcDir); err != nil {
		fmt.Fprintf(stderr, "エラー: ディレクトリ %s へのアクセスエラー: %v\n", config.SrcDir, err)
		return err
	}

	return nil
}
