package views

import (
	"context"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/actions"
	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/form"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateKeyModel struct {
	form       *form.Form
	help       *help.Auth
	actions    *actions.Actions
	key        string
	successBtn components.Button
}

const (
	updateKeyLoginInputName    = "login"
	updateKeyPasswordInputName = "password"
)

func NewUpdateKey(acts *actions.Actions) tea.Model {
	inputs := []components.TextInput{
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Login",
			IsFocus:     true,
			Name:        updateKeyLoginInputName,
		}),
		components.NewTextInput(components.TextInputOpts{
			Placeholder: "Auth password",
			IsPassword:  true,
			Name:        updateKeyPasswordInputName,
		}),
	}

	buttons := []components.Button{
		{
			Label: "Update and get key",
			Name:  "update",
		},
	}

	h := help.NewAuth()

	return &UpdateKeyModel{
		form:    form.NewForm(inputs, buttons, h.Keys.Up, h.Keys.Down),
		help:    h,
		actions: acts,
		successBtn: components.Button{
			Label:   "ok",
			Name:    "ok",
			Focused: true,
		},
	}
}

func (m *UpdateKeyModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *UpdateKeyModel) handleSubmit(msg msgs.FormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		req := models.UserLogin{
			Login:    msg.Values[updateKeyLoginInputName],
			Password: msg.Values[updateKeyPasswordInputName],
		}

		res, err := m.actions.UpdateKey(context.Background(), req)
		if err != nil {
			return msgs.ErrorMsg(err)
		}

		m.key = res.Key

		return ""
	}
}

func (m *UpdateKeyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typeMsg := msg.(type) {
	case msgs.FormSubmitMsg:
		return m, m.handleSubmit(typeMsg)
	case tea.KeyMsg:
		if m.key != "" && key.Matches(typeMsg, help.Select) {
			return m, func() tea.Msg {
				return msgs.UpdateKeySuccessMsg(m.key)
			}
		}
	}

	var newMsg tea.Cmd
	m.form, newMsg = m.form.Update(msg)

	return m, newMsg
}

func (m *UpdateKeyModel) View() string {
	var b strings.Builder

	if m.key != "" {
		b.WriteRune('\n')
		b.WriteString(styles.SuccessTitle.Render("✔ Token successfully generated"))
		b.WriteRune('\n')
		b.WriteString(styles.TokenBox.Render(m.key))
		b.WriteRune('\n')
		b.WriteString(styles.Instruction.Render("Copy this token and log in using it."))
		b.WriteString("\n")
		b.WriteString(m.successBtn.View())
		b.WriteString("\n\n")
	} else {
		b.WriteString(m.form.View())
		b.WriteRune('\n')
	}

	b.WriteString(m.help.View())

	return b.String()
}
