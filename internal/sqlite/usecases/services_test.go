package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockTableRepository struct {
	listTablesFunc func(ctx context.Context, dbPath string) ([]string, error)
}

func (m *mockTableRepository) ListTables(ctx context.Context, dbPath string) ([]string, error) {
	return m.listTablesFunc(ctx, dbPath)
}

func TestNewSQLiteService_NormalDefaultRepository(t *testing.T) {
	t.Parallel()

	service := NewSQLiteService(nil)
	require.NotNil(t, service)
	require.NotNil(t, service.repository)
}

func TestSQLiteService_HandleListTables_NormalText(t *testing.T) {
	t.Parallel()

	service := NewSQLiteService(&mockTableRepository{
		listTablesFunc: func(ctx context.Context, dbPath string) ([]string, error) {
			require.Equal(t, "/tmp/sample.db", dbPath)
			return []string{"posts", "users"}, nil
		},
	})

	result, err := service.HandleListTables(context.Background(), "/tmp/sample.db", "text")
	require.NoError(t, err)
	require.Equal(t, "posts\nusers", result)
}

func TestSQLiteService_HandleListTables_NormalJSON(t *testing.T) {
	t.Parallel()

	service := NewSQLiteService(&mockTableRepository{
		listTablesFunc: func(ctx context.Context, dbPath string) ([]string, error) {
			return []string{"posts", "users"}, nil
		},
	})

	result, err := service.HandleListTables(context.Background(), "/tmp/sample.db", "json")
	require.NoError(t, err)
	require.JSONEq(t, `["posts","users"]`, result)
}

func TestSQLiteService_HandleListTables_ErrorFromRepository(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("repository error")
	service := NewSQLiteService(&mockTableRepository{
		listTablesFunc: func(ctx context.Context, dbPath string) ([]string, error) {
			return nil, expectedErr
		},
	})

	_, err := service.HandleListTables(context.Background(), "/tmp/sample.db", "text")
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
}

func TestSQLiteService_HandleListTables_ErrorInvalidFormat(t *testing.T) {
	t.Parallel()

	service := NewSQLiteService(&mockTableRepository{
		listTablesFunc: func(ctx context.Context, dbPath string) ([]string, error) {
			return []string{"users"}, nil
		},
	})

	_, err := service.HandleListTables(context.Background(), "/tmp/sample.db", "xml")
	require.Error(t, err)
	require.Contains(t, err.Error(), "未対応の format")
}
