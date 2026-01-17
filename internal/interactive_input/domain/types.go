package domain

// InputType represents the supported interactive input modes.
type InputType string

const (
	// InputTypeText accepts free-form text responses.
	InputTypeText InputType = "text"
	// InputTypeChoice accepts one-character shortcut selections.
	InputTypeChoice InputType = "choice"
	// InputTypeConfirm accepts Y/N confirmation responses.
	InputTypeConfirm InputType = "confirm"
)

// ChoiceOption describes a selectable shortcut and the value it maps to.
type ChoiceOption struct {
	Shortcut           string
	NormalizedShortcut string
	Output             string
}
