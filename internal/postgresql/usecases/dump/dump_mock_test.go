package dump

import (
	"context"
	"errors"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
)

// #==============================================================#
// ##          Mock SQL Types                                    ##
// #==============================================================#

// MockRows は sql.Rows のモック
type MockRows struct {
	columns []string
	data    [][]any
	index   int
	closed  bool
}

func NewMockRows(columns []string, data [][]any) *MockRows {
	return &MockRows{
		columns: columns,
		data:    data,
		index:   -1,
		closed:  false,
	}
}

func (m *MockRows) Columns() ([]string, error) {
	return m.columns, nil
}

func (m *MockRows) Next() bool {
	if m.closed || m.index >= len(m.data)-1 {
		return false
	}
	m.index++
	return true
}

func (m *MockRows) Scan(dest ...any) error {
	if m.closed || m.index < 0 || m.index >= len(m.data) {
		return errors.New("no rows")
	}

	row := m.data[m.index]
	for i, val := range row {
		if i < len(dest) {
			switch d := dest[i].(type) {
			case *string:
				if val != nil {
					*d = val.(string)
				}
			case *any:
				*d = val
			}
		}
	}
	return nil
}

func (m *MockRows) Close() error {
	m.closed = true
	return nil
}

func (m *MockRows) Err() error {
	return nil
}

// MockRow は sql.Row のモック
type MockRow struct {
	data []any
	err  error
}

func NewMockRow(data []any, err error) *MockRow {
	return &MockRow{data: data, err: err}
}

func (m *MockRow) Scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}

	for i, val := range m.data {
		if i < len(dest) {
			switch d := dest[i].(type) {
			case *string:
				if val != nil {
					*d = val.(string)
				}
			case *any:
				*d = val
			}
		}
	}
	return nil
}

// #==============================================================#
// ##          Mock Database Query Executor                      ##
// #==============================================================#

// MockDatabaseQueryExecutor は新しいインターフェース用のモック
type MockDatabaseQueryExecutor struct {
	QueryContextRowsFunc   func(ctx context.Context, query string, args ...any) (model.RowsInterface, error)
	QueryRowContextRowFunc func(ctx context.Context, query string, args ...any) model.RowInterface
}

func (m *MockDatabaseQueryExecutor) QueryContextRows(ctx context.Context, query string, args ...any) (model.RowsInterface, error) {
	if m.QueryContextRowsFunc != nil {
		return m.QueryContextRowsFunc(ctx, query, args)
	}
	return nil, nil
}

func (m *MockDatabaseQueryExecutor) QueryRowContextRow(ctx context.Context, query string, args ...any) model.RowInterface {
	if m.QueryRowContextRowFunc != nil {
		return m.QueryRowContextRowFunc(ctx, query, args)
	}
	return nil
}
