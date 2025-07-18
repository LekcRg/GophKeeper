package views

import (
	"fmt"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	error      string
	help       components.AuthHelp
	inputs     []components.TextInput
	buttons    []components.Button
	focusIndex int
}

func NewAuth() tea.Model {
	m := model{
		inputs: []components.TextInput{
			components.NewTextInput(components.TextInputOpts{
				Placeholder: "Login",
				IsFocus:     true,
			}),
			components.NewTextInput(components.TextInputOpts{
				Placeholder: "Password",
				IsPassword:  true,
			}),
		},
		buttons: []components.Button{
			{
				Label: "Login",
			},
			{
				Label: "Register",
			},
		},
		help: components.NewAuthHelp(),
	}

	return m
}

func (m model) lastIndex() int {
	return len(m.inputs) + len(m.buttons) - 1
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) ChangeFocus(s string) (model, []tea.Cmd) {
	m.error = ""

	if s == "up" || s == "shift+tab" {
		m.focusIndex--
	} else {
		m.focusIndex++
	}

	if m.focusIndex > m.lastIndex() {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = m.lastIndex()
	}

	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}

	for i := range m.buttons {
		btnIndex := len(m.inputs) + i

		if btnIndex == m.focusIndex {
			m.buttons[i].Focus()
		} else {
			m.buttons[i].Blur()
		}
	}

	return m, cmds
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down", "enter":
			s := msg.String()
			if s == "enter" && m.focusIndex >= len(m.inputs) {
				// get values from m.Inputs
				// send values from inputs to auth
				m.error = "Login or password isn't valid"

				return m, tea.Batch()
			}

			newM, cmds := m.ChangeFocus(s)

			return newM, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	for _, input := range m.inputs {
		b.WriteString(input.View())
		b.WriteRune('\n')
	}

	fmt.Fprintf(&b, "%s\n\n", styles.ErrorStyle.Render(m.error))

	for _, button := range m.buttons {
		b.WriteString(button.View())
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
