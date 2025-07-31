package create

import (
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/components"
	"github.com/LekcRg/GophKeeper/internal/client/components/form"
	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	"github.com/LekcRg/GophKeeper/internal/client/msgs"
	"github.com/LekcRg/GophKeeper/internal/client/router"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectTypeModel struct {
	form *form.Form
	help *help.SelectAuth
}

func NewSelectType() tea.Model {
	buttons := []components.Button{
		{
			Label:   "Password",
			Name:    string(router.CreateVaultPassword),
			Focused: true,
		},
		{
			Label: "Note",
			Name:  string(router.CreateVaultNote),
		},
		{
			Label: "Card",
			Name:  string(router.CreateVaultCard),
		},
		{
			Label: "Binary",
			Name:  string(router.CreateVaultBinary),
		},
	}

	h := help.NewSelectAuth()

	return &SelectTypeModel{
		form: form.NewForm(form.FormOpts{
			Buttons: buttons,
			Up:      h.Keys.Up,
			Down:    h.Keys.Down,
		}),
		help: h,
	}
}

func (m *SelectTypeModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *SelectTypeModel) handleSubmit(msg msgs.FormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		return msgs.SelectTypeMsg{
			Selected: msg.ButtonName,
		}
	}
}

func (m *SelectTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typeMsg := msg.(type) {
	case msgs.FormSubmitMsg:
		return m, m.handleSubmit(typeMsg)
	default:
		var newMsg tea.Cmd
		m.form, newMsg = m.form.Update(msg)

		return m, newMsg
	}
}

func (m *SelectTypeModel) View() string {
	var b strings.Builder

	b.WriteString(m.form.View())
	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
