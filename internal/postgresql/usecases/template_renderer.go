package usecases

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// #==============================================================#
// ##          Template Definitions                              ##
// #==============================================================#

// テンプレート関連の定義
const describeTableDetailTemplate = `# テーブル: {{.Name}}{{if .Comment}} - {{.Comment}}{{end}}

## カラム{{range .Columns}}
{{formatColumn .}}{{end}}

## キー情報{{if .PrimaryKeys}}
[PK: {{formatPK .PrimaryKeys}}]{{end}}{{if .UniqueKeys}}
[UK: {{formatUK .UniqueKeys}}]{{end}}{{if .ForeignKeys}}
[FK: {{formatFK .ForeignKeys}}]{{end}}{{if .Indexes}}
[INDEX: {{formatIndex .Indexes}}]{{end}}
`

// listTablesTemplate はテーブル一覧の出力フォーマット
const listTablesTemplate = `# データベースのテーブル一覧 (全{{len .Tables}}件)
フォーマット:
テーブル名 — テーブルコメント
  ├─ PK: [主キー]
  ├─ UK: [一意キー1; 一意キー2; ...]
  └─ FK: [外部キー → 参照先テーブル.カラム; ...]

{{range .Tables -}}
- **{{.Name}}** — {{.Comment}}
  {{if len .PK}}
  - PK: [{{formatPK .PK}}]{{end}}
  {{if len .UK}}
  - UK: [{{formatUK .UK}}]{{end}}
  {{if len .FK}}
  - FK: [{{formatFK .FK}}]{{end}}
{{end -}}
`

var funcMap = template.FuncMap{
	"formatPK":     formatPK,
	"formatUK":     formatUK,
	"formatFK":     formatFK,
	"formatColumn": formatColumn,
	"formatIndex":  formatIndex,
}

// #==============================================================#
// ##          Template Functions                                ##
// #==============================================================#

// formatPK は主キー情報をフォーマットします
func formatPK(pk []string) string {
	if len(pk) == 0 {
		return ""
	}
	pkStr := strings.Join(pk, ", ")
	if len(pk) > 1 {
		pkStr = fmt.Sprintf("(%s)", pkStr)
	}
	return pkStr
}

// formatUK は一意キー情報をフォーマットします
func formatUK(uk []UniqueKey) string {
	if len(uk) == 0 {
		return ""
	}
	var ukInfo []string
	for _, k := range uk {
		if len(k.Columns) > 1 {
			ukInfo = append(ukInfo, fmt.Sprintf("(%s)", strings.Join(k.Columns, ", ")))
		} else {
			ukInfo = append(ukInfo, strings.Join(k.Columns, ", "))
		}
	}
	return strings.Join(ukInfo, "; ")
}

// formatFK は外部キー情報をフォーマットします
func formatFK(fk []ForeignKey) string {
	if len(fk) == 0 {
		return ""
	}
	var fkInfo []string
	for _, k := range fk {
		colStr := strings.Join(k.Columns, ", ")
		refColStr := strings.Join(k.RefColumns, ", ")

		if len(k.Columns) > 1 {
			colStr = fmt.Sprintf("(%s)", colStr)
		}

		if len(k.RefColumns) > 1 {
			refColStr = fmt.Sprintf("(%s)", refColStr)
		}

		fkInfo = append(fkInfo, fmt.Sprintf("%s -> %s.%s",
			colStr,
			k.RefTable,
			refColStr))
	}
	return strings.Join(fkInfo, "; ")
}

// formatColumn はカラム情報をフォーマットします
func formatColumn(col ColumnInfo) string {
	nullable := "NOT NULL"
	if col.IsNullable == "YES" {
		nullable = "NULL"
	}

	defaultValue := ""
	if col.Default.Valid {
		defaultValue = fmt.Sprintf(" DEFAULT %s", col.Default.String)
	}

	comment := ""
	if col.Comment != "" {
		comment = fmt.Sprintf(" [%s]", col.Comment)
	}

	return fmt.Sprintf("- %s: %s %s%s%s",
		col.Name, col.Type, nullable, defaultValue, comment)
}

// formatIndex はインデックス情報をフォーマットします
func formatIndex(idx []IndexInfo) string {
	if len(idx) == 0 {
		return ""
	}
	var idxInfo []string
	for _, i := range idx {
		if len(i.Columns) > 1 {
			idxInfo = append(idxInfo, fmt.Sprintf("(%s)", strings.Join(i.Columns, ", ")))
		} else {
			idxInfo = append(idxInfo, strings.Join(i.Columns, ", "))
		}
	}
	return strings.Join(idxInfo, "; ")
}

// #==============================================================#
// ##          DefaultTemplateRenderer Implementation            ##
// #==============================================================#

// DefaultTemplateRenderer は標準のtext/templateを使用する実装
type DefaultTemplateRenderer struct{}

// RenderTableDetail はテーブル詳細情報をテンプレートでレンダリングします
func (r *DefaultTemplateRenderer) RenderTableDetail(detail *TableDetail) (string, error) {
	var output bytes.Buffer
	tmpl, err := template.New("describeTableDetail").Funcs(funcMap).Parse(describeTableDetailTemplate)
	if err != nil {
		return "", fmt.Errorf("テンプレートの解析に失敗しました: %w", err)
	}

	if err := tmpl.Execute(&output, detail); err != nil {
		return "", fmt.Errorf("テンプレートの実行に失敗しました: %w", err)
	}

	return output.String(), nil
}

// RenderTableList はテーブル一覧をテンプレートでレンダリングします
func (r *DefaultTemplateRenderer) RenderTableList(data ListTablesData) (string, error) {
	var output bytes.Buffer
	tmpl, err := template.New("listTables").Funcs(funcMap).Parse(listTablesTemplate)
	if err != nil {
		return "", fmt.Errorf("テンプレートの解析に失敗しました: %w", err)
	}

	if err := tmpl.Execute(&output, data); err != nil {
		return "", fmt.Errorf("テンプレートの実行に失敗しました: %w", err)
	}

	return output.String(), nil
}
