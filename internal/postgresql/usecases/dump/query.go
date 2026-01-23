package dump

import (
	"context"
	"errors"
	"fmt"
	"strings"

	pq "github.com/lib/pq"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
	metaFetch "github.com/landmaster135/devbox/internal/postgresql/usecases/meta_fetch"
)

func quoteQualifiedTableName(tableName string) (string, []string, error) {
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

// buildQuery はダンプ用のクエリを構築します
func buildQuery(options *DumpOptions) (string, error) {
	if options == nil {
		return "", errors.New("オプションが指定されていません")
	}

	quotedTable, _, err := quoteQualifiedTableName(options.TableName)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("SELECT * FROM %s", quotedTable)

	if options.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *options.Limit)
	}

	return query, nil
}

func (d *TableDumper) queryRows(ctx context.Context, options *DumpOptions) (model.RowsInterface, error) {
	// クエリを構築
	query, err := buildQuery(options)
	if err != nil {
		return nil, fmt.Errorf("クエリ構築エラー: %w", err)
	}

	rows, err := d.executor.QueryContextRows(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("データ取得エラー: %w", err)
	}

	return rows, nil
}
