package components

import (
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type Textarea struct {
	textarea.Model
	Name string
}

type TextareaOpts struct {
	Placeholder string
	Value       string
	Name        string
	Focused     bool
}

func NewTextarea(opts TextareaOpts) Textarea {
	ta := textarea.New()
	ta.Placeholder = opts.Placeholder
	ta.FocusedStyle, ta.BlurredStyle = styles.GetTextareaStyles()
	ta.SetValue(opts.Value)

	if opts.Focused {
		ta.Focus()
	} else {
		ta.Blur()
	}

	return Textarea{
		Model: ta,
		Name:  opts.Name,
	}
}

func (ta *Textarea) Update(msg tea.Msg) tea.Cmd {
	model, cmd := ta.Model.Update(msg)
	ta.Model = model

	return cmd
}
