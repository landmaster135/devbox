package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 辞書の種類を表す定数
const (
	DictTypeKokugo   = "kokugo"   // 国語辞典
	DictTypeEiwa     = "eiwa"     // 英和辞典
	DictTypeWaei     = "waei"     // 和英辞典
	DictTypeDokuwa   = "dokuwa"   // 独和辞典
	DictTypeWadoku   = "wadoku"   // 和独辞典
	DictTypeRuigo    = "ruigo"    // 類語辞典
	DictTypeYojijuku = "yojijuku" // 四字熟語
)

// 単語情報のモデル
type Word struct {
	ConID           string     `json:"con_id"`
	Word            string     `json:"word"`
	DictType        string     `json:"dict_type"`
	Readings        []string   `json:"readings,omitempty"`
	Meanings        []string   `json:"meanings,omitempty"`
	Examples        []string   `json:"examples,omitempty"`
	Categories      []string   `json:"categories,omitempty"`
	LastMemorizedAt *time.Time `json:"last_memorized_at,omitempty"`
	Notice          string     `json:"notice,omitempty"`
}

// アプリケーション設定
type AppConfig struct {
	DictType    string
	IndexChar   string
	ApiURL      string
	ApiToken    string
	Debug       bool
	PrefixConID string
}

func main() {
	// コマンドライン引数の解析
	config := parseFlags()

	// フラグ検証
	if err := validateConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		printUsage()
		os.Exit(1)
	}

	// 索引ページのURLを構築
	indexURL := buildIndexURL(config.DictType, config.IndexChar)
	log.Printf("処理開始: %s", indexURL)

	// 索引ページから単語を取得して処理
	processIndex(config, indexURL, config.DictType)

	log.Println("処理完了")
}

// コマンドライン引数を解析する関数
func parseFlags() *AppConfig {
	config := &AppConfig{}

	// フラグ定義
	flag.StringVar(&config.DictType, "type", DictTypeKokugo,
		fmt.Sprintf("辞書の種類 (%s,%s,%s,%s,%s,%s,%s)",
			DictTypeKokugo, DictTypeEiwa, DictTypeWaei,
			DictTypeDokuwa, DictTypeWadoku,
			DictTypeRuigo, DictTypeYojijuku))
	flag.StringVar(&config.IndexChar, "char", "あ", "取得する索引の文字")
	flag.StringVar(&config.ApiURL, "api-url", "", "APIエンドポイントURL")
	flag.StringVar(&config.ApiToken, "api-token", "", "API認証トークン")
	flag.StringVar(&config.PrefixConID, "prefix", "VC", "con_idのプレフィックス")
	flag.BoolVar(&config.Debug, "debug", false, "デバッグモード")

	// ヘルプフラグ
	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "ヘルプメッセージを表示")
	flag.BoolVar(&showHelp, "help", false, "ヘルプメッセージを表示")

	flag.Parse()

	// ヘルプフラグが指定された場合
	if showHelp {
		printUsage()
		os.Exit(0)
	}

	if config.ApiURL == "" {
		config.ApiURL = fmt.Sprintf("%s%s", os.Getenv("DB_SERVER_URL"), os.Getenv("DB_SERVER_ENDPOINT_FOR_VOCABULARY_APPEND"))
	}

	return config
}

// 使用方法を表示する関数
func printUsage() {
	fmt.Println("使用方法: goo_dict_scraper [オプション]")
	fmt.Println("\nオプション:")
	flag.PrintDefaults()
	fmt.Println("\n例:")
	fmt.Println("  goo_dict_scraper -type kokugo -char あ")
	fmt.Println("  goo_dict_scraper -type eiwa -char a -api-url \"http://example.com/api/endpoint\"")
	fmt.Println("  goo_dict_scraper -type dokuwa -char a -debug")
}

// 設定を検証する関数
func validateConfig(config *AppConfig) error {
	// 辞書タイプの検証
	validTypes := map[string]bool{
		DictTypeKokugo:   true,
		DictTypeEiwa:     true,
		DictTypeWaei:     true,
		DictTypeDokuwa:   true,
		DictTypeWadoku:   true,
		DictTypeRuigo:    true,
		DictTypeYojijuku: true,
	}

	if _, ok := validTypes[config.DictType]; !ok {
		return fmt.Errorf("無効な辞書タイプ: %s", config.DictType)
	}

	// その他の検証ロジックをここに追加
	if config.IndexChar == "" {
		return fmt.Errorf("索引文字を指定してください")
	}

	// APIのURLの検証
	if config.ApiURL == "" {
		return fmt.Errorf("APIエンドポイントURLを指定してください")
	}

	return nil
}

// con_idを生成する関数
func generateConID(prefix string, index int) string {
	return fmt.Sprintf("%s%013d", prefix, index)
}

// 辞書タイプと索引文字に基づいて索引ページのURLを構築する関数
func buildIndexURL(dictType, indexChar string) string {
	encodedChar := url.QueryEscape(indexChar)

	switch dictType {
	case DictTypeKokugo:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/jn/index/%s/", encodedChar)
	case DictTypeEiwa:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/ej/index/%s/", encodedChar)
	case DictTypeWaei:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/je/index/%s/", encodedChar)
	case DictTypeDokuwa:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/dj/index/%s/", encodedChar)
	case DictTypeWadoku:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/jd/index/%s/", encodedChar)
	case DictTypeRuigo:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/thsrs/index/%s/", encodedChar)
	case DictTypeYojijuku:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/idiom/index/%s/", encodedChar)
	default:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/jn/index/%s/", encodedChar)
	}
}

// 索引ページから単語を取得して処理する関数
func processIndex(config *AppConfig, indexURL, dictType string) {
	log.Printf("索引ページの処理を開始: %s", indexURL)

	// 単語リストを取得
	words, nextPage, err := fetchWordsFromIndex(indexURL)
	if err != nil {
		log.Printf("索引からの単語取得に失敗: %v", err)
		return
	}

	log.Printf("索引から%d個の単語を取得しました", len(words))

	// 単語を処理
	var vocabularies []Word
	for i, word := range words {
		log.Printf("処理中 (%d/%d): %s", i+1, len(words), word)

		// 単語情報を取得
		info, err := fetchWordInfo(word, dictType, config.PrefixConID, i+1)
		if err != nil {
			log.Printf("単語「%s」の取得に失敗: %v", word, err)
			continue
		}

		// 単語情報を配列に追加
		vocabularies = append(vocabularies, info)

		log.Printf("単語「%s」を処理しました", word)

		// サーバーに負荷をかけないよう少し待機
		time.Sleep(2 * time.Second)
	}

	// APIにデータを送信
	if len(vocabularies) > 0 {
		err := sendToAPI(config, vocabularies)
		if err != nil {
			log.Printf("APIへのデータ送信に失敗: %v", err)
		} else {
			log.Printf("%d個の単語をAPIに送信しました", len(vocabularies))
		}
	}

	// 次のページがあれば処理
	if nextPage != "" {
		log.Printf("次のページを処理します: %s", nextPage)
		processIndex(config, nextPage, dictType)
	}
}

// APIにデータを送信する関数
func sendToAPI(config *AppConfig, vocabularies []Word) error {
	// リクエストデータの構築
	requestData := map[string]interface{}{
		"name":        fmt.Sprintf("Goo辞書スクレイピング - %s - %s", config.DictType, config.IndexChar),
		"description": fmt.Sprintf("Goo辞書から取得した%s辞典の「%s」で始まる単語", config.DictType, config.IndexChar),
		"data": map[string]interface{}{
			"vocabularies": vocabularies,
		},
	}

	// JSONに変換
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("JSONの生成に失敗: %w", err)
	}

	// HTTPリクエストの作成
	req, err := http.NewRequest("POST", config.ApiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗: %w", err)
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/json")
	if config.ApiToken != "" {
		req.Header.Set("Authorization", "Bearer "+config.ApiToken)
	}

	// リクエストの送信
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("APIリクエストの送信に失敗: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスの確認
	if resp.StatusCode >= 400 {
		return fmt.Errorf("APIエラー: ステータスコード %d", resp.StatusCode)
	}

	return nil
}

// 索引ページから単語リストと次ページのURLを取得する関数
func fetchWordsFromIndex(indexURL string) ([]string, string, error) {
	var words []string
	var nextPage string

	// HTTPリクエスト
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", indexURL, nil)
	if err != nil {
		return words, nextPage, err
	}

	// ユーザーエージェントの設定
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyDictionaryBot/1.0; +https://example.com)")

	resp, err := client.Do(req)
	if err != nil {
		return words, nextPage, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return words, nextPage, fmt.Errorf("HTTP エラー: %d", resp.StatusCode)
	}

	// レスポンスの解析
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return words, nextPage, err
	}

	// 単語リストの取得
	doc.Find(".content_list.idiom li a").Each(func(i int, s *goquery.Selection) {
		// タイトル部分を取得
		title := s.Find(".title").Text()

		if title == "" {
			// タイトルが見つからない場合はリンクのテキスト全体を使用
			title = s.Text()
		}

		// 括弧書きの部分を削除して単語だけを取得
		word := strings.Split(title, "【")[0]
		word = strings.TrimSpace(word)

		if word != "" {
			words = append(words, word)
		}
	})

	// 次のページのリンクを取得
	nextLink := doc.Find(".nav-paging .next a").AttrOr("href", "")
	if nextLink != "" {
		// 相対URLを絶対URLに変換
		nextPage = "https://dictionary.goo.ne.jp" + nextLink
	}

	return words, nextPage, nil
}

// 単語のページから詳細情報を取得する関数
func fetchWordInfo(word, dictType, prefixConID string, index int) (Word, error) {
	var info Word
	info.Word = word
	info.DictType = dictType
	info.ConID = generateConID(prefixConID, index)

	// URLエンコード
	encodedWord := url.QueryEscape(word)
	wordURL := buildWordURL(dictType, encodedWord)

	// HTTPリクエスト
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", wordURL, nil)
	if err != nil {
		return info, err
	}

	// ユーザーエージェントの設定
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyDictionaryBot/1.0; +https://example.com)")

	resp, err := client.Do(req)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return info, fmt.Errorf("HTTP エラー: %d", resp.StatusCode)
	}

	// レスポンスの解析
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return info, err
	}

	// 読み方の取得
	doc.Find(".content-box .text strong").Each(func(i int, s *goquery.Selection) {
		reading := strings.TrimSpace(s.Text())
		if reading != "" && !contains(info.Readings, reading) {
			info.Readings = append(info.Readings, reading)
		}
	})

	// 意味の取得
	doc.Find(".content-box .text").Each(func(i int, s *goquery.Selection) {
		meaning := strings.TrimSpace(s.Text())
		if meaning != "" && !contains(info.Meanings, meaning) {
			info.Meanings = append(info.Meanings, meaning)
		}
	})

	// 例文の取得
	doc.Find(".content-box .text-jp").Each(func(i int, s *goquery.Selection) {
		example := strings.TrimSpace(s.Text())
		if example != "" && !contains(info.Examples, example) {
			info.Examples = append(info.Examples, example)
		}
	})

	// カテゴリの取得
	doc.Find(".related_words_top_cat li a").Each(func(i int, s *goquery.Selection) {
		category := strings.TrimSpace(s.Text())
		if strings.HasPrefix(category, "#") {
			category = category[1:] // # を削除
		}
		if category != "" && !contains(info.Categories, category) {
			info.Categories = append(info.Categories, category)
		}
	})

	return info, nil
}

// 辞書タイプと単語に基づいて単語ページのURLを構築する関数
func buildWordURL(dictType, encodedWord string) string {
	switch dictType {
	case DictTypeKokugo:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/word/%s/", encodedWord)
	case DictTypeEiwa:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/ej/word/%s/", encodedWord)
	case DictTypeWaei:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/je/word/%s/", encodedWord)
	case DictTypeDokuwa:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/dj/word/%s/", encodedWord)
	case DictTypeWadoku:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/jd/word/%s/", encodedWord)
	case DictTypeRuigo:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/thsrs/word/%s/", encodedWord)
	case DictTypeYojijuku:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/idiom/word/%s/", encodedWord)
	default:
		return fmt.Sprintf("https://dictionary.goo.ne.jp/word/%s/", encodedWord)
	}
}

// スライスに特定の要素が含まれているかをチェックするヘルパー関数
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
