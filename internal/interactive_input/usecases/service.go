package usecases

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	domain "github.com/landmaster135/devbox/internal/interactive_input/domain"
)

// Config bundles the runtime parameters required to handle an interactive prompt.
type Config struct {
	Prompt          string
	InputType       domain.InputType
	Key             string
	DefaultValue    string
	DefaultProvided bool
	ChoiceOptions   []domain.ChoiceOption
	MaxAttempts     int
}

// Service orchestrates the interactive input workflow.
type Service struct {
	reader    *bufio.Reader
	errWriter io.Writer
}

// ErrExceededAttempts is returned when the user fails validation repeatedly.
var ErrExceededAttempts = errors.New("validation attempts exceeded")

// ErrUserCancelled is returned when the input stream ends unexpectedly.
var ErrUserCancelled = errors.New("input cancelled")

var displayReplacer = strings.NewReplacer("\\n", "\n")

// NewService creates a new Service instance.
func NewService(stdin io.Reader, stderr io.Writer) *Service {
	return &Service{
		reader:    bufio.NewReader(stdin),
		errWriter: stderr,
	}
}

// Run executes the configured interactive prompt and returns the formatted output string.
func (s *Service) Run(cfg Config) (string, error) {
	s.printPreface(cfg)

	switch cfg.InputType {
	case domain.InputTypeText:
		return s.handleTextInput(cfg)
	case domain.InputTypeChoice:
		return s.handleChoiceInput(cfg, true)
	case domain.InputTypeChoiceFlag:
		return s.handleChoiceInput(cfg, false)
	case domain.InputTypeConfirm:
		return s.handleConfirmInput(cfg)
	default:
		return "", fmt.Errorf("unsupported input type: %s", cfg.InputType)
	}
}

func (s *Service) handleTextInput(cfg Config) (string, error) {
	failed := 0
	prompt := s.render(cfg.Prompt)

	for {
		s.printPrompt(prompt)

		value, err := s.readLine()
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(value) == "" {
			if cfg.DefaultProvided {
				fmt.Fprintf(s.errWriter, "空入力のためデフォルト値 %q を使用します。\n", cfg.DefaultValue)
				return formatKeyValue(cfg.Key, cfg.DefaultValue), nil
			}

			failed++
			if !s.hasAttemptsRemaining(failed, cfg.MaxAttempts) {
				fmt.Fprintln(s.errWriter, "入力が空のまま最大試行回数に達したため終了します。")
				return "", ErrExceededAttempts
			}

			fmt.Fprintln(s.errWriter, "値を入力してください。")
			continue
		}

		return formatKeyValue(cfg.Key, value), nil
	}
}

func (s *Service) handleChoiceInput(cfg Config, includeKey bool) (string, error) {
	failed := 0
	prompt := s.render(cfg.Prompt)
	options := make(map[string]string, len(cfg.ChoiceOptions))

	for _, option := range cfg.ChoiceOptions {
		options[option.NormalizedShortcut] = option.Output
	}

	for {
		s.printPrompt(prompt)
		value, err := s.readLine()
		if err != nil {
			return "", err
		}

		normalized := strings.ToLower(strings.TrimSpace(value))
		if output, ok := options[normalized]; ok {
			return formatChoiceValue(cfg.Key, output, includeKey), nil
		}

		failed++
		if !s.hasAttemptsRemaining(failed, cfg.MaxAttempts) {
			fmt.Fprintln(s.errWriter, "有効な選択肢が入力されなかったため終了します。")
			return "", ErrExceededAttempts
		}

		fmt.Fprintf(s.errWriter, "選択肢 %q は無効です。指定されたショートカットを入力してください。\n", value)
	}
}

func (s *Service) handleConfirmInput(cfg Config) (string, error) {
	failed := 0
	prompt := s.render(cfg.Prompt)

	for {
		s.printPrompt(prompt)
		value, err := s.readLine()
		if err != nil {
			return "", err
		}

		normalized := strings.ToLower(strings.TrimSpace(value))
		switch normalized {
		case "y", "yes":
			if cfg.Key == "" {
				return "", nil
			}
			return formatKeyOnly(cfg.Key), nil
		case "n", "no":
			return "", nil
		}

		failed++
		if !s.hasAttemptsRemaining(failed, cfg.MaxAttempts) {
			fmt.Fprintln(s.errWriter, "Y か N が入力されなかったため終了します。")
			return "", ErrExceededAttempts
		}

		fmt.Fprintln(s.errWriter, "y または n を入力してください。")
	}
}

func (s *Service) readLine() (string, error) {
	line, err := s.reader.ReadString('\n')
	if errors.Is(err, io.EOF) {
		if len(line) == 0 {
			return "", ErrUserCancelled
		}
		// Treat partial input before EOF as the final value.
		return strings.TrimRight(line, "\r\n"), nil
	}
	if err != nil {
		return "", err
	}

	return strings.TrimRight(line, "\r\n"), nil
}

func (s *Service) printPreface(cfg Config) {
	if cfg.InputType == domain.InputTypeConfirm {
		fmt.Fprintln(s.errWriter, "[y] Yes / [n] No")
	}
}

func (s *Service) printPrompt(prompt string) {
	fmt.Fprint(s.errWriter, prompt)
}

func (s *Service) hasAttemptsRemaining(failed, max int) bool {
	if max == 0 {
		return true
	}

	return failed < max
}

func (s *Service) render(raw string) string {
	if raw == "" {
		return raw
	}

	return displayReplacer.Replace(raw)
}

func formatKeyValue(key, value string) string {
	return fmt.Sprintf("--%s=%s", key, value)
}

func formatKeyOnly(key string) string {
	return fmt.Sprintf("--%s", key)
}

func formatChoiceValue(key, value string, includeKey bool) string {
	if includeKey {
		return formatKeyValue(key, value)
	}

	return value
}
