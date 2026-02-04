package dump

import (
	"context"
	"errors"
	"fmt"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
	sql "github.com/landmaster135/devbox/internal/postgresql/domain/sql"
)

// buildQuery はダンプ用のクエリを構築します
func buildQuery(options *DumpOptions) (string, error) {
	if options == nil {
		return "", errors.New("オプションが指定されていません")
	}

	quotedTable, _, err := sql.QuoteQualifiedTableName(options.TableName)
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
