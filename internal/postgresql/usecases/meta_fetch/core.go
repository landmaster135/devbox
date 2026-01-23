package usecases

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	dbExecutor "github.com/landmaster135/devbox/internal/postgresql/domain/executor"
	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
)

const defaultTableSchema = "public"

var IdentifierPattern = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

func qualifyTableIdentifier(tableName string) (qualified string, schema string, name string, err error) {
	if tableName == "" {
		return "", "", "", errors.New("テーブル名が指定されていません")
	}

	parts := strings.Split(tableName, ".")
	if len(parts) > 2 {
		return "", "", "", fmt.Errorf("サポートされていないテーブル識別子です: %s", tableName)
	}

	schema = defaultTableSchema
	name = parts[len(parts)-1]
	if len(parts) == 2 {
		schema = parts[0]
	}

	if schema == "" {
		schema = defaultTableSchema
	}

	if !IdentifierPattern.MatchString(schema) {
		return "", "", "", fmt.Errorf("スキーマ名に使用できない文字が含まれています: %s", schema)
	}
	if !IdentifierPattern.MatchString(name) {
		return "", "", "", fmt.Errorf("テーブル名に使用できない文字が含まれています: %s", tableName)
	}

	qualified = fmt.Sprintf("%s.%s", pq.QuoteIdentifier(schema), pq.QuoteIdentifier(name))
	return qualified, schema, name, nil
}

// GetTableSchemaMinimum はテーブルの最小限のスキーマ情報を取得します
func GetTableSchemaMinimum(ctx context.Context, exec dbExecutor.DatabaseExecutor, tableName string) ([]model.Column, error) {
	_, schema, name, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_schema = $1
		  AND table_name = $2
	`

	rows, err := exec.QueryContext(ctx, query, schema, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []model.Column
	for rows.Next() {
		var column model.Column
		if err := rows.Scan(&column.Name, &column.DataType); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

// fetchTableWithComments はテーブル名とコメントを取得します
func fetchTableWithComments(ctx context.Context, exec dbExecutor.DatabaseExecutor, tableName string) (model.TableSummary, error) {
	_, schema, name, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return model.TableSummary{}, err
	}

	query := `
		SELECT t.table_name,
		       COALESCE(pg_catalog.obj_description(c.oid), '') AS table_comment
		FROM information_schema.tables t
		JOIN pg_catalog.pg_class c ON c.relname = t.table_name
		JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE t.table_schema = $1
		  AND t.table_name = $2
		  AND n.nspname = $1
	`

	var table model.TableSummary
	err = exec.QueryRowContext(ctx, query, schema, name).Scan(&table.Name, &table.Comment)
	if err != nil {
		return model.TableSummary{}, err
	}

	return table, nil
}

// fetchPrimaryKeys はテーブルの主キーカラムを取得します
func fetchPrimaryKeys(ctx context.Context, exec dbExecutor.DatabaseExecutor, tableName string) ([]string, error) {
	qualified, _, _, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT a.attname
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		WHERE i.indrelid = $1::regclass
		AND i.indisprimary
		ORDER BY a.attnum
	`

	rows, err := exec.QueryContext(ctx, query, qualified)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var primaryKeys []string
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, err
		}
		primaryKeys = append(primaryKeys, columnName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return primaryKeys, nil
}

// fetchUniqueKeys はテーブルの一意キー制約を取得します
func fetchUniqueKeys(ctx context.Context, exec dbExecutor.DatabaseExecutor, tableName string) ([]model.UniqueKey, error) {
	query := `
		SELECT
			c.conname AS constraint_name,
			a.attname AS column_name
		FROM
			pg_constraint c
		JOIN
			pg_attribute a ON a.attrelid = c.conrelid AND a.attnum = ANY(c.conkey)
		WHERE
			c.conrelid = $1::regclass
			AND c.contype = 'u'
		ORDER BY
			c.conname, a.attnum
	`

	qualified, _, _, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	rows, err := exec.QueryContext(ctx, query, qualified)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uniqueKeys []model.UniqueKey
	var currentUniqueKey *model.UniqueKey
	for rows.Next() {
		var constraintName, columnName string
		if err := rows.Scan(&constraintName, &columnName); err != nil {
			return nil, err
		}

		if currentUniqueKey == nil || currentUniqueKey.Name != constraintName {
			uniqueKeys = append(uniqueKeys, model.UniqueKey{
				Name:    constraintName,
				Columns: []string{},
			})
			currentUniqueKey = &uniqueKeys[len(uniqueKeys)-1]
		}

		currentUniqueKey.Columns = append(currentUniqueKey.Columns, columnName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return uniqueKeys, nil
}

// fetchForeignKeys はテーブルの外部キー制約を取得します
func fetchForeignKeys(ctx context.Context, exec dbExecutor.DatabaseExecutor, tableName string) ([]model.ForeignKey, error) {
	query := `
		SELECT
			c.conname AS constraint_name,
			a.attname AS column_name,
			cl.relname AS referenced_table,
			af.attname AS referenced_column
		FROM
			pg_constraint c
		JOIN
			pg_attribute a ON a.attrelid = c.conrelid AND a.attnum = ANY(c.conkey)
		JOIN
			pg_class cl ON cl.oid = c.confrelid
		JOIN
			pg_attribute af ON af.attrelid = c.confrelid AND af.attnum = ANY(c.confkey)
		WHERE
			c.conrelid = $1::regclass
			AND c.contype = 'f'
		ORDER BY
			c.conname, a.attnum
	`

	qualified, _, _, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	rows, err := exec.QueryContext(ctx, query, qualified)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foreignKeys []model.ForeignKey
	var currentFK *model.ForeignKey
	for rows.Next() {
		var constraintName, columnName, refTableName, refColumnName string
		if err := rows.Scan(&constraintName, &columnName, &refTableName, &refColumnName); err != nil {
			return nil, err
		}

		if currentFK == nil || currentFK.Name != constraintName {
			foreignKeys = append(foreignKeys, model.ForeignKey{
				Name:       constraintName,
				RefTable:   refTableName,
				Columns:    []string{},
				RefColumns: []string{},
			})
			currentFK = &foreignKeys[len(foreignKeys)-1]
		}

		currentFK.Columns = append(currentFK.Columns, columnName)
		currentFK.RefColumns = append(currentFK.RefColumns, refColumnName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return foreignKeys, nil
}

// fetchTableColumns はテーブルのカラム情報を取得します
func fetchTableColumns(ctx context.Context, exec dbExecutor.DatabaseExecutor, tableName string) ([]model.ColumnInfo, error) {
	_, schema, name, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			c.column_name,
			c.data_type,
			c.is_nullable,
			c.column_default,
			COALESCE(pg_catalog.col_description(
				pg_catalog.pg_class.oid,
				c.ordinal_position
			), '') AS column_comment
		FROM
			information_schema.columns c
		JOIN
			pg_catalog.pg_class ON pg_catalog.pg_class.relname = c.table_name
		WHERE
			c.table_schema = $1
			AND c.table_name = $2
		ORDER BY
			c.ordinal_position
	`

	rows, err := exec.QueryContext(ctx, query, schema, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []model.ColumnInfo
	for rows.Next() {
		var col model.ColumnInfo
		if err := rows.Scan(&col.Name, &col.Type, &col.IsNullable, &col.Default, &col.Comment); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

// fetchTableIndexes はテーブルのインデックス情報を取得します
func fetchTableIndexes(ctx context.Context, exec dbExecutor.DatabaseExecutor, tableName string) ([]model.IndexInfo, error) {
	_, schema, name, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			i.relname AS index_name,
			a.attname AS column_name,
			ix.indisunique AS is_unique
		FROM
			pg_index ix
		JOIN
			pg_class i ON i.oid = ix.indexrelid
		JOIN
			pg_attribute a ON a.attrelid = ix.indrelid AND a.attnum = ANY(ix.indkey)
		JOIN
			pg_class t ON t.oid = ix.indrelid
	JOIN
		pg_namespace n ON n.oid = t.relnamespace
		LEFT JOIN
			pg_constraint c ON c.conindid = ix.indexrelid
		WHERE
			t.relname = $1
			AND n.nspname = $2
			AND c.contype IS NULL
			AND NOT ix.indisprimary
		ORDER BY
			i.relname, a.attnum
	`

	rows, err := exec.QueryContext(ctx, query, name, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []model.IndexInfo
	var currentIdx *model.IndexInfo
	for rows.Next() {
		var indexName, columnName string
		var isUnique bool
		if err := rows.Scan(&indexName, &columnName, &isUnique); err != nil {
			return nil, err
		}

		if currentIdx == nil || currentIdx.Name != indexName {
			indexes = append(indexes, model.IndexInfo{
				Name:    indexName,
				Unique:  isUnique,
				Columns: []string{},
			})
			currentIdx = &indexes[len(indexes)-1]
		}

		currentIdx.Columns = append(currentIdx.Columns, columnName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return indexes, nil
}

// GetAllTableSummaries はデータベース内の全てのテーブルのサマリー情報を取得します
func GetAllTableSummaries(ctx context.Context, exec dbExecutor.DatabaseExecutor) ([]model.TableSummary, error) {
	// テーブル一覧を取得
	query := `
		SELECT t.table_name,
		       COALESCE(pg_catalog.obj_description(pg_catalog.pg_class.oid), '') AS table_comment
		FROM information_schema.tables t
		JOIN pg_catalog.pg_class ON pg_catalog.pg_class.relname = t.table_name
		WHERE t.table_schema = 'public'
		ORDER BY t.table_name
	`

	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []model.TableSummary
	for rows.Next() {
		var table model.TableSummary
		if err := rows.Scan(&table.Name, &table.Comment); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 各テーブルの追加情報を取得
	for i := range tables {
		// 主キー情報を取得
		tables[i].PK, err = fetchPrimaryKeys(ctx, exec, tables[i].Name)
		if err != nil {
			return nil, err
		}

		// 一意キー情報を取得
		tables[i].UK, err = fetchUniqueKeys(ctx, exec, tables[i].Name)
		if err != nil {
			return nil, err
		}

		// 外部キー情報を取得
		tables[i].FK, err = fetchForeignKeys(ctx, exec, tables[i].Name)
		if err != nil {
			return nil, err
		}
	}

	return tables, nil
}

// #==============================================================#
// ##          Table Information Methods                         ##
// #==============================================================#

// GetTablesMinimum はデータベース内のテーブル一覧を取得します
func GetTablesMinimum(ctx context.Context, exec dbExecutor.DatabaseExecutor) ([]model.Table, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
	`

	rows, err := exec.QueryContext(ctx, query)
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

// #==============================================================#
// ##          High-level Service Methods                        ##
// #==============================================================#

// GetTableDetail はテーブルの詳細情報を取得します
func GetTableDetail(ctx context.Context, exec dbExecutor.DatabaseExecutor, tableName string) (*model.TableDetail, error) {
	// テーブル情報を取得
	tableInfo, err := fetchTableWithComments(ctx, exec, tableName)
	if err != nil {
		return nil, fmt.Errorf("テーブル情報の取得に失敗しました: %w", err)
	}

	// 主キー情報を取得
	primaryKeys, err := fetchPrimaryKeys(ctx, exec, tableName)
	if err != nil {
		return nil, fmt.Errorf("主キー情報の取得に失敗しました: %w", err)
	}

	// 一意キー情報を取得
	uniqueKeys, err := fetchUniqueKeys(ctx, exec, tableName)
	if err != nil {
		return nil, fmt.Errorf("一意キー情報の取得に失敗しました: %w", err)
	}

	// 外部キー情報を取得
	foreignKeys, err := fetchForeignKeys(ctx, exec, tableName)
	if err != nil {
		return nil, fmt.Errorf("外部キー情報の取得に失敗しました: %w", err)
	}

	// カラム情報を取得
	columns, err := fetchTableColumns(ctx, exec, tableName)
	if err != nil {
		return nil, fmt.Errorf("カラム情報の取得に失敗しました: %w", err)
	}

	// インデックス情報を取得
	indexes, err := fetchTableIndexes(ctx, exec, tableName)
	if err != nil {
		return nil, fmt.Errorf("インデックス情報の取得に失敗しました: %w", err)
	}

	// 詳細情報を作成
	detail := &model.TableDetail{
		Name:        tableInfo.Name,
		Comment:     tableInfo.Comment,
		Columns:     columns,
		PrimaryKeys: primaryKeys,
		UniqueKeys:  uniqueKeys,
		ForeignKeys: foreignKeys,
		Indexes:     indexes,
	}

	return detail, nil
}
