package views

import (
	"context"
	"fmt"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/form"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/nav"
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type RegisterModel struct {
	help    *help.Register
	actions *actions.Actions
	nav     *nav.Navigation
	errors  *form.Errors
	log     *zap.Logger
}

func NewRegister(acts *actions.Actions, log *zap.Logger) tea.Model {
	inputs := []components.TextInput{
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Login",
			IsFocus:     true,
			Name:        "login",
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Auth password",
			IsPassword:  true,
			Name:        "password",
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Enctyption password",
			IsPassword:  true,
			Name:        "crypto-password",
		}),
	}

	buttons := []components.Button{
		{
			Label: "Register",
			Name:  "register",
		},
	}

	return &RegisterModel{
		help:    help.NewRegister(),
		actions: acts,
		log:     log,
		nav: &nav.Navigation{
			Inputs:  inputs,
			Buttons: buttons,
		},
		errors: form.NewErrors([]string{
			"login", "password", "crypto-password",
		}),
	}
}

func (m *RegisterModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *RegisterModel) handleRegister() tea.Cmd {
	return func() tea.Msg {
		m.errors.Clear()

		values := models.ClientAuthValues{
			Login:          m.nav.Inputs[0].Value(),
			Password:       m.nav.Inputs[1].Value(),
			CryptoPassword: m.nav.Inputs[2].Value(),
		}

		res, err := m.actions.Register(context.Background(), values)
		if err != nil {
			return msgs.RegisterErrorMsg{Err: err}
		}

		return msgs.RegisterSuccessMsg{Res: res}
	}
}

func (m *RegisterModel) updateInputs(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.nav.Inputs))
	for i := range m.nav.Inputs {
		cmds[i] = m.nav.Inputs[i].Update(msg)
	}

	return cmds
}

func (m *RegisterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		inputCmds := m.updateInputs(msg)
		navCmds := m.nav.HandleKeyPress(msg)

		if msg.Type == tea.KeyEnter && m.nav.IsOnButton() {
			btn := m.nav.GetCurrentButton()
			if btn != nil && btn.Name == "register" {
				return m, m.handleRegister()
			}
		}

		return m, tea.Batch(append(inputCmds, navCmds...)...)
	case msgs.RegisterErrorMsg:
		m.errors.HandleAPIError(msg.Err)
	}

	return m, nil
}

func (m *RegisterModel) View() string {
	var b strings.Builder

	for _, input := range m.nav.Inputs {
		b.WriteString(input.View())

		if err := m.errors.GetFieldError(input.Name); err != "" {
			b.WriteString(styles.ErrorStyle.Render(err))
		}

		b.WriteRune('\n')
	}

	fmt.Fprintf(&b, "\n%s\n", styles.ErrorStyle.Render(m.errors.Message))

	for _, btn := range m.nav.Buttons {
		b.WriteString(btn.View())
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
