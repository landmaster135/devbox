package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// OperationAuthURL は認可URLを生成する操作名です。
	OperationAuthURL = "auth-url"
	// OperationRequestToken は認可コードからトークンを交換する操作名です。
	OperationRequestToken = "request-token"
	// OperationRefreshToken はリフレッシュトークンでアクセストークンを再取得する操作名です。
	OperationRefreshToken = "refresh-token"
	// OperationDailySummary は日次サマリ取得操作名です。
	OperationDailySummary = "daily-summary"
	// DefaultOperation は CLI がサポートする標準操作名です。
	DefaultOperation = OperationDailySummary
	// DefaultOutputFormat は CLI の標準出力形式です。
	DefaultOutputFormat = "json"
	// defaultTimeout は API 呼び出しに利用するデフォルトのタイムアウトです。
	defaultTimeout = 15 * time.Second
)

// Config は CLI 実行時の設定値を保持します。
type Config struct {
	Operation         string
	AccessToken       string
	ClientID          string
	ClientSecret      string
	RefreshToken      string
	AuthorizationCode string
	RedirectURI       string
	Scope             string
	State             string
	Mode              string
	ResponseType      string
	UserID            int64
	StartDate         time.Time
	EndDate           time.Time
	MeasureTypes      []int
	IncludeActivity   bool
	OutputFormat      string
	Timeout           time.Duration
	Help              bool
}

// FlagParser は flag.FlagSet 相当の機能を抽象化します。
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	IntVar(p *int, name string, value int, usage string)
	Int64Var(p *int64, name string, value int64, usage string)
	DurationVar(p *time.Duration, name string, value time.Duration, usage string)
	Parse() error
	Args() []string
}

// standardFlagParser は標準 flag.FlagSet のラッパーです。
type standardFlagParser struct {
	fs *flag.FlagSet
}

// NewStandardFlagParser は標準フラグパーサーを生成します。
func NewStandardFlagParser() FlagParser {
	return &standardFlagParser{fs: flag.NewFlagSet(os.Args[0], flag.ContinueOnError)}
}

func (p *standardFlagParser) StringVar(dest *string, name, value, usage string) {
	p.fs.StringVar(dest, name, value, usage)
}

func (p *standardFlagParser) BoolVar(dest *bool, name string, value bool, usage string) {
	p.fs.BoolVar(dest, name, value, usage)
}

func (p *standardFlagParser) IntVar(dest *int, name string, value int, usage string) {
	p.fs.IntVar(dest, name, value, usage)
}

func (p *standardFlagParser) Int64Var(dest *int64, name string, value int64, usage string) {
	p.fs.Int64Var(dest, name, value, usage)
}

func (p *standardFlagParser) DurationVar(dest *time.Duration, name string, value time.Duration, usage string) {
	p.fs.DurationVar(dest, name, value, usage)
}

func (p *standardFlagParser) Parse() error {
	return p.fs.Parse(os.Args[1:])
}

func (p *standardFlagParser) Args() []string {
	return p.fs.Args()
}

// ParseFlags は標準パーサーで CLI フラグを解析します。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定したパーサーで CLI フラグを解析します。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation         = DefaultOperation
		accessToken       string
		clientID          string
		clientSecret      string
		refreshToken      string
		authorizationCode string
		redirectURI       string
		scope             = "user.metrics,user.activity"
		state             string
		mode              string
		responseType      = "code"
		userID            int64
		startDateStr      string
		endDateStr        string
		measureTypesStr   string
		includeActivity   = true
		outputFormat      = DefaultOutputFormat
		timeout           = defaultTimeout
		help              bool
	)

	parser.StringVar(&operation, "operation", operation, "実行する操作 (現在は daily-summary のみ対応)")
	parser.StringVar(&operation, "op", operation, "operation の短縮指定")
	parser.StringVar(&accessToken, "access-token", accessToken, "Withings API のアクセストークン (client-id/client-secret を指定しない場合は必須)")
	parser.StringVar(&accessToken, "token", accessToken, "access-token の短縮指定")
	parser.StringVar(&clientID, "client-id", clientID, "Withings OAuth クライアント ID (access-token を省略する場合に必須)")
	parser.StringVar(&clientSecret, "client-secret", clientSecret, "Withings OAuth クライアントシークレット (access-token を省略する場合に必須)")
	parser.StringVar(&refreshToken, "refresh-token", refreshToken, "Withings OAuth リフレッシュトークン (access-token を省略する場合に必須)")
	parser.StringVar(&authorizationCode, "authorization-code", authorizationCode, "Withings OAuth 認可コード (初回交換時に使用)")
	parser.StringVar(&authorizationCode, "auth-code", authorizationCode, "authorization-code の短縮指定")
	parser.StringVar(&redirectURI, "redirect-uri", redirectURI, "Withings OAuth 認可コード取得時に使用したリダイレクト URI")
	parser.StringVar(&scope, "scope", scope, "認可リクエストで要求するスコープ (カンマ区切り)")
	parser.StringVar(&state, "state", state, "認可リクエストに付加する state 値")
	parser.StringVar(&mode, "mode", mode, "Withings デモモード等を指定する mode パラメータ")
	parser.StringVar(&responseType, "response-type", responseType, "認可リクエストの response_type (デフォルト: code)")
	parser.Int64Var(&userID, "user-id", userID, "Withings ユーザー ID (必須)")
	parser.StringVar(&startDateStr, "start-date", startDateStr, "取得対象の開始日 (YYYY-MM-DD, 必須)")
	parser.StringVar(&startDateStr, "start", startDateStr, "start-date の短縮指定")
	parser.StringVar(&endDateStr, "end-date", endDateStr, "取得対象の終了日 (YYYY-MM-DD, 省略時は開始日と同じ)")
	parser.StringVar(&endDateStr, "end", endDateStr, "end-date の短縮指定")
	parser.StringVar(&measureTypesStr, "measure-types", measureTypesStr, "取得する measure type (カンマ区切り, 例: weight,fat_mass,diastolic)")
	parser.BoolVar(&includeActivity, "include-activity", includeActivity, "日次活動サマリを同時に取得するか")
	parser.StringVar(&outputFormat, "output", outputFormat, "出力フォーマット (json)")
	parser.DurationVar(&timeout, "timeout", timeout, "API 呼び出しのタイムアウト (例: 20s, 1m)")
	parser.BoolVar(&help, "help", help, "このヘルプを表示")
	parser.BoolVar(&help, "h", help, "help の短縮指定")

	if err := parser.Parse(); err != nil {
		return nil, err
	}

	// 位置引数のフォールバック: token, user-id, start, end
	args := parser.Args()
	if len(args) > 0 && accessToken == "" {
		accessToken = args[0]
	}
	if len(args) > 1 && userID == 0 {
		if parsed, err := strconv.ParseInt(args[1], 10, 64); err == nil {
			userID = parsed
		} else {
			return nil, fmt.Errorf("ユーザー ID の解析に失敗しました: %v", err)
		}
	}
	if len(args) > 2 && startDateStr == "" {
		startDateStr = args[2]
	}
	if len(args) > 3 && endDateStr == "" {
		endDateStr = args[3]
	}

	if help {
		return &Config{Help: true}, nil
	}

	operation = strings.TrimSpace(strings.ToLower(operation))
	if operation == "" {
		operation = DefaultOperation
	}

	if timeout <= 0 {
		timeout = defaultTimeout
	}

	cfg := &Config{
		Operation:         operation,
		AccessToken:       strings.TrimSpace(accessToken),
		ClientID:          strings.TrimSpace(clientID),
		ClientSecret:      strings.TrimSpace(clientSecret),
		RefreshToken:      strings.TrimSpace(refreshToken),
		AuthorizationCode: strings.TrimSpace(authorizationCode),
		RedirectURI:       strings.TrimSpace(redirectURI),
		Scope:             normalizeCommaSeparated(scope),
		State:             strings.TrimSpace(state),
		Mode:              strings.TrimSpace(mode),
		ResponseType:      strings.TrimSpace(strings.ToLower(responseType)),
		IncludeActivity:   includeActivity,
		OutputFormat:      strings.ToLower(strings.TrimSpace(outputFormat)),
		Timeout:           timeout,
	}

	switch cfg.Operation {
	case OperationDailySummary:
		if cfg.AccessToken == "" {
			return nil, errors.New("daily-summary では access-token を指定してください")
		}
		if userID <= 0 {
			return nil, errors.New("user-id は正の整数で指定してください")
		}
		if strings.TrimSpace(startDateStr) == "" {
			return nil, errors.New("start-date が指定されていません")
		}
		if strings.TrimSpace(endDateStr) == "" {
			endDateStr = startDateStr
		}

		startDate, err := parseDate(startDateStr)
		if err != nil {
			return nil, fmt.Errorf("start-date の解析に失敗しました: %w", err)
		}
		endDate, err := parseDate(endDateStr)
		if err != nil {
			return nil, fmt.Errorf("end-date の解析に失敗しました: %w", err)
		}
		if endDate.Before(startDate) {
			return nil, errors.New("end-date は start-date 以降の日付を指定してください")
		}
		measureTypes, err := parseMeasureTypes(measureTypesStr)
		if err != nil {
			return nil, err
		}
		if cfg.OutputFormat == "" {
			cfg.OutputFormat = DefaultOutputFormat
		}
		if cfg.OutputFormat != DefaultOutputFormat {
			return nil, fmt.Errorf("未対応の output フォーマットです: %s", cfg.OutputFormat)
		}
		cfg.UserID = userID
		cfg.StartDate = startDate
		cfg.EndDate = endDate
		cfg.MeasureTypes = measureTypes

	case OperationAuthURL:
		if cfg.ClientID == "" {
			return nil, errors.New("auth-url では client-id を指定してください")
		}
		if cfg.RedirectURI == "" {
			return nil, errors.New("auth-url では redirect-uri を指定してください")
		}
		if cfg.Scope == "" {
			cfg.Scope = "user.metrics,user.activity"
		}
		if cfg.ResponseType == "" {
			cfg.ResponseType = "code"
		}

	case OperationRequestToken:
		if cfg.ClientID == "" || cfg.ClientSecret == "" {
			return nil, errors.New("request-token では client-id と client-secret を指定してください")
		}
		if cfg.AuthorizationCode == "" {
			return nil, errors.New("request-token では authorization-code を指定してください")
		}
		if cfg.RedirectURI == "" {
			return nil, errors.New("request-token では redirect-uri を指定してください")
		}

	case OperationRefreshToken:
		if cfg.ClientID == "" || cfg.ClientSecret == "" {
			return nil, errors.New("refresh-token では client-id と client-secret を指定してください")
		}
		if cfg.RefreshToken == "" {
			return nil, errors.New("refresh-token では refresh-token を指定してください")
		}

	default:
		return nil, fmt.Errorf("未対応の operation が指定されました: %s", cfg.Operation)
	}

	return cfg, nil
}

// PrintUsage は CLI の使用方法を標準エラー出力に表示します。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `Withings CLI

使用例:
  # 認可URLを生成
  withings -operation auth-url -client-id YOUR_CLIENT_ID -redirect-uri https://yourapp/callback \
    -scope user.metrics,user.activity -state xyz

  # 認可コードをアクセストークンに交換
  withings -operation request-token -client-id YOUR_CLIENT_ID -client-secret YOUR_SECRET \
    -authorization-code CODE_FROM_CALLBACK -redirect-uri https://yourapp/callback

  # リフレッシュトークンでアクセストークンを更新
  withings -operation refresh-token -client-id YOUR_CLIENT_ID -client-secret YOUR_SECRET \
    -refresh-token STORED_REFRESH_TOKEN

  # 日次サマリを取得
  withings -operation daily-summary -access-token ACCESS_TOKEN -user-id 12345 \
    -start-date 2025-09-01 -end-date 2025-09-07 -measure-types weight,diastolic

主なオプション:
  -operation         実行する操作 (auth-url / request-token / refresh-token / daily-summary)
  -client-id         Withings OAuth クライアント ID
  -client-secret     Withings OAuth クライアントシークレット
  -redirect-uri      認可コードを受け取るリダイレクト URI
  -scope             認可リクエストで要求するスコープ (auth-url 用)
  -state             CSRF 対策の state 値 (auth-url 用)
  -authorization-code  認可コールバックで取得した code (request-token 用)
  -refresh-token     保存済みのリフレッシュトークン (refresh-token 用)
  -access-token      API 呼び出しで使用するアクセストークン (daily-summary 用)
  -user-id           Withings ユーザー ID (daily-summary 用)
  -start-date / -end-date   取得期間 (YYYY-MM-DD)
  -measure-types     取得対象 measure type (カンマ区切り)
  -include-activity  活動サマリを含めるか (daily-summary 用, デフォルト true)
  -timeout           API 呼び出しのタイムアウト (例: 20s)

`)
}

func parseDate(value string) (time.Time, error) {
	layout := "2006-01-02"
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}, errors.New("空の日付文字列です")
	}
	t, err := time.Parse(layout, trimmed)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func normalizeCommaSeparated(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	tokens := strings.Split(trimmed, ",")
	var normalized []string
	for _, token := range tokens {
		t := strings.TrimSpace(token)
		if t != "" {
			normalized = append(normalized, t)
		}
	}
	return strings.Join(normalized, ",")
}

var measureAliases = map[string]int{
	"weight":                  1,
	"weightkg":                1,
	"height":                  4,
	"fat_free_mass":           5,
	"fatfreemass":             5,
	"fat_ratio":               6,
	"fatpercentage":           6,
	"fat_mass":                8,
	"fatmass":                 8,
	"diastolic":               9,
	"diastolic_bp":            9,
	"systolic":                10,
	"systolic_bp":             10,
	"heart_pulse":             11,
	"heart_rate":              11,
	"temperature":             12,
	"spo2":                    54,
	"body_temperature":        71,
	"skin_temperature":        73,
	"muscle_mass":             76,
	"hydration":               77,
	"bone_mass":               88,
	"pulse_wave_velocity":     91,
	"vo2max":                  123,
	"qrs_interval":            135,
	"pr_interval":             136,
	"qt_interval":             137,
	"corrected_qt_interval":   138,
	"atrial_fibrillation":     139,
	"atrial_fibrillation_ppg": 139,
}

func parseMeasureTypes(value string) ([]int, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}

	tokens := strings.Split(trimmed, ",")
	if len(tokens) == 0 {
		return nil, nil
	}

	seen := make(map[int]struct{})
	var result []int

	for _, token := range tokens {
		t := normalizeToken(token)
		if t == "" {
			continue
		}
		if num, err := strconv.Atoi(t); err == nil {
			if num <= 0 {
				return nil, fmt.Errorf("measures type は正の整数で指定してください: %d", num)
			}
			if _, exists := seen[num]; !exists {
				seen[num] = struct{}{}
				result = append(result, num)
			}
			continue
		}
		if mapped, ok := measureAliases[t]; ok {
			if _, exists := seen[mapped]; !exists {
				seen[mapped] = struct{}{}
				result = append(result, mapped)
			}
			continue
		}
		return nil, fmt.Errorf("未知の measure type 指定です: %s", token)
	}

	sort.Ints(result)
	return result, nil
}

func normalizeToken(token string) string {
	lowered := strings.ToLower(strings.TrimSpace(token))
	lowered = strings.ReplaceAll(lowered, "-", "_")
	lowered = strings.ReplaceAll(lowered, " ", "_")
	return lowered
}
