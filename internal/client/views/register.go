package views

import (
	"context"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/form"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/LekcRg/GophKeeper/internal/server/service/valid"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type RegisterModel struct {
	help    *help.Register
	actions *actions.Actions
	log     *zap.Logger
	form    *form.Form
}

const (
	RegisterLoginInputName          = "login"
	RegisterPasswordInputName       = "password"
	RegisterCryptoPasswordInputName = "crypto-password"
)

func NewRegister(acts *actions.Actions, log *zap.Logger) tea.Model {
	inputs := []components.TextInput{
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Login",
			IsFocus:     true,
			Name:        RegisterLoginInputName,
			Valid:       valid.LoginRules,
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Auth password",
			IsPassword:  true,
			Name:        RegisterPasswordInputName,
			Valid:       valid.PasswordRules,
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Enctyption password",
			IsPassword:  true,
			Name:        RegisterCryptoPasswordInputName,
			Valid:       valid.PasswordRules,
		}),
	}

	buttons := []components.Button{
		{
			Label: "Register",
			Name:  "register",
		},
	}

	return &RegisterModel{
		actions: acts,
		form:    form.NewForm(inputs, buttons),
		help:    help.NewRegister(),
		log:     log,
	}
}

func (m *RegisterModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *RegisterModel) handleSubmit(msg msgs.FormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		if msg.Values[RegisterPasswordInputName] == msg.Values[RegisterCryptoPasswordInputName] {
			return msgs.ErrorMsg(errs.ErrEqualPasswords)
		}

		values := models.ClientRegisterValues{
			Login:          msg.Values[RegisterLoginInputName],
			Password:       msg.Values[RegisterPasswordInputName],
			CryptoPassword: msg.Values[RegisterCryptoPasswordInputName],
		}

		res, err := m.actions.Register(context.Background(), values)
		if err != nil {
			return err
		}

		return res
	}
}

func (m *RegisterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *RegisterModel) View() string {
	var b strings.Builder

	b.WriteString(m.form.View())
	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
