package config

import (
	"fmt"
	"os"
)

// Config はNotion同期CLIの設定を保持する構造体
type Config struct {
	Operation       string // 操作タイプ（必須）
	Token           string // Notionトークン
	ConID           string // コンテンツID（PageIDと排他）
	PageID          string // ページID（ConIDと排他）
	MarkdownContent string // マークダウンコンテンツ
	ToggleH1        bool   // H1ヘッダートグル
	ToggleH2        bool   // H2ヘッダートグル
	ToggleH3        bool   // H3ヘッダートグル
	EndpointURL     string // APIエンドポイントURL
	Date            string // 対象日付（YYYYMMDD形式、WebClip用）
	Title           string // 記事のタイトル（WebClip用）
	URL             string // 記事のURL（WebClip用）
	Help            bool   // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation, token, conID, pageID, markdownContent, endpointURL, date, title, url string, toggleH1, toggleH2, toggleH3 bool) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	if token == "" {
		return nil, fmt.Errorf("トークンが指定されていません")
	}

	if markdownContent == "" {
		return nil, fmt.Errorf("マークダウンコンテンツが指定されていません")
	}

	if endpointURL == "" {
		return nil, fmt.Errorf("エンドポイントURLが指定されていません")
	}

	// patch-web-clip操作の場合の追加バリデーション
	if operation == "patch-web-clip" {
		if title == "" {
			return nil, fmt.Errorf("patch-web-clip操作にはtitleが必要です")
		}
		if url == "" {
			return nil, fmt.Errorf("patch-web-clip操作にはurlが必要です")
		}
	} else {
		// patch操作の場合のConIDとPageIDの排他制御
		if conID == "" && pageID == "" {
			return nil, fmt.Errorf("con_id または page_id のいずれかを指定してください")
		}

		if conID != "" && pageID != "" {
			return nil, fmt.Errorf("con_id と page_id の両方を指定することはできません")
		}
	}

	return &Config{
		Operation:       operation,
		Token:           token,
		ConID:           conID,
		PageID:          pageID,
		MarkdownContent: markdownContent,
		ToggleH1:        toggleH1,
		ToggleH2:        toggleH2,
		ToggleH3:        toggleH3,
		EndpointURL:     endpointURL,
		Date:            date,
		Title:           title,
		URL:             url,
	}, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation       = ""
		token           = ""
		conID           = ""
		pageID          = ""
		markdownContent = ""
		toggleH1        = false
		toggleH2        = false
		toggleH3        = false
		endpointURL     = ""
		date            = ""
		title           = ""
		url             = ""
		help            = false
	)

	parser.StringVar(&operation, "operation", operation, "操作タイプ（必須）")
	parser.StringVar(&operation, "o", operation, "操作タイプの短縮形")

	parser.StringVar(&token, "token", token, "Notionトークン")
	parser.StringVar(&token, "t", token, "Notionトークンの短縮形")

	parser.StringVar(&conID, "con-id", conID, "コンテンツID")
	parser.StringVar(&conID, "c", conID, "コンテンツIDの短縮形")

	parser.StringVar(&pageID, "page-id", pageID, "ページID")
	parser.StringVar(&pageID, "p", pageID, "ページIDの短縮形")

	parser.StringVar(&markdownContent, "markdown", markdownContent, "マークダウンコンテンツ")
	parser.StringVar(&markdownContent, "m", markdownContent, "マークダウンコンテンツの短縮形")

	parser.BoolVar(&toggleH1, "toggle-h1", toggleH1, "H1ヘッダートグル")
	parser.BoolVar(&toggleH2, "toggle-h2", toggleH2, "H2ヘッダートグル")
	parser.BoolVar(&toggleH3, "toggle-h3", toggleH3, "H3ヘッダートグル")

	parser.StringVar(&endpointURL, "endpoint-url", endpointURL, "APIエンドポイントURL")
	parser.StringVar(&endpointURL, "u", endpointURL, "APIエンドポイントURLの短縮形")

	parser.StringVar(&date, "date", date, "対象日付（YYYYMMDD形式、WebClip用）")
	parser.StringVar(&date, "d", date, "対象日付の短縮形")

	parser.StringVar(&title, "title", title, "記事のタイトル（WebClip用）")

	parser.StringVar(&url, "url", url, "記事のURL（WebClip用）")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	// 残りの引数から値を取得（位置引数として）
	args := parser.Args()
	if len(args) >= 1 && token == "" {
		token = args[0]
	}
	if len(args) >= 2 && markdownContent == "" {
		markdownContent = args[1]
	}
	if len(args) >= 3 && endpointURL == "" {
		endpointURL = args[2]
	}

	return NewConfig(operation, token, conID, pageID, markdownContent, endpointURL, date, title, url, toggleH1, toggleH2, toggleH3)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `Notion同期CLIツール

使用方法:
  基本的な使用方法（ページID指定）:
    %s -operation "patch" -token "your_token" -page-id "page_id" -markdown "# Hello World" -endpoint-url "http://localhost:8080/task/patch"
    %s -o "patch" -t "your_token" -p "page_id" -m "# Hello World" -u "http://localhost:8080/task/patch"

  コンテンツID指定:
    %s -operation "patch" -token "your_token" -con-id "con_id" -markdown "# Hello World" -endpoint-url "http://localhost:8080/task/patch"
    %s -o "patch" -t "your_token" -c "con_id" -m "# Hello World" -u "http://localhost:8080/task/patch"

  ヘッダートグルオプション付き:
    %s -operation "patch" -token "your_token" -page-id "page_id" -markdown "# Hello World" -toggle-h1 -toggle-h2 -endpoint-url "http://localhost:8080/task/patch"

  WebClip操作:
    %s -operation "patch-web-clip" -token "your_token" -markdown "# Article Content" -title "Article Title" -url "https://example.com" -endpoint-url "http://localhost:8080/task/patch/web-clip"
    %s -o "patch-web-clip" -t "your_token" -m "# Article Content" -title "Article Title" -url "https://example.com" -date "20250817" -u "http://localhost:8080/task/patch/web-clip"

  位置引数での指定:
    %s "your_token" "# Hello World" "http://localhost:8080/task/patch" -operation "patch" -page-id "page_id"

オプション:
  -operation, -o    操作タイプ（必須）- 現在サポート: patch, patch-web-clip
  -token, -t        Notionトークン（必須）
  -con-id, -c       コンテンツID（page-idと排他、patch操作用）
  -page-id, -p      ページID（con-idと排他、patch操作用）
  -markdown, -m     マークダウンコンテンツ（必須）
  -toggle-h1        H1ヘッダートグル
  -toggle-h2        H2ヘッダートグル
  -toggle-h3        H3ヘッダートグル
  -endpoint-url, -u APIエンドポイントURL（必須）
  -date, -d         対象日付（YYYYMMDD形式、patch-web-clip操作用、省略可能）
  -title            記事のタイトル（patch-web-clip操作用、必須）
  -url              記事のURL（patch-web-clip操作用、必須）
  -help, -h         このヘルプを表示

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
