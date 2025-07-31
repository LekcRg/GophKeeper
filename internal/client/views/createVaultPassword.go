package views

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

type CreateVaultPasswordModel struct {
	help    *help.Auth
	actions *actions.Actions
	log     *zap.Logger
	form    *form.Form
}

var (
	CreateVaultNameInput     = "name"
	CreateVaultLoginInput    = "login"
	CreateVaultPasswordInput = "password"
	CreateVaultURLInput      = "url"
	CreateVaultCreateBtn     = "create"
	CreateVaultGenerateBtn   = "gen-password"
)

func NewCreateVaultPassword(acts *actions.Actions, log *zap.Logger) tea.Model {
	inputs := []components.TextInput{
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Name",
			IsFocus:     true,
			Name:        CreateVaultNameInput,
			Valid:       []validation.Rule{validation.Required},
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Login",
			Name:        CreateVaultLoginInput,
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Password",
			Name:        CreateVaultPasswordInput,
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "URL",
			Name:        CreateVaultURLInput,
		}),
	}

	buttons := []components.Button{
		{
			Label: "Create",
			Name:  CreateVaultCreateBtn,
		},
		// {
		// 	Label: "Generate password",
		// 	Name:  CreateVaultGenerateBtn,
		// },
	}

	h := help.NewAuth()

	return &CreateVaultPasswordModel{
		actions: acts,
		form:    form.NewForm(inputs, buttons, h.Keys.Up, h.Keys.Down),
		help:    h,
		log:     log,
	}
}

func (m *CreateVaultPasswordModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *CreateVaultPasswordModel) handleSubmit(msg msgs.FormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		name := msg.Values[CreateVaultNameInput]
		data := models.VaultItemDataPassword{
			Login:    msg.Values[CreateVaultLoginInput],
			Password: msg.Values[CreateVaultPasswordInput],
			URL:      msg.Values[CreateVaultURLInput],
		}

		res, err := m.actions.CreateVaultItem(context.Background(), name, "password", data)
		if err != nil {
			return msgs.ErrorMsg(err)
		}

		return msgs.CreateVaultSuccess{Item: res}
	}
}

func (m *CreateVaultPasswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *CreateVaultPasswordModel) View() string {
	var b strings.Builder

	b.WriteString(m.form.View())
	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
