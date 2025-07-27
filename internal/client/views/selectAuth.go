package views

import (
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectAuthModel struct {
	tea.Model
	input      components.TextInput
	buttons    []components.Button
	focusIndex int
	help       *help.SelectAuth
}

func NewSelectAuth() *SelectAuthModel {
	m := &SelectAuthModel{
		buttons: []components.Button{
			{
				Label: "Register",
				Name:  "register",
			},
			{
				Label: "Login with token",
				Name:  "token",
			},
			{
				Label: "Update and get new token",
				Name:  "update-token",
			},
		},
		input: components.NewTextInput(components.TextInputOpts{
			Placeholder: "Server address",
			Name:        "address",
			IsFocus:     true,
		}),
		help: help.NewSelectAuth(),
	}

	return m
}

func (m *SelectAuthModel) Init() tea.Cmd {
	return nil
}

func (m *SelectAuthModel) handleKeyPress(keyMsg tea.KeyMsg) tea.Cmd {
	switch keyMsg.String() {
	case "enter":
		if m.focusIndex != 0 {
			return func() tea.Msg {
				return msgs.SelectAuthMsg(m.buttons[m.focusIndex-1].Name)
			}
		}
	case "up":
		m.focusIndex--
	case "down":
		m.focusIndex++
	}

	lastIndex := len(m.buttons)
	if m.focusIndex > lastIndex {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = lastIndex
	}

	return nil
}

func (m *SelectAuthModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return tea.Batch(m.input.Update(msg), m.handleKeyPress(msg))
	default:
		return nil
	}
}

func (m *SelectAuthModel) View() string {
	var b strings.Builder

	if m.focusIndex == 0 {
		m.input.Focus()
	} else {
		m.input.Blur()
	}

	b.WriteRune('\n')
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	for i := 0; i < len(m.buttons); i++ {
		if i == m.focusIndex-1 {
			m.buttons[i].Focus()
		} else {
			m.buttons[i].Blur()
		}

		b.WriteString(m.buttons[i].View())
		b.WriteRune('\n')
	}

	b.WriteString("\n")

	b.WriteString(m.help.View())

	return b.String()
}
