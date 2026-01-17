package usecases

import (
	"bytes"
	"errors"
	"testing"

	domain "github.com/landmaster135/devbox/internal/interactive_input/domain"
)

func TestService_TextInputValue(t *testing.T) {
	stdin := bytes.NewBufferString("hello-world\n")
	var stderr bytes.Buffer

	svc := NewService(stdin, &stderr)
	output, err := svc.Run(Config{
		Prompt:      "Input value: ",
		InputType:   domain.InputTypeText,
		Key:         "text",
		MaxAttempts: 3,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "--text=hello-world" {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestService_TextInputUsesDefault(t *testing.T) {
	stdin := bytes.NewBufferString("\n")
	var stderr bytes.Buffer

	svc := NewService(stdin, &stderr)
	output, err := svc.Run(Config{
		Prompt:          "Input path: ",
		InputType:       domain.InputTypeText,
		Key:             "path",
		MaxAttempts:     2,
		DefaultValue:    ".",
		DefaultProvided: true,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "--path=." {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestService_ChoiceInput(t *testing.T) {
	stdin := bytes.NewBufferString("z\nw\n")
	var stderr bytes.Buffer

	svc := NewService(stdin, &stderr)
	output, err := svc.Run(Config{
		Prompt:    "Select: ",
		InputType: domain.InputTypeChoice,
		Key:       "operation",
		ChoiceOptions: []domain.ChoiceOption{
			{Shortcut: "v", NormalizedShortcut: "v", Output: "vlc"},
			{Shortcut: "w", NormalizedShortcut: "w", Output: "win"},
		},
		MaxAttempts: 3,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "--operation=win" {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestService_ChoiceFlagInput(t *testing.T) {
	stdin := bytes.NewBufferString("v\n")
	var stderr bytes.Buffer

	svc := NewService(stdin, &stderr)
	output, err := svc.Run(Config{
		Prompt:    "Select flag: ",
		InputType: domain.InputTypeChoiceFlag,
		Key:       "ignored",
		ChoiceOptions: []domain.ChoiceOption{
			{Shortcut: "v", NormalizedShortcut: "v", Output: "-operation=vlc"},
		},
		MaxAttempts: 2,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "-operation=vlc" {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestService_ConfirmNegative(t *testing.T) {
	stdin := bytes.NewBufferString("n\n")
	var stderr bytes.Buffer

	svc := NewService(stdin, &stderr)
	output, err := svc.Run(Config{
		Prompt:      "Move files? ",
		InputType:   domain.InputTypeConfirm,
		Key:         "move",
		MaxAttempts: 2,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != "" {
		t.Fatalf("expected empty output, got %q", output)
	}
}

func TestService_ConfirmExceededAttempts(t *testing.T) {
	stdin := bytes.NewBufferString("maybe\n")
	var stderr bytes.Buffer

	svc := NewService(stdin, &stderr)
	_, err := svc.Run(Config{
		Prompt:      "Move files? ",
		InputType:   domain.InputTypeConfirm,
		Key:         "move",
		MaxAttempts: 1,
	})

	if !errors.Is(err, ErrExceededAttempts) {
		t.Fatalf("expected ErrExceededAttempts, got %v", err)
	}
}

func TestService_UserCancelled(t *testing.T) {
	stdin := bytes.NewBuffer(nil)
	var stderr bytes.Buffer

	svc := NewService(stdin, &stderr)
	_, err := svc.Run(Config{
		Prompt:      "Input value? ",
		InputType:   domain.InputTypeText,
		Key:         "value",
		MaxAttempts: 1,
	})

	if !errors.Is(err, ErrUserCancelled) {
		t.Fatalf("expected ErrUserCancelled, got %v", err)
	}
}
