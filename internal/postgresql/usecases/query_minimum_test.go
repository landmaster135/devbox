package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          Test Helper Functions                             ##
// #==============================================================#

// setupMockDB はテスト用のモックデータベースをセットアップします
func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *PostgreSQLService) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
	}

	executor := &DefaultDatabaseExecutor{db: db}
	renderer := &DefaultTemplateRenderer{}
	marshaler := &DefaultJSONMarshaler{}
	tableDumper := NewTableDumper(executor)

	service := NewPostgreSQLServiceWithDependencies(
		executor,
		renderer,
		marshaler,
		tableDumper,
		"postgres://test:test@localhost/testdb",
		"postgres://test@localhost/testdb",
	)

	return db, mock, service
}

// newMockRows はモック用のRowsを作成します
func newMockRows(columns ...string) *sqlmock.Rows {
	return sqlmock.NewRows(columns)
}

// #==============================================================#
// ##          GetTablesMinimum Tests                            ##
// #==============================================================#

// TestGetTablesMinimum はGetTablesMinimum関数をテストします
func TestGetTablesMinimum_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// モックの期待値を設定
	rows := newMockRows("table_name").
		AddRow("users").
		AddRow("products").
		AddRow("orders")

	mock.ExpectQuery("SELECT table_name FROM information_schema.tables").
		WillReturnRows(rows)

	// 関数を実行
	ctx := context.Background()
	tables, err := service.GetTablesMinimum(ctx)

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, tables, 3)
	assert.Equal(t, "users", tables[0].Name)
	assert.Equal(t, "products", tables[1].Name)
	assert.Equal(t, "orders", tables[2].Name)

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestGetTablesMinimum_Error はGetTablesMinimumのエラーケースをテストします
func TestGetTablesMinimum_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT table_name FROM information_schema.tables").
		WillReturnError(fmt.Errorf("データベースエラー"))

	// 関数を実行
	ctx := context.Background()
	tables, err := service.GetTablesMinimum(ctx)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, tables)
	assert.Contains(t, err.Error(), "データベースエラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// #==============================================================#
// ##          GetTableSchemaMinimum Tests                       ##
// #==============================================================#

// TestGetTableSchemaMinimum はGetTableSchemaMinimum関数をテストします
func TestGetTableSchemaMinimum_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// モックの期待値を設定
	rows := newMockRows("column_name", "data_type").
		AddRow("id", "integer").
		AddRow("name", "character varying").
		AddRow("created_at", "timestamp")

	mock.ExpectQuery("SELECT column_name, data_type FROM information_schema.columns").
		WithArgs("public", "users").
		WillReturnRows(rows)

	// 関数を実行
	ctx := context.Background()
	columns, err := service.GetTableSchemaMinimum(ctx, "users")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, columns, 3)
	assert.Equal(t, "id", columns[0].Name)
	assert.Equal(t, "integer", columns[0].DataType)
	assert.Equal(t, "name", columns[1].Name)
	assert.Equal(t, "character varying", columns[1].DataType)
	assert.Equal(t, "created_at", columns[2].Name)
	assert.Equal(t, "timestamp", columns[2].DataType)

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestGetTableSchemaMinimum_Error はGetTableSchemaMinimumのエラーケースをテストします
func TestGetTableSchemaMinimum_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT column_name, data_type FROM information_schema.columns").
		WithArgs("public", "users").
		WillReturnError(fmt.Errorf("データベースエラー"))

	// 関数を実行
	ctx := context.Background()
	columns, err := service.GetTableSchemaMinimum(ctx, "users")

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, columns)
	assert.Contains(t, err.Error(), "データベースエラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// #==============================================================#
// ##          ExecuteQuery Tests                                ##
// #==============================================================#

// TestExecuteQuery はExecuteQuery関数をテストします
func TestExecuteQuery_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// トランザクションのモック
	mock.ExpectBegin()

	// クエリのモック
	rows := newMockRows("id", "name").
		AddRow(1, "John").
		AddRow(2, "Jane")

	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnRows(rows)

	mock.ExpectRollback()

	// 関数を実行
	ctx := context.Background()
	result, err := service.ExecuteQuery(ctx, "SELECT * FROM users")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "John", result[0]["name"])
	assert.Equal(t, "Jane", result[1]["name"])

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestExecuteQuery_BeginError はExecuteQueryのBeginエラーをテストします
func TestExecuteQuery_BeginError(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// トランザクション開始エラーのモック
	mock.ExpectBegin().WillReturnError(fmt.Errorf("トランザクション開始エラー"))

	// 関数を実行
	ctx := context.Background()
	result, err := service.ExecuteQuery(ctx, "SELECT * FROM users")

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "トランザクション開始エラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestExecuteQuery_QueryError はExecuteQueryのクエリエラーをテストします
func TestExecuteQuery_QueryError(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// トランザクションのモック
	mock.ExpectBegin()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnError(fmt.Errorf("クエリエラー"))

	mock.ExpectRollback()

	// 関数を実行
	ctx := context.Background()
	result, err := service.ExecuteQuery(ctx, "SELECT * FROM users")

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "クエリエラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// #==============================================================#
// ##          Handler Method Tests                              ##
// #==============================================================#

// TestHandleToQuery はHandleToQuery関数をテストします
func TestHandleToQuery_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// トランザクションのモック
	mock.ExpectBegin()

	// クエリのモック
	rows := newMockRows("id", "name").
		AddRow(1, "John").
		AddRow(2, "Jane")

	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnRows(rows)

	mock.ExpectRollback()

	// 関数を実行
	ctx := context.Background()
	result, err := service.HandleToQuery(ctx, "SELECT * FROM users")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "John", result[0]["name"])
	assert.Equal(t, "Jane", result[1]["name"])

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestHandleToQuery_Error はHandleToQueryのエラーケースをテストします
func TestHandleToQuery_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// トランザクションのモック
	mock.ExpectBegin()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnError(fmt.Errorf("クエリエラー"))

	mock.ExpectRollback()

	// 関数を実行
	ctx := context.Background()
	result, err := service.HandleToQuery(ctx, "SELECT * FROM users")

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "クエリエラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestHandleToListTablesMinimum はHandleToListTablesMinimum関数をテストします
func TestHandleToListTablesMinimum_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// モックの期待値を設定
	rows := newMockRows("table_name").
		AddRow("users").
		AddRow("products").
		AddRow("orders")

	mock.ExpectQuery("SELECT table_name FROM information_schema.tables").
		WillReturnRows(rows)

	// 関数を実行
	ctx := context.Background()
	result, err := service.HandleToListTablesMinimum(ctx)

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "users", result[0].Name)
	assert.Equal(t, "products", result[1].Name)
	assert.Equal(t, "orders", result[2].Name)

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestHandleToListTablesMinimum_Error はHandleToListTablesMinimumのエラーケースをテストします
func TestHandleToListTablesMinimum_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT table_name FROM information_schema.tables").
		WillReturnError(fmt.Errorf("データベースエラー"))

	// 関数を実行
	ctx := context.Background()
	result, err := service.HandleToListTablesMinimum(ctx)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "データベースエラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestHandleToGetTableSchemaMinimum はHandleToGetTableSchemaMinimum関数をテストします
func TestHandleToGetTableSchemaMinimum_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// モックの期待値を設定
	rows := newMockRows("column_name", "data_type").
		AddRow("id", "integer").
		AddRow("name", "character varying").
		AddRow("created_at", "timestamp")

	mock.ExpectQuery("SELECT column_name, data_type FROM information_schema.columns").
		WithArgs("public", "users").
		WillReturnRows(rows)

	// 関数を実行
	ctx := context.Background()
	result, err := service.HandleToGetTableSchemaMinimum(ctx, "users")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "id", result[0].Name)
	assert.Equal(t, "integer", result[0].DataType)
	assert.Equal(t, "name", result[1].Name)
	assert.Equal(t, "character varying", result[1].DataType)
	assert.Equal(t, "created_at", result[2].Name)
	assert.Equal(t, "timestamp", result[2].DataType)

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestHandleToGetTableSchemaMinimum_Error はHandleToGetTableSchemaMinimumのエラーケースをテストします
func TestHandleToGetTableSchemaMinimum_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT column_name, data_type FROM information_schema.columns").
		WithArgs("public", "users").
		WillReturnError(fmt.Errorf("データベースエラー"))

	// 関数を実行
	ctx := context.Background()
	result, err := service.HandleToGetTableSchemaMinimum(ctx, "users")

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "データベースエラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}
