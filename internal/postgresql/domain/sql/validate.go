package sql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"

	metaFetch "github.com/landmaster135/devbox/internal/postgresql/usecases/meta_fetch"
)

func QuoteQualifiedTableName(tableName string) (string, []string, error) {
	if tableName == "" {
		return "", nil, errors.New("テーブル名が指定されていません")
	}

	parts := strings.Split(tableName, ".")
	if len(parts) > 2 {
		return "", nil, fmt.Errorf("サポートされていないテーブル識別子です: %s", tableName)
	}

	quotedParts := make([]string, len(parts))
	for i, part := range parts {
		if part == "" || !metaFetch.IdentifierPattern.MatchString(part) {
			return "", nil, fmt.Errorf("テーブル名に使用できない文字が含まれています: %s", tableName)
		}
		quotedParts[i] = pq.QuoteIdentifier(part)
	}

	return strings.Join(quotedParts, "."), parts, nil
}
