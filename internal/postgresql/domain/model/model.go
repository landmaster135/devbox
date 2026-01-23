package domain

import (
	"database/sql"
)

// #==============================================================#
// ##          Data Structures                                   ##
// #==============================================================#

// ColumnInfo はデータベースのカラム詳細情報を表します
type ColumnInfo struct {
	Name       string         `json:"column_name"`
	Type       string         `json:"data_type"`
	IsNullable string         `json:"is_nullable"`
	Default    sql.NullString `json:"column_default"`
	Comment    string         `json:"column_comment"`
}

// TableSummary はテーブルのサマリー情報を表します
type TableSummary struct {
	Name    string       `json:"table_name"`
	Comment string       `json:"table_comment"`
	PK      []string     `json:"primary_keys"`
	UK      []UniqueKey  `json:"unique_keys"`
	FK      []ForeignKey `json:"foreign_keys"`
}

// UniqueKey は一意キー制約を表します
type UniqueKey struct {
	Name    string   `json:"constraint_name"`
	Columns []string `json:"columns"`
}

// ForeignKey は外部キー制約を表します
type ForeignKey struct {
	Name       string   `json:"constraint_name"`
	Columns    []string `json:"columns"`
	RefTable   string   `json:"referenced_table"`
	RefColumns []string `json:"referenced_columns"`
}

// IndexInfo はインデックス情報を表します
type IndexInfo struct {
	Name    string   `json:"index_name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"is_unique"`
}

// ListTablesData はテーブル一覧のテンプレートに渡すデータ構造
type ListTablesData struct {
	Tables []TableSummary
}

// TableDetail はテーブルの詳細情報を表します
type TableDetail struct {
	Name        string       `json:"table_name"`
	Comment     string       `json:"table_comment"`
	Columns     []ColumnInfo `json:"columns"`
	PrimaryKeys []string     `json:"primary_keys"`
	UniqueKeys  []UniqueKey  `json:"unique_keys"`
	ForeignKeys []ForeignKey `json:"foreign_keys"`
	Indexes     []IndexInfo  `json:"indexes"`
}

// Column はデータベースのカラム情報を表します（最小限の情報）
type Column struct {
	Name     string `json:"column_name"`
	DataType string `json:"data_type"`
}

// Table はデータベースのテーブル情報を表します
type Table struct {
	Name string `json:"table_name"`
}

// #==============================================================#
// ##          Interfaces                                        ##
// #==============================================================#

// RowsInterface は sql.Rows の操作を抽象化するインターフェースです
type RowsInterface interface {
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

// RowInterface は sql.Row の操作を抽象化するインターフェースです
type RowInterface interface {
	Scan(dest ...any) error
}
