package usecases

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          Format Function Tests                             ##
// #==============================================================#

// TestFormatPK はformatPK関数をテストします
func TestFormatPK_Normal(t *testing.T) {
	// 単一主キー
	result := formatPK([]string{"id"})
	assert.Equal(t, "id", result)

	// 複合主キー
	result = formatPK([]string{"id", "email"})
	assert.Equal(t, "(id, email)", result)

	// 空の主キー
	result = formatPK([]string{})
	assert.Equal(t, "", result)
}

// TestFormatUK はformatUK関数をテストします
func TestFormatUK_Normal(t *testing.T) {
	// 単一一意キー
	uk1 := UniqueKey{Name: "uk1", Columns: []string{"email"}}
	result := formatUK([]UniqueKey{uk1})
	assert.Equal(t, "email", result)

	// 複合一意キー
	uk2 := UniqueKey{Name: "uk2", Columns: []string{"username", "domain"}}
	result = formatUK([]UniqueKey{uk2})
	assert.Equal(t, "(username, domain)", result)

	// 複数の一意キー
	result = formatUK([]UniqueKey{uk1, uk2})
	assert.Equal(t, "email; (username, domain)", result)

	// 空の一意キー
	result = formatUK([]UniqueKey{})
	assert.Equal(t, "", result)
}

// TestFormatFK はformatFK関数をテストします
func TestFormatFK_Normal(t *testing.T) {
	// 単一外部キー
	fk1 := ForeignKey{
		Name:       "fk1",
		Columns:    []string{"role_id"},
		RefTable:   "roles",
		RefColumns: []string{"id"},
	}
	result := formatFK([]ForeignKey{fk1})
	assert.Equal(t, "role_id -> roles.id", result)

	// 複合外部キー
	fk2 := ForeignKey{
		Name:       "fk2",
		Columns:    []string{"country_id", "region_id"},
		RefTable:   "locations",
		RefColumns: []string{"country_id", "region_id"},
	}
	result = formatFK([]ForeignKey{fk2})
	assert.Equal(t, "(country_id, region_id) -> locations.(country_id, region_id)", result)

	// 複数の外部キー
	result = formatFK([]ForeignKey{fk1, fk2})
	assert.Equal(t, "role_id -> roles.id; (country_id, region_id) -> locations.(country_id, region_id)", result)

	// 空の外部キー
	result = formatFK([]ForeignKey{})
	assert.Equal(t, "", result)
}

// TestFormatColumn はformatColumn関数をテストします
func TestFormatColumn_Normal(t *testing.T) {
	// NOT NULLカラム（デフォルト値あり）
	col1 := ColumnInfo{
		Name:       "id",
		Type:       "integer",
		IsNullable: "NO",
		Default:    sql.NullString{String: "nextval('users_id_seq'::regclass)", Valid: true},
		Comment:    "ID",
	}
	result := formatColumn(col1)
	assert.Equal(t, "- id: integer NOT NULL DEFAULT nextval('users_id_seq'::regclass) [ID]", result)

	// NULLカラム（デフォルト値なし）
	col2 := ColumnInfo{
		Name:       "name",
		Type:       "character varying",
		IsNullable: "YES",
		Default:    sql.NullString{Valid: false},
		Comment:    "名前",
	}
	result = formatColumn(col2)
	assert.Equal(t, "- name: character varying NULL [名前]", result)

	// コメントなし
	col3 := ColumnInfo{
		Name:       "status",
		Type:       "boolean",
		IsNullable: "NO",
		Default:    sql.NullString{String: "false", Valid: true},
		Comment:    "",
	}
	result = formatColumn(col3)
	assert.Equal(t, "- status: boolean NOT NULL DEFAULT false", result)
}

// TestFormatIndex はformatIndex関数をテストします
func TestFormatIndex_Normal(t *testing.T) {
	// 単一インデックス
	idx1 := IndexInfo{
		Name:    "users_email_idx",
		Columns: []string{"email"},
		Unique:  true,
	}
	result := formatIndex([]IndexInfo{idx1})
	assert.Equal(t, "email", result)

	// 複合インデックス
	idx2 := IndexInfo{
		Name:    "users_name_idx",
		Columns: []string{"first_name", "last_name"},
		Unique:  false,
	}
	result = formatIndex([]IndexInfo{idx2})
	assert.Equal(t, "(first_name, last_name)", result)

	// 複数のインデックス
	result = formatIndex([]IndexInfo{idx1, idx2})
	assert.Equal(t, "email; (first_name, last_name)", result)

	// 空のインデックス
	result = formatIndex([]IndexInfo{})
	assert.Equal(t, "", result)
}

// #==============================================================#
// ##          DefaultTemplateRenderer Tests                     ##
// #==============================================================#

// TestDefaultTemplateRenderer_RenderTableDetail はRenderTableDetail関数をテストします
func TestDefaultTemplateRenderer_RenderTableDetail_Normal(t *testing.T) {
	// Arrange
	renderer := &DefaultTemplateRenderer{}
	detail := &TableDetail{
		Name:    "users",
		Comment: "ユーザーテーブル",
		Columns: []ColumnInfo{
			{
				Name:       "id",
				Type:       "integer",
				IsNullable: "NO",
				Default:    sql.NullString{String: "nextval('users_id_seq'::regclass)", Valid: true},
				Comment:    "ID",
			},
			{
				Name:       "name",
				Type:       "character varying",
				IsNullable: "YES",
				Default:    sql.NullString{Valid: false},
				Comment:    "名前",
			},
		},
		PrimaryKeys: []string{"id"},
		UniqueKeys: []UniqueKey{
			{Name: "users_email_key", Columns: []string{"email"}},
		},
		ForeignKeys: []ForeignKey{
			{
				Name:       "users_role_id_fkey",
				Columns:    []string{"role_id"},
				RefTable:   "roles",
				RefColumns: []string{"id"},
			},
		},
		Indexes: []IndexInfo{
			{Name: "users_email_idx", Columns: []string{"email"}, Unique: true},
		},
	}

	// Act
	result, err := renderer.RenderTableDetail(detail)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "# テーブル: users - ユーザーテーブル")
	assert.Contains(t, result, "## カラム")
	assert.Contains(t, result, "- id: integer NOT NULL DEFAULT nextval('users_id_seq'::regclass) [ID]")
	assert.Contains(t, result, "- name: character varying NULL [名前]")
	assert.Contains(t, result, "## キー情報")
	assert.Contains(t, result, "[PK: id]")
	assert.Contains(t, result, "[UK: email]")
	assert.Contains(t, result, "[FK: role_id -> roles.id]")
	assert.Contains(t, result, "[INDEX: email]")
}

// TestDefaultTemplateRenderer_RenderTableDetail_MinimalData は最小限のデータでのテストです
func TestDefaultTemplateRenderer_RenderTableDetail_MinimalData(t *testing.T) {
	// Arrange
	renderer := &DefaultTemplateRenderer{}
	detail := &TableDetail{
		Name: "simple_table",
		Columns: []ColumnInfo{
			{
				Name:       "id",
				Type:       "integer",
				IsNullable: "NO",
				Default:    sql.NullString{Valid: false},
				Comment:    "",
			},
		},
		PrimaryKeys: []string{},
		UniqueKeys:  []UniqueKey{},
		ForeignKeys: []ForeignKey{},
		Indexes:     []IndexInfo{},
	}

	// Act
	result, err := renderer.RenderTableDetail(detail)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "# テーブル: simple_table")
	assert.Contains(t, result, "## カラム")
	assert.Contains(t, result, "- id: integer NOT NULL")
	assert.Contains(t, result, "## キー情報")
	// 空のキー情報は表示されない
	assert.NotContains(t, result, "[PK:")
	assert.NotContains(t, result, "[UK:")
	assert.NotContains(t, result, "[FK:")
	assert.NotContains(t, result, "[INDEX:")
}

// TestDefaultTemplateRenderer_RenderTableList はRenderTableList関数をテストします
func TestDefaultTemplateRenderer_RenderTableList_Normal(t *testing.T) {
	// Arrange
	renderer := &DefaultTemplateRenderer{}
	data := ListTablesData{
		Tables: []TableSummary{
			{
				Name:    "users",
				Comment: "ユーザーテーブル",
				PK:      []string{"id"},
				UK: []UniqueKey{
					{Name: "users_email_key", Columns: []string{"email"}},
				},
				FK: []ForeignKey{
					{
						Name:       "users_role_id_fkey",
						Columns:    []string{"role_id"},
						RefTable:   "roles",
						RefColumns: []string{"id"},
					},
				},
			},
			{
				Name:    "products",
				Comment: "商品テーブル",
				PK:      []string{"id"},
				UK:      []UniqueKey{},
				FK:      []ForeignKey{},
			},
		},
	}

	// Act
	result, err := renderer.RenderTableList(data)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "# データベースのテーブル一覧 (全2件)")
	assert.Contains(t, result, "- **users** — ユーザーテーブル")
	assert.Contains(t, result, "- PK: [id]")
	assert.Contains(t, result, "- UK: [email]")
	assert.Contains(t, result, "- FK: [role_id -> roles.id]")
	assert.Contains(t, result, "- **products** — 商品テーブル")
}

// TestDefaultTemplateRenderer_RenderTableList_EmptyData は空のデータでのテストです
func TestDefaultTemplateRenderer_RenderTableList_EmptyData(t *testing.T) {
	// Arrange
	renderer := &DefaultTemplateRenderer{}
	data := ListTablesData{
		Tables: []TableSummary{},
	}

	// Act
	result, err := renderer.RenderTableList(data)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "# データベースのテーブル一覧 (全0件)")
}

// TestDefaultTemplateRenderer_RenderTableList_MinimalTableData は最小限のテーブルデータでのテストです
func TestDefaultTemplateRenderer_RenderTableList_MinimalTableData(t *testing.T) {
	// Arrange
	renderer := &DefaultTemplateRenderer{}
	data := ListTablesData{
		Tables: []TableSummary{
			{
				Name:    "simple_table",
				Comment: "",
				PK:      []string{},
				UK:      []UniqueKey{},
				FK:      []ForeignKey{},
			},
		},
	}

	// Act
	result, err := renderer.RenderTableList(data)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "# データベースのテーブル一覧 (全1件)")
	assert.Contains(t, result, "- **simple_table** — ")
	// 空のキー情報は表示されない
	assert.NotContains(t, result, "- PK:")
	assert.NotContains(t, result, "- UK:")
	assert.NotContains(t, result, "- FK:")
}

// #==============================================================#
// ##          Error Case Tests                                  ##
// #==============================================================#

// TestDefaultTemplateRenderer_RenderTableDetail_NilData はnilデータでのエラーテストです
func TestDefaultTemplateRenderer_RenderTableDetail_NilData(t *testing.T) {
	// Arrange
	renderer := &DefaultTemplateRenderer{}

	// Act
	result, err := renderer.RenderTableDetail(nil)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "テンプレートの実行に失敗しました")
}

// #==============================================================#
// ##          Edge Case Tests                                   ##
// #==============================================================#

// TestFormatFunctions_EdgeCases はフォーマット関数のエッジケースをテストします
func TestFormatFunctions_EdgeCases(t *testing.T) {
	// formatPK with nil slice
	result := formatPK(nil)
	assert.Equal(t, "", result)

	// formatUK with nil slice
	result = formatUK(nil)
	assert.Equal(t, "", result)

	// formatFK with nil slice
	result = formatFK(nil)
	assert.Equal(t, "", result)

	// formatIndex with nil slice
	result = formatIndex(nil)
	assert.Equal(t, "", result)

	// formatColumn with empty values
	col := ColumnInfo{
		Name:       "",
		Type:       "",
		IsNullable: "",
		Default:    sql.NullString{Valid: false},
		Comment:    "",
	}
	result = formatColumn(col)
	assert.Equal(t, "- :  NOT NULL", result)
}

// TestComplexDataStructures は複雑なデータ構造でのテストです
func TestComplexDataStructures_Normal(t *testing.T) {
	// 複合主キー、複数の一意キー、複数の外部キーを持つテーブル
	detail := &TableDetail{
		Name:    "complex_table",
		Comment: "複雑なテーブル",
		Columns: []ColumnInfo{
			{
				Name:       "id1",
				Type:       "integer",
				IsNullable: "NO",
				Default:    sql.NullString{Valid: false},
				Comment:    "ID1",
			},
			{
				Name:       "id2",
				Type:       "integer",
				IsNullable: "NO",
				Default:    sql.NullString{Valid: false},
				Comment:    "ID2",
			},
		},
		PrimaryKeys: []string{"id1", "id2"},
		UniqueKeys: []UniqueKey{
			{Name: "uk1", Columns: []string{"col1"}},
			{Name: "uk2", Columns: []string{"col2", "col3"}},
		},
		ForeignKeys: []ForeignKey{
			{
				Name:       "fk1",
				Columns:    []string{"ref1"},
				RefTable:   "table1",
				RefColumns: []string{"id"},
			},
			{
				Name:       "fk2",
				Columns:    []string{"ref2", "ref3"},
				RefTable:   "table2",
				RefColumns: []string{"id1", "id2"},
			},
		},
		Indexes: []IndexInfo{
			{Name: "idx1", Columns: []string{"col1"}, Unique: false},
			{Name: "idx2", Columns: []string{"col2", "col3"}, Unique: true},
		},
	}

	renderer := &DefaultTemplateRenderer{}
	result, err := renderer.RenderTableDetail(detail)

	assert.NoError(t, err)
	assert.Contains(t, result, "[PK: (id1, id2)]")
	assert.Contains(t, result, "[UK: col1; (col2, col3)]")
	assert.Contains(t, result, "[FK: ref1 -> table1.id; (ref2, ref3) -> table2.(id1, id2)]")
	assert.Contains(t, result, "[INDEX: col1; (col2, col3)]")
}
