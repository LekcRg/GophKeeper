package components

import (
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/components/help"
	tea "github.com/charmbracelet/bubbletea"
)

type DetailModel struct {
	help   *help.Auth
	fields []Field
}

func NewDetail(field []Field) tea.Model {
	return &DetailModel{
		help:   help.NewAuth(),
		fields: field,
	}
}

func (m *DetailModel) Init() tea.Cmd {
	return nil
}

func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *DetailModel) View() string {
	var b strings.Builder

	b.WriteRune('\n')

	for i := range m.fields {
		field := &m.fields[i]
		b.WriteString(field.View())
		b.WriteRune('\n')
	}

	b.WriteRune('\n')
	b.WriteString(m.help.View())

	return b.String()
}
