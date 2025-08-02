package components

import (
	"github.com/LekcRg/GophKeeper/internal/client/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Field struct {
	Label string
	Value string
}

func (m *Field) Init() tea.Cmd {
	return nil
}

func (m *Field) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Field) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		styles.FieldLabel.Render(m.Label),
		styles.NoStyle.Render(m.Value),
	)
}
