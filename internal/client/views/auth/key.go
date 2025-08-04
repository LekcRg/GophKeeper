package auth

import (
	"context"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/form"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyModel struct {
	form    *form.Form
	help    *help.Auth
	actions *actions.Actions
}

func NewKey(acts *actions.Actions) tea.Model {
	inputWidth := 50

	inputs := []components.TextInput{
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Key",
			Name:        "key",
			IsFocus:     true,
			Width:       inputWidth,
		}),
	}

	buttons := []components.Button{
		{
			Label: "Login",
			Name:  "login",
		},
	}

	h := help.NewAuth()

	return &KeyModel{
		form: form.NewForm(form.Opts{
			Inputs:  inputs,
			Buttons: buttons,
			Up:      h.Keys.Up,
			Down:    h.Keys.Down,
		}),
		help:    h,
		actions: acts,
	}
}

func (m *KeyModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *KeyModel) handleSubmit(msg msgs.FormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		res, err := m.actions.GetCredentials(context.Background(), msg.Values["key"])
		if err != nil {
			return msgs.ErrorMsg(err)
		}

		return res
	}
}

func (m *KeyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typeMsg := msg.(type) {
	case msgs.FormSubmitMsg:
		return m, m.handleSubmit(typeMsg)
	default:
		var newMsg tea.Cmd
		m.form, newMsg = m.form.Update(msg)

		return m, newMsg
	}
}

func (m *KeyModel) View() string {
	var b strings.Builder

	b.WriteString(m.form.View())
	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
