package dump

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
)

// validateOptions はオプションを検証します
func validateOptions(options *DumpOptions) error {
	if options == nil {
		return errors.New("オプションが指定されていません")
	}

	if _, _, err := quoteQualifiedTableName(options.TableName); err != nil {
		return err
	}

	if options.OutputPath == "" {
		options.OutputPath = "."
	}

	// パストラバーサル攻撃を防ぐ
	cleanPath := filepath.Clean(options.OutputPath)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("無効なパスが指定されました: %s", options.OutputPath)
	}

	// サポートされているフォーマットかチェック
	validFormats := map[string]bool{
		"json": true,
		"csv":  true,
		"sql":  true,
	}
	if !validFormats[options.Format] {
		return fmt.Errorf("サポートされていないフォーマットです: %s", options.Format)
	}

	if options.Limit != nil && *options.Limit <= 0 {
		return fmt.Errorf("limitは正の数である必要があります: %d", *options.Limit)
	}

	return nil
}

// getAllTables はデータベース内の全テーブル一覧を取得します
func (d *TableDumper) getAllTables(ctx context.Context) ([]model.Table, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`

	rows, err := d.executor.QueryContextRows(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []model.Table
	for rows.Next() {
		var table model.Table
		if err := rows.Scan(&table.Name); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// ensureAllowedTable はテーブルが存在し、ホワイトリストに一致することを確認します
func (d *TableDumper) ensureAllowedTable(ctx context.Context, tableName string) error {
	_, parts, err := quoteQualifiedTableName(tableName)
	if err != nil {
		return err
	}

	if len(parts) == 2 && parts[0] != model.DefaultTableSchema {
		return fmt.Errorf("サポートされていないスキーマが指定されました: %s", parts[0])
	}

	tables, err := d.getAllTables(ctx)
	if err != nil {
		return fmt.Errorf("テーブル一覧取得エラー: %w", err)
	}

	targetName := parts[len(parts)-1]
	for _, tbl := range tables {
		if tbl.Name == targetName {
			return nil
		}
	}

	return fmt.Errorf("指定されたテーブルが存在しません: %s", tableName)
}

// ensureOutputDirectory は出力ディレクトリが存在することを確認し、必要に応じて作成します
func (d *TableDumper) ensureOutputDirectory(outputPath string) error {
	if outputPath == "" || outputPath == "." {
		return nil
	}

	return d.fileWriter.MkdirAll(outputPath, 0755)
}

func (d *TableDumper) ensureOptions(ctx context.Context, options *DumpOptions) error {
	// パラメータ検証
	if err := validateOptions(options); err != nil {
		return fmt.Errorf("オプション検証エラー: %w", err)
	}

	if err := d.ensureAllowedTable(ctx, options.TableName); err != nil {
		return fmt.Errorf("テーブル検証エラー: %w", err)
	}

	// 出力ディレクトリを作成
	if err := d.ensureOutputDirectory(options.OutputPath); err != nil {
		return fmt.Errorf("出力ディレクトリ作成エラー: %w", err)
	}

	return nil
}
