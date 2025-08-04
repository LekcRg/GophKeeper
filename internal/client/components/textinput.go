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
	Valid       []validation.Rule
	Type        InputType
	CharLimit   int
	Width       int
	IsFocus     bool
	IsPassword  bool
}

type TextInput struct {
	Name  string
	Valid []validation.Rule
	textinput.Model
	Type InputType
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
	m := TextInput{
		Model: textinput.New(),
		Name:  opts.Name,
		Valid: opts.Valid,
		Type:  opts.Type,
	}
	m.Cursor.Style = styles.CursorStyle

	if opts.Value != "" {
		m.SetValue(opts.Value)
	}

	if opts.CharLimit > 0 {
		m.CharLimit = opts.CharLimit
	} else {
		m.CharLimit = textCharLimit
	}

	m.Placeholder = opts.Placeholder

	m.Width = textWidth
	if opts.Width > 0 {
		m.Width = opts.Width
	}

	if opts.IsFocus {
		m.Focus()
	}

	if opts.IsPassword {
		m.EchoMode = textinput.EchoPassword
		m.EchoCharacter = textPasswordEchoChar
	}

	return m
}

func digitsOnly(s string) []rune {
	re := regexp.MustCompile(`\D`)
	return []rune(re.ReplaceAllString(s, ""))
}

func (m *TextInput) setInputValue(value string) {
	m.SetValue(value)
	m.SetCursor(len(value))
}

func (m *TextInput) formatCardNumber() {
	val := m.Value()
	if val == "" {
		return
	}

	nums := digitsOnly(val)
	partLen := 4

	if len(nums) <= partLen {
		m.setInputValue(string(nums))

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
		m.setInputValue(formatted)
	}
}

func (m *TextInput) formatCardExp() {
	val := m.Value()
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
		m.setInputValue(formatted)
	}
}

func (m *TextInput) formatNums() {
	val := m.Value()
	if val == "" {
		return
	}

	formatted := string(digitsOnly(val))

	if formatted != val {
		m.setInputValue(formatted)
	}
}

func (m *TextInput) format() {
	switch m.Type {
	case InputTypeNums:
		m.formatNums()
	case InputTypeCardNumber:
		m.formatCardNumber()
	case InputTypeCardExpire:
		m.formatCardExp()
	}
}

func (m *TextInput) Update(msg tea.Msg) tea.Cmd {
	model, cmd := m.Model.Update(msg)
	m.Model = model

	if m.Type != InputTypeText {
		m.format()
	}

	return cmd
}

func (m *TextInput) View() string {
	if m.Focused() {
		m.PromptStyle = styles.FocusedStyle
		m.TextStyle = styles.FocusedStyle
	} else {
		m.PromptStyle = styles.NoStyle
		m.TextStyle = styles.NoStyle
	}

	return m.Model.View()
}
