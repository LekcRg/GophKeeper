package components

import tea "github.com/charmbracelet/bubbletea"

type FormComopnents interface {
	View() string
	Update() tea.Cmd
	Focus() tea.Cmd
	Blur()
}
