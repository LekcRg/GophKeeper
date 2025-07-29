package views

import (
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/form"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type KeyAuthModel struct {
	form    *form.Form
	help    *help.Register
	log     *zap.Logger
	actions *actions.Actions
}

func NewKeyAuth(acts *actions.Actions) tea.Model {
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

	return &KeyAuthModel{
		form:    form.NewForm(inputs, buttons),
		help:    help.NewRegister(),
		actions: acts,
	}
}

func (m *KeyAuthModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *KeyAuthModel) handleSubmit(msg msgs.FormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		// req := models.UserLogin{
		// 	Login:    msg.Values[updateTokenLoginInputName],
		// 	Password: msg.Values[updateTokenPasswordInputName],
		// }

		// res, err := m.actions.UpdateKey(context.Background(), req)
		// if err != nil {
		// 	return msgs.ErrorMsg(err)
		// }

		return nil
	}
}

func (m *KeyAuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typeMsg := msg.(type) {
	case msgs.FormSubmitMsg:
		return m, m.handleSubmit(typeMsg)
	default:
		var newMsg tea.Cmd
		m.form, newMsg = m.form.Update(msg)

		return m, newMsg
	}
}

func (m *KeyAuthModel) View() string {
	var b strings.Builder

	b.WriteString(m.form.View())
	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
