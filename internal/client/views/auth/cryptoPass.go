package auth

import (
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/form"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	tea "github.com/charmbracelet/bubbletea"
)

type CryptoPassModel struct {
	form    *form.Form
	help    *help.Auth
	actions *actions.Actions
}

const cryptoPassInputName = "crypto-password"

func NewCryptoPass(acts *actions.Actions) tea.Model {
	inputs := []components.TextInput{
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Crypto password",
			IsPassword:  true,
			Name:        cryptoPassInputName,
			IsFocus:     true,
		}),
	}

	buttons := []components.Button{
		{
			Label: "Login",
			Name:  "login",
		},
	}

	h := help.NewAuth()

	return &CryptoPassModel{
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

func (m *CryptoPassModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *CryptoPassModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typeMsg := msg.(type) {
	case msgs.FormSubmitMsg:
		return m, m.actions.CheckCryptoPassword(typeMsg, cryptoPassInputName)
	default:
		var newMsg tea.Cmd
		m.form, newMsg = m.form.Update(msg)

		return m, newMsg
	}
}

func (m *CryptoPassModel) View() string {
	var b strings.Builder

	b.WriteString(m.form.View())
	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
