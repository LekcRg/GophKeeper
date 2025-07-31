package create

import (
	"context"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/form"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
)

type PasswordModel struct {
	help    *help.Auth
	actions *actions.Actions
	log     *zap.Logger
	form    *form.Form
}

var (
	passwordNameInput     = "name"
	passwordLoginInput    = "login"
	passwordPasswordInput = "password"
	passwordURLInput      = "url"
	passwordCreateBtn     = "create"
	passwordGenerateBtn   = "gen-password"
)

func NewPassword(acts *actions.Actions, log *zap.Logger) tea.Model {
	inputs := []components.TextInput{
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Name",
			IsFocus:     true,
			Name:        passwordNameInput,
			Valid:       []validation.Rule{validation.Required},
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Login",
			Name:        passwordLoginInput,
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Password",
			Name:        passwordPasswordInput,
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "URL",
			Name:        passwordURLInput,
		}),
	}

	buttons := []components.Button{
		{
			Label: "Create",
			Name:  passwordCreateBtn,
		},
		// {
		// 	Label: "Generate password",
		// 	Name:  passwordGenerateBtn,
		// },
	}

	h := help.NewAuth()

	return &PasswordModel{
		actions: acts,
		form: form.NewForm(form.FormOpts{
			Inputs:  inputs,
			Buttons: buttons,
			Up:      h.Keys.Up,
			Down:    h.Keys.Down,
		}),
		help: h,
		log:  log,
	}
}

func (m *PasswordModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *PasswordModel) handleSubmit(msg msgs.FormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		name := msg.Values[passwordNameInput]
		data := models.VaultItemDataPassword{
			Login:    msg.Values[passwordLoginInput],
			Password: msg.Values[passwordPasswordInput],
			URL:      msg.Values[passwordURLInput],
		}

		res, err := m.actions.CreateVaultItem(context.Background(), name, "password", data)
		if err != nil {
			return msgs.ErrorMsg(err)
		}

		return msgs.CreateVaultSuccess{Item: res}
	}
}

func (m *PasswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	default:
		switch typeMsg := msg.(type) {
		case msgs.FormSubmitMsg:
			return m, m.handleSubmit(typeMsg)
		default:
			var newMsg tea.Cmd
			m.form, newMsg = m.form.Update(msg)

			return m, newMsg
		}
	}
}

func (m *PasswordModel) View() string {
	var b strings.Builder

	b.WriteString(m.form.View())
	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
