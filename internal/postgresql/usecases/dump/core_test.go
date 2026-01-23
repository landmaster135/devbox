package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	dbExecutor "github.com/landmaster135/devbox/internal/postgresql/domain/executor"
	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
	templateRenderer "github.com/landmaster135/devbox/internal/postgresql/usecases/template_renderer"
)

// #==============================================================#
// ##          Test Helper Types & Functions                     ##
// #==============================================================#

type tableRenderer interface {
	RenderTableDetail(detail *model.TableDetail) (string, error)
	RenderTableList(data model.ListTablesData) (string, error)
}

type testDumpService struct {
	exec     dbExecutor.DatabaseExecutor
	renderer tableRenderer
}

// setupMockDB はテスト用のモックデータベースをセットアップします
func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *testDumpService) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("モックデータベースの作成に失敗しました: %v", err)
	}

	executor := &dbExecutor.DefaultDatabaseExecutor{DB: db}
	service := &testDumpService{
		exec:     executor,
		renderer: &templateRenderer.DefaultTemplateRenderer{},
	}

	return db, mock, service
}

// newMockRows はモック行セットを作成します
func newMockRows(columns ...string) *sqlmock.Rows {
	return sqlmock.NewRows(columns)
}

// #==============================================================#
// ##          FetchTableWithComments Tests                      ##
// #==============================================================#

// TestFetchTableWithComments はfetchTableWithComments関数をテストします
func TestFetchTableWithComments_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// モックの期待値を設定
	rows := newMockRows("table_name", "table_comment").
		AddRow("users", "ユーザーテーブル")

	mock.ExpectQuery("SELECT t.table_name, COALESCE").
		WithArgs("public", "users").
		WillReturnRows(rows)

	// 関数を実行
	ctx := context.Background()
	table, err := fetchTableWithComments(ctx, service.exec, "users")

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, "users", table.Name)
	assert.Equal(t, "ユーザーテーブル", table.Comment)

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestFetchTableWithComments_Error はfetchTableWithCommentsのエラーケースをテストします
func TestFetchTableWithComments_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT t.table_name, COALESCE").
		WithArgs("public", "users").
		WillReturnError(fmt.Errorf("データベースエラー"))

	// 関数を実行
	ctx := context.Background()
	table, err := fetchTableWithComments(ctx, service.exec, "users")

	// アサーション
	assert.Error(t, err)
	assert.Empty(t, table.Name)
	assert.Contains(t, err.Error(), "データベースエラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// #==============================================================#
// ##          FetchPrimaryKeys Tests                            ##
// #==============================================================#

// TestFetchPrimaryKeys はfetchPrimaryKeys関数をテストします
func TestFetchPrimaryKeys_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// モックの期待値を設定
	rows := newMockRows("attname").
		AddRow("id").
		AddRow("email")

	mock.ExpectQuery("SELECT a.attname FROM pg_index").
		WithArgs("\"public\".\"users\"").
		WillReturnRows(rows)

	// 関数を実行
	ctx := context.Background()
	primaryKeys, err := fetchPrimaryKeys(ctx, service.exec, "users")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, primaryKeys, 2)
	assert.Equal(t, "id", primaryKeys[0])
	assert.Equal(t, "email", primaryKeys[1])

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestFetchPrimaryKeys_Error はfetchPrimaryKeysのエラーケースをテストします
func TestFetchPrimaryKeys_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT a.attname FROM pg_index").
		WithArgs("\"public\".\"users\"").
		WillReturnError(fmt.Errorf("データベースエラー"))

	// 関数を実行
	ctx := context.Background()
	primaryKeys, err := fetchPrimaryKeys(ctx, service.exec, "users")

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, primaryKeys)
	assert.Contains(t, err.Error(), "データベースエラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// #==============================================================#
// ##          FetchUniqueKeys Tests                             ##
// #==============================================================#

// TestFetchUniqueKeys はfetchUniqueKeys関数をテストします
func TestFetchUniqueKeys_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// モックの期待値を設定
	rows := newMockRows("constraint_name", "column_name").
		AddRow("users_email_key", "email").
		AddRow("users_username_key", "username")

	mock.ExpectQuery("SELECT c.conname AS constraint_name, a.attname AS column_name FROM pg_constraint").
		WithArgs("\"public\".\"users\"").
		WillReturnRows(rows)

	// 関数を実行
	ctx := context.Background()
	uniqueKeys, err := fetchUniqueKeys(ctx, service.exec, "users")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, uniqueKeys, 2)
	assert.Equal(t, "users_email_key", uniqueKeys[0].Name)
	assert.Equal(t, "email", uniqueKeys[0].Columns[0])
	assert.Equal(t, "users_username_key", uniqueKeys[1].Name)
	assert.Equal(t, "username", uniqueKeys[1].Columns[0])

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestFetchUniqueKeys_Error はfetchUniqueKeysのエラーケースをテストします
func TestFetchUniqueKeys_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT c.conname AS constraint_name, a.attname AS column_name FROM pg_constraint").
		WithArgs("\"public\".\"users\"").
		WillReturnError(fmt.Errorf("データベースエラー"))

	// 関数を実行
	ctx := context.Background()
	uniqueKeys, err := fetchUniqueKeys(ctx, service.exec, "users")

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, uniqueKeys)
	assert.Contains(t, err.Error(), "データベースエラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// #==============================================================#
// ##          FetchForeignKeys Tests                            ##
// #==============================================================#

// TestFetchForeignKeys はfetchForeignKeys関数をテストします
func TestFetchForeignKeys_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// モックの期待値を設定
	rows := newMockRows("constraint_name", "column_name", "referenced_table", "referenced_column").
		AddRow("users_role_id_fkey", "role_id", "roles", "id")

	mock.ExpectQuery("SELECT c.conname AS constraint_name, a.attname AS column_name").
		WithArgs("\"public\".\"users\"").
		WillReturnRows(rows)

	// 関数を実行
	ctx := context.Background()
	foreignKeys, err := fetchForeignKeys(ctx, service.exec, "users")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, foreignKeys, 1)
	assert.Equal(t, "users_role_id_fkey", foreignKeys[0].Name)
	assert.Equal(t, "role_id", foreignKeys[0].Columns[0])
	assert.Equal(t, "roles", foreignKeys[0].RefTable)
	assert.Equal(t, "id", foreignKeys[0].RefColumns[0])

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestFetchForeignKeys_Error はfetchForeignKeysのエラーケースをテストします
func TestFetchForeignKeys_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT c.conname AS constraint_name, a.attname AS column_name").
		WithArgs("\"public\".\"users\"").
		WillReturnError(fmt.Errorf("データベースエラー"))

	// 関数を実行
	ctx := context.Background()
	foreignKeys, err := fetchForeignKeys(ctx, service.exec, "users")

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, foreignKeys)
	assert.Contains(t, err.Error(), "データベースエラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// #==============================================================#
// ##          FetchTableColumns Tests                           ##
// #==============================================================#

// TestFetchTableColumns はfetchTableColumns関数をテストします
func TestFetchTableColumns_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// モックの期待値を設定
	rows := newMockRows("column_name", "data_type", "is_nullable", "column_default", "column_comment").
		AddRow("id", "integer", "NO", sql.NullString{String: "nextval('users_id_seq'::regclass)", Valid: true}, "ID").
		AddRow("name", "character varying", "YES", sql.NullString{Valid: false}, "名前")

	mock.ExpectQuery("SELECT c.column_name, c.data_type, c.is_nullable").
		WithArgs("public", "users").
		WillReturnRows(rows)

	// 関数を実行
	ctx := context.Background()
	columns, err := fetchTableColumns(ctx, service.exec, "users")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, columns, 2)
	assert.Equal(t, "id", columns[0].Name)
	assert.Equal(t, "integer", columns[0].Type)
	assert.Equal(t, "NO", columns[0].IsNullable)
	assert.Equal(t, "ID", columns[0].Comment)
	assert.Equal(t, "name", columns[1].Name)
	assert.Equal(t, "character varying", columns[1].Type)
	assert.Equal(t, "YES", columns[1].IsNullable)
	assert.Equal(t, "名前", columns[1].Comment)

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestFetchTableColumns_Error はfetchTableColumnsのエラーケースをテストします
func TestFetchTableColumns_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT c.column_name, c.data_type, c.is_nullable").
		WithArgs("public", "users").
		WillReturnError(fmt.Errorf("データベースエラー"))

	// 関数を実行
	ctx := context.Background()
	columns, err := fetchTableColumns(ctx, service.exec, "users")

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
// ##          FetchTableIndexes Tests                           ##
// #==============================================================#

// TestFetchTableIndexes はfetchTableIndexes関数をテストします
func TestFetchTableIndexes_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// モックの期待値を設定
	rows := newMockRows("index_name", "column_name", "is_unique").
		AddRow("users_email_idx", "email", true).
		AddRow("users_name_idx", "name", false)

	mock.ExpectQuery("SELECT i.relname AS index_name, a.attname AS column_name").
		WithArgs("users", "public").
		WillReturnRows(rows)

	// 関数を実行
	ctx := context.Background()
	indexes, err := fetchTableIndexes(ctx, service.exec, "users")

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, indexes, 2)
	assert.Equal(t, "users_email_idx", indexes[0].Name)
	assert.Equal(t, "email", indexes[0].Columns[0])
	assert.True(t, indexes[0].Unique)
	assert.Equal(t, "users_name_idx", indexes[1].Name)
	assert.Equal(t, "name", indexes[1].Columns[0])
	assert.False(t, indexes[1].Unique)

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestFetchTableIndexes_Error はfetchTableIndexesのエラーケースをテストします
func TestFetchTableIndexes_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// クエリエラーのモック
	mock.ExpectQuery("SELECT i.relname AS index_name, a.attname AS column_name").
		WithArgs("users", "public").
		WillReturnError(fmt.Errorf("データベースエラー"))

	// 関数を実行
	ctx := context.Background()
	indexes, err := fetchTableIndexes(ctx, service.exec, "users")

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, indexes)
	assert.Contains(t, err.Error(), "データベースエラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// #==============================================================#
// ##          GetTableDetail Tests                              ##
// #==============================================================#

// TestGetTableDetail はGetTableDetail関数をテストします
func TestGetTableDetail_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// fetchTableWithCommentsのモック
	tableRows := newMockRows("table_name", "table_comment").
		AddRow("users", "ユーザーテーブル")
	mock.ExpectQuery("SELECT t.table_name, COALESCE").
		WithArgs("public", "users").
		WillReturnRows(tableRows)

	// fetchPrimaryKeysのモック
	pkRows := newMockRows("attname").
		AddRow("id")
	mock.ExpectQuery("SELECT a.attname FROM pg_index").
		WithArgs("\"public\".\"users\"").
		WillReturnRows(pkRows)

	// fetchUniqueKeysのモック
	ukRows := newMockRows("constraint_name", "column_name").
		AddRow("users_email_key", "email")
	mock.ExpectQuery("SELECT c.conname AS constraint_name, a.attname AS column_name FROM pg_constraint").
		WithArgs("\"public\".\"users\"").
		WillReturnRows(ukRows)

	// fetchForeignKeysのモック
	fkRows := newMockRows("constraint_name", "column_name", "referenced_table", "referenced_column").
		AddRow("users_role_id_fkey", "role_id", "roles", "id")
	mock.ExpectQuery("SELECT c.conname AS constraint_name, a.attname AS column_name").
		WithArgs("\"public\".\"users\"").
		WillReturnRows(fkRows)

	// fetchTableColumnsのモック
	colRows := newMockRows("column_name", "data_type", "is_nullable", "column_default", "column_comment").
		AddRow("id", "integer", "NO", sql.NullString{String: "nextval('users_id_seq'::regclass)", Valid: true}, "ID").
		AddRow("name", "character varying", "YES", sql.NullString{Valid: false}, "名前")
	mock.ExpectQuery("SELECT c.column_name, c.data_type, c.is_nullable").
		WithArgs("public", "users").
		WillReturnRows(colRows)

	// fetchTableIndexesのモック
	idxRows := newMockRows("index_name", "column_name", "is_unique").
		AddRow("users_email_idx", "email", true)
	mock.ExpectQuery("SELECT i.relname AS index_name, a.attname AS column_name").
		WithArgs("users", "public").
		WillReturnRows(idxRows)

	// 関数を実行
	ctx := context.Background()
	detail, err := GetTableDetail(ctx, service.exec, "users")

	// アサーション
	assert.NoError(t, err)
	assert.NotNil(t, detail)
	assert.Equal(t, "users", detail.Name)
	assert.Equal(t, "ユーザーテーブル", detail.Comment)
	assert.Len(t, detail.PrimaryKeys, 1)
	assert.Equal(t, "id", detail.PrimaryKeys[0])
	assert.Len(t, detail.UniqueKeys, 1)
	assert.Equal(t, "users_email_key", detail.UniqueKeys[0].Name)
	assert.Len(t, detail.ForeignKeys, 1)
	assert.Equal(t, "users_role_id_fkey", detail.ForeignKeys[0].Name)
	assert.Len(t, detail.Columns, 2)
	assert.Equal(t, "id", detail.Columns[0].Name)
	assert.Len(t, detail.Indexes, 1)
	assert.Equal(t, "users_email_idx", detail.Indexes[0].Name)

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestGetTableDetail_FetchTableWithCommentsError はGetTableDetailのfetchTableWithCommentsエラーをテストします
func TestGetTableDetail_FetchTableWithCommentsError(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// fetchTableWithCommentsのエラーモック
	mock.ExpectQuery("SELECT t.table_name, COALESCE").
		WithArgs("public", "users").
		WillReturnError(fmt.Errorf("テーブル情報の取得エラー"))

	// 関数を実行
	ctx := context.Background()
	detail, err := GetTableDetail(ctx, service.exec, "users")

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, detail)
	assert.Contains(t, err.Error(), "テーブル情報の取得に失敗しました")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// #==============================================================#
// ##          GetAllTableSummaries Tests                        ##
// #==============================================================#

// TestGetAllTableSummaries はGetAllTableSummaries関数をテストします
func TestGetAllTableSummaries_Normal(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// テーブル一覧のモック
	tableRows := newMockRows("table_name", "table_comment").
		AddRow("users", "ユーザーテーブル").
		AddRow("products", "商品テーブル")
	mock.ExpectQuery("SELECT t.table_name, COALESCE").
		WillReturnRows(tableRows)

	// users テーブルの主キー情報のモック
	usersPkRows := newMockRows("attname").
		AddRow("id")
	mock.ExpectQuery("SELECT a.attname FROM pg_index").
		WithArgs("\"public\".\"users\"").
		WillReturnRows(usersPkRows)

	// users テーブルの一意キー情報のモック
	usersUkRows := newMockRows("constraint_name", "column_name").
		AddRow("users_email_key", "email")
	mock.ExpectQuery("SELECT c.conname AS constraint_name, a.attname AS column_name FROM pg_constraint").
		WithArgs("\"public\".\"users\"").
		WillReturnRows(usersUkRows)

	// users テーブルの外部キー情報のモック
	usersFkRows := newMockRows("constraint_name", "column_name", "referenced_table", "referenced_column").
		AddRow("users_role_id_fkey", "role_id", "roles", "id")
	mock.ExpectQuery("SELECT c.conname AS constraint_name, a.attname AS column_name").
		WithArgs("\"public\".\"users\"").
		WillReturnRows(usersFkRows)

	// products テーブルの主キー情報のモック
	productsPkRows := newMockRows("attname").
		AddRow("id")
	mock.ExpectQuery("SELECT a.attname FROM pg_index").
		WithArgs("\"public\".\"products\"").
		WillReturnRows(productsPkRows)

	// products テーブルの一意キー情報のモック
	productsUkRows := newMockRows("constraint_name", "column_name").
		AddRow("products_code_key", "code")
	mock.ExpectQuery("SELECT c.conname AS constraint_name, a.attname AS column_name FROM pg_constraint").
		WithArgs("\"public\".\"products\"").
		WillReturnRows(productsUkRows)

	// products テーブルの外部キー情報のモック
	productsFkRows := newMockRows("constraint_name", "column_name", "referenced_table", "referenced_column").
		AddRow("products_category_id_fkey", "category_id", "categories", "id")
	mock.ExpectQuery("SELECT c.conname AS constraint_name, a.attname AS column_name").
		WithArgs("\"public\".\"products\"").
		WillReturnRows(productsFkRows)

	// 関数を実行
	ctx := context.Background()
	tables, err := GetAllTableSummaries(ctx, service.exec)

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, tables, 2)
	assert.Equal(t, "users", tables[0].Name)
	assert.Equal(t, "ユーザーテーブル", tables[0].Comment)
	assert.Len(t, tables[0].PK, 1)
	assert.Equal(t, "id", tables[0].PK[0])
	assert.Len(t, tables[0].UK, 1)
	assert.Equal(t, "users_email_key", tables[0].UK[0].Name)
	assert.Len(t, tables[0].FK, 1)
	assert.Equal(t, "users_role_id_fkey", tables[0].FK[0].Name)

	assert.Equal(t, "products", tables[1].Name)
	assert.Equal(t, "商品テーブル", tables[1].Comment)
	assert.Len(t, tables[1].PK, 1)
	assert.Equal(t, "id", tables[1].PK[0])
	assert.Len(t, tables[1].UK, 1)
	assert.Equal(t, "products_code_key", tables[1].UK[0].Name)
	assert.Len(t, tables[1].FK, 1)
	assert.Equal(t, "products_category_id_fkey", tables[1].FK[0].Name)

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// TestGetAllTableSummaries_Error はGetAllTableSummariesのエラーケースをテストします
func TestGetAllTableSummaries_Error(t *testing.T) {
	// モックデータベースをセットアップ
	db, mock, service := setupMockDB(t)
	defer db.Close()

	// テーブル一覧のエラーモック
	mock.ExpectQuery("SELECT t.table_name, COALESCE").
		WillReturnError(fmt.Errorf("テーブル一覧の取得エラー"))

	// 関数を実行
	ctx := context.Background()
	tables, err := GetAllTableSummaries(ctx, service.exec)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, tables)
	assert.Contains(t, err.Error(), "テーブル一覧の取得エラー")

	// モックの期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}
