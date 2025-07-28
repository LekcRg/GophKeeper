package views

import (
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/nav"
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	"github.com/LekcRg/GophKeeper/internal/server/service/valid"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectAuthModel struct {
	help  *help.SelectAuth
	error string
	nav   nav.Navigation
}

func NewSelectAuth(addr string) tea.Model {
	m := &SelectAuthModel{
		nav: nav.Navigation{
			Inputs: []components.TextInput{
				components.NewTextInput(components.TextInputOpts{
					Placeholder: "Server address",
					Name:        "address",
					IsFocus:     true,
					Value:       addr,
				}),
			},
			Buttons: []components.Button{
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
		},
		help: help.NewSelectAuth(),
	}

	return m
}

func (m *SelectAuthModel) Init() tea.Cmd {
	return nil
}

func (m *SelectAuthModel) Submit() tea.Cmd {
	addr := m.nav.Inputs[0].Value()

	err := valid.ValidAddr(addr)
	if err != nil {
		m.error = err.Error()

		return nil
	}

	return func() tea.Msg {
		return msgs.SelectAuthMsg{
			Selected: m.nav.GetCurrentButton().Name,
			Address:  addr,
		}
	}
}

func (m *SelectAuthModel) updateInputs(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.nav.Inputs))
	for i := range m.nav.Inputs {
		cmds[i] = m.nav.Inputs[i].Update(msg)
	}

	return cmds
}

func (m *SelectAuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		inputCmds := m.updateInputs(msg)
		navCmds := m.nav.HandleKeyPress(msg)

		if msg.Type == tea.KeyEnter && m.nav.IsOnButton() {
			btn := m.nav.GetCurrentButton()
			if btn != nil {
				return m, m.Submit()
			}
		}

		return m, tea.Batch(append(inputCmds, navCmds...)...)
	default:
		return m, nil
	}
}

func (m *SelectAuthModel) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(m.nav.Inputs[0].View())
	b.WriteString(styles.ErrorStyle.Render(m.error))
	b.WriteString("\n\n")

	for _, btn := range m.nav.Buttons {
		b.WriteString(btn.View())
		b.WriteRune('\n')
	}

	b.WriteString("\n")

	b.WriteString(m.help.View())

	return b.String()
}
