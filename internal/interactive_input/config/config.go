package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	domain "github.com/landmaster135/devbox/internal/interactive_input/domain"
)

// Config represents parsed CLI flags for the interactive-input tool.
type Config struct {
	Prompt          string
	InputType       domain.InputType
	Key             string
	DefaultValue    string
	DefaultProvided bool
	ChoiceOptions   []domain.ChoiceOption
	MaxAttempts     int
	Help            bool
}

// ParseFlags parses command-line arguments into Config.
func ParseFlags(args []string) (*Config, error) {
	flagSet := flag.NewFlagSet("interactive-input", flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)

	prompt := flagSet.String("prompt", "", "ユーザーに表示する質問文（\\nで改行可）")
	inputType := flagSet.String("input-type", "", "入力タイプ text|choice|confirm")
	key := flagSet.String("key", "", "標準出力に使うキー名（--<key>=value）")
	maxAttempts := flagSet.Int("max-attempts", 3, "バリデーション失敗時の再入力許可回数 (0で無制限)")
	help := flagSet.Bool("help", false, "使い方を表示")
	flagSet.BoolVar(help, "h", false, "使い方を表示")

	var defaultValue stringFlag
	flagSet.Var(&defaultValue, "default", "text入力でEnterのみだった場合に採用するデフォルト値")

	var choiceOptions multiValue
	flagSet.Var(&choiceOptions, "choice-option", "choice入力の選択肢を 'shortcut|output' 形式で指定（複数可）")

	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	if *help {
		return &Config{Help: true}, nil
	}

	if strings.TrimSpace(*prompt) == "" {
		return nil, errors.New("--prompt は必須です")
	}

	inputTypeValue := strings.ToLower(strings.TrimSpace(*inputType))
	var parsedType domain.InputType

	switch inputTypeValue {
	case string(domain.InputTypeText):
		parsedType = domain.InputTypeText
	case string(domain.InputTypeChoice):
		parsedType = domain.InputTypeChoice
	case string(domain.InputTypeChoiceFlag):
		parsedType = domain.InputTypeChoiceFlag
	case string(domain.InputTypeConfirm):
		parsedType = domain.InputTypeConfirm
	default:
		return nil, fmt.Errorf("--input-type には text / choice / choice-flag / confirm のいずれかを指定してください: %s", *inputType)
	}

	requiresKey := parsedType != domain.InputTypeChoiceFlag
	trimmedKey := strings.TrimSpace(*key)

	if requiresKey {
		if trimmedKey == "" {
			return nil, errors.New("--key はこの入力タイプで必須です")
		}
		if err := validateKey(trimmedKey); err != nil {
			return nil, err
		}
	} else if trimmedKey != "" {
		if err := validateKey(trimmedKey); err != nil {
			return nil, err
		}
	}

	if *maxAttempts < 0 {
		return nil, errors.New("--max-attempts は0以上で指定してください")
	}

	parsedChoices, err := parseChoiceOptions(choiceOptions)
	if err != nil {
		return nil, err
	}

	if parsedType == domain.InputTypeChoice && len(parsedChoices) == 0 {
		return nil, errors.New("--input-type choice のときは --choice-option を少なくとも1件指定してください")
	}

	cfg := &Config{
		Prompt:          *prompt,
		InputType:       parsedType,
		Key:             trimmedKey,
		DefaultValue:    defaultValue.value,
		DefaultProvided: defaultValue.set,
		ChoiceOptions:   parsedChoices,
		MaxAttempts:     *maxAttempts,
		Help:            *help,
	}

	return cfg, nil
}

// PrintUsage prints help text for the CLI.
func PrintUsage() {
	usage := `interactive-input はプロンプトごとに1回起動し、ユーザー入力をキー付き文字列または指定のフラグとして標準出力へ返すCLIです。

使用例:
  # 任意テキスト入力（Enterのみならデフォルト"."）
  interactive-input \
    --prompt "Input path (blank = current dir): " \
    --input-type text \
    --key path \
    --default "."

  # 選択肢
  interactive-input \
    --prompt "Select operation: " \
    --input-type choice \
    --key operation \
    --choice-option "v|vlc" \
    --choice-option "w|win"

  # フラグ選択肢（--choice-option の output をそのまま出力）
  interactive-input \
    --prompt "Select rename flag: " \
    --input-type choice-flag \
    --choice-option "v|-operation=vlc" \
    --choice-option "w|-operation=win"

  # 確認
  interactive-input \
    --prompt "Move originals? " \
    --input-type confirm \
    --key move

主要フラグ:
  --prompt string           質問文（必須）
  --input-type string       text / choice / choice-flag / confirm から選択（必須）
  --key string              出力のキー名（choice-flag以外で必須。空白・=は不可）
  --default string          text入力で空行時に採用する値
  --choice-option string    choice/choice-flag専用。"shortcut|output" 形式。複数指定可
  --max-attempts int        バリデーション失敗時の再入力回数。0で無制限（既定:3）
  --help (-h)               このヘルプを表示
`

	fmt.Fprint(os.Stdout, usage)
}

func validateKey(key string) error {
	if strings.ContainsAny(key, " \t\r\n") {
		return fmt.Errorf("--key に空白は使用できません: %s", key)
	}

	if strings.ContainsRune(key, '=') {
		return fmt.Errorf("--key に '=' は使用できません: %s", key)
	}

	if strings.HasPrefix(key, "-") {
		return fmt.Errorf("--key から先頭のハイフンは取り除いて指定してください: %s", key)
	}

	validKey := regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-_]*$`)
	if !validKey.MatchString(key) {
		return fmt.Errorf("--key には英数字と - _ のみ使用できます: %s", key)
	}

	return nil
}

func parseChoiceOptions(values []string) ([]domain.ChoiceOption, error) {
	if len(values) == 0 {
		return nil, nil
	}

	options := make([]domain.ChoiceOption, 0, len(values))
	seen := make(map[string]struct{})

	for _, raw := range values {
		parts := strings.SplitN(raw, "|", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("--choice-option は 'shortcut|output' 形式で指定してください: %s", raw)
		}

		shortcut := strings.TrimSpace(parts[0])
		output := strings.TrimSpace(parts[1])

		if shortcut == "" {
			return nil, fmt.Errorf("--choice-option の shortcut が空です: %s", raw)
		}

		if len([]rune(shortcut)) != 1 {
			return nil, fmt.Errorf("--choice-option の shortcut は1文字で指定してください: %s", shortcut)
		}

		if output == "" {
			return nil, fmt.Errorf("--choice-option の output が空です: %s", raw)
		}

		normalized := strings.ToLower(shortcut)
		if _, exists := seen[normalized]; exists {
			return nil, fmt.Errorf("--choice-option の shortcut が重複しています: %s", shortcut)
		}

		seen[normalized] = struct{}{}

		options = append(options, domain.ChoiceOption{
			Shortcut:           shortcut,
			NormalizedShortcut: normalized,
			Output:             output,
		})
	}

	return options, nil
}

type multiValue []string

func (m *multiValue) String() string {
	return strings.Join(*m, ",")
}

func (m *multiValue) Set(value string) error {
	*m = append(*m, value)
	return nil
}

type stringFlag struct {
	value string
	set   bool
}

func (s *stringFlag) String() string {
	return s.value
}

func (s *stringFlag) Set(value string) error {
	s.value = value
	s.set = true
	return nil
}
