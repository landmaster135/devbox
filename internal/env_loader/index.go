package envloader

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Load(keys []string) (map[string]string, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf(".envの読み込みに失敗しました: %w", err)
	}

	values := make(map[string]string, len(keys))
	for _, key := range keys {
		value := os.Getenv(key)
		if value == "" {
			return nil, fmt.Errorf("環境変数 %s が設定されていません", key)
		}
		values[key] = value
	}

	return values, nil
}
