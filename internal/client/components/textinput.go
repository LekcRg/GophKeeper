package components

import (
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputOpts struct {
	Placeholder string
	Name        string
	Value       string
	CharLimit   int
	IsFocus     bool
	IsPassword  bool
}

type TextInput struct {
	Name string
	textinput.Model
}

const (
	textCharLimit        = 32
	textWidth            = 30
	textPasswordEchoChar = '•'
)

func NewTextInput(opts TextInputOpts) TextInput {
	ti := TextInput{
		Model: textinput.New(),
		Name:  opts.Name,
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

	if opts.IsFocus {
		ti.Focus()
	}

	if opts.IsPassword {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = textPasswordEchoChar
	}

	return ti
}

func (ti *TextInput) Update(msg tea.Msg) tea.Cmd {
	model, cmd := ti.Model.Update(msg)
	ti.Model = model

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
