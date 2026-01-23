package executor

import (
	"context"
	"database/sql"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
)

// #==============================================================#
// ##          Interfaces                                        ##
// #==============================================================#

// DatabaseExecutor はデータベース操作のインターフェースです
type DatabaseExecutor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Ping() error
	Close() error
}

// #==============================================================#
// ##          Default Implementations                           ##
// #==============================================================#

// SQLRowWrapper は *sql.Row を RowInterface として扱うためのラッパー
type SQLRowWrapper struct {
	row *sql.Row
}

func (w *SQLRowWrapper) Scan(dest ...any) error {
	return w.row.Scan(dest...)
}

// SQLRowsWrapper は *sql.Rows を RowsInterface として扱うためのラッパー
type SQLRowsWrapper struct {
	rows *sql.Rows
}

func (w *SQLRowsWrapper) Columns() ([]string, error) {
	return w.rows.Columns()
}

func (w *SQLRowsWrapper) Next() bool {
	return w.rows.Next()
}

func (w *SQLRowsWrapper) Scan(dest ...any) error {
	return w.rows.Scan(dest...)
}

func (w *SQLRowsWrapper) Close() error {
	return w.rows.Close()
}

func (w *SQLRowsWrapper) Err() error {
	return w.rows.Err()
}

// DefaultDatabaseExecutor は標準のsql.DBを使用する実装
type DefaultDatabaseExecutor struct {
	DB *sql.DB
}

func (d *DefaultDatabaseExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.DB.QueryContext(ctx, query, args...)
}

func (d *DefaultDatabaseExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.DB.QueryRowContext(ctx, query, args...)
}

func (d *DefaultDatabaseExecutor) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.DB.BeginTx(ctx, opts)
}

func (d *DefaultDatabaseExecutor) Ping() error {
	return d.DB.Ping()
}

func (d *DefaultDatabaseExecutor) Close() error {
	return d.DB.Close()
}

// QueryContextRows は新しいインターフェース用のメソッド
func (d *DefaultDatabaseExecutor) QueryContextRows(ctx context.Context, query string, args ...any) (model.RowsInterface, error) {
	rows, err := d.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLRowsWrapper{rows: rows}, nil
}

// QueryRowContextRow は新しいインターフェース用のメソッド
func (d *DefaultDatabaseExecutor) QueryRowContextRow(ctx context.Context, query string, args ...any) model.RowInterface {
	row := d.DB.QueryRowContext(ctx, query, args...)
	return &SQLRowWrapper{row: row}
}
