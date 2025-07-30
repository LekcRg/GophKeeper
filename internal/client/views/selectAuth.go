package views

import (
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/form"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/router"
	"github.com/LekcRg/GophKeeper/internal/server/service/valid"
	tea "github.com/charmbracelet/bubbletea"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type SelectAuthModel struct {
	form *form.Form
	help *help.SelectAuth
}

const addrInputName = "address"

func NewSelectAuth(addr string) tea.Model {
	inputs := []components.TextInput{
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Server address",
			Name:        addrInputName,
			IsFocus:     true,
			Value:       addr,
			Valid: []validation.Rule{
				validation.Required,
				is.URL,
				validation.By(valid.IsContainsHTTP),
			},
		}),
	}

	buttons := []components.Button{
		{
			Label: "Register",
			Name:  string(router.RegisterView),
		},
		{
			Label: "Login with token",
			Name:  string(router.TokenAuthView),
		},
		{
			Label: "Update and get new token",
			Name:  string(router.UpdateTokenView),
		},
	}

	h := help.NewSelectAuth()

	return &SelectAuthModel{
		form: form.NewForm(inputs, buttons, h.Keys.Up, h.Keys.Down),
		help: h,
	}
}

func (m *SelectAuthModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *SelectAuthModel) handleSubmit(msg msgs.FormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		return msgs.SelectAuthMsg{
			Selected: msg.ButtonName,
			Address:  msg.Values[addrInputName],
		}
	}
}

func (m *SelectAuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typeMsg := msg.(type) {
	case msgs.FormSubmitMsg:
		return m, m.handleSubmit(typeMsg)
	default:
		var newMsg tea.Cmd
		m.form, newMsg = m.form.Update(msg)

		return m, newMsg
	}
}

func (m *SelectAuthModel) View() string {
	var b strings.Builder

	b.WriteString(m.form.View())
	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
