package components

import (
	"regexp"

	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type TextInputOpts struct {
	Placeholder string
	Name        string
	Value       string
	Type        InputType
	Valid       []validation.Rule
	CharLimit   int
	Width       int
	IsFocus     bool
	IsPassword  bool
}

type TextInput struct {
	Name  string
	Type  InputType
	Valid []validation.Rule
	textinput.Model
}

type InputType int

const (
	InputTypeText = iota
	InputTypeNums
	InputTypeCardNumber
	InputTypeCardExpire
)

const (
	textCharLimit        = 60
	textWidth            = 30
	textPasswordEchoChar = '•'
)

func NewTextInput(opts TextInputOpts) TextInput {
	ti := TextInput{
		Model: textinput.New(),
		Name:  opts.Name,
		Valid: opts.Valid,
		Type:  opts.Type,
	}
	ti.Cursor.Style = styles.CursorStyle

	if opts.Value != "" {
		ti.SetValue(opts.Value)
	}

	if opts.CharLimit > 0 {
		ti.CharLimit = opts.CharLimit
	} else {
		ti.CharLimit = textCharLimit
	}

	ti.Placeholder = opts.Placeholder

	ti.Width = textWidth
	if opts.Width > 0 {
		ti.Width = opts.Width
	}

	if opts.IsFocus {
		ti.Focus()
	}

	if opts.IsPassword {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = textPasswordEchoChar
	}

	return ti
}

func digitsOnly(s string) []rune {
	re := regexp.MustCompile(`\D`)
	return []rune(re.ReplaceAllString(s, ""))
}

func (ti *TextInput) setInputValue(value string) {
	ti.SetValue(value)
	ti.SetCursor(len([]rune(value)))
}

func (ti *TextInput) formatCardNumber() {
	val := ti.Value()
	if val == "" {
		return
	}

	nums := digitsOnly(val)
	partLen := 4

	if len(nums) <= partLen {
		ti.setInputValue(string(nums))

		return
	}

	var res []rune

	for i := 0; i < len(nums); i += partLen {
		end := i + partLen

		if end > len(nums) {
			end = len(nums)
		}

		res = append(res, nums[i:end]...)
		if end < len(nums) {
			res = append(res, ' ')
		}
	}

	formatted := string(res)

	if val != formatted {
		ti.setInputValue(formatted)
	}
}

func (ti *TextInput) formatCardExp() {
	val := ti.Value()
	if val == "" {
		return
	}

	nums := digitsOnly(val)
	partLen := 2

	var formatted string

	if len([]rune(val)) <= partLen {
		formatted = string(nums)
	} else {
		formatted = string(nums[:partLen]) + "/" + string(nums[partLen:])
	}

	if formatted != val {
		ti.setInputValue(formatted)
	}
}

func (ti *TextInput) formatNums() {
	val := ti.Value()
	if val == "" {
		return
	}

	formatted := string(digitsOnly(val))

	if formatted != val {
		ti.setInputValue(formatted)
	}
}

func (ti *TextInput) format() {
	switch ti.Type {
	case InputTypeNums:
		ti.formatNums()
	case InputTypeCardNumber:
		ti.formatCardNumber()
	case InputTypeCardExpire:
		ti.formatCardExp()
	}
}

func (ti *TextInput) Update(msg tea.Msg) tea.Cmd {
	model, cmd := ti.Model.Update(msg)
	ti.Model = model

	if ti.Type != InputTypeText {
		ti.format()
	}

	return cmd
}

func (ti *TextInput) View() string {
	if ti.Focused() {
		ti.PromptStyle = styles.FocusedStyle
		ti.TextStyle = styles.FocusedStyle
	} else {
		ti.PromptStyle = styles.NoStyle
		ti.TextStyle = styles.NoStyle
	}

	return ti.Model.View()
}
