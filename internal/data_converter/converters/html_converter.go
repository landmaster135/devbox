package converters

import (
	"fmt"
	"strings"
)

// HTMLConverter はHTMLテーブル変換を行う構造体
type HTMLConverter struct{}

// NewHTMLConverter は新しいHTMLConverterを作成する
func NewHTMLConverter() *HTMLConverter {
	return &HTMLConverter{}
}

// ConvertToHTML は2次元配列をHTMLテーブルに変換する
// isTheadContained: 最初の行をヘッダーとして扱うかどうか
// textReplacingIfBlank: 空のセルに表示するテキスト
func (h *HTMLConverter) ConvertToHTML(values [][]string, isTheadContained bool, textReplacingIfBlank string) string {
	if len(values) == 0 {
		return "<table></table>"
	}

	if textReplacingIfBlank == "" {
		textReplacingIfBlank = "💩"
	}

	var html strings.Builder
	html.WriteString("<table>\n")

	isTheadAdded := false
	isTbodyAdded := false

	if !isTheadContained {
		isTheadAdded = true
		isTbodyAdded = true
	}

	for i, row := range values {
		isTh := false
		if i == 0 && isTheadContained {
			isTh = true
			if !isTheadAdded {
				html.WriteString("<thead>\n")
			}
		}

		if isTheadAdded && !isTbodyAdded {
			html.WriteString("<tbody>\n")
			isTbodyAdded = true
		}

		trContent := h.getTrByRow(row, isTh, textReplacingIfBlank)
		html.WriteString(trContent)
		html.WriteString("\n")

		if !isTheadAdded && i == 0 {
			html.WriteString("</thead>\n")
			isTheadAdded = true
		}
	}

	if len(values) != 1 && isTheadContained {
		html.WriteString("</tbody>\n")
	}

	html.WriteString("</table>")
	return html.String()
}

// getTrByRow は行データをtr要素に変換する
func (h *HTMLConverter) getTrByRow(row []string, isTh bool, textReplacingIfBlank string) string {
	var tr strings.Builder
	tagTd := "td"
	replacingText := textReplacingIfBlank

	if isTh {
		tagTd = "th"
		replacingText = ""
	}

	for _, cell := range row {
		element := h.closeByTag(tagTd, cell, replacingText)
		tr.WriteString(element)
	}

	return h.closeByTag("tr", tr.String(), "")
}

// closeByTag はタグで囲んだ要素を作成する
func (h *HTMLConverter) closeByTag(tag, innerText, textReplacingIfBlank string) string {
	if strings.HasPrefix(tag, "<") {
		panic("Initial of tag mustn't be \"<\".")
	}
	if strings.HasSuffix(tag, ">") {
		panic("End of tag mustn't be \">\".")
	}

	displayText := innerText
	if displayText == "" {
		displayText = textReplacingIfBlank
	}

	return fmt.Sprintf("<%s>%s</%s>", tag, displayText, tag)
}
