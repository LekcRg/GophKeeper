package views

import (
	"fmt"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AuthModel struct {
	tea.Model
	error      string
	help       *components.AuthHelp
	inputs     []components.TextInput
	buttons    []components.Button
	focusIndex int
}

func NewAuth() *AuthModel {
	m := AuthModel{
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

	return &m
}

func (m *AuthModel) lastIndex() int {
	return len(m.inputs) + len(m.buttons) - 1
}

func (m *AuthModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *AuthModel) ChangeFocus(msg tea.Msg) []tea.Cmd {
	m.error = ""

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyUp || keyMsg.Type == tea.KeyShiftTab {
			m.focusIndex--
		} else {
			m.focusIndex++
		}
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

	return cmds
}

func (m *AuthModel) Update(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)

	key := keyMsg.Type
	if ok && (key == tea.KeyTab || key == tea.KeyShiftTab ||
		key == tea.KeyUp || key == tea.KeyDown || key == tea.KeyEnter) {
		if key == tea.KeyEnter && m.focusIndex >= len(m.inputs) {
			// get values from m.Inputs
			// send values from inputs to auth
			m.error = "Login or password isn't valid"

			return tea.Batch()
		}

		cmds := m.ChangeFocus(msg)

		return tea.Batch(cmds...)
	}

	cmd := m.updateInputs(msg)

	return cmd
}

func (m *AuthModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *AuthModel) View() string {
	var b strings.Builder

	for i := 0; i < len(m.inputs); i++ {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}

	fmt.Fprintf(&b, "%s\n\n", styles.ErrorStyle.Render(m.error))

	for i := 0; i < len(m.buttons); i++ {
		b.WriteString(m.buttons[i].View())
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
